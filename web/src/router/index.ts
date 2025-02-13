import type { App } from "vue";
import { createRouter, createWebHashHistory, type RouteRecordRaw } from "vue-router";

export const Layout = () => import("@/layout/index.vue");

export const asyncRoutes = [];
// 静态路由
export const constantRoutes: RouteRecordRaw[] = [
  {
    path: "/redirect",
    component: Layout,
    meta: { hidden: true },
    children: [
      {
        path: "/redirect/:path(.*)",
        component: () => import("@/views/redirect/index.vue"),
      },
    ],
  },

  {
    path: "/login",
    component: () => import("@/views/login/index.vue"),
    meta: { hidden: true },
  },
  {
    path: "/",
    name: "/",
    component: Layout,
    redirect: "/dashboard",
    children: [
      {
        path: "/profile",
        name: "profile",
        component: () => import("@/views/profile/index.vue"),
        meta: { hidden: true },
      },
      {
        path: "dashboard",
        component: () => import("@/views/dashboard/index.vue"),
        name: "Dashboard",
        meta: {
          title: "首页",
          icon: "homepage",
          affix: true,
          keepAlive: true,
        },
      },
      {
        path: "401",
        component: () => import("@/views/error/401.vue"),
        meta: { hidden: true },
      },
      {
        path: "404",
        component: () => import("@/views/error/404.vue"),
        meta: { hidden: true },
      },
    ],
  },
  {
    path: "/system",
    component: Layout,
    redirect: "/system/user",
    name: "/system",
    meta: {
      title: "系统管理",
      icon: "system",
      hidden: false,
      alwaysShow: false,
      params: null,
    },
    children: [
      {
        path: "user",
        component: () => import("@/views/system/user/index.vue"),
        name: "User",
        meta: {
          title: "用户管理",
          icon: "el-icon-User",
          hidden: false,
          keepAlive: true,
          alwaysShow: false,
          params: null,
        },
      },
      {
        path: "logs",
        component: () => import("@/views/system/logs/index.vue"),
        name: "UserLogs",
        meta: {
          title: "用户日志",
          icon: "el-icon-document",
          hidden: false,
          keepAlive: true,
          alwaysShow: false,
          params: null,
        },
      },
      {
        path: "project/list",
        component: () => import("@/views/project/index.vue"),
        name: "ProjectList",
        meta: {
          title: "项目管理",
          icon: "el-icon-document",
          hidden: false,
          keepAlive: true,
          alwaysShow: false,
          params: null,
        },
      },
      {
        path: "list",
        component: () => import("@/views/node/list/index.vue"),
        name: "NodeList",
        meta: {
          title: "节点列表",
          icon: "el-icon-document",
          hidden: false,
          keepAlive: true,
          alwaysShow: false,
          params: null,
        },
      },
    ],
  },
  {
    path: "/task",
    component: Layout,
    redirect: "/task/task",
    name: "/task",
    meta: {
      title: "任务管理",
      icon: "el-icon-Operation",
      hidden: false,
      alwaysShow: true,
      params: null,
    },
    children: [
      {
        path: "task",
        component: () => import("@/views/task/task/index.vue"),
        name: "Task",
        meta: {
          title: "任务管理",
          // icon: "el-icon-User",
          hidden: false,
          keepAlive: true,
          alwaysShow: false,
          params: null,
        },
      },
      {
        path: "process",
        component: () => import("@/views/task/process/index.vue"),
        name: "Process",
        meta: {
          title: "进程管理",
          // icon: "el-icon-User",
          hidden: false,
          keepAlive: true,
          alwaysShow: false,
          params: null,
        },
      },
    ],
  },
  {
    path: "/logs",
    component: Layout,
    redirect: "/task/logs",
    name: "/task/logs",
    meta: {
      title: "日志管理",
      icon: "el-icon-Operation",
      hidden: false,
      alwaysShow: true,
      params: null,
    },
    children: [
      {
        path: "task/logs",
        component: () => import("@/views/task/task/logs.vue"),
        name: "TaskLogs",
        meta: {
          title: "任务日志",
          // icon: "el-icon-User",
          hidden: false,
          keepAlive: true,
          alwaysShow: false,
          params: null,
        },
      },
      {
        path: "process/logs",
        component: () => import("@/views/task/process/logs.vue"),
        name: "ProcessLogs",
        meta: {
          title: "进程日志",
          // icon: "el-icon-User",
          hidden: false,
          keepAlive: true,
          alwaysShow: false,
          params: null,
        },
      },
    ],
  },
];

/**
 * 创建路由
 */
const router = createRouter({
  history: createWebHashHistory(),
  routes: constantRoutes,
  // 刷新时，滚动条位置还原
  scrollBehavior: () => ({ left: 0, top: 0 }),
});

// 全局注册 router
export function setupRouter(app: App<Element>) {
  app.use(router);
}

export default router;
