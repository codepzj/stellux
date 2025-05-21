import type { ActionButton } from "@/types/action-button";

export const useDocumentStore = defineStore(
  "document",
  () => {
    const mode = ref<"preview" | "edit">("preview");

    const setMode = (newMode: "preview" | "edit") => {
      mode.value = newMode;
    };

    const editDocumentTitle = ref("");

    const setEditDocumentTitle = (newEditDocumentTitle: string) => {
      editDocumentTitle.value = newEditDocumentTitle;
    };

    const actionButtonList = ref<ActionButton[]>([]);

    const setActionButtonList = (newActionButtonList: ActionButton[]) => {
      actionButtonList.value = newActionButtonList;
    };

    return {
      mode,
      setMode,
      editDocumentTitle,
      setEditDocumentTitle,
      actionButtonList,
      setActionButtonList,
    };
  },
  {
    persist: true,
  }
);
