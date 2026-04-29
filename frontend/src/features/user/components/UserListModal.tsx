"use client";

import { X } from "lucide-react";
import type { User } from "@/shared/api/types";
import Avatar from "./Avatar";

interface UserListModalProps {
  title: string;
  users: User[];
  total: number;
  onClose: () => void;
  onUserClick: (userId: number) => void;
  isLoading?: boolean;
}

export function UserListModal({
  title,
  users,
  total,
  onClose,
  onUserClick,
  isLoading,
}: UserListModalProps) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
      onClick={onClose}
    >
      <div
        className="bg-base-100 rounded-2xl shadow-xl max-w-md w-full max-h-[80vh] flex flex-col"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between p-4 border-b border-base-200">
          <h3 className="text-lg font-semibold">
            {title} ({total})
          </h3>
          <button onClick={onClose} className="btn btn-sm btn-ghost btn-circle">
            <X className="w-4 h-4" />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto p-2">
          {isLoading ? (
            <div className="space-y-2 p-4">
              {[...Array(5)].map((_, i) => (
                <div key={i} className="flex items-center gap-3">
                  <div className="skeleton w-10 h-10 rounded-full" />
                  <div className="flex-1">
                    <div className="skeleton h-4 w-24 mb-1" />
                    <div className="skeleton h-3 w-32" />
                  </div>
                </div>
              ))}
            </div>
          ) : users.length === 0 ? (
            <div className="text-center py-8 text-base-content/40">
              暂无数据
            </div>
          ) : (
            <div className="space-y-1">
              {users.map((user) => (
                <button
                  key={user.id}
                  onClick={() => {
                    onUserClick(user.id);
                    onClose();
                  }}
                  className="w-full flex items-center gap-3 p-3 rounded-xl hover:bg-base-200 transition-colors text-left"
                >
                  <div className="avatar">
                    <div className="w-10 h-10 rounded-full">
                      <Avatar
                        username={user.username}
                        avatarUrl={user.avatar}
                        size="md"
                      />
                    </div>
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-sm truncate">
                      {user.username}
                    </p>
                    {user.bio && (
                      <p className="text-xs text-base-content/40 truncate">
                        {user.bio}
                      </p>
                    )}
                  </div>
                </button>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
