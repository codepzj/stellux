<!-- 选择图片组件，用于选择图片或上传并返回图片的url -->
<template>
  <div>
    <a-modal
      width="100%"
      wrap-class-name="full-modal"
      :open="thumbnailModalOpen"
      @update:open="fetchThumbnails"
      title="选择图片"
      @ok="confirmSelectedPicture"
      @cancel="thumbnailModalOpen = false"
    >
      <div class="flex justify-end my-4">
        <Uploader @close="fetchThumbnails" />
      </div>
      <PhotoWall
        mode="picture-card"
        type="select"
        :list="page.list"
        @selected-photos="handleSelectPicture"
      />
      <div class="flex justify-end my-4">
        <!-- 分页器 只有当总页数大于1时才显示 -->
        <a-pagination
          v-if="page.total_count > page.page_size"
          v-model:current="page.page_no"
          :total="page.total_count"
          show-less-items
          @change="fetchThumbnails"
        />
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { queryFileList } from "@/api/file";
import PhotoWall from "@/components/PhotoWall/index.vue";
import type { PageData } from "@/types/dto";
import type { FileVO } from "@/types/file";
import { useVModel } from "@vueuse/core";
import Uploader from "@/components/Uploader/index.vue"; // 图片上传组件

/**
 * Props 接收外部传入的 open 属性，控制 Modal 是否打开
 */
const props = defineProps<{
  open: boolean;
}>();

/**
 * Emits 事件定义
 * @event update:open - 通知父组件更新 Modal 的打开状态
 * @event selected-picture - 通知父组件已选择的图片 URL
 */
const emit = defineEmits<{
  /**
   * 更新 Modal 打开状态
   * @param open - 是否打开
   */
  (e: "update:open", open: boolean): void;

  /**
   * 返回已选中的图片 URL
   * @param picture - 图片 URL 字符串
   */
  (e: "selected-picture", picture: string): void;
}>();

/**
 * thumbnailModalOpen 绑定 open 属性，实现 Modal 显隐的双向绑定
 */
const thumbnailModalOpen = useVModel(props, "open", emit);

/**
 * 分页数据对象，包含当前页码、每页条数、总数、总页数和图片列表
 */
const page = reactive<PageData<FileVO>>({
  page_no: 1,
  page_size: 10,
  total_count: 0,
  total_page: 0,
  list: [],
});

const selectedPicture = ref("");
/**
 * 用户选择图片时触发，只返回第一个图片 URL 给父组件
 * @param pictures - 用户选择的图片 URL 数组
 */
const handleSelectPicture = (pictures: string[]) => {
  selectedPicture.value = pictures[0];
};

const confirmSelectedPicture = () => {
  emit("selected-picture", selectedPicture.value);
  // 清空选中的图片
  selectedPicture.value = "";
  thumbnailModalOpen.value = false;
};

/**
 * 拉取缩略图列表（分页）
 * 根据当前页码与分页大小调用接口，并更新图片列表与总数
 */
const fetchThumbnails = async () => {
  const res = await queryFileList({
    page_no: page.page_no,     // 当前页码
    page_size: page.page_size, // 每页数量
  });

  page.list = res.data.list;
  page.total_count = res.data.total_count;
};

watch(thumbnailModalOpen, (newVal) => {
  // 如果打开，则拉取图片列表
  if (newVal) {
    fetchThumbnails();
  }
});
</script>
