package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	system "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DataPermissionMiddleware 数据权限中间件
type DataPermissionMiddleware struct{}

// ApplyDataPermission 应用数据权限到GORM查询
func (m *DataPermissionMiddleware) ApplyDataPermission(db *gorm.DB, c *gin.Context, tableName string) *gorm.DB {
	// 从上下文获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		return db
	}

	// 获取用户的所有角色ID
	authorityIds, err := m.getUserAuthorityIds(userID.(uint))
	if err != nil {
		global.GVA_LOG.Error("获取用户角色失败", zap.Error(err))
		return db
	}

	// 如果用户没有角色，返回原始查询
	if len(authorityIds) == 0 {
		return db
	}

	// 获取所有角色的数据权限配置
	allConditions := make([]string, 0)
	for _, authorityId := range authorityIds {
		condition, err := m.getDataPermissionCondition(tableName, authorityId, userID.(uint))
		if err != nil {
			global.GVA_LOG.Error("获取数据权限配置失败", zap.Error(err))
			continue
		}
		if condition == "1=1" {
			allConditions = []string{"1=1"}
			break
		}
		if condition != "" {
			allConditions = append(allConditions, condition)
		}
	}

	// 如果有多个条件，使用OR连接
	if len(allConditions) > 0 {
		finalCondition := strings.Join(allConditions, " OR ")
		// 检查是否已经应用了相同的条件，避免重复
		if db.Statement != nil && db.Statement.SQL.String() != "" {
			// 如果SQL已经包含了相同的条件，跳过
			if strings.Contains(db.Statement.SQL.String(), finalCondition) {
				return db
			}
		}
		db = db.Where("(" + finalCondition + ")")
	}

	return db
}

// ApplyFieldPermission 应用字段权限到查询结果
func (m *DataPermissionMiddleware) ApplyFieldPermission(data interface{}, c *gin.Context, tableName string, operation string) interface{} {
	// 从上下文获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		return data
	}

	// 获取用户的所有角色ID
	authorityIds, err := m.getUserAuthorityIds(userID.(uint))
	if err != nil {
		global.GVA_LOG.Error("获取用户角色失败", zap.Error(err))
		return data
	}

	// 如果用户没有角色，返回原始数据
	if len(authorityIds) == 0 {
		return data
	}

	// 获取所有角色的字段权限配置
	allFieldPermissions := make(map[string]model.RoleFieldPermission)
	for _, authorityId := range authorityIds {
		fieldPermissions, err := m.getFieldPermissions(tableName, authorityId)
		if err != nil {
			global.GVA_LOG.Error("获取字段权限配置失败", zap.Error(err))
			continue
		}
		// 合并字段权限（取最宽松的权限）
		for fieldName, perm := range fieldPermissions {
			allFieldPermissions[fieldName] = perm
		}
	}

	// 根据操作类型过滤字段
	return m.filterFieldsByPermission(data, allFieldPermissions, operation)
}

// getUserAuthorityIds 获取用户的所有角色ID
func (m *DataPermissionMiddleware) getUserAuthorityIds(userID uint) ([]uint, error) {
	var user system.SysUser
	err := global.GVA_DB.Preload("Authorities").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	// 收集所有角色ID
	authorityIds := make([]uint, 0)

	// 添加主角色ID
	if user.AuthorityId != 0 {
		authorityIds = append(authorityIds, user.AuthorityId)
	}

	// 添加多角色ID
	for _, authority := range user.Authorities {
		// 避免重复添加
		alreadyExists := false
		for _, existingId := range authorityIds {
			if existingId == authority.AuthorityId {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			authorityIds = append(authorityIds, authority.AuthorityId)
		}
	}

	return authorityIds, nil
}

// getDataPermissionCondition 获取数据权限条件
func (m *DataPermissionMiddleware) getDataPermissionCondition(tableName string, authorityID uint, userID uint) (string, error) {
	// 获取受控表信息
	var controlledTable model.ControlledTable
	if err := SkipDataPermission(global.GVA_DB).Where("table_name = ? AND enabled = ?", tableName, true).First(&controlledTable).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果表不受控制，返回无限制条件
			return "1=1", nil
		}
		return "", err
	}

	// 获取角色数据权限配置
	var roleDataPermission model.RoleDataPermission
	if err := SkipDataPermission(global.GVA_DB).Where("authority_id = ? AND controlled_table_id = ? AND enabled = ?", authorityID, controlledTable.ID, true).First(&roleDataPermission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有配置权限，使用默认的自己数据权限
			return m.generateSQLCondition("self", "", controlledTable, userID), nil
		}
		return "", err
	}

	// 生成SQL条件
	return m.generateSQLCondition(roleDataPermission.DataScope, roleDataPermission.CustomCondition, controlledTable, userID), nil
}

// getFieldPermissions 获取字段权限配置
func (m *DataPermissionMiddleware) getFieldPermissions(tableName string, authorityID uint) (map[string]model.RoleFieldPermission, error) {
	// 获取受控表信息
	var controlledTable model.ControlledTable
	if err := SkipDataPermission(global.GVA_DB).Where("table_name = ? AND enabled = ?", tableName, true).First(&controlledTable).Error; err != nil {
		return nil, err
	}

	// 获取字段权限配置
	var fieldPermissions []model.RoleFieldPermission
	if err := SkipDataPermission(global.GVA_DB).Where("authority_id = ? AND controlled_table_id = ? AND enabled = ?", authorityID, controlledTable.ID, true).Find(&fieldPermissions).Error; err != nil {
		return nil, err
	}

	// 转换为map
	permissionMap := make(map[string]model.RoleFieldPermission)
	for _, perm := range fieldPermissions {
		permissionMap[perm.FieldName] = perm
	}

	return permissionMap, nil
}

// generateSQLCondition 生成SQL条件
func (m *DataPermissionMiddleware) generateSQLCondition(dataScope, customCondition string, table model.ControlledTable, userID uint) string {
	switch dataScope {
	case "self":
		if table.UserField != "" {
			return fmt.Sprintf("%s = %d", table.UserField, userID)
		}
		return ""
	case "dept":
		if table.DeptField != "" {
			// 获取用户部门ID，可以和组织机构插件结合
			return fmt.Sprintf("%s = (SELECT %s FROM sys_users WHERE id = %d)", table.DeptField, table.DeptField, userID)
		}
		return ""
	case "dept_and_child":
		if table.DeptField != "" {
			// 获取用户部门及子部门，可以和组织机构插件结合
			return fmt.Sprintf("%s IN (SELECT id FROM sys_dept WHERE find_in_set(id, get_child_dept((SELECT dept_id FROM sys_users WHERE id = %d))))", table.DeptField, userID)
		}
		return ""
	case "all":
		return "1=1"
	case "custom":
		if customCondition != "" {
			// 替换自定义条件中的占位符
			condition := strings.ReplaceAll(customCondition, "${USER_ID}", fmt.Sprintf("%d", userID))
			return condition
		}
		return ""
	default:
		return ""
	}
}

// filterFieldsByPermission 根据权限过滤字段
func (m *DataPermissionMiddleware) filterFieldsByPermission(data interface{}, permissions map[string]model.RoleFieldPermission, operation string) interface{} {
	// 这里需要根据具体的数据结构来实现字段过滤
	// 由于Go的反射比较复杂，这里提供一个简化的实现思路
	// 实际使用时可能需要根据具体的业务模型来调整

	// 如果没有权限配置，返回原数据
	if len(permissions) == 0 {
		return data
	}

	// TODO: 实现具体的字段过滤逻辑
	// 可以使用反射来动态过滤字段，或者在业务层面进行处理

	return data
}

// CheckFieldPermission 检查字段权限
func (m *DataPermissionMiddleware) CheckFieldPermission(c *gin.Context, tableName, fieldName, operation string) *bool {
	// 从上下文获取用户信息
	authorityID, exists := c.Get("authorityId")
	if !exists {
		// 修改后
		var True = true
		return &True // 如果没有权限信息，默认允许
	}

	// 获取字段权限配置
	fieldPermissions, err := m.getFieldPermissions(tableName, authorityID.(uint))
	if err != nil {
		var True = true
		return &True // 如果获取权限失败，默认允许
	}

	// 检查字段权限
	perm, exists := fieldPermissions[fieldName]
	if !exists {
		var True = true
		return &True // 如果没有配置该字段权限，默认允许
	}

	// 根据操作类型检查权限
	switch operation {
	case "view":
		return perm.Visibility
	case "edit":
		return perm.EditPermission
	case "export":
		return perm.Exportable
	case "query":
		return perm.Queryable
	default:
		var True = true
		return &True
	}
}

// DataPermissionGinMiddleware Gin中间件，用于自动应用数据权限（传统方式）
// 注意：推荐使用 DataPermissionInterceptorMiddleware() 获得更好的自动化体验
func DataPermissionGinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 将数据权限中间件实例添加到上下文
		c.Set("dataPermissionMiddleware", &DataPermissionMiddleware{})
		c.Next()
	}
}

// DataPermissionAutoMiddleware 自动数据权限中间件（推荐使用）
// 这是 DataPermissionInterceptorMiddleware 的别名，提供更直观的命名
func DataPermissionAutoMiddleware() gin.HandlerFunc {
	return DataPermissionInterceptorMiddleware()
}

// GetDataPermissionMiddleware 从Gin上下文获取数据权限中间件
func GetDataPermissionMiddleware(c *gin.Context) *DataPermissionMiddleware {
	if middleware, exists := c.Get("dataPermissionMiddleware"); exists {
		return middleware.(*DataPermissionMiddleware)
	}
	return &DataPermissionMiddleware{}
}
