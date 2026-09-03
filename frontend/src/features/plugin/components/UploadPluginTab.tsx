import { CheckCircle2 } from "lucide-react";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { useUpload } from "../hooks/useUpload";
import { useAdminPlugins } from "../useAdminPlugins";
import toast from "react-hot-toast";

export function UploadPluginTab() {
  const t = useTranslations("Plugin");
  const [file, setFile] = useState<File | null>(null);
  const [uploaded, setUploaded] = useState(false);
  const { isUploading, error, resetError } = useUpload();
  const { upload } = useAdminPlugins();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) {
      toast.error(t("select_zip"));
      return;
    }
    resetError();
    const result = await upload(file);
    if (result) {
      setUploaded(true);
      setFile(null);
    } else {
      toast.error(error || t("upload_failed"));
    }
  };

  if (uploaded) {
    return (
      <div className="card bg-base-200 p-6 text-center">
        <CheckCircle2 className="w-12 h-12 text-success mx-auto mb-3" />
        <h3 className="text-lg font-semibold">{t("submit_success")}</h3>
        <p className="text-sm text-base-content/60">
          {t("submit_success_desc")}
        </p>
        <button
          className="btn btn-sm btn-outline mt-4"
          onClick={() => setUploaded(false)}
        >
          {t("continue_upload")}
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-2xl">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <input
            type="file"
            accept=".zip"
            required
            onChange={(e) => setFile(e.target.files?.[0] || null)}
            className="file-input file-input-bordered w-full cursor-pointer"
          />
          <p className="text-xs text-base-content/50 mt-1">
            {t("zip_instruction")}
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
            {t("submit_review")}
          </button>
        </div>
      </form>
    </div>
  );
}
