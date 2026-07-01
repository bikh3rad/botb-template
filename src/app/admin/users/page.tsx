import type { Metadata } from "next"
import { Ticket, TrendingUp, Users } from "lucide-react"

import { Card, CardContent } from "@/components/ui/card"
import { UsersTable } from "@/components/admin/UsersTable"
import { adminUsers, formatNumber, formatPence } from "@/lib/admin-data"

export const metadata: Metadata = {
  title: "Users & Tickets | Admin Console",
  description: "Search and manage registered players, their tickets and spend.",
}

/** Headline totals derived from the full user set for the summary chips. */
const totalUsers = adminUsers.length
const totalTickets = adminUsers.reduce((sum, user) => sum + user.ticketsOwned, 0)
const totalRevenuePence = adminUsers.reduce(
  (sum, user) => sum + user.totalSpentPence,
  0
)

const summaryStats = [
  {
    id: "users",
    label: "Total players",
    value: formatNumber(totalUsers),
    icon: Users,
    accent: "text-botb-blue",
    tint: "bg-botb-blue/10",
  },
  {
    id: "tickets",
    label: "Tickets owned",
    value: formatNumber(totalTickets),
    icon: Ticket,
    accent: "text-primary",
    tint: "bg-primary/10",
  },
  {
    id: "revenue",
    label: "Lifetime revenue",
    value: formatPence(totalRevenuePence),
    icon: TrendingUp,
    accent: "text-botb-green",
    tint: "bg-botb-green/10",
  },
] as const

export default function UsersPage() {
  return (
    <div className="space-y-6">
      <header className="space-y-1">
        <h2 className="font-heading text-2xl font-bold">Users &amp; Tickets</h2>
        <p className="text-sm text-muted-foreground">
          Search and manage registered players, their tickets and spend.
        </p>
      </header>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        {summaryStats.map((stat) => {
          const Icon = stat.icon
          return (
            <Card key={stat.id}>
              <CardContent className="flex items-center gap-4">
                <span
                  className={`flex size-11 shrink-0 items-center justify-center rounded-lg ${stat.tint}`}
                >
                  <Icon className={`size-5 ${stat.accent}`} />
                </span>
                <div className="min-w-0">
                  <p className="text-xs font-medium tracking-wide text-muted-foreground uppercase">
                    {stat.label}
                  </p>
                  <p className="font-heading text-xl font-semibold tabular-nums">
                    {stat.value}
                  </p>
                </div>
              </CardContent>
            </Card>
          )
        })}
      </div>

      <UsersTable />
    </div>
  )
}
