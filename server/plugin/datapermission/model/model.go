package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// ControlledTable 权限受控表
type ControlledTable struct {
	global.GVA_MODEL
	Table       string `json:"tableName" form:"tableName" gorm:"column:table_name;comment:表名;not null;uniqueIndex"`
	Description string `json:"description" form:"description" gorm:"column:description;comment:表描述;"`
	Enabled     bool   `json:"enabled" form:"enabled" gorm:"column:enabled;comment:是否启用;default:true"`
	DataScope   string `json:"dataScope" form:"dataScope" gorm:"column:data_scope;comment:数据范围;default:'all'"`
	UserField   string `json:"userField" form:"userField" gorm:"column:user_field;comment:用户字段;"`
	DeptField   string `json:"deptField" form:"deptField" gorm:"column:dept_field;comment:部门字段;"`
}

// TableName ControlledTable 表名
func (ControlledTable) TableName() string {
	return "data_permission_controlled_tables"
}

// RoleDataPermission 角色数据权限配置
type RoleDataPermission struct {
	global.GVA_MODEL
	AuthorityID       uint                `json:"authorityID" form:"authorityID" gorm:"column:authority_id;comment:角色ID;not null"`
	Authority         system.SysAuthority `json:"authority" gorm:"foreignKey:AuthorityID;references:AuthorityId"`
	ControlledTableID uint                `json:"controlledTableID" form:"controlledTableID" gorm:"column:controlled_table_id;comment:受控表ID;not null"`
	ControlledTable   ControlledTable     `json:"controlledTable" gorm:"foreignKey:ControlledTableID;references:ID"`
	DataScope         string              `json:"dataScope" form:"dataScope" gorm:"column:data_scope;comment:数据范围;not null"`
	CustomCondition   string              `json:"customCondition" form:"customCondition" gorm:"column:custom_condition;comment:自定义SQL条件;type:text"`
	Priority          int                 `json:"priority" form:"priority" gorm:"column:priority;comment:优先级;default:50"`
	Enabled           bool                `json:"enabled" form:"enabled" gorm:"column:enabled;comment:是否启用;default:true"`
}

// TableName RoleDataPermission 表名
func (RoleDataPermission) TableName() string {
	return "data_permission_role_data"
}

// RoleFieldPermission 角色字段权限配置
type RoleFieldPermission struct {
	global.GVA_MODEL
	AuthorityID       uint                `json:"authorityID" form:"authorityID" gorm:"column:authority_id;comment:角色ID;not null"`
	Authority         system.SysAuthority `json:"authority" gorm:"foreignKey:AuthorityID;references:AuthorityId"`
	ControlledTableID uint                `json:"controlledTableID" form:"controlledTableID" gorm:"column:controlled_table_id;comment:受控表ID;not null"`
	ControlledTable   ControlledTable     `json:"controlledTable" gorm:"foreignKey:ControlledTableID;references:ID"`
	FieldName         string              `json:"fieldName" form:"fieldName" gorm:"column:field_name;comment:字段名;not null"`
	FieldChinese      string              `json:"fieldChinese" form:"fieldChinese" gorm:"column:field_chinese;comment:字段中文名;"`
	FieldDesc         string              `json:"fieldDesc" form:"fieldDesc" gorm:"column:field_desc;comment:字段描述;"`
	Visibility        *bool               `json:"visibility" form:"visibility" gorm:"column:visibility;comment:可见性;default:true"`
	EditPermission    *bool               `json:"editPermission" form:"editPermission" gorm:"column:edit_permission;comment:编辑权限;default:true"`
	Exportable        *bool               `json:"exportable" form:"exportable" gorm:"column:exportable;comment:导出权限;default:true"`
	Queryable         *bool               `json:"queryable" form:"queryable" gorm:"column:queryable;comment:查询权限;default:true"`
	Enabled           bool                `json:"enabled" form:"enabled" gorm:"column:enabled;comment:是否启用;default:true"`
}

// TableName RoleFieldPermission 表名
func (RoleFieldPermission) TableName() string {
	return "data_permission_role_fields"
}

// DataPermissionConfig 数据权限配置（用于前端保存配置）
type DataPermissionConfig struct {
	AuthorityID uint                  `json:"authorityID"`
	Table       string                `json:"tableName"`
	DataScope   string                `json:"dataScope"`
	CustomSQL   string                `json:"customSql"`
	Fields      []FieldPermissionItem `json:"fields"`
}

// FieldPermissionItem 字段权限项
type FieldPermissionItem struct {
	FieldName    string `json:"fieldName"`
	FieldChinese string `json:"fieldChinese"`
	FieldDesc    string `json:"fieldDesc"`
	Visible      *bool  `json:"visible"`
	Editable     *bool  `json:"editable"`
	Exportable   *bool  `json:"exportable"`
	Queryable    *bool  `json:"queryable"`
}

// TableInfo 表信息
type TableInfo struct {
	Table        string `json:"tableName"`
	TableComment string `json:"tableComment"`
}

// FieldInfo 字段信息
type FieldInfo struct {
	FieldName       string `json:"fieldName"`
	FieldChinese    string `json:"fieldChinese"`
	FieldDesc       string `json:"fieldDesc"`
	DataType        string `json:"dataType"`
	DefaultReadonly bool   `json:"defaultReadonly"`
}

// PermissionTestRequest 权限测试请求
type PermissionTestRequest struct {
	AuthorityID uint   `json:"authorityID"`
	TableName   string `json:"tableName"`
	UserID      uint   `json:"userID"`
}

// PermissionTestResponse 权限测试响应
type PermissionTestResponse struct {
	DataScope        string                         `json:"dataScope"`
	CustomCondition  string                         `json:"customCondition"`
	FieldPermissions map[string]FieldPermissionItem `json:"fieldPermissions"`
	SQLCondition     string                         `json:"sqlCondition"`
}
