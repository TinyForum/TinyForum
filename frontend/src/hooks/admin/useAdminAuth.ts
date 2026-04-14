
import { useEffect, useState } from "react";
import { useAuthStore } from "@/store/auth";
import { useRouter } from "next/navigation";

export function useAdminAuth() {
  const { user, isAuthenticated } = useAuthStore();
  const router = useRouter();
  const [isCheckingAuth, setIsCheckingAuth] = useState(true);

  useEffect(() => {
    const checkAuth = async () => {
      await new Promise(resolve => setTimeout(resolve, 100));
      setIsCheckingAuth(false);
      
      if (!isAuthenticated) {
        router.push('/auth/login');
      } else if (user?.role !== "admin" && user?.role !== "super_admin") {
        router.push('/');
      }
    };
    
    checkAuth();
  }, [isAuthenticated, user, router]);

  const isAdmin = isAuthenticated && 
    (user?.role === "admin" || user?.role === "super_admin");

  return { isCheckingAuth, isAdmin, user };
}
