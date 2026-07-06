"use client"

import * as React from "react"
import { usePathname } from "next/navigation"
import { Bell, Menu, Search, X } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { adminNav } from "@/lib/admin-data"
import { AdminSidebar } from "@/components/admin/AdminSidebar"
import { useAdminAuth } from "@/lib/admin/auth-context"

/** Two-letter initials from a name, e.g. "Ada Lovelace" -> "AL". */
function initials(name: string | undefined): string {
  if (!name) return "AD"
  const parts = name.trim().split(/\s+/)
  const first = parts[0]?.[0] ?? ""
  const last = parts.length > 1 ? parts[parts.length - 1][0] : ""
  return (first + last).toUpperCase() || "AD"
}

function usePageTitle(): string {
  const pathname = usePathname()
  const match = [...adminNav]
    .sort((a, b) => b.href.length - a.href.length)
    .find((item) =>
      item.href === "/admin"
        ? pathname === "/admin"
        : pathname === item.href || pathname.startsWith(`${item.href}/`)
    )
  return match?.label ?? "Admin"
}

export function AdminShell({ children }: { children: React.ReactNode }) {
  const [mobileOpen, setMobileOpen] = React.useState(false)
  const pathname = usePathname()
  const title = usePageTitle()

  // Close the mobile drawer whenever the route changes.
  React.useEffect(() => {
    setMobileOpen(false)
  }, [pathname])

  // The login page owns its own full-screen layout — render it without the
  // authenticated console chrome.
  if (pathname === "/admin/login") {
    return <>{children}</>
  }

  return (
    <div className="min-h-screen bg-muted/40 text-foreground">
      {/* Desktop sidebar */}
      <aside className="fixed inset-y-0 left-0 z-30 hidden w-64 border-r border-sidebar-border lg:block">
        <AdminSidebar />
      </aside>

      {/* Mobile drawer */}
      {mobileOpen && (
        <div className="fixed inset-0 z-50 lg:hidden">
          <div
            className="absolute inset-0 bg-black/50"
            onClick={() => setMobileOpen(false)}
            aria-hidden
          />
          <div className="absolute inset-y-0 left-0 w-64 max-w-[80%] border-r border-sidebar-border shadow-xl">
            <Button
              variant="ghost"
              size="icon"
              className="absolute top-3.5 right-3 z-10"
              onClick={() => setMobileOpen(false)}
              aria-label="Close menu"
            >
              <X className="size-5" />
            </Button>
            <AdminSidebar onNavigate={() => setMobileOpen(false)} />
          </div>
        </div>
      )}

      <div className="lg:pl-64">
        {/* Top bar */}
        <header className="sticky top-0 z-20 flex h-16 items-center gap-3 border-b border-border bg-background/95 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/80 sm:px-6">
          <Button
            variant="ghost"
            size="icon"
            className="lg:hidden"
            onClick={() => setMobileOpen(true)}
            aria-label="Open menu"
          >
            <Menu className="size-5" />
          </Button>

          <h1 className="font-heading text-lg font-semibold">{title}</h1>

          <div className="relative ml-auto hidden max-w-xs flex-1 sm:block">
            <Search className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="search"
              placeholder="Search…"
              className="pl-9"
              aria-label="Search"
            />
          </div>

          <Button
            variant="ghost"
            size="icon"
            className="relative ml-auto sm:ml-0"
            aria-label="Notifications"
          >
            <Bell className="size-5" />
            <span className="absolute top-2 right-2 size-2 rounded-full bg-primary ring-2 ring-background" />
          </Button>

          <AdminAvatar />
        </header>

        <main className="p-4 sm:p-6">{children}</main>
      </div>
    </div>
  )
}

function AdminAvatar() {
  const { admin } = useAdminAuth()
  return (
    <Avatar className="bg-primary/10" title={admin?.email}>
      <AvatarFallback className="bg-transparent text-botb-orange-hover dark:text-primary">
        {initials(admin?.name)}
      </AvatarFallback>
    </Avatar>
  )
}
