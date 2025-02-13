<!-- 用户管理 -->
<template>
  <div class="app-container">
    <el-row :gutter="20">
      <!-- 部门树 -->
      <!--      <el-col :lg="4" :xs="24" class="mb-[12px]">
        <DeptTree v-model="queryParams.deptId" @node-click="handleQuery" />
      </el-col>-->

      <!-- 用户列表 -->
      <el-col :lg="24" :xs="24">
        <div class="search-bar">
          <el-form ref="queryFormRef" :model="queryParams" :inline="true" class="form-container">
            <el-form-item label="任务名称：" prop="taskName">
              <el-input
                v-model="queryParams.taskName"
                placeholder=""
                clearable
                style="width: 200px"
                @keyup.enter="handleQuery"
              />
            </el-form-item>
            <el-form-item label="任务ID：" prop="taskId">
              <el-input
                v-model="queryParams.taskId"
                placeholder=""
                clearable
                style="width: 200px"
                @keyup.enter="handleQuery"
              />
            </el-form-item>

            <el-form-item label="状态：" prop="status">
              <el-select
                v-model="queryParams.status"
                placeholder="全部"
                clearable
                class="!w-[180px]"
              >
                <el-option label="成功" :value="1" />
                <el-option label="进行中" :value="0" />
                <el-option label="失败" :value="-1" />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" icon="search" @click="handleQuery">搜索</el-button>
              <el-button icon="refresh" @click="handleResetQuery">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <el-card shadow="never">
          <div class="flex-x-between mb-10px">
            <div>
            </div>
            <div>
              <el-button
                type="primary"
                icon="refresh"
                @click="handleQuery"
              >
                刷新
              </el-button>
            </div>
          </div>

          <el-table v-loading="loading" :data="pageData">
            <el-table-column label="ID" prop="id" width="60" />
            <el-table-column label="任务名称" prop="taskName" />
            <el-table-column label="任务ID" prop="taskId" />

            <el-table-column label="状态" prop="status" width="100">
<!--              <template #default="scope">
                <el-tag :type="scope.row.status == 1 ? 'success' : 'info'">
                  {{ scope.row.status == 1 ? "成功" : "失败" }}
                </el-tag>
              </template>-->

              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>

            <el-table-column label="备注信息" prop="message" />

            <el-table-column label="调度时间" prop="createdAt" />
            <el-table-column label="更新时间" prop="updatedAt" />

          </el-table>
          <div class="pagination-container">
            <pagination
              v-if="total > 0"
              v-model:total="total"
              v-model:page="queryParams.page"
              v-model:limit="queryParams.pageSize"
              @pagination="handleQuery"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>


  </div>
</template>



<script setup lang="ts">
import TaskLogsAPI, { TaskLogsPageQuery, TaskLogPageVO } from "@/api/task/logs";
import router from "@/router";

defineOptions({
  name: "TaskLogs",
  inheritAttrs: false,
});

const queryFormRef = ref(ElForm);
const route = useRoute();

const queryParams = reactive<TaskLogsPageQuery>({
  page: 1,
  pageSize: 10,
});

const pageData = ref<TaskLogPageVO[]>();
const total = ref(0);
const loading = ref(false);



// 查询
function handleQuery() {
  loading.value = true;
  TaskLogsAPI.getPage(queryParams)
    .then((data) => {
      pageData.value = data.items;
      total.value = data.total;
    })
    .finally(() => {
      loading.value = false;
    });
}

// 重置查询
function handleResetQuery() {
  queryFormRef.value.resetFields();
  queryParams.page = 1;

  router.push({ name: "TaskLogs" }); // 更新路由

  // const query = { ...route.query };
  // delete query.taskId; // 移除 taskId 参数

  handleQuery();
}

// 获取状态的类型
function getStatusType(status: number): string {
  switch (status) {
    case 1:
      return "success"; // 成功
    case 0:
      return "warning"; // 进行中
    case -1:
      return "danger"; // 失败
    default:
      return "info"; // 默认
  }
}

// 获取状态的文本
function getStatusText(status: number): string {
  switch (status) {
    case 1:
      return "成功";
    case 0:
      return "进行中";
    case -1:
      return "失败";
    default:
      return "未知状态";
  }
}


onMounted(() => {
  const taskId = route.query.taskId;
  if (taskId) {
    queryParams.taskId = parseInt(taskId.toString(), 10); // 将 taskId 转换为数字
  }

  handleQuery();
});
</script>

<style scoped>
.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px; /* 根据需要调整 */
}
</style>
