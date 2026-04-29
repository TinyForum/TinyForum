"use client";

import PostLoginHandler from "@/features/auth/PostLoginHandler";
import { ReactNode } from "react";

export default function AuthProvider({ children }: { children: ReactNode }) {
  return (
    <>
      <PostLoginHandler />
      {children}
    </>
  );
}
