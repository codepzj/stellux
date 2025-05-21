<template>
  <div class="flex h-14 w-full items-center justify-between px-4 dark:bg-dark">
    <!-- 左侧 Logo 区域 -->
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
          sidebarStore.collapse ? 'mx-auto' : 'w-8',
          logoLoaded ? 'opacity-100' : 'opacity-0',
        ]"
      />
      <!-- 知识库按钮 -->
      <div
        class="text-sm font-bold cursor-pointer ml-2 transition-all duration-300 ease-in-out dark:text-white hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-md px-2 py-1"
        @click="$router.push({ name: 'DocumentOverview' })"
      >
        知识库
      </div>
      <!-- 编辑 / 预览切换 -->
      <div
        v-if="$route.name === 'DocumentContent'"
        class="h-full flex items-center"
      >
        <a-select
          v-model:value="documentStore.mode"
          class="w-24 text-center"
          :bordered="false"
        >
          <a-select-option value="preview" @click="setMode('preview')">
            <span class="text-xs">预览模式</span>
          </a-select-option>
          <a-select-option value="edit" @click="setMode('edit')">
            <span class="text-xs">编辑模式</span>
          </a-select-option>
        </a-select>
      </div>
      <div
        v-if="$route.name === 'DocumentContent' && documentStore.editDocumentTitle"
        class="ml-6 flex items-center gap-1 text-sm text-zinc-600 dark:text-zinc-300 px-2 py-1 rounded hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-all duration-300"
        :title="documentStore.editDocumentTitle"
      >
        ✏️
        <span class="font-medium text-zinc-700 dark:text-white">
          {{ documentStore.editDocumentTitle }}
        </span>
      </div>
    </div>

    <!-- 右侧操作区域 -->
    <div class="flex items-center gap-2">
      <!-- 编辑操作按钮 -->
      <div v-if="$route.name === 'DocumentContent'">
        <a-button
          v-for="button in documentStore.actionButtonList"
          :disabled="button.disabled"
          :key="button.name"
          type="primary"
          @click="button.action"
        >
          {{ button.name }}
        </a-button>
      </div>

      <!-- 全屏/设置按钮 -->
      <FullScreen class="hidden md:block" />
      <Setting class="hidden md:block" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { FullScreen, Setting } from "./components";
import { useSidebarStore, useSystemStore } from "@/store";
import { useDocumentStore } from "@/store/modules/document";

const documentStore = useDocumentStore();
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

const setMode = (newMode: "preview" | "edit") => {
  documentStore.setMode(newMode);
};
</script>
