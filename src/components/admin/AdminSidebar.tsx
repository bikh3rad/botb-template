"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import {
  LayoutDashboard,
  Trophy,
  Users,
  Ticket,
  Tags,
  Images,
  FileText,
  ShieldCheck,
  UserCog,
  LogOut,
  type LucideIcon,
} from "lucide-react"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { useAdminAuth } from "@/lib/admin/auth-context"

interface NavItem {
  label: string
  href: string
  icon: LucideIcon
  /** Only shown to superadmins. */
  superadmin?: boolean
}

const navItems: NavItem[] = [
  { label: "Dashboard", href: "/admin", icon: LayoutDashboard },
  { label: "Competitions", href: "/admin/competitions", icon: Trophy },
  { label: "Categories", href: "/admin/categories", icon: Tags },
  { label: "Users & Tickets", href: "/admin/users", icon: Users },
  { label: "Winners & Draws", href: "/admin/winners", icon: Ticket },
  { label: "Media Library", href: "/admin/media", icon: Images },
  { label: "Site Texts", href: "/admin/texts", icon: FileText },
  { label: "Admin Accounts", href: "/admin/accounts", icon: UserCog, superadmin: true },
]

function isActive(pathname: string, href: string): boolean {
  if (href === "/admin") return pathname === "/admin"
  return pathname === href || pathname.startsWith(`${href}/`)
}

export function AdminSidebar({ onNavigate }: { onNavigate?: () => void }) {
  const pathname = usePathname()
  const { admin, isSuperadmin, logout } = useAdminAuth()

  const visible = navItems.filter((item) => !item.superadmin || isSuperadmin)

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
        {visible.map((item) => {
          const Icon = item.icon
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

      <div className="space-y-2 border-t border-sidebar-border p-4">
        <div className="rounded-lg bg-sidebar-accent px-3 py-2.5 text-xs text-muted-foreground">
          Signed in as{" "}
          <span className="font-medium text-sidebar-foreground">
            {admin?.name ?? "…"}
          </span>
          {admin?.role && (
            <span className="mt-0.5 block text-[11px] capitalize text-muted-foreground/80">
              {admin.role}
            </span>
          )}
        </div>
        <Button
          variant="ghost"
          size="sm"
          className="w-full justify-start gap-2 text-sidebar-foreground/70"
          onClick={() => void logout()}
        >
          <LogOut className="size-4" />
          Sign out
        </Button>
      </div>
    </div>
  )
}
