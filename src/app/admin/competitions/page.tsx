import type { Metadata } from "next"

import { CompetitionsTable } from "@/components/admin/CompetitionsTable"

export const metadata: Metadata = {
  title: "Competitions | Admin Console",
  description: "Manage active, upcoming and drawn competitions.",
}

export default function CompetitionsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="font-heading text-2xl font-semibold tracking-tight">
          Competitions
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Manage active, upcoming and drawn competitions.
        </p>
      </div>

      <CompetitionsTable />
    </div>
  )
}
