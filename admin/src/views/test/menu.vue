<template>
  <div>
    <a-menu :items="items"> </a-menu>
  </div>
</template>

<script setup>
const items = [
  { label: "个人中心", key: "1" },
  { label: "设置", key: "2" },
  { label: "退出", key: "3" },
];
</script>
