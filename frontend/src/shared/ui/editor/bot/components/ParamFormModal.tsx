"use client";

import { useState } from "react";
import { ParamMeta } from "@/shared/api/types/nocode.model";

interface ParamFormModalProps {
  title: string;
  params: Record<string, unknown>;
  paramMetas: ParamMeta[];
  onSave: (newParams: Record<string, unknown>) => void;
  onClose: () => void;
}

function resolveDefault(meta: ParamMeta): unknown {
  if (meta.default !== undefined && meta.default !== null) return meta.default;
  switch (meta.type) {
    case "number":
      return 0;
    case "boolean":
      return false;
    case "tags":
      return [];
    default:
      return "";
  }
}

export function ParamFormModal({
  title,
  params,
  paramMetas,
  onSave,
  onClose,
}: ParamFormModalProps) {
  const [formValues, setFormValues] = useState<Record<string, unknown>>(() => {
    const initial: Record<string, unknown> = {};
    paramMetas.forEach((meta) => {
      initial[meta.key] =
        params[meta.key] !== undefined ? params[meta.key] : resolveDefault(meta);
    });
    return initial;
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const setField = (key: string, value: unknown) => {
    setFormValues((prev) => ({ ...prev, [key]: value }));
    if (errors[key]) {
      setErrors((prev) => {
        const next = { ...prev };
        delete next[key];
        return next;
      });
    }
  };

  const handleSave = () => {
    const newErrors: Record<string, string> = {};
    paramMetas.forEach((meta) => {
      if (meta.required) {
        const val = formValues[meta.key];
        if (val === undefined || val === null || val === "") {
          newErrors[meta.key] = `${meta.label} 为必填项`;
        }
        if (meta.type === "tags" && Array.isArray(val) && val.length === 0) {
          newErrors[meta.key] = `${meta.label} 至少需要一个`;
        }
      }
    });
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }
    onSave(formValues);
  };

  if (paramMetas.length === 0) {
    return (
      <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
        <div className="bg-white rounded-xl shadow-2xl w-96 p-6">
          <h3 className="font-bold text-lg mb-4">{title}</h3>
          <p className="text-gray-500 text-sm mb-4">该节点无需额外配置</p>
          <div className="flex justify-end gap-2">
            <button
              className="px-4 py-2 border rounded-lg hover:bg-gray-100 text-sm"
              onClick={onClose}
            >
              取消
            </button>
            <button
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm"
              onClick={() => onSave({})}
            >
              确认
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <div className="bg-white rounded-xl shadow-2xl w-[480px] max-h-[80vh] flex flex-col">
        <div className="p-4 border-b">
          <h3 className="font-bold text-lg">配置 - {title}</h3>
        </div>

        <div className="p-4 space-y-4 overflow-y-auto flex-1">
          {paramMetas.map((meta) => (
            <div key={meta.key}>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                {meta.label}
                {meta.required && (
                  <span className="text-red-500 ml-0.5">*</span>
                )}
              </label>

              {meta.type === "textarea" ? (
                <textarea
                  className="w-full border rounded-lg p-2 text-sm font-mono
                    focus:outline-none focus:ring-2 focus:ring-blue-500
                    disabled:bg-gray-100"
                  rows={3}
                  value={String(formValues[meta.key] ?? "")}
                  onChange={(e) => setField(meta.key, e.target.value)}
                  placeholder={meta.placeholder}
                />
              ) : meta.type === "select" ? (
                <select
                  className="w-full border rounded-lg p-2 text-sm
                    focus:outline-none focus:ring-2 focus:ring-blue-500
                    disabled:bg-gray-100"
                  value={String(formValues[meta.key] ?? "")}
                  onChange={(e) => setField(meta.key, e.target.value)}
                >
                  <option value="">请选择</option>
                  {meta.options?.map((opt) => (
                    <option key={String(opt.value)} value={String(opt.value)}>
                      {opt.label}
                    </option>
                  ))}
                </select>
              ) : meta.type === "number" ? (
                <input
                  type="number"
                  className="w-full border rounded-lg p-2 text-sm
                    focus:outline-none focus:ring-2 focus:ring-blue-500
                    disabled:bg-gray-100"
                  value={Number(formValues[meta.key] ?? 0)}
                  onChange={(e) =>
                    setField(meta.key, Number(e.target.value))
                  }
                  placeholder={meta.placeholder}
                />
              ) : meta.type === "boolean" ? (
                <label className="inline-flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-sm"
                    checked={Boolean(formValues[meta.key])}
                    onChange={(e) => setField(meta.key, e.target.checked)}
                  />
                  <span className="text-sm text-gray-600">
                    {meta.placeholder || "启用"}
                  </span>
                </label>
              ) : meta.type === "tags" ? (
                <TagsInput
                  value={Array.isArray(formValues[meta.key])
                    ? (formValues[meta.key] as string[])
                    : []}
                  onChange={(tags) => setField(meta.key, tags)}
                  placeholder={meta.placeholder}
                />
              ) : (
                <input
                  type="text"
                  className="w-full border rounded-lg p-2 text-sm
                    focus:outline-none focus:ring-2 focus:ring-blue-500
                    disabled:bg-gray-100"
                  value={String(formValues[meta.key] ?? "")}
                  onChange={(e) => setField(meta.key, e.target.value)}
                  placeholder={meta.placeholder}
                />
              )}

              {errors[meta.key] && (
                <p className="text-red-500 text-xs mt-1">{errors[meta.key]}</p>
              )}
            </div>
          ))}
        </div>

        <div className="p-4 border-t flex justify-end gap-2">
          <button
            className="px-4 py-2 border rounded-lg hover:bg-gray-100 text-sm"
            onClick={onClose}
          >
            取消
          </button>
          <button
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm"
            onClick={handleSave}
          >
            保存
          </button>
        </div>
      </div>
    </div>
  );
}

// ---------- 标签输入组件 ----------
function TagsInput({
  value,
  onChange,
  placeholder,
}: {
  value: string[];
  onChange: (tags: string[]) => void;
  placeholder?: string;
}) {
  const [input, setInput] = useState("");

  const addTag = () => {
    const trimmed = input.trim();
    if (trimmed && !value.includes(trimmed)) {
      onChange([...value, trimmed]);
    }
    setInput("");
  };

  const removeTag = (tag: string) => {
    onChange(value.filter((t) => t !== tag));
  };

  return (
    <div>
      <div className="flex flex-wrap gap-1 mb-1">
        {value.map((tag) => (
          <span
            key={tag}
            className="inline-flex items-center gap-1 px-2 py-0.5 bg-blue-100 text-blue-800 text-xs rounded"
          >
            {tag}
            <button
              type="button"
              className="hover:text-red-600"
              onClick={() => removeTag(tag)}
            >
              ×
            </button>
          </span>
        ))}
      </div>
      <div className="flex gap-1">
        <input
          type="text"
          className="flex-1 border rounded-lg p-1.5 text-sm
            focus:outline-none focus:ring-2 focus:ring-blue-500"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              addTag();
            }
          }}
          placeholder={placeholder || "输入后按回车添加"}
        />
        <button
          type="button"
          className="px-3 py-1.5 border rounded-lg text-sm hover:bg-gray-100"
          onClick={addTag}
        >
          添加
        </button>
      </div>
    </div>
  );
}
