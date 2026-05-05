"use client";

import { useState } from "react";
import {
  Globe,
  Tag,
  User,
  FileCode2,
  AlignLeft,
  Hash,
  X,
  Plus,
} from "lucide-react";
import { PluginMeta, SLOT_NAMES } from "@/shared/type/plugin.type";
import { CreatePluginPayload } from "@/shared/api/modules/plugin/plugins";

interface PluginUploadFormProps {
  initial?: PluginMeta;
  onSubmit: (payload: CreatePluginPayload) => Promise<void>;
  onCancel: () => void;
  isLoading: boolean;
}

const EMPTY_FORM: CreatePluginPayload = {
  name: "",
  version: "1.0.0",
  description: "",
  author: "",
  scriptUrl: "",
  enabled: true,
  slots: [],
};

export function PluginUploadForm({
  initial,
  onSubmit,
  onCancel,
  isLoading,
}: PluginUploadFormProps) {
  const [form, setForm] = useState<CreatePluginPayload>(
    initial
      ? {
          name: initial.name,
          version: initial.version,
          description: initial.description,
          author: initial.author,
          scriptUrl: initial.scriptUrl,
          enabled: initial.enabled,
          slots: initial.slots ?? [],
        }
      : EMPTY_FORM,
  );
  const [errors, setErrors] = useState<
    Partial<Record<keyof CreatePluginPayload, string>>
  >({});

  const set = (key: keyof CreatePluginPayload, value: unknown) => {
    setForm((prev) => ({ ...prev, [key]: value }));
    if (errors[key]) setErrors((prev) => ({ ...prev, [key]: undefined }));
  };

  const toggleSlot = (slot: string) => {
    const slots = form.slots ?? [];
    set(
      "slots",
      slots.includes(slot) ? slots.filter((s) => s !== slot) : [...slots, slot],
    );
  };

  const validate = (): boolean => {
    const newErrors: typeof errors = {};
    if (!form.name.trim()) newErrors.name = "Name is required";
    if (!form.version.trim()) newErrors.version = "Version is required";
    if (!form.author.trim()) newErrors.author = "Author is required";
    if (!form.scriptUrl.trim()) {
      newErrors.scriptUrl = "Script URL is required";
    } else {
      try {
        new URL(form.scriptUrl);
      } catch {
        newErrors.scriptUrl = "Must be a valid URL (https://...)";
      }
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    await onSubmit(form);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Name */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-medium">Plugin Name *</span>
        </label>
        <label className="input input-bordered flex items-center gap-2">
          <Tag className="w-4 h-4 text-base-content/40 shrink-0" />
          <input
            type="text"
            placeholder="My Awesome Plugin"
            value={form.name}
            onChange={(e) => set("name", e.target.value)}
            className="grow"
          />
        </label>
        {errors.name && (
          <span className="label-text-alt text-error mt-1">{errors.name}</span>
        )}
      </div>

      {/* Version + Author row */}
      <div className="grid grid-cols-2 gap-3">
        <div className="form-control">
          <label className="label">
            <span className="label-text font-medium">Version *</span>
          </label>
          <label className="input input-bordered flex items-center gap-2">
            <Hash className="w-4 h-4 text-base-content/40 shrink-0" />
            <input
              type="text"
              placeholder="1.0.0"
              value={form.version}
              onChange={(e) => set("version", e.target.value)}
              className="grow"
            />
          </label>
          {errors.version && (
            <span className="label-text-alt text-error mt-1">
              {errors.version}
            </span>
          )}
        </div>
        <div className="form-control">
          <label className="label">
            <span className="label-text font-medium">Author *</span>
          </label>
          <label className="input input-bordered flex items-center gap-2">
            <User className="w-4 h-4 text-base-content/40 shrink-0" />
            <input
              type="text"
              placeholder="Your name"
              value={form.author}
              onChange={(e) => set("author", e.target.value)}
              className="grow"
            />
          </label>
          {errors.author && (
            <span className="label-text-alt text-error mt-1">
              {errors.author}
            </span>
          )}
        </div>
      </div>

      {/* Script URL */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-medium">Script URL *</span>
          <span className="label-text-alt text-base-content/40">
            Publicly accessible JS bundle
          </span>
        </label>
        <label className="input input-bordered flex items-center gap-2">
          <Globe className="w-4 h-4 text-base-content/40 shrink-0" />
          <input
            type="text"
            placeholder="https://cdn.example.com/my-plugin.js"
            value={form.scriptUrl}
            onChange={(e) => set("scriptUrl", e.target.value)}
            className="grow font-mono text-sm"
          />
        </label>
        {errors.scriptUrl && (
          <span className="label-text-alt text-error mt-1">
            {errors.scriptUrl}
          </span>
        )}
      </div>

      {/* Description */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-medium">Description</span>
        </label>
        <label className="flex items-start gap-2 textarea textarea-bordered">
          <AlignLeft className="w-4 h-4 text-base-content/40 shrink-0 mt-1" />
          <textarea
            placeholder="What does this plugin do?"
            value={form.description}
            onChange={(e) => set("description", e.target.value)}
            rows={3}
            className="grow resize-none bg-transparent outline-none"
          />
        </label>
      </div>

      {/* Slots */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-medium flex items-center gap-1">
            <FileCode2 className="w-4 h-4" /> Injection Slots
          </span>
          <span className="label-text-alt text-base-content/40">
            Declare which slots this plugin uses
          </span>
        </label>
        <div className="flex flex-wrap gap-2 p-3 bg-base-200 rounded-lg border border-base-300">
          {SLOT_NAMES.map((slot) => {
            const selected = (form.slots ?? []).includes(slot);
            return (
              <button
                key={slot}
                type="button"
                onClick={() => toggleSlot(slot)}
                className={`badge gap-1 cursor-pointer transition-all ${
                  selected ? "badge-primary" : "badge-outline hover:badge-ghost"
                }`}
              >
                {selected ? (
                  <X className="w-3 h-3" />
                ) : (
                  <Plus className="w-3 h-3" />
                )}
                {slot}
              </button>
            );
          })}
        </div>
      </div>

      {/* Enabled toggle */}
      <div className="form-control">
        <label className="label cursor-pointer justify-start gap-3">
          <input
            type="checkbox"
            className="toggle toggle-primary"
            checked={form.enabled}
            onChange={(e) => set("enabled", e.target.checked)}
          />
          <span className="label-text font-medium">
            Enable immediately after install
          </span>
        </label>
      </div>

      {/* Actions */}
      <div className="flex justify-end gap-2 pt-2 border-t border-base-300">
        <button
          type="button"
          onClick={onCancel}
          className="btn btn-ghost btn-sm"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={isLoading}
          className="btn btn-primary btn-sm gap-2"
        >
          {isLoading && <span className="loading loading-spinner loading-xs" />}
          {initial ? "Save Changes" : "Install Plugin"}
        </button>
      </div>
    </form>
  );
}
