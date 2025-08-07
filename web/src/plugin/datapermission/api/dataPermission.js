import service from '@/utils/request'

// 获取数据权限配置
export const getDataPermissionConfig = (data) => {
  return service({
    url: '/datapermission/getConfig?tableName=' + data.tableName + '&authorityId=' + data.authorityID,
    method: 'get',
  })
}

// 保存数据权限配置
export const saveDataPermissionConfig = (data) => {
  return service({
    url: '/datapermission/saveConfig',
    method: 'post',
    data
  })
}

// 获取数据库表列表
export const getDatabaseTables = () => {
  return service({
    url: '/datapermission/tableList',
    method: 'get'
  })
}

// 获取表字段信息
export const getTableColumns = (tableName) => {
  return service({
    url: '/datapermission/fieldList',
    method: 'get',
    params: { tableName }
  })
}

// 获取受控表列表
export const getControlledTables = () => {
  return service({
    url: '/datapermission/controlledTableList',
    method: 'get'
  })
}

// 创建受控表
export const createControlledTable = (data) => {
  return service({
    url: '/datapermission/controlledTable',
    method: 'post',
    data
  })
}

// 删除受控表
export const deleteControlledTable = (data) => {
  return service({
    url: '/datapermission/controlledTable',
    method: 'delete',
    data
  })
}