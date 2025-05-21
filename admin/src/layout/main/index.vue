<template>
  <div class="h-full">
    <router-view v-slot="{ Component }">
      <template v-if="Component">
        <Suspense>
          <Transition name="fade-slide" mode="out-in" appear>
            <keep-alive :include="cacheNames">
              <component :is="Component" />
            </keep-alive>
          </Transition>
          <template #fallback>
            <a-skeleton />
          </template>
        </Suspense>
      </template>
    </router-view>
  </div>
</template>

<script lang="ts" setup>
import { cacheNames } from "@/router";
</script>
<style lang="scss" scoped>
.fade-slide-leave-active,
.fade-slide-enter-active {
  transition: all 0.3s;
}

.fade-slide-enter-from {
  transform: translateX(-30px);
  opacity: 0;
}

.fade-slide-leave-to {
  transform: translateX(30px);
  opacity: 0;
}
</style>
