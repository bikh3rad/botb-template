"use client"

import * as React from "react"
import { Pencil, Plus, Search, Trash2, Trophy } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
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
import {
  adminCompetitions,
  formatDate,
  formatNumber,
  formatPence,
  soldPercent,
} from "@/lib/admin-data"
import type { AdminCompetition, CompetitionStatus } from "@/types/admin"

type BadgeVariant = React.ComponentProps<typeof Badge>["variant"]

/**
 * Display metadata for each lifecycle status: a nicely-cased label and the
 * Badge variant used to colour it. Keeps table + filter labels consistent.
 */
const STATUS_META: Record<
  CompetitionStatus,
  { label: string; variant: BadgeVariant }
> = {
  live: { label: "Live", variant: "success" },
  "ending-soon": { label: "Ending soon", variant: "warning" },
  "sold-out": { label: "Sold out", variant: "secondary" },
  drawn: { label: "Drawn", variant: "muted" },
  draft: { label: "Draft", variant: "outline" },
}

/** Render order for status filter options and form select. */
const STATUS_ORDER: CompetitionStatus[] = [
  "live",
  "ending-soon",
  "sold-out",
  "drawn",
  "draft",
]

/** Editable fields captured by the create/edit dialog. */
interface CompetitionForm {
  title: string
  prize: string
  /** Ticket price in pounds, as typed (empty or "0" means FREE). */
  price: string
  /** Total tickets available, as typed. */
  total: string
  category: string
  status: CompetitionStatus
}

const EMPTY_FORM: CompetitionForm = {
  title: "",
  prize: "",
  price: "",
  total: "",
  category: "",
  status: "draft",
}

/** Build a form snapshot from an existing competition for editing. */
function toForm(competition: AdminCompetition): CompetitionForm {
  return {
    title: competition.title,
    prize: competition.prize,
    price:
      competition.ticketPricePence === 0
        ? ""
        : (competition.ticketPricePence / 100).toFixed(2),
    total: String(competition.ticketsTotal),
    category: competition.category,
    status: competition.status,
  }
}

/** Today's date as an ISO YYYY-MM-DD string, used as a draft's default close date. */
function todayIso(): string {
  return new Date().toISOString().slice(0, 10)
}

export function CompetitionsTable() {
  // Copy the seed data into local state so create/edit/delete mutate locally.
  const [competitions, setCompetitions] =
    React.useState<AdminCompetition[]>(adminCompetitions)
  const [query, setQuery] = React.useState("")
  const [statusFilter, setStatusFilter] = React.useState<"all" | CompetitionStatus>(
    "all"
  )

  // Create/edit dialog. `editingId === null` means we are creating a new row.
  const [formOpen, setFormOpen] = React.useState(false)
  const [editingId, setEditingId] = React.useState<string | null>(null)
  const [form, setForm] = React.useState<CompetitionForm>(EMPTY_FORM)

  // Delete confirmation dialog target.
  const [deleteTarget, setDeleteTarget] = React.useState<AdminCompetition | null>(
    null
  )

  const filtered = React.useMemo(() => {
    const term = query.trim().toLowerCase()
    return competitions.filter((competition) => {
      const matchesStatus =
        statusFilter === "all" || competition.status === statusFilter
      const matchesQuery =
        term === "" ||
        competition.title.toLowerCase().includes(term) ||
        competition.prize.toLowerCase().includes(term)
      return matchesStatus && matchesQuery
    })
  }, [competitions, query, statusFilter])

  function updateForm<K extends keyof CompetitionForm>(
    key: K,
    value: CompetitionForm[K]
  ) {
    setForm((prev) => ({ ...prev, [key]: value }))
  }

  function openCreate() {
    setEditingId(null)
    setForm(EMPTY_FORM)
    setFormOpen(true)
  }

  function openEdit(competition: AdminCompetition) {
    setEditingId(competition.id)
    setForm(toForm(competition))
    setFormOpen(true)
  }

  function handleSave(event: React.FormEvent) {
    event.preventDefault()

    const ticketPricePence = Math.max(
      0,
      Math.round((Number.parseFloat(form.price) || 0) * 100)
    )
    const ticketsTotal = Math.max(0, Number.parseInt(form.total, 10) || 0)

    if (editingId === null) {
      // Create: brand-new competition with a locally generated id.
      const created: AdminCompetition = {
        id: `c-${Date.now()}`,
        title: form.title.trim(),
        prize: form.prize.trim(),
        ticketPricePence,
        ticketsSold: 0,
        ticketsTotal,
        status: form.status,
        endDate: todayIso(),
        category: form.category.trim(),
      }
      setCompetitions((prev) => [created, ...prev])
    } else {
      // Edit: replace the matching row, preserving sales + close date.
      setCompetitions((prev) =>
        prev.map((competition) =>
          competition.id === editingId
            ? {
                ...competition,
                title: form.title.trim(),
                prize: form.prize.trim(),
                ticketPricePence,
                ticketsTotal,
                status: form.status,
                category: form.category.trim(),
              }
            : competition
        )
      )
    }

    setFormOpen(false)
  }

  function confirmDelete() {
    if (deleteTarget === null) return
    setCompetitions((prev) =>
      prev.filter((competition) => competition.id !== deleteTarget.id)
    )
    setDeleteTarget(null)
  }

  return (
    <div className="space-y-4">
      {/* Toolbar: search (left) + status filter and create action (right). */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="relative w-full sm:max-w-xs">
          <Search className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            type="search"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search title or prize…"
            className="pl-9"
            aria-label="Search competitions"
          />
        </div>

        <div className="flex items-center gap-2">
          <Select
            value={statusFilter}
            onChange={(event) =>
              setStatusFilter(event.target.value as "all" | CompetitionStatus)
            }
            className="sm:w-44"
            aria-label="Filter by status"
          >
            <option value="all">All statuses</option>
            {STATUS_ORDER.map((status) => (
              <option key={status} value={status}>
                {STATUS_META[status].label}
              </option>
            ))}
          </Select>

          <Button onClick={openCreate} className="shrink-0">
            <Plus data-icon="inline-start" />
            New Competition
          </Button>
        </div>
      </div>

      <p className="text-xs text-muted-foreground">
        Showing {formatNumber(filtered.length)} of{" "}
        {formatNumber(competitions.length)} competitions
      </p>

      <Card className="overflow-hidden py-0">
        <CardContent className="px-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Competition</TableHead>
                <TableHead>Prize</TableHead>
                <TableHead>Ticket Price</TableHead>
                <TableHead>Tickets Sold</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>End Date</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={7}
                    className="h-32 text-center text-sm text-muted-foreground"
                  >
                    <span className="inline-flex flex-col items-center gap-2">
                      <Trophy className="size-6 text-muted-foreground/60" />
                      No competitions match your filters.
                    </span>
                  </TableCell>
                </TableRow>
              ) : (
                filtered.map((competition) => {
                  const percent = soldPercent(
                    competition.ticketsSold,
                    competition.ticketsTotal
                  )
                  const meta = STATUS_META[competition.status]

                  return (
                    <TableRow key={competition.id}>
                      <TableCell>
                        <div className="font-medium">{competition.title}</div>
                        <div className="text-xs text-muted-foreground">
                          {competition.category}
                        </div>
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {competition.prize}
                      </TableCell>
                      <TableCell>
                        {competition.ticketPricePence === 0 ? (
                          <span className="font-medium text-botb-green">
                            FREE
                          </span>
                        ) : (
                          formatPence(competition.ticketPricePence)
                        )}
                      </TableCell>
                      <TableCell>
                        <div className="min-w-32 space-y-1.5">
                          <div className="text-xs text-muted-foreground">
                            {formatNumber(competition.ticketsSold)} /{" "}
                            {formatNumber(competition.ticketsTotal)}
                          </div>
                          <div
                            className="h-1.5 w-full overflow-hidden rounded-full bg-muted"
                            role="progressbar"
                            aria-valuenow={percent}
                            aria-valuemin={0}
                            aria-valuemax={100}
                            aria-label={`${percent}% of tickets sold`}
                          >
                            <div
                              className="h-full rounded-full bg-primary"
                              style={{ width: `${percent}%` }}
                            />
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant={meta.variant}>{meta.label}</Badge>
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {formatDate(competition.endDate)}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center justify-end gap-1">
                          <Button
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => openEdit(competition)}
                            aria-label={`Edit ${competition.title}`}
                          >
                            <Pencil />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => setDeleteTarget(competition)}
                            className="text-muted-foreground hover:text-destructive"
                            aria-label={`Delete ${competition.title}`}
                          >
                            <Trash2 />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  )
                })
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Create / edit dialog — reused for both flows. */}
      <Dialog open={formOpen} onOpenChange={setFormOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editingId === null
                ? "Create New Competition"
                : "Edit Competition"}
            </DialogTitle>
            <DialogDescription>
              {editingId === null
                ? "Add a new competition to the platform."
                : "Update the details for this competition."}
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleSave} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="competition-title">Title</Label>
              <Input
                id="competition-title"
                value={form.title}
                onChange={(event) => updateForm("title", event.target.value)}
                placeholder="Win a Defender D350 X"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="competition-prize">Prize</Label>
              <Input
                id="competition-prize"
                value={form.prize}
                onChange={(event) => updateForm("prize", event.target.value)}
                placeholder="Land Rover Defender D350 X"
                required
              />
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="competition-price">Ticket price (£)</Label>
                <Input
                  id="competition-price"
                  type="number"
                  min={0}
                  step={0.01}
                  inputMode="decimal"
                  value={form.price}
                  onChange={(event) => updateForm("price", event.target.value)}
                  placeholder="0.00 (FREE)"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="competition-total">Total tickets</Label>
                <Input
                  id="competition-total"
                  type="number"
                  min={0}
                  step={1}
                  inputMode="numeric"
                  value={form.total}
                  onChange={(event) => updateForm("total", event.target.value)}
                  placeholder="250000"
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="competition-category">Category</Label>
                <Input
                  id="competition-category"
                  value={form.category}
                  onChange={(event) =>
                    updateForm("category", event.target.value)
                  }
                  placeholder="Cars"
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="competition-status">Status</Label>
                <Select
                  id="competition-status"
                  value={form.status}
                  onChange={(event) =>
                    updateForm(
                      "status",
                      event.target.value as CompetitionStatus
                    )
                  }
                >
                  {STATUS_ORDER.map((status) => (
                    <option key={status} value={status}>
                      {STATUS_META[status].label}
                    </option>
                  ))}
                </Select>
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setFormOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit">
                {editingId === null ? "Create Competition" : "Save Changes"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Delete confirmation dialog. */}
      <Dialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null)
        }}
      >
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Delete competition?</DialogTitle>
            <DialogDescription>
              {deleteTarget
                ? `“${deleteTarget.title}” will be permanently removed. This can’t be undone.`
                : "This can’t be undone."}
            </DialogDescription>
          </DialogHeader>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setDeleteTarget(null)}
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={confirmDelete}
              className="gap-1.5"
            >
              <Trash2 data-icon="inline-start" />
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
