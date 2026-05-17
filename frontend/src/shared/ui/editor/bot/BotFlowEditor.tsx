import { useState, useCallback, useMemo } from "react";
import {
  DndContext,
  closestCenter,
  DragEndEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
  useSortable,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { CreateBotRequest } from "@/shared/api/types/bot.model";
import { useBotActions } from "@/features/bot/hooks/bot";
import {
  useNocodeMetadata,
  useValidateFlow,
} from "@/features/bot/hooks/useNocodeMetadata";
import { NodeMeta, NocodeMetadata } from "@/features/bot/noco.type";

// ---------- 类型定义（与后端 Go 模型对齐）----------
export type TriggerType =
  | "on_schedule"
  | "on_new_post"
  | "on_new_comment"
  | "on_user_register"
  | "on_keyword"
  | "on_manual";

export type CondType =
  | "post_title_contains"
  | "post_content_contains"
  | "user_role_is"
  | "user_post_count_gte"
  | "board_id_in"
  | "time_range"
  | "custom_expr";

export type ActionType =
  | "reply_post"
  | "delete_post"
  | "hide_post"
  | "pin_post"
  | "lock_post"
  | "create_post"
  | "delete_comment"
  | "ban_user"
  | "send_message"
  | "webhook"
  | "notify_admin"
  | "wait"
  | "set_variable"
  | "stop_if";

export interface TriggerNode {
  type: TriggerType;
  params?: Record<string, unknown>;
}

export interface CondNode {
  type: CondType;
  negate?: boolean;
  params: Record<string, unknown>;
}

export interface ActionNode {
  type: ActionType;
  params: Record<string, unknown>;
}

export interface Flow {
  version: string;
  trigger: TriggerNode;
  conditions?: CondNode[];
  actions: ActionNode[];
}

// ---------- 辅助函数 ----------
// function _getNodeLabel(nodeMeta: NodeMeta): string {
//   console.log("获取节点标签: ", nodeMeta);
//   return nodeMeta.label;
// }

/**
 * 创建默认节点参数的函数
 * @param nodeMeta - 节点的元数据信息，包含schema等信息
 * @returns 返回一个包含默认参数的Record对象
 */
function createDefaultParams(nodeMeta: NodeMeta): Record<string, unknown> {
  console.log("创建默认节点参数: ", nodeMeta); // 输出创建默认参数时的节点元数据信息
  // 可从 nodeMeta.schema 生成默认值，此处简单返回空对象
  return {}; // 返回一个空对象作为默认参数
}

// ---------- 可拖拽卡片组件（SortableItem）----------
interface SortableItemProps {
  id: string;
  typeLabel: string;
  onEdit: () => void;
  onDelete: () => void;
}

function SortableItem({ id, typeLabel, onEdit, onDelete }: SortableItemProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className="bg-white border rounded p-3 flex items-center justify-between shadow-sm"
    >
      <div className="flex items-center gap-2 flex-1">
        <div
          {...attributes}
          {...listeners}
          className="cursor-grab text-gray-400 hover:text-gray-600"
        >
          ⋮⋮
        </div>
        <span className="font-medium">{typeLabel}</span>
      </div>
      <div className="flex gap-2">
        <button
          onClick={onEdit}
          className="text-blue-600 hover:text-blue-800 text-sm"
        >
          配置
        </button>
        <button
          onClick={onDelete}
          className="text-red-600 hover:text-red-800 text-sm"
        >
          删除
        </button>
      </div>
    </div>
  );
}

// ---------- 配置弹窗（编辑节点的 params）----------
interface ConfigModalProps {
  title: string;
  params: Record<string, unknown>;
  onSave: (newParams: Record<string, unknown>) => void;
  onClose: () => void;
}

function ConfigModal({ title, params, onSave, onClose }: ConfigModalProps) {
  const [config, setConfig] = useState<Record<string, unknown>>(params);

  const handleSave = () => {
    onSave(config);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-4 rounded shadow-lg w-96">
        <h3 className="font-bold text-lg mb-2">配置 {title}</h3>
        <textarea
          className="w-full border p-2 mt-2 font-mono text-sm"
          rows={6}
          value={JSON.stringify(config, null, 2)}
          onChange={(e) => {
            try {
              const parsed = JSON.parse(e.target.value) as Record<
                string,
                unknown
              >;
              setConfig(parsed);
            } catch {
              // 保留原文本，不更新
            }
          }}
        />
        <div className="flex justify-end mt-4 space-x-2">
          <button
            className="px-3 py-1 border rounded hover:bg-gray-100"
            onClick={onClose}
          >
            取消
          </button>
          <button
            className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
            onClick={handleSave}
          >
            保存
          </button>
        </div>
      </div>
    </div>
  );
}

// ---------- 折叠面板组件 ----------
interface CollapsibleSectionProps {
  title: string;
  children: React.ReactNode;
  defaultOpen?: boolean;
}

function CollapsibleSection({
  title,
  children,
  defaultOpen = true,
}: CollapsibleSectionProps) {
  const [isOpen, setIsOpen] = useState(defaultOpen);
  return (
    <div className="mb-4">
      <div
        className="flex items-center justify-between cursor-pointer py-1 select-none"
        onClick={() => setIsOpen(!isOpen)}
      >
        <h3 className="font-bold text-gray-800">{title}</h3>
        <span className="text-gray-500 text-sm">{isOpen ? "▼" : "▶"}</span>
      </div>
      {isOpen && <div className="mt-2 space-y-2">{children}</div>}
    </div>
  );
}

// ---------- 主组件 ----------
export function BotFlowEditor() {
  const {
    metadata,
    loading: metaLoading,
    error: metaError,
  } = useNocodeMetadata();
  const { validate, loading: validating } = useValidateFlow();
  const { createBot, loading: saving } = useBotActions();

  // 线性数据结构
  const [trigger, setTrigger] = useState<TriggerNode | null>(null);
  const [conditions, setConditions] = useState<CondNode[]>([]);
  const [actions, setActions] = useState<ActionNode[]>([]);
  const [savingError, setSavingError] = useState<string | null>(null);

  // 配置弹窗状态
  const [editingItem, setEditingItem] = useState<{
    type: "trigger" | "condition" | "action";
    index?: number;
    node: TriggerNode | CondNode | ActionNode;
    label: string;
  } | null>(null);

  // 拖拽传感器
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
  );

  const {
    triggers,
    conditions: condMetas,
    actions: actionMetas,
  } = useMemo<NocodeMetadata>(() => {
    if (!metadata) return { triggers: [], conditions: [], actions: [] };
    return metadata;
  }, [metadata]);

  // 添加触发器（替换）
  const addTrigger = useCallback((nodeMeta: NodeMeta) => {
    const newTrigger: TriggerNode = {
      type: nodeMeta.type as TriggerType,
      params: createDefaultParams(nodeMeta),
    };
    setTrigger(newTrigger);
  }, []);

  // 添加条件（追加）
  const addCondition = useCallback((nodeMeta: NodeMeta) => {
    const newCond: CondNode = {
      type: nodeMeta.type as CondType,
      negate: false,
      params: createDefaultParams(nodeMeta),
    };
    setConditions((prev) => [...prev, newCond]);
  }, []);

  // 添加动作（追加）
  const addAction = useCallback((nodeMeta: NodeMeta) => {
    const newAction: ActionNode = {
      type: nodeMeta.type as ActionType,
      params: createDefaultParams(nodeMeta),
    };
    setActions((prev) => [...prev, newAction]);
  }, []);

  // 更新条件配置
  const updateCondition = useCallback(
    (index: number, params: Record<string, unknown>) => {
      setConditions((prev) =>
        prev.map((c, i) => (i === index ? { ...c, params } : c)),
      );
    },
    [],
  );

  // 更新动作配置
  const updateAction = useCallback(
    (index: number, params: Record<string, unknown>) => {
      setActions((prev) =>
        prev.map((a, i) => (i === index ? { ...a, params } : a)),
      );
    },
    [],
  );

  // 更新触发器配置
  const updateTriggerConfig = useCallback((params: Record<string, unknown>) => {
    setTrigger((prev) => (prev ? { ...prev, params } : prev));
  }, []);

  // 删除条件
  const deleteCondition = useCallback((index: number) => {
    setConditions((prev) => prev.filter((_, i) => i !== index));
  }, []);

  // 删除动作
  const deleteAction = useCallback((index: number) => {
    setActions((prev) => prev.filter((_, i) => i !== index));
  }, []);

  // 条件拖拽排序结束
  const handleConditionDragEnd = useCallback(
    (event: DragEndEvent) => {
      const { active, over } = event;
      if (over && active.id !== over.id) {
        const oldIndex = conditions.findIndex(
          (_, i) => String(i) === active.id,
        );
        const newIndex = conditions.findIndex((_, i) => String(i) === over.id);
        if (oldIndex !== -1 && newIndex !== -1) {
          const newConditions = [...conditions];
          const [moved] = newConditions.splice(oldIndex, 1);
          newConditions.splice(newIndex, 0, moved);
          setConditions(newConditions);
        }
      }
    },
    [conditions],
  );

  // 动作拖拽排序结束
  const handleActionDragEnd = useCallback(
    (event: DragEndEvent) => {
      const { active, over } = event;
      if (over && active.id !== over.id) {
        const oldIndex = actions.findIndex((_, i) => String(i) === active.id);
        const newIndex = actions.findIndex((_, i) => String(i) === over.id);
        if (oldIndex !== -1 && newIndex !== -1) {
          const newActions = [...actions];
          const [moved] = newActions.splice(oldIndex, 1);
          newActions.splice(newIndex, 0, moved);
          setActions(newActions);
        }
      }
    },
    [actions],
  );

  // 构建线性 Flow 对象
  const buildFlowPayload = useCallback((): Flow | null => {
    if (!trigger) return null;
    return {
      version: "1",
      trigger,
      conditions: conditions.length > 0 ? conditions : undefined,
      actions,
    };
  }, [trigger, conditions, actions]);

  // 校验流程
  const handleValidate = async () => {
    const flow = buildFlowPayload();
    if (!flow) {
      alert("请先配置触发器");
      return false;
    }
    const result = await validate(flow);
    if (!result.valid) {
      alert(`校验失败：\n${result.errors?.join("\n") || "未知错误"}`);
    } else {
      alert("流程校验通过！");
    }
    return result.valid;
  };

  // 保存机器人
  const handleSave = async () => {
    setSavingError(null);
    const flow = buildFlowPayload();
    if (!flow) {
      alert("请先配置触发器");
      return;
    }
    const validation = await validate(flow);
    if (!validation.valid) {
      alert(`保存前校验失败：\n${validation.errors?.join("\n")}`);
      return;
    }

    const requestData: CreateBotRequest = {
      name: "未命名零代码机器人",
      version: "1.0.0",
      description: "通过线性流程创建",
      type: "task",
      scriptCode: "",
      triggerType: "manual",
      configValues: { flow },
    };

    const id = await createBot(requestData);
    if (id) {
      alert(`机器人创建成功！ID: ${id}`);
      // 重置状态
      setTrigger(null);
      setConditions([]);
      setActions([]);
    } else {
      setSavingError("保存失败，请重试");
    }
  };

  if (metaLoading) {
    return <div className="p-4 text-center">加载节点定义中...</div>;
  }
  if (metaError) {
    return (
      <div className="p-4 text-center text-red-500">加载失败：{metaError}</div>
    );
  }

  return (
    <div className="flex flex-col h-[calc(100vh-24rem)] border rounded-lg overflow-hidden">
      <div className="flex flex-1 min-h-0">
        {/* 左侧节点库 */}
        <div className="w-64 bg-gray-100 p-3 overflow-y-auto border-r h-full">
          <CollapsibleSection title="触发器" defaultOpen>
            {triggers.map((t) => (
              <div
                key={t.type}
                className="p-2 bg-white rounded shadow cursor-pointer hover:bg-blue-50"
                onClick={() => addTrigger(t)}
              >
                {t.label}
              </div>
            ))}
          </CollapsibleSection>

          <CollapsibleSection title="条件" defaultOpen>
            {condMetas.map((c) => (
              <div
                key={c.type}
                className="p-2 bg-white rounded shadow cursor-pointer hover:bg-blue-50"
                onClick={() => addCondition(c)}
              >
                {c.label}
              </div>
            ))}
          </CollapsibleSection>

          <CollapsibleSection title="动作" defaultOpen>
            {actionMetas.map((a) => (
              <div
                key={a.type}
                className="p-2 bg-white rounded shadow cursor-pointer hover:bg-blue-50"
                onClick={() => addAction(a)}
              >
                {a.label}
              </div>
            ))}
          </CollapsibleSection>
        </div>

        {/* 右侧流程编辑区 */}
        <div className="flex-1 p-4 overflow-y-auto space-y-6">
          {/* 触发器区块 */}
          <div>
            <h2 className="font-bold text-lg mb-2">触发器</h2>
            {trigger ? (
              <div className="bg-gray-50 border rounded p-3 flex justify-between items-center">
                <span className="font-medium">{trigger.type}</span>
                <button
                  onClick={() =>
                    setEditingItem({
                      type: "trigger",
                      node: trigger,
                      label: trigger.type,
                    })
                  }
                  className="text-blue-600 hover:text-blue-800 text-sm"
                >
                  配置
                </button>
              </div>
            ) : (
              <div className="text-gray-400 italic">请从左侧点击添加触发器</div>
            )}
          </div>

          {/* 条件区块（可拖拽排序） */}
          <div>
            <h2 className="font-bold text-lg mb-2">条件（全部满足）</h2>
            {conditions.length === 0 ? (
              <div className="text-gray-400 italic">暂无条件，点击左侧添加</div>
            ) : (
              <DndContext
                sensors={sensors}
                collisionDetection={closestCenter}
                onDragEnd={handleConditionDragEnd}
              >
                <SortableContext
                  items={conditions.map((_, i) => String(i))}
                  strategy={verticalListSortingStrategy}
                >
                  <div className="space-y-2">
                    {conditions.map((cond, idx) => (
                      <SortableItem
                        key={idx}
                        id={String(idx)}
                        typeLabel={cond.type}
                        onEdit={() =>
                          setEditingItem({
                            type: "condition",
                            index: idx,
                            node: cond,
                            label: cond.type,
                          })
                        }
                        onDelete={() => deleteCondition(idx)}
                      />
                    ))}
                  </div>
                </SortableContext>
              </DndContext>
            )}
          </div>

          {/* 动作区块（可拖拽排序） */}
          <div>
            <h2 className="font-bold text-lg mb-2">动作（顺序执行）</h2>
            {actions.length === 0 ? (
              <div className="text-gray-400 italic">暂无动作，点击左侧添加</div>
            ) : (
              <DndContext
                sensors={sensors}
                collisionDetection={closestCenter}
                onDragEnd={handleActionDragEnd}
              >
                <SortableContext
                  items={actions.map((_, i) => String(i))}
                  strategy={verticalListSortingStrategy}
                >
                  <div className="space-y-2">
                    {actions.map((action, idx) => (
                      <SortableItem
                        key={idx}
                        id={String(idx)}
                        typeLabel={action.type}
                        onEdit={() =>
                          setEditingItem({
                            type: "action",
                            index: idx,
                            node: action,
                            label: action.type,
                          })
                        }
                        onDelete={() => deleteAction(idx)}
                      />
                    ))}
                  </div>
                </SortableContext>
              </DndContext>
            )}
          </div>
        </div>
      </div>

      {/* 底部按钮 */}
      <div className="p-3 border-t flex justify-end gap-2">
        <button
          className="px-3 py-1 border rounded text-sm hover:bg-gray-50"
          onClick={handleValidate}
          disabled={validating}
        >
          {validating ? "校验中..." : "校验流程"}
        </button>
        <button
          className="px-3 py-1 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 disabled:bg-gray-400"
          onClick={handleSave}
          disabled={saving}
        >
          {saving ? "保存中..." : "保存机器人"}
        </button>
      </div>

      {savingError && (
        <div className="text-red-500 text-center p-2">{savingError}</div>
      )}

      {/* 配置弹窗 */}
      {editingItem && (
        <ConfigModal
          title={editingItem.label}
          params={
            editingItem.type === "trigger"
              ? (editingItem.node as TriggerNode).params || {}
              : (editingItem.node as CondNode | ActionNode).params
          }
          onSave={(newParams) => {
            if (editingItem.type === "trigger") {
              updateTriggerConfig(newParams);
            } else if (
              editingItem.type === "condition" &&
              editingItem.index !== undefined
            ) {
              updateCondition(editingItem.index, newParams);
            } else if (
              editingItem.type === "action" &&
              editingItem.index !== undefined
            ) {
              updateAction(editingItem.index, newParams);
            }
            setEditingItem(null);
          }}
          onClose={() => setEditingItem(null)}
        />
      )}
    </div>
  );
}
