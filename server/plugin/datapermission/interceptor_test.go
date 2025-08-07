package datapermission

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/utils"
)

// 测试用户模型
type TestUser struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Username  string    `json:"username" gorm:"column:username"`
	Email     string    `json:"email" gorm:"column:email"`
	Phone     string    `json:"phone" gorm:"column:phone"`
	DeptID    uint      `json:"dept_id" gorm:"column:dept_id"`
	CreatedBy uint      `json:"created_by" gorm:"column:created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TestUser) TableName() string {
	return "test_users"
}

// 测试部门模型
type TestDept struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Name string `json:"name" gorm:"column:name"`
}

func (TestDept) TableName() string {
	return "test_depts"
}

// 测试权限配置模型
type TestDataPermission struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	AuthorityID string `json:"authority_id" gorm:"column:authority_id"`
	Table       string `json:"table_name" gorm:"column:table_name"`
	DataScope   int    `json:"data_scope" gorm:"column:data_scope"`
	CustomSQL   string `json:"custom_sql" gorm:"column:custom_sql"`
	FieldConfig string `json:"field_config" gorm:"column:field_config"`
}

func (TestDataPermission) TableName() string {
	return "test_data_permissions"
}

// 测试用户信息
type TestUserInfo struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	AuthorityID string `json:"authority_id"`
	DeptID      uint   `json:"dept_id"`
}

// 测试设置
type InterceptorTestSuite struct {
	db     *gorm.DB
	gin    *gin.Engine
	helper *utils.EnhancedPermissionHelper
}

// 设置测试环境
func setupTestSuite(t *testing.T) *InterceptorTestSuite {
	// 创建内存数据库
	// 使用内存MySQL数据库进行测试
	dsn := "root:@tcp(localhost:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		// 如果MySQL连接失败，跳过测试
		t.Skip("MySQL database not available for testing")
		return nil
	}
	assert.NoError(t, err)

	// 自动迁移表结构
	err = db.AutoMigrate(&TestUser{}, &TestDept{}, &TestDataPermission{})
	assert.NoError(t, err)

	// 初始化测试数据
	setupTestData(t, db)

	// 设置Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 创建增强权限助手
	helper := &utils.EnhancedPermissionHelper{
		PermissionHelper: utils.PermissionHelperApp,
	}

	return &InterceptorTestSuite{
		db:     db,
		gin:    router,
		helper: helper,
	}
}

// 设置测试数据
func setupTestData(t *testing.T, db *gorm.DB) {
	// 创建测试部门
	depts := []TestDept{
		{ID: 1, Name: "技术部"},
		{ID: 2, Name: "市场部"},
		{ID: 3, Name: "人事部"},
	}
	for _, dept := range depts {
		err := db.Create(&dept).Error
		assert.NoError(t, err)
	}

	// 创建测试用户
	users := []TestUser{
		{ID: 1, Username: "admin", Email: "admin@test.com", Phone: "13800000001", DeptID: 1, CreatedBy: 1},
		{ID: 2, Username: "user1", Email: "user1@test.com", Phone: "13800000002", DeptID: 1, CreatedBy: 1},
		{ID: 3, Username: "user2", Email: "user2@test.com", Phone: "13800000003", DeptID: 2, CreatedBy: 2},
		{ID: 4, Username: "user3", Email: "user3@test.com", Phone: "13800000004", DeptID: 3, CreatedBy: 3},
	}
	for _, user := range users {
		err := db.Create(&user).Error
		assert.NoError(t, err)
	}

	// 创建测试权限配置
	fieldConfig := map[string]interface{}{
		"username": map[string]bool{"view": true, "edit": true, "export": true, "query": true},
		"email":    map[string]bool{"view": true, "edit": false, "export": true, "query": true},
		"phone":    map[string]bool{"view": false, "edit": false, "export": false, "query": false},
	}
	fieldConfigJSON, _ := json.Marshal(fieldConfig)

	permissions := []TestDataPermission{
		{
			ID:          1,
			AuthorityID: "admin",
			Table:       "test_users",
			DataScope:   1, // 全部数据
			FieldConfig: string(fieldConfigJSON),
		},
		{
			ID:          2,
			AuthorityID: "user",
			Table:       "test_users",
			DataScope:   2, // 本部门数据
			FieldConfig: string(fieldConfigJSON),
		},
		{
			ID:          3,
			AuthorityID: "guest",
			Table:       "test_users",
			DataScope:   3, // 仅本人数据
			FieldConfig: string(fieldConfigJSON),
		},
	}
	for _, permission := range permissions {
		err := db.Create(&permission).Error
		assert.NoError(t, err)
	}
}

// 创建测试上下文
func createTestContext(userInfo TestUserInfo) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// 设置用户信息到上下文
	c.Set("userInfo", userInfo)

	return c, w
}

// 测试拦截器中间件
func TestDataPermissionInterceptorMiddleware(t *testing.T) {
	suite := setupTestSuite(t)

	// 设置路由
	suite.gin.Use(utils.DataPermissionInterceptorMiddleware())
	suite.gin.GET("/test", func(c *gin.Context) {
		// 获取拦截器数据库实例
		db := utils.GetInterceptorDB(c)
		assert.NotNil(t, db)

		c.JSON(200, gin.H{"status": "ok"})
	})

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 设置用户信息
	userInfo := TestUserInfo{
		ID:          1,
		Username:    "admin",
		AuthorityID: "admin",
		DeptID:      1,
	}
	req.Header.Set("User-Info", encodeUserInfo(userInfo))

	suite.gin.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// 测试查询权限
func TestQueryWithPermission(t *testing.T) {
	suite := setupTestSuite(t)

	// 测试管理员权限（全部数据）
	c, _ := createTestContext(TestUserInfo{
		ID:          1,
		Username:    "admin",
		AuthorityID: "admin",
		DeptID:      1,
	})

	var users []TestUser
	err := suite.helper.QueryWithPermission(c, "test_users", &users)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(users)) // 应该能看到所有用户

	// 测试普通用户权限（本部门数据）
	c, _ = createTestContext(TestUserInfo{
		ID:          2,
		Username:    "user1",
		AuthorityID: "user",
		DeptID:      1,
	})

	users = []TestUser{}
	err = suite.helper.QueryWithPermission(c, "test_users", &users)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users)) // 应该只能看到本部门用户

	// 测试访客权限（仅本人数据）
	c, _ = createTestContext(TestUserInfo{
		ID:          3,
		Username:    "user2",
		AuthorityID: "guest",
		DeptID:      2,
	})

	users = []TestUser{}
	err = suite.helper.QueryWithPermission(c, "test_users", &users)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users)) // 应该只能看到自己
	assert.Equal(t, uint(3), users[0].ID)
}

// 测试创建权限
func TestCreateWithPermission(t *testing.T) {
	suite := setupTestSuite(t)

	c, _ := createTestContext(TestUserInfo{
		ID:          1,
		Username:    "admin",
		AuthorityID: "admin",
		DeptID:      1,
	})

	newUser := TestUser{
		Username: "newuser",
		Email:    "newuser@test.com",
		Phone:    "13800000005",
	}

	err := suite.helper.CreateWithPermission(c, "test_users", &newUser)
	assert.NoError(t, err)

	// 验证自动填充的字段
	assert.Equal(t, uint(1), newUser.CreatedBy) // 应该自动设置创建者
	assert.Equal(t, uint(1), newUser.DeptID)    // 应该自动设置部门
}

// 测试更新权限
func TestUpdateWithPermission(t *testing.T) {
	suite := setupTestSuite(t)

	c, _ := createTestContext(TestUserInfo{
		ID:          1,
		Username:    "admin",
		AuthorityID: "admin",
		DeptID:      1,
	})

	updateData := TestUser{
		Username: "updated_user",
		Email:    "updated@test.com",
	}

	err := suite.helper.UpdateWithPermission(c, "test_users", &updateData, "id = ?", 2)
	assert.NoError(t, err)

	// 验证更新结果
	var user TestUser
	err = suite.db.Where("id = ?", 2).First(&user).Error
	assert.NoError(t, err)
	assert.Equal(t, "updated_user", user.Username)
}

// 测试删除权限
func TestDeleteWithPermission(t *testing.T) {
	suite := setupTestSuite(t)

	c, _ := createTestContext(TestUserInfo{
		ID:          1,
		Username:    "admin",
		AuthorityID: "admin",
		DeptID:      1,
	})

	err := suite.helper.DeleteWithPermission(c, "test_users", "id = ?", 4)
	assert.NoError(t, err)

	// 验证删除结果
	var count int64
	suite.db.Model(&TestUser{}).Where("id = ?", 4).Count(&count)
	assert.Equal(t, int64(0), count)
}

// 测试字段过滤
func TestFilterResponseFields(t *testing.T) {
	suite := setupTestSuite(t)

	c, _ := createTestContext(TestUserInfo{
		ID:          2,
		Username:    "user1",
		AuthorityID: "user",
		DeptID:      1,
	})

	users := []TestUser{
		{ID: 1, Username: "admin", Email: "admin@test.com", Phone: "13800000001"},
		{ID: 2, Username: "user1", Email: "user1@test.com", Phone: "13800000002"},
	}

	filteredUsers := suite.helper.FilterResponseFields(c, users, "test_users")
	assert.NotNil(t, filteredUsers)

	// 根据字段配置，phone字段应该被过滤掉
	// 这里需要根据实际的字段过滤逻辑进行验证
}

// 测试跳过权限
func TestSkipDataPermission(t *testing.T) {
	suite := setupTestSuite(t)

	c, _ := createTestContext(TestUserInfo{
		ID:          3,
		Username:    "user2",
		AuthorityID: "guest",
		DeptID:      2,
	})

	// 正常情况下只能看到自己的数据
	var users []TestUser
	err := suite.helper.QueryWithPermission(c, "test_users", &users)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))

	// 跳过权限后应该能看到所有数据
	db := suite.helper.GetDBSkipPermission(c)
	users = []TestUser{}
	err = db.Find(&users).Error
	assert.NoError(t, err)
	assert.Equal(t, 4, len(users))
}

// 测试系统操作
func TestSystemOperation(t *testing.T) {
	suite := setupTestSuite(t)

	c, _ := createTestContext(TestUserInfo{
		ID:          3,
		Username:    "user2",
		AuthorityID: "guest",
		DeptID:      2,
	})

	// 系统操作不受权限限制
	db := suite.helper.GetDBSystemOperation(c)
	var users []TestUser
	err := db.Find(&users).Error
	assert.NoError(t, err)
	assert.Equal(t, 4, len(users)) // 系统操作应该能看到所有数据
}

// 测试计数功能
func TestCountWithPermission(t *testing.T) {
	suite := setupTestSuite(t)

	// 管理员应该能统计所有用户
	c, _ := createTestContext(TestUserInfo{
		ID:          1,
		Username:    "admin",
		AuthorityID: "admin",
		DeptID:      1,
	})

	count, err := suite.helper.CountWithPermission(c, "test_users")
	assert.NoError(t, err)
	assert.Equal(t, int64(4), count)

	// 普通用户只能统计本部门用户
	c, _ = createTestContext(TestUserInfo{
		ID:          2,
		Username:    "user1",
		AuthorityID: "user",
		DeptID:      1,
	})

	count, err = suite.helper.CountWithPermission(c, "test_users")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// 测试存在性检查
func TestExistsWithPermission(t *testing.T) {
	suite := setupTestSuite(t)

	c, _ := createTestContext(TestUserInfo{
		ID:          3,
		Username:    "user2",
		AuthorityID: "guest",
		DeptID:      2,
	})

	// 检查自己的记录应该存在
	exists, err := suite.helper.ExistsWithPermission(c, "test_users", "id = ?", 3)
	assert.NoError(t, err)
	assert.True(t, exists)

	// 检查其他用户的记录应该不存在（因为权限限制）
	exists, err = suite.helper.ExistsWithPermission(c, "test_users", "id = ?", 1)
	assert.NoError(t, err)
	assert.False(t, exists)
}

// 辅助函数：编码用户信息
func encodeUserInfo(userInfo TestUserInfo) string {
	data, _ := json.Marshal(userInfo)
	return string(data)
}

// 基准测试：对比传统模式和拦截器模式的性能
func BenchmarkTraditionalMode(b *testing.B) {
	// 这里可以实现传统模式的基准测试
}

func BenchmarkInterceptorMode(b *testing.B) {
	// 这里可以实现拦截器模式的基准测试
}

// 集成测试：测试完整的HTTP请求流程
func TestFullHTTPFlow(t *testing.T) {
	suite := setupTestSuite(t)

	// 设置路由
	suite.gin.Use(utils.DataPermissionInterceptorMiddleware())
	suite.gin.GET("/users", func(c *gin.Context) {
		var users []TestUser
		err := suite.helper.QueryWithPermission(c, "test_users", &users)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// 过滤响应字段
		filteredUsers := suite.helper.FilterResponseFields(c, users, "test_users")
		c.JSON(200, gin.H{"data": filteredUsers})
	})

	// 测试不同权限用户的请求
	testCases := []struct {
		name          string
		userInfo      TestUserInfo
		expectedCode  int
		expectedCount int
	}{
		{
			name: "管理员用户",
			userInfo: TestUserInfo{
				ID:          1,
				Username:    "admin",
				AuthorityID: "admin",
				DeptID:      1,
			},
			expectedCode:  200,
			expectedCount: 4,
		},
		{
			name: "普通用户",
			userInfo: TestUserInfo{
				ID:          2,
				Username:    "user1",
				AuthorityID: "user",
				DeptID:      1,
			},
			expectedCode:  200,
			expectedCount: 2,
		},
		{
			name: "访客用户",
			userInfo: TestUserInfo{
				ID:          3,
				Username:    "user2",
				AuthorityID: "guest",
				DeptID:      2,
			},
			expectedCode:  200,
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/users", nil)
			req.Header.Set("User-Info", encodeUserInfo(tc.userInfo))

			w := httptest.NewRecorder()
			suite.gin.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			if w.Code == 200 {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				data, ok := response["data"].([]interface{})
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCount, len(data))
			}
		})
	}
}
