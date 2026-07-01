import type {
  ActivityItem,
  AdminCompetition,
  AdminDraw,
  AdminNavItem,
  AdminUser,
  DashboardStat,
  RevenuePoint,
} from "@/types/admin";

// ---------------------------------------------------------------------------
// Formatters — shared display helpers so every table/card renders consistently.
// ---------------------------------------------------------------------------

const gbp = new Intl.NumberFormat("en-GB", {
  style: "currency",
  currency: "GBP",
  minimumFractionDigits: 0,
  maximumFractionDigits: 0,
});

const gbpPrecise = new Intl.NumberFormat("en-GB", {
  style: "currency",
  currency: "GBP",
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
});

/** Format whole pence into a GBP string, e.g. 125 -> "£1.25". */
export function formatPence(pence: number): string {
  return gbpPrecise.format(pence / 100);
}

/** Format whole GBP into a rounded currency string, e.g. 482910 -> "£482,910". */
export function formatCurrency(pounds: number): string {
  return gbp.format(pounds);
}

/** Format a number with thousands separators, e.g. 12480 -> "12,480". */
export function formatNumber(value: number): string {
  return new Intl.NumberFormat("en-GB").format(value);
}

/** Format an ISO date (YYYY-MM-DD) as "6 Jul 2026". */
export function formatDate(iso: string): string {
  const [y, m, d] = iso.split("-").map(Number);
  const date = new Date(Date.UTC(y, m - 1, d));
  return date.toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
    timeZone: "UTC",
  });
}

/** Percentage of tickets sold, clamped to 0–100. */
export function soldPercent(sold: number, total: number): number {
  if (total <= 0) return 0;
  return Math.min(100, Math.round((sold / total) * 100));
}

// ---------------------------------------------------------------------------
// Navigation
// ---------------------------------------------------------------------------

export const adminNav: AdminNavItem[] = [
  { label: "Dashboard", href: "/admin", icon: "dashboard" },
  { label: "Competitions", href: "/admin/competitions", icon: "competitions" },
  { label: "Users & Tickets", href: "/admin/users", icon: "users" },
  { label: "Winners & Draws", href: "/admin/winners", icon: "winners" },
];

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

export const dashboardStats: DashboardStat[] = [
  {
    id: "revenue",
    label: "Total Revenue",
    value: formatCurrency(482910),
    deltaPct: 12.4,
    deltaLabel: "vs last month",
    icon: "revenue",
  },
  {
    id: "competitions",
    label: "Active Competitions",
    value: "24",
    deltaPct: 4.2,
    deltaLabel: "vs last month",
    icon: "competitions",
  },
  {
    id: "tickets",
    label: "Tickets Sold",
    value: formatNumber(1284630),
    deltaPct: 18.9,
    deltaLabel: "vs last month",
    icon: "tickets",
  },
  {
    id: "users",
    label: "Registered Users",
    value: formatNumber(86214),
    deltaPct: -2.1,
    deltaLabel: "vs last month",
    icon: "users",
  },
];

export const revenueSeries: RevenuePoint[] = [
  { month: "Jan", revenue: 268400, tickets: 742000 },
  { month: "Feb", revenue: 291200, tickets: 803000 },
  { month: "Mar", revenue: 324800, tickets: 889000 },
  { month: "Apr", revenue: 312500, tickets: 861000 },
  { month: "May", revenue: 358900, tickets: 972000 },
  { month: "Jun", revenue: 402300, tickets: 1094000 },
  { month: "Jul", revenue: 431700, tickets: 1168000 },
  { month: "Aug", revenue: 418600, tickets: 1131000 },
  { month: "Sep", revenue: 447900, tickets: 1210000 },
  { month: "Oct", revenue: 462100, tickets: 1247000 },
  { month: "Nov", revenue: 471300, tickets: 1268000 },
  { month: "Dec", revenue: 482910, tickets: 1284630 },
];

export const recentActivity: ActivityItem[] = [
  {
    id: "a1",
    type: "entry",
    message: "Olivia Bennett bought 25 tickets for “Audi RS3 Carbon Black”",
    timestamp: "2 min ago",
  },
  {
    id: "a2",
    type: "signup",
    message: "New user registered — james.okafor@gmail.com",
    timestamp: "9 min ago",
  },
  {
    id: "a3",
    type: "draw",
    message: "Draw completed for “£50,000 Tax-Free Cash” — winner notified",
    timestamp: "41 min ago",
  },
  {
    id: "a4",
    type: "payout",
    message: "Prize payout of £50,000 marked as sent to Marcus Reid",
    timestamp: "1 hour ago",
  },
  {
    id: "a5",
    type: "competition",
    message: "“Defender D350 X” competition went live",
    timestamp: "3 hours ago",
  },
  {
    id: "a6",
    type: "entry",
    message: "Priya Sharma bought 100 tickets for “Golden Boot Supercar”",
    timestamp: "4 hours ago",
  },
  {
    id: "a7",
    type: "signup",
    message: "New user registered — thomas.wright@outlook.com",
    timestamp: "5 hours ago",
  },
  {
    id: "a8",
    type: "competition",
    message: "“Summer Festival Instant Wins” sold out",
    timestamp: "6 hours ago",
  },
];

// ---------------------------------------------------------------------------
// Competitions
// ---------------------------------------------------------------------------

export const adminCompetitions: AdminCompetition[] = [
  {
    id: "c-13736",
    title: "Win a Defender D350 X",
    prize: "Land Rover Defender D350 X",
    ticketPricePence: 20,
    ticketsSold: 184300,
    ticketsTotal: 250000,
    status: "live",
    endDate: "2026-07-18",
    category: "Cars",
  },
  {
    id: "c-13589",
    title: "Audi RS3 Carbon Black",
    prize: "Audi RS3 Carbon Black Edition",
    ticketPricePence: 6,
    ticketsSold: 421900,
    ticketsTotal: 500000,
    status: "live",
    endDate: "2026-07-22",
    category: "Cars",
  },
  {
    id: "c-13344",
    title: "£1.2M London Home in Zone 1",
    prize: "£1,200,000 London Home",
    ticketPricePence: 100,
    ticketsSold: 968200,
    ticketsTotal: 1000000,
    status: "ending-soon",
    endDate: "2026-07-06",
    category: "Property",
  },
  {
    id: "c-13728",
    title: "The Golden Boot Supercar",
    prize: "6 Supercars + £100k Cash",
    ticketPricePence: 250,
    ticketsSold: 312400,
    ticketsTotal: 450000,
    status: "live",
    endDate: "2026-07-30",
    category: "Cars",
  },
  {
    id: "c-13551",
    title: "Summer Festival Instant Wins",
    prize: "500+ Instant Prizes",
    ticketPricePence: 119,
    ticketsSold: 275000,
    ticketsTotal: 275000,
    status: "sold-out",
    endDate: "2026-07-04",
    category: "Instant Wins",
  },
  {
    id: "c-13548",
    title: "Wimbledon Wonders Instant Wins",
    prize: "Centre Court Package + Tech",
    ticketPricePence: 119,
    ticketsSold: 143700,
    ticketsTotal: 300000,
    status: "live",
    endDate: "2026-07-14",
    category: "Instant Wins",
  },
  {
    id: "c-13982",
    title: "Lifestyle Competition",
    prize: "Cars, Cash, Watches & Tech",
    ticketPricePence: 50,
    ticketsSold: 58900,
    ticketsTotal: 400000,
    status: "live",
    endDate: "2026-08-02",
    category: "Lifestyle",
  },
  {
    id: "c-13734",
    title: "Win Free Tech in App!",
    prize: "AirPods, Switch 2 & Ninja Creami",
    ticketPricePence: 0,
    ticketsSold: 96400,
    ticketsTotal: 150000,
    status: "live",
    endDate: "2026-07-11",
    category: "Tech",
  },
  {
    id: "c-13210",
    title: "£50,000 Tax-Free Cash",
    prize: "£50,000 Cash",
    ticketPricePence: 89,
    ticketsSold: 210000,
    ticketsTotal: 210000,
    status: "drawn",
    endDate: "2026-06-28",
    category: "Cash",
  },
  {
    id: "c-13099",
    title: "Toshiba 65\" 4K TV Bundle",
    prize: "Toshiba 65\" 4K TV + Soundbar",
    ticketPricePence: 45,
    ticketsSold: 84100,
    ticketsTotal: 120000,
    status: "drawn",
    endDate: "2026-06-20",
    category: "Tech",
  },
  {
    id: "c-14002",
    title: "Porsche 911 GT3 RS",
    prize: "Porsche 911 GT3 RS (992)",
    ticketPricePence: 199,
    ticketsSold: 0,
    ticketsTotal: 600000,
    status: "draft",
    endDate: "2026-08-15",
    category: "Cars",
  },
  {
    id: "c-14005",
    title: "Dream Kitchen Makeover",
    prize: "£25,000 Kitchen Renovation",
    ticketPricePence: 75,
    ticketsSold: 0,
    ticketsTotal: 180000,
    status: "draft",
    endDate: "2026-08-20",
    category: "Lifestyle",
  },
];

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

export const adminUsers: AdminUser[] = [
  { id: "u-1042", name: "Olivia Bennett", email: "olivia.bennett@gmail.com", ticketsOwned: 1284, totalSpentPence: 214500, joinDate: "2024-02-11", status: "vip" },
  { id: "u-1043", name: "James Okafor", email: "james.okafor@gmail.com", ticketsOwned: 12, totalSpentPence: 1800, joinDate: "2026-06-30", status: "active" },
  { id: "u-1044", name: "Priya Sharma", email: "priya.sharma@outlook.com", ticketsOwned: 940, totalSpentPence: 167300, joinDate: "2024-08-03", status: "vip" },
  { id: "u-1045", name: "Marcus Reid", email: "marcus.reid@yahoo.com", ticketsOwned: 356, totalSpentPence: 58900, joinDate: "2025-01-19", status: "active" },
  { id: "u-1046", name: "Sophie Turner", email: "sophie.turner@gmail.com", ticketsOwned: 78, totalSpentPence: 9200, joinDate: "2025-11-22", status: "active" },
  { id: "u-1047", name: "Thomas Wright", email: "thomas.wright@outlook.com", ticketsOwned: 5, totalSpentPence: 600, joinDate: "2026-06-29", status: "active" },
  { id: "u-1048", name: "Amara Nwosu", email: "amara.nwosu@gmail.com", ticketsOwned: 2140, totalSpentPence: 389400, joinDate: "2023-05-14", status: "vip" },
  { id: "u-1049", name: "Daniel Cohen", email: "daniel.cohen@icloud.com", ticketsOwned: 214, totalSpentPence: 31200, joinDate: "2025-03-27", status: "active" },
  { id: "u-1050", name: "Isabella Rossi", email: "isabella.rossi@gmail.com", ticketsOwned: 640, totalSpentPence: 98700, joinDate: "2024-10-08", status: "active" },
  { id: "u-1051", name: "Liam Murphy", email: "liam.murphy@hotmail.com", ticketsOwned: 0, totalSpentPence: 0, joinDate: "2026-06-25", status: "suspended" },
  { id: "u-1052", name: "Chloe Adams", email: "chloe.adams@gmail.com", ticketsOwned: 428, totalSpentPence: 72100, joinDate: "2024-12-01", status: "active" },
  { id: "u-1053", name: "Noah Patel", email: "noah.patel@gmail.com", ticketsOwned: 1890, totalSpentPence: 301500, joinDate: "2023-09-30", status: "vip" },
  { id: "u-1054", name: "Emily Clarke", email: "emily.clarke@outlook.com", ticketsOwned: 96, totalSpentPence: 12400, joinDate: "2025-07-16", status: "active" },
  { id: "u-1055", name: "Oliver Hughes", email: "oliver.hughes@yahoo.com", ticketsOwned: 52, totalSpentPence: 6800, joinDate: "2026-01-04", status: "active" },
  { id: "u-1056", name: "Grace Kim", email: "grace.kim@gmail.com", ticketsOwned: 774, totalSpentPence: 129900, joinDate: "2024-04-22", status: "active" },
  { id: "u-1057", name: "Ethan Baker", email: "ethan.baker@icloud.com", ticketsOwned: 18, totalSpentPence: 2300, joinDate: "2026-05-18", status: "active" },
  { id: "u-1058", name: "Ava Morgan", email: "ava.morgan@gmail.com", ticketsOwned: 1520, totalSpentPence: 248600, joinDate: "2023-11-11", status: "vip" },
  { id: "u-1059", name: "Lucas Silva", email: "lucas.silva@hotmail.com", ticketsOwned: 240, totalSpentPence: 35700, joinDate: "2025-02-08", status: "active" },
  { id: "u-1060", name: "Mia Fischer", email: "mia.fischer@gmail.com", ticketsOwned: 61, totalSpentPence: 7900, joinDate: "2025-09-13", status: "active" },
  { id: "u-1061", name: "Benjamin Cole", email: "benjamin.cole@outlook.com", ticketsOwned: 405, totalSpentPence: 66200, joinDate: "2024-06-27", status: "active" },
  { id: "u-1062", name: "Zara Ahmed", email: "zara.ahmed@gmail.com", ticketsOwned: 988, totalSpentPence: 158400, joinDate: "2024-01-30", status: "vip" },
  { id: "u-1063", name: "Henry Foster", email: "henry.foster@yahoo.com", ticketsOwned: 33, totalSpentPence: 4100, joinDate: "2026-04-02", status: "active" },
  { id: "u-1064", name: "Lily Nguyen", email: "lily.nguyen@gmail.com", ticketsOwned: 712, totalSpentPence: 114800, joinDate: "2024-07-19", status: "active" },
  { id: "u-1065", name: "Jack Robinson", email: "jack.robinson@icloud.com", ticketsOwned: 0, totalSpentPence: 0, joinDate: "2026-06-20", status: "suspended" },
  { id: "u-1066", name: "Freya Andersson", email: "freya.andersson@gmail.com", ticketsOwned: 1345, totalSpentPence: 221000, joinDate: "2023-12-05", status: "vip" },
  { id: "u-1067", name: "Samuel Green", email: "samuel.green@outlook.com", ticketsOwned: 129, totalSpentPence: 17600, joinDate: "2025-05-21", status: "active" },
  { id: "u-1068", name: "Ruby Evans", email: "ruby.evans@gmail.com", ticketsOwned: 486, totalSpentPence: 79300, joinDate: "2024-09-17", status: "active" },
  { id: "u-1069", name: "Leo Martinez", email: "leo.martinez@hotmail.com", ticketsOwned: 67, totalSpentPence: 8500, joinDate: "2025-12-29", status: "active" },
  { id: "u-1070", name: "Hannah Scott", email: "hannah.scott@gmail.com", ticketsOwned: 2035, totalSpentPence: 333900, joinDate: "2023-06-08", status: "vip" },
  { id: "u-1071", name: "Adam Walsh", email: "adam.walsh@yahoo.com", ticketsOwned: 154, totalSpentPence: 21400, joinDate: "2025-04-11", status: "active" },
];

// ---------------------------------------------------------------------------
// Draws & Winners
// ---------------------------------------------------------------------------

export const adminDraws: AdminDraw[] = [
  { id: "d-13210", competition: "£50,000 Tax-Free Cash", winner: "Marcus Reid", prize: "£50,000 Cash", drawDate: "2026-06-28", ticketNumber: 148230, status: "completed" },
  { id: "d-13099", competition: "Toshiba 65\" 4K TV Bundle", winner: "Grace Kim", prize: "Toshiba 65\" 4K TV + Soundbar", drawDate: "2026-06-20", ticketNumber: 51902, status: "completed" },
  { id: "d-12988", competition: "BMW M4 Competition", winner: "Amara Nwosu", prize: "BMW M4 Competition xDrive", drawDate: "2026-06-13", ticketNumber: 298471, status: "completed" },
  { id: "d-12905", competition: "£25,000 Holiday Fund", winner: "Noah Patel", prize: "£25,000 Travel Voucher", drawDate: "2026-06-06", ticketNumber: 77310, status: "completed" },
  { id: "d-12844", competition: "Rolex Submariner Bundle", winner: "Hannah Scott", prize: "Rolex Submariner + £5k Cash", drawDate: "2026-05-30", ticketNumber: 19045, status: "completed" },
  { id: "d-12790", competition: "Ford Mustang Dark Horse", winner: "Isabella Rossi", prize: "Ford Mustang Dark Horse", drawDate: "2026-05-23", ticketNumber: 205118, status: "completed" },
  { id: "d-12722", competition: "MacBook Pro + iPhone Bundle", winner: "Daniel Cohen", prize: "MacBook Pro M4 + iPhone 17", drawDate: "2026-05-16", ticketNumber: 40233, status: "completed" },
  { id: "d-12655", competition: "£100,000 Life-Changing Cash", winner: "Freya Andersson", prize: "£100,000 Cash", drawDate: "2026-05-09", ticketNumber: 512884, status: "completed" },
  // Pending — ready for the admin to run the draw.
  { id: "d-13551", competition: "Summer Festival Instant Wins", winner: null, prize: "500+ Instant Prizes", drawDate: "2026-07-04", ticketNumber: null, status: "pending" },
  { id: "d-13344", competition: "£1.2M London Home in Zone 1", winner: null, prize: "£1,200,000 London Home", drawDate: "2026-07-06", ticketNumber: null, status: "pending" },
];
