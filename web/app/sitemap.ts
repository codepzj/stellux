import type { MetadataRoute } from "next";
import { getSiteMapAPI } from "@/api/post";

const BASE_URL = process.env.NEXT_PUBLIC_SITE_URL;

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const res = await getSiteMapAPI();

  const staticPages: MetadataRoute.Sitemap = [
    {
      url: `${BASE_URL}/`,
      lastModified: new Date(),
      changeFrequency: "weekly",
      priority: 1.0,
    }
  ];

  const dynamicPages: MetadataRoute.Sitemap = res.data.map((item) => ({
    url: `${BASE_URL}/post/${item.id}`,
    lastModified: new Date(item.updated_at),
    changeFrequency: "weekly",
    priority: 0.8,
  }));

  return [...staticPages, ...dynamicPages];
}
