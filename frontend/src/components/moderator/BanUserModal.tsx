// components/moderator/BanUserModal.tsx
import { useState } from "react";
import { Ban } from "lucide-react";

interface BanUserModalProps {
  boardId: number;
  onBan: (data: { userId: number; reason: string; expiresAt?: string }) => void;
  isBanning: boolean;
  t: (key: string) => string;
}

export function BanUserModal({
  boardId,
  onBan,
  isBanning,
  t,
}: BanUserModalProps) {
  const [userId, setUserId] = useState("");
  const [reason, setReason] = useState("");
  const [expiresAt, setExpiresAt] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onBan({
      userId: Number(userId),
      reason,
      expiresAt: expiresAt || undefined,
    });
    setUserId("");
    setReason("");
    setExpiresAt("");
  };

  return (
    <>
      <button
        className="btn btn-warning btn-sm"
        onClick={() =>
          (
            document.getElementById("ban_modal") as HTMLDialogElement
          )?.showModal()
        }
      >
        <Ban className="w-4 h-4" />
        {t("ban_user")}
      </button>

      <dialog id="ban_modal" className="modal">
        <div className="modal-box">
          <h3 className="font-bold text-lg mb-4">{t("ban_user")}</h3>
          <form onSubmit={handleSubmit} className="space-y-4">
            <input
              type="number"
              placeholder={t("user_id")}
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              className="input input-bordered w-full"
              required
            />
            <textarea
              placeholder={t("ban_reason")}
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              className="textarea textarea-bordered w-full"
              required
            />
            <input
              type="datetime-local"
              value={expiresAt}
              onChange={(e) => setExpiresAt(e.target.value)}
              className="input input-bordered w-full"
            />
            <div className="modal-action">
              <button
                type="button"
                className="btn"
                onClick={() =>
                  (
                    document.getElementById("ban_modal") as HTMLDialogElement
                  )?.close()
                }
              >
                {t("cancel")}
              </button>
              <button
                type="submit"
                className="btn btn-warning"
                disabled={isBanning}
              >
                {isBanning ? t("banning") : t("confirm_ban")}
              </button>
            </div>
          </form>
        </div>
        <form method="dialog" className="modal-backdrop">
          <button>{t("close")}</button>
        </form>
      </dialog>
    </>
  );
}
