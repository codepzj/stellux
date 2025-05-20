import { PageVO, Page } from "@/types/page";
import { PostVO, PostSearchVO, SiteMapVO } from "@/types/post";
import request from "@/utils/request";

// 获取首页文章列表
export const getPostDetailListAPI = (page: Page) => {
  return request.get<PageVO<PostVO>>("/post/detail/list", page);
};

// 获取文章详情
export const getPostDetailAPI = (id: string) => {
  return request.get<PostVO>(`/post/detail/${id}`);
};

// 搜索框获取文章列表
export const getPostByKeyWordAPI = (keyword: string) => {
  return request.get<PostSearchVO[]>(`/post/search?keyword=${keyword}`);
};

// 获取站点地图
export const getSiteMapAPI = () => {
  return request.get<SiteMapVO[]>(`/post/sitemap`);
};
