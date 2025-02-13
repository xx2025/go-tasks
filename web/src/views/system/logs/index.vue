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
          <el-form ref="queryFormRef" :model="queryParams" :inline="true">
            <el-form-item label="用户名：" prop="userId">
              <el-autocomplete
                v-model="selectedUsername"
                :fetch-suggestions="fetchUserOptions"
                :suggestions="userOptions"
                placeholder="请输入用户名"
                clearable
                @select="handleSelect"
              >
                <template #default="{ item }">
                  {{ item?.username }}
                  <!-- 显示节点名称 -->
                </template>
              </el-autocomplete>
            </el-form-item>

            <el-form-item label="URI：" prop="uri">
              <el-input
                v-model="queryParams.uri"
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
            <!--            <el-table-column type="selection" width="50" align="center" />-->
            <el-table-column label="ID" prop="id" />
            <el-table-column label="用户名" prop="username" />
            <el-table-column label="URI" prop="uri" />
            <el-table-column label="请求数据" prop="data" width="350">
              <template #default="scope">
                <div class="text-ellipsis">
                  {{ scope.row.data }}
                </div>
              </template>
            </el-table-column>

            <el-table-column label="操作时间" prop="createdAt" />
            <el-table-column label="操作" fixed="right">
              <template #default="scope">
                <el-button
                  v-hasPerm="'sys:user:password:reset'"
                  type="primary"
                  icon="RefreshLeft"
                  size="small"
                  link
                  @click="handleOpenDialog(scope.row)"
                >
                  查看数据
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <pagination class="pagination-container"
            v-if="total > 0"
            v-model:total="total"
            v-model:page="queryParams.page"
            v-model:limit="queryParams.pageSize"
            @pagination="handleQuery"
          />
        </el-card>
      </el-col>
    </el-row>
    <el-dialog
      v-model="dialog.visible"
      :title="dialog.title"
      width="50%"
      :close-on-click-modal="false"
      append-to-body
      @close="handleCloseDialog"
    >
      <div style="height: 40vh" v-html="dialog.content" />
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import LogsAPI, { LogsPageQuery, LogsPageVO } from "@/api/system/logs";
import SelectorsAPI, { UserSelector } from "@/api/system/selectors";

defineOptions({
  name: "UserLogs",
  inheritAttrs: false,
});

const queryFormRef = ref<typeof ElForm>();
const userOptions = ref<UserSelector[]>([]);
const queryParams = reactive<LogsPageQuery>({
  page: 1,
  pageSize: 10,
});
const selectedUsername = ref("");
const pageData = ref<LogsPageVO[]>();
const total = ref(0);
const loading = ref(false);

const dialog = reactive({
  visible: false,
  title: "请求数据详情",
  content: "",
});

/**
 * 打开弹窗
 *
 * @param rowData 用户信息
 */
async function handleOpenDialog(rowData?: LogsPageVO) {
  dialog.visible = true;
  dialog.content = rowData?.data ?? "";
}

// 初始化下拉数据
function initUserOptions() {
  SelectorsAPI.getUserSelectors().then((data) => {
    if (Array.isArray(data)) {
      userOptions.value = data.map((user) => ({
        id: user.id, // 用于存储 id
        username: user.username, // 用于显示 username
      }));
    }else{
      console.error("返回的数据不是数组类型：", data);
    }

  });
}

// 获取下拉选项数据
function fetchUserOptions(queryString: string, cb: Function) {
  // 过滤用户列表并返回匹配的结果
  const results = userOptions.value.filter((user) =>
    String(user.id).toLowerCase().includes(queryString.toLowerCase())
  );
  cb(results);
}

// 处理下拉选项被选中
function handleSelect(item: {  username: string; id: number }) {
  queryParams.userId = item.id; // 更新 userId 为选中用户的 id
  selectedUsername.value = item.username; // 更新选中的用户名
}

// 查询
function handleQuery() {
  LogsAPI.getPage(queryParams)
    .then((data) => {
      pageData.value = data.items;
      total.value = data.total;
    })
    .finally(() => {
      loading.value = false;
    });
}

// 关闭弹窗
function handleCloseDialog() {
  dialog.visible = false;
  dialog.content = "";
}
// 重置查询
function handleResetQuery() {
  queryFormRef.value?.resetFields();
  queryParams.page = 1;
  selectedUsername.value = "";
  handleQuery();
}
onMounted(() => {
  initUserOptions();
  handleQuery();
});
</script>

<style scoped>
.text-ellipsis {
  white-space: nowrap; /* 禁止文本换行 */
  overflow: hidden; /* 隐藏超出部分 */
  text-overflow: ellipsis; /* 使用省略号表示被截断的文本 */
  max-width: 350px; /* 设置最大宽度，根据需要调整 */
}
.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px; /* 根据需要调整 */
}
</style>
