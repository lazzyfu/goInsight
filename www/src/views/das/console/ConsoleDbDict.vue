<template>
  <a-modal :open="open" title="数据字典" width="65%" :footer="null" @cancel="handleCancel">
    <a-spin tip="加载中..." :spinning="loading">
      <div>
        <strong class="dict-title">一、表名索引</strong>
        <div class="table-container" v-html="dictIndex"></div>
      </div>
      <div>
        <strong class="dict-title">二、表结构详情</strong>
        <div v-html="dictData"></div>
      </div>
    </a-spin>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['update:open'])
defineProps({
  open: Boolean,
})

const dictIndex = ref('')
const dictData = ref('')
const loading = ref(false)


const handleCancel = () => {
  emit('update:open')
}

const render = (data) => {
  loading.value = true;

  const generatorTableRows = (items, separator) => {
    let html = '';
    items.forEach((item, index) => {
      html += `<tr><td>${index + 1}</td>`;
      item.split(separator).forEach((field) => {
        html += `<td>${field}</td>`;
      });
      html += '</tr>';
    });
    return html;
  };

  let num = 0;
  dictIndex.value = '';
  dictData.value = '';

  data.forEach((row) => {
    num += 1;
    const { TABLE_NAME, TABLE_COMMENT, CREATE_TIME, COLUMNS_INFO, INDEXES_INFO } = row;
    const tableComment = TABLE_COMMENT || 'None';

    // 生成索引链接
    dictIndex.value += `
      <div style="margin-top: 8px; padding-left: 12px;">
        <a href="#${TABLE_NAME}" style="color: #1890ff; text-decoration: none; font-weight: 500;">${num}、${TABLE_NAME} ............ ${tableComment}</a><br>
      </div>
    `;

    // 生成表格内容
    const tableRows = generatorTableRows(COLUMNS_INFO.split('<a>'), '<b>');
    const indexRows = generatorTableRows(INDEXES_INFO.split('<a>'), '<b>');

    dictData.value += `
     <div class="table-container">
        <a style="color: #262626; font-weight: 600; font-size: 14px;" name="${TABLE_NAME}">
          ${num}、表名: ${TABLE_NAME} 备注: ${tableComment} 创建时间: ${CREATE_TIME}
        </a>
        <table class="modern-table">
          <thead>
            <tr>
              <th>序列</th>
              <th>列名</th>
              <th>数据类型</th>
              <th>可空</th>
              <th>默认值</th>
              <th>字符集</th>
              <th>排序规则</th>
              <th>备注</th>
            </tr>
          </thead>
          <tbody>
            ${tableRows}
          </tbody>
        </table>
        <table class="modern-table index-table">
          <thead>
            <tr>
              <th>序列</th>
              <th>索引名</th>
              <th>唯一</th>
              <th>基数</th>
              <th>类型</th>
              <th>包含字段</th>
            </tr>
          </thead>
          <tbody>
            ${indexRows}
          </tbody>
        </table>
      </div>
    `;
  });

  loading.value = false;
};

defineExpose({
  render
})

</script>


<style scoped>
:deep(.dict-title) {
  color: #262626;
  font-size: 18px;
  font-weight: 600;
}

/* 添加现代化表格样式，替代Bootstrap */
:deep(.table-container) {
  height: auto;
  overflow: auto;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  padding: 16px;
  margin: 12px 0;
  background: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

:deep(.modern-table) {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
  margin-top: 12px;
  margin-bottom: 8px;
  background: #ffffff;
  border-radius: 6px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

:deep(.modern-table thead) {
  background: linear-gradient(135deg, #f5f5f5 0%, #e8e8e8 100%);
}

:deep(.modern-table th) {
  padding: 12px 8px;
  text-align: left;
  font-weight: 600;
  color: #262626;
  border-bottom: 2px solid #d9d9d9;
  font-size: 13px;
  white-space: nowrap;
}

:deep(.modern-table td) {
  padding: 10px 8px;
  border-bottom: 1px solid #f0f0f0;
  color: #595959;
  font-size: 12px;
  transition: background-color 0.2s ease;
}

:deep(.modern-table tbody tr:nth-child(even)) {
  background-color: #fafafa;
}

:deep(.modern-table tbody tr:hover) {
  background-color: #e6f7ff;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

:deep(.modern-table tbody tr:last-child td) {
  border-bottom: none;
}

:deep(.index-table) {
  margin-top: 20px;
}

/* 响应式设计优化 */
@media (max-width: 768px) {
  :deep(.modern-table) {
    font-size: 11px;
  }

  :deep(.modern-table th,
    .modern-table td) {
    padding: 8px 6px;
  }

  :deep(.table-container) {
    padding: 12px;
    margin: 8px 0;
  }
}

/* 滚动条美化 */
:deep(.table-container::-webkit-scrollbar) {
  width: 6px;
  height: 6px;
}

:deep(.table-container::-webkit-scrollbar-track) {
  background: #f1f1f1;
  border-radius: 3px;
}

:deep(.table-container::-webkit-scrollbar-thumb) {
  background: #c1c1c1;
  border-radius: 3px;
}

:deep(.table-container::-webkit-scrollbar-thumb:hover) {
  background: #a8a8a8;
}
</style>
