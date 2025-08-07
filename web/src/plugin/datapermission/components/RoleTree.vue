<template>
  <div class="role-tree">
    <div class="tree-header">
      <h3 class="tree-title">角色列表</h3>
      <el-button type="primary" size="small" @click="refreshRoles">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>
    
    <div class="tree-search">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索角色名称"
        prefix-icon="Search"
        clearable
        @input="handleSearch"
      />
    </div>

    <div class="tree-container">
      <el-tree
        ref="treeRef"
        :data="filteredRoleData"
        :props="treeProps"
        :highlight-current="true"
        :expand-on-click-node="false"
        :default-expand-all="true"
        node-key="authorityId"
        @node-click="handleNodeClick"
        class="role-tree-component"
      >
        <template #default="{ node, data }">
          <div class="tree-node">
            <div class="node-content">
              <el-icon class="node-icon">
                <User v-if="!data.children || data.children.length === 0" />
                <UserFilled v-else />
              </el-icon>
              <span class="node-label">{{ data.authorityName }}</span>
            </div>
            <div class="node-info">
              <el-tag size="small" type="info">ID: {{ data.authorityId }}</el-tag>
            </div>
          </div>
        </template>
      </el-tree>
    </div>

    <div v-if="filteredRoleData.length === 0" class="empty-state">
      <el-empty description="暂无角色数据">
        <el-button type="primary" @click="refreshRoles">重新加载</el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, User, UserFilled } from '@element-plus/icons-vue'
import { getAuthorityList } from '@/api/authority'

defineOptions({
  name: 'RoleTree'
})

const emit = defineEmits(['role-selected'])

// 响应式数据
const treeRef = ref(null)
const roleData = ref([])
const searchKeyword = ref('')
const selectedRole = ref(null)

// 树形组件配置
const treeProps = {
  children: 'children',
  label: 'authorityName',
  value: 'authorityId'
}

// 计算属性 - 过滤后的角色数据
const filteredRoleData = computed(() => {
  if (!searchKeyword.value) {
    return roleData.value
  }
  
  const filterTree = (nodes) => {
    return nodes.filter(node => {
      const matchesSearch = node.authorityName.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
                           node.authorityId.toString().includes(searchKeyword.value)
      
      if (node.children && node.children.length > 0) {
        const filteredChildren = filterTree(node.children)
        if (filteredChildren.length > 0) {
          return {
            ...node,
            children: filteredChildren
          }
        }
      }
      
      return matchesSearch
    }).map(node => {
      if (node.children && node.children.length > 0) {
        return {
          ...node,
          children: filterTree(node.children)
        }
      }
      return node
    })
  }
  
  return filterTree(roleData.value)
})

// 初始化
onMounted(() => {
  loadRoleData()
})

// 方法
const loadRoleData = async () => {
  try {
    const res = await getAuthorityList()
    if (res.code === 0) {
      roleData.value = res.data || []
    } else {
      ElMessage.error(res.msg || '获取角色列表失败')
    }
  } catch (error) {
    console.error('加载角色数据失败:', error)
    ElMessage.error('加载角色数据失败')
  }
}

const refreshRoles = () => {
  loadRoleData()
}

const handleSearch = () => {
  // 搜索时展开所有节点以便查看结果
  if (searchKeyword.value && treeRef.value) {
    setTimeout(() => {
      treeRef.value.setExpandedKeys(getAllNodeKeys(filteredRoleData.value))
    }, 100)
  }
}

const getAllNodeKeys = (nodes) => {
  let keys = []
  nodes.forEach(node => {
    keys.push(node.authorityId)
    if (node.children && node.children.length > 0) {
      keys = keys.concat(getAllNodeKeys(node.children))
    }
  })
  return keys
}

const handleNodeClick = (data, node) => {
  selectedRole.value = data
  emit('role-selected', data)
}

// 暴露方法给父组件
defineExpose({
  refreshRoles,
  getSelectedRole: () => selectedRole.value,
  setSelectedRole: (role) => {
    selectedRole.value = role
    if (treeRef.value && role) {
      treeRef.value.setCurrentKey(role.authorityId)
    }
  }
})
</script>

<style scoped>
.role-tree {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.tree-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #ebeef5;
}

.tree-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.tree-search {
  padding: 15px 20px;
  border-bottom: 1px solid #ebeef5;
}

.tree-container {
  flex: 1;
  padding: 10px;
  overflow-y: auto;
}

.role-tree-component {
  background: transparent;
}

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 5px 0;
}

.node-content {
  display: flex;
  align-items: center;
  flex: 1;
}

.node-icon {
  margin-right: 8px;
  color: #409eff;
}

.node-label {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

.node-info {
  margin-left: 10px;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

/* 自定义树形组件样式 */
:deep(.el-tree-node__content) {
  padding: 8px 5px;
  border-radius: 6px;
  margin: 2px 0;
  transition: all 0.3s ease;
}

:deep(.el-tree-node__content:hover) {
  background-color: #f5f7fa;
}

:deep(.el-tree-node.is-current > .el-tree-node__content) {
  background-color: #e6f7ff;
  border: 1px solid #409eff;
}

:deep(.el-tree-node__expand-icon) {
  color: #409eff;
}

:deep(.el-tree-node__expand-icon.is-leaf) {
  color: transparent;
}

/* 滚动条样式 */
.tree-container::-webkit-scrollbar {
  width: 6px;
}

.tree-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.tree-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.tree-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>