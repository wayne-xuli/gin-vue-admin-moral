package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/model"
)

// ControlledTableResponse 受控表响应
type ControlledTableResponse struct {
	model.ControlledTable
	RolePermissionCount  int `json:"rolePermissionCount"`  // 角色权限数量
	FieldPermissionCount int `json:"fieldPermissionCount"` // 字段权限数量
}

// RoleDataPermissionResponse 角色数据权限响应
type RoleDataPermissionResponse struct {
	model.RoleDataPermission
	AuthorityName string `json:"authorityName"` // 角色名称
	TableName     string `json:"tableName"`     // 表名
}

// RoleFieldPermissionResponse 角色字段权限响应
type RoleFieldPermissionResponse struct {
	model.RoleFieldPermission
	AuthorityName string `json:"authorityName"` // 角色名称
	TableName     string `json:"tableName"`     // 表名
}

// DataPermissionConfigResponse 数据权限配置响应
type DataPermissionConfigResponse struct {
	AuthorityID   uint                        `json:"authorityID"`
	AuthorityName string                      `json:"authorityName"`
	Table         string                      `json:"table"`
	DataScope     string                      `json:"dataScope"`
	CustomSQL     string                      `json:"customSql"`
	Fields        []model.FieldPermissionItem `json:"fields"`
}

// TableListResponse 表列表响应
type TableListResponse struct {
	Tables []model.TableInfo `json:"tables"`
}

// FieldListResponse 字段列表响应
type FieldListResponse struct {
	Fields []model.FieldInfo `json:"fields"`
}

// PermissionTestResponse 权限测试响应
type PermissionTestResponse struct {
	model.PermissionTestResponse
	Message string `json:"message"` // 测试结果消息
}

// DataScopeOption 数据范围选项
type DataScopeOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// DataScopeOptionsResponse 数据范围选项响应
type DataScopeOptionsResponse struct {
	Options []DataScopeOption `json:"options"`
}

// AuthorityListResponse 角色列表响应
type AuthorityListResponse struct {
	Authorities []AuthorityInfo `json:"authorities"`
}

// AuthorityInfo 角色信息
type AuthorityInfo struct {
	AuthorityID   uint   `json:"authorityID"`
	AuthorityName string `json:"authorityName"`
	ParentID      uint   `json:"parentID"`
}

// StatisticsResponse 统计信息响应
type StatisticsResponse struct {
	ControlledTableCount     int `json:"controlledTableCount"`     // 受控表数量
	RoleDataPermissionCount  int `json:"roleDataPermissionCount"`  // 角色数据权限数量
	RoleFieldPermissionCount int `json:"roleFieldPermissionCount"` // 角色字段权限数量
	ActiveAuthorityCount     int `json:"activeAuthorityCount"`     // 活跃角色数量
}
