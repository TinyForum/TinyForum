"use client";

import Link from "next/link";
import { Fragment } from "react";
import { Menu, MenuButton, MenuItems, Transition } from "@headlessui/react";
import {
  Flame,
  Clock,
  PenSquare,
  Award,
  RefreshCw,
  Megaphone,
  ChevronDown,
} from "lucide-react";
import { useTranslations } from "next-intl";
import { SortBy } from "@/shared/type/posts.types";
import { useMediaQuery } from "@/features/common/hooks/useMediaQuery";

interface PostFilterBarProps {
  sortBy: SortBy;
  onSortChange: (sortBy: SortBy) => void;
  isAuthenticated: boolean;
  onRefetch: () => void;
  isLoading?: boolean;
  totalCount?: number;
}

export default function PostFilterBar({
  sortBy,
  onSortChange,
  isAuthenticated,
  onRefetch,
  isLoading = false,
  totalCount,
}: PostFilterBarProps) {
  const t = useTranslations("Post");
  const isDesktop = useMediaQuery("(min-width: 1024px)");

  const sortOptions: { value: SortBy; label: string; icon: React.ReactNode }[] =
    [
      {
        value: "random",
        label: t("random"),
        icon: <Clock className="w-4 h-4" />,
      },
      { value: "hot", label: t("hot"), icon: <Flame className="w-4 h-4" /> },
      { value: "like", label: t("like"), icon: <Award className="w-4 h-4" /> },
      {
        value: "latest",
        label: t("latest"),
        icon: <Megaphone className="w-4 h-4" />,
      },
    ];

  const currentOption = sortOptions.find((opt) => opt.value === sortBy);

  const DesktopView = () => (
    <div className="join">
      {sortOptions.map((option) => (
        <button
          key={option.value}
          className={`btn btn-sm join-item transition-all duration-200 ${
            sortBy === option.value
              ? "btn-primary shadow-md"
              : "btn-ghost hover:bg-base-200"
          }`}
          onClick={() => onSortChange(option.value)}
        >
          {option.icon}
          <span className="hidden sm:inline">{option.label}</span>
        </button>
      ))}
    </div>
  );

  const MobileView = () => (
    <div className="dropdown dropdown-end">
      <Menu>
        <MenuButton className="btn btn-sm btn-primary gap-1 whitespace-nowrap">
          {currentOption?.icon}
          <span>{currentOption?.label}</span>
          <ChevronDown className="w-3 h-3" />
        </MenuButton>
        <Transition
          as={Fragment}
          enter="transition ease-out duration-100"
          enterFrom="transform opacity-0 scale-95"
          enterTo="transform opacity-100 scale-100"
          leave="transition ease-in duration-75"
          leaveFrom="transform opacity-100 scale-100"
          leaveTo="transform opacity-0 scale-95"
        >
          <MenuItems className="dropdown-content z-50 mt-2 w-48 rounded-box bg-base-100 p-2 shadow-lg ring-1 ring-base-300 focus:outline-none">
            {sortOptions.map((option) => (
              <Menu.Item key={option.value}>
                {({ active }: { active: boolean }) => (
                  <button
                    className={`flex w-full items-center gap-2 rounded-btn px-4 py-2 text-sm transition-colors
                      ${active ? "bg-base-200" : ""}
                      ${
                        sortBy === option.value
                          ? "text-primary font-medium"
                          : "text-base-content"
                      }
                    `}
                    onClick={() => onSortChange(option.value)}
                  >
                    {option.icon}
                    {option.label}
                    {sortBy === option.value && (
                      <span className="ml-auto text-primary">✓</span>
                    )}
                  </button>
                )}
              </Menu.Item>
            ))}
          </MenuItems>
        </Transition>
      </Menu>
    </div>
  );

  return (
    <div className="relative mb-4 flex items-center gap-3 rounded-box border border-base-300 bg-base-100 p-3 shadow-sm">
      {/* 排序和刷新区域 */}
      <div className="flex flex-1 items-center gap-2 min-w-0">
        {/* 排序组件 */}
        {isDesktop ? <DesktopView /> : <MobileView />}

        {/* 刷新按钮 */}
        <button
          className="btn btn-sm btn-ghost gap-1 shrink-0"
          onClick={() => onRefetch()}
          disabled={isLoading}
          aria-label={t("refresh")}
        >
          <RefreshCw className={`h-4 w-4 ${isLoading ? "animate-spin" : ""}`} />
          <span className="hidden sm:inline">{t("refresh")}</span>
        </button>

        {/* 统计信息 */}
        {totalCount !== undefined && totalCount > 0 && (
          <div className="text-xs text-base-content/60 whitespace-nowrap hidden sm:block">
            {isDesktop ? t("total_posts", { count: totalCount }) : totalCount}
          </div>
        )}
      </div>

      {/* 发帖按钮 - 固定宽度 1/5 */}
      {isAuthenticated && (
        <div className="w-1/5 shrink-0 min-w-[80px]">
          <Link
            href="/posts/new"
            className="btn btn-primary btn-sm gap-1 w-full"
          >
            <PenSquare className="h-4 w-4" />
            <span className="hidden sm:inline">{t("create")}</span>
          </Link>
        </div>
      )}
    </div>
  );
}
