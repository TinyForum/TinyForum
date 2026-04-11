// src/app/users/[id]/page.tsx (服务端组件)
import { Suspense } from 'react';
import UserProfileClient from './UserProfileClient';
import { Metadata } from 'next';
export const metadata :Metadata = {
  title: "user",
  description: "User page"
}
export default async function UserProfilePage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const userId = Number(id);
  
  return (
    <Suspense fallback={<UserProfileSkeleton />}>
      <UserProfileClient userId={userId} />
    </Suspense>
  );
}

function UserProfileSkeleton() {
  return (
    <div className="max-w-4xl mx-auto space-y-4">
      <div className="skeleton h-40 w-full rounded-xl" />
      <div className="skeleton h-20 w-full rounded-xl" />
    </div>
  );
}