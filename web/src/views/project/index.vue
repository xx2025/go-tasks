<!-- 用户管理 -->
<template>
  <div class="app-container">
    <el-row :gutter="20">
      <!-- 列表 -->
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
            <el-table-column label="描述" prop="describe" />
            <el-table-column label="任务数" prop="taskNum" />
            <el-table-column label="进程数" prop="processNum" />

            <el-table-column label="创建时间" prop="createdAt" />
            <el-table-column label="操作" fixed="right">
              <template #default="scope">
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
                  v-hasPerm="'sys:user:delete'"
                  type="danger"
                  icon="delete"
                  link
                  size="small"
                  @click="handleDelete(scope.row.id)"
                >
                  删除
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
      <el-form ref="ProjectFormRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item label="名称：" prop="name">
          <el-input v-model="formData.name" :readonly="!!formData.id" placeholder="请输入名称" />
        </el-form-item>

        <el-form-item label="描述：" prop="describe">
          <el-input
            v-model="formData.describe"
            type="textarea"
            :rows="4"
            placeholder="请输入项目描述"
          />
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
import ProjectAPI, { ProjectForm, ProjectPageQuery, ProjectPageVO } from "@/api/project";

defineOptions({
  name: "ProjectList",
  inheritAttrs: false,
});

const queryFormRef = ref(ElForm);
const ProjectFormRef = ref(ElForm);

const queryParams = reactive<ProjectPageQuery>({
  page: 1,
  pageSize: 10,
});

const pageData = ref<ProjectPageVO[]>();
const total = ref(0);
const loading = ref(false);

const dialog = reactive({
  visible: false,
  title: "新增项目",
});

const formData = reactive<ProjectForm>({
  id: 0,
  name: "",
  describe: "",
});

const rules = reactive({
  name: [{ required: true, message: "项目名称不能为空", trigger: "blur" }],
  describe: [{ required: true, message: "项目描述不能为空", trigger: "blur" }],
});

// 查询
function handleQuery() {
  loading.value = true;
  ProjectAPI.getPage(queryParams)
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
    dialog.title = "修改项目";
    Object.assign(formData, { ...rowData });
  } else {
    dialog.title = "新增项目";
  }
}

// 关闭弹窗
function handleCloseDialog() {
  dialog.visible = false;
  ProjectFormRef.value.resetFields();
  ProjectFormRef.value.clearValidate();
  formData.id = undefined;
  formData.name = "";
  formData.describe = "";
}

// 提交表单（防抖）
const handleSubmit = useDebounceFn(() => {
  ProjectFormRef.value.validate((valid: boolean) => {
    if (valid) {
      const id = formData.id;
      loading.value = true;
      if (id) {
        ProjectAPI.update({
          id: formData.id,
          name: formData.name,
          describe: formData.describe,
        })
          .then(() => {
            ElMessage.success("修改项目成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      } else {
        ProjectAPI.add(formData)
          .then(() => {
            ElMessage.success("新增项目成功");
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
  ElMessageBox.confirm("确认删除项目?", "警告", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(
    function () {
      loading.value = true;
      ProjectAPI.deleteById(id)
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
