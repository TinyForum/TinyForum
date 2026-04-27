"use client";

import { ReactNode } from "react";
import PostLoginHandler from "../auth/PostLoginHandler";

export default function AuthProvider({ children }: { children: ReactNode }) {
  return (
    <>
      <PostLoginHandler />
      {children}
    </>
  );
}
