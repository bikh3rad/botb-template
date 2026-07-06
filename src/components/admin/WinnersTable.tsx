"use client"

import * as React from "react"
import {
  AlertTriangle,
  Ban,
  Loader2,
  Pencil,
  Play,
  Plus,
  ShieldAlert,
  Sparkles,
  Trophy,
} from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { apiGet, apiPost, apiPut, ApiError } from "@/lib/admin/client"
import { formatDateTime } from "@/lib/admin/format"
import type { Competition, Draw, DrawApiStatus } from "@/types/admin-api"

/** Response shape for the admin draws list endpoint. */
interface DrawsResponse {
  draws: Draw[]
  total: number
  count: number
  limit: number
  offset: number
}

/** Response shape for the competitions list endpoint. */
interface CompetitionsResponse {
  competitions: Competition[]
}

/** Display metadata for each draw status. */
const STATUS_META: Record<
  DrawApiStatus,
  { label: string; variant: React.ComponentProps<typeof Badge>["variant"] }
> = {
  pending: { label: "Pending", variant: "warning" },
  drawn: { label: "Drawn", variant: "success" },
  void: { label: "Void", variant: "muted" },
}

/** Which dialog (if any) is currently open, and for which draw. */
type DialogState =
  | { kind: "create" }
  | { kind: "edit"; draw: Draw }
  | { kind: "void"; draw: Draw }
  | { kind: "run"; draw: Draw }
  | { kind: "reassign"; draw: Draw }
  | null

export function WinnersTable() {
  const [draws, setDraws] = React.useState<Draw[]>([])
  const [competitions, setCompetitions] = React.useState<Competition[]>([])
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)

  const [dialog, setDialog] = React.useState<DialogState>(null)
  const [submitting, setSubmitting] = React.useState(false)
  const [formError, setFormError] = React.useState<string | null>(null)

  // Create-draw form fields.
  const [createCompetitionId, setCreateCompetitionId] = React.useState("")
  const [createPrize, setCreatePrize] = React.useState("")

  // Edit-prize form field.
  const [editPrize, setEditPrize] = React.useState("")

  // Void-draw form field.
  const [voidReason, setVoidReason] = React.useState("")

  // Reassign-winner form fields + two-step confirmation.
  const [reassignStep, setReassignStep] = React.useState<1 | 2>(1)
  const [reassignTicketId, setReassignTicketId] = React.useState("")
  const [reassignReason, setReassignReason] = React.useState("")

  const load = React.useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const [drawsRes, liveRes, closedRes] = await Promise.all([
        apiGet<DrawsResponse>("/apis/draw/v1/admin/draws?limit=100&offset=0"),
        apiGet<CompetitionsResponse>("/apis/competition/v1/competitions?status=live"),
        apiGet<CompetitionsResponse>("/apis/competition/v1/competitions?status=closed"),
      ])
      setDraws(drawsRes.draws ?? [])
      setCompetitions([...(liveRes.competitions ?? []), ...(closedRes.competitions ?? [])])
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setLoading(false)
    }
  }, [])

  React.useEffect(() => {
    void load()
  }, [load])

  const competitionTitle = React.useCallback(
    (id: string) => competitions.find((c) => c.id === id)?.title ?? id,
    [competitions]
  )

  const drawnCount = draws.filter((d) => d.status === "drawn").length
  const pendingCount = draws.filter((d) => d.status === "pending").length
  const voidCount = draws.filter((d) => d.status === "void").length

  function closeDialog() {
    setDialog(null)
    setFormError(null)
    setSubmitting(false)
    setCreateCompetitionId("")
    setCreatePrize("")
    setEditPrize("")
    setVoidReason("")
    setReassignStep(1)
    setReassignTicketId("")
    setReassignReason("")
  }

  function openCreate() {
    setCreateCompetitionId(competitions[0]?.id ?? "")
    setCreatePrize("")
    setFormError(null)
    setDialog({ kind: "create" })
  }

  function openEdit(draw: Draw) {
    setEditPrize(draw.prize)
    setFormError(null)
    setDialog({ kind: "edit", draw })
  }

  function openVoid(draw: Draw) {
    setVoidReason("")
    setFormError(null)
    setDialog({ kind: "void", draw })
  }

  function openRun(draw: Draw) {
    setFormError(null)
    setDialog({ kind: "run", draw })
  }

  function openReassign(draw: Draw) {
    setReassignStep(1)
    setReassignTicketId("")
    setReassignReason("")
    setFormError(null)
    setDialog({ kind: "reassign", draw })
  }

  async function handleCreate(event: React.FormEvent) {
    event.preventDefault()
    if (!createCompetitionId) return
    setSubmitting(true)
    setFormError(null)
    try {
      await apiPost("/apis/draw/v1/admin/draws", {
        competition_id: createCompetitionId,
        prize: createPrize.trim(),
      })
      closeDialog()
      await load()
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : "Something went wrong")
      setSubmitting(false)
    }
  }

  async function handleEditPrize(event: React.FormEvent) {
    event.preventDefault()
    if (dialog?.kind !== "edit") return
    setSubmitting(true)
    setFormError(null)
    try {
      await apiPut(`/apis/draw/v1/admin/draws/${dialog.draw.id}`, {
        prize: editPrize.trim(),
      })
      closeDialog()
      await load()
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : "Something went wrong")
      setSubmitting(false)
    }
  }

  async function handleVoid(event: React.FormEvent) {
    event.preventDefault()
    if (dialog?.kind !== "void" || voidReason.trim() === "") return
    setSubmitting(true)
    setFormError(null)
    try {
      await apiPost(`/apis/draw/v1/admin/draws/${dialog.draw.id}/void`, {
        reason: voidReason.trim(),
      })
      closeDialog()
      await load()
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : "Something went wrong")
      setSubmitting(false)
    }
  }

  async function handleRun() {
    if (dialog?.kind !== "run") return
    setSubmitting(true)
    setFormError(null)
    try {
      await apiPost(`/apis/draw/v1/admin/draws/${dialog.draw.id}/run`)
      closeDialog()
      await load()
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : "Something went wrong")
      setSubmitting(false)
    }
  }

  async function handleReassign(event: React.FormEvent) {
    event.preventDefault()
    if (
      dialog?.kind !== "reassign" ||
      reassignTicketId.trim() === "" ||
      reassignReason.trim() === ""
    )
      return
    setSubmitting(true)
    setFormError(null)
    try {
      await apiPost(`/apis/draw/v1/admin/draws/${dialog.draw.id}/reassign`, {
        winner_ticket_id: reassignTicketId.trim(),
        reason: reassignReason.trim(),
      })
      closeDialog()
      await load()
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : "Something went wrong")
      setSubmitting(false)
    }
  }

  return (
    <div className="space-y-6">
      {/* Summary chips + create action. */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex flex-wrap gap-3">
          <span className="inline-flex items-center gap-2 rounded-lg border border-border bg-card px-3 py-1.5 text-sm shadow-sm">
            <Trophy className="size-4 text-botb-green" />
            <span className="font-medium tabular-nums">{drawnCount}</span>
            <span className="text-muted-foreground">draws completed</span>
          </span>
          <span className="inline-flex items-center gap-2 rounded-lg border border-border bg-card px-3 py-1.5 text-sm shadow-sm">
            <Sparkles className="size-4 text-primary" />
            <span className="font-medium tabular-nums">{pendingCount}</span>
            <span className="text-muted-foreground">awaiting draw</span>
          </span>
          <span className="inline-flex items-center gap-2 rounded-lg border border-border bg-card px-3 py-1.5 text-sm shadow-sm">
            <Ban className="size-4 text-muted-foreground" />
            <span className="font-medium tabular-nums">{voidCount}</span>
            <span className="text-muted-foreground">voided</span>
          </span>
        </div>

        <Button onClick={openCreate} className="shrink-0">
          <Plus data-icon="inline-start" />
          Create draw
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Draw history</CardTitle>
        </CardHeader>

        <CardContent className="px-0">
          {loading ? (
            <div className="flex items-center justify-center gap-2 py-16 text-sm text-muted-foreground">
              <Loader2 className="size-4 animate-spin" />
              Loading draws…
            </div>
          ) : error ? (
            <div className="flex flex-col items-center gap-2 py-16 text-sm text-destructive">
              <AlertTriangle className="size-5" />
              {error}
            </div>
          ) : draws.length === 0 ? (
            <div className="flex flex-col items-center gap-2 py-16 text-sm text-muted-foreground">
              <Trophy className="size-6 text-muted-foreground/60" />
              No draws yet.
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="px-6">Competition</TableHead>
                  <TableHead>Prize</TableHead>
                  <TableHead>Winner</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Drawn At</TableHead>
                  <TableHead className="px-6 text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {draws.map((draw) => {
                  const meta = STATUS_META[draw.status]
                  return (
                    <TableRow key={draw.id}>
                      <TableCell className="px-6 font-medium">
                        {competitionTitle(draw.competition_id)}
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {draw.prize}
                      </TableCell>
                      <TableCell>
                        {draw.status === "drawn" ? (
                          <div className="space-y-0.5 font-mono text-xs text-muted-foreground">
                            <div>
                              {draw.winner_ticket_id
                                ? `Ticket ${draw.winner_ticket_id}`
                                : "—"}
                            </div>
                            <div>
                              {draw.winner_user_id
                                ? `User ${draw.winner_user_id}`
                                : "—"}
                            </div>
                          </div>
                        ) : (
                          <span className="text-sm text-muted-foreground">—</span>
                        )}
                      </TableCell>
                      <TableCell>
                        <Badge variant={meta.variant}>{meta.label}</Badge>
                        {draw.status === "void" && draw.void_reason && (
                          <p className="mt-1 max-w-48 text-xs text-muted-foreground">
                            {draw.void_reason}
                          </p>
                        )}
                      </TableCell>
                      <TableCell className="text-right text-muted-foreground tabular-nums">
                        {formatDateTime(draw.drawn_at)}
                      </TableCell>
                      <TableCell className="px-6 text-right">
                        <div className="flex items-center justify-end gap-1">
                          {draw.status === "pending" && (
                            <>
                              <Button
                                size="sm"
                                onClick={() => openRun(draw)}
                                aria-label={`Run draw for ${competitionTitle(draw.competition_id)}`}
                              >
                                <Play data-icon="inline-start" />
                                Run
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon-sm"
                                onClick={() => openEdit(draw)}
                                aria-label={`Edit prize for ${competitionTitle(draw.competition_id)}`}
                              >
                                <Pencil />
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon-sm"
                                onClick={() => openVoid(draw)}
                                className="text-muted-foreground hover:text-destructive"
                                aria-label={`Void draw for ${competitionTitle(draw.competition_id)}`}
                              >
                                <Ban />
                              </Button>
                            </>
                          )}
                          {draw.status === "drawn" && (
                            <>
                              <Button
                                variant="ghost"
                                size="icon-sm"
                                onClick={() => openEdit(draw)}
                                aria-label={`Edit prize for ${competitionTitle(draw.competition_id)}`}
                              >
                                <Pencil />
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon-sm"
                                onClick={() => openVoid(draw)}
                                className="text-muted-foreground hover:text-destructive"
                                aria-label={`Void draw for ${competitionTitle(draw.competition_id)}`}
                              >
                                <Ban />
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon-sm"
                                onClick={() => openReassign(draw)}
                                className="text-muted-foreground hover:text-destructive"
                                aria-label={`Reassign winner for ${competitionTitle(draw.competition_id)}`}
                              >
                                <ShieldAlert />
                              </Button>
                            </>
                          )}
                          {draw.status === "void" && (
                            <span className="text-sm text-muted-foreground">—</span>
                          )}
                        </div>
                      </TableCell>
                    </TableRow>
                  )
                })}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {/* Create draw dialog. */}
      <Dialog
        open={dialog?.kind === "create"}
        onOpenChange={(open) => !open && closeDialog()}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create draw</DialogTitle>
            <DialogDescription>
              Set up a new pending draw for a competition. Run it once you&apos;re
              ready to pick a winner.
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleCreate} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="draw-competition">Competition</Label>
              <Select
                id="draw-competition"
                value={createCompetitionId}
                onChange={(event) => setCreateCompetitionId(event.target.value)}
                required
              >
                {competitions.length === 0 && <option value="">No competitions</option>}
                {competitions.map((competition) => (
                  <option key={competition.id} value={competition.id}>
                    {competition.title}
                  </option>
                ))}
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="draw-prize">Prize</Label>
              <Input
                id="draw-prize"
                value={createPrize}
                onChange={(event) => setCreatePrize(event.target.value)}
                placeholder="Land Rover Defender D350 X"
                required
              />
            </div>

            {formError && (
              <p className="text-sm text-destructive">{formError}</p>
            )}

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={closeDialog}
                disabled={submitting}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={submitting || !createCompetitionId}>
                {submitting ? <Loader2 className="animate-spin" /> : <Plus />}
                Create draw
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Edit prize dialog. */}
      <Dialog
        open={dialog?.kind === "edit"}
        onOpenChange={(open) => !open && closeDialog()}
      >
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Edit prize</DialogTitle>
            <DialogDescription>
              Update the prize description for this draw.
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleEditPrize} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="edit-prize">Prize</Label>
              <Input
                id="edit-prize"
                value={editPrize}
                onChange={(event) => setEditPrize(event.target.value)}
                required
              />
            </div>

            {formError && (
              <p className="text-sm text-destructive">{formError}</p>
            )}

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={closeDialog}
                disabled={submitting}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={submitting}>
                {submitting && <Loader2 className="animate-spin" />}
                Save changes
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Void draw dialog. */}
      <Dialog
        open={dialog?.kind === "void"}
        onOpenChange={(open) => !open && closeDialog()}
      >
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Void draw?</DialogTitle>
            <DialogDescription>
              This freezes the draw permanently — it can no longer be run or
              edited. To re-draw, void this draw then create a new one.
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleVoid} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="void-reason">Reason</Label>
              <textarea
                id="void-reason"
                value={voidReason}
                onChange={(event) => setVoidReason(event.target.value)}
                placeholder="Explain why this draw is being voided…"
                className="flex min-h-20 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 dark:bg-input/30"
                required
              />
            </div>

            {formError && (
              <p className="text-sm text-destructive">{formError}</p>
            )}

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={closeDialog}
                disabled={submitting}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                variant="destructive"
                disabled={submitting || voidReason.trim() === ""}
              >
                {submitting && <Loader2 className="animate-spin" />}
                Void draw
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Run draw confirmation dialog. */}
      <Dialog
        open={dialog?.kind === "run"}
        onOpenChange={(open) => !open && closeDialog()}
      >
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <div className="mb-1 flex size-11 items-center justify-center rounded-full bg-primary/12 text-primary">
              <Trophy className="size-5" />
            </div>
            <DialogTitle>Run this draw?</DialogTitle>
            <DialogDescription>
              A random winning ticket will be selected
              {dialog?.kind === "run"
                ? ` for “${competitionTitle(dialog.draw.competition_id)}”`
                : ""}
              . This can&apos;t be undone.
            </DialogDescription>
          </DialogHeader>

          {formError && (
            <p className="text-sm text-destructive">{formError}</p>
          )}

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={closeDialog}
              disabled={submitting}
            >
              Cancel
            </Button>
            <Button onClick={handleRun} disabled={submitting}>
              {submitting ? (
                <Loader2 className="animate-spin" />
              ) : (
                <Play data-icon="inline-start" />
              )}
              Run draw
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Reassign winner dialog — two-step confirmation. */}
      <Dialog
        open={dialog?.kind === "reassign"}
        onOpenChange={(open) => !open && closeDialog()}
      >
        <DialogContent className="max-w-sm">
          {reassignStep === 1 ? (
            <>
              <DialogHeader>
                <div className="mb-1 flex size-11 items-center justify-center rounded-full bg-destructive/12 text-destructive">
                  <ShieldAlert className="size-5" />
                </div>
                <DialogTitle>Reassign winner?</DialogTitle>
                <DialogDescription>
                  This directly overrides the recorded winner for this draw and
                  is audited. Prefer voiding and re-drawing when possible — only
                  continue if you need to correct this specific draw&apos;s
                  winner.
                </DialogDescription>
              </DialogHeader>

              <DialogFooter>
                <Button type="button" variant="outline" onClick={closeDialog}>
                  Cancel
                </Button>
                <Button
                  type="button"
                  variant="destructive"
                  onClick={() => setReassignStep(2)}
                >
                  Continue
                </Button>
              </DialogFooter>
            </>
          ) : (
            <>
              <DialogHeader>
                <DialogTitle>Reassign winner</DialogTitle>
                <DialogDescription>
                  Provide the new winning ticket and a reason for the audit
                  log.
                </DialogDescription>
              </DialogHeader>

              <form onSubmit={handleReassign} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="reassign-ticket">Winning ticket ID</Label>
                  <Input
                    id="reassign-ticket"
                    value={reassignTicketId}
                    onChange={(event) => setReassignTicketId(event.target.value)}
                    placeholder="Ticket ID"
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="reassign-reason">Reason</Label>
                  <textarea
                    id="reassign-reason"
                    value={reassignReason}
                    onChange={(event) => setReassignReason(event.target.value)}
                    placeholder="Explain why the winner is being reassigned…"
                    className="flex min-h-20 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 dark:bg-input/30"
                    required
                  />
                </div>

                {formError && (
                  <p className="text-sm text-destructive">{formError}</p>
                )}

                <DialogFooter>
                  <Button
                    type="button"
                    variant="outline"
                    onClick={closeDialog}
                    disabled={submitting}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    variant="destructive"
                    disabled={
                      submitting ||
                      reassignTicketId.trim() === "" ||
                      reassignReason.trim() === ""
                    }
                  >
                    {submitting && <Loader2 className="animate-spin" />}
                    Reassign winner
                  </Button>
                </DialogFooter>
              </form>
            </>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}
