// components/topic/TopicCard.tsx
"use client";

import { useAuthStore } from "@/store/auth";
import {
  DocumentTextIcon,
  UserGroupIcon,
  GlobeAltIcon,
  LockClosedIcon,
  HeartIcon,
} from "@heroicons/react/24/outline";
import { HeartIcon as HeartSolidIcon } from "@heroicons/react/24/solid";
import Link from "next/link";
import Image from "next/image";
import { useState } from "react";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { Topic } from "@/shared/api/types/topic.model";
import { useFollowTopic } from "@/features/topic/hooks/useFollowTopic";

// 类型定义
interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

// 话题卡片组件
export function TopicCard({
  topic,
  onFollowChange,
}: {
  topic: Topic;
  onFollowChange?: () => void;
}) {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const { follow: followMutation, unfollow: unfollowMutation } =
    useFollowTopic();
  // 修复：topic.is_public 表示话题是否公开，而不是是否已关注
  // 应该从 topic 中获取关注状态，或者通过 API 获取
  const [following, setFollowing] = useState<boolean>(false);
  const [followLoading, setFollowLoading] = useState<boolean>(false);

  const t = useTranslations("Topic");

  const handleFollow = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (!isAuthenticated) {
      toast.error(t("please_login"));
      router.push("/login");
      return;
    }

    setFollowLoading(true);
    try {
      if (following) {
        await unfollowMutation.mutateAsync(topic.id);
        setFollowing(false);
        toast.success(t("unfollow_success"));
      } else {
        await followMutation.mutateAsync(topic.id);
        setFollowing(true);
        toast.success(t("follow_success"));
      }
      onFollowChange?.();
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      toast.error(error.response?.data?.message || t("operation_failed"));
    } finally {
      setFollowLoading(false);
    }
  };

  return (
    <Link href={`/topics/${topic.id}`}>
      <div className="card bg-base-100 shadow-sm hover:shadow-lg transition-all duration-300 cursor-pointer group">
        {/* 话题封面 */}
        {topic.cover && (
          <figure className="relative overflow-hidden">
            <Image
              src={topic.cover}
              alt={topic.title}
              fill
              sizes="(max-width: 768px) 100vw, 33vw"
              className="object-cover group-hover:scale-105 transition-transform duration-500"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
          </figure>
        )}

        <div className="card-body p-5">
          {/* 话题标题 */}
          <h3 className="card-title text-xl font-bold text-base-content group-hover:text-primary transition-colors">
            {topic.title}
          </h3>

          {/* 话题描述 */}
          {topic.description && (
            <p className="text-base-content/60 line-clamp-2 leading-relaxed">
              {topic.description}
            </p>
          )}

          {/* 统计信息 */}
          <div className="flex flex-wrap items-center gap-3 text-sm text-base-content/50 mt-2">
            <div className="flex items-center gap-1.5">
              <DocumentTextIcon className="w-4 h-4" />
              <span>{t("post_count", { count: topic.post_count || 0 })}</span>
            </div>
            <div className="flex items-center gap-1.5">
              <UserGroupIcon className="w-4 h-4" />
              <span>{t("follower_count", { count: topic.follower_count || 0 })}</span>
            </div>
            <div className="flex items-center gap-1.5">
              {topic.is_public ? (
                <>
                  <GlobeAltIcon className="w-4 h-4" />
                  <span className="badge badge-ghost badge-xs">{t("public")}</span>
                </>
              ) : (
                <>
                  <LockClosedIcon className="w-4 h-4" />
                  <span className="badge badge-ghost badge-xs">{t("private")}</span>
                </>
              )}
            </div>
          </div>

          {/* 创建者信息 */}
          {topic.creator && (
            <div className="flex items-center gap-2 mt-3 pt-3 border-t border-base-200">
              <div className="avatar placeholder">
                <div className="w-5 h-5 rounded-full bg-primary/10 text-primary">
                  <span className="text-xs font-medium">
                    {topic.creator.username?.[0]?.toUpperCase() || "U"}
                  </span>
                </div>
              </div>
              <span className="text-xs text-base-content/40">
                {t("created_by", { username: topic.creator.username || t("user_id_format", { id: topic.creator_id }) })}
              </span>
            </div>
          )}

          {/* 操作按钮区域 */}
          <div className="card-actions justify-end mt-4">
            <button
              onClick={handleFollow}
              disabled={followLoading}
              className={`btn btn-sm gap-1.5 ${
                following ? "btn-primary btn-outline" : "btn-primary"
              }`}
            >
              {followLoading ? (
                <span className="loading loading-spinner loading-xs"></span>
              ) : following ? (
                <HeartSolidIcon className="w-4 h-4" />
              ) : (
                <HeartIcon className="w-4 h-4" />
              )}
              {following ? t("followed") : t("follow")}
            </button>
          </div>
        </div>
      </div>
    </Link>
  );
}
