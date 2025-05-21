<template>
  <div class="h-full px-8 overflow-x-hidden">
    <Viewer v-model:value="content" :plugins="mdPlugins" />
  </div>
</template>

<script setup lang="ts">
// @ts-ignore
import { Viewer } from "@bytemd/vue-next";
import gfm from "@bytemd/plugin-gfm";
import gemoji from "@bytemd/plugin-gemoji";
import highlight from "@bytemd/plugin-highlight";
import frontmatter from "@bytemd/plugin-frontmatter";
import mediumZoom from "@bytemd/plugin-medium-zoom";
import breaks from "@bytemd/plugin-breaks";
import mermaid from "@bytemd/plugin-mermaid";
import { useVModel } from "@vueuse/core";

const props = defineProps<{
  content: string;
}>();

const emit = defineEmits<{
  (e: "update:content", value: string): void;
}>();

const content = useVModel(props, "content", emit);
const mdPlugins = ref([
  gfm(),
  gemoji(),
  highlight(),
  frontmatter(),
  mediumZoom(),
  breaks(),
  mermaid({
    locale: {
      mermaid: "图表",
      flowchart: "流程图",
      sequence: "时序图",
      class: "类图",
      state: "状态图",
      er: "实体关系图",
      uj: "用户旅程图",
      gantt: "甘特图",
      pie: "饼图",
      mindmap: "思维导图",
      timeline: "时间线",
    },
  }),
]);
</script>
<style lang="scss" scoped></style>