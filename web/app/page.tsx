import { getPostDetailListAPI } from "@/api/post";
import { PostList } from "@/components/PostList";
import type { PageVO, Page } from "@/types/page";
import type { PostVO } from "@/types/post";

// SSR 获取首屏数据
export default async function Home() {
  const page: Page = {
    page_no: 1,
    page_size: 10,
  };

  const res = await getPostDetailListAPI(page);
  const listData: PageVO<PostVO> = res?.data || {
    list: [],
    page_no: 1,
    size: 10,
    total_count: 0,
    total_page: 1,
    field: "",
    order: "DESC",
  };

  return (
    <div className="flex flex-col gap-4 max-w-2xl mx-auto">
      <PostList initialData={listData} />
    </div>
  );
}
