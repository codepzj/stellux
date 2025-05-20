<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter, type RouteRecordRaw } from "vue-router";

defineOptions({
  name: "BreadCrumb",
});

const router = useRouter();
const route = useRoute();

// 点击菜单
const clickMenuItem = (menuItem: RouteRecordRaw) => {
  const { isExt, extOpenMode, type } = menuItem?.meta || {};

  if (menuItem.name === "__index") {
    router.push({ name: "Dashboard" });
    return;
  }

  if (type === 0 && !menuItem.redirect) return;

  if (isExt && extOpenMode === 1) {
    window.open(menuItem.path);
  } else {
    const to =
      typeof menuItem.redirect === "string" ? menuItem.redirect : menuItem;
    router.push(to);
  }
};

// 构建面包屑数据
const menus = computed(() => {
  const matched = route.matched.filter(
    r => r.meta?.title && r.meta?.hideInBreadcrumb !== true
  );

  const lastRoute = route.matched[route.matched.length - 1];

  // 如果当前路由设置了 hideInBreadcrumb 但我们想显示标题
  if (
    lastRoute &&
    lastRoute.meta?.title &&
    lastRoute.meta?.hideInBreadcrumb === true &&
    !matched.some(m => m.name === lastRoute.name)
  ) {
    matched.push(lastRoute);
  }

  return matched.map(routeItem => {
    const children = (routeItem.children || []).filter(
      c => c.meta?.hideInBreadcrumb !== true
    );
    return { ...routeItem, children };
  });
});

// 获取下拉菜单选中项
const getSelectKeys = (routeIndex: number) => {
  return [menus.value[routeIndex + 1]?.name] as string[];
};
</script>

<template>
  <a-breadcrumb>
    <template v-for="(routeItem, routeIndex) in menus" :key="routeItem?.name">
      <a-breadcrumb-item>
        <span class="cursor-pointer" @click="clickMenuItem(routeItem)">
          {{ routeItem?.meta?.title }}
        </span>
        <!-- 有子项才显示下拉菜单 -->
        <template v-if="routeItem?.children?.length" #overlay>
          <a-menu :selected-keys="getSelectKeys(routeIndex)">
            <template
              v-for="childItem in routeItem?.children"
              :key="childItem.name"
            >
              <a-menu-item
                v-if="
                  !childItem.meta?.hideInSidebar &&
                  !childItem.meta?.hideInBreadcrumb
                "
                :key="childItem.name"
                @click="clickMenuItem(childItem)"
              >
                {{ childItem.meta?.title }}
              </a-menu-item>
            </template>
          </a-menu>
        </template>
      </a-breadcrumb-item>
    </template>
  </a-breadcrumb>
</template>
