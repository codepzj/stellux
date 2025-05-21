<template>
  <div class="h-full">
    <div class="h-full">
      <MdWriter v-if="documentStore.mode === 'edit'" v-model:content="content" mode="auto" />
      <MdViewer v-else :content="content" />
    </div>
  </div>
</template>

<script setup lang="ts">
import MdWriter from "@/components/MdWriter/index.vue";
import MdViewer from "@/components/MdViewer/index.vue";
import { useDocumentStore } from "@/store/modules/document";
import { useVModel } from "@vueuse/core";

const documentStore = useDocumentStore();
const props = defineProps<{
  content: string;
}>();

const emit = defineEmits<{
  (e: "update:content", value: string): void;
}>();

const content = useVModel(props, "content", emit);
</script>
