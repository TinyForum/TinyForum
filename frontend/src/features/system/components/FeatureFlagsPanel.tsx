"use client";

import { Switch } from "@headlessui/react";
import { CheckCircle2, Circle } from "lucide-react";
import type { Feature } from "../hooks/useFeatureFlags";

interface FeatureFlagsPanelProps {
  grouped: Record<string, Feature[]>;
  enabledCount: number;
  total: number;
  togglingId: string | null;
  onToggle: (id: string, enabled: boolean) => void;
  onEnableAll: () => void;
}

export function FeatureFlagsPanel({
  grouped,
  enabledCount,
  total,
  togglingId,
  onToggle,
  onEnableAll,
}: FeatureFlagsPanelProps) {
  const pct = total > 0 ? Math.round((enabledCount / total) * 100) : 0;

  return (
    <div className="space-y-5">
      {/* Stats bar */}
      <div className="card bg-base-100 border border-base-300 shadow-sm">
        <div className="card-body py-4 flex-row items-center gap-6">
          <div className="flex-1 space-y-1.5">
            <div className="flex items-center justify-between text-sm">
              <span className="text-base-content/60">功能启用率</span>
              <span className="font-bold tabular-nums">
                {enabledCount} / {total}
              </span>
            </div>
            <div className="h-2 w-full bg-base-200 rounded-full overflow-hidden">
              <div
                className="h-full bg-primary rounded-full transition-all duration-500"
                style={{ width: `${pct}%` }}
              />
            </div>
          </div>
          <div className="text-right shrink-0">
            <p className="text-2xl font-bold text-primary tabular-nums">
              {pct}%
            </p>
            <p className="text-xs text-base-content/40">启用率</p>
          </div>
          <button onClick={onEnableAll} className="btn btn-ghost btn-sm">
            全部启用
          </button>
        </div>
      </div>

      {/* Grouped feature cards */}
      {Object.entries(grouped).map(([category, features]) => (
        <div
          key={category}
          className="card bg-base-100 border border-base-300 shadow-sm"
        >
          <div className="card-body gap-3">
            <div className="flex items-center gap-2 pb-2 border-b border-base-200">
              <h3 className="font-semibold text-sm">{category}</h3>
              <span className="badge badge-ghost badge-xs">
                {features.filter((f) => f.enabled).length}/{features.length}
              </span>
            </div>
            <div className="space-y-2">
              {features.map((feature) => (
                <FeatureRow
                  key={feature.id}
                  feature={feature}
                  isToggling={togglingId === feature.id}
                  onToggle={onToggle}
                />
              ))}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

function FeatureRow({
  feature,
  isToggling,
  onToggle,
}: {
  feature: Feature;
  isToggling: boolean;
  onToggle: (id: string, enabled: boolean) => void;
}) {
  return (
    <div
      className={`flex items-center justify-between rounded-lg px-3 py-2.5 transition-colors ${
        feature.enabled ? "bg-primary/5" : "bg-base-200/60"
      }`}
    >
      <div className="flex items-start gap-2.5 flex-1 min-w-0">
        <span
          className={`mt-0.5 shrink-0 ${feature.enabled ? "text-primary" : "text-base-content/20"}`}
        >
          {feature.enabled ? (
            <CheckCircle2 className="w-4 h-4" />
          ) : (
            <Circle className="w-4 h-4" />
          )}
        </span>
        <div className="min-w-0">
          <p className="text-sm font-medium truncate">{feature.name}</p>
          <p className="text-xs text-base-content/50 truncate">
            {feature.description}
          </p>
        </div>
      </div>
      <div className="shrink-0 ml-3">
        {isToggling ? (
          <span className="loading loading-spinner loading-xs text-primary" />
        ) : (
          <Switch
            checked={feature.enabled}
            onChange={(v) => onToggle(feature.id, v)}
            className={`${
              feature.enabled ? "bg-primary" : "bg-base-300"
            } relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-primary/30 focus:ring-offset-1`}
          >
            <span
              className={`${
                feature.enabled ? "translate-x-5" : "translate-x-1"
              } inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow transition-transform`}
            />
          </Switch>
        )}
      </div>
    </div>
  );
}
