package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/api"
	"github.com/gin-gonic/gin"
)

type DataPermissionRouter struct{}

// InitDataPermissionRouter 初始化数据权限路由信息
func (s *DataPermissionRouter) InitDataPermissionRouter(Router *gin.RouterGroup) {
	dataPermissionRouter := Router.Use(middleware.OperationRecord())
	dataPermissionRouterWithoutRecord := Router
	dataPermissionApi := api.ApiGroupApp.DataPermissionApi
	{
		// 受控表管理
		dataPermissionRouter.POST("controlledTable", dataPermissionApi.CreateControlledTable)   // 创建受控表
		dataPermissionRouter.PUT("controlledTable", dataPermissionApi.UpdateControlledTable)    // 更新受控表
		dataPermissionRouter.DELETE("controlledTable", dataPermissionApi.DeleteControlledTable) // 删除受控表
	}
	{
		// 查询接口（不记录操作日志）
		dataPermissionRouterWithoutRecord.GET("controlledTableList", dataPermissionApi.GetControlledTableList) // 获取受控表列表
		dataPermissionRouterWithoutRecord.GET("getConfig", dataPermissionApi.GetDataPermissionConfig)          // 获取数据权限配置
		dataPermissionRouterWithoutRecord.GET("tableList", dataPermissionApi.GetTableList)                     // 获取数据库表列表
		dataPermissionRouterWithoutRecord.GET("fieldList", dataPermissionApi.GetFieldList)                     // 获取表字段列表
		dataPermissionRouterWithoutRecord.GET("statistics", dataPermissionApi.GetStatistics)                   // 获取统计信息
	}
	{
		// 配置管理
		dataPermissionRouter.POST("saveConfig", dataPermissionApi.SaveDataPermissionConfig) // 保存数据权限配置
	}
}
