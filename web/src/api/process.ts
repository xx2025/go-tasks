import request from "@/utils/request";
import {TaskDetail} from "@/api/task";

const ProcessAPI = {
  /**
   *
   * @param queryParams 查询参数
   */
  getPage(queryParams: ProcessPageQuery) {
    return request<any, PageResult<ProcessPageVO[]>>({
      // url: `${USER_BASE_URL}/page`,
      url: `/process/list`,
      method: "get",
      params: queryParams,
    });
  },

  getDetail(queryParams: {id: number}) {
    return request<any, ProcessDetail>({
      // url: `${USER_BASE_URL}/page`,
      url: `/process/detail`,
      method: "get",
      params: queryParams,
    });
  },

  /**
   *
   * @param data 表单数据
   */
  add(data: ProcessForm) {
    return request({
      url: `/process/save`,
      method: "post",
      data: data,
    });
  },

  /**
   *
   * @param data 表单数据
   */
  update(data: ProcessForm) {
    return request({
      url: `/process/save`,
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
      url: `/process/delete`,
      method: "post",
      data: { id: id },
    });
  },

  startById(id: number) {
    return request({
      url: `/process/start`,
      method: "post",
      data: { id: id },
    });
  },

  stopById(id: number) {
    return request({
      url: `/process/stop`,
      method: "post",
      data: { id: id },
    });
  },

  followingById(id: number) {
    return request({
      url: `/process/following`,
      method: "post",
      data: { id: id },
    });
  },
  unFollowingById(id: number) {
    return request({
      url: `/process/un/following`,
      method: "post",
      data: { id: id },
    });
  },
};

export default ProcessAPI;

/**
 * 分页查询对象
 */
export interface ProcessPageQuery extends PageQuery {
  name?: string;
  projectId?: number;
  nodeId?: number;
  status?: number;
  following: boolean;
}

export interface ProcessPageVO {
  id: number;
  createdAt?: Date;
  updatedAt?: Date;
  name?: string;
  projectName?: string;
  nodeName?: number;
  isFollowing: boolean; // 用户是否已关注
}

export interface ProcessForm {
  id?: number;
  name?: string;
  status?: number;
  projectId?: number;
  nodeId?: number;
  describe?: string;
  maxRetries?: number;
}

export interface ProcessDetail {
  id: number;
  name: string;
  projectName: string;
  nodeName: string;
  describe: string;
  status: number;
  pid: number;
  runningStatus: boolean;
}

