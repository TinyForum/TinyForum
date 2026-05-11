// components/BotFormModal.tsx

import {
  BotConfigField,
  BotPricingModel,
  BotTriggerType,
  BotType,
  BotVO,
  CreateBotRequest,
  UpdateBotRequest,
} from "@/shared/api/types/bot.model";
import React, { useState, useEffect } from "react";
import { toast } from "react-hot-toast";

// 动态配置表单组件
function DynamicConfigForm({
  schema,
  values,
  onChange,
}: {
  schema: BotConfigField[];
  values: Record<string, any>;
  onChange: (key: string, value: any) => void;
}) {
  if (!schema?.length) return null;
  return (
    <div className="space-y-3">
      {schema.map((field) => (
        <div key={field.key} className="form-control">
          <label className="label">
            <span className="label-text">
              {field.label}
              {field.required && <span className="text-error ml-1">*</span>}
            </span>
            {field.description && (
              <span className="label-text-alt text-base-content/60">
                {field.description}
              </span>
            )}
          </label>
          {field.type === "textarea" ? (
            <textarea
              className="textarea textarea-bordered"
              placeholder={field.placeholder}
              value={values[field.key] ?? field.defaultValue ?? ""}
              onChange={(e) => onChange(field.key, e.target.value)}
            />
          ) : field.type === "select" ? (
            <select
              className="select select-bordered"
              value={values[field.key] ?? field.defaultValue ?? ""}
              onChange={(e) => onChange(field.key, e.target.value)}
            >
              <option value="">请选择</option>
              {field.options?.map((opt) => (
                <option key={String(opt.value)} value={String(opt.value)}>
                  {opt.label}
                </option>
              ))}
            </select>
          ) : field.type === "boolean" ? (
            <input
              type="checkbox"
              className="toggle"
              checked={!!(values[field.key] ?? field.defaultValue ?? false)}
              onChange={(e) => onChange(field.key, e.target.checked)}
            />
          ) : (
            <input
              type={field.type === "secret" ? "password" : field.type}
              className="input input-bordered"
              placeholder={field.placeholder}
              value={values[field.key] ?? field.defaultValue ?? ""}
              onChange={(e) => onChange(field.key, e.target.value)}
            />
          )}
        </div>
      ))}
    </div>
  );
}

// 获取预设 Cron 表达式
const getCronFromSchedule = (schedule: string): string => {
  switch (schedule) {
    case "hourly":
      return "0 * * * *";
    case "daily":
      return "0 0 * * *";
    case "weekly":
      return "0 0 * * 0";
    default:
      return "0 */6 * * *";
  }
};

interface BotFormModalProps {
  isOpen: boolean;
  editingBot: BotVO | null;
  onClose: () => void;
  onSave: (
    data: CreateBotRequest | UpdateBotRequest,
    isEdit: boolean,
    botId?: number,
  ) => Promise<void>;
  isLoading?: boolean;
}

export function BotFormModal({
  isOpen,
  editingBot,
  onClose,
  onSave,
  isLoading,
}: BotFormModalProps) {
  const [formData, setFormData] = useState<Partial<CreateBotRequest>>({
    name: "",
    version: "1.0.0",
    type: "task",
    triggerType: "schedule",
    timeoutSec: 10,
    retryTimes: 0,
    pricing: { model: "free" },
    permissions: [],
    configSchema: [],
    configValues: {},
  });
  const [schedulePreset, setSchedulePreset] = useState("");

  useEffect(() => {
    if (editingBot) {
      // 编辑模式：预填充
      setFormData({
        name: editingBot.name,
        version: editingBot.version,
        description: editingBot.description,
        summary: editingBot.summary,
        avatarUrl: editingBot.avatarUrl,
        scriptCode: editingBot.scriptCode,
        scriptUrl: editingBot.scriptUrl,
        type: editingBot.type,
        triggerType: editingBot.triggerType,
        cronExpr: editingBot.cronExpr,
        eventFilter: editingBot.eventFilter,
        timeoutSec: editingBot.timeoutSec,
        retryTimes: editingBot.retryTimes,
        envVars: editingBot.envVars,
        resourceLimit: editingBot.resourceLimit,
        pricing: editingBot.pricing,
        permissions: editingBot.permissions,
        configSchema: editingBot.configSchema,
        configValues: editingBot.configValues,
        enabled: editingBot.enabled, // 用于编辑时显示
      });
      setSchedulePreset("");
    } else {
      // 创建模式：重置
      setFormData({
        name: "",
        version: "1.0.0",
        type: "task",
        triggerType: "schedule",
        timeoutSec: 10,
        retryTimes: 0,
        pricing: { model: "free" },
        permissions: [],
        configSchema: [],
        configValues: {},
      });
      setSchedulePreset("");
    }
  }, [editingBot, isOpen]);

  const handleSubmit = async () => {
    if (!formData.name || !formData.scriptCode) {
      toast.error("请填写名称和脚本代码");
      return;
    }
    if (
      formData.triggerType === "schedule" &&
      !formData.cronExpr &&
      !schedulePreset
    ) {
      toast.error("定时任务必须填写 Cron 表达式或选择一个预设频率");
      return;
    }

    let cron = formData.cronExpr;
    if (schedulePreset && !cron) {
      cron = getCronFromSchedule(schedulePreset);
    }

    if (editingBot) {
      // 更新机器人
      const updateData: UpdateBotRequest = {
        name: formData.name,
        version: formData.version,
        description: formData.description,
        summary: formData.summary,
        avatarUrl: formData.avatarUrl,
        scriptCode: formData.scriptCode,
        scriptUrl: formData.scriptUrl,
        type: formData.type,
        triggerType: formData.triggerType,
        cronExpr: cron,
        eventFilter: formData.eventFilter,
        timeoutSec: formData.timeoutSec,
        retryTimes: formData.retryTimes,
        envVars: formData.envVars,
        resourceLimit: formData.resourceLimit,
        pricing: formData.pricing,
        permissions: formData.permissions,
        configSchema: formData.configSchema,
        configValues: formData.configValues,
        enabled: formData.enabled,
      };
      await onSave(updateData, true, editingBot.id);
    } else {
      // 创建机器人
      const createData: CreateBotRequest = {
        name: formData.name!,
        version: formData.version || "1.0.0",
        description: formData.description,
        summary: formData.summary,
        avatarUrl: formData.avatarUrl,
        type: formData.type!,
        scriptCode: formData.scriptCode!,
        scriptUrl: formData.scriptUrl,
        triggerType: formData.triggerType!,
        cronExpr: cron,
        eventFilter: formData.eventFilter,
        timeoutSec: formData.timeoutSec,
        retryTimes: formData.retryTimes,
        envVars: formData.envVars,
        resourceLimit: formData.resourceLimit,
        pricing: formData.pricing,
        permissions: formData.permissions,
        configSchema: formData.configSchema,
        configValues: formData.configValues,
      };
      await onSave(createData, false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-2xl">
        <h3 className="font-bold text-lg">
          {editingBot ? "编辑机器人" : "创建机器人"}
        </h3>
        <div className="py-4 space-y-4">
          {/* 基本信息 */}
          <div className="form-control">
            <label className="label">
              <span className="label-text">名称 *</span>
            </label>
            <input
              type="text"
              className="input input-bordered"
              value={formData.name || ""}
              onChange={(e) =>
                setFormData({ ...formData, name: e.target.value })
              }
            />
          </div>
          <div className="form-control">
            <label className="label">
              <span className="label-text">描述</span>
            </label>
            <textarea
              className="textarea textarea-bordered"
              rows={2}
              value={formData.description || ""}
              onChange={(e) =>
                setFormData({ ...formData, description: e.target.value })
              }
            />
          </div>
          <div className="form-control">
            <label className="label">
              <span className="label-text">类型 *</span>
            </label>
            <select
              className="select select-bordered"
              value={formData.type}
              onChange={(e) =>
                setFormData({ ...formData, type: e.target.value as BotType })
              }
            >
              <option value="task">任务 (task)</option>
              <option value="chat">聊天 (chat)</option>
              <option value="moderate">审核 (moderate)</option>
              <option value="notify">通知 (notify)</option>
              <option value="sync">同步 (sync)</option>
              <option value="webhook">Webhook (webhook)</option>
              <option value="analysis">分析 (analysis)</option>
            </select>
          </div>
          {/* Lua 脚本 */}
          <div className="form-control">
            <label className="label">
              <span className="label-text">Lua 脚本 *</span>
            </label>
            <textarea
              className="textarea textarea-bordered font-mono"
              rows={8}
              placeholder="function main()&#10;    log('Hello')&#10;end"
              value={formData.scriptCode || ""}
              onChange={(e) =>
                setFormData({ ...formData, scriptCode: e.target.value })
              }
            />
          </div>
          {/* 触发方式 */}
          <div className="form-control">
            <label className="label">
              <span className="label-text">触发方式</span>
            </label>
            <select
              className="select select-bordered"
              value={formData.triggerType}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  triggerType: e.target.value as BotTriggerType,
                })
              }
            >
              <option value="schedule">定时执行</option>
              <option value="event">事件触发</option>
              <option value="webhook">Webhook</option>
              <option value="manual">手动触发</option>
            </select>
          </div>
          {/* 定时配置 */}
          {formData.triggerType === "schedule" && (
            <>
              <div className="form-control">
                <label className="label">
                  <span className="label-text">预设频率</span>
                </label>
                <select
                  className="select select-bordered"
                  value={schedulePreset}
                  onChange={(e) => setSchedulePreset(e.target.value)}
                >
                  <option value="">自定义 Cron</option>
                  <option value="hourly">每小时</option>
                  <option value="daily">每天</option>
                  <option value="weekly">每周</option>
                </select>
              </div>
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Cron 表达式</span>
                  <span className="label-text-alt">(使用预设可留空)</span>
                </label>
                <input
                  type="text"
                  className="input input-bordered"
                  placeholder="0 */6 * * *"
                  value={formData.cronExpr || ""}
                  onChange={(e) =>
                    setFormData({ ...formData, cronExpr: e.target.value })
                  }
                />
              </div>
            </>
          )}
          {/* 事件过滤 */}
          {formData.triggerType === "event" && (
            <div className="form-control">
              <label className="label">
                <span className="label-text">事件过滤</span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                placeholder="post:create, comment:create"
                value={formData.eventFilter || ""}
                onChange={(e) =>
                  setFormData({ ...formData, eventFilter: e.target.value })
                }
              />
            </div>
          )}
          {/* 超时与重试 */}
          <div className="grid grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">超时(秒)</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={formData.timeoutSec ?? 10}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    timeoutSec: parseInt(e.target.value) || 10,
                  })
                }
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">重试次数</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={formData.retryTimes ?? 0}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    retryTimes: parseInt(e.target.value) || 0,
                  })
                }
              />
            </div>
          </div>
          {/* 动态配置 */}
          {formData.configSchema && formData.configSchema.length > 0 && (
            <div>
              <label className="label">
                <span className="label-text">机器人配置</span>
              </label>
              <DynamicConfigForm
                schema={formData.configSchema}
                values={formData.configValues || {}}
                onChange={(key, value) =>
                  setFormData({
                    ...formData,
                    configValues: {
                      ...(formData.configValues || {}),
                      [key]: value,
                    },
                  })
                }
              />
            </div>
          )}
          {/* 定价模型 */}
          <div className="form-control">
            <label className="label">
              <span className="label-text">定价模型</span>
            </label>
            <select
              className="select select-bordered"
              value={formData.pricing?.model ?? "free"}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  pricing: {
                    ...formData.pricing,
                    model: e.target.value as BotPricingModel,
                  },
                })
              }
            >
              <option value="free">免费</option>
              <option value="freemium">免费增值</option>
              <option value="paid">付费</option>
              <option value="subscription">订阅</option>
            </select>
          </div>
          {/* 编辑模式下显示启用开关 */}
          {editingBot && (
            <div className="form-control">
              <label className="label cursor-pointer">
                <span className="label-text">启用机器人</span>
                <input
                  type="checkbox"
                  className="toggle"
                  checked={!!formData.enabled}
                  onChange={(e) =>
                    setFormData({ ...formData, enabled: e.target.checked })
                  }
                />
              </label>
            </div>
          )}
        </div>
        <div className="modal-action">
          <button className="btn" onClick={onClose}>
            取消
          </button>
          <button
            className="btn btn-primary"
            onClick={handleSubmit}
            disabled={isLoading}
          >
            {isLoading ? "保存中..." : "保存"}
          </button>
        </div>
      </div>
    </div>
  );
}
