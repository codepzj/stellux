<template>
  <div class="flex h-14 w-full items-center justify-between px-4 dark:bg-dark">
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
          'transition-opacity cursor-pointer duration-700 ease-in-out my-[15px]',
          sidebarStore.collapse ? 'mx-auto' : 'ml-4 w-8',
          logoLoaded ? 'opacity-100' : 'opacity-0',
        ]"
      />
      <a-button
        type="default"
        class="text-sm font-bold"
        @click="$router.push({ name: 'DocumentOverview' })"
      >
        <span class="text-sm font-bold">知识库</span>
      </a-button>
    </div>
    <div class="flex items-center gap-2">
      <FullScreen class="hidden md:block" />
      <Setting class="hidden md:block" />
      <Dropdown placement="bottomRight">
        <Avatar
          class="cursor-pointer"
          :src="userStore.userInfo?.avatar"
          :alt="userStore.userInfo?.username"
        >
          {{ userStore.userInfo?.nickname }}
        </Avatar>
        <template #overlay>
          <Menu>
            <Menu.Item>
              <div
                @click.prevent="$router.push({ name: 'UserEdit' })"
                class="flex items-center gap-2"
              >
                <EditOutlined />
                编辑资料
              </div>
            </Menu.Item>
            <Menu.Item>
              <div @click.prevent="Logout" class="flex items-center gap-2">
                <PoweroffOutlined />
                退出登录
              </div>
            </Menu.Item>
          </Menu>
        </template>
      </Dropdown>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { FullScreen, Setting } from "./components";
import { useUserStore } from "@/store/modules/user";
import { Menu, Dropdown, Avatar, message } from "ant-design-vue";
import { PoweroffOutlined, EditOutlined } from "@ant-design/icons-vue";
import { clearStore } from "@/utils/clear";
import { useRouter } from "vue-router";
import { useSidebarStore, useSystemStore } from "@/store";
const userStore = useUserStore();
const router = useRouter();
const sidebarStore = useSidebarStore();
const systemStore = useSystemStore();
const logoLoaded = ref(false);
const logoSrc = computed(
  () => `/logo-sm-${systemStore.themeMode === "dark" ? "dark" : "light"}.png`
);
const onLoad = () => {
  logoLoaded.value = true;
};
watch(logoSrc, () => {
  logoLoaded.value = false;
});

function Logout() {
  clearStore();
  message.info("退出成功,请重新登录");
  router.push("/login");
}
</script>
