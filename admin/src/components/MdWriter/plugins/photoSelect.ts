import type { BytemdPlugin } from "bytemd";
import type { Ref } from "vue";

/**
 * @param modalVisibleRef 外部定义的 ref，用于控制弹窗显示
 */
export default function imageUploadPlugin(
  modalVisibleRef: Ref<boolean>
): BytemdPlugin {
  return {
    actions: [
      {
        icon: '<svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" fill="none" viewBox="0 0 48 48"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="4" d="M5 10a2 2 0 0 1 2-2h34a2 2 0 0 1 2 2v28a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V10Z" clip-rule="evenodd"></path><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="4" d="M14.5 18a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3Z" clip-rule="evenodd"></path><path stroke="currentColor" stroke-linejoin="round" stroke-width="4" d="m15 24 5 4 6-7 17 13v4a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2v-4l10-10Z"></path></svg>',
        title: "选择图片",
        handler: {
          type: "action",
          click() {
            modalVisibleRef.value = true;
          },
        },
      },
    ],
  };
}
