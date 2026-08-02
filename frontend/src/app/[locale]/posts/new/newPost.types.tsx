import { z } from "zod";

export type PostForm = z.infer<typeof postSchema>;
// ---------- 表单验证 ----------
export const postSchema = z
  .object({
    title: z.string().min(2, "标题至少2个字符").max(200, "标题最多200个字符"),
    content: z.string(),
    summary: z.string().max(500).optional(),
    cover: z.string().url("请输入有效的图片URL").optional().or(z.literal("")),
    type: z.enum(["post", "article", "topic"]),
    board_id: z.number().min(1, "请选择板块"),
    tag_ids: z.array(z.number()).max(5, "最多选择5个标签"),
    status: z
      .enum(["draft", "published", "pending", "hidden"])
      .default("published"),
  })
  .superRefine((data, ctx) => {
    if (
      data.status !== "draft" &&
      (!data.content || data.content.length < 10)
    ) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["content"],
        message: "内容至少10个字符",
      });
    }
  });
