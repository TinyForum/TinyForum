// ImageUploader.tsx
import { uploadApi } from "@/shared/api/modules/uploads";
import Image from "next/image";
import React, { useState, useRef, useEffect } from "react";
import toast from "react-hot-toast";

export type LayoutMode = "grid" | "waterfall" | "horizontal" | "tile";
export type GridSize = 2 | 3 | 4;

export interface ImageItem {
  id: string;
  url: string;
  file?: File;
  isCover?: boolean;
  uploading?: boolean;
  error?: string;
}

interface ImageUploaderProps {
  initialImages?: Array<{ url: string; isCover?: boolean }>;
  uploadFn?: (file: File) => Promise<{ url: string }>;
  maxCount?: number;
  supportCover?: boolean;
  layout?: LayoutMode;
  gridSize?: GridSize;
  defaultCollapsed?: boolean;
  onChange?: (images: ImageItem[]) => void;
  onDelete?: (image: ImageItem) => void;
  className?: string;
  onLayoutChange?: (layout: LayoutMode) => void; // 新增回调
}

export const ImageUploader: React.FC<ImageUploaderProps> = ({
  initialImages = [],
  uploadFn,
  maxCount = 9,
  supportCover = true,
  layout: initialLayout = "grid",
  gridSize = 3,
  defaultCollapsed = false,
  onChange,
  onDelete,
  className = "",
  onLayoutChange,
}) => {
  const [images, setImages] = useState<ImageItem[]>(() =>
    initialImages.map((img, idx) => ({
      id: `initial-${idx}`,
      url: img.url,
      isCover: img.isCover || (idx === 0 && supportCover),
    })),
  );
  const [collapsed, setCollapsed] = useState(defaultCollapsed);
  const [layout, setLayout] = useState<LayoutMode>(initialLayout);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // 确保最多只有一张封面
  useEffect(() => {
    if (supportCover) {
      const coverCount = images.filter((img) => img.isCover).length;
      if (coverCount === 0 && images.length > 0) {
        setImages((prev) =>
          prev.map((img, idx) => ({ ...img, isCover: idx === 0 })),
        );
      } else if (coverCount > 1) {
        let firstCoverSet = false;
        setImages((prev) =>
          prev.map((img) => {
            if (img.isCover) {
              if (!firstCoverSet) {
                firstCoverSet = true;
                return img;
              }
              return { ...img, isCover: false };
            }
            return img;
          }),
        );
      }
    }
  }, [images, supportCover]);

  useEffect(() => {
    onChange?.(images);
  }, [images, onChange]);

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);
    if (images.length + files.length > maxCount) {
      toast.error(`最多上传 ${maxCount} 张图片`);
      return;
    }

    const newImages: ImageItem[] = files.map((file, idx) => ({
      id: `temp-${Date.now()}-${idx}`,
      url: URL.createObjectURL(file),
      file,
      uploading: true,
    }));
    setImages((prev) => [...prev, ...newImages]);

    for (const item of newImages) {
      try {
        const upload =
          uploadFn ||
          ((f: File) =>
            uploadApi
              .uploadPostFile("postId", f)
              .then((res) => ({ url: res.data.data })));
        const { url } = await upload(item.file!);
        setImages((prev) =>
          prev.map((img) =>
            img.id === item.id
              ? { ...img, url, uploading: false, file: undefined }
              : img,
          ),
        );
      } catch {
        setImages((prev) =>
          prev.map((img) =>
            img.id === item.id
              ? { ...img, uploading: false, error: "上传失败" }
              : img,
          ),
        );
      }
    }
  };

  const handleDelete = (image: ImageItem) => {
    if (image.url && image.url.startsWith("blob:")) {
      URL.revokeObjectURL(image.url);
    }
    const newImages = images.filter((img) => img.id !== image.id);
    setImages(newImages);
    onDelete?.(image);
  };

  const handleSetCover = (imageId: string) => {
    if (!supportCover) return;
    setImages((prev) =>
      prev.map((img) => ({ ...img, isCover: img.id === imageId })),
    );
  };

  const handleDragStart = (e: React.DragEvent, index: number) => {
    e.dataTransfer.setData("text/plain", String(index));
  };
  const handleDragOver = (e: React.DragEvent) => e.preventDefault();
  const handleDrop = (e: React.DragEvent, targetIndex: number) => {
    const sourceIndex = Number(e.dataTransfer.getData("text/plain"));
    if (sourceIndex === targetIndex) return;
    setImages((prev) => {
      const newImages = [...prev];
      const [moved] = newImages.splice(sourceIndex, 1);
      newImages.splice(targetIndex, 0, moved);
      return newImages;
    });
  };

  const handleLayoutChange = (newLayout: LayoutMode) => {
    setLayout(newLayout);
    onLayoutChange?.(newLayout);
  };

  const getLayoutStyle = (): React.CSSProperties => {
    switch (layout) {
      case "grid":
        return {
          display: "grid",
          gridTemplateColumns: `repeat(${gridSize}, minmax(0, 1fr))`,
          gap: "1rem",
        };
      case "waterfall":
        return {
          columnCount: gridSize,
          columnGap: "1rem",
        };
      case "horizontal":
        return {
          display: "flex",
          overflowX: "auto",
          gap: "1rem",
          paddingBottom: "0.5rem",
        };
      case "tile":
        return {
          display: "flex",
          flexWrap: "wrap",
          gap: "1rem",
        };
      default:
        return {};
    }
  };

  const renderImageItem = (image: ImageItem, index: number) => {
    const isHorizontal = layout === "horizontal";
    const itemStyle: React.CSSProperties = isHorizontal
      ? { flex: "0 0 auto", width: "200px" }
      : { breakInside: "avoid", marginBottom: "1rem" };

    return (
      <div
        key={image.id}
        style={itemStyle}
        draggable={layout !== "waterfall"}
        onDragStart={(e) => handleDragStart(e, index)}
        onDragOver={handleDragOver}
        onDrop={(e) => handleDrop(e, index)}
        className="relative group rounded-lg overflow-hidden border border-base-300 bg-base-100 shadow-sm"
      >
        <div className="relative aspect-square">
          <Image
            src={image.url}
            alt="预览"
            className="w-full h-full object-cover"
          />
          {image.uploading && (
            <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
              <span className="loading loading-spinner loading-md text-white"></span>
            </div>
          )}
          {image.error && (
            <div className="absolute inset-0 bg-error/80 flex items-center justify-center text-white text-sm">
              {image.error}
            </div>
          )}
        </div>

        <div className="absolute top-2 right-2 flex gap-1 opacity-0 group-hover:opacity-100 transition">
          {supportCover && !image.isCover && (
            <button
              type="button"
              onClick={() => handleSetCover(image.id)}
              className="btn btn-xs btn-circle btn-primary"
              title="设为封面"
            >
              🌟
            </button>
          )}
          <button
            type="button"
            onClick={() => handleDelete(image)}
            className="btn btn-xs btn-circle btn-error"
            title="删除"
          >
            ✕
          </button>
        </div>

        {supportCover && image.isCover && (
          <div className="absolute bottom-2 left-2 bg-primary text-primary-content text-xs px-2 py-0.5 rounded-full">
            封面
          </div>
        )}
      </div>
    );
  };

  const addButton = (
    <button
      type="button"
      onClick={() => fileInputRef.current?.click()}
      className="flex flex-col items-center justify-center aspect-square rounded-lg border-2 border-dashed border-base-300 hover:border-primary transition bg-base-100"
    >
      <svg
        className="w-8 h-8 text-base-content/40"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M12 4v16m8-8H4"
        />
      </svg>
      <span className="text-xs text-base-content/60 mt-1">上传图片</span>
    </button>
  );

  return (
    <div className={`w-full ${className}`}>
      <div className="flex justify-between items-center mb-3">
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setCollapsed(!collapsed)}
            className="btn btn-sm btn-ghost"
          >
            {collapsed ? "展开" : "折叠"}
          </button>
          <select
            value={layout}
            onChange={(e) => handleLayoutChange(e.target.value as LayoutMode)}
            className="select select-bordered select-sm w-32"
          >
            <option value="grid">网格</option>
            <option value="waterfall">瀑布流</option>
            <option value="horizontal">横向滚动</option>
            <option value="tile">平铺</option>
          </select>
        </div>
        <span className="text-sm text-base-content/60">
          {images.length}/{maxCount}
        </span>
      </div>

      {!collapsed && (
        <div>
          <div style={getLayoutStyle()}>
            {images.map((img, idx) => renderImageItem(img, idx))}
            {images.length < maxCount && addButton}
          </div>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            multiple
            onChange={handleFileSelect}
            className="hidden"
          />
        </div>
      )}
    </div>
  );
};
