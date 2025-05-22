import { getPostDetailAPI } from "@/api/post";
import Md from "@/components/Md";
import { Toc } from "@/components/Toc";
import { getTableOfContents } from "@/lib/toc";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Metadata } from "next";

type Props = {
  params: Promise<{ id: string }>
}

export default async function PostPage({ params }: Props) {
  const { id } = await params;
  const post = await getPostDetailAPI(id);
  const toc = await getTableOfContents(post.data.content);
  return (
    <div className="relative max-w-6xl mx-auto px-4">
      <div className="w-full flex flex-col justify-end md:flex-row">
        {/* Markdown 内容 */}
        <div className="w-full md:w-3/4 md:pr-12">
          <h1 className="text-4xl font-bold mb-10">{post.data.title}</h1>
          <Md content={post.data.content} />
        </div>
        {/* 目录 */}
        <div className="hidden md:block md:w-1/5">
          <div className="sticky top-32">
            <ScrollShadow className="w-full max-h-[400px] overflow-y-auto" size={10} >
              <Toc toc={toc} />
            </ScrollShadow>
          </div>
        </div>
      </div>
    </div>
  );
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { id } = await params;
  const post = await getPostDetailAPI(id);
  const data = post.data;

  const title = data.title || "默认标题";
  const description = data.description || "默认描述";
  const image = data.thumbnail || "/default-thumbnail.jpg";
  const keywords = [data.category, ...(data.tags || [])].filter(Boolean);

  const url = `${process.env.NEXT_PUBLIC_SITE_URL}/post/${id}`;

  return {
    title,
    description,
    keywords,
    openGraph: {
      title,
      description,
      url,
      type: "article",
      images: [{ url: image }],
    },
    twitter: {
      card: "summary_large_image",
      title,
      description,
      images: [image],
    },
    authors: [{ name: data.author }],
    metadataBase: new URL(url),
  };
}
