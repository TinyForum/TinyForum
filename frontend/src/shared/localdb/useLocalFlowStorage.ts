// hooks/useLocalFlowStorage.ts
import { useState } from "react";
import Dexie, { Table } from "dexie";
import { Node, Edge } from "reactflow";

// 定义数据库结构
class FlowDraftDB extends Dexie {
  drafts!: Table<{
    id: string;
    nodes: Node[];
    edges: Edge[];
    updatedAt: number;
  }>;

  constructor() {
    super("NocoFlowDB");
    this.version(1).stores({
      drafts: "id, updatedAt",
    });
  }
}

const db = new FlowDraftDB();
const DRAFT_ID = "current_draft";

export function useLocalFlowStorage() {
  const [saving, setSaving] = useState(false);

  // 加载草稿
  const loadDraft = async (): Promise<{
    nodes: Node[];
    edges: Edge[];
  } | null> => {
    try {
      const draft = await db.drafts.get(DRAFT_ID);
      if (draft) {
        return { nodes: draft.nodes, edges: draft.edges };
      }
      return null;
    } catch (error) {
      console.error("加载草稿失败:", error);
      return null;
    }
  };

  // 保存草稿（去抖可自行添加）
  const saveDraft = async (nodes: Node[], edges: Edge[]) => {
    setSaving(true);
    try {
      await db.drafts.put({
        id: DRAFT_ID,
        nodes,
        edges,
        updatedAt: Date.now(),
      });
    } catch (error) {
      console.error("保存草稿失败:", error);
    } finally {
      setSaving(false);
    }
  };

  // 清空草稿
  const clearDraft = async () => {
    await db.drafts.delete(DRAFT_ID);
  };

  return { loadDraft, saveDraft, clearDraft, saving };
}
