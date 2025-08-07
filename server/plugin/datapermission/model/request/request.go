package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/model"
)

// GetByIDRequest 根据ID获取请求
type GetByIDRequest struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

// ControlledTableSearch 受控表搜索条件
type ControlledTableSearch struct {
	request.PageInfo
	Table       string `json:"table" form:"table"`
	Description string `json:"description" form:"description"`
	Enabled     *bool  `json:"enabled" form:"enabled"`
}

// CreateControlledTableRequest 创建受控表请求
type CreateControlledTableRequest struct {
	Table            string                      `json:"table" binding:"required"`
	Description      string                      `json:"description"`
	Enabled          bool                        `json:"enabled"`
	DataScope        string                      `json:"dataScope"`
	UserField        string                      `json:"userField"`
	DeptField        string                      `json:"deptField"`
	FieldPermissions []model.FieldPermissionItem `json:"fieldPermissions"`
}

// UpdateControlledTableRequest 更新受控表请求
type UpdateControlledTableRequest struct {
	ID               uint                        `json:"id" binding:"required"`
	Table            string                      `json:"table" binding:"required"`
	Description      string                      `json:"description"`
	Enabled          bool                        `json:"enabled"`
	DataScope        string                      `json:"dataScope"`
	UserField        string                      `json:"userField"`
	DeptField        string                      `json:"deptField"`
	FieldPermissions []model.FieldPermissionItem `json:"fieldPermissions"`
}

// RoleDataPermissionSearch 角色数据权限搜索条件
type RoleDataPermissionSearch struct {
	request.PageInfo
	AuthorityID       uint   `json:"authorityID" form:"authorityID"`
	ControlledTableID uint   `json:"controlledTableID" form:"controlledTableID"`
	DataScope         string `json:"dataScope" form:"dataScope"`
	Enabled           *bool  `json:"enabled" form:"enabled"`
}

// CreateRoleDataPermissionRequest 创建角色数据权限请求
type CreateRoleDataPermissionRequest struct {
	AuthorityID       uint   `json:"authorityID" binding:"required"`
	ControlledTableID uint   `json:"controlledTableID" binding:"required"`
	DataScope         string `json:"dataScope" binding:"required"`
	CustomCondition   string `json:"customCondition"`
	Priority          int    `json:"priority"`
	Enabled           bool   `json:"enabled"`
}

// UpdateRoleDataPermissionRequest 更新角色数据权限请求
type UpdateRoleDataPermissionRequest struct {
	ID                uint   `json:"id" binding:"required"`
	AuthorityID       uint   `json:"authorityID" binding:"required"`
	ControlledTableID uint   `json:"controlledTableID" binding:"required"`
	DataScope         string `json:"dataScope" binding:"required"`
	CustomCondition   string `json:"customCondition"`
	Priority          int    `json:"priority"`
	Enabled           bool   `json:"enabled"`
}

// RoleFieldPermissionSearch 角色字段权限搜索条件
type RoleFieldPermissionSearch struct {
	request.PageInfo
	AuthorityID       uint   `json:"authorityID" form:"authorityID"`
	ControlledTableID uint   `json:"controlledTableID" form:"controlledTableID"`
	FieldName         string `json:"fieldName" form:"fieldName"`
	Enabled           *bool  `json:"enabled" form:"enabled"`
}

// CreateRoleFieldPermissionRequest 创建角色字段权限请求
type CreateRoleFieldPermissionRequest struct {
	AuthorityID       uint   `json:"authorityID" binding:"required"`
	ControlledTableID uint   `json:"controlledTableID" binding:"required"`
	FieldName         string `json:"fieldName" binding:"required"`
	FieldChinese      string `json:"fieldChinese"`
	FieldDesc         string `json:"fieldDesc"`
	Visibility        string `json:"visibility"`
	EditPermission    string `json:"editPermission"`
	Exportable        bool   `json:"exportable"`
	Queryable         bool   `json:"queryable"`
	Enabled           bool   `json:"enabled"`
}

// UpdateRoleFieldPermissionRequest 更新角色字段权限请求
type UpdateRoleFieldPermissionRequest struct {
	ID                uint   `json:"id" binding:"required"`
	AuthorityID       uint   `json:"authorityID" binding:"required"`
	ControlledTableID uint   `json:"controlledTableID" binding:"required"`
	FieldName         string `json:"fieldName" binding:"required"`
	FieldChinese      string `json:"fieldChinese"`
	FieldDesc         string `json:"fieldDesc"`
	Visibility        string `json:"visibility"`
	EditPermission    string `json:"editPermission"`
	Exportable        bool   `json:"exportable"`
	Queryable         bool   `json:"queryable"`
	Enabled           bool   `json:"enabled"`
}

// SaveDataPermissionConfigRequest 保存数据权限配置请求
type SaveDataPermissionConfigRequest struct {
	AuthorityID uint                        `json:"authorityID" binding:"required"`
	Table       string                      `json:"table" binding:"required"`
	DataScope   string                      `json:"dataScope" binding:"required"`
	CustomSQL   string                      `json:"customSql"`
	Fields      []model.FieldPermissionItem `json:"fields"`
	UserField   string                      `json:"userField"`
	DeptField   string                      `json:"deptField"`
}

// GetDataPermissionConfigRequest 获取数据权限配置请求
type GetDataPermissionConfigRequest struct {
	AuthorityID uint   `json:"authorityID" form:"authorityID" binding:"required"`
	Table       string `json:"tableName" form:"tableName" binding:"required"`
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// PermissionTestRequest 权限测试请求
type PermissionTestRequest struct {
	AuthorityID uint   `json:"authorityId" form:"authorityId"` // 角色ID
	Table       string `json:"tableName" form:"tableName"`     // 表名
	UserID      uint   `json:"userId" form:"userId"`           // 用户ID
}

// PermissionTestResponse 权限测试响应
type PermissionTestResponse struct {
	DataScope        string                               `json:"dataScope"`        // 数据范围
	CustomCondition  string                               `json:"customCondition"`  // 自定义条件
	SQLCondition     string                               `json:"sqlCondition"`     // 生成的SQL条件
	FieldPermissions map[string]model.FieldPermissionItem `json:"fieldPermissions"` // 字段权限
}
