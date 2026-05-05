"use client";

import React, { useSyncExternalStore } from "react";
import { pluginRegistry } from "../PluginRegistry";

interface PluginSlotProps {
  name: string;
  // 可以传额外 props 给插槽内所有组件
  slotProps?: Record<string, unknown>;
  // 没有插件时的 fallback
  fallback?: React.ReactNode;
  className?: string;
}

function subscribe(cb: () => void) {
  return pluginRegistry.subscribe(cb);
}

function getSnapshot(slotName: string) {
  return pluginRegistry.getSlotComponents(slotName);
}

/**
 * 在页面中预埋插槽，插件注册的组件会渲染在此处
 *
 * 用法：
 *   <PluginSlot name="sidebar-top" />
 *   <PluginSlot name="post-detail-bottom" slotProps={{ postId: "123" }} />
 */
export function PluginSlot({
  name,
  slotProps = {},
  fallback = null,
  className,
}: PluginSlotProps) {
  const components = useSyncExternalStore(
    subscribe,
    () => getSnapshot(name),
    () => [], // server snapshot
  );

  if (components.length === 0) return <>{fallback}</>;

  return (
    <div className={className} data-plugin-slot={name}>
      {components.map((item) => {
        const Component = item.component;
        return (
          <PluginComponentWrapper
            key={`${item.pluginId}-${name}`}
            pluginId={item.pluginId}
            pluginName={item.pluginName}
          >
            <Component {...(item.props ?? {})} {...slotProps} />
          </PluginComponentWrapper>
        );
      })}
    </div>
  );
}

/**
 * 错误边界包裹每个插件组件，防止单个插件崩溃影响整个页面
 */
class PluginComponentWrapper extends React.Component<
  { children: React.ReactNode; pluginId: string; pluginName: string },
  { hasError: boolean; error: string }
> {
  constructor(props: {
    children: React.ReactNode;
    pluginId: string;
    pluginName: string;
  }) {
    super(props);
    this.state = { hasError: false, error: "" };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error: error.message };
  }

  componentDidCatch(error: Error) {
    console.error(
      `[PluginSlot] Plugin "${this.props.pluginName}" crashed:`,
      error,
    );
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="text-xs text-error/60 px-2 py-1 border border-error/20 rounded bg-error/5">
          Plugin "{this.props.pluginName}" encountered an error
        </div>
      );
    }
    return this.props.children;
  }
}
