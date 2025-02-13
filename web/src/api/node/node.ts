import request from "@/utils/request";

const NodeAPI = {
  /**
   * 获取节点分页列表
   *
   * @param queryParams 查询参数
   */
  getPage(queryParams: NodePageQuery) {
    return request<any, PageResult<NodePageVO[]>>({
      // url: `${USER_BASE_URL}/page`,
      url: `/node/list`,
      method: "get",
      params: queryParams,
    });
  },

  /**
   *
   * @param data 表单数据
   */
  add(data: NodeForm) {
    return request({
      url: `/node/add`,
      method: "post",
      data: data,
    });
  },

  /**
   *
   * @param data 表单数据
   */
  update(data: NodeForm) {
    return request({
      url: `/node/save`,
      method: "post",
      data: data,
    });
  },

  /**
   *
   * @param id
   */
  deleteById(id: number) {
    return request({
      url: `/node/delete`,
      method: "post",
      data: { id: id },
    });
  },

  checkHealth(id: number) {
    return request({
      url: `/node/health`,
      method: "post",
      data: { id: id },
    });
  },
};

export default NodeAPI;

/**
 * 分页查询对象
 */
export interface NodePageQuery extends PageQuery {
  /** 搜索关键字 */
  name?: string;
  url?: string;
}

/** 节点分页对象 */
export interface NodePageVO {
  id: number;
  /** 更新时间 */
  updatedAt?: Date;
  name?: string;
  url?: string;
  taskNum?: number;
  processNum?: number;
}

/** 用户表单类型 */
export interface NodeForm {
  /** 用户ID */
  id?: number;
  name?: string;
  url?: string;
}
