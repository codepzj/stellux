<template>
  <div class="h-full overflow-x-hidden">
    <a-card class="h-full" :bordered="false">
      <div class="h-full flex">
        <SplitPanel>
          <template #left-content>
            <a-spin class="w-full mx-auto" :spinning="loading">
              <a-tree
                v-if="treeData?.length"
                v-model:expandedKeys="expandedKeys"
                v-model:selectedKeys="selectedKeys"
                :tree-data="treeData"
                @select="onHandleSelect"
                block-node
                show-icon
              >
                <template #title="{ title, key }">
                  <a-dropdown :trigger="['contextmenu']">
                    <div class="cursor-pointer w-full m-1">{{ title }}</div>
                    <template #overlay>
                      <a-menu v-if="isRoot(key)">
                        <a-menu-item key="1" @click="handleRenameModalOpen(key)"
                          >重命名</a-menu-item
                        >
                        <a-menu-item
                          key="2"
                          @click="openCreateModal('parent', key, document_id)"
                          >新增目录</a-menu-item
                        >
                        <a-menu-item
                          key="3"
                          @click="openCreateModal('leaf', key, document_id)"
                          >新增文档</a-menu-item
                        >
                      </a-menu>
                      <a-menu v-else-if="isParent(key)">
                        <a-menu-item key="1" @click="handleRenameModalOpen(key)"
                          >重命名</a-menu-item
                        >
                        <a-menu-item
                          key="2"
                          @click="openCreateModal('parent', key, document_id)"
                          >新增目录</a-menu-item
                        >
                        <a-menu-item
                          key="3"
                          @click="openCreateModal('leaf', key, document_id)"
                          >新增文档</a-menu-item
                        >
                        <a-menu-item key="4" @click="deleteDocumentParent(key)"
                          >删除目录</a-menu-item
                        >
                      </a-menu>
                      <a-menu v-else>
                        <a-menu-item key="1" @click="handleRenameModalOpen(key)"
                          >重命名</a-menu-item
                        >
                        <a-menu-item key="2" @click="deleteDocumentLeaf(key)"
                          >删除文档</a-menu-item
                        >
                      </a-menu>
                    </template>
                  </a-dropdown>
                </template>

                <template #switcherIcon="{ dataRef, expanded }">
                  <SvgIcon
                    v-if="isRoot(dataRef.key) || isParent(dataRef.key)"
                    :size="expanded ? 18 : 16"
                    :name="expanded ? 'folder-open' : 'folder'"
                    class="cursor-pointer"
                  />
                </template>

                <template #icon="{ key }">
                  <template v-if="isLeaf(key)">
                    <SvgIcon :size="16" name="md" />
                  </template>
                  <template
                    v-if="
                      (isParent(key) || isRoot(key)) &&
                      getChildrenLengthByKey(treeData, key) === 0
                    "
                  >
                    <SvgIcon :size="16" name="folder" />
                  </template>
                </template>
              </a-tree>
            </a-spin>
          </template>

          <template #right-content>
            <a-skeleton active v-if="loading" />
            <div v-else>
              <div class="pr-2.5">
                <a-alert
                  v-if="isRoot(selectedKeys[0])"
                  class="!py-3 !my-3"
                  message="该页面用于展示文档首页"
                  type="info"
                  closable
                  show-icon
                />
                <DocumentDetail class="mt-3" v-model:content="documentDetail.content" />
              </div>
            </div>
          </template>
        </SplitPanel>
      </div>
    </a-card>

    <a-modal
      v-model:open="modal.visible"
      :title="modal.title"
      @ok="handleCreate"
      ok-text="创建"
      cancel-text="取消"
    >
      <div class="flex flex-col gap-2 py-4">
        <a-input
          v-model:value="modal.input"
          placeholder="请输入标题"
          @keyup.enter="handleCreate"
        />
      </div>
    </a-modal>

    <a-modal
      v-model:open="renameModal.visible"
      title="重命名"
      @ok="handleRename"
      ok-text="确认"
      cancel-text="取消"
    >
      <div class="flex flex-col gap-2 py-4">
        <a-input
          v-model:value="renameModal.title"
          placeholder="请输入标题"
          @keyup.enter="handleRename"
        />
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, nextTick, computed, watchEffect } from "vue";
import { useRoute } from "vue-router";
import { message, type TreeProps } from "ant-design-vue";

import SvgIcon from "@/components/SvgIcon/index.vue";
import SplitPanel from "@/components/SplitPanel/index.vue";
import DocumentDetail from "./components/DocumentDetail.vue";

import {
  createDocumentAPI,
  deleteDocumentByIDListAPI,
  deleteDocumentLeafByIDAPI,
  getAllParentDocAPI,
  getDocumentByIDAPI,
  getDocumentTreeByIDAPI,
  renameDocumentAPI,
  saveDocumentAPI,
} from "@/api/document";

import { useDocumentStore } from "@/store/modules/document";
import {
  convertToTreeData,
  getChildrenLengthByKey,
  getAllChildIdByParentId,
} from "@/utils/convert";

import type { DocumentRequest, DocumentTreeVO } from "@/types/document";

const route = useRoute();
const document_id = ref(route.params.id as string);

const expandedKeys = ref<string[]>([document_id.value]);
const selectedKeys = ref<string[]>([]);
const title = ref("");
const loading = ref(true);
const docTreeData = ref<DocumentTreeVO[]>([]);
const treeData = ref<TreeProps["treeData"]>([]);
const documentStore = useDocumentStore();

const documentDetail = ref({
  id: document_id.value,
  title: "",
  content: "",
});

const originDocument = ref({
  id: document_id.value,
  title: "",
  content: "",
});

const modal = reactive({
  visible: false,
  title: "",
  type: "", // 'parent' or 'leaf'
  parent_id: "",
  document_id: "",
  input: "",
  expandedKey: "",
});

const renameModal = reactive({
  visible: false,
  id: "",
  title: "",
});

const openCreateModal = (
  type: string,
  parent_id: string,
  document_id: string
) => {
  modal.visible = true;
  modal.type = type;
  modal.parent_id = parent_id;
  modal.document_id = document_id;
  modal.input = "";
  modal.title = type === "parent" ? "新增目录" : "新增文档";
  modal.expandedKey = parent_id;
};

const handleCreate = async () => {
  if (!modal.input.trim()) {
    message.warning("请输入标题");
    return;
  }

  const newDocument: DocumentRequest = {
    title: modal.input,
    content: "",
    is_public: true,
    parent_id: modal.parent_id,
    document_id: modal.document_id,
    document_type: modal.type,
  };

  await createDocumentAPI(newDocument);
  expandedKeys.value = [...expandedKeys.value, modal.expandedKey];
  await getDocumentTree(route.params.id as string);
  modal.visible = false;
  message.success("创建成功");
};

const getDocumentTree = async (id: string) => {
  loading.value = true;
  const res = await getDocumentTreeByIDAPI(id);
  docTreeData.value = res.data;
  treeData.value = convertToTreeData(docTreeData.value);
  loading.value = false;
};

const getParentDoc = async () => {
  const res = await getAllParentDocAPI(document_id.value);
  expandedKeys.value = [
    ...expandedKeys.value,
    ...res.data.map(item => item.id),
  ];
};

const deleteDocumentLeaf = async (id: string) => {
  await deleteDocumentLeafByIDAPI(id);
  await getDocumentTree(route.params.id as string);
};

const deleteDocumentParent = async (id: string) => {
  const childIds = getAllChildIdByParentId(treeData.value, id);
  await deleteDocumentByIDListAPI(childIds);
  await getDocumentTree(route.params.id as string);
};

const handleRenameModalOpen = (id: string) => {
  renameModal.visible = true;
  renameModal.id = id;
  renameModal.title = docTreeData.value.find(item => item.id === id)
    ?.title as string;
};

const handleRename = async () => {
  await renameDocumentAPI(renameModal.id, renameModal.title);
  await getDocumentTree(route.params.id as string);
  renameModal.visible = false;
  message.success("重命名成功");
};

const saveDocument = async () => {
  await saveDocumentAPI({
    id: documentDetail.value.id,
    title: documentDetail.value.title,
    content: documentDetail.value.content,
  });
  originDocument.value.title = documentDetail.value.title;
  originDocument.value.content = documentDetail.value.content;
  message.success("保存成功");
};

const onHandleSelect = async (keys: string[]) => {
  const key = keys[0];
  if (isLeaf(key) || isRoot(key)) {
    const res = await getDocumentByIDAPI(key);
    documentDetail.value = {
      id: res.data.id,
      title: res.data.title,
      content: res.data.content,
    };
    originDocument.value = {
      id: res.data.id,
      title: res.data.title,
      content: res.data.content,
    };
  }
};

const isRoot = (id: string): boolean =>
  docTreeData.value.find(item => item.id === id)?.document_type === "root";

const isParent = (id: string): boolean =>
  docTreeData.value.find(item => item.id === id)?.document_type === "parent";

const isLeaf = (id: string): boolean =>
  docTreeData.value.find(item => item.id === id)?.document_type === "leaf";

const isDocumentAlreadyEdit = computed(() => {
  return (
    documentDetail.value.title !== originDocument.value.title ||
    documentDetail.value.content !== originDocument.value.content
  );
});

watchEffect(() => {
  documentStore.setActionButtonList([
    {
      name: "保存",
      disabled: !isDocumentAlreadyEdit.value,
      action: saveDocument,
    },
  ]);
  documentStore.setEditDocumentTitle(documentDetail.value.title);
});

// 页面加载逻辑（先加载文档树，然后加载首页）
onMounted(async () => {
  await getDocumentTree(document_id.value);
  await nextTick();

  const rootNode = docTreeData.value.find(
    item => item.document_type === "root"
  );
  if (rootNode) {
    selectedKeys.value = [rootNode.id];
    const res = await getDocumentByIDAPI(rootNode.id);
    documentDetail.value = {
      id: res.data.id,
      title: res.data.title,
      content: res.data.content,
    };
    originDocument.value = {
      id: res.data.id,
      title: res.data.title,
      content: res.data.content,
    };
    title.value = res.data.title;
  }

  await getParentDoc();
});
</script>

<style lang="scss" scoped>
:deep(.ant-card-body) {
  height: 100%;
  padding: 10px 0 10px 10px;
}

:deep(.ant-tree .ant-tree-switcher) {
  width: 18px;
  display: flex;
  align-items: center;
}
:deep(.ant-tree .ant-tree-switcher-noop) {
  width: 0;
}
:deep(.ant-tree .ant-tree-node-content-wrapper) {
  display: flex;
  align-items: center;
  line-height: 22px;
  user-select: none;
  padding: 0;
}
:deep(.ant-tree .ant-tree-title) {
  width: 100%;
}
:deep(.ant-tree .ant-tree-indent-unit) {
  display: inline-block;
  width: 18px;
}
</style>
