<template>
  <div class="overflow-hidden py-8">
    <a-tabs
      v-model:activeKey="activeKey"
      tab-position="left"
      @change="(key: string) => $router.push({ name: key })"
    >
      <a-tab-pane v-for="item in tabMeta" :key="item.key" :tab="item.title">
        <router-view></router-view>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script lang="ts" setup>
const activeKey = ref("UserEditBasic");
const tabMeta = [
  {
    key: "UserEditBasic",
    title: "基本信息",
  },
  {
    key: "UserEditResetPassword",
    title: "重置密码",
  },
];
</script>
