<template>
  <span class="cursor-pointer" @click="toggle">
    <SvgIcon name="fullscreen" class="rounded-md p-1" :size="28" />
  </span>
</template>
<script lang="ts" setup>
import { useFullscreen } from "@vueuse/core";

import SvgIcon from "@/components/SvgIcon/index.vue";
const { toggle } = useFullscreen();
</script>
