import request from "@/utils/request";
import type { Response } from "@/types/dto";
import type {
  DocumentRootRequest,
  DocumentRequest,
  DocumentTreeVO,
  DocumentRootVO,
  DocumentVO,
  DocumentSaveRequest,
} from "@/types/document";

export const createDocumentAPI = (
  data: DocumentRootRequest | DocumentRequest
): Promise<Response<any>> => {
  return request.post("/admin-api/document/create", data);
};

export const getDocumentTreeByIDAPI = (
  id: string
): Promise<Response<DocumentTreeVO[]>> => {
  return request.get(`/document/tree?document_id=${id}`);
};

export const getAllRootDocAPI = (): Promise<Response<DocumentRootVO[]>> => {
  return request.get("/admin-api/document/list");
};

export const getAllParentDocAPI = (document_id: string): Promise<Response<DocumentRootVO[]>> => {
  return request.get(`/admin-api/document/parent-list?document_id=${document_id}`);
};

export const getDocumentDetailByIDAPI = (id: string): Promise<Response<string>> => {
  return request.get(`/admin-api/document/detail?document_id=${id}`);
};

// 删除叶子节点
export const deleteDocumentLeafByIDAPI = (document_id: string): Promise<Response<any>> => {
  return request.delete(`/admin-api/document/delete`, {
    data: {
      document_id,
    },
  });
};

// 删除父节点
export const deleteDocumentByIDListAPI = (document_id_list: string[]): Promise<Response<any>> => {
  return request.delete(`/admin-api/document/delete-list`, {
    data: {
      document_id_list,
    },
  });
};

// 获取特定文档详情
export const getDocumentByIDAPI = (id: string): Promise<Response<DocumentVO>> => {
  return request.get(`/admin-api/document/${id}`);
};

// 保存文档
export const saveDocumentAPI = (data: DocumentSaveRequest): Promise<Response<any>> => {
  return request.put("/admin-api/document/save", data);
};

// 重命名文档
export const renameDocumentAPI = (id: string, title: string): Promise<Response<any>> => {
  return request.put("/admin-api/document/rename", { id, title });
};
