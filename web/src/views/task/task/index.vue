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
            <el-form-item label="名称：" prop="name">
              <el-input
                v-model="queryParams.name"
                placeholder=""
                clearable
                style="width: 200px"
                @keyup.enter="handleQuery"
              />
            </el-form-item>
            <el-form-item label="节点：" prop="nodeId">
              <el-autocomplete
                v-model="selectedNodeName"
                :fetch-suggestions="fetchNodeOptions"
                :suggestions="nodeOptions"
                placeholder="请输入节点名称"
                clearable
                @select="nodeSelector"
              >
                <template #default="{ item }">
                  {{ item?.name }}
                  <!-- 显示节点名称 -->
                </template>
              </el-autocomplete>
            </el-form-item>

            <el-form-item label="项目：" prop="projectId">
              <el-autocomplete
                v-model="selectedProjectName"
                :fetch-suggestions="fetchProjectOptions"
                :suggestions="projectOptions"
                placeholder="请输入项目名称"
                clearable
                @select="projectSelector"
              >
                <template #default="{ item }">
                  {{ item?.name }}
                  <!-- 显示节点名称 -->
                </template>
              </el-autocomplete>
            </el-form-item>

            <el-form-item label="状态：" prop="status">
              <el-select
                v-model="queryParams.status"
                placeholder="全部"
                clearable
                class="!w-[180px]"
              >
                <el-option label="启用" :value="1" />
                <el-option label="禁用" :value="0" />
              </el-select>
            </el-form-item>
            <el-form-item label="我关注的：">
              <el-switch v-model="queryParams.following" @change="handleQuery" />
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
              <el-button
                type="success"
                icon="plus"
                @click="handleOpenDialog()"
              >
                新增
              </el-button>
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
            <el-table-column label="名称" prop="name" />
            <el-table-column label="频率" prop="spec" />
            <el-table-column label="节点" prop="nodeId">
              <template #default="scope">
                {{ getNodeName(scope.row.nodeId) }}
              </template>
            </el-table-column>
            <el-table-column label="项目" prop="projectId">
              <template #default="scope">
                {{ getProjectName(scope.row.projectId) }}
              </template>
            </el-table-column>
            <el-table-column label="状态" prop="status" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status == 1 ? 'success' : 'info'">
                  {{ scope.row.status == 1 ? "启用" : "禁用" }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="描述" prop="describe" width="300">
              <template #default="scope">
                <div class="cell-content">{{ scope.row.describe }}</div>
              </template>
            </el-table-column>
            <el-table-column label="更新时间" prop="updatedAt" />
            <el-table-column label="操作" fixed="right">
              <template #default="scope">
                <el-button
                  type="primary"
                  icon="Plus"
                  size="small"
                  link
                  @click="handleFollow(scope.row.id)"
                >
                  关注
                </el-button>
                <el-button
                  v-if="scope.row.isFollowing"
                  type="danger"
                  icon="Minus"
                  size="small"
                  link
                  @click="handleUnFollow(scope.row.id)"
                >
                  取消关注
                </el-button>
                <el-button
                  type="primary"
                  icon="edit"
                  link
                  size="small"
                  @click="handleOpenDetail(scope.row.id)"
                >
                  详情
                </el-button>
                <el-button
                  type="primary"
                  icon="edit"
                  link
                  size="small"
                  @click="handleOpenDialog(scope.row.id, scope.row)"
                >
                  编辑
                </el-button>
                <el-button
                  type="danger"
                  icon="delete"
                  link
                  size="small"
                  @click="handleDelete(scope.row.id)"
                >
                  删除
                </el-button>

                <el-button
                  type="primary"
                  icon="edit"
                  link
                  size="small"
                  @click="handleLogs(scope.row.id)"
                >
                  查看日志
                </el-button>
              </template>
            </el-table-column>
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

    <!-- 任务表单 -->
    <el-dialog
      v-model="dialog.visible"
      :title="dialog.title"
      width="55%"
      :close-on-click-modal="false"
      append-to-body
      @close="handleCloseDialog"
    >
      <el-form ref="TaskFormRef" :model="formData" :rules="rules" label-width="150px">
        <el-form-item label="名称：" prop="name">
          <el-input
            v-model="formData.name"
            :readonly="!!formData.id"
            :disabled="!!formData.id"
            placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="频率：" prop="spec">
          <el-input v-model="formData.spec" placeholder="每分钟执行示例：* * * * *" />
        </el-form-item>

        <el-form-item label="状态：" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="未结束时跳过：" prop="isSingle">
          <el-radio-group v-model="formData.isSingle">
            <el-radio :value="1">是</el-radio>
            <el-radio :value="0">否</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="节点：" prop="nodeId">
          <el-autocomplete
            v-model="selectedNodeName"
            :fetch-suggestions="fetchNodeOptions"
            :suggestions="nodeOptions"
            placeholder="请选择节点"
            clearable
            @select="nodeSelector"
          >
            <template #default="{ item }">
              {{ item?.name }}
              <!-- 显示节点名称 -->
            </template>
          </el-autocomplete>
        </el-form-item>

        <el-form-item label="项目：" prop="projectId">
          <el-autocomplete
            v-model="selectedProjectName"
            :fetch-suggestions="fetchProjectOptions"
            :suggestions="projectOptions"
            placeholder="请选择项目"
            clearable
            @select="projectSelector"
          >
            <template #default="{ item }">
              {{ item?.name }}
              <!-- 显示节点名称 -->
            </template>
          </el-autocomplete>
        </el-form-item>
        <el-form-item label="描述：" prop="describe">
          <el-input
            v-model="formData.describe"
            type="textarea"
            :rows="3"
            placeholder="请输入项目描述"
          />
        </el-form-item>
        <div class="dialog-footer" style="text-align: right">
          <el-button type="primary" @click="handleSubmit">确 定</el-button>
          <el-button @click="handleCloseDialog">取 消</el-button>
        </div>
      </el-form>
    </el-dialog>
    <!-- 详情 -->
    <el-dialog
      v-model="detail.visible"
      :title="detail.title"
      width="55%"
      :close-on-click-modal="false"
      append-to-body
      @close="handleCloseDetail"
    >
      <el-col>
        <el-descriptions :column="1" border label-width="140px">
          <el-descriptions-item label="任务名称">
            {{ (detail.taskDetail as TaskDetail)?.name || "" }}
          </el-descriptions-item>
          <el-descriptions-item label="任务描述">
            {{ (detail.taskDetail as TaskDetail)?.describe || "" }}
          </el-descriptions-item>
          <el-descriptions-item label="调度频率">
            {{ (detail.taskDetail as TaskDetail)?.spec || "" }}
          </el-descriptions-item>
          <el-descriptions-item label="调度状态">
<!--            {{ (detail.taskDetail as TaskDetail)?.scheduleState == 1 ? "调度中" :  "已停止" }}-->
            <template v-if="(detail.taskDetail as TaskDetail)?.scheduleState == 1">
              <el-button
                type="primary"
                :style="{ backgroundColor: '#00FF00', color: '#000'}"
                size="small"
              >
                调度中
              </el-button>
            </template>
            <template v-else>
              <el-button
                type="primary"
                :style="{ backgroundColor: '#FFD700', color: '#000'}"
                size="small"
              >
                已停止
              </el-button>
            </template>
          </el-descriptions-item>
          <el-descriptions-item label="节点">
            {{ (detail.taskDetail as TaskDetail)?.nodeName || "" }}
          </el-descriptions-item>
          <el-descriptions-item label="项目">
            {{ (detail.taskDetail as TaskDetail)?.projectName || "" }}
          </el-descriptions-item>
          <el-descriptions-item label="PID">
            {{ (detail.taskDetail as TaskDetail)?.pid || 0 }}
          </el-descriptions-item>
          <el-descriptions-item label="操作">
            <el-button
              v-if="(detail.taskDetail as TaskDetail)?.pid > 0"
              type="primary"
              :style="{ backgroundColor: '#FFD700', color: '#000'}"
              size="small"
              @click="handleStop((detail.taskDetail as TaskDetail)?.id)"
            >
              停止本次运行
            </el-button>
            <el-button
              v-if="(detail.taskDetail as TaskDetail)?.pid <= 0"
              type="primary"
              :style="{ backgroundColor: '#FFD700', color: '#000'}"
              size="small"
              @click="handleExec((detail.taskDetail as TaskDetail)?.id)"
            >
              执行一次
            </el-button>
          </el-descriptions-item>



        </el-descriptions>
        <div class="dialog-footer" style="text-align: right; margin-top: 10px">
          <el-button @click="handleCloseDetail">关闭</el-button>
        </div>
      </el-col>
    </el-dialog>
  </div>
</template>



<script setup lang="ts">
import TaskAPI, { TaskForm, TaskPageQuery, TaskPageVO, TaskDetail } from "@/api/task";
import SelectorsAPI, { NodeSelector, ProjectSelector } from "@/api/system/selectors";
import router from "@/router";
defineOptions({
  name: "Task",
  inheritAttrs: false,
});

const queryFormRef = ref(ElForm);
const TaskFormRef = ref(ElForm);

const nodeOptions = ref<NodeSelector[]>([]);
const projectOptions = ref<ProjectSelector[]>([]);

const selectedNodeName = ref("");
const selectedProjectName = ref("");

const queryParams = reactive<TaskPageQuery>({
  page: 1,
  pageSize: 10,
  following: false,
});

const pageData = ref<TaskPageVO[]>();
const total = ref(0);
const loading = ref(false);

const dialog = reactive({
  visible: false,
  title: "新增任务",
});

const detail = reactive({
  visible: false,
  title: "任务详情",
  taskDetail: null as TaskDetail | null, // 用于存储任务的详细信息
});

const formData = reactive<TaskForm>({
  id: 0,
  name: "",
  spec: "",
  status: 1,
  isSingle: 1,
  projectId: 0,
  nodeId: 0,
  describe: "",
});


const rules = reactive({
  name: [{ required: true, message: "名称不能为空", trigger: "blur" }],
  spec: [{ required: true, message: "执行频率不能为空", trigger: "blur" }],
  nodeId: [
    { required: true, message: "请选择节点", trigger: "blur" },
    {
      validator: (rule: any, value: any, callback: Function) => {
        if (value <= 0) {
          callback(new Error("请选择节点"));
        } else {
          callback();
        }
      },
      trigger: "blur"
    }
  ],
  projectId: [
    { required: true, message: "请选择项目", trigger: "blur" },
    {
      validator: (rule: any, value: any, callback: Function) => {
        if (value <= 0) {
          callback(new Error("请选择项目"));
        } else {
          callback();
        }
      },
      trigger: "blur"
    }
  ],
});

// 查询节点名称
function getNodeName(nodeId: number) {
  const node = nodeOptions.value.find((p) => p.id === nodeId);
  return node ? node.name : "未知节点";
}

// 查询项目名称
function getProjectName(projectId: number) {
  const project = projectOptions.value.find((p) => p.id === projectId);
  return project ? project.name : "未知项目";
}

function fetchNodeOptions(queryString: string, cb: Function) {
  // 过滤用户列表并返回匹配的结果
  const results = nodeOptions.value.filter((node) =>
    String(node.id).toLowerCase().includes(queryString.toLowerCase())
  );
  cb(results);
}

function nodeSelector(item: { id: number; name: string}) {
  queryParams.nodeId = item.id; //
  formData.nodeId = item.id;
  selectedNodeName.value = item.name; // 更新选中的用户名
}

// 初始化下拉数据
function initNodeSelector() {
  SelectorsAPI.getNodeSelectors().then((data) => {
    if (Array.isArray(data)) {
      nodeOptions.value = data.map((node: NodeSelector) => ({
        id: node.id,
        name: node.name,
      }));
    }else{
      console.error("返回的数据不是数组类型：", data);
    }

  })
}

function fetchProjectOptions(queryString: string, cb: Function) {
  // 过滤用户列表并返回匹配的结果
  const results = projectOptions.value.filter((project) =>
    String(project.id).toLowerCase().includes(queryString.toLowerCase())
  );
  cb(results);
}

function projectSelector(item: { name: string; id: number }) {
  queryParams.projectId = item.id; //
  formData.projectId = item.id;
  selectedProjectName.value = item.name;
}

// 初始化下拉数据
function initProjectSelector() {
  SelectorsAPI.getProjectSelectors().then((data) => {
    if (Array.isArray(data)){
      projectOptions.value = data.map((project: ProjectSelector) => ({
        id: project.id,
        name: project.name,
      }));
    }

  });
}

// 查询
function handleQuery() {
  loading.value = true;
  TaskAPI.getPage(queryParams)
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
  queryParams.following = false;
  selectedNodeName.value = "";
  selectedProjectName.value = "";
  handleQuery();
}

/**
 * 打开弹窗
 *
 * @param id
 * @param rowData 用户信息
 */
async function handleOpenDialog(id?: number, rowData?: TaskForm) {
  dialog.visible = true;
  if (id && rowData) {
    dialog.title = "修改节点";
    Object.assign(formData, { ...rowData });

    // 查找并设置项目名称
    const project = projectOptions.value.find((p) => p.id === rowData.projectId);
    if (project) {
      selectedProjectName.value = project.name;
    }

    // 查找并设置节点名称
    const node = nodeOptions.value.find((n) => n.id === rowData.nodeId);
    if (node) {
      selectedNodeName.value = node.name;
    }
  } else {
    dialog.title = "新增节点";
  }
}

// 关闭弹窗
function handleCloseDialog() {
  dialog.visible = false;
  TaskFormRef.value.resetFields();
  TaskFormRef.value.clearValidate();
  formData.id = undefined;
  formData.status = 1;
  formData.name = "";
  formData.spec = "";
  formData.projectId = 0;
  formData.nodeId = 0;
  formData.describe = "";
  selectedProjectName.value = "";
  selectedNodeName.value = "";
}


async function handleOpenDetail(id: number) {
  detail.visible = true;
  detail.title = "任务详情";
  try {
    loading.value = true;
    const data = await TaskAPI.getDetail({id: id}); // 假设后端接口返回任务的详细信息
    detail.taskDetail = data; // 将获取到的详情数据存储到detail对象中
  } catch (error) {
    ElMessage.error("获取任务详情失败");
  } finally {
    loading.value = false;
  }
}

// 关闭弹窗
function handleCloseDetail() {
  detail.visible = false;
  detail.taskDetail = null;
}




// 提交表单（防抖）
const handleSubmit = useDebounceFn(() => {
  TaskFormRef.value.validate((valid: boolean) => {
    if (valid) {
      const id = formData.id;
      loading.value = true;
      if (id) {
        TaskAPI.update(formData)
          .then(() => {
            ElMessage.success("修改任务成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      } else {
        TaskAPI.add(formData)
          .then(() => {
            ElMessage.success("新增任务成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      }
    }
  });
}, 1000);

function handleFollow(id: number) {
  if (!id) {
    ElMessage.warning("请选择");
    return;
  }
  loading.value = true;
  TaskAPI.followingById(id)
    .then(() => {
      ElMessage.success("关注成功");
      handleResetQuery();
    })
    .finally(() => (loading.value = false));
}

function handleUnFollow(id: number) {
  loading.value = true;
  TaskAPI.unFollowingById(id)
    .then(() => {
      ElMessage.success("操作成功");
      handleResetQuery();
    })
    .finally(() => {
      loading.value = false;
    });
}

function handleStop(id: number) {

  ElMessageBox.confirm("确认停止本次运行?", "警告", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(
    function () {
      loading.value = true;
      TaskAPI.stopById(id)
        .then(() => {
          ElMessage.success("操作成功");
          handleCloseDetail();
        })
        .finally(() => (loading.value = false));
    },
    function () {
      ElMessage.info("已取消操作");
    }
  );
}



function handleExec(id: number) {

  ElMessageBox.confirm("确认执行该任务?", "警告", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(
    function () {
      loading.value = true;
      TaskAPI.execById(id)
        .then(() => {
          ElMessage.success("操作成功");
        })
        .finally(() => (loading.value = false));
    },
    function () {
      ElMessage.info("已取消操作");
    }
  );
}



/**
 *
 * @param id  用户ID
 */
function handleDelete(id?: number) {
  if (!id) {
    ElMessage.warning("请选择删除项");
    return;
  }
  ElMessageBox.confirm("确认删除该任务?", "警告", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(
    function () {
      loading.value = true;
      TaskAPI.deleteById(id)
        .then(() => {
          ElMessage.success("删除成功");
          handleResetQuery();
        })
        .finally(() => (loading.value = false));
    },
    function () {
      ElMessage.info("已取消删除");
    }
  );
}

function handleLogs(taskId: number) {
  // 跳转到日志页面，并传递任务 ID 和任务名称作为查询参数
  router.push({ name: "TaskLogs", query: { taskId: taskId.toString() } });
}


onMounted(() => {
  handleQuery();
  initNodeSelector();
  initProjectSelector();
});
</script>
<style scoped>
.cell-content {
  max-width: 300px; /* 设置最大宽度 */
  overflow: hidden; /* 隐藏溢出内容 */
  text-overflow: ellipsis; /* 使用省略号表示溢出内容 */
  white-space: nowrap; /* 防止文本换行 */
}
</style>

<style scoped>
.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px; /* 根据需要调整 */
}

</style>
