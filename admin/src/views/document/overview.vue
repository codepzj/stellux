<template>
  <div class="h-full p-4 overflow-auto">
    <div class="flex justify-end">
      <a-button type="primary" @click="showModal = true">新增文档</a-button>
    </div>
    <a-skeleton :loading="loading" active>
      <div class="py-6">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <a-card
            v-for="doc in docList"
            :key="doc.id"
            class="cursor-pointer hover:translate-y-[-5px] transition-all duration-300"
            @click="
              $router.push({ name: 'DocumentContent', params: { id: doc.id } })
            "
          >
            <div class="flex justify-start items-center gap-2 px-2 py-4">
              <SvgIcon name="document" />
              <span class="text-lg font-bold">{{ doc.title }}</span>
              <SvgIcon name="global" :size="14" v-if="doc.is_public" />
              <SvgIcon name="lock" :size="14" v-else />
            </div>
            <div class="flex justify-start gap-2 px-2 py-4">
              <div
                class="text-sm text-zinc-500 dark:text-zinc-200 line-clamp-2"
              >
                {{ doc.description }}
              </div>
            </div>
          </a-card>
        </div>
      </div>
    </a-skeleton>

    <!-- 新增文档弹窗 -->
    <a-modal
      v-model:open="showModal"
      title="新增文档"
      @ok="handleOk"
      @cancel="handleCancel"
    >
      <a-form layout="vertical" :model="createDoc" :rules="rules" ref="formRef">
        <a-form-item label="标题" name="title">
          <a-input v-model:value="createDoc.title" placeholder="请输入标题" />
        </a-form-item>
        <a-form-item label="描述" name="description">
          <a-input
            v-model:value="createDoc.description"
            placeholder="请输入描述"
          />
        </a-form-item>
        <a-form-item label="封面" name="thumbnail">
          <div class="flex flex-row gap-2 items-end">
            <div
              v-if="createDoc.thumbnail"
              class="w-[200px] h-[112px] flex justify-center relative group"
            >
              <a-image
                :src="createDoc.thumbnail"
                class="rounded-md cursor-pointer object-cover"
                :preview="false"
                @click="thumbnailModalOpen = true"
              />
              <div
                class="absolute top-1 right-1 z-10 opacity-0 group-hover:opacity-100 transition-opacity"
              >
                <a-button
                  type="primary"
                  shape="circle"
                  danger
                  size="small"
                  @click.stop="createDoc.thumbnail = ''"
                >
                  <CloseOutlined />
                </a-button>
              </div>
            </div>
            <div
              v-else
              class="w-[200px] h-[112px] flex items-center justify-center border-1 border-dashed border-gray-300 rounded-md cursor-pointer text-zinc-400 dark:text-zinc-600"
              @click="thumbnailModalOpen = true"
            >
              <span class="text-sm">选择图片</span>
            </div>
          </div>
        </a-form-item>
        <a-form-item label="是否公开">
          <a-switch v-model:checked="createDoc.is_public" />
        </a-form-item>
      </a-form>
    </a-modal>
    <PhotoSelect
      v-model:open="thumbnailModalOpen"
      @selected-picture="handleSelectedPicture"
    />
  </div>
</template>

<script setup lang="ts">
import { createDocumentAPI, getAllRootDocAPI } from "@/api/document";
import type { DocumentRootRequest, DocumentRootVO } from "@/types/document";
import { message } from "ant-design-vue";
import { CloseOutlined } from "@ant-design/icons-vue";
import SvgIcon from "@/components/SvgIcon/index.vue";
import PhotoSelect from "@/components/PhotoSelect/index.vue";
import type { FormInstance } from "ant-design-vue";
const docList = ref<DocumentRootVO[]>([]);
const showModal = ref(false);
const thumbnailModalOpen = ref(false);
const loading = ref(false);
const formRef = ref<FormInstance>();
const rules: Record<string, any> = {
  title: [{ required: true, message: "请输入标题" }],
  description: [{ required: true, message: "请输入描述" }],
};
const createDoc = ref<DocumentRootRequest>({
  title: "",
  content: "",
  description: "",
  thumbnail: "",
  document_type: "root",
  is_public: false,
});

const getAllDoc = async () => {
  loading.value = true;
  const res = await getAllRootDocAPI();
  docList.value = res.data;
  loading.value = false;
};

const clearCreateDoc = () => {
  createDoc.value.title = "";
  createDoc.value.is_public = false;
  createDoc.value.description = "";
  createDoc.value.thumbnail = "";
};

const handleOk = async () => {
  await formRef.value?.validate();
  await createDocumentAPI(createDoc.value);

  message.success("文档创建成功");
  showModal.value = false;
  clearCreateDoc(); // 重新获取文档列表
  await getAllDoc();
};

const handleSelectedPicture = (picture: string) => {
  createDoc.value.thumbnail = picture;
};

const handleCancel = () => {
  showModal.value = false;
};

onMounted(getAllDoc);
</script>
