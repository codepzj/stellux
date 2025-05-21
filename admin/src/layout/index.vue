<template>
  <a-layout class="h-screen flex flex-col overflow-hidden">
    <a-layout-sider
      width="180"
      :trigger="null"
      collapsible
      v-model:collapsed="sidebarStore.collapse"
    >
      <!-- 侧边栏内容 -->
      <SideBar :collapsed="sidebarStore.collapse" />
    </a-layout-sider>

    <a-layout class="flex-1">
      <a-layout-header class="!h-12 !px-0">
        <Header />
      </a-layout-header>
      <a-layout-content
        class="p-4 overflow-y-auto overflow-x-hidden"
      >
        <Main />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import SideBar from "@/components/SideBar/index.vue";
import Header from "@/components/Header/index.vue";
import Main from "./main/index.vue";
import { useSidebarStore, useUserStore } from "@/store";
import { getUserInfoAPI } from "@/api/user";
import { useWindowSize } from "@vueuse/core";

// 判断是否为移动端设备
const { width } = useWindowSize();

const userStore = useUserStore();

const sidebarStore = useSidebarStore();

// 加载界面时初始化用户信息
const getUserInfo = async () => {
  const res = await getUserInfoAPI();
  userStore.setUserInfo(res.data);
};

onMounted(async () => {
  await getUserInfo();
});

watch(width, newWidth => {
  sidebarStore.setCollapse(newWidth < 768);
});
</script>

<style lang="scss">
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 5px;
}

::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>
