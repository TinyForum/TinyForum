// 骨架屏组件
export function PostListSkeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 5 }).map((_: unknown, i: number) => (
        <div key={i} className="skeleton h-28 w-full rounded-xl" />
      ))}
    </div>
  );
}
