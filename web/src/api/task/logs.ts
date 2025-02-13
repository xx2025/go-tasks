import request from "@/utils/request";

const TaskLogsAPI = {
  /**
   *
   * @param queryParams 查询参数
   */
  getPage(queryParams: TaskLogsPageQuery) {
    return request<any, PageResult<TaskLogPageVO[]>>({
      // url: `${USER_BASE_URL}/page`,
      url: `/task/logs`,
      method: "get",
      params: queryParams,
    });
  },
};

export default TaskLogsAPI;

/**
 * 分页查询对象
 */
export interface TaskLogsPageQuery extends PageQuery {
  /** 搜索关键字 */
  taskName?: string;
  status?: number;
  taskId?: number;
}

export interface TaskLogPageVO {
  /** 用户ID */
  id: number;
  createdAt: Date;
  updatedAt: Date;
  taskId: number;
  taskName: string;
  status: number;
  message: string;
}
