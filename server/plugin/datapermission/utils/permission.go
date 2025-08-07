package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PermissionHelper 权限助手
type PermissionHelper struct{}

var PermissionHelperApp = &PermissionHelper{}

// ApplyDataPermissionToQuery 为GORM查询应用数据权限
// 使用示例：
// db := global.GVA_DB.Model(&User{})
// db = PermissionHelperApp.ApplyDataPermissionToQuery(db, c, "users")
// var users []User
// db.Find(&users)
func (p *PermissionHelper) ApplyDataPermissionToQuery(db *gorm.DB, c *gin.Context, tableName string) *gorm.DB {
	// 检查是否已经使用了拦截器模式
	if _, exists := c.Get("interceptor_db"); exists {
		// 如果已经使用了拦截器，记录警告并直接返回，避免重复应用权限
		global.GVA_LOG.Warn("检测到同时使用拦截器和手动权限应用，建议统一使用拦截器模式。请移除手动调用 ApplyDataPermissionToQuery")
		return db
	}
	middleware := GetDataPermissionMiddleware(c)
	return middleware.ApplyDataPermission(db, c, tableName)
}

// CheckTablePermission 检查表级权限
func (p *PermissionHelper) CheckTablePermission(c *gin.Context, tableName string) bool {
	// 从上下文获取用户信息
	authorityID, exists := c.Get("authorityId")
	if !exists {
		return true // 如果没有权限信息，默认允许
	}

	// 检查表是否受控制
	var controlledTable model.ControlledTable
	if err := global.GVA_DB.Where("table_name = ? AND enabled = ?", tableName, true).First(&controlledTable).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true // 如果表不受控制，默认允许
		}
		return false
	}

	// 检查角色是否有权限配置
	var count int64
	global.GVA_DB.Model(&model.RoleDataPermission{}).Where("authority_id = ? AND controlled_table_id = ? AND enabled = ?", authorityID, controlledTable.ID, true).Count(&count)
	return count > 0
}

// CheckFieldPermission 检查字段权限
func (p *PermissionHelper) CheckFieldPermission(c *gin.Context, tableName, fieldName, operation string) bool {
	middleware := GetDataPermissionMiddleware(c)
	return *middleware.CheckFieldPermission(c, tableName, fieldName, operation)
}

// GetUserDataScope 获取用户在指定表的数据范围
func (p *PermissionHelper) GetUserDataScope(c *gin.Context, tableName string) (string, error) {
	// 从上下文获取用户信息
	authorityID, exists := c.Get("authorityId")
	if !exists {
		return "self", errors.New("无法获取用户权限信息")
	}

	// 获取受控表信息
	var controlledTable model.ControlledTable
	if err := global.GVA_DB.Where("table_name = ? AND enabled = ?", tableName, true).First(&controlledTable).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "all", nil // 如果表不受控制，返回全部权限
		}
		return "self", err
	}

	// 获取角色数据权限配置
	var roleDataPermission model.RoleDataPermission
	if err := global.GVA_DB.Where("authority_id = ? AND controlled_table_id = ? AND enabled = ?", authorityID, controlledTable.ID, true).First(&roleDataPermission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "self", nil // 如果没有配置权限，默认返回自己数据权限
		}
		return "self", err
	}

	return roleDataPermission.DataScope, nil
}

// GetDataPermissionSQL 获取数据权限SQL条件
func (p *PermissionHelper) GetDataPermissionSQL(c *gin.Context, tableName string) (string, error) {
	// 从上下文获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		return "1=1", errors.New("无法获取用户ID")
	}

	authorityIds, exists := c.Get("authorityIds")
	if !exists {
		return "1=1", errors.New("无法获取用户权限信息")
	}

	middleware := &DataPermissionMiddleware{}
	for authorityID, _ := range authorityIds.([]uint) {
		condition, err := middleware.getDataPermissionCondition(tableName, uint(authorityID), userID.(uint))
		if err != nil {
			return "1=1", err
		}

		return condition, nil
	}
	return "1=1", nil
}

// FilterFieldsByPermission 根据权限过滤结构体字段
func (p *PermissionHelper) FilterFieldsByPermission(c *gin.Context, data interface{}, tableName, operation string) interface{} {
	// 从上下文获取用户信息
	authorityID, exists := c.Get("authorityId")
	if !exists {
		return data // 如果没有权限信息，返回原数据
	}

	// 获取字段权限配置
	middleware := &DataPermissionMiddleware{}
	fieldPermissions, err := middleware.getFieldPermissions(tableName, authorityID.(uint))
	if err != nil || len(fieldPermissions) == 0 {
		return data // 如果获取权限失败或没有权限配置，返回原数据
	}

	// 使用反射过滤字段
	return p.filterStructFields(data, fieldPermissions, operation)
}

// filterStructFields 使用反射过滤结构体字段
func (p *PermissionHelper) filterStructFields(data interface{}, permissions map[string]model.RoleFieldPermission, operation string) interface{} {
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	// 如果是指针，获取其指向的值
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// 如果是切片，递归处理每个元素
	if val.Kind() == reflect.Slice {
		result := reflect.MakeSlice(typ, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			filtered := p.filterStructFields(val.Index(i).Interface(), permissions, operation)
			result.Index(i).Set(reflect.ValueOf(filtered))
		}
		return result.Interface()
	}

	// 如果不是结构体，直接返回
	if val.Kind() != reflect.Struct {
		return data
	}

	// 创建新的结构体实例
	newVal := reflect.New(typ).Elem()

	// 遍历结构体字段
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// 获取字段的JSON标签作为字段名
		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			// 处理json标签，去掉omitempty等选项
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				fieldName = jsonTag[:idx]
			} else {
				fieldName = jsonTag
			}
		}

		// 检查字段权限
		if perm, exists := permissions[fieldName]; exists {
			allowed := false
			switch operation {
			case "view":
				allowed = *perm.Visibility
			case "edit":
				allowed = *perm.EditPermission
			case "export":
				allowed = *perm.Exportable
			case "query":
				allowed = *perm.Queryable
			default:
				allowed = true
			}

			if !allowed {
				// 如果没有权限，设置为零值或跳过
				continue
			}
		}

		// 复制字段值
		if fieldVal.CanSet() && newVal.Field(i).CanSet() {
			newVal.Field(i).Set(fieldVal)
		}
	}

	return newVal.Interface()
}

// GetVisibleFields 获取用户可见的字段列表
func (p *PermissionHelper) GetVisibleFields(c *gin.Context, tableName string) ([]string, error) {
	// 从上下文获取用户信息
	authorityID, exists := c.Get("authorityId")
	if !exists {
		return nil, errors.New("无法获取用户权限信息")
	}

	// 获取字段权限配置
	middleware := &DataPermissionMiddleware{}
	fieldPermissions, err := middleware.getFieldPermissions(tableName, authorityID.(uint))
	if err != nil {
		return nil, err
	}

	var visibleFields []string
	for fieldName, perm := range fieldPermissions {
		if *perm.Visibility {
			visibleFields = append(visibleFields, fieldName)
		}
	}

	return visibleFields, nil
}

// GetEditableFields 获取用户可编辑的字段列表
func (p *PermissionHelper) GetEditableFields(c *gin.Context, tableName string) ([]string, error) {
	// 从上下文获取用户信息
	authorityID, exists := c.Get("authorityId")
	if !exists {
		return nil, errors.New("无法获取用户权限信息")
	}

	// 获取字段权限配置
	middleware := &DataPermissionMiddleware{}
	fieldPermissions, err := middleware.getFieldPermissions(tableName, authorityID.(uint))
	if err != nil {
		return nil, err
	}

	var editableFields []string
	for fieldName, perm := range fieldPermissions {
		if *perm.EditPermission {
			editableFields = append(editableFields, fieldName)
		}
	}

	return editableFields, nil
}

// ValidateDataAccess 验证用户是否有权限访问指定数据
func (p *PermissionHelper) ValidateDataAccess(c *gin.Context, tableName string, recordID uint) (bool, error) {
	// 获取数据权限SQL条件
	condition, err := p.GetDataPermissionSQL(c, tableName)
	if err != nil {
		return false, err
	}

	// 如果条件是1=1，表示有全部权限
	if condition == "1=1" {
		return true, nil
	}

	// 检查记录是否满足权限条件
	var count int64
	sql := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ? AND (%s)", tableName, condition)
	if err := global.GVA_DB.Raw(sql, recordID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// BuildSelectFields 根据权限构建SELECT字段列表
func (p *PermissionHelper) BuildSelectFields(c *gin.Context, tableName string, allFields []string) []string {
	visibleFields, err := p.GetVisibleFields(c, tableName)
	if err != nil {
		// 如果获取权限失败，返回所有字段
		return allFields
	}

	if len(visibleFields) == 0 {
		// 如果没有权限配置，返回所有字段
		return allFields
	}

	// 过滤出可见字段
	var result []string
	visibleMap := make(map[string]bool)
	for _, field := range visibleFields {
		visibleMap[field] = true
	}

	for _, field := range allFields {
		if visibleMap[field] {
			result = append(result, field)
		}
	}

	return result
}
