import { userApi } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";

 export function useUserProfile() {
 const { data: profile, isLoading } = useQuery({
    queryKey: ['user', userId],
    queryFn: () => userApi.getCurrentRole().then((r) => r.data.data),
  });
}