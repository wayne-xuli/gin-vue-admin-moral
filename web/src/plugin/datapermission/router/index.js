const DataPermissionRouter = {
  path: '/datapermission',
  name: 'datapermission',
  component: () => import('@/view/layout/index.vue'),
  meta: {
    title: '数据权限管理',
    icon: 'lock',
    keepAlive: true
  },
  children: [
    {
      path: 'index',
      name: 'DataPermissionIndex',
      component: () => import('@/plugin/datapermission/view/index.vue'),
      meta: {
        title: '数据权限配置',
        icon: 'setting',
        keepAlive: true
      }
    }
  ]
}

export default DataPermissionRouter