import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import rehypeHighlight from "rehype-highlight";
import "./md.scss";
import CopyButton from "./copy";
import React from "react";
import { cn } from "@/lib/utils";
import { Code } from "@heroui/code";
import { Link } from "@heroui/link";
export default function Md({ content, className }: { content: string, className?: string }) {
  let index = 1;

  return (
    <div className={cn("markdown-body", className)}>
      <ReactMarkdown
        rehypePlugins={[rehypeRaw, rehypeHighlight]}
        remarkPlugins={[remarkGfm]}
        components={{
          h2: ({ children }) => (
            <h2 id={`header-${index++}`}>{children}</h2>
          ),
          h3: ({ children }) => (
            <h3 id={`header-${index++}`} className="!border-none">{children}</h3>
          ),
          h4: ({ children }) => (
            <h4 className="!border-none">{children}</h4>
          ),
          h5: ({ children }) => (
            <h5 className="!border-none">{children}</h5>
          ),
          h6: ({ children }) => (
            <h6 className="!border-none">{children}</h6>
          ),
          p: ({ children }) => (
            <p className="leading-7">{children}</p>
          ),
          a: ({ children, href }) => (
            <Link
              href={href}
              isExternal
              className="text-primary"
            >
              {children}
            </Link>
          ),
          ul: ({ children }) => (
            <ul className="list-disc !ml-0 !pl-4">{children}</ul>
          ),
          pre: ({ children }) => (
            <pre className="rounded-md dark:bg-[#141414] !p-0">{children}</pre>
          ),
          code: ({ className, children }) => {
            const match = /language-(\w+)/.exec(className || "");
            const count =
              React.Children.toArray(children).length === 1
                ? (
                  React.Children.toArray(children)[0]
                    .toString()
                    .match(/\n/g) || []
                ).length
                : 0;

            if (match?.length || count > 0) {
              const id = Math.random().toString(36).slice(2, 11);
              // 多行代码块
              return (
                <div className="not-prose relative rounded-md">
                  <CopyButton id={id} />
                  <div className="overflow-x-auto p-4" id={id}>
                    {children}
                  </div>
                </div>
              );
            }
            // 单行代码块
            return (
              <Code radius="md">
                {children}
              </Code>
            );
          },
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}
