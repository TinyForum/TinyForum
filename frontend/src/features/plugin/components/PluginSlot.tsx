"use client";

import React, { useSyncExternalStore, useCallback } from "react";
import { pluginRegistry } from "../PluginRegistry";
import type { SlotComponent } from "@/shared/type/plugin.type";
import { useAdminPlugins } from "../useAdminPlugins";

interface PluginSlotProps {
  name: string;
  slotProps?: Record<string, unknown>;
  fallback?: React.ReactNode;
  className?: string;
}

// 全局唯一的空数组引用，解决无组件时每次返回新数组的问题
const STABLE_EMPTY_ARRAY: SlotComponent[] = [];

/**
 * 获取插槽组件的快照，保证返回稳定的引用
 */
function getSnapshot(slotName: string): SlotComponent[] {
  const components = pluginRegistry.getSlotComponents(slotName);
  // 如果没有组件，返回全局空数组（稳定引用）
  if (components.length === 0) {
    return STABLE_EMPTY_ARRAY;
  }
  // 有组件时直接返回注册表里缓存的数组（注册表内部会维护引用，只有变化时才更新）
  return components;
}

/**
 * 服务端渲染时使用的快照 – 始终返回稳定的空数组
 */
function getServerSnapshot(): SlotComponent[] {
  return STABLE_EMPTY_ARRAY;
}

export function PluginSlot({
  name,
  slotProps = {},
  fallback = null,
  className,
}: PluginSlotProps) {
  // 订阅函数：依赖 name 变化时会重新创建，但内部使用全局的 registry.subscribe
  console.log("PluginSlot render");
  const subscribe = useCallback(
    (onStoreChange: () => void) => {
      // 全局订阅，任何插件变化都会触发，但这是可接受的简洁方案
      const unsubscribe = pluginRegistry.subscribe(onStoreChange);
      return unsubscribe;
    },
    // 注意：name 变化时不需重新订阅，因为 registry 是全局的，只需订阅一次
    // 但为了符合规范，可以依赖空数组，但 useCallback 的依赖是 []
    [],
  );

  const components = useSyncExternalStore(
    subscribe,
    () => getSnapshot(name),
    getServerSnapshot,
  );

  if (components.length === 0) {
    return <>{fallback}</>;
  }

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

// 错误边界（保持不变）
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
