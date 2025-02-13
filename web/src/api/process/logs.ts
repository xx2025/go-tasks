import request from "@/utils/request";

const ProcessLogsAPI = {
  /**
   *
   * @param queryParams 查询参数
   */
  getPage(queryParams: ProcessLogsPageQuery) {
    return request<any, PageResult<ProcessLogPageVO[]>>({
      // url: `${USER_BASE_URL}/page`,
      url: `/process/logs`,
      method: "get",
      params: queryParams,
    });
  },
};

export default ProcessLogsAPI;

/**
 * 分页查询对象
 */
export interface ProcessLogsPageQuery extends PageQuery {
  /** 搜索关键字 */
  processName?: string;
  processId?: number;
}

export interface ProcessLogPageVO {
  /** 用户ID */
  id: number;
  createdAt: Date;
  updatedAt: Date;
  processId: number;
  processName: string;
  message: string;
}
