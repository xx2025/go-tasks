import request from "@/utils/request";

const TaskAPI = {
  /**
   *
   * @param queryParams 查询参数
   */
  getPage(queryParams: TaskPageQuery) {
    return request<any, PageResult<TaskPageVO[]>>({
      // url: `${USER_BASE_URL}/page`,
      url: `/task/list`,
      method: "get",
      params: queryParams,
    });
  },

  getDetail(queryParams: {id: number}) {
    return request<any, TaskDetail>({
      // url: `${USER_BASE_URL}/page`,
      url: `/task/detail`,
      method: "get",
      params: queryParams,
    });
  },

  /**
   *
   * @param data 表单数据
   */
  add(data: TaskForm) {
    return request({
      url: `/task/save`,
      method: "post",
      data: data,
    });
  },

  /**
   *
   * @param data 表单数据
   */
  update(data: TaskForm) {
    return request({
      url: `/task/save`,
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
      url: `/task/delete`,
      method: "post",
      data: { id: id },
    });
  },

  execById(id: number) {
    return request({
      url: `/task/exec`,
      method: "post",
      data: { id: id },
    });
  },

  stopById(id: number) {
    return request({
      url: `/task/stop`,
      method: "post",
      data: { id: id },
    });
  },

  followingById(id: number) {
    return request({
      url: `/task/following`,
      method: "post",
      data: { id: id },
    });
  },
  unFollowingById(id: number) {
    return request({
      url: `/task/un/following`,
      method: "post",
      data: { id: id },
    });
  },
};

export default TaskAPI;

/**
 * 分页查询对象
 */
export interface TaskPageQuery extends PageQuery {
  /** 搜索关键字 */
  name?: string;
  projectId?: number;
  nodeId?: number;
  status?: number;
  following?: boolean;
}

export interface TaskPageVO {
  /** 用户ID */
  id: number;
  createdAt?: Date;
  updatedAt?: Date;
  name?: string;
  spec?: string;
  projectName?: string;
  nodeName?: number;
  isFollowing: boolean; // 用户是否已关注
  isSingle: number;
  status: number;
}

export interface TaskForm {
  /** 用户ID */
  id?: number;
  name?: string;
  spec?: string;
  status?: number;
  isSingle?: number;
  projectId?: number;
  nodeId?: number;
  describe?: string;
}

export interface TaskDetail {
  id: number;
  name: string;
  spec: string;
  status: number;
  isSingle: number;
  projectName: string;
  nodeName: string;
  describe: string;
  pid: number;
  scheduleState: number;
}
