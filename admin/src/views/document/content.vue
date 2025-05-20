<template>
  <div class="h-full">
    <a-card class="h-full overflow-scroll" :bordered="false">
      <div class="h-full flex">
        <SplitPanel>
          <template #left-content>
            <a-spin class="w-full mx-auto" :spinning="loading">
              <a-tree
                v-if="treeData?.length"
                v-model:selectedKeys="selectedKeys"
                :tree-data="treeData"
                @select="onHandleSelect"
                block-node
                show-icon
                default-expand-all
              >
                <template #title="{ title, key }">
                  <a-dropdown :trigger="['contextmenu']">
                    <div class="cursor-pointer w-full p-1">{{ title }}</div>
                    <template #overlay>
                      <a-menu v-if="isRoot(key)">
                        <a-menu-item key="1">重命名</a-menu-item>
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
                        <a-menu-item key="1">重命名</a-menu-item>
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
                        <a-menu-item key="4">删除目录</a-menu-item>
                      </a-menu>
                      <a-menu v-else>
                        <a-menu-item key="1">重命名</a-menu-item>
                        <a-menu-item key="2">删除文档</a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
                </template>
                <template #switcherIcon="{ dataRef }">
                  <SvgIcon
                    v-if="isRoot(dataRef.key) || isParent(dataRef.key)"
                    :size="16"
                    name="folder"
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
            <DocumentDetail v-model:title="title" />
          </template>
        </SplitPanel>
      </div>
    </a-card>

    <!-- 创建目录或文档的弹窗 -->
    <a-modal
      v-model:open="modal.visible"
      :title="modal.title"
      @ok="handleCreate"
      ok-text="创建"
      cancel-text="取消"
    >
      <a-input
        v-model:value="modal.input"
        placeholder="请输入标题"
        @keyup.enter="handleCreate"
      />
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import SvgIcon from "@/components/SvgIcon/index.vue";
import SplitPanel from "@/components/SplitPanel/index.vue";
import { message, type TreeProps } from "ant-design-vue";
import { createDocumentAPI, getDocumentTreeByIDAPI } from "@/api/document";
import { convertToTreeData, getChildrenLengthByKey } from "@/utils/convert";
import type { DocumentRequest, DocumentTreeVO } from "@/types/document";
import DocumentDetail from "./components/DocumentDetail.vue";

const route = useRoute();

const selectedKeys = ref<string[]>([]);
const title = ref("");
const document_id = ref(route.params.id as string);

const loading = ref(true);
const docTreeData = ref<DocumentTreeVO[]>([]);
const treeData = ref<TreeProps["treeData"]>([]);

// 弹窗状态管理
const modal = reactive({
  visible: false,
  title: "",
  type: "", // 'parent' or 'leaf'
  parent_id: "",
  document_id: "",
  input: "",
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
  await getDocumentTree(route.params.id as string);
  modal.visible = false;
  message.success("创建成功");
};

const getDocumentTree = async (id: string) => {
  loading.value = true;
  const res = await getDocumentTreeByIDAPI(id);
  docTreeData.value = res.data;
  treeData.value = convertToTreeData(docTreeData.value);
  console.log("treeData", treeData.value);
  loading.value = false;
};

onMounted(async () => {
  await getDocumentTree(route.params.id as string);
  const rootNode = docTreeData.value.find(
    item => item.document_type === "root"
  );
  if (rootNode) {
    selectedKeys.value = [rootNode.id];
    title.value = rootNode.title;
  }
});

const isRoot = (id: string): boolean => {
  return (
    docTreeData.value.find(item => item.id === id)?.document_type === "root"
  );
};

const isParent = (id: string): boolean => {
  return (
    docTreeData.value.find(item => item.id === id)?.document_type === "parent"
  );
};

const isLeaf = (id: string): boolean => {
  return (
    docTreeData.value.find(item => item.id === id)?.document_type === "leaf"
  );
};

const onHandleSelect = (keys: string[], event: any) => {
  console.log("keys", keys);
  if (isLeaf(keys[0]) || isRoot(keys[0])) {
    title.value = event.node.title;
  }
};
</script>

<style lang="scss" scoped>
:deep(.ant-card-body) {
  height: 100%;
  padding: 10px;
}
:deep(.ant-tree .ant-tree-switcher) {
  width: 16px;
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
:deep(.ant-tree .ant-tree-switcher) {
  display: flex;
  align-items: center;
}
:deep(.ant-tree .ant-tree-title) {
  width: 100%;
}
:deep(.ant-tree .ant-tree-indent-unit) {
  display: inline-block;
  width: 18px;
}
</style>
