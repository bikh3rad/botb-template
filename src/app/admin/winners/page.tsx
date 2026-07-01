import type { Metadata } from "next"

import { WinnersTable } from "@/components/admin/WinnersTable"

export const metadata: Metadata = {
  title: "Winners & Draws | Admin Console",
  description: "Past draws, winners and prizes awarded.",
}

export default function WinnersPage() {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="font-heading text-2xl font-semibold tracking-tight">
          Winners &amp; Draws
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Past draws, winners and prizes awarded.
        </p>
      </div>

      <WinnersTable />
    </div>
  )
}
