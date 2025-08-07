# 数据权限管理插件

## 功能介绍

数据权限管理插件是基于 gin-vue-admin 框架开发的权限管理扩展，提供了细粒度的数据访问控制功能。

### 主要特性

- **角色树管理**: 左侧显示完整的角色层级结构，支持搜索和筛选
- **数据范围控制**: 支持全部数据、部门数据、个人数据等多种数据范围
- **字段级权限**: 可配置字段的可见性、编辑权限、导出权限、查询权限
- **受控表管理**: 动态添加需要进行权限控制的数据库表
- **自定义SQL**: 支持自定义SQL条件进行更灵活的数据过滤

## 目录结构

```
datapermission/
├── api/                    # API接口
│   └── dataPermission.js   # 数据权限相关接口
├── components/             # 组件
│   ├── DataPermissionConfig.vue  # 数据权限配置组件
│   └── RoleTree.vue       # 角色树组件
├── view/                   # 页面
│   └── index.vue          # 主页面
├── router/                 # 路由配置
│   └── index.js           # 路由定义
├── index.js               # 插件入口文件
└── README.md              # 说明文档
```

## 组件说明

### RoleTree 角色树组件

**功能**: 显示角色层级结构，支持角色选择和搜索

**Props**: 无

**Events**:
- `role-selected`: 角色选择事件，参数为选中的角色对象

**Methods**:
- `refreshRoles()`: 刷新角色数据
- `getSelectedRole()`: 获取当前选中的角色
- `setSelectedRole(role)`: 设置选中的角色

### DataPermissionConfig 数据权限配置组件

**功能**: 配置指定角色的数据权限

**Props**:
- `currentRole`: 当前选中的角色对象

**主要功能**:
- 权限受控表管理
- 数据范围配置
- 字段权限设置
- 配置保存

## API 接口

### 数据权限配置相关

- `getDataPermissionConfig(data)`: 获取数据权限配置
- `saveDataPermissionConfig(data)`: 保存数据权限配置
- `getDataScopeOptions()`: 获取数据范围选项

### 数据库表相关

- `getDatabaseTables()`: 获取数据库表列表
- `getTableColumns(tableName)`: 获取表字段信息
- `getControlledTables()`: 获取受控表列表
- `createControlledTable(data)`: 创建受控表

## 使用方法

### 1. 路由配置

在主应用的路由配置中引入插件路由：

```javascript
import DataPermissionRouter from '@/plugin/datapermission/router/index'

// 添加到路由配置中
const routes = [
  // ... 其他路由
  DataPermissionRouter
]
```

### 2. 组件使用

```vue
<template>
  <div>
    <!-- 使用完整的数据权限管理页面 -->
    <DataPermissionIndex />
    
    <!-- 或者单独使用组件 -->
    <div class="layout">
      <div class="left">
        <RoleTree @role-selected="handleRoleSelected" />
      </div>
      <div class="right">
        <DataPermissionConfig :current-role="selectedRole" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { DataPermissionIndex, RoleTree, DataPermissionConfig } from '@/plugin/datapermission'

const selectedRole = ref(null)

const handleRoleSelected = (role) => {
  selectedRole.value = role
}
</script>
```

### 3. API 使用

```javascript
import { dataPermissionApi } from '@/plugin/datapermission'

// 获取数据权限配置
const config = await dataPermissionApi.getDataPermissionConfig({
  authorityId: 1
})

// 保存配置
await dataPermissionApi.saveDataPermissionConfig({
  authorityId: 1,
  dataScope: 'dept',
  fieldList: [...]
})
```

## 配置说明

### 数据范围类型

- `all`: 全部数据权限
- `custom`: 自定义数据权限（需要配置SQL条件）
- `dept`: 本部门数据权限
- `dept_and_child`: 本部门及以下数据权限
- `self`: 仅本人数据权限

### 字段权限配置

每个字段可以配置以下权限：

- **可见性**: `visible`(可见) | `hidden`(隐藏) | `readonly`(只读)
- **编辑权限**: `editable`(可编辑) | `readonly`(只读) | `disabled`(禁用)
- **导出权限**: `true`(可导出) | `false`(禁止导出)
- **查询权限**: `true`(可查询) | `false`(禁止查询)

## 注意事项

1. 使用前需要确保后端已实现相应的API接口
2. 数据权限的实际生效需要后端配合实现数据过滤逻辑
3. 建议在生产环境中对敏感操作添加二次确认
4. 定期备份权限配置数据

## 版本历史

### v1.0.0
- 初始版本发布
- 实现基础的数据权限管理功能
- 支持角色树、权限配置、字段权限等核心功能