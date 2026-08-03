"use client";

import { UseFormRegister, FieldErrors } from "react-hook-form";

import toast from "react-hot-toast";
import { Board } from "@/shared/api/types/board.model";
import { Tag } from "@/shared/api/types/tag.model";
import { ImageUploader } from "@/shared/ui/upload/ImageUploader";
import { VideoUploader, VideoItem } from "@/shared/ui/upload/VideoUploader";
import { useCoverUpload } from "@/features/upload/hooks/useCoverUpload";
import { FolderOpen, X } from "lucide-react";
import { PostForm, WORK_TYPES } from "./newPost.types";
import { ImageItem } from "@/shared/ui/upload/upload.types";
import { PostType } from "@/shared/api/types/post.model";
import { uploadApi } from "@/shared/api/modules/uploads";

function getTypeConfig(type: PostType) {
  const found = WORK_TYPES.find((w) => w.value === type);
  return { maxImages: (found as { maxImages?: number })?.maxImages ?? 0 };
}

function isImageGridType(type: PostType) {
  return type === "image_text" || type === "image" || type === "topic" || type === "post";
}

function isVideoType(type: PostType) {
  return type === "short_video" || type === "long_video";
}

// ---------- 左侧：帖子设置组件（类型优化）----------
export interface PostSettingsProps {
  register: UseFormRegister<PostForm>;
  errors: FieldErrors<PostForm>;
  boards: Board[];
  tags: Tag[];
  boardsLoading: boolean;
  selectedTagIds: number[];
  selectedStatus: string;
  selectedType: PostType;
  coverValue: string;
  onToggleTag: (tagId: number) => void;
  onCoverChange: (url: string) => void;
  onVideoChange: (videoUrl: string) => void;
  onImageUrlsChange: (urls: string[]) => void;
  t: (key: string, params?: Record<string, string | number>) => string;
}

export function PostSettings({
  register,
  errors,
  boards,
  tags,
  boardsLoading,
  selectedTagIds,
  selectedStatus,
  selectedType,
  coverValue,
  onToggleTag,
  onCoverChange,
  onVideoChange,
  onImageUrlsChange,
  t,
}: PostSettingsProps) {
  const { uploadCover } = useCoverUpload();
  const typeConfig = getTypeConfig(selectedType);

  const statusOptions = [
    { value: "draft", label: t("status_draft"), desc: t("status_draft_desc") },
    {
      value: "published",
      label: t("status_published"),
      desc: t("status_published_desc"),
    },
    {
      value: "pending",
      label: t("status_pending"),
      desc: t("status_pending_desc"),
    },
    {
      value: "hidden",
      label: t("status_hidden"),
      desc: t("status_hidden_desc"),
    },
  ];

  // 封面上传函数（无需 postId）
  const handleUploadCover = async (file: File): Promise<{ url: string; file_id?: string }> => {
    try {
      return await uploadCover(file);
    } catch (error) {
      toast.error(t("cover_upload_failed"));
      throw error;
    }
  };

  return (
    <div className="space-y-5">
      {/* 板块选择 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">
            <FolderOpen className="w-4 h-4 inline mr-1" />
            {t("board")} <span className="text-error">*</span>
          </span>
        </label>
        <select
          {...register("board_id", {
            required: "请选择板块",
            valueAsNumber: true,
          })}
          className={`select select-bordered w-full focus:outline-none focus:border-primary ${
            errors.board_id ? "select-error" : ""
          }`}
          defaultValue=""
        >
          <option value="" disabled>
            {t("select_board")}
          </option>
          {boardsLoading ? (
            <option disabled>{t("loading")}</option>
          ) : (
            boards.map((board) => (
              <option key={board.id} value={board.id}>
                {board.name} {board.description ? `- ${board.description}` : ""}
              </option>
            ))
          )}
        </select>
        {errors.board_id && (
          <label className="label pt-1">
            <span className="label-text-alt text-error">
              {errors.board_id.message}
            </span>
          </label>
        )}
      </div>

      {/* 作品类型 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">{t("post_type")}</span>
        </label>
        <div className="grid grid-cols-3 gap-2">
          {WORK_TYPES.map((typeOption) => (
            <label key={typeOption.value} className="cursor-pointer">
              <input
                {...register("type")}
                type="radio"
                value={typeOption.value}
                className="hidden peer"
              />
              <div className="border-2 border-base-300 rounded-lg py-2 text-center peer-checked:border-primary peer-checked:bg-primary/5 transition-all">
                <div className="font-medium text-sm">{typeOption.label}</div>
              </div>
            </label>
          ))}
        </div>
      </div>

      {/* 状态 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">
            {t("status")} <span className="text-error">*</span>
          </span>
        </label>
        <select
          {...register("status")}
          className="select select-bordered w-full"
        >
          {statusOptions.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label} - {opt.desc}
            </option>
          ))}
        </select>
        {selectedStatus === "draft" && (
          <label className="label pt-1">
            <span className="label-text-alt text-info">{t("draft_hint")}</span>
          </label>
        )}
      </div>

      {/* 标题 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">
            {t("post_title")} <span className="text-error">*</span>
          </span>
        </label>
        <input
          {...register("title")}
          type="text"
          placeholder={t("post_title_placeholder")}
          className={`input input-bordered focus:outline-none focus:border-primary ${
            errors.title ? "input-error" : ""
          }`}
        />
        {errors.title && (
          <label className="label pt-1">
            <span className="label-text-alt text-error">
              {errors.title.message}
            </span>
          </label>
        )}
      </div>

      {/* 标签 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">
            {t("tags")}
            <span className="text-base-content/40 text-xs ml-2">
              {t("select_up_to_tags")}
            </span>
          </span>
        </label>
        <div className="flex flex-wrap gap-2">
          {tags.map((tag) => {
            const selected = selectedTagIds?.includes(tag.id);
            return (
              <button
                key={tag.id}
                type="button"
                onClick={() => onToggleTag(tag.id)}
                className={`badge badge-lg cursor-pointer transition-all ${
                  selected ? "ring-2" : "opacity-60 hover:opacity-100"
                }`}
                style={{
                  backgroundColor: selected
                    ? tag.color + "30"
                    : tag.color + "15",
                  color: tag.color,
                  borderColor: tag.color + "60",
                }}
              >
                {selected && <X className="w-3 h-3 mr-1" />}
                {tag.name}
              </button>
            );
          })}
        </div>
      </div>

      {/* 视频上传（短视频 / 长视频） */}
      {isVideoType(selectedType) && (
        <div className="form-control">
          <label className="label pb-1">
            <span className="label-text font-medium">{t("video_upload")}</span>
          </label>
          <VideoUploader
            onChange={(video: VideoItem | null) => {
              onVideoChange(video?.url || "");
            }}
          />
          <label className="label">
            <span className="label-text-alt text-base-content/40">
              {t("video_upload_hint")}
            </span>
          </label>
        </div>
      )}

      {/* 图片矩阵上传（图文 / 图片 / 话题 / 帖子） */}
      {isImageGridType(selectedType) && typeConfig.maxImages > 0 && (
        <div className="form-control">
          <label className="label pb-1">
            <span className="label-text font-medium">{t("image_upload")}</span>
          </label>
          <ImageUploader
            initialImages={[]}
            uploadFn={async (file: File) => {
              const res = await uploadApi.uploadFile({ file, type: "post_image" });
              return { url: res.data.data?.url ?? "", file_id: res.data.data?.file_id };
            }}
            maxCount={typeConfig.maxImages}
            supportCover={false}
            layout="grid"
            gridSize={3}
            onChange={(images: ImageItem[]) => {
              const urls = images
                .map((img) => img.url)
                .filter((url) => url && !url.startsWith("blob:"));
              onImageUrlsChange(urls);
            }}
          />
          <label className="label">
            <span className="label-text-alt text-base-content/40">
              {t("image_upload_hint", { max: typeConfig.maxImages })}
            </span>
          </label>
        </div>
      )}

      {/* 封面图（文章 / 问答） */}
      {!isVideoType(selectedType) && !isImageGridType(selectedType) && (
        <div className="form-control">
          <label className="label pb-1">
            <span className="label-text font-medium">{t("cover")}</span>
          </label>
          <ImageUploader
            initialImages={coverValue ? [{ url: coverValue }] : []}
            uploadFn={handleUploadCover}
            maxCount={1}
            supportCover={false}
            layout="grid"
            gridSize={2}
            onChange={(images: ImageItem[]) => {
              const coverUrl = images.length > 0 ? images[0].url : "";
              onCoverChange(coverUrl);
            }}
          />
          <label className="label">
            <span className="label-text-alt text-base-content/40">
              {t("cover_hint")}
            </span>
          </label>
        </div>
      )}

      {/* 摘要 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">
            {t("summary")}
            <span className="text-base-content/40 text-xs ml-2">
              {t("summary_desc")}
            </span>
          </span>
        </label>
        <textarea
          {...register("summary")}
          rows={2}
          placeholder={t("summary_placeholder")}
          className="textarea textarea-bordered focus:outline-none focus:border-primary resize-none"
        />
      </div>
    </div>
  );
}
