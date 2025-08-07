// 数据权限管理插件入口文件

// 导入组件
import DataPermissionIndex from './view/index.vue'
import DataPermissionConfig from './components/DataPermissionConfig.vue'
import RoleTree from './components/RoleTree.vue'

// 导入API
import * as dataPermissionApi from './api/dataPermission'

// 导入路由
import DataPermissionRouter from './router/index'

// 插件信息
const pluginInfo = {
  name: 'datapermission',
  version: '1.0.0',
  description: '数据权限管理插件',
  author: 'GVA Team',
  components: {
    DataPermissionIndex,
    DataPermissionConfig,
    RoleTree
  },
  api: dataPermissionApi,
  router: DataPermissionRouter
}

// 安装插件的方法
const install = (app) => {
  // 注册全局组件
  app.component('DataPermissionIndex', DataPermissionIndex)
  app.component('DataPermissionConfig', DataPermissionConfig)
  app.component('RoleTree', RoleTree)
}

// 导出插件
export default {
  install,
  ...pluginInfo
}

// 单独导出组件
export {
  DataPermissionIndex,
  DataPermissionConfig,
  RoleTree,
  dataPermissionApi,
  DataPermissionRouter
}