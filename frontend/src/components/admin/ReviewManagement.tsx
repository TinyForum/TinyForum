// components/admin/PostModeration.tsx
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import toast from "react-hot-toast";
import { CheckCircle, XCircle, Clock, Eye, User as UserIcon, FileText, AlertTriangle } from "lucide-react";
import { adminApi } from "@/lib/api";
import type { Post, User as ApiUser } from "@/lib/api/types";
import DOMPurify from "dompurify";
// 扩展 Post 类型以包含风险信息（根据实际后端字段调整）
interface PostWithRisk extends Post {
  risk_score?: number;
  risk_reason?: string;
  risk_logs?: Array<{
    id: number;
    level: string;
    rule: string;
    matched_content: string;
    created_at: string;
  }>;
}

// 定义后端返回的分页数据结构
interface PendingPostsResponse {
  list: PostWithRisk[];
  total: number;
}

// 获取待审核帖子列表
const fetchPendingPosts = async (params: { page: number; page_size: number; keyword?: string }): Promise<PendingPostsResponse> => {
  const res = await adminApi.listPendingPosts(params);
  return res.data.data; // 假设后端返回 { list: Post[], total: number }
};

// 审核帖子
const reviewPost = async (id: number, status: "approved" | "rejected" | "pending") => {
  const res = await adminApi.reviewPosts(id, { status });
  return res.data;
};

export function ReviewManagement() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [keyword, setKeyword] = useState("");
  const [selectedPost, setSelectedPost] = useState<PostWithRisk | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  // 查询待审核列表
  const { data, isLoading, refetch } = useQuery<PendingPostsResponse>({
    queryKey: ["admin-pending-posts", page, keyword],
    queryFn: () => fetchPendingPosts({ page, page_size: 20, keyword }),
    placeholderData: (prev) => prev, // 相当于 v4 的 keepPreviousData: true
  });

  // 审核 mutation
  const reviewMutation = useMutation({
    mutationFn: ({ id, status }: { id: number; status: "approved" | "rejected" | "pending" }) =>
      reviewPost(id, status),
    onSuccess: (_, variables) => {
      toast.success(`帖子已${variables.status === "approved" ? "通过" : variables.status === "rejected" ? "拒绝" : "转为待审"}`);
      queryClient.invalidateQueries({ queryKey: ["admin-pending-posts"] });
      if (selectedPost && selectedPost.id === variables.id) {
        setIsModalOpen(false);
        setSelectedPost(null);
      }
    },
    onError: () => toast.error("操作失败"),
  });

  const handleOpenModal = (post: PostWithRisk) => {
    setSelectedPost(post);
    setIsModalOpen(true);
  };

  const handleReview = (status: "approved" | "rejected" | "pending") => {
    if (!selectedPost) return;
    reviewMutation.mutate({ id: selectedPost.id, status });
  };

  const posts = data?.list ?? [];
  const total = data?.total ?? 0;
  const totalPages = Math.ceil(total / 20);
  
  const sanitizeHtml = (html: string) => ({
  __html: DOMPurify.sanitize(html, {
    ALLOWED_TAGS: ["b", "i", "em", "strong", "a", "p", "br", "ul", "ol", "li", "img", "h1", "h2", "h3", "h4", "blockquote", "code", "pre"],
    ALLOWED_ATTR: ["href", "target", "src", "alt", "class", "id"],
  }),
});

  return (
    <div className="space-y-4">
      {/* 搜索栏 */}
      <div className="flex justify-between items-center gap-4">
        <div className="flex gap-2 flex-1">
          <input
            type="text"
            placeholder="搜索标题或作者"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            className="input input-bordered flex-1"
          />
          <button className="btn btn-outline" onClick={() => setPage(1)}>搜索</button>
          <button className="btn btn-ghost" onClick={() => refetch()}>刷新</button>
        </div>
      </div>

      {/* 列表 */}
      {isLoading ? (
        <div className="flex justify-center py-8"><span className="loading loading-spinner loading-lg"></span></div>
      ) : (
        <div className="space-y-3">
          {posts.map((post) => (
            <div key={post.id} className="card bg-base-100 shadow-sm border border-base-200 cursor-pointer hover:shadow-md transition-shadow" onClick={() => handleOpenModal(post)}>
              <div className="card-body p-4">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="font-medium">{post.author?.username || `用户${post.author_id}`}</span>
                      <span className="text-xs text-gray-400">#{post.author_id}</span>
                      {post.risk_score && (
                        <span className="badge badge-warning badge-sm">
                          <AlertTriangle className="w-3 h-3 mr-1" />
                          风险分 {post.risk_score}
                        </span>
                      )}
                    </div>
                    <h3 className="font-semibold">{post.title}</h3>
                    {post.risk_reason && (
                      <p className="text-sm text-orange-600 mt-1">⚠️ {post.risk_reason}</p>
                    )}
                  </div>
                  <Eye className="w-5 h-5 text-gray-400" />
                </div>
              </div>
            </div>
          ))}
          {posts.length === 0 && (
            <div className="text-center py-8 text-gray-500">暂无待审核帖子</div>
          )}
        </div>
      )}

      {/* 分页 */}
      {totalPages > 1 && (
        <div className="flex justify-center gap-2 pt-4">
          <button className="btn btn-sm" disabled={page === 1} onClick={() => setPage(p => p - 1)}>上一页</button>
          <span className="btn btn-sm btn-ghost">第 {page} / {totalPages} 页</span>
          <button className="btn btn-sm" disabled={page === totalPages} onClick={() => setPage(p => p + 1)}>下一页</button>
        </div>
      )}

      {/* 模态框 - 详情 */}
      {isModalOpen && selectedPost && (
        <dialog className="modal modal-open" onClick={() => setIsModalOpen(false)}>
          <div className="modal-box max-w-3xl" onClick={(e) => e.stopPropagation()}>
            <h3 className="font-bold text-lg mb-4">帖子审核</h3>
            
            {/* 用户信息 */}
            <div className="mb-4 p-3 bg-base-200 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <UserIcon className="w-4 h-4" />
                <span className="font-medium">发布人信息</span>
              </div>
              <div className="grid grid-cols-2 gap-2 text-sm">
                <div>用户名：{selectedPost.author?.username || `ID ${selectedPost.author_id}`}</div>
                <div>用户ID：{selectedPost.author_id}</div>
                {selectedPost.author?.email && <div>邮箱：{selectedPost.author.email}</div>}
                {selectedPost.author?.role && <div>角色：{selectedPost.author.role}</div>}
                {selectedPost.author?.created_at && <div>注册时间：{new Date(selectedPost.author.created_at).toLocaleString()}</div>}
              </div>
            </div>

            {/* 风险日志 */}
            <div className="mb-4 p-3 bg-base-200 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <AlertTriangle className="w-4 h-4 text-warning" />
                <span className="font-medium">风险检测日志</span>
              </div>
              {selectedPost.risk_logs && selectedPost.risk_logs.length > 0 ? (
                <div className="space-y-2 max-h-40 overflow-y-auto">
                  {selectedPost.risk_logs.map((log) => (
                    <div key={log.id} className="text-sm border-l-2 border-warning pl-2">
                      <div className="flex justify-between">
                        <span className="font-mono text-xs">{log.level}</span>
                        <span className="text-xs text-gray-500">{new Date(log.created_at).toLocaleString()}</span>
                      </div>
                      <div className="text-gray-700">规则：{log.rule}</div>
                      <div className="text-gray-500 text-xs break-all">匹配内容：{log.matched_content}</div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-sm text-gray-500">暂无风险日志</div>
              )}
            </div>

            {/* 帖子内容 */}
            <div className="mb-4 p-3 bg-base-200 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <FileText className="w-4 h-4" />
                <span className="font-medium">帖子内容</span>
              </div>
              <h4 className="font-semibold mb-2">{selectedPost.title}</h4>
              <div
  className="text-sm whitespace-pre-wrap max-h-60 overflow-y-auto bg-base-100 p-2 rounded"
  dangerouslySetInnerHTML={sanitizeHtml(selectedPost.content)}
/>
              {selectedPost.cover && (
                <img src={selectedPost.cover} alt="封面" className="mt-2 max-h-32 rounded object-cover" />
              )}
            </div>

            {/* 操作按钮 */}
            <div className="modal-action">
              <button className="btn" onClick={() => setIsModalOpen(false)}>关闭</button>
              <button
                className="btn btn-success"
                onClick={() => handleReview("approved")}
                disabled={reviewMutation.isPending}
              >
                <CheckCircle className="w-4 h-4 mr-1" />
                通过
              </button>
              <button
                className="btn btn-error"
                onClick={() => handleReview("rejected")}
                disabled={reviewMutation.isPending}
              >
                <XCircle className="w-4 h-4 mr-1" />
                拒绝
              </button>
              <button
                className="btn btn-warning"
                onClick={() => handleReview("pending")}
                disabled={reviewMutation.isPending}
              >
                <Clock className="w-4 h-4 mr-1" />
                待审
              </button>
            </div>
          </div>
        </dialog>
      )}
    </div>
  );
}