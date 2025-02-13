<template>
  <div class="dashboard">
    <el-card shadow="never">
      <el-row justify="space-between">
        <el-col :span="18" :xs="24">
          <div class="flex h-full items-center">
            <img
              class="w-20 h-20 mr-5 rounded-full"
              :src="userStore.userInfo.avatar + '?imageView2/1/w/80/h/80'"
             alt=""/>
            <div>
              <p>{{ greetings }}</p>
              <p class="text-sm text-gray">今日天气晴朗，气温在15℃至25℃之间，东南风。</p>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <el-row class="mt-3">
      <el-col>
        <el-card>
          <el-descriptions :column="1" border label-width="200px">
            <el-descriptions-item label="节点数量">
              <span class="font-bold">{{ dashboardData?.nodeCount }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="项目数量">
              {{ dashboardData?.projectCount }}
            </el-descriptions-item>
            <el-descriptions-item label="任务数量">
              {{ dashboardData?.taskCount }}
            </el-descriptions-item>
            <el-descriptions-item label="进程数量">
              {{ dashboardData?.processCount }}
            </el-descriptions-item>
            <el-descriptions-item label="CPU使用率">
              {{ dashboardData?.cpu }}
            </el-descriptions-item>
            <el-descriptions-item label="总内存">
              {{ dashboardData?.totalMem }} G
            </el-descriptions-item>
            <el-descriptions-item label="可用内存">
              {{ dashboardData?.freeMem }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>


  </div>
</template>

<script setup lang="ts">

defineOptions({
  name: "Dashboard",
  inheritAttrs: false,
});
import { useUserStore } from "@/store/modules/user";
import DashboardAPI, { DashboardData } from "@/api/dashboard";

const loading = ref(false);

const userStore = useUserStore();
const date: Date = new Date();
const greetings = computed(() => {
  const hours = date.getHours();
  if (hours >= 6 && hours < 8) {
    return "晨起披衣出草堂，轩窗已自喜微凉🌅！";
  } else if (hours >= 8 && hours < 12) {
    return "上午好，" + userStore.userInfo.nickname + "！";
  } else if (hours >= 12 && hours < 18) {
    return "下午好，" + userStore.userInfo.nickname + "！";
  } else if (hours >= 18 && hours < 24) {
    return "晚上好，" + userStore.userInfo.nickname + "！";
  } else {
    return "偷偷向银河要了一把碎星，只等你闭上眼睛撒入你的梦中，晚安🌛！";
  }
});


const dashboardData = ref<DashboardData>();

// 查询
function handleQuery() {
  loading.value = true;
  DashboardAPI.getData()
    .then((data) => {
      dashboardData.value = data
    })
    .finally(() => {
      loading.value = false;
    });
}

onMounted(() => {
  handleQuery();
});

</script>

<style lang="scss" scoped>
.dashboard {
  padding: 10px;
}
</style>
