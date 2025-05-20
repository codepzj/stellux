import request from "@/utils/request";
import type { Response } from "@/types/dto";
import type {
  DocumentRootRequest,
  DocumentRequest,
  DocumentTreeVO,
  DocumentRootVO,
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
