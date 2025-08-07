<template>
  <div class="data-permission-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-main">
          <div class="title-section">
            <h3 class="page-title">
              <el-icon class="title-icon"><Lock /></el-icon>
              数据权限管理
            </h3>
            <p class="page-description">配置不同角色的数据访问权限，包括数据范围、字段权限等</p>
          </div>
          
          <div class="help-section">
            <el-collapse v-model="helpVisible" class="help-collapse">
              <el-collapse-item name="help">
                <template #title>
                  <div class="help-title">
                    <el-icon class="help-icon"><QuestionFilled /></el-icon>
                    <span>使用说明</span>
                  </div>
                </template>
                <div class="help-content-compact">
                  <div class="help-steps">
                    <div class="step-item" v-for="(step, index) in helpSteps" :key="index">
                      <span class="step-number">{{ index + 1 }}</span>
                      <span class="step-text">{{ step }}</span>
                    </div>
                  </div>
                </div>
              </el-collapse-item>
            </el-collapse>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="page-content">
      <div class="content-layout">
        <!-- 左侧角色树 -->
        <div class="left-panel">
          <RoleTree @role-selected="handleRoleSelected" ref="roleTreeRef" />
        </div>

        <!-- 右侧数据权限配置 -->
        <div class="right-panel">
          <div class="panel-container">
            <!-- 未选择角色时的提示 -->
            <div v-if="!selectedRole" class="no-selection">
              <el-empty description="请从左侧选择一个角色来配置数据权限">
                <el-icon class="empty-icon"><UserFilled /></el-icon>
              </el-empty>
            </div>

            <!-- 数据权限配置组件 -->
            <DataPermissionConfig 
              v-else 
              :current-role="selectedRole" 
              :key="selectedRole.authorityId"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Lock, UserFilled, QuestionFilled } from '@element-plus/icons-vue'
import RoleTree from '../components/RoleTree.vue'
import DataPermissionConfig from '../components/DataPermissionConfig.vue'

defineOptions({
  name: 'DataPermissionIndex'
})

// 响应式数据
const selectedRole = ref(null)
const roleTreeRef = ref(null)
const helpVisible = ref([])

// 帮助步骤数据
const helpSteps = ref([
  '从左侧角色树中选择要配置的角色',
  '添加或选择权限受控表',
  '配置表的用户字段和部门字段',
  '设置数据访问范围（全部、部门、个人等）',
  '配置字段级别的权限（可见性、编辑权限等）',
  '保存配置'
])

// 方法
const handleRoleSelected = (role) => {
  selectedRole.value = role
  console.log('选中角色:', role)
}

// 暴露方法
defineExpose({
  refreshRoles: () => {
    if (roleTreeRef.value) {
      roleTreeRef.value.refreshRoles()
    }
  },
  setSelectedRole: (role) => {
    selectedRole.value = role
    if (roleTreeRef.value) {
      roleTreeRef.value.setSelectedRole(role)
    }
  }
})
</script>

<style scoped>
.data-permission-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 20px;
}

.page-header {
  margin-bottom: 10px;
}

.header-content {
  background: white;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  border-left: 4px solid #409eff;
}

.header-main {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 30px;
}

.title-section {
  flex: 1;
}

.page-title {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  display: flex;
  align-items: center;
}

.title-icon {
  margin-right: 12px;
  color: #409eff;
  font-size: 28px;
}

.page-description {
  margin: 0;
  color: #606266;
  font-size: 14px;
  line-height: 1.5;
}

.help-section {
  width: 400px;
  flex-shrink: 0;
}

.help-collapse {
  border: none;
  background: transparent;
}

.help-collapse :deep(.el-collapse-item__header) {
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border: 1px solid #e1f5fe;
  border-radius: 8px;
  padding: 12px 16px;
  font-size: 14px;
  font-weight: 500;
  color: #0277bd;
  transition: all 0.3s ease;
}

.help-collapse :deep(.el-collapse-item__header:hover) {
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  border-color: #81d4fa;
}

.help-collapse :deep(.el-collapse-item__content) {
  padding: 0;
  border: none;
}

.help-collapse :deep(.el-collapse-item__wrap) {
  border: none;
}

.help-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.help-icon {
  color: #0277bd;
  font-size: 16px;
}

.help-content-compact {
  background: white;
  border: 1px solid #e1f5fe;
  border-top: none;
  border-radius: 0 0 8px 8px;
  padding: 16px;
}

.help-steps {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.step-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 6px 0;
}

.step-number {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  background: linear-gradient(135deg, #409eff 0%, #1976d2 100%);
  color: white;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
  margin-top: 1px;
}

.step-text {
  color: #606266;
  font-size: 13px;
  line-height: 1.4;
  flex: 1;
}

.page-content {
  margin-bottom: 20px;
}

.content-layout {
  display: flex;
  gap: 20px;
  height: calc(100vh - 200px);
  min-height: 600px;
}

.left-panel {
  width: 320px;
  flex-shrink: 0;
}

.right-panel {
  flex: 1;
  min-width: 0;
}

.panel-container {
  height: 100%;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.no-selection {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
}

.empty-icon {
  font-size: 48px;
  color: #409eff;
  margin-bottom: 16px;
}

.help-info {
  margin-top: 20px;
}

.help-content {
  line-height: 1.6;
}

.help-content p {
  margin: 0 0 8px 0;
}

.help-content ol {
  margin: 8px 0 0 0;
  padding-left: 20px;
}

.help-content li {
  margin: 4px 0;
  color: #606266;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .header-main {
    flex-direction: column;
    gap: 20px;
  }
  
  .help-section {
    width: 100%;
  }
  
  .content-layout {
    flex-direction: column;
    height: auto;
  }
  
  .left-panel {
    width: 100%;
    height: 300px;
  }
  
  .right-panel {
    height: 600px;
  }
}

@media (max-width: 768px) {
  .data-permission-page {
    padding: 10px;
  }
  
  .header-content {
    padding: 16px;
  }
  
  .header-main {
    gap: 15px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .title-icon {
    font-size: 24px;
  }
  
  .help-section {
    width: 100%;
  }
  
  .step-text {
    font-size: 12px;
  }
  
  .content-layout {
    gap: 10px;
  }
  
  .left-panel {
    height: 250px;
  }
  
  .right-panel {
    height: 500px;
  }
}

/* 动画效果 */
.panel-container {
  transition: all 0.3s ease;
}

.panel-container:hover {
  box-shadow: 0 6px 30px rgba(0, 0, 0, 0.15);
}

/* 滚动条样式 */
:deep(.el-scrollbar__wrap) {
  scrollbar-width: thin;
  scrollbar-color: #c1c1c1 #f1f1f1;
}

:deep(.el-scrollbar__wrap::-webkit-scrollbar) {
  width: 6px;
  height: 6px;
}

:deep(.el-scrollbar__wrap::-webkit-scrollbar-track) {
  background: #f1f1f1;
  border-radius: 3px;
}

:deep(.el-scrollbar__wrap::-webkit-scrollbar-thumb) {
  background: #c1c1c1;
  border-radius: 3px;
}

:deep(.el-scrollbar__wrap::-webkit-scrollbar-thumb:hover) {
  background: #a8a8a8;
}
</style>