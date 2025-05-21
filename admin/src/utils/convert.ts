import type { TreeProps } from "ant-design-vue";
import type { DocumentTreeVO } from "@/types/document";

export function convertToTreeData(
  data: DocumentTreeVO[]
): TreeProps["treeData"] {
  const map = new Map<string, any>();
  const roots: any[] = [];

  data.forEach(item => {
    map.set(item.id, {
      key: item.id,
      title: item.title,
      children: [],
    });
  });

  data.forEach(item => {
    const node = map.get(item.id);
    if (item.parent_id && map.has(item.parent_id)) {
      map.get(item.parent_id).children.push(node);
    } else {
      roots.push(node);
    }
  });

  return roots;
}

/**
 * 在 treeData 中递归查找指定 key 的节点
 */
export function findNodeByKey(
  tree: TreeProps["treeData"],
  key: string
): any | null {
  if (!tree) return null;
  for (const node of tree) {
    if (node.key === key) {
      return node;
    }
    if (node.children?.length) {
      const found = findNodeByKey(node.children, key);
      if (found) return found;
    }
  }
  return null;
}

/**
 * 获取指定 key 节点的 children 长度
 */
export function getChildrenLengthByKey(
  tree: TreeProps["treeData"],
  key: string
): number {
  const node = findNodeByKey(tree, key);
  return node?.children?.length || 0;
}

/**
 * 获取指定 key 节点及其所有子孙节点的 key 数组(删除子目录时使用)
 */
export function getAllChildIdByParentId(
  tree: TreeProps["treeData"],
  parentKey: string
): string[] {
  const result: string[] = [];

  const node = findNodeByKey(tree, parentKey);
  if (!node) return result;

  function dfs(currentNode: any) {
    result.push(currentNode.key);
    if (currentNode.children?.length) {
      currentNode.children.forEach((child: any) => dfs(child));
    }
  }

  dfs(node);
  return result;
}
