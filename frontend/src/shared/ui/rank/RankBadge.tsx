import { Crown, Medal } from "lucide-react";

// 排名徽章组件
export function RankBadge({
  rank,
  isTopThree,
}: {
  rank: number;
  isTopThree: boolean;
}) {
  if (!isTopThree) {
    return (
      <div className="w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm bg-base-200 text-base-content/50">
        {rank}
      </div>
    );
  }

  const config = {
    1: {
      bg: "bg-gradient-to-br from-yellow-400 to-yellow-500",
      text: "text-yellow-900",
      icon: Crown,
    },
    2: {
      bg: "bg-gradient-to-br from-gray-300 to-gray-400",
      text: "text-gray-700",
      icon: Medal,
    },
    3: {
      bg: "bg-gradient-to-br from-amber-500 to-amber-600",
      text: "text-white",
      icon: Medal,
    },
  };

  const { bg, text, icon: Icon } = config[rank as 1 | 2 | 3];

  return (
    <div
      className={`w-10 h-10 rounded-full flex items-center justify-center ${bg} ${text} shadow-md`}
    >
      <Icon className="w-5 h-5" />
    </div>
  );
}
