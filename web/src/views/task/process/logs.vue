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
            <el-form-item label="进程ID：" prop="processId">
              <el-input
                v-model="queryParams.processId"
                placeholder=""
                clearable
                style="width: 200px"
                @keyup.enter="handleQuery"
              />
            </el-form-item>
            <el-form-item label="进程名称：" prop="processName">
              <el-input
                v-model="queryParams.processName"
                placeholder=""
                clearable
                style="width: 200px"
                @keyup.enter="handleQuery"
              />
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
              <!--              <el-button-->
              <!--                v-hasPerm="['sys:user:add']"-->
              <!--                type="success"-->
              <!--                icon="plus"-->
              <!--                @click="handleOpenDialog()"-->
              <!--              >-->
              <!--                新增-->
              <!--              </el-button>-->
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
            <el-table-column label="进程名称" prop="processName" />
            <el-table-column label="进程ID" prop="processId" />


            <el-table-column label="备注信息" prop="message" />

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
import ProcessLogsAPI, { ProcessLogsPageQuery, ProcessLogPageVO } from "@/api/process/logs";
import router from "@/router";
defineOptions({
  name: "ProcessLogs",
  inheritAttrs: false,
});

const queryFormRef = ref(ElForm);
const route = useRoute();

const queryParams = reactive<ProcessLogsPageQuery>({
  page: 1,
  pageSize: 10,
});

const pageData = ref<ProcessLogPageVO[]>();
const total = ref(0);
const loading = ref(false);



// 查询
function handleQuery() {
  loading.value = true;
  ProcessLogsAPI.getPage(queryParams)
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

  router.push({ name: "ProcessLogs" }); // 更新路由

  handleQuery();
}


onMounted(() => {

  const processId = route.query.processId;
  if (processId) {
    queryParams.processId = parseInt(processId.toString(), 10); // 将 taskId 转换为数字
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
