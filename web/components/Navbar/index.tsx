"use client";

import {
  Navbar,
  NavbarContent,
  NavbarBrand,
  NavbarItem,
} from "@heroui/navbar";
import { Link } from "@heroui/link";
import NextLink from "next/link";
import { Input } from "@heroui/input";
import { Modal, ModalContent } from "@heroui/modal";

import { siteConfig } from "@/config/site";
import { ThemeSwitch } from "@/components/ThemeSwitcher";
import { GithubIcon, SearchLinearIcon } from "@/components/SvgIcon";
import { Logo } from "@/components/Logo";

import { useState, useEffect } from "react";
import { PostSearchVO } from "@/types/post";
import { getPostByKeyWordAPI } from "@/api/post";

export default () => {
  const [isOpen, setIsOpen] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [postList, setPostList] = useState<PostSearchVO[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  // 搜索文章（防抖）
  useEffect(() => {
    const delayDebounce = setTimeout(async () => {
      if (!keyword.trim()) {
        setPostList([]);
        setIsLoading(false);
        return;
      }
      setIsLoading(true);
      try {
        const res = await getPostByKeyWordAPI(keyword);
        setPostList(res.data || []);
      } finally {
        setIsLoading(false);
      }
    }, 300);

    return () => clearTimeout(delayDebounce);
  }, [keyword]);

  const handleClose = (open: boolean) => {
    setIsOpen(open);
    if (!open) {
      setKeyword("");
      setPostList([]);
      setIsLoading(false);
    }
  };

  const handleItemClick = () => {
    handleClose(false);
  };

  // 高亮关键词
  const highlightKeyword = (text: string, keyword: string) => {
    if (!keyword) return text;
    const parts = text.split(new RegExp(`(${keyword})`, "gi"));
    return parts.map((part, index) =>
      part.toLowerCase() === keyword.toLowerCase() ? (
        <mark key={index} className="bg-yellow-300 text-black px-1 rounded">
          {part}
        </mark>
      ) : (
        <span key={index}>{part}</span>
      )
    );
  };

  return (
    <>
      {/* 顶部导航栏 */}
      <Navbar maxWidth="lg">
        <NavbarContent className="basis-1/5 sm:basis-full" justify="start">
          <NavbarBrand as="li" className="gap-3 max-w-fit">
            <NextLink
              className="flex justify-start items-center gap-1"
              href="/"
            >
              <Logo />
            </NextLink>
          </NavbarBrand>
        </NavbarContent>

        <NavbarContent
          className="hidden h-11 rounded-full border-small border-default-200/20 bg-background/60 px-4 shadow-xs backdrop-blur-md backdrop-saturate-150 dark:bg-default-100/50 lg:flex"
          justify="center"
        >
          <NavbarItem isActive>
            <Link className="text-default-500" href="/" size="sm">
              首页
            </Link>
          </NavbarItem>
          <NavbarItem>
            <Link className="text-default-500" href="/" size="sm">
              文档
            </Link>
          </NavbarItem>
        </NavbarContent>

        <NavbarContent className="flex basis-1/5 sm:basis-full" justify="end">
          <NavbarItem className="flex gap-2 items-center">
            {/* 大屏显示搜索框 */}
            <div
              onClick={() => setIsOpen(true)}
              className="cursor-pointer w-[10rem] hidden sm:block"
            >
              <Input
                readOnly
                placeholder="搜索"
                size="sm"
                startContent={<SearchLinearIcon size={18} />}
                classNames={{
                  base: "h-10",
                  inputWrapper:
                    "h-full font-normal text-default-500 bg-default-400/20 dark:bg-default-500/20",
                  input: "text-small",
                }}
              />
            </div>

            {/* 小屏显示图标 */}
            <div
              onClick={() => setIsOpen(true)}
              className="sm:hidden cursor-pointer"
            >
              <SearchLinearIcon size={20} className="text-default-500" />
            </div>

            <Link
              isExternal
              aria-label="Github"
              href={siteConfig.links.github}
            >
              <GithubIcon className="text-default-500" />
            </Link>
            <ThemeSwitch />
          </NavbarItem>
        </NavbarContent>
      </Navbar>

      {/* 搜索弹窗 */}
      <Modal isOpen={isOpen} onOpenChange={handleClose} placement="top" size="xl">
        <ModalContent>
          <div className="px-4 pt-10 pb-4 space-y-6 bg-background rounded-xl shadow-sm">
            {/* 搜索框 */}
            <Input
              placeholder="输入关键词搜索..."
              startContent={<SearchLinearIcon size={18} />}
              size="lg"
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              classNames={{
                base: "rounded-md shadow-none",
                inputWrapper: "bg-default-100 dark:bg-default-500/20",
                input: "text-base",
              }}
            />

            {/* 搜索结果 */}
            <div className="max-h-[300px] overflow-y-auto flex flex-col">
              {isLoading ? (
                <div className="text-default-500 text-sm text-center py-4 animate-pulse">
                  正在加载...
                </div>
              ) : postList.length > 0 ? (
                postList.map((post) => (
                  <NextLink
                    key={post.id}
                    href={`/post/${post.id}`}
                    onClick={handleItemClick}
                  >
                    <div
                      className="
                        p-3 mb-2 rounded-xl bg-default-100/70 dark:bg-default-500/10
                        hover:bg-default-200 dark:hover:bg-white/10
                        transition-colors cursor-pointer shadow-sm
                        border border-transparent hover:border-default-300 dark:hover:border-zinc-800
                      "
                    >
                      <div className="text-base font-medium text-foreground line-clamp-2">
                        {highlightKeyword(post.title, keyword)}
                      </div>
                      {post.description && (
                        <div className="text-sm text-default-500 line-clamp-2 mt-1">
                          {highlightKeyword(post.description, keyword)}
                        </div>
                      )}
                    </div>
                  </NextLink>
                ))
              ) : (
                keyword && (
                  <div className="text-default-500 text-sm text-center py-4">
                    暂无匹配结果 🕵️‍♀️
                  </div>
                )
              )}
            </div>
          </div>
        </ModalContent>
      </Modal>
    </>
  );
};
