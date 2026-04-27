// src/app/posts/[id]/page.tsx (服务端组件)
import { Suspense } from "react";
import PostDetailClient from "./PostDetailClient";
import CommentSection from "@/components/post/CommentSection";

export default async function PostDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  console.log(params);
  const { id } = await params;
  const postId = Number(id);

  return (
    <div className="max-w-3xl mx-auto">
      <Suspense fallback={<PostDetailSkeleton />}>
        <PostDetailClient postId={postId} />
      </Suspense>

      {/* Comments section can be client component too */}
      <div className="card bg-base-100 border border-base-300 shadow-sm">
        <div className="card-body p-6 lg:p-8">
          <CommentSection postId={postId} />
        </div>
      </div>
    </div>
  );
}

function PostDetailSkeleton() {
  return (
    <div className="space-y-4">
      <div className="skeleton h-10 w-3/4" />
      <div className="skeleton h-4 w-1/2" />
      <div className="skeleton h-64 w-full" />
    </div>
  );
}
