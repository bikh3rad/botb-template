import type { Metadata } from "next"

import { UsersTable } from "@/components/admin/UsersTable"

export const metadata: Metadata = {
  title: "Users & Tickets | Admin Console",
  description: "Search and manage registered players, their tickets and spend.",
}

export default function UsersPage() {
  return (
    <div className="space-y-6">
      <header className="space-y-1">
        <h2 className="font-heading text-2xl font-bold">Users &amp; Tickets</h2>
        <p className="text-sm text-muted-foreground">
          Search and manage registered players, their tickets and spend.
        </p>
      </header>

      <UsersTable />
    </div>
  )
}
