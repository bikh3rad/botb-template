"use client"

import * as React from "react"
import {
  Banknote,
  Trophy,
  Ticket,
  Users,
  Loader2,
  AlertCircle,
} from "lucide-react"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { cn } from "@/lib/utils"
import { apiGet, ApiError } from "@/lib/admin/client"
import { formatPence, formatNumber, formatDateTime } from "@/lib/admin/format"
import type { Competition } from "@/types/admin-api"

interface CompetitionsResponse {
  competitions: Competition[]
}
interface CountResponse {
  total: number
}
interface AuditResponse {
  entries: {
    id: string
    actor_email: string
    action: string
    entity_type: string
    entity_id: string
    reason?: string
    created_at: string
  }[]
}

interface Stat {
  key: string
  label: string
  value: string
  hint: string
  icon: typeof Banknote
  badge: string
}

// Human labels for audit action codes, e.g. "competition.update" -> "Updated competition".
const ACTION_VERB: Record<string, string> = {
  create: "Created",
  update: "Updated",
  delete: "Deleted",
  suspend: "Suspended",
  activate: "Activated",
  void: "Voided",
  run: "Ran draw for",
  reassign: "Reassigned",
  upload: "Uploaded media for",
}

function describeAction(action: string, entityType: string): string {
  const verb = action.split(".").pop() ?? action
  return `${ACTION_VERB[verb] ?? verb} ${entityType.replace(/_/g, " ")}`
}

export function DashboardOverview() {
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)
  const [stats, setStats] = React.useState<Stat[]>([])
  const [byCategory, setByCategory] = React.useState<
    { name: string; revenuePence: number }[]
  >([])
  const [activity, setActivity] = React.useState<AuditResponse["entries"]>([])

  React.useEffect(() => {
    let cancelled = false

    async function load() {
      setLoading(true)
      setError(null)
      try {
        // All numbers come from the real backend — no mock data.
        const [draft, live, closed, usersPage, drawsPage, audit] =
          await Promise.all([
            apiGet<CompetitionsResponse>(
              "/apis/competition/v1/competitions?status=draft",
            ),
            apiGet<CompetitionsResponse>(
              "/apis/competition/v1/competitions?status=live",
            ),
            apiGet<CompetitionsResponse>(
              "/apis/competition/v1/competitions?status=closed",
            ),
            apiGet<CountResponse>("/apis/user/v1/admin/users?limit=1"),
            apiGet<CountResponse>("/apis/draw/v1/admin/draws?limit=1"),
            apiGet<AuditResponse>("/apis/competition/v1/admin/audit?limit=8"),
          ])

        if (cancelled) return

        const all = [
          ...draft.competitions,
          ...live.competitions,
          ...closed.competitions,
        ]
        const ticketsSold = all.reduce((n, c) => n + c.tickets_sold, 0)
        const revenuePence = all.reduce(
          (n, c) => n + c.tickets_sold * c.ticket_price_pence,
          0,
        )

        setStats([
          {
            key: "revenue",
            label: "Revenue",
            value: formatPence(revenuePence),
            hint: "Tickets sold × price, all competitions",
            icon: Banknote,
            badge: "bg-primary/10 text-primary",
          },
          {
            key: "active",
            label: "Active competitions",
            value: formatNumber(live.competitions.length),
            hint: `${all.length} total (${draft.competitions.length} draft, ${closed.competitions.length} closed)`,
            icon: Trophy,
            badge: "bg-botb-purple/12 text-botb-purple",
          },
          {
            key: "tickets",
            label: "Tickets sold",
            value: formatNumber(ticketsSold),
            hint: "Across every competition",
            icon: Ticket,
            badge: "bg-botb-blue/12 text-botb-blue",
          },
          {
            key: "users",
            label: "Registered users",
            value: formatNumber(usersPage.total),
            hint: `${drawsPage.total} draws run`,
            icon: Users,
            badge: "bg-botb-teal/12 text-botb-teal",
          },
        ])

        // Real revenue distribution by category.
        const cat = new Map<string, number>()
        for (const c of all) {
          const name = c.category_name ?? "Uncategorised"
          cat.set(name, (cat.get(name) ?? 0) + c.tickets_sold * c.ticket_price_pence)
        }
        setByCategory(
          [...cat.entries()]
            .map(([name, revenuePence]) => ({ name, revenuePence }))
            .sort((a, b) => b.revenuePence - a.revenuePence),
        )

        setActivity(audit.entries)
      } catch (err) {
        if (!cancelled) {
          setError(
            err instanceof ApiError ? err.message : "Failed to load dashboard",
          )
        }
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    void load()
    return () => {
      cancelled = true
    }
  }, [])

  if (loading) {
    return (
      <div className="flex h-64 items-center justify-center text-muted-foreground">
        <Loader2 className="mr-2 size-5 animate-spin" /> Loading dashboard…
      </div>
    )
  }

  if (error) {
    return (
      <p className="inline-flex items-center gap-2 rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">
        <AlertCircle className="size-4 shrink-0" />
        {error}
      </p>
    )
  }

  const maxRevenue = Math.max(1, ...byCategory.map((c) => c.revenuePence))

  return (
    <div className="space-y-6">
      <section
        aria-label="Key metrics"
        className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4"
      >
        {stats.map((stat) => {
          const Icon = stat.icon
          return (
            <Card key={stat.key}>
              <CardContent className="flex flex-col gap-4">
                <div className="flex items-start justify-between gap-3">
                  <span className="text-sm font-medium text-muted-foreground">
                    {stat.label}
                  </span>
                  <span
                    className={cn(
                      "flex size-10 shrink-0 items-center justify-center rounded-lg",
                      stat.badge,
                    )}
                  >
                    <Icon className="size-5" />
                  </span>
                </div>
                <div className="space-y-1">
                  <p className="font-heading text-2xl font-semibold tracking-tight">
                    {stat.value}
                  </p>
                  <p className="text-xs text-muted-foreground">{stat.hint}</p>
                </div>
              </CardContent>
            </Card>
          )
        })}
      </section>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle className="text-base">Revenue by category</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {byCategory.length === 0 && (
              <p className="text-sm text-muted-foreground">No revenue yet.</p>
            )}
            {byCategory.map((c) => (
              <div key={c.name} className="space-y-1">
                <div className="flex items-center justify-between text-sm">
                  <span className="font-medium">{c.name}</span>
                  <span className="text-muted-foreground">
                    {formatPence(c.revenuePence)}
                  </span>
                </div>
                <div className="h-2 overflow-hidden rounded-full bg-muted">
                  <div
                    className="h-full rounded-full bg-primary"
                    style={{
                      width: `${(c.revenuePence / maxRevenue) * 100}%`,
                    }}
                  />
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Recent admin activity</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {activity.length === 0 && (
              <p className="text-sm text-muted-foreground">
                No admin actions yet.
              </p>
            )}
            {activity.map((a) => (
              <div key={a.id} className="flex flex-col gap-0.5 text-sm">
                <span className="font-medium">
                  {describeAction(a.action, a.entity_type)}
                </span>
                <span className="text-xs text-muted-foreground">
                  {a.actor_email || "system"} · {formatDateTime(a.created_at)}
                </span>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
