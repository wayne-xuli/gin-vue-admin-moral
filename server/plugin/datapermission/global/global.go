package global

import (
	"fmt"
	"sync"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/config"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/model"
	"go.uber.org/zap"
)

// 全局变量
var (
	// GVA_DP_CONFIG 数据权限配置
	GVA_DP_CONFIG config.DataPermission

	// GVA_DP_LOGGER 数据权限日志记录器
	GVA_DP_LOGGER *zap.Logger

	// GVA_DP_CACHE 权限缓存
	GVA_DP_CACHE *PermissionCache

	// GVA_DP_INIT_FLAG 初始化标志
	GVA_DP_INIT_FLAG bool
)

// InitConfig 初始化插件配置
func InitConfig() {
	// 从独立配置文件加载配置
	GVA_DP_CONFIG = config.LoadConfig()

	// 初始化缓存
	GVA_DP_CACHE = NewPermissionCache()

	// 设置初始化标志
	GVA_DP_INIT_FLAG = true
}

// PermissionCache 权限缓存结构
type PermissionCache struct {
	mu sync.RWMutex
	// 数据权限缓存 key: authorityId_tableName
	dataPermissions map[string]*model.RoleDataPermission
	// 字段权限缓存 key: authorityId_tableName_fieldName
	fieldPermissions map[string]*model.RoleFieldPermission
	// 受控表缓存 key: tableName
	controlledTables map[string]*model.ControlledTable
}

// NewPermissionCache 创建新的权限缓存
func NewPermissionCache() *PermissionCache {
	return &PermissionCache{
		dataPermissions:  make(map[string]*model.RoleDataPermission),
		fieldPermissions: make(map[string]*model.RoleFieldPermission),
		controlledTables: make(map[string]*model.ControlledTable),
	}
}

// GetDataPermission 获取数据权限缓存
func (c *PermissionCache) GetDataPermission(authorityId uint, tableName string) (*model.RoleDataPermission, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	key := c.buildDataPermissionKey(authorityId, tableName)
	perm, exists := c.dataPermissions[key]
	return perm, exists
}

// SetDataPermission 设置数据权限缓存
func (c *PermissionCache) SetDataPermission(authorityId uint, tableName string, permission *model.RoleDataPermission) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := c.buildDataPermissionKey(authorityId, tableName)
	c.dataPermissions[key] = permission
}

// GetFieldPermission 获取字段权限缓存
func (c *PermissionCache) GetFieldPermission(authorityId uint, tableName, fieldName string) (*model.RoleFieldPermission, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	key := c.buildFieldPermissionKey(authorityId, tableName, fieldName)
	perm, exists := c.fieldPermissions[key]
	return perm, exists
}

// SetFieldPermission 设置字段权限缓存
func (c *PermissionCache) SetFieldPermission(authorityId uint, tableName, fieldName string, permission *model.RoleFieldPermission) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := c.buildFieldPermissionKey(authorityId, tableName, fieldName)
	c.fieldPermissions[key] = permission
}

// GetControlledTable 获取受控表缓存
func (c *PermissionCache) GetControlledTable(tableName string) (*model.ControlledTable, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	table, exists := c.controlledTables[tableName]
	return table, exists
}

// SetControlledTable 设置受控表缓存
func (c *PermissionCache) SetControlledTable(tableName string, table *model.ControlledTable) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.controlledTables[tableName] = table
}

// ClearDataPermissionCache 清除数据权限缓存
func (c *PermissionCache) ClearDataPermissionCache(authorityId uint, tableName string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := c.buildDataPermissionKey(authorityId, tableName)
	delete(c.dataPermissions, key)
}

// ClearFieldPermissionCache 清除字段权限缓存
func (c *PermissionCache) ClearFieldPermissionCache(authorityId uint, tableName string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 清除该角色和表的所有字段权限缓存
	prefix := c.buildFieldPermissionPrefix(authorityId, tableName)
	for key := range c.fieldPermissions {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.fieldPermissions, key)
		}
	}
}

// ClearControlledTableCache 清除受控表缓存
func (c *PermissionCache) ClearControlledTableCache(tableName string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.controlledTables, tableName)
}

// ClearAllCache 清除所有缓存
func (c *PermissionCache) ClearAllCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dataPermissions = make(map[string]*model.RoleDataPermission)
	c.fieldPermissions = make(map[string]*model.RoleFieldPermission)
	c.controlledTables = make(map[string]*model.ControlledTable)
}

// buildDataPermissionKey 构建数据权限缓存键
func (c *PermissionCache) buildDataPermissionKey(authorityId uint, tableName string) string {
	return fmt.Sprintf("%d_%s", authorityId, tableName)
}

// buildFieldPermissionKey 构建字段权限缓存键
func (c *PermissionCache) buildFieldPermissionKey(authorityId uint, tableName, fieldName string) string {
	return fmt.Sprintf("%d_%s_%s", authorityId, tableName, fieldName)
}

// buildFieldPermissionPrefix 构建字段权限缓存键前缀
func (c *PermissionCache) buildFieldPermissionPrefix(authorityId uint, tableName string) string {
	return fmt.Sprintf("%d_%s_", authorityId, tableName)
}

// InitGlobal 初始化全局变量
func InitGlobal() {
	if GVA_DP_INIT_FLAG {
		return
	}

	// 初始化配置
	GVA_DP_CONFIG = config.GetDefaultConfig()

	// 初始化缓存
	GVA_DP_CACHE = NewPermissionCache()

	// 设置初始化标志
	GVA_DP_INIT_FLAG = true
}

// GetConfig 获取配置
func GetConfig() config.DataPermission {
	return GVA_DP_CONFIG
}

// SetConfig 设置配置
func SetConfig(cfg config.DataPermission) {
	GVA_DP_CONFIG = cfg
}

// GetCache 获取缓存
func GetCache() *PermissionCache {
	return GVA_DP_CACHE
}

// IsEnabled 检查数据权限是否启用
func IsEnabled() bool {
	return GVA_DP_CONFIG.Enabled
}

// IsFieldPermissionEnabled 检查字段权限是否启用
func IsFieldPermissionEnabled() bool {
	return GVA_DP_CONFIG.FieldPermissionEnabled
}

// IsCacheEnabled 检查缓存是否启用
func IsCacheEnabled() bool {
	return GVA_DP_CONFIG.Cache.Enabled
}

// IsLogEnabled 检查日志是否启用
func IsLogEnabled() bool {
	return GVA_DP_CONFIG.Log.Enabled
}
