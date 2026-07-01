"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import {
  LayoutDashboard,
  Trophy,
  Users,
  Ticket,
  ShieldCheck,
  type LucideIcon,
} from "lucide-react"

import { cn } from "@/lib/utils"
import { adminNav } from "@/lib/admin-data"
import type { AdminNavItem } from "@/types/admin"

const iconMap: Record<AdminNavItem["icon"], LucideIcon> = {
  dashboard: LayoutDashboard,
  competitions: Trophy,
  users: Users,
  winners: Ticket,
}

function isActive(pathname: string, href: string): boolean {
  if (href === "/admin") return pathname === "/admin"
  return pathname === href || pathname.startsWith(`${href}/`)
}

export function AdminSidebar({ onNavigate }: { onNavigate?: () => void }) {
  const pathname = usePathname()

  return (
    <div className="flex h-full flex-col bg-sidebar text-sidebar-foreground">
      <div className="flex h-16 items-center gap-2.5 border-b border-sidebar-border px-5">
        <span className="flex size-9 items-center justify-center rounded-lg bg-primary text-primary-foreground">
          <ShieldCheck className="size-5" />
        </span>
        <div className="flex flex-col leading-tight">
          <span className="font-heading text-sm font-semibold">Admin Console</span>
          <span className="text-xs text-muted-foreground">Competitions Platform</span>
        </div>
      </div>

      <nav className="flex-1 space-y-1 overflow-y-auto p-3">
        {adminNav.map((item) => {
          const Icon = iconMap[item.icon]
          const active = isActive(pathname, item.href)
          return (
            <Link
              key={item.href}
              href={item.href}
              onClick={onNavigate}
              aria-current={active ? "page" : undefined}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                active
                  ? "bg-primary/10 text-botb-orange-hover dark:bg-primary/15 dark:text-primary"
                  : "text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-foreground"
              )}
            >
              <Icon className="size-4.5 shrink-0" />
              {item.label}
            </Link>
          )
        })}
      </nav>

      <div className="border-t border-sidebar-border p-4">
        <div className="rounded-lg bg-sidebar-accent px-3 py-2.5 text-xs text-muted-foreground">
          Signed in as{" "}
          <span className="font-medium text-sidebar-foreground">Admin</span>
        </div>
      </div>
    </div>
  )
}
