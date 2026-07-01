import { Banknote, Ticket, Trophy, Users, TrendingDown, TrendingUp } from "lucide-react"
import type { LucideIcon } from "lucide-react"

import { Card, CardContent } from "@/components/ui/card"
import { cn } from "@/lib/utils"
import type { DashboardStat } from "@/types/admin"

/**
 * Visual treatment for each stat kind: which Lucide icon to render and the
 * tint applied to its rounded badge. Keeps the four cards distinct at a glance.
 */
const ICON_CONFIG: Record<
  DashboardStat["icon"],
  { Icon: LucideIcon; badge: string }
> = {
  revenue: { Icon: Banknote, badge: "bg-primary/10 text-primary" },
  competitions: { Icon: Trophy, badge: "bg-botb-purple/12 text-botb-purple" },
  tickets: { Icon: Ticket, badge: "bg-botb-blue/12 text-botb-blue" },
  users: { Icon: Users, badge: "bg-botb-teal/12 text-botb-teal" },
}

/** Format a signed percentage for display, e.g. 12.4 -> "+12.4%". */
function formatDelta(deltaPct: number): string {
  const sign = deltaPct >= 0 ? "+" : ""
  return `${sign}${deltaPct.toFixed(1)}%`
}

export function StatCard({ stat }: { stat: DashboardStat }) {
  const { Icon, badge } = ICON_CONFIG[stat.icon]
  const isPositive = stat.deltaPct >= 0
  const DeltaIcon = isPositive ? TrendingUp : TrendingDown

  return (
    <Card>
      <CardContent className="flex flex-col gap-4">
        <div className="flex items-start justify-between gap-3">
          <span className="text-sm font-medium text-muted-foreground">
            {stat.label}
          </span>
          <span
            className={cn(
              "flex size-10 shrink-0 items-center justify-center rounded-lg",
              badge
            )}
          >
            <Icon className="size-5" />
          </span>
        </div>

        <div className="space-y-2">
          <p className="font-heading text-2xl font-semibold tracking-tight">
            {stat.value}
          </p>
          <p className="flex flex-wrap items-center gap-1.5 text-xs">
            <span
              className={cn(
                "inline-flex items-center gap-1 font-medium",
                isPositive ? "text-botb-green" : "text-destructive"
              )}
            >
              <DeltaIcon className="size-3.5" />
              {formatDelta(stat.deltaPct)}
            </span>
            <span className="text-muted-foreground">{stat.deltaLabel}</span>
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
