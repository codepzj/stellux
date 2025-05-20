export type DocumentRootRequest = {
  title: string;
  content: string;
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