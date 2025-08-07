package utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EnhancedPermissionHelper 增强的权限助手
// 提供更简洁的API，配合拦截器使用
type EnhancedPermissionHelper struct {
	*PermissionHelper
}

var EnhancedPermissionHelperApp = &EnhancedPermissionHelper{
	PermissionHelper: PermissionHelperApp,
}

// GetDB 获取带有数据权限拦截器的数据库实例
// 使用示例：
// db := utils.EnhancedPermissionHelperApp.GetDB(c)
// var users []User
// db.Find(&users) // 自动应用数据权限
func (h *EnhancedPermissionHelper) GetDB(c *gin.Context) *gorm.DB {
	return GetInterceptorDB(c)
}

// GetDBWithTable 获取带有数据权限拦截器的数据库实例，并指定表名
// 使用示例：
// db := utils.EnhancedPermissionHelperApp.GetDBWithTable(c, "users")
// var users []User
// db.Find(&users) // 自动应用数据权限
func (h *EnhancedPermissionHelper) GetDBWithTable(c *gin.Context, tableName string) *gorm.DB {
	db := GetInterceptorDB(c)
	return SetTableName(db, tableName)
}

// GetDBSkipPermission 获取跳过数据权限检查的数据库实例
// 使用示例：
// db := utils.EnhancedPermissionHelperApp.GetDBSkipPermission(c)
// var users []User
// db.Find(&users) // 跳过数据权限检查
func (h *EnhancedPermissionHelper) GetDBSkipPermission(c *gin.Context) *gorm.DB {
	db := GetInterceptorDB(c)
	return SkipDataPermission(db)
}

// GetDBSystemOperation 获取标记为系统操作的数据库实例
// 使用示例：
// db := utils.EnhancedPermissionHelperApp.GetDBSystemOperation(c)
// var users []User
// db.Find(&users) // 系统操作，跳过权限检查
func (h *EnhancedPermissionHelper) GetDBSystemOperation(c *gin.Context) *gorm.DB {
	db := GetInterceptorDB(c)
	return SetSystemOperation(db)
}

// QueryWithPermission 带权限的查询
// 使用示例：
// var users []User
// err := utils.EnhancedPermissionHelperApp.QueryWithPermission(c, "users", &users)
func (h *EnhancedPermissionHelper) QueryWithPermission(c *gin.Context, tableName string, dest interface{}) error {
	db := h.GetDBWithTable(c, tableName)
	return db.Find(dest).Error
}

// QueryOneWithPermission 带权限的单条查询
// 使用示例：
// var user User
// err := utils.EnhancedPermissionHelperApp.QueryOneWithPermission(c, "users", &user, "id = ?", 1)
func (h *EnhancedPermissionHelper) QueryOneWithPermission(c *gin.Context, tableName string, dest interface{}, query interface{}, args ...interface{}) error {
	db := h.GetDBWithTable(c, tableName)
	return db.Where(query, args...).First(dest).Error
}

// CreateWithPermission 带权限的创建
// 使用示例：
// user := User{Name: "test"}
// err := utils.EnhancedPermissionHelperApp.CreateWithPermission(c, "users", &user)
func (h *EnhancedPermissionHelper) CreateWithPermission(c *gin.Context, tableName string, value interface{}) error {
	db := h.GetDBWithTable(c, tableName)
	return db.Create(value).Error
}

// UpdateWithPermission 带权限的更新
// 使用示例：
// err := utils.EnhancedPermissionHelperApp.UpdateWithPermission(c, "users", map[string]interface{}{"name": "new name"}, "id = ?", 1)
func (h *EnhancedPermissionHelper) UpdateWithPermission(c *gin.Context, tableName string, values interface{}, query interface{}, args ...interface{}) error {
	db := h.GetDBWithTable(c, tableName)
	return db.Where(query, args...).Updates(values).Error
}

// DeleteWithPermission 带权限的删除
// 使用示例：
// err := utils.EnhancedPermissionHelperApp.DeleteWithPermission(c, "users", "id = ?", 1)
func (h *EnhancedPermissionHelper) DeleteWithPermission(c *gin.Context, tableName string, query interface{}, args ...interface{}) error {
	db := h.GetDBWithTable(c, tableName)
	return db.Where(query, args...).Delete(nil).Error
}

// CountWithPermission 带权限的计数
// 使用示例：
// count, err := utils.EnhancedPermissionHelperApp.CountWithPermission(c, "users")
func (h *EnhancedPermissionHelper) CountWithPermission(c *gin.Context, tableName string) (int64, error) {
	db := h.GetDBWithTable(c, tableName)
	var count int64
	err := db.Count(&count).Error
	return count, err
}

// ExistsWithPermission 带权限的存在性检查
// 使用示例：
// exists, err := utils.EnhancedPermissionHelperApp.ExistsWithPermission(c, "users", "id = ?", 1)
func (h *EnhancedPermissionHelper) ExistsWithPermission(c *gin.Context, tableName string, query interface{}, args ...interface{}) (bool, error) {
	count, err := h.CountWithPermissionWhere(c, tableName, query, args...)
	return count > 0, err
}

// CountWithPermissionWhere 带权限和条件的计数
// 使用示例：
// count, err := utils.EnhancedPermissionHelperApp.CountWithPermissionWhere(c, "users", "status = ?", "active")
func (h *EnhancedPermissionHelper) CountWithPermissionWhere(c *gin.Context, tableName string, query interface{}, args ...interface{}) (int64, error) {
	db := h.GetDBWithTable(c, tableName)
	var count int64
	err := db.Where(query, args...).Count(&count).Error
	return count, err
}

// FilterResponseFields 过滤响应字段
// 使用示例：
// filteredData := utils.EnhancedPermissionHelperApp.FilterResponseFields(c, users, "users")
func (h *EnhancedPermissionHelper) FilterResponseFields(c *gin.Context, data interface{}, tableName string) interface{} {
	return h.FilterFieldsByPermission(c, data, tableName, "view")
}

// FilterEditableFields 过滤可编辑字段
// 使用示例：
// editableData := utils.EnhancedPermissionHelperApp.FilterEditableFields(c, updateData, "users")
func (h *EnhancedPermissionHelper) FilterEditableFields(c *gin.Context, data map[string]interface{}, tableName string) map[string]interface{} {
	result := make(map[string]interface{})
	for field, value := range data {
		if h.CheckFieldPermission(c, tableName, field, "edit") {
			result[field] = value
		}
	}
	return result
}

// FilterExportableFields 过滤可导出字段
// 使用示例：
// exportableFields := utils.EnhancedPermissionHelperApp.FilterExportableFields(c, allFields, "users")
func (h *EnhancedPermissionHelper) FilterExportableFields(c *gin.Context, fields []string, tableName string) []string {
	var result []string
	for _, field := range fields {
		if h.CheckFieldPermission(c, tableName, field, "export") {
			result = append(result, field)
		}
	}
	return result
}

// FilterQueryableFields 过滤可查询字段
// 使用示例：
// queryableFields := utils.EnhancedPermissionHelperApp.FilterQueryableFields(c, searchFields, "users")
func (h *EnhancedPermissionHelper) FilterQueryableFields(c *gin.Context, fields []string, tableName string) []string {
	var result []string
	for _, field := range fields {
		if h.CheckFieldPermission(c, tableName, field, "query") {
			result = append(result, field)
		}
	}
	return result
}

// BuildSelectClause 构建SELECT子句
// 使用示例：
// selectFields := utils.EnhancedPermissionHelperApp.BuildSelectClause(c, "users", []string{"id", "name", "email"})
// db.Select(selectFields).Find(&users)
func (h *EnhancedPermissionHelper) BuildSelectClause(c *gin.Context, tableName string, allFields []string) []string {
	return h.BuildSelectFields(c, tableName, allFields)
}

// ValidateFieldAccess 验证字段访问权限
// 使用示例：
//
//	if utils.EnhancedPermissionHelperApp.ValidateFieldAccess(c, "users", "salary", "view") {
//	    // 显示薪资字段
//	}
func (h *EnhancedPermissionHelper) ValidateFieldAccess(c *gin.Context, tableName, fieldName, operation string) bool {
	return h.CheckFieldPermission(c, tableName, fieldName, operation)
}

// GetUserPermissionInfo 获取用户权限信息
// 使用示例：
// info := utils.EnhancedPermissionHelperApp.GetUserPermissionInfo(c, "users")
func (h *EnhancedPermissionHelper) GetUserPermissionInfo(c *gin.Context, tableName string) map[string]interface{} {
	dataScope, _ := h.GetUserDataScope(c, tableName)
	visibleFields, _ := h.GetVisibleFields(c, tableName)
	editableFields, _ := h.GetEditableFields(c, tableName)

	return map[string]interface{}{
		"dataScope":      dataScope,
		"visibleFields":  visibleFields,
		"editableFields": editableFields,
		"tableName":      tableName,
	}
}
