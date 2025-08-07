package api

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DataPermissionApi struct{}

var dataPermissionService = service.ServiceGroupApp.DataPermissionService

// CreateControlledTable 创建受控表
// @Tags DataPermission
// @Summary 创建受控表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CreateControlledTableRequest true "创建受控表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /datapermission/controlledTable [post]
func (api *DataPermissionApi) CreateControlledTable(c *gin.Context) {
	var req request.CreateControlledTableRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = dataPermissionService.CreateControlledTable(req)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateControlledTable 更新受控表
// @Tags DataPermission
// @Summary 更新受控表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UpdateControlledTableRequest true "更新受控表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /datapermission/controlledTable [put]
func (api *DataPermissionApi) UpdateControlledTable(c *gin.Context) {
	var req request.UpdateControlledTableRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = dataPermissionService.UpdateControlledTable(req)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteControlledTable 删除受控表
// @Tags DataPermission
// @Summary 删除受控表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetByIDRequest true "删除受控表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /datapermission/controlledTable [delete]
func (api *DataPermissionApi) DeleteControlledTable(c *gin.Context) {
	var req request.GetByIDRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = dataPermissionService.DeleteControlledTable(req.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetControlledTableList 分页获取受控表列表
// @Tags DataPermission
// @Summary 分页获取受控表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.ControlledTableSearch true "分页获取受控表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /datapermission/controlledTableList [get]
func (api *DataPermissionApi) GetControlledTableList(c *gin.Context) {
	var pageInfo request.ControlledTableSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := dataPermissionService.GetControlledTableList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// SaveDataPermissionConfig 保存数据权限配置
// @Tags DataPermission
// @Summary 保存数据权限配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SaveDataPermissionConfigRequest true "保存数据权限配置"
// @Success 200 {object} response.Response{msg=string} "保存成功"
// @Router /datapermission/saveConfig [post]
func (api *DataPermissionApi) SaveDataPermissionConfig(c *gin.Context) {
	var req request.SaveDataPermissionConfigRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = dataPermissionService.SaveDataPermissionConfig(req)
	if err != nil {
		global.GVA_LOG.Error("保存失败!", zap.Error(err))
		response.FailWithMessage("保存失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("保存成功", c)
}

// GetDataPermissionConfig 获取数据权限配置
// @Tags DataPermission
// @Summary 获取数据权限配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param authorityId query string true "角色ID"
// @Param tableName query string true "表名"
// @Success 200 {object} response.Response{data=response.DataPermissionConfigResponse,msg=string} "获取成功"
// @Router /datapermission/getConfig [get]
func (api *DataPermissionApi) GetDataPermissionConfig(c *gin.Context) {
	authorityIdStr := c.Query("authorityId")
	tableName := c.Query("tableName")

	if authorityIdStr == "" || tableName == "" {
		response.FailWithMessage("参数不能为空", c)
		return
	}

	authorityId, err := strconv.ParseUint(authorityIdStr, 10, 32)
	if err != nil {
		response.FailWithMessage("角色ID格式错误", c)
		return
	}

	req := request.GetDataPermissionConfigRequest{
		AuthorityID: uint(authorityId),
		Table:       tableName,
	}

	config, err := dataPermissionService.GetDataPermissionConfig(req)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(config, "获取成功", c)
}

// GetTableList 获取数据库表列表
// @Tags DataPermission
// @Summary 获取数据库表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=response.TableListResponse,msg=string} "获取成功"
// @Router /datapermission/tableList [get]
func (api *DataPermissionApi) GetTableList(c *gin.Context) {
	tables, err := dataPermissionService.GetTableList()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(tables, "获取成功", c)
}

// GetFieldList 获取表字段列表
// @Tags DataPermission
// @Summary 获取表字段列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param tableName query string true "表名"
// @Success 200 {object} response.Response{data=response.FieldListResponse,msg=string} "获取成功"
// @Router /datapermission/fieldList [get]
func (api *DataPermissionApi) GetFieldList(c *gin.Context) {
	tableName := c.Query("tableName")
	if tableName == "" {
		response.FailWithMessage("表名不能为空", c)
		return
	}

	fields, err := dataPermissionService.GetFieldList(tableName)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(fields, "获取成功", c)
}

// GetStatistics 获取统计信息
// @Tags DataPermission
// @Summary 获取统计信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=response.StatisticsResponse,msg=string} "获取成功"
// @Router /datapermission/statistics [get]
func (api *DataPermissionApi) GetStatistics(c *gin.Context) {
	stats, err := dataPermissionService.GetStatistics()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(stats, "获取成功", c)
}
