"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { LucideIcon } from "lucide-react";

interface NavItem {
  key: string;
  name: string;
  href: string;
  icon: LucideIcon;
  requiresAuth: boolean;
}

interface NavLinksProps {
  items: NavItem[];
}

export default function NavLinks({ items }: NavLinksProps) {
  const pathname = usePathname();

  const isActive = (href: string) => {
    if (href === "/") {
      return pathname === href;
    }
    return pathname.startsWith(href);
  };

  return (
    <>
      {items.map((item) => {
        const active = isActive(item.href);
        const Icon = item.icon;

        return (
          <Link
            key={item.key}
            href={item.href}
            className={`btn btn-sm gap-2 transition-all duration-200 ${
              active 
                ? "btn-primary shadow-md" 
                : "btn-ghost hover:bg-primary/10"
            }`}
          >
            <Icon className={`w-4 h-4 ${active ? "animate-pulse" : ""}`} />
            <span>{item.name.charAt(0).toUpperCase() + item.name.slice(1)}</span>
            {active && (
              <span className="absolute bottom-0 left-1/2 transform -translate-x-1/2 w-1 h-1 bg-primary rounded-full" />
            )}
          </Link>
        );
      })}
    </>
  );
}