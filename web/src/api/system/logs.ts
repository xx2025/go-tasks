import request from "@/utils/request";

const LogsAPI = {
  /**
   * 获取用户分页列表
   *
   * @param queryParams 查询参数
   */
  getPage(queryParams: LogsPageQuery) {
    return request<any, PageResult<LogsPageVO[]>>({
      // url: `${USER_BASE_URL}/page`,
      url: `/user/logs`,
      method: "get",
      params: queryParams,
    });
  },
};

export default LogsAPI;

/**
 * 分页查询对象
 */
export interface LogsPageQuery extends PageQuery {
  /** 搜索关键字 */
  userId?: number;
  uri?: string;
}

/** 分页对象 */
export interface LogsPageVO {
  /** 用户ID */
  id: number;
  username?: string;
  data?: string;
  /** 更新时间 */
  updatedAt?: Date;
  createdAt?: Date;
}
