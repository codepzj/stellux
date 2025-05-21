export type DocumentRootRequest = {
  title: string;
  content: string;
  description: string;
  thumbnail: string;
  document_type: string;
  is_public: boolean;
};

export type DocumentRequest = {
  title: string;
  content: string;
  document_type: string;
  is_public: boolean;
  parent_id: string;
  document_id: string;
};

export type DocumentSaveRequest = {
  id: string;
  title: string;
  content: string;
};

export type DocumentVO = {
  id: string;
  created_at: string;
  updated_at: string;
  title: string;
  content: string;
  document_type: string;
  is_public: boolean;
  parent_id: string;
  document_id: string;
};

export type DocumentRootVO = {
  id: string;
  created_at: string;
  updated_at: string;
  title: string;
  description: string;
  thumbnail: string;
  is_public: boolean;
};

export type DocumentTreeVO = {
  id: string;
  created_at: string;
  updated_at: string;
  title: string;
  document_type: string;
  is_public: boolean;
  parent_id: string;
  document_id: string;
};
