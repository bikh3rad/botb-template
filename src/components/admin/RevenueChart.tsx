import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card"
import { formatCurrency, revenueSeries } from "@/lib/admin-data"
import type { RevenuePoint } from "@/types/admin"

// Chart geometry, expressed in viewBox units. The SVG scales fluidly to its
// container via `w-full h-auto`; month labels are rendered as HTML below so
// they stay crisp and evenly sized at every breakpoint.
const VIEW_W = 720
const VIEW_H = 300
const PAD_X = 8
const PAD_TOP = 16
const PAD_BOTTOM = 12
const INNER_W = VIEW_W - PAD_X * 2
const INNER_H = VIEW_H - PAD_TOP - PAD_BOTTOM
const GRID_LINES = 4

interface PlotPoint {
  x: number
  y: number
  point: RevenuePoint
}

/**
 * Project the revenue series into SVG coordinates. Values are mapped into a
 * padded range so the line sits comfortably within the plot rather than
 * clipping the top or hugging the baseline.
 */
function buildPlot(series: RevenuePoint[]): {
  points: PlotPoint[]
  baseline: number
} {
  const values = series.map((p) => p.revenue)
  const dataMin = Math.min(...values)
  const dataMax = Math.max(...values)
  const range = dataMax - dataMin || 1
  const ceil = dataMax + range * 0.15
  const floor = dataMin - range * 0.35
  const step = INNER_W / series.length

  const points = series.map((point, i) => {
    const x = PAD_X + step * (i + 0.5)
    const ratio = (point.revenue - floor) / (ceil - floor)
    const y = PAD_TOP + INNER_H * (1 - ratio)
    return { x, y, point }
  })

  return { points, baseline: PAD_TOP + INNER_H }
}

export function RevenueChart() {
  const { points, baseline } = buildPlot(revenueSeries)

  const linePath = points
    .map((p, i) => `${i === 0 ? "M" : "L"} ${p.x} ${p.y}`)
    .join(" ")

  const areaPath = [
    `M ${points[0].x} ${baseline}`,
    ...points.map((p) => `L ${p.x} ${p.y}`),
    `L ${points[points.length - 1].x} ${baseline}`,
    "Z",
  ].join(" ")

  const latest = revenueSeries[revenueSeries.length - 1]

  return (
    <Card>
      <CardHeader>
        <CardTitle>Revenue</CardTitle>
        <CardDescription>Last 12 months</CardDescription>
      </CardHeader>
      <CardContent className="space-y-3">
        <p className="font-heading text-2xl font-semibold tracking-tight">
          {formatCurrency(latest.revenue)}
          <span className="ml-2 align-middle text-xs font-normal text-muted-foreground">
            in {latest.month}
          </span>
        </p>

        <svg
          viewBox={`0 0 ${VIEW_W} ${VIEW_H}`}
          className="h-auto w-full text-primary"
          role="img"
          aria-label="Line chart of monthly revenue over the last 12 months"
        >
          <defs>
            <linearGradient id="revenue-fill" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="currentColor" stopOpacity="0.22" />
              <stop offset="100%" stopColor="currentColor" stopOpacity="0" />
            </linearGradient>
          </defs>

          {/* Horizontal gridlines for vertical reference. */}
          {Array.from({ length: GRID_LINES + 1 }, (_, i) => {
            const y = PAD_TOP + (INNER_H / GRID_LINES) * i
            return (
              <line
                key={i}
                x1={PAD_X}
                x2={VIEW_W - PAD_X}
                y1={y}
                y2={y}
                className="stroke-border"
                strokeWidth={1}
                vectorEffect="non-scaling-stroke"
              />
            )
          })}

          <path d={areaPath} fill="url(#revenue-fill)" />
          <path
            d={linePath}
            fill="none"
            stroke="currentColor"
            strokeWidth={2.5}
            strokeLinecap="round"
            strokeLinejoin="round"
            vectorEffect="non-scaling-stroke"
          />

          {points.map((p, i) => (
            <circle
              key={p.point.month}
              cx={p.x}
              cy={p.y}
              r={i === points.length - 1 ? 5 : 3}
              className="fill-primary"
              stroke="var(--card)"
              strokeWidth={2}
              vectorEffect="non-scaling-stroke"
            />
          ))}
        </svg>

        {/* Month axis labels — column-aligned with the plotted points. */}
        <div className="flex px-2" aria-hidden>
          {revenueSeries.map((point) => (
            <span
              key={point.month}
              className="flex-1 text-center text-[10px] font-medium text-muted-foreground sm:text-xs"
            >
              {point.month}
            </span>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
