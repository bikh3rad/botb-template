"use client"

import * as React from "react"
import {
  Calendar,
  Check,
  Loader2,
  PartyPopper,
  Play,
  Sparkles,
  Ticket,
  Trophy,
} from "lucide-react"

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { adminDraws, formatDate, formatNumber } from "@/lib/admin-data"
import { cn } from "@/lib/utils"
import type { AdminDraw } from "@/types/admin"

/**
 * Realistic winner names used to simulate a live draw. Kept intentionally
 * distinct from the seeded draw data so a "run draw" result reads as fresh.
 */
const WINNER_POOL = [
  "Charlotte Hayes",
  "Ryan O'Sullivan",
  "Aisha Rahman",
  "Daniel Whitfield",
  "Elena Petrova",
  "Marcus Boateng",
  "Sofia Delgado",
  "Nathan Brooks",
  "Priya Kapoor",
  "Callum Fraser",
] as const

/** How long the simulated draw "spins" before revealing a winner. */
const DRAW_DURATION_MS = 1200

/** How long a freshly-drawn row stays highlighted after the dialog closes. */
const HIGHLIGHT_DURATION_MS = 4000

/** The stage of the run-draw dialog flow. */
type DrawPhase = "confirm" | "drawing" | "success"

/** The outcome produced by a simulated draw. */
interface DrawResult {
  winner: string
  ticketNumber: number
}

/** Derive up to two uppercase initials from a full name, e.g. "Marcus Reid" -> "MR". */
function getInitials(name: string): string {
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("")
}

/**
 * Order draws so the actionable pending draws surface first, then completed
 * draws by most-recent date. Sorting a copy keeps the source order intact.
 */
function sortDraws(draws: AdminDraw[]): AdminDraw[] {
  return [...draws].sort((a, b) => {
    if (a.status !== b.status) return a.status === "pending" ? -1 : 1
    return b.drawDate.localeCompare(a.drawDate)
  })
}

export function WinnersTable() {
  // Local copy so simulated draws persist for the session without an API.
  const [draws, setDraws] = React.useState<AdminDraw[]>(() => sortDraws(adminDraws))

  // Dialog flow state. `activeDrawId` doubles as the open/closed signal.
  const [activeDrawId, setActiveDrawId] = React.useState<string | null>(null)
  const [phase, setPhase] = React.useState<DrawPhase>("confirm")
  const [result, setResult] = React.useState<DrawResult | null>(null)
  const [recentlyDrawnId, setRecentlyDrawnId] = React.useState<string | null>(null)

  const activeDraw = React.useMemo(
    () => draws.find((draw) => draw.id === activeDrawId) ?? null,
    [draws, activeDrawId]
  )

  const completedCount = draws.filter((draw) => draw.status === "completed").length
  const pendingCount = draws.filter((draw) => draw.status === "pending").length
  const pendingDraws = draws.filter((draw) => draw.status === "pending")

  // Clear the row highlight a short while after a draw completes.
  React.useEffect(() => {
    if (!recentlyDrawnId) return
    const timer = window.setTimeout(
      () => setRecentlyDrawnId(null),
      HIGHLIGHT_DURATION_MS
    )
    return () => window.clearTimeout(timer)
  }, [recentlyDrawnId])

  function openRunDraw(id: string) {
    setActiveDrawId(id)
    setPhase("confirm")
    setResult(null)
  }

  function closeDialog() {
    setActiveDrawId(null)
    setPhase("confirm")
    setResult(null)
  }

  // Base-ui requests a close on backdrop/escape; block it mid-draw so the
  // simulation can't be interrupted, otherwise reset the flow.
  function handleOpenChange(open: boolean) {
    if (open || phase === "drawing") return
    closeDialog()
  }

  function confirmDraw() {
    if (!activeDraw) return
    const drawId = activeDraw.id
    setPhase("drawing")

    // Randomness lives only inside the handler/timeout — never during render —
    // so server and client markup stay identical and hydration is safe.
    window.setTimeout(() => {
      const winner = WINNER_POOL[Math.floor(Math.random() * WINNER_POOL.length)]
      const ticketNumber = Math.floor(Math.random() * 900000) + 1

      setDraws((prev) =>
        sortDraws(
          prev.map((draw) =>
            draw.id === drawId
              ? { ...draw, winner, ticketNumber, status: "completed" }
              : draw
          )
        )
      )
      setResult({ winner, ticketNumber })
      setRecentlyDrawnId(drawId)
      setPhase("success")
    }, DRAW_DURATION_MS)
  }

  return (
    <div className="space-y-6">
      {/* Summary chips — a quick read on draw throughput. */}
      <div className="flex flex-wrap gap-3">
        <span className="inline-flex items-center gap-2 rounded-lg border border-border bg-card px-3 py-1.5 text-sm shadow-sm">
          <Trophy className="size-4 text-botb-green" />
          <span className="font-medium tabular-nums">{completedCount}</span>
          <span className="text-muted-foreground">draws completed</span>
        </span>
        <span className="inline-flex items-center gap-2 rounded-lg border border-border bg-card px-3 py-1.5 text-sm shadow-sm">
          <Sparkles className="size-4 text-primary" />
          <span className="font-medium tabular-nums">{pendingCount}</span>
          <span className="text-muted-foreground">awaiting draw</span>
        </span>
      </div>

      {/* Pending draws surfaced as highlighted, actionable cards. */}
      {pendingDraws.length > 0 && (
        <section aria-label="Draws ready to run" className="space-y-3">
          <h3 className="font-heading text-sm font-semibold tracking-tight">
            Ready to draw
          </h3>
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            {pendingDraws.map((draw) => (
              <Card
                key={draw.id}
                className="gap-4 border-primary/30 bg-primary/[0.04] py-4"
              >
                <CardContent className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                  <div className="min-w-0 space-y-1">
                    <div className="flex items-center gap-2">
                      <Badge variant="warning">Pending</Badge>
                      <span className="inline-flex items-center gap-1 text-xs text-muted-foreground">
                        <Calendar className="size-3.5" />
                        {formatDate(draw.drawDate)}
                      </span>
                    </div>
                    <p className="truncate font-medium">{draw.competition}</p>
                    <p className="truncate text-sm text-muted-foreground">
                      {draw.prize}
                    </p>
                  </div>
                  <Button
                    onClick={() => openRunDraw(draw.id)}
                    className="shrink-0"
                    aria-label={`Run draw for ${draw.competition}`}
                  >
                    <Play />
                    Run draw
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>
      )}

      {/* Full draw history. */}
      <Card>
        <CardHeader>
          <CardTitle>Draw history</CardTitle>
        </CardHeader>

        <CardContent className="px-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="px-6">Competition</TableHead>
                <TableHead>Winner</TableHead>
                <TableHead>Prize</TableHead>
                <TableHead className="text-right">Winning Ticket</TableHead>
                <TableHead className="text-right">Draw Date</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="px-6 text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {draws.map((draw) => {
                const isCompleted = draw.status === "completed"
                const isHighlighted = draw.id === recentlyDrawnId
                return (
                  <TableRow
                    key={draw.id}
                    className={cn(
                      isHighlighted && "bg-botb-green/[0.07] hover:bg-botb-green/10"
                    )}
                  >
                    <TableCell className="px-6 font-medium">
                      {draw.competition}
                    </TableCell>
                    <TableCell>
                      {isCompleted && draw.winner ? (
                        <div className="flex items-center gap-3">
                          <Avatar className="bg-primary/10">
                            <AvatarFallback className="bg-transparent text-botb-orange-hover dark:text-primary">
                              {getInitials(draw.winner)}
                            </AvatarFallback>
                          </Avatar>
                          <span className="font-medium">{draw.winner}</span>
                        </div>
                      ) : (
                        <span className="text-sm text-muted-foreground">
                          — Not drawn —
                        </span>
                      )}
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {draw.prize}
                    </TableCell>
                    <TableCell className="text-right tabular-nums">
                      {draw.ticketNumber !== null ? (
                        <span className="inline-flex items-center justify-end gap-1.5">
                          <Ticket className="size-3.5 text-muted-foreground" />
                          #{formatNumber(draw.ticketNumber)}
                        </span>
                      ) : (
                        <span className="text-muted-foreground">—</span>
                      )}
                    </TableCell>
                    <TableCell className="text-right text-muted-foreground tabular-nums">
                      {formatDate(draw.drawDate)}
                    </TableCell>
                    <TableCell>
                      {isCompleted ? (
                        <Badge variant="success">Completed</Badge>
                      ) : (
                        <Badge variant="warning">Pending</Badge>
                      )}
                    </TableCell>
                    <TableCell className="px-6 text-right">
                      {isCompleted ? (
                        <span className="text-sm text-muted-foreground">—</span>
                      ) : (
                        <Button
                          size="sm"
                          onClick={() => openRunDraw(draw.id)}
                          aria-label={`Run draw for ${draw.competition}`}
                        >
                          <Play />
                          Run draw
                        </Button>
                      )}
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Controlled run-draw dialog: confirm → drawing → success. */}
      <Dialog open={activeDrawId !== null} onOpenChange={handleOpenChange}>
        <DialogContent showCloseButton={phase !== "drawing"}>
          {phase === "success" && result ? (
            <>
              <DialogHeader>
                <div className="mb-1 flex size-11 items-center justify-center rounded-full bg-botb-green/12 text-botb-green">
                  <PartyPopper className="size-5" />
                </div>
                <DialogTitle>Winner drawn!</DialogTitle>
                <DialogDescription>
                  {activeDraw
                    ? `The draw for “${activeDraw.competition}” is complete.`
                    : "The draw is complete."}
                </DialogDescription>
              </DialogHeader>
              <div className="rounded-lg border border-border bg-muted/50 p-4">
                <div className="flex items-center gap-3">
                  <Avatar className="size-10 bg-primary/10">
                    <AvatarFallback className="bg-transparent text-botb-orange-hover dark:text-primary">
                      {getInitials(result.winner)}
                    </AvatarFallback>
                  </Avatar>
                  <div className="min-w-0">
                    <p className="font-medium">{result.winner}</p>
                    <p className="inline-flex items-center gap-1.5 text-sm text-muted-foreground tabular-nums">
                      <Ticket className="size-3.5" />
                      Ticket #{formatNumber(result.ticketNumber)}
                    </p>
                  </div>
                </div>
              </div>
              <DialogFooter>
                <Button onClick={closeDialog}>
                  <Check />
                  Done
                </Button>
              </DialogFooter>
            </>
          ) : (
            <>
              <DialogHeader>
                <div className="mb-1 flex size-11 items-center justify-center rounded-full bg-primary/12 text-primary">
                  {phase === "drawing" ? (
                    <Loader2 className="size-5 animate-spin" />
                  ) : (
                    <Trophy className="size-5" />
                  )}
                </div>
                <DialogTitle>
                  {phase === "drawing" ? "Drawing a winner…" : "Run this draw?"}
                </DialogTitle>
                <DialogDescription>
                  {phase === "drawing"
                    ? "Selecting a random winning ticket. This only takes a moment."
                    : activeDraw
                      ? `A random winning ticket will be selected for “${activeDraw.competition}”. This can’t be undone.`
                      : "A random winning ticket will be selected. This can’t be undone."}
                </DialogDescription>
              </DialogHeader>
              {activeDraw && (
                <div className="rounded-lg border border-border bg-muted/50 p-4 text-sm">
                  <div className="flex items-center justify-between gap-3">
                    <span className="text-muted-foreground">Prize</span>
                    <span className="font-medium">{activeDraw.prize}</span>
                  </div>
                  <div className="mt-2 flex items-center justify-between gap-3">
                    <span className="text-muted-foreground">Draw date</span>
                    <span className="font-medium tabular-nums">
                      {formatDate(activeDraw.drawDate)}
                    </span>
                  </div>
                </div>
              )}
              <DialogFooter>
                <Button
                  variant="outline"
                  onClick={closeDialog}
                  disabled={phase === "drawing"}
                >
                  Cancel
                </Button>
                <Button onClick={confirmDraw} disabled={phase === "drawing"}>
                  {phase === "drawing" ? (
                    <>
                      <Loader2 className="animate-spin" />
                      Drawing…
                    </>
                  ) : (
                    <>
                      <Play />
                      Run draw
                    </>
                  )}
                </Button>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}
