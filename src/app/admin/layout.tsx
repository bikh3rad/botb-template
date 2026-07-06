import type { Metadata } from "next"

import { AdminShell } from "@/components/admin/AdminShell"
import { AdminAuthProvider } from "@/lib/admin/auth-context"

export const metadata: Metadata = {
  title: "Admin Console | Competitions Platform",
  description: "Manage competitions, users, tickets and prize draws.",
  robots: { index: false, follow: false },
}

export default function AdminLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <AdminAuthProvider>
      <AdminShell>{children}</AdminShell>
    </AdminAuthProvider>
  )
}
