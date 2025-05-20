<template>
  <Tooltip title="搜索" placement="bottom" :mouse-enter-delay="0.5">
    <slot
      ><SvgIcon
        name="search"
        class="cursor-pointer rounded-md p-1"
        :size="28"
        @click="visible = true"
    /></slot>
    <a-modal
      v-model:open="visible"
      title="搜索菜单"
      :keyboard="false"
      @cancel="handleClose"
    >
      <a-input
        ref="inputRef"
        v-model:value="keyword"
        clearable
        placeholder="请输入关键词搜索"
        @change="handleSearch"
      >
        <template #prefix>
          <SvgIcon name="search" :size="16" />
        </template>
      </a-input>
      <div class="mt-20px">
        <Empty v-if="resultOptions.length === 0" description="暂无搜索结果" />
        <search-result
          v-else
          v-model:value="activePath"
          :options="resultOptions"
          @enter="handleEnter"
        />
      </div>
      <template #footer>
        <search-footer />
      </template>
    </a-modal>
  </Tooltip>
</template>

<script lang="ts" setup>
import { ref, shallowRef, computed, watch, nextTick } from "vue";
import { useDebounceFn, onKeyStroke } from "@vueuse/core";
import SvgIcon from "@/components/SvgIcon/index.vue";
import { Empty, Tooltip } from "ant-design-vue";
import SearchResult from "./components/SearchResult.vue";
import SearchFooter from "./components/SearchFooter.vue";
import type { RouteRecordRaw } from "vue-router";
import router, { routes } from "@/router";

const keyword = ref("");
const activePath = ref("");
const menusList = computed(() => transformRouteToList(routes));
const resultOptions = shallowRef<RouteRecordRaw[]>([]);
const inputRef = ref<HTMLInputElement | null>(null);
const handleSearch = useDebounceFn(search, 300);

const visible = ref(false);

watch(visible, async val => {
  if (val) {
    /** 自动聚焦 */
    await nextTick();
    inputRef.value?.focus();
  }
});

/** 查询 */
function search() {
  resultOptions.value = menusList.value.filter(menu => {
    const title = menu.meta?.title as string;
    return (
      keyword.value &&
      title
        ?.toLocaleLowerCase()
        .includes(keyword.value.toLocaleLowerCase().trim())
    );
  });
  if (resultOptions.value?.length > 0) {
    activePath.value = resultOptions.value[0].name as string;
  } else {
    activePath.value = "";
  }
}

/** 将路由转换成菜单列表 */
function transformRouteToList(
  routes: RouteRecordRaw[],
  treeMap: RouteRecordRaw[] = []
) {
  if (routes && routes.length === 0) return [];
  return routes.reduce((acc, cur) => {
    /** 允许在菜单内显示并且无子路由，且 meta.hidden 为 false 或 undefined */
    if (!cur.meta?.hideInMenu && !cur.meta?.hidden && !cur.children) {
      acc.push(cur);
    }
    if (cur.children && cur.children.length > 0) {
      transformRouteToList(cur.children, treeMap);
    }
    return acc;
  }, treeMap);
}

function handleClose() {
  visible.value = false;
  /** 延时处理防止用户看到某些操作 */
  setTimeout(() => {
    resultOptions.value = [];
    keyword.value = "";
  }, 200);
}

/** key up */
function handleUp() {
  const { length } = resultOptions.value;
  if (length === 0) return;
  const index = resultOptions.value.findIndex(
    item => item.name === activePath.value
  );
  if (index === 0) {
    activePath.value = resultOptions.value[length - 1].name as string;
  } else {
    activePath.value = resultOptions.value[index - 1].name as string;
  }
}

/** key down */
function handleDown() {
  const { length } = resultOptions.value;
  if (length === 0) return;
  const index = resultOptions.value.findIndex(
    item => item.name === activePath.value
  );
  if (index + 1 === length) {
    activePath.value = resultOptions.value[0].name as string;
  } else {
    activePath.value = resultOptions.value[index + 1].name as string;
  }
}

/** key enter */
function handleEnter() {
  const { length } = resultOptions.value;
  if (length === 0 || activePath.value === "") return;
  if (/http(s)?:/.test(activePath.value)) {
    window.open(activePath.value);
  } else {
    router.push({ name: activePath.value });
    handleClose();
  }
}

onKeyStroke("Escape", handleClose);
onKeyStroke("Enter", handleEnter);
onKeyStroke("ArrowUp", handleUp);
onKeyStroke("ArrowDown", handleDown);
</script>
<style lang="scss" scoped></style>
