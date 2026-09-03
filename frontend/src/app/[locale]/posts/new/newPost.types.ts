import { z } from "zod";

function stripHtml(html: string): string {
  return html.replace(/<[^>]*>/g, "").trim();
}

function isValidImageUrl(val: string): boolean {
  if (!val) return true;
  return /^(\/|https?:\/\/)/.test(val);
}

export const WORK_TYPES = [
  { value: "image_text", labelKey: "work_type_image_text", maxImages: 9 },
  { value: "short_video", labelKey: "work_type_short_video" },
  { value: "long_video", labelKey: "work_type_long_video" },
  { value: "image", labelKey: "work_type_image", maxImages: 36 },
  { value: "article", labelKey: "article" },
  { value: "question", labelKey: "work_type_question" },
  { value: "topic", labelKey: "topic", maxImages: 9 },
  { value: "post", labelKey: "post", maxImages: 9 },
] as const;

export type PostForm = z.infer<ReturnType<typeof createPostSchema>>;

export function createPostSchema(t: (key: string) => string) {
  return z
    .object({
      title: z
        .string()
        .min(2, t("title_min_length"))
        .max(200, t("title_max_length")),
      content: z.string(),
      summary: z.string().max(500).optional(),
      cover: z
        .string()
        .refine(isValidImageUrl, t("invalid_image_url"))
        .optional()
        .or(z.literal("")),
      video_url: z.string().optional().or(z.literal("")),
      type: z.enum(["image_text", "short_video", "long_video", "image", "article", "question", "topic", "post"]),
      board_id: z.number().min(1, t("select_board_error")),
      tag_ids: z.array(z.number()).max(5, t("max_tags_error")),
      status: z.enum(["draft", "published", "pending", "hidden"]).default("published"),
    })
    .superRefine((data, ctx) => {
      if (data.status !== "draft" && stripHtml(data.content).length < 10) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["content"],
          message: t("content_min_length"),
        });
      }
    });
}
