import {
  createWebHistory,
  createRouter,
  type RouteRecordRaw,
} from "vue-router";

import Layout from "@/layout/index.vue";
import Dashboard from "@/views/dashboard/index.vue";
import PostCreate from "@/views/post/create.vue";
import PostEdit from "@/views/post/edit.vue";
import Login from "@/views/auth/login.vue";
import {
  HomeOutlined,
  UserOutlined,
  AppstoreOutlined,
  ExperimentOutlined,
  CloudUploadOutlined,
  BookOutlined,
} from "@ant-design/icons-vue";

export const routes: RouteRecordRaw[] = [
  {
    path: "/",
    component: Layout,
    children: [
      {
        path: "",
        name: "Dashboard",
        component: Dashboard,
        meta: { title: "仪表盘", icon: () => h(HomeOutlined) },
      },
      {
        path: "/user",
        name: "UserManagement",
        meta: { title: "用户", icon: () => h(UserOutlined) },
        children: [
          {
            path: "list",
            component: () => import("@/views/user/list.vue"),
            name: "UserList",
            meta: { title: "用户列表" },
          },
          {
            path: "edit",
            name: "UserEdit",
            meta: { title: "编辑用户" },
            component: () => import("@/views/user/edit/index.vue"),
            redirect: { name: "UserEditBasic" },
            children: [
              {
                path: "basic",
                component: () => import("@/views/user/edit/pages/basic.vue"),
                name: "UserEditBasic",
                meta: { title: "基本信息", hideInSideBar: true },
              },
              {
                path: "reset-password",
                component: () =>
                  import("@/views/user/edit/pages/ResetPassword.vue"),
                name: "UserEditResetPassword",
                meta: { title: "重置密码", hideInSideBar: true },
              },
            ],
          },
        ],
      },
      {
        path: "/post",
        name: "Post",
        meta: { title: "文章", icon: () => h(AppstoreOutlined) },
        children: [
          {
            path: "create",
            component: PostCreate,
            name: "PostCreate",
            meta: {
              title: "发布文章",
              keepAlive: true,
            },
          },
          {
            path: "edit/:id",
            component: PostEdit,
            name: "PostEdit",
            meta: {
              title: "编辑文章",
              hideInBreadcrumb: true,
              hideInSideBar: true,
            },
          },
          {
            path: "list",
            component: () => import("@/views/post/list.vue"),
            name: "PostList",
            meta: {
              title: "文章列表",
            },
          },
          {
            path: "draft",
            component: () => import("@/views/post/draft.vue"),
            name: "PostDraft",
            meta: {
              title: "草稿箱",
              hideInSideBar: true,
            },
          },
          {
            path: "bin",
            component: () => import("@/views/post/bin.vue"),
            name: "PostBin",
            meta: { title: "回收箱", hideInSideBar: true },
          },
          {
            path: "label",
            name: "PostLabel",
            meta: { title: "文章标签" },
            component: () => import("@/views/label/index.vue"),
          },
        ],
      },
      {
        path: "document",
        name: "Document",
        meta: { title: "文档", icon: () => h(BookOutlined) },
        redirect: { name: "DocumentOverview" },
      },

      {
        path: "file",
        name: "File",
        meta: {
          title: "附件",
          icon: () => h(CloudUploadOutlined),
        },
        component: () => import("@/views/file/index.vue"),
      },
      {
        path: "/test",
        component: () => import("@/views/test/index.vue"),
        name: "Test",
        meta: { title: "测试", icon: () => h(ExperimentOutlined) },
      },
    ],
  },
  {
    path: "/document",
    name: "Document",
    component: () => import("@/views/document/index.vue"),
    redirect: { name: "DocumentOverview" },
    children: [
      {
        path: "overview",
        name: "DocumentOverview",
        component: () => import("@/views/document/overview.vue"),
      },
      {
        path: "content/:id",
        name: "DocumentContent",
        component: () => import("@/views/document/content.vue"),
      },
    ],
  },
  {
    path: "/login",
    component: Login,
    name: "Login",
  },
  {
    path: "/error",
    redirect: { name: "404" },
    children: ["401", "403", "404", "500"].map(code => {
      return {
        path: code,
        name: code,
        component: () => import(`@/views/auth/error.vue`),
      };
    }),
  },
  {
    path: "/:pathMatch(.*)",
    redirect: { name: "404" },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

function getCacheNames(routes: RouteRecordRaw[]): string[] {
  return routes.reduce((names: string[], route) => {
    if (route.meta?.keepAlive) {
      names.push(route.name as string);
    }
    if (route.children) {
      names.push(...getCacheNames(route.children));
    }
    return names;
  }, []);
}

export const cacheNames = getCacheNames(routes);

export default router;
