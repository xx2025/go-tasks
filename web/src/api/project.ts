import request from "@/utils/request";

const ProjectAPI = {
  /**
   * 获取节点分页列表
   *
   * @param queryParams 查询参数
   */
  getPage(queryParams: ProjectPageQuery) {
    return request<any, PageResult<ProjectPageVO[]>>({
      // url: `${USER_BASE_URL}/page`,
      url: `/project/list`,
      method: "get",
      params: queryParams,
    });
  },

  /**
   *
   * @param data 表单数据
   */
  add(data: ProjectForm) {
    return request({
      url: `/project/save`,
      method: "post",
      data: data,
    });
  },

  /**
   *
   * @param data 表单数据
   */
  update(data: ProjectForm) {
    return request({
      url: `/project/save`,
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
      url: `/project/delete`,
      method: "post",
      data: { id: id },
    });
  },
};

export default ProjectAPI;

/**
 * 分页查询对象
 */
export interface ProjectPageQuery extends PageQuery {
  /** 搜索关键字 */
  name?: string;
}

/** 节点分页对象 */
export interface ProjectPageVO {
  /** 用户ID */
  id: number;
  createdAt?: Date;
  name?: string;
  describe?: string;
  taskNum?: number;
  processNum?: number;
}

export interface ProjectForm {
  /** 用户ID */
  id?: number;
  name?: string;
  describe?: string;
}
