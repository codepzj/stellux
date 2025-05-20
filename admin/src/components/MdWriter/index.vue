<template>
  <div class="h-full">
    <div class="md-editor markdown-body">
      <Editor
        :value="content"
        :locale="zhHans"
        :plugins="mdPlugins"
        @change="content = $event"
        :mode="props.mode"
        placeholder="请输入内容..."
      />
    </div>
    <PhotoSelect
      v-model:open="photoSelectOpen"
      @selected-picture="selectedPicture"
    />
  </div>
</template>

<script lang="ts" setup>
// @ts-ignore
import { Editor } from "@bytemd/vue-next";
import zhHans from "bytemd/locales/zh_Hans.json";

import gfm from "@bytemd/plugin-gfm";
import gemoji from "@bytemd/plugin-gemoji";
import highlight from "@bytemd/plugin-highlight";
import frontmatter from "@bytemd/plugin-frontmatter";
import mediumZoom from "@bytemd/plugin-medium-zoom";
import mermaid from "@bytemd/plugin-mermaid";
import breaks from "@bytemd/plugin-breaks";
import photoSelect from "./plugins/photoSelect";
import { useVModel } from "@vueuse/core";
import PhotoSelect from "@/components/PhotoSelect/index.vue";

const photoSelectOpen = ref(false);

const props = withDefaults(
  defineProps<{
    content: string;
    mode?: "split" | "tab" | "auto";
  }>(),
  {
    mode: "auto",
  }
);

const emit = defineEmits<{
  (e: "update:content", value: string): void;
  (e: "update:mode", value: "split" | "tab" | "auto"): void;
}>();

const content = useVModel(props, "content", emit);
const selectedPicture = (picture: string) => {
  // 如果图片为空，则不添加
  if (!picture) return;
  content.value = content.value + `![图片](${picture})`;
};

const mdPlugins = ref([
  gfm({
    locale: {
      strike: "删除线",
      strikeText: "删除线",
      task: "任务",
      taskText: "任务",
      table: "表格",
      tableHeading: "表格标题",
    },
  }),
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
  photoSelect(photoSelectOpen),
]);
</script>

<style lang="scss" scoped></style>
