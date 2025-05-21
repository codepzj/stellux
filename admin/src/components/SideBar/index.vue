<template>
  <div class="flex flex-col h-full">
    <div
      :class="[
        'flex items-center',
        sidebarStore.collapse ? 'justify-center' : 'justify-start',
      ]"
    >
      <img
        :src="logoSrc"
        :key="logoSrc"
        alt="logo"
        @load="onLoad"
        @click="$router.push({ name: 'Dashboard' })"
        :class="[
          'transition-opacity duration-700 ease-in-out my-2 cursor-pointer',
          sidebarStore.collapse ? 'mx-auto size-8' : 'ml-4 w-32',
          logoLoaded ? 'opacity-100' : 'opacity-0',
        ]"
      />
    </div>
    <a-menu
      mode="inline"
      :selectedKeys="selectedKeys"
      :openKeys="openKeys"
      :items="menuItems"
      @select="onSelect"
      @openChange="onOpenChange"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { routes } from "@/router";
import type { VNodeChild } from "vue";
import type { ItemType } from "ant-design-vue";
import { useSidebarStore, useSystemStore } from "@/store";

const sidebarStore = useSidebarStore();
const systemStore = useSystemStore();
const router = useRouter();
const route = useRoute();
const selectedKeys = ref<string[]>([]);
const openKeys = ref<string[]>([]);

const logoLoaded = ref(false);
const logoSrc = computed(
  () =>
    `/logo${sidebarStore.collapse ? "-sm" : ""}-${systemStore.themeMode === "dark" ? "dark" : "light"}.png`
);
const onLoad = () => {
  logoLoaded.value = true;
};
watch(logoSrc, () => {
  logoLoaded.value = false;
});

type MenuItemType = {
  key: string;
  label: string;
  children?: MenuItemType[];
  icon?: () => VNodeChild;
  type?: "group";
} & Omit<ItemType, "key" | "label" | "children" | "icon" | "type">;

function getItem(
  label: string,
  key: string,
  icon?: () => VNodeChild,
  children?: MenuItemType[],
  type?: "group"
): MenuItemType {
  return { key, icon, children, label, type };
}

const mainRoute = routes.find(r => r.path === "/");

const menuItems = computed(() => {
  if (!mainRoute?.children) return [];
  const generateMenuItems = (route: any): MenuItemType | null => {
    if (!route.name || !route.meta?.title || route.meta.hideInSideBar === true)
      return null;
    const children = route.children
      ?.map(generateMenuItems)
      .filter((i: any): i is MenuItemType => i !== null);
    return getItem(
      route.meta.title,
      route.name.toString(),
      route.meta.icon,
      children?.length ? children : undefined
    );
  };
  return mainRoute.children
    .map(generateMenuItems)
    .filter((i): i is MenuItemType => i !== null);
});

const rootSubmenuKeys = computed(() =>
  menuItems.value.filter(item => item.children?.length).map(item => item.key)
);

const onSelect = ({ key }: { key: string }) => {
  selectedKeys.value = [key];
  router.push({ name: key });
};

const onOpenChange = (keys: string[]) => {
  const latest = keys.find(key => !openKeys.value.includes(key));
  if (latest && rootSubmenuKeys.value.includes(latest)) {
    openKeys.value = [latest];
  } else {
    openKeys.value = keys;
  }
};

const findRouteNameForHighlight = (
  name: string | symbol | undefined,
  routeList: any[],
  parentName?: string
): string | undefined => {
  for (const route of routeList) {
    if (route.name === name) {
      return route.meta?.hideInSideBar ? parentName : route.name?.toString();
    }
    if (route.children) {
      const found = findRouteNameForHighlight(
        name,
        route.children,
        route.meta?.hideInSideBar ? parentName : route.name?.toString()
      );
      if (found) return found;
    }
  }
  return undefined;
};

watch(
  () => route.name,
  () => {
    const highlightName = findRouteNameForHighlight(route.name, routes);
    selectedKeys.value = highlightName ? [highlightName] : [];
    const segments = route.path.split("/").filter(Boolean);
    let currentPath = "";
    const keys: string[] = [];
    const baseRoute = routes.find(r => r.path === "/");
    segments.forEach(seg => {
      currentPath += `/${seg}`;
      const parent = baseRoute?.children?.find(
        r => r.path === currentPath && r.name && r.meta?.hideInSideBar !== true
      );
      if (parent?.name) keys.push(parent.name.toString());
    });
    if (keys.length && keys[keys.length - 1] === route.name?.toString())
      keys.pop();
  },
  { immediate: true }
);
</script>
