<template>
  <div class="flex items-center">
    <span class="mr-4">
      <EnterOutlined class="icon font-[15px] p-0.5 mr-1" />
      确认
    </span>
    <span class="mr-4">
      <ArrowUpOutlined class="icon font-[15px] p-0.5 mr-1" />
      <ArrowDownOutlined class="icon font-[15px] p-0.5 mr-1" />
      切换
    </span>
    <span>
      <CloseOutlined class="icon font-[15px] p-0.5 mr-1" />
      关闭
    </span>
  </div>
</template>

<script lang="ts" setup>
import {
  EnterOutlined,
  ArrowDownOutlined,
  ArrowUpOutlined,
  CloseOutlined,
} from "@ant-design/icons-vue";
</script>
<style lang="scss" scoped>
.icon {
  box-shadow:
    inset 0 -2px #cdcde6,
    inset 0 0 1px 1px #fff,
    0 1px 2px 1px #1e235a66;
}
</style>
