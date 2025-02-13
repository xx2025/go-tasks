import request from "@/utils/request";

const DashboardAPI = {

  getData() {
    return request<any, DashboardData>({
      // url: `${USER_BASE_URL}/page`,
      url: `/dashboard`,
      method: "get",
    });
  },
};

export default DashboardAPI;

export interface DashboardData {
  nodeCount: number;
  projectCount: number;
  taskCount: number;
  processCount: number;
  cpu: string;
  totalMem: string;
  freeMem: string;
}
