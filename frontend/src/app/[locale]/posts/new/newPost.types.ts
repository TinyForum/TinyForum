import { z } from "zod";

function stripHtml(html: string): string {
  return html.replace(/<[^>]*>/g, "").trim();
}

function isValidImageUrl(val: string): boolean {
  if (!val) return true;
  return /^(\/|https?:\/\/)/.test(val);
}

export const WORK_TYPES = [
  { value: "image_text", label: "图文", maxImages: 9 },
  { value: "short_video", label: "短视频" },
  { value: "long_video", label: "长视频" },
  { value: "image", label: "图片", maxImages: 36 },
  { value: "article", label: "文章" },
  { value: "question", label: "问答" },
  { value: "topic", label: "话题", maxImages: 9 },
  { value: "post", label: "帖子", maxImages: 9 },
] as const;

export type PostForm = z.infer<typeof postSchema>;
export const postSchema = z
  .object({
    title: z.string().min(2, "标题至少2个字符").max(200, "标题最多200个字符"),
    content: z.string(),
    summary: z.string().max(500).optional(),
    cover: z.string().refine(isValidImageUrl, "请输入有效的图片URL").optional().or(z.literal("")),
    video_url: z.string().optional().or(z.literal("")),
    type: z.enum(["image_text", "short_video", "long_video", "image", "article", "question", "topic", "post"]),
    board_id: z.number().min(1, "请选择板块"),
    tag_ids: z.array(z.number()).max(5, "最多选择5个标签"),
    status: z.enum(["draft", "published", "pending", "hidden"]).default("published"),
  })
  .superRefine((data, ctx) => {
    if (data.status !== "draft" && stripHtml(data.content).length < 10) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["content"],
        message: "内容至少10个字符",
      });
    }
  });
