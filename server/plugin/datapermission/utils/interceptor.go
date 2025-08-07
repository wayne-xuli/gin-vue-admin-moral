package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DataPermissionInterceptor 数据权限拦截器
type DataPermissionInterceptor struct {
	ginContext *gin.Context
}

// NewDataPermissionInterceptor 创建数据权限拦截器
func NewDataPermissionInterceptor(c *gin.Context) *DataPermissionInterceptor {
	return &DataPermissionInterceptor{
		ginContext: c,
	}
}

// Name 插件名称
func (interceptor *DataPermissionInterceptor) Name() string {
	return "data_permission_interceptor"
}

// Initialize 初始化插件
func (interceptor *DataPermissionInterceptor) Initialize(db *gorm.DB) error {
	// 注册查询回调
	err := db.Callback().Query().Before("gorm:query").Register("data_permission:before_query", interceptor.beforeQuery)
	if err != nil {
		return err
	}

	// 注册创建回调
	err = db.Callback().Create().Before("gorm:create").Register("data_permission:before_create", interceptor.beforeCreate)
	if err != nil {
		return err
	}

	// 注册更新回调
	err = db.Callback().Update().Before("gorm:update").Register("data_permission:before_update", interceptor.beforeUpdate)
	if err != nil {
		return err
	}

	// 注册删除回调
	err = db.Callback().Delete().Before("gorm:delete").Register("data_permission:before_delete", interceptor.beforeDelete)
	if err != nil {
		return err
	}

	// 注册查询后回调，用于字段过滤
	err = db.Callback().Query().After("gorm:after_query").Register("data_permission:after_query", interceptor.afterQuery)
	if err != nil {
		return err
	}

	return nil
}

// beforeQuery 查询前拦截
func (interceptor *DataPermissionInterceptor) beforeQuery(db *gorm.DB) {
	// 检查是否需要跳过数据权限
	if interceptor.shouldSkipDataPermission(db) {
		return
	}

	// 获取表名
	tableName := interceptor.getTableName(db)
	if tableName == "" {
		return
	}

	// 应用数据权限条件
	condition, err := interceptor.getDataPermissionCondition(tableName)
	if err != nil {
		// 记录错误但不中断查询
		global.GVA_LOG.Error(fmt.Sprintf("获取数据权限条件失败: %v", err))
		return
	}

	if condition == "1=1" {
		return
	}

	// 应用权限条件
	if condition != "" {
		// 检查是否已经应用了数据权限
		if applied, exists := db.Get("data_permission_applied"); exists && applied.(bool) {
			return
		}
		// 标记已应用数据权限
		db = db.Set("data_permission_applied", true)
		db.Where(condition)
	}
	// 应用字段权限过滤
	interceptor.applyFieldPermissionFilter(db, tableName)
}

// beforeCreate 创建前拦截
func (interceptor *DataPermissionInterceptor) beforeCreate(db *gorm.DB) {
	// 检查是否需要跳过数据权限
	if interceptor.shouldSkipDataPermission(db) {
		return
	}

	// 获取表名
	tableName := interceptor.getTableName(db)
	if tableName == "" {
		return
	}

	// 检查创建权限
	if !interceptor.checkCreatePermission(tableName) {
		db.AddError(errors.New("没有创建权限"))
		return
	}

	// 自动填充用户相关字段
	interceptor.autoFillUserFields(db, tableName)
}

// beforeUpdate 更新前拦截
func (interceptor *DataPermissionInterceptor) beforeUpdate(db *gorm.DB) {
	// 检查是否需要跳过数据权限
	if interceptor.shouldSkipDataPermission(db) {
		return
	}

	// 获取表名
	tableName := interceptor.getTableName(db)
	if tableName == "" {
		return
	}

	// 应用数据权限条件（确保只能更新有权限的数据）
	condition, err := interceptor.getDataPermissionCondition(tableName)
	if err != nil {
		db.AddError(fmt.Errorf("获取数据权限条件失败: %v", err))
		return
	}

	// 应用权限条件
	if condition != "" && condition != "1=1" {
		// 检查是否已经应用了数据权限
		if applied, exists := db.Get("data_permission_applied"); exists && applied.(bool) {
			goto checkFieldPermission
		}
		// 标记已应用数据权限
		db = db.Set("data_permission_applied", true)
		db.Where(condition)
	}

checkFieldPermission:
	// 检查字段编辑权限
	interceptor.checkFieldEditPermission(db, tableName)
}

// beforeDelete 删除前拦截
func (interceptor *DataPermissionInterceptor) beforeDelete(db *gorm.DB) {
	// 检查是否需要跳过数据权限
	if interceptor.shouldSkipDataPermission(db) {
		return
	}

	// 获取表名
	tableName := interceptor.getTableName(db)
	if tableName == "" {
		return
	}

	// 应用数据权限条件（确保只能删除有权限的数据）
	condition, err := interceptor.getDataPermissionCondition(tableName)
	if err != nil {
		db.AddError(fmt.Errorf("获取数据权限条件失败: %v", err))
		return
	}

	// 应用权限条件
	if condition != "" && condition != "1=1" {
		// 检查是否已经应用了数据权限
		if applied, exists := db.Get("data_permission_applied"); exists && applied.(bool) {
			return
		}
		// 标记已应用数据权限
		db = db.Set("data_permission_applied", true)
		db.Where(condition)
	}
}

// afterQuery 查询后拦截，用于字段过滤
func (interceptor *DataPermissionInterceptor) afterQuery(db *gorm.DB) {
	// 检查是否需要跳过数据权限
	if interceptor.shouldSkipDataPermission(db) {
		return
	}

	// 获取隐藏字段信息
	hiddenFieldsValue, exists := db.Get("hidden_fields")
	if !exists {
		return
	}

	hiddenFields, ok := hiddenFieldsValue.([]string)
	if !ok || len(hiddenFields) == 0 {
		return
	}

	// 对查询结果进行字段过滤
	interceptor.filterResultFields(db, hiddenFields)
}

// shouldSkipDataPermission 检查是否应该跳过数据权限
func (interceptor *DataPermissionInterceptor) shouldSkipDataPermission(db *gorm.DB) bool {
	// 如果没有Gin上下文，跳过权限检查
	if interceptor.ginContext == nil {
		return true
	}

	// 检查是否显式跳过数据权限
	if skipValue, exists := db.Get("skip_data_permission"); exists {
		if skip, ok := skipValue.(bool); ok && skip {
			return true
		}
	}

	// 检查是否为系统内部操作
	if isSystemOperation, exists := db.Get("system_operation"); exists {
		if system, ok := isSystemOperation.(bool); ok && system {
			return true
		}
	}

	return false
}

// getTableName 获取表名
func (interceptor *DataPermissionInterceptor) getTableName(db *gorm.DB) string {
	if db.Statement == nil {
		return ""
	}

	// 优先使用显式设置的表名
	if tableName, exists := db.Get("table_name"); exists {
		if name, ok := tableName.(string); ok {
			return name
		}
	}

	// 从Statement获取表名
	if db.Statement.Table != "" {
		return db.Statement.Table
	}

	// 从模型获取表名
	if db.Statement.Model != nil {
		return db.NamingStrategy.TableName(reflect.TypeOf(db.Statement.Model).Elem().Name())
	}

	return ""
}

// getDataPermissionCondition 获取数据权限条件
func (interceptor *DataPermissionInterceptor) getDataPermissionCondition(tableName string) (string, error) {
	// 从上下文获取用户信息
	userID, exists := interceptor.ginContext.Get("userID")
	if !exists {
		return "1=1", nil
	}
	middleware := &DataPermissionMiddleware{}
	allConditions := make([]string, 0)
	if authorityIds, exists := interceptor.ginContext.Get("authorityIds"); exists {
		if roleIds, ok := authorityIds.([]uint); ok {
			for _, authorityID := range roleIds {
				condition, err := middleware.getDataPermissionCondition(tableName, authorityID, userID.(uint))
				if err != nil {
					return "", err
				}
				// allConditions 中已经包含 condition则不执行下面的逻辑
				if strings.Contains(strings.Join(allConditions, " OR "), condition) {
					continue
				}
				allConditions = append(allConditions, condition)
			}
		}
	}
	if len(allConditions) > 0 {
		finalCondition := strings.Join(allConditions, " OR ")
		return finalCondition, nil
	}
	return "", nil
}

// checkCreatePermission 检查创建权限
func (interceptor *DataPermissionInterceptor) checkCreatePermission(tableName string) bool {
	// 检查表是否受控制
	var controlledTable model.ControlledTable
	if err := global.GVA_DB.Where("table_name = ? AND enabled = ?", tableName, true).First(&controlledTable).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true // 如果表不受控制，默认允许
		}
		return false
	}

	// 获取用户权限信息
	authorityID, exists := interceptor.ginContext.Get("authorityId")
	if !exists {
		return true
	}

	// 检查角色是否有权限配置
	var count int64
	global.GVA_DB.Model(&model.RoleDataPermission{}).Where("authority_id = ? AND controlled_table_id = ? AND enabled = ?", authorityID, controlledTable.ID, true).Count(&count)
	return count > 0
}

// autoFillUserFields 自动填充用户相关字段
func (interceptor *DataPermissionInterceptor) autoFillUserFields(db *gorm.DB, tableName string) {
	// 获取受控表配置
	var controlledTable model.ControlledTable
	if err := global.GVA_DB.Where("table_name = ? AND enabled = ?", tableName, true).First(&controlledTable).Error; err != nil {
		return
	}

	// 获取用户ID
	userID, exists := interceptor.ginContext.Get("userID")
	if !exists {
		return
	}

	// 自动填充用户字段
	if controlledTable.UserField != "" {
		db.Set(controlledTable.UserField, userID)
	}

	// 自动填充部门字段（如果需要）
	if controlledTable.DeptField != "" {
		// 获取用户部门ID
		var userDeptID uint
		global.GVA_DB.Raw("SELECT dept_id FROM sys_users WHERE id = ?", userID).Scan(&userDeptID)
		if userDeptID > 0 {
			db.Set(controlledTable.DeptField, userDeptID)
		}
	}
}

// getFieldPermissions 获取字段权限配置（在interceptor中的实现）
func (interceptor *DataPermissionInterceptor) getFieldPermissions(tableName string, authorityID uint) (map[string]model.RoleFieldPermission, error) {
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

// checkFieldEditPermission 检查字段编辑权限
func (interceptor *DataPermissionInterceptor) checkFieldEditPermission(db *gorm.DB, tableName string) {
	// 获取用户权限信息
	authorityID, exists := interceptor.ginContext.Get("authorityId")
	if !exists {
		return
	}

	// 获取字段权限配置
	fieldPermissions, err := interceptor.getFieldPermissions(tableName, authorityID.(uint))
	if err != nil || len(fieldPermissions) == 0 {
		return
	}

	// 检查更新的字段是否有编辑权限
	if db.Statement.Dest != nil {
		// 这里可以进一步实现字段级别的权限检查
		// 由于GORM的复杂性，这里提供基础框架
		// 具体实现可能需要根据业务需求调整
	}
}

// applyFieldPermissionFilter 应用字段权限过滤
func (interceptor *DataPermissionInterceptor) applyFieldPermissionFilter(db *gorm.DB, tableName string) {
	// 获取用户权限信息
	allFieldPermissions := make(map[string]model.RoleFieldPermission)
	if authorityIds, exists := interceptor.ginContext.Get("authorityIds"); exists {
		if roleIds, ok := authorityIds.([]uint); ok {
			for _, roleId := range roleIds {
				fieldPermissions, err := interceptor.getFieldPermissions(tableName, roleId)
				if err != nil || len(fieldPermissions) == 0 {
					continue
				}
				// 合并字段权限，采用最宽松的权限策略
				for fieldName, permission := range fieldPermissions {
					if existingPerm, exists := allFieldPermissions[fieldName]; exists {
						// 合并权限：任一角色有权限则允许
						mergedPerm := interceptor.mergeFieldPermissions(existingPerm, permission)
						allFieldPermissions[fieldName] = mergedPerm
					} else {
						allFieldPermissions[fieldName] = permission
					}
				}
			}
		}
	}

	// 应用字段权限过滤
	interceptor.applyFieldVisibilityFilter(db, allFieldPermissions)
	interceptor.applyFieldEditFilter(db, allFieldPermissions)
	interceptor.applyFieldQueryFilter(db, allFieldPermissions)
}

// mergeFieldPermissions 合并字段权限（采用最宽松策略）
func (interceptor *DataPermissionInterceptor) mergeFieldPermissions(perm1, perm2 model.RoleFieldPermission) model.RoleFieldPermission {
	merged := perm1

	// 可见性：任一为true则为true
	if perm1.Visibility != nil && perm2.Visibility != nil {
		visible := *perm1.Visibility || *perm2.Visibility
		merged.Visibility = &visible
	} else if perm1.Visibility == nil && perm2.Visibility != nil {
		merged.Visibility = perm2.Visibility
	}

	// 编辑权限：任一为true则为true
	if perm1.EditPermission != nil && perm2.EditPermission != nil {
		editable := *perm1.EditPermission || *perm2.EditPermission
		merged.EditPermission = &editable
	} else if perm1.EditPermission == nil && perm2.EditPermission != nil {
		merged.EditPermission = perm2.EditPermission
	}

	// 导出权限：任一为true则为true
	if perm1.Exportable != nil && perm2.Exportable != nil {
		exportable := *perm1.Exportable || *perm2.Exportable
		merged.Exportable = &exportable
	} else if perm1.Exportable == nil && perm2.Exportable != nil {
		merged.Exportable = perm2.Exportable
	}

	// 查询权限：任一为true则为true
	if perm1.Queryable != nil && perm2.Queryable != nil {
		queryable := *perm1.Queryable || *perm2.Queryable
		merged.Queryable = &queryable
	} else if perm1.Queryable == nil && perm2.Queryable != nil {
		merged.Queryable = perm2.Queryable
	}

	return merged
}

// applyFieldVisibilityFilter 应用字段可见性过滤
func (interceptor *DataPermissionInterceptor) applyFieldVisibilityFilter(db *gorm.DB, fieldPermissions map[string]model.RoleFieldPermission) {
	// 对于查询操作，过滤不可见字段
	if db.Statement.SQL.String() == "" { // 只在构建查询时应用
		// 检查是否有字段权限配置，如果没有则跳过
		if len(fieldPermissions) == 0 {
			return
		}

		hiddenFields := make([]string, 0)
		for fieldName, permission := range fieldPermissions {
			if permission.Visibility != nil && !*permission.Visibility {
				hiddenFields = append(hiddenFields, fieldName)
			}
		}

		// 如果没有隐藏字段，则不需要应用过滤
		if len(hiddenFields) == 0 {
			return
		}

		// 只存储隐藏字段信息，不直接修改SELECT语句
		// 这样可以避免影响其他查询，由具体的业务逻辑决定如何处理
		db.Set("hidden_fields", hiddenFields)
		db.Set("field_permissions", fieldPermissions)
	}
}

// getModelFields 获取模型的所有字段名
func (interceptor *DataPermissionInterceptor) getModelFields(db *gorm.DB) []string {
	fields := make([]string, 0)

	// 通过GORM的Schema获取字段信息
	if db.Statement.Schema != nil {
		for _, field := range db.Statement.Schema.Fields {
			// 只包含数据库字段，排除关联字段
			if field.DBName != "" {
				fields = append(fields, field.DBName)
			}
		}
	} else if db.Statement.Model != nil {
		// 如果Schema为空，尝试解析模型
		err := db.Statement.Parse(db.Statement.Model)
		if err == nil && db.Statement.Schema != nil {
			for _, field := range db.Statement.Schema.Fields {
				if field.DBName != "" {
					fields = append(fields, field.DBName)
				}
			}
		}
	}

	return fields
}

// intersectSelectFields 将现有SELECT字段与可见字段取交集
func (interceptor *DataPermissionInterceptor) intersectSelectFields(db *gorm.DB, visibleFields []string) {
	visibleFieldsMap := make(map[string]bool)
	for _, field := range visibleFields {
		visibleFieldsMap[field] = true
	}

	// 过滤现有的SELECT字段
	filteredSelects := make([]string, 0)
	for _, selectField := range db.Statement.Selects {
		// 简单的字段名匹配，可能需要更复杂的解析
		fieldName := strings.TrimSpace(selectField)
		// 移除表别名（如果有）
		if dotIndex := strings.LastIndex(fieldName, "."); dotIndex != -1 {
			fieldName = fieldName[dotIndex+1:]
		}
		// 移除AS别名（如果有）
		if asIndex := strings.Index(strings.ToUpper(fieldName), " AS "); asIndex != -1 {
			fieldName = strings.TrimSpace(fieldName[:asIndex])
		}

		if visibleFieldsMap[fieldName] {
			filteredSelects = append(filteredSelects, selectField)
		}
	}

	// 重新设置SELECT字段
	db.Statement.Selects = filteredSelects
}

// applyFieldEditFilter 应用字段编辑权限过滤
func (interceptor *DataPermissionInterceptor) applyFieldEditFilter(db *gorm.DB, fieldPermissions map[string]model.RoleFieldPermission) {
	// 对于更新操作，检查编辑权限
	if db.Statement.SQL.String() != "" && strings.Contains(strings.ToUpper(db.Statement.SQL.String()), "UPDATE") {
		readonlyFields := make([]string, 0)
		for fieldName, permission := range fieldPermissions {
			if permission.EditPermission != nil && !*permission.EditPermission {
				readonlyFields = append(readonlyFields, fieldName)
			}
		}

		// 存储只读字段信息
		if len(readonlyFields) > 0 {
			db.Set("readonly_fields", readonlyFields)
			// 从更新字段中移除只读字段
			for _, field := range readonlyFields {
				if db.Statement.Dest != nil {
					// 通过反射移除只读字段的值
					interceptor.removeFieldFromUpdate(db, field)
				}
			}
		}
	}
}

// applyFieldQueryFilter 应用字段查询权限过滤
func (interceptor *DataPermissionInterceptor) applyFieldQueryFilter(db *gorm.DB, fieldPermissions map[string]model.RoleFieldPermission) {
	// 检查WHERE条件中是否包含不可查询的字段
	for fieldName, permission := range fieldPermissions {
		if permission.Queryable != nil && !*permission.Queryable {
			// 检查WHERE条件中是否使用了该字段
			if db.Statement.SQL.String() != "" && strings.Contains(db.Statement.SQL.String(), fieldName) {
				// 记录警告或阻止查询
				global.GVA_LOG.Warn(fmt.Sprintf("Field %s is not queryable for current user", fieldName))
			}
		}
	}
}

// filterResultFields 过滤查询结果中的隐藏字段
func (interceptor *DataPermissionInterceptor) filterResultFields(db *gorm.DB, hiddenFields []string) {
	if db.Statement.Dest == nil {
		return
	}

	// 创建隐藏字段映射
	hiddenFieldsMap := make(map[string]bool)
	for _, field := range hiddenFields {
		hiddenFieldsMap[field] = true
	}

	// 获取目标值
	destValue := reflect.ValueOf(db.Statement.Dest)
	if destValue.Kind() != reflect.Ptr {
		return
	}
	destValue = destValue.Elem()

	// 处理不同类型的结果
	switch destValue.Kind() {
	case reflect.Slice:
		// 处理切片类型（多条记录）
		for i := 0; i < destValue.Len(); i++ {
			item := destValue.Index(i)
			interceptor.filterStructFields(item, hiddenFieldsMap)
		}
	case reflect.Struct:
		// 处理单个结构体
		interceptor.filterStructFields(destValue, hiddenFieldsMap)
	}
}

// filterStructFields 过滤结构体中的隐藏字段
func (interceptor *DataPermissionInterceptor) filterStructFields(structValue reflect.Value, hiddenFieldsMap map[string]bool) {
	if structValue.Kind() == reflect.Ptr {
		if structValue.IsNil() {
			return
		}
		structValue = structValue.Elem()
	}

	if structValue.Kind() != reflect.Struct {
		return
	}

	structType := structValue.Type()
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// 跳过不可设置的字段
		if !field.CanSet() {
			continue
		}

		// 获取数据库字段名
		dbFieldName := interceptor.getDBFieldName(fieldType)
		if dbFieldName == "" {
			continue
		}

		// 如果是隐藏字段，设置为零值
		if hiddenFieldsMap[dbFieldName] {
			field.Set(reflect.Zero(field.Type()))
		}
	}
}

// getDBFieldName 获取字段的数据库名称
func (interceptor *DataPermissionInterceptor) getDBFieldName(field reflect.StructField) string {
	// 检查gorm标签
	if gormTag := field.Tag.Get("gorm"); gormTag != "" {
		// 解析gorm标签中的column名称
		parts := strings.Split(gormTag, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "column:") {
				return strings.TrimPrefix(part, "column:")
			}
		}
	}

	// 检查json标签作为备选
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		parts := strings.Split(jsonTag, ",")
		if len(parts) > 0 && parts[0] != "-" {
			return parts[0]
		}
	}

	// 使用字段名的snake_case形式
	return interceptor.toSnakeCase(field.Name)
}

// toSnakeCase 将驼峰命名转换为蛇形命名
func (interceptor *DataPermissionInterceptor) toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// removeFieldFromUpdate 从更新操作中移除指定字段
func (interceptor *DataPermissionInterceptor) removeFieldFromUpdate(db *gorm.DB, fieldName string) {
	// 如果是map类型的更新
	if updateMap, ok := db.Statement.Dest.(map[string]interface{}); ok {
		delete(updateMap, fieldName)
		return
	}

	// 如果是结构体类型的更新，通过反射处理
	if db.Statement.Dest != nil {
		destValue := reflect.ValueOf(db.Statement.Dest)
		if destValue.Kind() == reflect.Ptr {
			destValue = destValue.Elem()
		}

		if destValue.Kind() == reflect.Struct {
			// 查找对应的字段
			for i := 0; i < destValue.NumField(); i++ {
				field := destValue.Type().Field(i)
				// 检查gorm标签或字段名
				gormTag := field.Tag.Get("gorm")
				columnName := fieldName
				if strings.Contains(gormTag, "column:") {
					// 解析column名称
					parts := strings.Split(gormTag, ";")
					for _, part := range parts {
						if strings.HasPrefix(part, "column:") {
							columnName = strings.TrimPrefix(part, "column:")
							break
						}
					}
				}

				if columnName == fieldName || strings.ToLower(field.Name) == strings.ToLower(fieldName) {
					// 将字段设置为零值
					fieldValue := destValue.Field(i)
					if fieldValue.CanSet() {
						fieldValue.Set(reflect.Zero(fieldValue.Type()))
					}
					break
				}
			}
		}
	}
}

// DataPermissionInterceptorMiddleware 数据权限拦截器中间件
// 这个中间件会为每个请求创建一个数据权限拦截器，并注册到GORM中
func DataPermissionInterceptorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 创建数据权限拦截器
		interceptor := NewDataPermissionInterceptor(c)

		// 为当前请求创建一个带有拦截器的数据库实例
		db := global.GVA_DB.Session(&gorm.Session{})

		// 注册拦截器
		if err := interceptor.Initialize(db); err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("初始化数据权限拦截器失败: %v", err))
		}

		// 将带有拦截器的数据库实例存储到上下文中
		c.Set("interceptor_db", db)
		c.Set("data_permission_interceptor", interceptor)

		c.Next()
	}
}

// GetInterceptorDB 从上下文获取带有拦截器的数据库实例
func GetInterceptorDB(c *gin.Context) *gorm.DB {
	if db, exists := c.Get("interceptor_db"); exists {
		return db.(*gorm.DB)
	}
	return global.GVA_DB
}

// SkipDataPermission 跳过数据权限检查
func SkipDataPermission(db *gorm.DB) *gorm.DB {
	return db.Set("skip_data_permission", true)
}

// SetSystemOperation 标记为系统操作
func SetSystemOperation(db *gorm.DB) *gorm.DB {
	return db.Set("system_operation", true)
}

// SetTableName 显式设置表名
func SetTableName(db *gorm.DB, tableName string) *gorm.DB {
	return db.Set("table_name", tableName)
}
