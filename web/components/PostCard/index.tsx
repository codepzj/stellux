"use client";

import { IPostCard } from "@/types/post";
import { Card, CardBody, CardHeader, CardFooter } from "@heroui/card";
import { useRouter } from "next/navigation";
import { cn } from "@/lib/utils";
import { Bookmark, Tag, Calendar } from "lucide-react";
import dayjs from "dayjs";
import { Image } from "@heroui/image";
import { TopIcon } from "@/components/SvgIcon";

export function PostCard({
  post,
  className,
}: {
  post: IPostCard;
  className?: string;
}) {
  const router = useRouter();
  const tags = post.tags?.join(", ");

  return (
    <Card
      key={post.id}
      className={cn(
        "p-6 shadow transition-all hover:shadow-md overflow-hidden max-w-full",
        className
      )}
      onPress={() => router.push(`/post/${post.id}`)}
      isHoverable
      isPressable
    >
      <div className="flex flex-wrap md:flex-nowrap gap-4">
        {/* Left side: Title, Description, Footer */}
        <div className="flex-1 min-w-0">
          <CardHeader className="p-0 mb-4 flex items-center gap-2">
            <span
              className="text-base md:text-lg font-bold truncate px-2 max-w-full"
              title={post.title}
            >
              {post.title}
            </span>
            {post.is_top && (
              <TopIcon className="text-red-500 flex-shrink-0" size={24} />
            )}
          </CardHeader>

          <CardBody className="p-0 mb-4 md:min-h-8">
            <p
              className="text-sm text-muted-foreground leading-relaxed px-2 max-w-full line-clamp-2 overflow-hidden text-ellipsis"
              title={post.description}
            >
              {post.description}
            </p>
          </CardBody>

          <CardFooter className="hidden md:flex flex-wrap gap-2 p-0 text-sm text-muted-foreground">
            {post.category && (
              <span
                className="flex items-center gap-1.5 px-2 py-0.5 truncate"
                title={`分类：${post.category}`}
              >
                <Bookmark size={14} />
                {post.category}
              </span>
            )}
            {tags && (
              <span
                className="flex items-center gap-1.5 px-2 py-0.5 truncate"
                title={`标签：${tags}`}
              >
                <Tag size={14} />
                {tags}
              </span>
            )}
            <span
              className="flex items-center gap-1.5 px-2 py-0.5"
              title={post.created_at}
            >
              <Calendar size={14} />
              {dayjs(post.created_at).format("YYYY-MM-DD")}
            </span>
          </CardFooter>
        </div>

        {/* Right side: Image */}
        <div className="hidden md:flex justify-center items-center w-48 h-auto flex-shrink-0">
          {post.thumbnail && (
            <Image
              src={post.thumbnail}
              alt={post.title}
              width={180}
              height={120}
              shadow="none"
              loading="lazy"
              isZoomed
              className="object-cover"
            />
          )}
        </div>
      </div>
    </Card>
  );
}