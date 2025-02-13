import request from "@/utils/request";

const SelectorsAPI = {
  getUserSelectors() {
    return request<UserSelector[]>({
      // url: `${USER_BASE_URL}/page`,
      url: `/user/selector`,
      method: "get",
    });
  },

  getNodeSelectors() {
    return request<UserSelector[]>({
      // url: `${USER_BASE_URL}/page`,
      url: `/node/selector`,
      method: "get",
    });
  },

  getProjectSelectors() {
    return request<UserSelector[]>({
      // url: `${USER_BASE_URL}/page`,
      url: `/project/selector`,
      method: "get",
    });
  },
};

export default SelectorsAPI;

export interface UserSelector {
  /** 用户ID */
  id?: number;
  username?: string;
}

export interface NodeSelector {
  id: number;
  name: string;
}

export interface ProjectSelector {
  /** 用户ID */
  id: number;
  name: string;
}



