package datapermission

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dpGlobal "github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/router"
	interfaces "github.com/flipped-aurora/gin-vue-admin/server/utils/plugin"
	"github.com/gin-gonic/gin"
)

var _ interfaces.Plugin = (*plugin)(nil)

var Plugin = new(plugin)

type plugin struct{}

func (p *plugin) Register(group *gin.RouterGroup) {
	// 自动初始化插件配置和数据库表
	p.autoInit()
	// 注册路由
	router.RouterGroupApp.DataPermissionRouter.InitDataPermissionRouter(group)
}

func (p *plugin) RouterPath() string {
	return "datapermission"
}

// Install 保留Install方法以兼容手动安装场景
func (p *plugin) Install(ctx *gin.Context) error {
	// 初始化插件独立配置
	dpGlobal.InitConfig()
	return p.migrate(ctx)
}

func (p *plugin) Uninstall(ctx *gin.Context) error {
	// 卸载时可以选择是否删除数据表
	// 这里暂时不删除，避免数据丢失
	return nil
}

// autoInit 自动初始化插件（在Register时调用）
func (p *plugin) autoInit() {
	// 初始化插件独立配置
	dpGlobal.InitConfig()

	// 自动迁移数据库表
	p.autoMigrate()
}

// autoMigrate 自动迁移数据库表
func (p *plugin) autoMigrate() {
	// 检查数据库连接是否可用
	if global.GVA_DB == nil {
		global.GVA_LOG.Warn("数据权限插件: 数据库未初始化，跳过表迁移")
		return
	}

	// 自动迁移数据库表
	err := global.GVA_DB.AutoMigrate(
		model.ControlledTable{},
		model.RoleDataPermission{},
		model.RoleFieldPermission{},
	)
	if err != nil {
		global.GVA_LOG.Error("数据权限插件: 数据库表迁移失败")
	} else {
		global.GVA_LOG.Info("数据权限插件: 数据库表迁移成功")
	}
}

// migrate 保留原有migrate方法以兼容Install调用
func (p *plugin) migrate(ctx *gin.Context) error {
	// 检查数据库连接是否可用
	if global.GVA_DB == nil {
		return nil // 数据库未初始化，跳过迁移
	}

	// 自动迁移数据库表
	err := global.GVA_DB.AutoMigrate(
		model.ControlledTable{},
		model.RoleDataPermission{},
		model.RoleFieldPermission{},
	)
	return err
}
