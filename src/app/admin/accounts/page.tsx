import type { Metadata } from "next"

import { AdminAccountsManager } from "@/components/admin/AdminAccountsManager"

export const metadata: Metadata = {
  title: "Admin Accounts | Admin Console",
  description: "Manage admin console accounts and permissions.",
}

export default function AdminAccountsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="font-heading text-2xl font-semibold tracking-tight">
          Admin Accounts
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Manage admin console accounts and permissions.
        </p>
      </div>

      <AdminAccountsManager />
    </div>
  )
}
