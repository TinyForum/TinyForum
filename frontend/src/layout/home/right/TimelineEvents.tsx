"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import { formatDistanceToNow } from "date-fns";
import { zhCN } from "date-fns/locale";
import {
  Zap,
  FileText,
  MessageSquare,
  Heart,
  UserPlus,
  CheckCircle,
  Activity,
} from "lucide-react";

interface TimelineEvent {
  id: number;
  action: string;
  target_type: string;
  target_id: number;
  created_at: string;
  actor?: {
    id: number;
    username: string;
    avatar: string;
  };
}

interface TimelineEventsProps {
  timelineEvents: TimelineEvent[];
}

export function TimelineEvents({ timelineEvents }: TimelineEventsProps) {
  const t = useTranslations("Sidebar");

  const getEventIcon = (action: string) => {
    switch (action) {
      case "create_post":
        return <FileText className="w-3 h-3" />;
      case "create_comment":
        return <MessageSquare className="w-3 h-3" />;
      case "like_post":
      case "like_comment":
        return <Heart className="w-3 h-3" />;
      case "follow_user":
        return <UserPlus className="w-3 h-3" />;
      case "accept_answer":
        return <CheckCircle className="w-3 h-3" />;
      default:
        return <Activity className="w-3 h-3" />;
    }
  };

  const getEventText = (event: TimelineEvent) => {
    const actor = event.actor?.username || t("user");
    switch (event.action) {
      case "create_post":
        return `${actor} ${t("posted_a_new_thread")}`;
      case "create_comment":
        return `${actor} ${t("commented")}`;
      case "like_post":
        return `${actor} ${t("liked_a_post")}`;
      case "like_comment":
        return `${actor} ${t("liked_a_comment")}`;
      case "follow_user":
        return `${actor} ${t("followed_you")}`;
      case "accept_answer":
        return `${actor} ${t("accepted_your_answer")}`;
      default:
        return `${actor} ${t("has_a_new_activity")}`;
    }
  };

  return (
    <div className="rounded-lg border bg-card">
      <div className="p-3 border-b">
        <h3 className="font-semibold flex items-center gap-2">
          <Zap className="w-4 h-4 text-blue-500" />
          {t("recent_updates")}
        </h3>
      </div>
      <div className="p-2 space-y-1 max-h-[300px] overflow-y-auto">
        {timelineEvents.slice(0, 5).map((event) => (
          <div
            key={event.id}
            className="flex items-start gap-2 p-2 rounded-lg hover:bg-muted transition-colors text-sm"
          >
            <div className="flex-shrink-0 mt-0.5 text-muted-foreground">
              {getEventIcon(event.action)}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm truncate">{getEventText(event)}</p>
              <p className="text-xs text-muted-foreground">
                {formatDistanceToNow(new Date(event.created_at), {
                  addSuffix: true,
                  locale: zhCN,
                })}
              </p>
            </div>
          </div>
        ))}
      </div>
      <div className="p-2 border-t">
        <Link
          href="/timeline"
          className="block text-xs text-center text-muted-foreground hover:text-primary"
        >
          查看更多动态 →
        </Link>
      </div>
    </div>
  );
}
