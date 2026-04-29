import type { LucideIcon } from "lucide-react";

export interface NavItem {
  key: string;
  name: string;
  href: string;
  icon: LucideIcon;
  requiresAuth: boolean;
}
