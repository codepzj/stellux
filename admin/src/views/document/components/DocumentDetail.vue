<template>
  <div class="h-full">
    <div class="flex justify-between items-center">
      <div class="w-[300px] md:w-2/3 mb-4">
        <a-input
          v-model:value="title"
          placeholder="请输入标题"
          addon-before="标题"
          show-count
          :maxlength="50"
        />
      </div>
    </div>
    <MdWriter v-model:content="content" mode="auto" />
  </div>
</template>

<script setup lang="ts">
import MdWriter from "@/components/MdWriter/index.vue";
import { useVModel } from "@vueuse/core";

const props = defineProps<{
  title: string;
}>();

const emit = defineEmits<{
  (e: "update:title", value: string): void;
}>();

const title = useVModel(props, "title", emit);
const content = ref("");
</script>
