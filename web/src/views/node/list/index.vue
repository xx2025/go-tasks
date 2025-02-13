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
            <el-form-item label="url：" prop="url">
              <el-input
                v-model="queryParams.url"
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
              <el-button
                v-hasPerm="['sys:user:add']"
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
            <!--            <el-table-column type="selection" width="50" align="center" />-->
            <el-table-column label="ID" prop="id" width="60" />
            <el-table-column label="名称" prop="name" />
            <el-table-column label="URL" prop="url" />
            <el-table-column label="任务数" prop="taskNum" />
            <el-table-column label="进程数" prop="processNum" />
            <el-table-column label="更新时间" prop="updatedAt" />
            <el-table-column label="操作" fixed="right">
              <template #default="scope">

                <el-button
                  v-hasPerm="'sys:user:delete'"
                  type="danger"
                  icon="delete"
                  link
                  size="small"
                  @click="handleDelete(scope.row.id)"
                >
                  删除
                </el-button>
                <el-button
                  v-hasPerm="'sys:user:delete'"
                  type="primary"
                  icon="RefreshLeft"
                  size="small"
                  link
                  @click="checkHealth(scope.row.id)"
                >
                  健康检查
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

    <!-- 节点表单 -->
    <el-dialog
      v-model="dialog.visible"
      :title="dialog.title"
      width="40%"
      :close-on-click-modal="false"
      append-to-body
      @close="handleCloseDialog"
    >
      <el-form ref="NodeFormRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item label="名称：" prop="name">
          <el-input v-model="formData.name" :readonly="!!formData.id" placeholder="请输入名称" />
        </el-form-item>

        <el-form-item label="url：" prop="url">
          <el-input v-model="formData.url" placeholder="请输入地址" />
        </el-form-item>


        <div class="dialog-footer" style="text-align: right">
          <el-button type="primary" @click="handleSubmit">确 定</el-button>
          <el-button @click="handleCloseDialog">取 消</el-button>
        </div>
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import NodeAPI, { NodeForm, NodePageQuery, NodePageVO } from "@/api/node/node";

defineOptions({
  name: "NodeList",
  inheritAttrs: false,
});

const queryFormRef = ref(ElForm);
const NodeFormRef = ref(ElForm);

const queryParams = reactive<NodePageQuery>({
  page: 1,
  pageSize: 10,
});

const pageData = ref<NodePageVO[]>();
const total = ref(0);
const loading = ref(false);

const dialog = reactive({
  visible: false,
  title: "新增节点",
});

const formData = reactive<NodeForm>({
  id: 0,
  name: "",
  url: "",
});

const rules = reactive({
  name: [{ required: true, message: "节点名称不能为空", trigger: "blur" }],
  url: [{ required: true, message: "URL不能为空", trigger: "blur" }],
});

// 查询
function handleQuery() {
  loading.value = true;
  NodeAPI.getPage(queryParams)
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
  handleQuery();
}

/**
 * 打开弹窗
 *
 * @param id
 * @param rowData 用户信息
 */
async function handleOpenDialog(id?: number, rowData?: NodeForm) {
  dialog.visible = true;
  if (id && rowData) {
    dialog.title = "修改节点";
    Object.assign(formData, { ...rowData });
  } else {
    dialog.title = "新增节点";
  }
}

// 关闭弹窗
function handleCloseDialog() {
  dialog.visible = false;
  NodeFormRef.value.resetFields();
  NodeFormRef.value.clearValidate();
  formData.id = undefined;
  formData.name = "";
  formData.url = "";
}

// 提交表单（防抖）
const handleSubmit = useDebounceFn(() => {
  NodeFormRef.value.validate((valid: boolean) => {
    if (valid) {
      console.log(formData);
      const id = formData.id;
      loading.value = true;
      if (id) {
        NodeAPI.update({
          id: formData.id,
          name: formData.name,
          url: formData.url,
          status: formData.status,
        })
          .then(() => {
            ElMessage.success("修改节点成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      } else {
        NodeAPI.add(formData)
          .then(() => {
            ElMessage.success("新增节点成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      }
    }
  });
}, 1000);

/**
 *
 * @param id  用户ID
 */
function handleDelete(id?: number) {
  console.log(id);
  // const userIds = [id || selectIds.value].join(",");
  if (!id) {
    ElMessage.warning("请选择删除项");
    return;
  }
  ElMessageBox.confirm("确认删除用户?", "警告", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(
    function () {
      loading.value = true;
      NodeAPI.deleteById(id)
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

// 健康检查请求方法
function checkHealth(id: number) {
  NodeAPI.checkHealth(id)
    .then(() => {
      ElMessageBox.alert("节点状态正常", "健康检查结果", {
        type: "success",
        confirmButtonText: "确定",
      });
    })
    .catch((error) => {
      ElMessageBox.alert("节点异常 ", "健康检查结果", {
        type: "error",
        confirmButtonText: "确定",
      });
    })
    .finally(() => {
      loading.value = false;
    });
}

onMounted(() => {
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
