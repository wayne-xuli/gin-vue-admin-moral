package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/model"
	req "github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/model/request"
	resp "github.com/flipped-aurora/gin-vue-admin/server/plugin/datapermission/model/response"
	"gorm.io/gorm"
)

type DataPermissionService struct{}

// CreateControlledTable 创建受控表
func (s *DataPermissionService) CreateControlledTable(request req.CreateControlledTableRequest) error {
	// 检查表名是否已存在
	var count int64
	global.GVA_DB.Model(&model.ControlledTable{}).Where("table_name = ?", request.Table).Count(&count)
	if count > 0 {
		return errors.New("表名已存在")
	}

	// 创建受控表记录
	controlledTable := model.ControlledTable{
		Table:       request.Table,
		Description: request.Description,
		Enabled:     request.Enabled,
		DataScope:   request.DataScope,
		UserField:   request.UserField,
		DeptField:   request.DeptField,
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 创建受控表
		if err := tx.Create(&controlledTable).Error; err != nil {
			return err
		}

		// 如果有字段权限配置，批量创建
		if len(request.FieldPermissions) > 0 {
			for _, fieldPerm := range request.FieldPermissions {
				fieldPermission := model.RoleFieldPermission{
					ControlledTableID: controlledTable.ID,
					FieldName:         fieldPerm.FieldName,
					FieldChinese:      fieldPerm.FieldChinese,
					FieldDesc:         fieldPerm.FieldDesc,
					Visibility:        fieldPerm.Visible,
					EditPermission:    fieldPerm.Editable,
					Exportable:        fieldPerm.Exportable,
					Queryable:         fieldPerm.Queryable,
					Enabled:           true,
				}
				if err := tx.Create(&fieldPermission).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// UpdateControlledTable 更新受控表
func (s *DataPermissionService) UpdateControlledTable(request req.UpdateControlledTableRequest) error {
	// 检查表是否存在
	var controlledTable model.ControlledTable
	if err := global.GVA_DB.First(&controlledTable, request.ID).Error; err != nil {
		return errors.New("受控表不存在")
	}

	// 检查表名是否被其他记录使用
	var count int64
	global.GVA_DB.Model(&model.ControlledTable{}).Where("table_name = ? AND id != ?", request.Table, request.ID).Count(&count)
	if count > 0 {
		return errors.New("表名已被其他记录使用")
	}

	// 更新受控表
	controlledTable.Table = request.Table
	controlledTable.Description = request.Description
	controlledTable.Enabled = request.Enabled
	controlledTable.DataScope = request.DataScope
	controlledTable.UserField = request.UserField
	controlledTable.DeptField = request.DeptField

	return global.GVA_DB.Save(&controlledTable).Error
}

// DeleteControlledTable 删除受控表
func (s *DataPermissionService) DeleteControlledTable(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除相关的角色数据权限
		if err := tx.Where("controlled_table_id = ?", id).Unscoped().Delete(&model.RoleDataPermission{}).Error; err != nil {
			return err
		}

		// 删除相关的角色字段权限
		if err := tx.Where("controlled_table_id = ?", id).Unscoped().Delete(&model.RoleFieldPermission{}).Error; err != nil {
			return err
		}

		// 删除受控表
		return tx.Delete(&model.ControlledTable{}, id).Unscoped().Error
	})
}

// GetControlledTableList 获取受控表列表
func (s *DataPermissionService) GetControlledTableList(info req.ControlledTableSearch) (list []resp.ControlledTableResponse, total int64, err error) {

	db := global.GVA_DB.Model(&model.ControlledTable{})

	// 添加搜索条件
	if info.Table != "" {
		db = db.Where("table_name LIKE ?", "%"+info.Table+"%")
	}
	if info.Description != "" {
		db = db.Where("description LIKE ?", "%"+info.Description+"%")
	}
	if info.Enabled != nil {
		db = db.Where("enabled = ?", *info.Enabled)
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 获取列表数据
	var controlledTables []model.ControlledTable
	err = db.Order("created_at DESC").Find(&controlledTables).Error
	if err != nil {
		return
	}

	// 构建响应数据
	for _, table := range controlledTables {
		// 统计角色权限数量
		var rolePermCount int64
		global.GVA_DB.Model(&model.RoleDataPermission{}).Where("controlled_table_id = ?", table.ID).Count(&rolePermCount)

		// 统计字段权限数量
		var fieldPermCount int64
		global.GVA_DB.Model(&model.RoleFieldPermission{}).Where("controlled_table_id = ?", table.ID).Count(&fieldPermCount)

		list = append(list, resp.ControlledTableResponse{
			ControlledTable:      table,
			RolePermissionCount:  int(rolePermCount),
			FieldPermissionCount: int(fieldPermCount),
		})
	}

	return list, total, nil
}

// SaveDataPermissionConfig 保存数据权限配置
func (s *DataPermissionService) SaveDataPermissionConfig(request req.SaveDataPermissionConfigRequest) error {
	// 检查受控表是否存在
	var controlledTable model.ControlledTable
	if err := global.GVA_DB.Where("table_name = ?", request.Table).First(&controlledTable).Error; err != nil {
		return errors.New("受控表不存在，请先添加受控表")
	}
	// 更新受控表的用戶字段和部门字段
	controlledTable.UserField = request.UserField
	controlledTable.DeptField = request.DeptField
	if err := global.GVA_DB.Save(&controlledTable).Error; err != nil {
		return errors.New("受控表的用戶字段和部门字段失敗")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除现有的角色数据权限配置
		if err := tx.Where("authority_id = ? AND controlled_table_id = ?", request.AuthorityID, controlledTable.ID).Unscoped().Delete(&model.RoleDataPermission{}).Error; err != nil {
			return err
		}

		// 删除现有的角色字段权限配置
		if err := tx.Where("authority_id = ? AND controlled_table_id = ?", request.AuthorityID, controlledTable.ID).Unscoped().Delete(&model.RoleFieldPermission{}).Error; err != nil {
			return err
		}

		// 创建新的角色数据权限配置
		roleDataPermission := model.RoleDataPermission{
			AuthorityID:       request.AuthorityID,
			ControlledTableID: controlledTable.ID,
			DataScope:         request.DataScope,
			CustomCondition:   request.CustomSQL,
			Priority:          50,
			Enabled:           true,
		}
		if err := tx.Create(&roleDataPermission).Error; err != nil {
			return err
		}

		// 创建新的角色字段权限配置
		for _, field := range request.Fields {
			roleFieldPermission := model.RoleFieldPermission{
				AuthorityID:       request.AuthorityID,
				ControlledTableID: controlledTable.ID,
				FieldName:         field.FieldName,
				FieldChinese:      field.FieldChinese,
				FieldDesc:         field.FieldDesc,
				Visibility:        field.Visible,
				EditPermission:    field.Editable,
				Exportable:        field.Exportable,
				Queryable:         field.Queryable,
				Enabled:           true,
			}
			if err := tx.Create(&roleFieldPermission).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetDataPermissionConfig 获取数据权限配置
func (s *DataPermissionService) GetDataPermissionConfig(request req.GetDataPermissionConfigRequest) (resp.DataPermissionConfigResponse, error) {
	var response resp.DataPermissionConfigResponse

	// 获取受控表信息
	var controlledTable model.ControlledTable
	if err := global.GVA_DB.Where("table_name = ?", request.Table).First(&controlledTable).Error; err != nil {
		return response, errors.New("受控表不存在")
	}

	// 获取角色数据权限配置
	var roleDataPermission model.RoleDataPermission
	if err := global.GVA_DB.Preload("Authority").Where("authority_id = ? AND controlled_table_id = ?", request.AuthorityID, controlledTable.ID).First(&roleDataPermission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有配置，返回默认配置
			response.AuthorityID = request.AuthorityID
			response.Table = request.Table
			response.DataScope = "self"
			response.CustomSQL = ""
			response.Fields = []model.FieldPermissionItem{}
			return response, nil
		}
		return response, err
	}

	// 获取角色字段权限配置
	var roleFieldPermissions []model.RoleFieldPermission
	global.GVA_DB.Where("authority_id = ? AND controlled_table_id = ?", request.AuthorityID, controlledTable.ID).Find(&roleFieldPermissions)

	// 构建响应数据
	response.AuthorityID = request.AuthorityID
	response.AuthorityName = roleDataPermission.Authority.AuthorityName
	response.Table = request.Table
	response.DataScope = roleDataPermission.DataScope
	response.CustomSQL = roleDataPermission.CustomCondition

	for _, fieldPerm := range roleFieldPermissions {
		response.Fields = append(response.Fields, model.FieldPermissionItem{
			FieldName:    fieldPerm.FieldName,
			FieldChinese: fieldPerm.FieldChinese,
			FieldDesc:    fieldPerm.FieldDesc,
			Visible:      fieldPerm.Visibility,
			Editable:     fieldPerm.EditPermission,
			Exportable:   fieldPerm.Exportable,
			Queryable:    fieldPerm.Queryable,
		})
	}

	return response, nil
}

// GetTableList 获取数据库表列表
func (s *DataPermissionService) GetTableList() (resp.TableListResponse, error) {
	var response resp.TableListResponse

	// 获取数据库中的所有表
	var tables []string
	if err := global.GVA_DB.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'").Scan(&tables).Error; err != nil {
		return response, err
	}

	// 获取表注释
	for _, tableName := range tables {
		//var tableComment string
		// 查询表注释
		//global.GVA_DB.Raw("SELECT table_comment FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", tableName).Scan(&tableComment)

		response.Tables = append(response.Tables, model.TableInfo{
			Table:        tableName,
			TableComment: tableName,
		})
	}

	return response, nil
}

// GetFieldList 获取表字段列表
func (s *DataPermissionService) GetFieldList(tableName string) (resp.FieldListResponse, error) {
	var response resp.FieldListResponse

	// 获取表字段信息
	type ColumnInfo struct {
		ColumnName    string `json:"column_name"`
		DataType      string `json:"data_type"`
		ColumnComment string `json:"column_comment"`
		ColumnKey     string `json:"column_key"`
	}

	var columns []ColumnInfo
	if err := global.GVA_DB.Raw(`
		SELECT 
			"column_name",
			"data_type"
		FROM information_schema.columns 
		WHERE table_schema = 'public' AND table_name = ?
		ORDER BY ordinal_position
	`, tableName).Scan(&columns).Error; err != nil {
		return response, err
	}

	// 构建字段信息
	for _, column := range columns {
		defaultReadonly := column.ColumnKey == "PRI" || strings.Contains(column.ColumnName, "_id") ||
			column.ColumnName == "created_at" || column.ColumnName == "updated_at" || column.ColumnName == "deleted_at"

		response.Fields = append(response.Fields, model.FieldInfo{
			FieldName:       column.ColumnName,
			FieldChinese:    column.ColumnComment,
			FieldDesc:       fmt.Sprintf("%s %s", column.DataType, column.ColumnComment),
			DataType:        column.DataType,
			DefaultReadonly: defaultReadonly,
		})
	}

	return response, nil
}

// GetStatistics 获取统计信息
func (s *DataPermissionService) GetStatistics() (resp.StatisticsResponse, error) {
	var response resp.StatisticsResponse

	// 统计受控表数量
	var controlledTableCount int64
	global.GVA_DB.Model(&model.ControlledTable{}).Count(&controlledTableCount)
	response.ControlledTableCount = int(controlledTableCount)

	// 统计角色数据权限数量
	var roleDataPermissionCount int64
	global.GVA_DB.Model(&model.RoleDataPermission{}).Count(&roleDataPermissionCount)
	response.RoleDataPermissionCount = int(roleDataPermissionCount)

	// 统计角色字段权限数量
	var roleFieldPermissionCount int64
	global.GVA_DB.Model(&model.RoleFieldPermission{}).Count(&roleFieldPermissionCount)
	response.RoleFieldPermissionCount = int(roleFieldPermissionCount)

	// 统计活跃角色数量（有权限配置的角色）
	var activeAuthorityCount int64
	global.GVA_DB.Model(&model.RoleDataPermission{}).Distinct("authority_id").Count(&activeAuthorityCount)
	response.ActiveAuthorityCount = int(activeAuthorityCount)

	return response, nil
}
