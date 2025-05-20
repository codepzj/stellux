<template>
  <div>
    <div class="flex justify-between items-center my-4 gap-2">
      <div class="w-[300px] md:w-[50%]">
        <a-input
          v-model:value="postForm.title"
          placeholder="请输入标题"
          addon-before="标题"
          show-count
          :maxlength="50"
        />
      </div>
      <a-button type="primary" @click="open = true"> 发布文章 </a-button>
    </div>

    <MdWriter v-model:content="postForm.content" />

    <PostFormDrawer
      :mode="props.mode"
      v-model:open="open"
      v-model:post-form="postForm"
    />
  </div>
</template>

<script lang="ts" setup>
import MdWriter from "@/components/MdWriter/index.vue";

import PostFormDrawer from "./PostFormDrawer.vue";
import type { PostReq } from "@/types/post";
import { useVModel } from "@vueuse/core";

const open = ref(false);

const props = defineProps<{
  mode: "create" | "edit";
  postForm: PostReq;
}>();

const emit = defineEmits<{
  (e: "update:postForm", value: PostReq): void;
}>();

const postForm = useVModel(props, "postForm", emit);
</script>

<style lang="scss" scoped></style>
