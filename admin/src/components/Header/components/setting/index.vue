<template>
  <div>
    <Tooltip title="设置" placement="bottom" :mouse-enter-delay="0.5">
      <SvgIcon
        name="setting"
        :size="28"
        class="cursor-pointer rounded-md p-1"
        @click="showDrawer"
      />
    </Tooltip>
    <Drawer v-model:open="visible" placement="right" :closable="false">
      <Descriptions title="主题风格" :column="5">
        <Descriptions.Item v-for="theme in themeStyle" :key="theme.value">
          <Tooltip :title="theme.label">
            <div
              class="style-checbox-item"
              :class="{ active: themeMode === theme.value }"
              @click="systemStore.setThemeMode(theme.value as 'light' | 'dark')"
            >
              <img :src="theme.value === 'light' ? LightTheme : DarkTheme" />
            </div>
          </Tooltip>
        </Descriptions.Item>
      </Descriptions>
      <Descriptions title="主题颜色" :column="9">
        <Descriptions.Item v-for="item in themeColors" :key="item.key">
          <div class="style-checbox-item">
            <Tooltip :title="item.title">
              <Tag
                :color="item.value"
                @click="systemStore.setThemeColor(item.value)"
              >
                <span :style="{ visibility: getThemeColorVisible(item.value) }">
                  ✔
                </span>
              </Tag>
            </Tooltip>
          </div>
        </Descriptions.Item>
      </Descriptions>
    </Drawer>
  </div>
</template>

<script lang="ts" setup>
import SvgIcon from "@/components/SvgIcon/index.vue";
import { storeToRefs } from "pinia";
import { Drawer, Descriptions, Tag, Tooltip } from "ant-design-vue";
import { themeColors, themeStyle } from "./constant";
import { useSystemStore } from "@/store";
import LightTheme from "@/assets/svg/light.svg";
import DarkTheme from "@/assets/svg/dark.svg";

defineOptions({
  name: "ProjectSetting",
});

const systemStore = useSystemStore();
const { themeMode, themeColor } = storeToRefs(systemStore);
const visible = ref(false);

const getThemeColorVisible = (color: string) =>
  themeColor.value === color ? "visible" : "hidden";

const showDrawer = () => {
  visible.value = true;
};
</script>

<style lang="scss" scoped>
.style-checbox-item {
  position: relative;
  cursor: pointer;

  &.active::after {
    content: "✔";
    position: absolute;
    right: 12px;
    bottom: 10px;
    color: var(--app-primary-color);
  }
}

input[type="color"] {
  width: 40px;
  height: 40px;
  padding: 0;
  border: 0;
  outline: none;
  appearance: none;

  &::-webkit-color-swatch-wrapper {
    background: var(--custom-color);
  }

  &::-webkit-color-swatch {
    display: none;
  }
}
</style>
