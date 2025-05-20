import { Page } from "./page";

export interface PostVO {
  id: string;
  created_at: string;
  title: string;
  content: string;
  description: string;
  author: string;
  category: string;
  tags: string[];
  thumbnail: string;
  is_top: boolean;
  is_publish: boolean;
}

export interface PostSearchVO {
  id: string;
  title: string;
  description: string;
  created_at: string;
  author: string;
  thumbnail: string;
  category_id: string;
  tags_id: string[];
  is_publish: boolean;
  is_top: boolean;
}

export interface IPostCard {
  id: string;
  title: string;
  content: string;
  description: string;
  created_at: string;
  author: string;
  thumbnail: string;
  category: string;
  tags: string[];
  is_publish: boolean;
  is_top: boolean;
}

export interface PostPageReq extends Page {
  post_type: string;
}

export interface SiteMapVO {
  id: string;
  updated_at: string;
}
