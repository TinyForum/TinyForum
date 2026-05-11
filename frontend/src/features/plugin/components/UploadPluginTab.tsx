import { CheckCircle2 } from "lucide-react";
import { useState } from "react";
import { useUpload } from "../hooks/useUpload";
import { useAdminPlugins } from "../useAdminPlugins";

export function UploadPluginTab() {
  const [file, setFile] = useState<File | null>(null);
  const [uploaded, setUploaded] = useState(false);
  const { isUploading, error, resetError } = useUpload();
  const { upload } = useAdminPlugins();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) {
      alert("请选择插件 ZIP 包");
      return;
    }
    resetError();
    const result = await upload(file);
    if (result) {
      setUploaded(true);
      setFile(null);
    } else {
      alert(error || "上传失败，请稍后重试");
    }
  };

  if (uploaded) {
    return (
      <div className="card bg-base-200 p-6 text-center">
        <CheckCircle2 className="w-12 h-12 text-success mx-auto mb-3" />
        <h3 className="text-lg font-semibold">提交成功！</h3>
        <p className="text-sm text-base-content/60">
          插件已提交审核，通过后将在插件市场上架。
        </p>
        <button
          className="btn btn-sm btn-outline mt-4"
          onClick={() => setUploaded(false)}
        >
          继续上传
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-2xl">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="label-text">插件 ZIP 包 *</label>
          <input
            type="file"
            accept=".zip"
            required
            onChange={(e) => setFile(e.target.files?.[0] || null)}
            className="file-input file-input-bordered w-full"
          />
          <p className="text-xs text-base-content/50 mt-1">
            请打包符合规范的插件目录为 ZIP 文件，系统将自动读取插件信息。
          </p>
        </div>
        <div className="flex justify-end gap-2">
          <button
            type="submit"
            className="btn btn-secondary btn-sm"
            disabled={isUploading}
          >
            {isUploading && (
              <span className="loading loading-spinner loading-xs" />
            )}
            提交审核
          </button>
        </div>
      </form>
    </div>
  );
}
