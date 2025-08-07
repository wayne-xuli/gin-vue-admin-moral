<template>
  <div class="data-permission-config">
    <!-- 数据权限管理标题栏 -->
    <div class="permission-header">
      <div class="header-actions">
        <span class="selected-role-info">当前角色：{{ currentRole?.authorityName || '未选择' }}</span>
        <el-button type="primary" @click="saveConfig" :disabled="!currentRole">保存配置</el-button>
      </div>
    </div>
    
    <!-- 主要内容区域 -->
    <div class="main-content">

    <!-- 权限模块 -->
    <div class="permission-modules" v-if="currentRole">
      <div class="module-header">
        <h3 class="section-title">权限受控表</h3>
        <el-button v-if="moduleOptions.length > 0" type="primary" size="small" @click="addModule">添加权限受控表</el-button>
      </div>
      <div class="module-checkboxes">
        <div v-if="moduleOptions.length === 0" class="no-data-tip">
          <el-empty description="暂无权限受控表">
            <el-button type="primary" @click="addModule">添加权限受控表</el-button>
          </el-empty>
        </div>
        <div v-else class="module-cards">
          <div 
            v-for="module in moduleOptions" 
            :key="module.value" 
            class="module-card"
            :class="{ 'active': selectedModule === module.value }"
            @click="selectedModule = module.value"
          >
            <div class="module-info">
              <div class="module-name">{{ module.label }}</div>
              <div class="module-table">{{ module.value }}</div>
            </div>
            <div class="module-actions">
              <el-button 
                type="danger" 
                size="small" 
                :icon="Delete" 
                @click.stop="deleteModule(module.value, module.ID)"
                title="删除权限受控表"
              >
                删除
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 权限字段设置 -->
    <div v-if="selectedModule && currentRole" class="table-field-config">
      <div class="section-header">
        <h3 class="section-title">权限字段设置</h3>
        <el-tooltip :content="helpTexts.dataScope" placement="top">
            <el-icon class="help-icon"><QuestionFilled /></el-icon>
          </el-tooltip>
      </div>
      <div class="field-config-form">
        <el-form :model="tableFieldConfig" label-width="70px" label-position="left" class="config-form">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="用户字段">
                <el-select 
                  v-model="tableFieldConfig.userField" 
                  placeholder="请选择用户字段"
                  filterable
                  clearable
                  allow-create
                  style="width: 100%"
                >
                  <el-option 
                    v-for="field in userFieldOptions" 
                    :key="field.value" 
                    :label="field.label" 
                    :value="field.value"
                  />
                </el-select>
              </el-form-item>
              <div class="field-tip">用于标识数据创建者，默认为 created_by</div>
            </el-col>
            <el-col :span="12">
              <el-form-item label="部门字段">
                <el-select 
                  v-model="tableFieldConfig.deptField" 
                  placeholder="请选择部门字段"
                  filterable
                  clearable
                  allow-create
                  style="width: 100%"
                >
                  <el-option 
                    v-for="field in deptFieldOptions" 
                    :key="field.value" 
                    :label="field.label" 
                    :value="field.value"
                  />
                </el-select>
              </el-form-item>
              <div class="field-tip">用于标识数据所属部门，默认为 dept_code</div>
            </el-col>
          </el-row>
        </el-form>
      </div>
    </div>

    <!-- 数据范围 -->
    <div v-if="currentRole" class="data-scope">
      <div class="section-header">
        <h3 class="section-title">数据范围</h3>
        <el-tooltip :content="helpTexts.dataScope" placement="top">
          <el-icon class="help-icon"><QuestionFilled /></el-icon>
        </el-tooltip>
      </div>
      <div class="scope-options">
        <el-radio-group v-model="dataScope">
          <el-radio 
            v-for="item in dataScopeOptions" 
            :key="item.value" 
            :label="item.value"
          >
            {{ item.label }}
          </el-radio>
        </el-radio-group>
      </div>
      
      <!-- 自定义SQL输入 -->
      <div v-if="dataScope === 'custom'" class="custom-sql">
        <h4>自定义SQL条件</h4>
        <el-input
          v-model="customSql"
          type="textarea"
          :rows="4"
          placeholder="请输入自定义SQL条件，例如：dept_id IN (1,2,3)"
        />
      </div>
    </div>

    <!-- 字段权限配置 -->
    <div v-if="selectedModule && currentRole" class="field-permissions">
      <div class="field-header">
        <div class="section-header">
          <h3 class="section-title">字段权限配置</h3>
          <el-tooltip :content="helpTexts.fieldPermission" placement="top">
            <el-icon class="help-icon"><QuestionFilled /></el-icon>
          </el-tooltip>
        </div>
        <div class="field-actions">
          <div class="field-operations">
            <el-button type="primary" size="small" @click="addField" :disabled="!selectedModule">添加字段</el-button>
            <el-button type="success" size="small" @click="refreshFields" :disabled="!selectedModule">刷新字段</el-button>
          </div>
        </div>
      </div>
      
      <el-table 
        :data="fieldList" 
        border
        style="width: 100%" 
        :empty-text="selectedModule ? '暂无字段数据' : '请先选择权限受控表'"
        class="field-table"
        height="500"
        :show-overflow-tooltip="true"
      >
        <el-table-column prop="fieldName" label="字段名称" min-width="50" show-overflow-tooltip />
        <el-table-column prop="fieldChinese" label="字段中文" min-width="60" show-overflow-tooltip />
        <el-table-column prop="fieldDesc" label="字段描述" min-width="80" show-overflow-tooltip />
        <el-table-column label="加密" width="80" align="center">
          <template #default="scope">
            <el-switch 
              v-model="scope.row.encrypted" 
              size="small"
              inline-prompt
              active-text="加密"
              inactive-text="明文"
              active-color="#f56c6c"
              @change="handleEncryptionChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="可见性" width="90" align="center">
          <template #default="scope">
            <el-switch 
              v-model="scope.row.visible" 
              size="small"
              inline-prompt
              active-text="可见"
              inactive-text="隐藏"
              :disabled="scope.row.encrypted"
              @change="handleVisibilityChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="编辑" width="80" align="center">
          <template #default="scope">
            <el-switch 
              v-model="scope.row.editable" 
              size="small"
              inline-prompt
              active-text="可编辑"
              inactive-text="只读"
              :disabled="scope.row.encrypted || !scope.row.visible"
            />
          </template>
        </el-table-column>
        <el-table-column label="导出" width="80" align="center">
          <template #default="scope">
            <el-switch 
              v-model="scope.row.exportable" 
              size="small"
              inline-prompt
              active-text="可导出"
              inactive-text="禁导出"
              :disabled="scope.row.encrypted || !scope.row.visible"
            />
          </template>
        </el-table-column>
        <el-table-column label="查询" width="80" align="center">
          <template #default="scope">
            <el-switch 
              v-model="scope.row.queryable" 
              size="small"
              inline-prompt
              active-text="可查询"
              inactive-text="禁查询"
              :disabled="scope.row.encrypted || !scope.row.visible"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="editField(scope.row)">编辑</el-button>
            <el-button type="danger" size="small" @click="deleteField(scope.$index)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 添加模块对话框 -->
    <el-dialog 
      v-model="moduleDialogVisible" 
      title="选择数据库表作为权限模块"
      width="700px"
    >
      <div class="table-selection">
        <el-input 
          v-model="tableSearchKeyword" 
          placeholder="搜索表名或中文名称"
          prefix-icon="Search"
          clearable
          style="margin-bottom: 15px;"
        />
        <el-table 
          :data="filteredTables" 
          height="400"
          @selection-change="handleTableSelection"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="tableName" label="表名" min-width="150" show-overflow-tooltip />
          <el-table-column prop="tableComment" label="中文名称" min-width="150" show-overflow-tooltip />
        </el-table>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="moduleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveSelectedTables" :disabled="selectedTables.length === 0">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 添加字段对话框 -->
    <el-dialog 
      v-model="showFieldDialog" 
      title="添加字段" 
      width="800px" 
      :close-on-click-modal="false"
    >
      <div class="field-dialog-content">
        <div class="dialog-header">
          <el-input 
            v-model="fieldSearchKeyword" 
            placeholder="搜索字段名称或描述" 
            prefix-icon="Search" 
            clearable 
            style="width: 300px; margin-bottom: 16px;"
          />
          <span class="field-count">可用字段：{{ filteredAvailableFields.length }} 个</span>
        </div>
        
        <el-table 
          ref="fieldSelectionTable"
          :data="filteredAvailableFields" 
          border
          height="400px"
          @selection-change="handleFieldSelection"
          :empty-text="availableFields.length === 0 ? '正在加载字段...' : '没有找到匹配的字段'"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="fieldName" label="字段名称" min-width="80" show-overflow-tooltip />
          <el-table-column prop="fieldChinese" label="字段中文" min-width="80" show-overflow-tooltip />
          <el-table-column prop="dataType" label="字段类型" width="70" />
          <el-table-column prop="fieldDesc" label="字段描述" min-width="150" show-overflow-tooltip />
          <el-table-column prop="isRequired" label="必填" width="60" align="center">
            <template #default="scope">
              <el-tag :type="scope.row.isRequired ? 'danger' : 'info'" size="small">
                {{ scope.row.isRequired ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showFieldDialog = false">取消</el-button>
          <el-button type="primary" @click="saveField" :disabled="selectedFields.length === 0">
            添加选中字段 ({{ selectedFields.length }})
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 编辑字段对话框 -->
    <el-dialog 
      v-model="showEditFieldDialog" 
      title="编辑字段权限" 
      width="600px" 
      :close-on-click-modal="false"
    >
      <el-form :model="currentField" label-width="100px" v-if="currentField">
        <el-form-item label="字段名称">
          <el-input v-model="currentField.fieldName" disabled />
        </el-form-item>
        <el-form-item label="字段中文">
          <el-input v-model="currentField.fieldChinese" />
        </el-form-item>
        <el-form-item label="字段描述">
          <el-input v-model="currentField.fieldDesc" type="textarea" :rows="2" />
        </el-form-item>
        <div class="permission-levels">
          <div class="permission-level level-1">
            <el-form-item label="加密权限">
              <el-switch 
                v-model="currentField.encrypted" 
                active-text="加密" 
                inactive-text="明文"
                active-color="#f56c6c"
                @change="handleEditFieldEncryptionChange"
              />
              <div class="permission-tip">最高级别保护，加密字段对用户完全不可见</div>
            </el-form-item>
          </div>
          
          <div class="permission-level level-2">
            <el-form-item label="可见性">
              <el-switch 
                v-model="currentField.visible" 
                active-text="可见" 
                inactive-text="隐藏"
                :disabled="currentField.encrypted"
                @change="handleEditFieldVisibilityChange"
              />
              <div class="permission-tip">基础权限，控制字段是否对用户显示</div>
            </el-form-item>
          </div>
          
          <div class="permission-level level-3">
            <div class="function-permissions">
              <el-form-item label="编辑权限">
                <el-switch 
                  v-model="currentField.editable" 
                  active-text="可编辑" 
                  inactive-text="只读"
                  :disabled="currentField.encrypted || !currentField.visible"
                />
              </el-form-item>
              <el-form-item label="导出权限">
                <el-switch 
                  v-model="currentField.exportable" 
                  active-text="允许导出" 
                  inactive-text="禁止导出"
                  :disabled="currentField.encrypted || !currentField.visible"
                />
              </el-form-item>
              <el-form-item label="查询权限">
                <el-switch 
                  v-model="currentField.queryable" 
                  active-text="允许查询" 
                  inactive-text="禁止查询"
                  :disabled="currentField.encrypted || !currentField.visible"
                />
              </el-form-item>
            </div>
            <div class="permission-tip">功能权限，依赖于可见性权限</div>
          </div>
        </div>
        
        <!-- 权限预设模板 -->
        <div class="permission-templates">
          <el-form-item label="权限模板">
            <el-button-group>
              <el-button size="small" @click="applyTemplate('full_access')">完全开放</el-button>
              <el-button size="small" @click="applyTemplate('readonly')">只读模式</el-button>
              <el-button size="small" @click="applyTemplate('limited_edit')">受限编辑</el-button>
              <el-button size="small" @click="applyTemplate('encrypted')" type="danger">完全加密</el-button>
            </el-button-group>
          </el-form-item>
        </div>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showEditFieldDialog = false">取消</el-button>
          <el-button type="primary" @click="updateField">保存</el-button>
        </div>
      </template>
    </el-dialog>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { QuestionFilled, ArrowDown, Delete } from '@element-plus/icons-vue'
import {
  getDataPermissionConfig,
  saveDataPermissionConfig,
  getDatabaseTables,
  getTableColumns,
  getControlledTables,
  createControlledTable,
  deleteControlledTable
} from '../api/dataPermission'
import { getDict } from '@/utils/dictionary'

defineOptions({
  name: 'DataPermissionConfig'
})

const props = defineProps({
  currentRole: {
    type: Object,
    default: null
  }
})

// 响应式数据
const selectedModule = ref('')
const moduleOptions = ref([])
const dataScope = ref('')
const dataScopeOptions = ref([])
const customSql = ref('')
const fieldList = ref([])
const availableFields = ref([])
const selectedFields = ref([])
const showFieldDialog = ref(false)
const showEditFieldDialog = ref(false)
const fieldSearchKeyword = ref('')
// 字段对话框相关变量已在上面定义
const moduleDialogVisible = ref(false)
const databaseTables = ref([])
const selectedTables = ref([])
const tableSearchKeyword = ref('')

// 受控表字段配置
const tableFieldConfig = reactive({
  userField: 'created_by',
  deptField: 'dept_code'
})

// 用户字段选项
const userFieldOptions = ref([])

// 部门字段选项
const deptFieldOptions = ref([])

// 说明信息
const helpTexts = ref({
  permissionField: '权限字段作为数据范围的过滤条件，可以在字典中配置：用户字段、部门字段',
  dataScope: '数据范围用于控制用户可以访问的数据范围，可以在字典中配置：数据范围',
  fieldPermission: '字段权限用于控制用户对特定字段的访问和操作权限'
})

// 当前编辑的字段
const currentField = ref(null)
const editingFieldIndex = ref(-1)

// 计算属性
const filteredTables = computed(() => {
  if (!tableSearchKeyword.value) return databaseTables.value
  return databaseTables.value.filter(table => 
    table.tableName.toLowerCase().includes(tableSearchKeyword.value.toLowerCase()) ||
    (table.tableComment && table.tableComment.includes(tableSearchKeyword.value))
  )
})

const filteredAvailableFields = computed(() => {
  if (!fieldSearchKeyword.value) return availableFields.value
  return availableFields.value.filter(field => 
    field.fieldName.toLowerCase().includes(fieldSearchKeyword.value.toLowerCase()) ||
    (field.fieldChinese && field.fieldChinese.includes(fieldSearchKeyword.value)) ||
    (field.fieldDesc && field.fieldDesc.includes(fieldSearchKeyword.value))
  )
})

// 方法
const loadDataPermissionConfig = async () => {
  if (!props.currentRole) return
  console.log('加载数据权限配置', selectedModule.value, props.currentRole)
  try {
    const res = await getDataPermissionConfig({ authorityID: props.currentRole.authorityId, tableName: selectedModule.value })
    if (res.code === 0 && res.data) {
      const config = res.data
      selectedModule.value = config.table || ''
      dataScope.value = config.dataScope || ''
      customSql.value = config.customSql || ''
      // 转换字段数据格式，兼容新的五级权限体系
      fieldList.value = (config.fields || []).map(field => ({
        ...field,
        encrypted: field.encrypted === true,
        visible: field.visibility === 'visible' || field.visible === true,
        editable: field.editPermission === 'editable' || field.editable === true,
        exportable: field.exportable !== false,
        queryable: field.queryable !== false
      }))
      Object.assign(tableFieldConfig, config.tableFieldConfig || {})
    }
  } catch (error) {
    console.error('加载数据权限配置失败:', error)
  }
}

const resetConfig = () => {
  selectedModule.value = ''
  dataScope.value = ''
  customSql.value = ''
  fieldList.value = []
  Object.assign(tableFieldConfig, {
    userField: 'created_by',
    deptField: 'dept_code'
  })
}

const loadDatabaseTables = async () => {
  try {
    const res = await getDatabaseTables()
    if (res.code === 0) {
      databaseTables.value = res.data.tables || []
    }
  } catch (error) {
    console.error('加载数据库表失败:', error)
  }
}

const loadControlledTables = async () => {
  try {
    const res = await getControlledTables()
    if (res.code === 0) {
      moduleOptions.value = (res.data.list || []).map(table => ({
        label: table.description || table.tableName,
        value: table.tableName,
        ID: table.ID,
        description: table.description,
        userField: table.userField,
        deptField: table.deptField
      }))
      if (moduleOptions.value.length > 0) {
        selectedModule.value = moduleOptions.value[0].value
        await loadDataPermissionConfig()
      }
    }
  } catch (error) {
    console.error('加载受控表失败:', error)
  }
}

const saveConfig = async () => {
  if (!props.currentRole) {
    ElMessage.warning('请先选择角色')
    return
  }

  try {
    // 转换字段数据格式，保持向后兼容
    const convertedFieldList = fieldList.value.map(field => ({
      ...field,
      encrypted: field.encrypted || false,
      visibility: field.visible ? 'visible' : 'hidden',
      editPermission: field.editable ? 'editable' : 'readonly',
      exportable: field.exportable || false,
      queryable: field.queryable || false
    }))

    const config = {
      authorityId: props.currentRole.authorityId,
      table: selectedModule.value,
      dataScope: dataScope.value,
      customSql: customSql.value,
      fields: convertedFieldList,
       ...tableFieldConfig
    }

    const res = await saveDataPermissionConfig(config)
    if (res.code === 0) {
      ElMessage.success('保存成功')
      // 重新加载受控权限表
      await loadControlledTables()
    } else {
      ElMessage.error(res.msg || '保存失败')
    }
  } catch (error) {
    console.error('保存配置失败:', error)
    ElMessage.error('保存失败')
  }
}

const addModule = () => {
  moduleDialogVisible.value = true
  selectedTables.value = []
}

const deleteModule = (tableName, ID) => {
  ElMessageBox.confirm(
    `确定要删除权限受控表「${tableName}」吗？删除后该表的所有权限配置将被清除。`,
    '删除权限受控表',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
      dangerouslyUseHTMLString: false
    }
  ).then(async () => {
    try {
      // 如果当前选中的就是要删除的表，先清空选择
      if (selectedModule.value === tableName) {
        selectedModule.value = ''
        fieldList.value = []
        dataScope.value = ''
        customSql.value = ''
      }
      
      // 从模块选项中移除
      const index = moduleOptions.value.findIndex(module => module.value === tableName)
      if (index > -1) {
        moduleOptions.value.splice(index, 1)
      }
      
      // 如果还有其他受控表，自动选择第一个
      if (moduleOptions.value.length > 0 && !selectedModule.value) {
        selectedModule.value = moduleOptions.value[0].value
      }
      
      const res = await deleteControlledTable({ id: ID })
      if (res.code == 0) {
        ElMessage.success('删除成功')
      } else {
        ElMessage.error(res.msg || '删除失败')
      }
      // 重新加载配置
      if (selectedModule.value) {
        await loadDataPermissionConfig()
      }
    } catch (error) {
      console.error('删除受控表失败:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    // 用户取消删除
  })
}

const handleTableSelection = (selection) => {
  selectedTables.value = selection
}

const saveSelectedTables = async () => {
  if (selectedTables.value.length === 0) {
    ElMessage.warning('请选择至少一个表')
    return
  }

  try {
    for (const table of selectedTables.value) {
      await createControlledTable({
        table: table.tableName,
        tableComment: table.tableComment
      })
    }
    ElMessage.success('添加成功')
    moduleDialogVisible.value = false
    await loadControlledTables()
  } catch (error) {
    console.error('添加受控表失败:', error)
    ElMessage.error('添加失败')
  }
}

const addField = async () => {
  if (!selectedModule.value) {
    ElMessage.warning('请先选择权限受控表')
    return
  }
  
  try {
    const res = await getTableColumns(selectedModule.value)
    if (res.code == 0) {
      availableFields.value = (res.data.fields || []).filter(field => 
        !fieldList.value.some(existingField => existingField.fieldName === field.fieldName)
      )
      selectedFields.value = []
      fieldSearchKeyword.value = ''
      showFieldDialog.value = true
    }
  } catch (error) {
    console.error('获取表字段失败:', error)
    ElMessage.error('获取表字段失败')
  }
}

const editField = (row) => {
  const index = fieldList.value.findIndex(field => field.fieldName === row.fieldName)
  editingFieldIndex.value = index
  currentField.value = { ...row }
  showEditFieldDialog.value = true
}

const saveField = () => {
  if (selectedFields.value.length === 0) {
    ElMessage.warning('请选择要添加的字段')
    return
  }

  selectedFields.value.forEach(field => {
    fieldList.value.push({
      fieldName: field.fieldName,
      fieldChinese: field.fieldChinese || field.fieldName,
      fieldDesc: field.fieldDesc || '',
      fieldType: field.fieldType || '',
      encrypted: false,
      visible: true,
      editable: true,
      exportable: true,
      queryable: true
    })
  })
  
  showFieldDialog.value = false
  ElMessage.success(`成功添加 ${selectedFields.value.length} 个字段`)
}

const updateField = () => {
  if (!currentField.value) return
  
  fieldList.value[editingFieldIndex.value] = { ...currentField.value }
  showEditFieldDialog.value = false
  ElMessage.success('字段更新成功')
}

const refreshFields = async () => {
  if (!selectedModule.value) {
    ElMessage.warning('请先选择权限受控表')
    return
  }
  
  try {
    const res = await getTableColumns(selectedModule.value)
    if (res.code === 0) {
      const tableFields = res.data || []
      // 更新现有字段的信息，保留权限配置
      fieldList.value.forEach(field => {
        const tableField = tableFields.find(tf => tf.fieldName === field.fieldName)
        if (tableField) {
          field.fieldChinese = tableField.fieldChinese || field.fieldChinese
          field.fieldDesc = tableField.fieldDesc || field.fieldDesc
          field.fieldType = tableField.fieldType || field.fieldType
        }
      })
      ElMessage.success('字段信息刷新成功')
    }
  } catch (error) {
    console.error('刷新字段失败:', error)
    ElMessage.error('刷新字段失败')
  }
}

const handleFieldSelection = (selection) => {
  selectedFields.value = selection
}

const deleteField = (index) => {
  ElMessageBox.confirm(
    '确定要删除这个字段权限配置吗？',
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(() => {
    fieldList.value.splice(index, 1)
    ElMessage.success('删除成功')
  }).catch(() => {
    // 用户取消删除
  })
}

// 处理加密状态变化
const handleEncryptionChange = (row) => {
  if (row.encrypted) {
    // 如果字段加密，自动隐藏并禁用所有其他权限
    row.visible = false
    row.editable = false
    row.exportable = false
    row.queryable = false
  }
}

// 处理可见性变化
const handleVisibilityChange = (row) => {
  if (!row.visible) {
    // 如果字段隐藏，自动禁用所有功能权限
    row.editable = false
    row.exportable = false
    row.queryable = false
  }
}

// 处理编辑字段对话框中的加密状态变化
const handleEditFieldEncryptionChange = () => {
  if (currentField.value.encrypted) {
    currentField.value.visible = false
    currentField.value.editable = false
    currentField.value.exportable = false
    currentField.value.queryable = false
  }
}

// 处理编辑字段对话框中的可见性变化
const handleEditFieldVisibilityChange = () => {
  if (!currentField.value.visible) {
    currentField.value.editable = false
    currentField.value.exportable = false
    currentField.value.queryable = false
  }
}

// 应用权限模板
const applyTemplate = (templateType) => {
  if (!currentField.value) return
  
  switch (templateType) {
    case 'full_access':
      // 完全开放：除加密外全部开启
      currentField.value.encrypted = false
      currentField.value.visible = true
      currentField.value.editable = true
      currentField.value.exportable = true
      currentField.value.queryable = true
      break
    case 'readonly':
      // 只读模式：可见+查询
      currentField.value.encrypted = false
      currentField.value.visible = true
      currentField.value.editable = false
      currentField.value.exportable = false
      currentField.value.queryable = true
      break
    case 'limited_edit':
      // 受限编辑：可见+编辑+查询
      currentField.value.encrypted = false
      currentField.value.visible = true
      currentField.value.editable = true
      currentField.value.exportable = false
      currentField.value.queryable = true
      break
    case 'encrypted':
      // 完全加密：仅加密开启
      currentField.value.encrypted = true
      currentField.value.visible = false
      currentField.value.editable = false
      currentField.value.exportable = false
      currentField.value.queryable = false
      break
  }
  
  ElMessage.success(`已应用${getTemplateName(templateType)}模板`)
}

// 获取模板名称
const getTemplateName = (templateType) => {
  const names = {
    'full_access': '完全开放',
    'readonly': '只读模式',
    'limited_edit': '受限编辑',
    'encrypted': '完全加密',
    'reverse_encrypted': '取消加密',
    'reverse_visible': '设为可见',
    'reverse_editable': '设为可编辑',
    'reverse_exportable': '设为可导出',
    'reverse_queryable': '设为可查询'
  }
  return names[templateType] || '未知模板'
}

// 监听当前角色变化
watch(() => props.currentRole, async (newRole) => {
  if (newRole) {
    // 确保受控表已加载，然后默认选中第一个受控表
    await loadControlledTables()
    console.log('当前角色:', moduleOptions.value[0])
    if (moduleOptions.value.length > 0) {
      selectedModule.value = moduleOptions.value[0].value
      setPermissionField(moduleOptions.value[0])
      // 加载配置
      loadDataPermissionConfig()
    }
    
  } else {
    resetConfig()
  }
}, { immediate: true })

// 监听选中模块变化
watch(selectedModule, async (newValue) => {
  if (newValue && props.currentRole) {
    // 先清空当前字段列表，避免显示上一个表的字段
    fieldList.value = []
    // 加载该表的权限配置
    await loadDataPermissionConfig()
    setPermissionField(moduleOptions.value.find(module => module.value == newValue))
  } else {
    // 重置配置
    dataScope.value = ''
    customSql.value = ''
    fieldList.value = []
    // 重置受控表字段配置为默认值
    tableFieldConfig.userField = 'created_by'
    tableFieldConfig.deptField = 'dept_code'
  }
})
// 设置权限字段
const setPermissionField = (module) => {
  console.log('设置权限字段', module)
  if (module) {
    tableFieldConfig.userField = module.userField
    tableFieldConfig.deptField = module.deptField
  }
}
// 加载数据范围字典
const loadDataScopeOptions = async () => {
  try {
    const dictData = await getDict('data_scope')
    if (dictData && dictData.length > 0) {
      dataScopeOptions.value = dictData.map(item => ({
        label: item.label,
        value: item.value
      }))
    } else {
      // 如果字典不存在，使用默认值
      dataScopeOptions.value = [
        { label: '全部数据权限', value: 'all' },
        { label: '自定义数据权限', value: 'custom' },
        { label: '本部门数据权限', value: 'dept' },
        { label: '本部门及以下数据权限', value: 'dept_and_child' },
        { label: '仅本人数据权限', value: 'self' }
      ]
    }
  } catch (error) {
    console.error('加载数据范围字典失败:', error)
    // 加载失败时使用默认值
    dataScopeOptions.value = [
      { label: '全部数据权限', value: 'all' },
      { label: '自定义数据权限', value: 'custom' },
      { label: '本部门数据权限', value: 'dept' },
      { label: '本部门及以下数据权限', value: 'dept_and_child' },
      { label: '仅本人数据权限', value: 'self' }
    ]
  }
}

// 加载用户字段字典
const loadUserFieldOptions = async () => {
  try {
    const dictData = await getDict('user_field')
    if (dictData && dictData.length > 0) {
      userFieldOptions.value = dictData.map(item => ({
        label: item.label,
        value: item.value
      }))
    } else {
      // 如果字典不存在，使用默认值
      userFieldOptions.value = [
        { label: 'created_by (创建者)', value: 'created_by' },
        { label: 'user_id (用户ID)', value: 'user_id' },
        { label: 'creator_id (创建者ID)', value: 'creator_id' },
        { label: 'owner_id (拥有者ID)', value: 'owner_id' }
      ]
    }
  } catch (error) {
    console.error('加载用户字段字典失败:', error)
    // 加载失败时使用默认值
    userFieldOptions.value = [
      { label: 'created_by (创建者)', value: 'created_by' },
      { label: 'user_id (用户ID)', value: 'user_id' },
      { label: 'creator_id (创建者ID)', value: 'creator_id' },
      { label: 'owner_id (拥有者ID)', value: 'owner_id' }
    ]
  }
}

// 加载部门字段字典
const loadDeptFieldOptions = async () => {
  try {
    const dictData = await getDict('dept_field')
    if (dictData && dictData.length > 0) {
      deptFieldOptions.value = dictData.map(item => ({
        label: item.label,
        value: item.value
      }))
    } else {
      // 如果字典不存在，使用默认值
      deptFieldOptions.value = [
        { label: 'dept_code (部门编码)', value: 'dept_code' },
        { label: 'dept_id (部门ID)', value: 'dept_id' },
        { label: 'department_id (部门ID)', value: 'department_id' },
        { label: 'org_id (组织ID)', value: 'org_id' }
      ]
    }
  } catch (error) {
    console.error('加载部门字段字典失败:', error)
    // 加载失败时使用默认值
    deptFieldOptions.value = [
      { label: 'dept_code (部门编码)', value: 'dept_code' },
      { label: 'dept_id (部门ID)', value: 'dept_id' },
      { label: 'department_id (部门ID)', value: 'department_id' },
      { label: 'org_id (组织ID)', value: 'org_id' }
    ]
  }
}


// 初始化
onMounted(() => {
  loadDatabaseTables()
  loadDataScopeOptions()
  loadUserFieldOptions()
  loadDeptFieldOptions()
})
</script>

<style scoped>
.data-permission-config {
  background: #f5f7fa;
  height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
}

.permission-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: white;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
  flex-shrink: 0;
}

.page-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.selected-role-info {
  font-size: 16px;
  font-weight: 600;
  color: #409eff;
}

.permission-modules,
.table-field-config,
.data-scope,
.field-permissions {
  margin-bottom: 12px;
  padding: 12px;
  background: white;
  border-radius: 6px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.section-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  border-bottom: 2px solid #409eff;
  padding-bottom: 8px;
  line-height: 1;
}

.help-icon {
  color: #909399;
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
  transition: color 0.3s;
}

.help-icon:hover {
  color: #409eff;
}

.module-header,
.field-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  padding: 8px;
  background: #f8f9fa;
  border-radius: 4px;
}

.field-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.batch-operations {
  display: flex;
  align-items: center;
}

.field-operations {
  display: flex;
  align-items: center;
  gap: 8px;
}

.selected-table-info {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.field-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.custom-sql {
  margin-top: 15px;
}

.custom-sql h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #606266;
}

.field-table {
  margin-top: 15px;
}

/* 表格滚动样式优化 */
.field-table .el-table__body-wrapper {
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: #c1c1c1 #f1f1f1;
}

.field-table .el-table__body-wrapper::-webkit-scrollbar {
  width: 8px;
}

.field-table .el-table__body-wrapper::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.field-table .el-table__body-wrapper::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.field-table .el-table__body-wrapper::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 表格固定高度时的样式调整 */
.field-table .el-table--border .el-table__header-wrapper {
  border-bottom: 1px solid #ebeef5;
}

.field-table .el-table__fixed-right {
  box-shadow: -1px 0 8px rgba(0, 0, 0, 0.1);
}

.field-dialog-content {
  padding: 0;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.field-count {
  font-size: 14px;
  color: #909399;
}

.dialog-footer {
  text-align: right;
}

.no-data-tip {
  text-align: center;
  padding: 40px 0;
}

/* 权限层次样式 */
.permission-levels {
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  overflow: hidden;
  margin-bottom: 16px;
}

.permission-level {
  padding: 16px;
  border-bottom: 1px solid #f0f2f5;
}

.permission-level:last-child {
  border-bottom: none;
}

.level-1 {
  background: #fef0f0;
  border-left: 4px solid #f56c6c;
}

.level-2 {
  background: #f0f9ff;
  border-left: 4px solid #409eff;
}

.level-3 {
  background: #f6ffed;
  border-left: 4px solid #67c23a;
}

.function-permissions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 8px;
}

.permission-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  font-style: italic;
}

/* 权限模板样式 */
.permission-templates {
  background: #fafafa;
  padding: 16px;
  border-radius: 6px;
  border: 1px dashed #d9d9d9;
}

.permission-templates .el-form-item {
  margin-bottom: 0;
}

.permission-templates .el-button-group .el-button {
  margin-right: 8px;
}

/* 模块卡片样式 */
.module-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.module-card {
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  padding: 16px;
  background: #fff;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 80px;
}

.module-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.15);
  transform: translateY(-2px);
}

.module-card.active {
  border-color: #409eff;
  background: #f0f9ff;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.2);
}

.module-info {
  flex: 1;
  min-width: 0;
}

.module-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
  word-break: break-all;
}

.module-table {
  font-size: 14px;
  color: #909399;
  font-family: 'Courier New', monospace;
  word-break: break-all;
}

.module-actions {
  margin-left: 16px;
  flex-shrink: 0;
}

.module-actions .el-button {
  opacity: 0.7;
  transition: opacity 0.3s ease;
}

.module-card:hover .module-actions .el-button {
  opacity: 1;
}

.module-card.active .module-actions .el-button {
  opacity: 1;
}

.permission-templates .el-button-group .el-button:last-child {
  margin-right: 0;
}

/* 表格中的开关样式优化 */
.field-table .el-switch {
  --el-switch-on-color: #67c23a;
  --el-switch-off-color: #dcdfe6;
}

.field-table .el-switch.is-disabled {
  opacity: 0.5;
}

/* 加密开关特殊样式 */
.field-table .el-switch[aria-checked="true"] {
  --el-switch-on-color: #f56c6c;
}

/* 权限层次指示器 */
.permission-level::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
}

.level-1::before {
  background: #f56c6c;
}

.level-2::before {
  background: #409eff;
}

.level-3::before {
  background: #67c23a;
}
</style>