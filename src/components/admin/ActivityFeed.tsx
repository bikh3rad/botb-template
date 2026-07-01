import { Banknote, Sparkles, Ticket, Trophy, UserPlus } from "lucide-react"
import type { LucideIcon } from "lucide-react"

import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card"
import { recentActivity } from "@/lib/admin-data"
import { cn } from "@/lib/utils"
import type { ActivityItem } from "@/types/admin"

/**
 * Icon + tint per activity type, so each event kind is recognisable at a glance
 * down the feed (ticket entries, sign-ups, draws, launches, payouts).
 */
const TYPE_CONFIG: Record<
  ActivityItem["type"],
  { Icon: LucideIcon; badge: string }
> = {
  entry: { Icon: Ticket, badge: "bg-primary/10 text-primary" },
  signup: { Icon: UserPlus, badge: "bg-botb-blue/12 text-botb-blue" },
  draw: { Icon: Trophy, badge: "bg-botb-purple/12 text-botb-purple" },
  competition: { Icon: Sparkles, badge: "bg-botb-teal/12 text-botb-teal" },
  payout: { Icon: Banknote, badge: "bg-botb-green/12 text-botb-green" },
}

export function ActivityFeed() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Activity</CardTitle>
        <CardDescription>Latest events across the platform</CardDescription>
      </CardHeader>
      <CardContent>
        <ul className="divide-y divide-border">
          {recentActivity.map((item) => {
            const { Icon, badge } = TYPE_CONFIG[item.type]
            return (
              <li key={item.id} className="flex items-start gap-3 py-3 first:pt-0 last:pb-0">
                <span
                  className={cn(
                    "flex size-8 shrink-0 items-center justify-center rounded-full",
                    badge
                  )}
                >
                  <Icon className="size-4" />
                </span>
                <div className="flex min-w-0 flex-1 items-start justify-between gap-3">
                  <p className="min-w-0 text-sm leading-snug text-foreground">
                    {item.message}
                  </p>
                  <time className="shrink-0 whitespace-nowrap pt-0.5 text-xs text-muted-foreground">
                    {item.timestamp}
                  </time>
                </div>
              </li>
            )
          })}
        </ul>
      </CardContent>
    </Card>
  )
}
