<template>
  <div class="p-6">
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <a-card
        v-for="doc in docList"
        :key="doc.id"
        class="cursor-pointer hover:translate-y-[-5px] transition-all duration-300"
        @click="$router.push({ name: 'DocumentContent', params: { id: doc.id } })"
      >
        <h3>{{ doc.title }}</h3>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getAllRootDocAPI } from '@/api/document';
import type { DocumentTreeVO } from '@/types/document';

const docList = ref<DocumentTreeVO[]>([]);

const getAllDoc = async () => {
  const res = await getAllRootDocAPI();
  docList.value = res.data;
};

onMounted(getAllDoc);
</script>
