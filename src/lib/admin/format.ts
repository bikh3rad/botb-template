// Display + conversion helpers shared across the admin editor UIs. Pure
// functions only — safe to import anywhere without pulling in mock data.

const gbp = new Intl.NumberFormat("en-GB", {
  style: "currency",
  currency: "GBP",
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
})

const num = new Intl.NumberFormat("en-GB")

/** Format whole pence into a GBP string, e.g. 125 -> "£1.25". */
export function formatPence(pence: number): string {
  return gbp.format((pence || 0) / 100)
}

/** Format a number with thousands separators, e.g. 12480 -> "12,480". */
export function formatNumber(value: number): string {
  return num.format(value || 0)
}

/** Parse a pounds string (e.g. "1.25") into whole pence, clamped to >= 0. */
export function poundsToPence(input: string): number {
  const n = Number.parseFloat(input)
  if (!Number.isFinite(n)) return 0
  return Math.max(0, Math.round(n * 100))
}

/** Render whole pence as an editable pounds string, e.g. 125 -> "1.25". */
export function penceToPounds(pence: number): string {
  if (!pence) return ""
  return (pence / 100).toFixed(2)
}

/** Convert an RFC3339 string to a value for <input type="datetime-local">. */
export function rfc3339ToLocalInput(iso: string): string {
  if (!iso) return ""
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ""
  const pad = (n: number) => String(n).padStart(2, "0")
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(
    d.getHours(),
  )}:${pad(d.getMinutes())}`
}

/** Convert a datetime-local input value back to an RFC3339 string ("" stays ""). */
export function localInputToRfc3339(value: string): string {
  if (!value) return ""
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ""
  return d.toISOString()
}

/** Format an RFC3339 timestamp for display, e.g. "6 Jul 2026, 14:30". */
export function formatDateTime(iso?: string | null): string {
  if (!iso) return "—"
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return "—"
  return d.toLocaleString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  })
}

/** Format an RFC3339 timestamp as a date only, e.g. "6 Jul 2026". */
export function formatDate(iso?: string | null): string {
  if (!iso) return "—"
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return "—"
  return d.toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
  })
}

/** Slugify a title into a URL-safe slug, e.g. "Win a Car!" -> "win-a-car". */
export function slugify(input: string): string {
  return input
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "")
}

/** Percentage of tickets sold, clamped to 0–100. */
export function soldPercent(sold: number, total: number): number {
  if (total <= 0) return 0
  return Math.min(100, Math.round((sold / total) * 100))
}

/** Two uppercase initials from a name, e.g. "Ada Lovelace" -> "AL". */
export function initials(name: string): string {
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("")
}
