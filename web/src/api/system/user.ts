import request from "@/utils/request";

const UserAPI = {
  /**
   * 获取当前登录用户信息
   *
   * @returns 登录用户昵称、头像信息，包括角色和权限
   */
  getInfo() {
    return request<any, UserInfo>({
      url: `/my/info`,
      method: "get",
    });
  },

  /**
   * 获取用户分页列表
   *
   * @param queryParams 查询参数
   */
  getPage(queryParams: UserPageQuery) {
    return request<any, PageResult<UserForm[]>>({
      // url: `${USER_BASE_URL}/page`,
      url: `/user/list`,
      method: "get",
      params: queryParams,
    });
  },

  /**
   * 添加用户
   *
   * @param data 用户表单数据
   */
  add(data: UserForm) {
    return request({
      url: `/user/add`,
      method: "post",
      data: data,
    });
  },

  /**
   * 修改用户
   *
   * @param id 用户ID
   * @param data 用户表单数据
   */
  update(data: UserForm) {
    return request({
      url: `/user/save`,
      method: "post",
      data: data,
    });
  },

  /**
   * 修改用户密码
   *
   * @param id 用户ID
   * @param password 新密码
   */
  resetPassword(id: number | undefined, password: string) {
    return request({
      url: `/user/password/reset`,
      method: "post",
      data: { password: password, id: id },
    });
  },

  /**
   * 删除用户
   *
   * @param id
   */
  deleteById(id: number) {
    return request({
      url: `/user/delete`,
      method: "post",
      data: { id: id },
    });
  },

  /** 获取个人中心用户信息 */
  getProfile() {
    return request<any, UserProfileVO>({
      url: `/my/info`,
      method: "get",
    });
  },

  /** 修改个人中心用户信息 */
  updateAvatar(data: UserProfileForm) {
    return request({
      url: `/save/my/avatar`,
      method: "post",
      data: data,
    });
  },

  /** 修改个人中心用户信息 */
  updateNickname(data: UserProfileForm) {
    return request({
      url: `/save/my/nickname`,
      method: "post",
      data: data,
    });
  },

  /** 修改个人中心用户密码 */
  changePassword(data: PasswordChangeForm) {
    return request({
      url: `/save/my/password`,
      method: "post",
      data: data,
    });
  },
};

export default UserAPI;

/** 登录用户信息 */
export interface UserInfo {
  /** 用户ID */
  id?: number;

  /** 用户名 */
  username?: string;

  /** 昵称 */
  nickname?: string;

  /** 头像URL */
  avatar?: string;

  /** 角色 */
  roleId: number;
}

/**
 * 用户分页查询对象
 */
export interface UserPageQuery extends PageQuery {
  /** 搜索关键字 */
  username?: string;
  nickname?: string;
  roleId?: number;

  /** 用户状态 */
  status?: number;
}

/** 用户分页对象 */
export interface UserPageVO {
  /** 用户ID */
  id: number;
  /** 更新时间 */
  updatedAt?: Date;
  /** 用户昵称 */
  nickname?: string;
  /** 用户状态(1:启用;0:禁用) */
  status?: number;
  /** 用户名 */
  username?: string;
  roleId?: number;
  password?: string;
}

/** 用户表单类型 */
export interface UserForm {
  /** 用户ID */
  id?: number;
  /** 用户名 */
  username?: string;
  /** 昵称 */
  nickname?: string;
  /** 密码 */
  password?: string;
  /**用户等级*/
  roleId?: number;
  /** 用户状态(1:正常;0:禁用) */
  status?: number;
  avatar?: string;
}

/** 个人中心用户信息 */
export interface UserProfileVO {
  /** 用户ID */
  id?: number;
  status?: number;

  /** 用户名 */
  username?: string;

  /** 昵称 */
  nickname?: string;

  /** 头像URL */
  avatar?: string;
  roleId?: number;
  roleName?: string;

  /** 创建时间 */
  createdAt?: Date;
  updatedAt?: Date;
}

/** 个人中心用户信息表单 */
export interface UserProfileForm {
  /** 用户ID */
  id?: number;

  /** 昵称 */
  nickname?: string;

  /** 头像URL */
  avatar?: string;
}

/** 修改密码表单 */
export interface PasswordChangeForm {
  /** 原密码 */
  oldPassword?: string;
  /** 新密码 */
  newPassword?: string;
  /** 确认新密码 */
  confirmPassword?: string;
}
