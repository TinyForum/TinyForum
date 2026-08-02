"use client";

import { useState, useRef, useCallback } from "react";
import toast from "react-hot-toast";
import { uploadApi } from "@/shared/api/modules/uploads";
import { Play, X } from "lucide-react";

function normalizeUrl(url: string): string {
  if (!url) return url;
  if (url.startsWith("http://") || url.startsWith("https://") || url.startsWith("/")) return url;
  return "/" + url;
}

export interface VideoItem {
  id: string;
  url: string;
  file_id?: string;
  file?: File;
  uploading?: boolean;
  error?: string;
}

interface VideoUploaderProps {
  initialVideo?: { url: string; file_id?: string } | null;
  uploadFn?: (file: File) => Promise<{ url: string; file_id?: string }>;
  maxSizeMB?: number;
  onChange?: (video: VideoItem | null) => void;
  className?: string;
}

export function VideoUploader({
  initialVideo,
  uploadFn,
  maxSizeMB = 500,
  onChange,
  className = "",
}: VideoUploaderProps) {
  const [video, setVideo] = useState<VideoItem | null>(() =>
    initialVideo?.url
      ? { id: "initial-video", url: initialVideo.url, file_id: initialVideo.file_id }
      : null,
  );
  const fileInputRef = useRef<HTMLInputElement>(null);

  const notifyChange = useCallback(
    (v: VideoItem | null) => {
      onChange?.(v);
    },
    [onChange],
  );

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > maxSizeMB * 1024 * 1024) {
      toast.error(`视频文件不能超过 ${maxSizeMB}MB`);
      return;
    }

    const newVideo: VideoItem = {
      id: `temp-${Date.now()}`,
      url: URL.createObjectURL(file),
      file,
      uploading: true,
    };
    setVideo(newVideo);
    notifyChange(newVideo);

    const upload = uploadFn || (async (f: File) => {
      const res = await uploadApi.uploadFile({ file: f, type: "video" });
      return { url: res.data.data?.url ?? "", file_id: res.data.data?.file_id };
    });

    try {
      const { url, file_id } = await upload(file);
      const uploaded: VideoItem = { ...newVideo, url, file_id, uploading: false, file: undefined };
      setVideo(uploaded);
      notifyChange(uploaded);
    } catch {
      const failed: VideoItem = { ...newVideo, uploading: false, error: "上传失败" };
      setVideo(failed);
      notifyChange(failed);
    }
  };

  const handleRemove = () => {
    if (video?.url?.startsWith("blob:")) URL.revokeObjectURL(video.url);
    setVideo(null);
    notifyChange(null);
  };

  return (
    <div className={className}>
      {video ? (
        <div className="relative rounded-lg overflow-hidden border border-base-300 bg-base-100">
          <div className="relative aspect-video bg-base-200">
            <video
              src={normalizeUrl(video.url)}
              controls
              className="w-full h-full object-contain"
              preload="metadata"
            />
            {video.uploading && (
              <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                <span className="loading loading-spinner loading-lg text-white" />
              </div>
            )}
            {video.error && (
              <div className="absolute inset-0 bg-error/80 flex items-center justify-center text-white text-sm">
                {video.error}
              </div>
            )}
          </div>
          <button
            type="button"
            onClick={handleRemove}
            className="absolute top-2 right-2 btn btn-xs btn-circle btn-error"
          >
            <X className="w-3 h-3" />
          </button>
        </div>
      ) : (
        <button
          type="button"
          onClick={() => fileInputRef.current?.click()}
          className="w-full flex flex-col items-center justify-center rounded-lg border-2 border-dashed border-base-300 hover:border-primary transition bg-base-100 py-10 gap-2"
        >
          <Play className="w-10 h-10 text-base-content/30" />
          <span className="text-sm text-base-content/50">点击上传视频</span>
          <span className="text-xs text-base-content/30">支持 MP4、WebM 等格式，最大 {maxSizeMB}MB</span>
        </button>
      )}
      <input
        ref={fileInputRef}
        type="file"
        accept="video/*"
        onChange={handleFileSelect}
        className="hidden"
      />
    </div>
  );
}
