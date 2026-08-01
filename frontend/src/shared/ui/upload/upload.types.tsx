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

export interface ImageUploaderProps {
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
