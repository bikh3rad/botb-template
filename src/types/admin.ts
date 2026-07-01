// Data shapes for the admin panel. These mirror what a real API/database layer
// would return, so the mock data in `src/lib/admin-data.ts` can be swapped for
// live queries without touching the UI.

/** Lifecycle status of a competition. */
export type CompetitionStatus =
  | "live"
  | "ending-soon"
  | "sold-out"
  | "drawn"
  | "draft";

/** Status of a scheduled/completed prize draw. */
export type DrawStatus = "completed" | "pending";

/** Account standing of a registered user. */
export type UserStatus = "active" | "vip" | "suspended";

/** A single navigation entry in the admin sidebar. */
export interface AdminNavItem {
  label: string;
  href: string;
  /** Lucide icon name key, resolved in the sidebar component. */
  icon: "dashboard" | "competitions" | "users" | "winners";
}

/** A headline metric rendered as a stat card on the dashboard. */
export interface DashboardStat {
  id: string;
  label: string;
  /** Preformatted display value, e.g. "£482,910" or "12,480". */
  value: string;
  /** Percentage change vs. the previous period, e.g. 12.4 or -3.1. */
  deltaPct: number;
  /** Context for the delta, e.g. "vs last month". */
  deltaLabel: string;
  icon: "revenue" | "competitions" | "tickets" | "users";
}

/** One point on the revenue/sales chart. */
export interface RevenuePoint {
  /** Short month label, e.g. "Jan". */
  month: string;
  /** Revenue for the period, in whole GBP. */
  revenue: number;
  /** Tickets sold for the period. */
  tickets: number;
}

/** An entry in the dashboard "recent activity" feed. */
export interface ActivityItem {
  id: string;
  type: "entry" | "signup" | "draw" | "competition" | "payout";
  /** Human-readable summary of the event. */
  message: string;
  /** Relative time label, e.g. "2 min ago". */
  timestamp: string;
}

/** A competition managed from the admin panel. */
export interface AdminCompetition {
  id: string;
  title: string;
  prize: string;
  /** Ticket price in whole pence, e.g. 125 = £1.25. */
  ticketPricePence: number;
  ticketsSold: number;
  ticketsTotal: number;
  status: CompetitionStatus;
  /** ISO date string (YYYY-MM-DD) the competition closes. */
  endDate: string;
  category: string;
}

/** A registered user/customer. */
export interface AdminUser {
  id: string;
  name: string;
  email: string;
  ticketsOwned: number;
  /** Lifetime spend in whole pence. */
  totalSpentPence: number;
  /** ISO date string (YYYY-MM-DD) the user registered. */
  joinDate: string;
  status: UserStatus;
}

/** A past or scheduled prize draw and its winner. */
export interface AdminDraw {
  id: string;
  competition: string;
  /** Winner name; null when the draw has not been run yet. */
  winner: string | null;
  prize: string;
  /** ISO date string (YYYY-MM-DD) of the draw. */
  drawDate: string;
  /** Winning ticket number; null until drawn. */
  ticketNumber: number | null;
  status: DrawStatus;
}
