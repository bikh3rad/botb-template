"use client"

import * as React from "react"
import {
  AlertCircle,
  Loader2,
  Pencil,
  Plus,
  RefreshCw,
  Search,
  Trash2,
  Trophy,
} from "lucide-react"

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
import { CompetitionMediaManager } from "@/components/admin/CompetitionMediaManager"
import { apiDelete, apiGet, apiPost, apiPut, ApiError } from "@/lib/admin/client"
import {
  formatNumber,
  formatPence,
  localInputToRfc3339,
  penceToPounds,
  poundsToPence,
  rfc3339ToLocalInput,
  slugify,
  soldPercent,
} from "@/lib/admin/format"
import type {
  Category,
  Competition,
  CompetitionApiStatus,
  CompetitionInput,
} from "@/types/admin-api"

type BadgeVariant = React.ComponentProps<typeof Badge>["variant"]

/** Display metadata for each lifecycle status: label + Badge variant. */
const STATUS_META: Record<CompetitionApiStatus, { label: string; variant: BadgeVariant }> = {
  draft: { label: "Draft", variant: "outline" },
  live: { label: "Live", variant: "success" },
  closed: { label: "Closed", variant: "muted" },
}

/** Render order for the status filter dropdown. */
const STATUS_FILTER_ORDER: CompetitionApiStatus[] = ["draft", "live", "closed"]

/** Valid forward-only status transitions available from the current status. */
function statusOptionsFor(current: CompetitionApiStatus | null): CompetitionApiStatus[] {
  if (current === null) return ["draft", "live"]
  if (current === "draft") return ["draft", "live"]
  if (current === "live") return ["live", "closed"]
  return ["closed"]
}

/** Editable fields captured by the create/edit dialog (all strings for inputs). */
interface CompetitionForm {
  title: string
  slug: string
  description: string
  prize: string
  /** Ticket price in pounds, as typed (empty or "0" means FREE). */
  price: string
  /** Total tickets available, as typed. */
  total: string
  categoryId: string
  status: CompetitionApiStatus
  startsAt: string
  endsAt: string
}

const EMPTY_FORM: CompetitionForm = {
  title: "",
  slug: "",
  description: "",
  prize: "",
  price: "",
  total: "",
  categoryId: "",
  status: "draft",
  startsAt: "",
  endsAt: "",
}

/** Build a form snapshot from an existing competition for editing. */
function toForm(competition: Competition): CompetitionForm {
  return {
    title: competition.title,
    slug: competition.slug,
    description: competition.description,
    prize: competition.prize,
    price: penceToPounds(competition.ticket_price_pence),
    total: String(competition.tickets_total),
    categoryId: competition.category_id ?? "",
    status: competition.status,
    startsAt: rfc3339ToLocalInput(competition.starts_at),
    endsAt: rfc3339ToLocalInput(competition.ends_at),
  }
}

/** Format an end date for the table, guarding against empty/invalid values. */
function formatEndDate(iso: string): string {
  if (!iso) return "—"
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return "—"
  return d.toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
  })
}

interface CompetitionsListResponse {
  competitions: Competition[]
}

interface CategoriesListResponse {
  categories: Category[]
}

export function CompetitionsTable() {
  const [competitions, setCompetitions] = React.useState<Competition[]>([])
  const [categories, setCategories] = React.useState<Category[]>([])
  const [loading, setLoading] = React.useState(true)
  const [loadError, setLoadError] = React.useState<string | null>(null)

  const [query, setQuery] = React.useState("")
  const [statusFilter, setStatusFilter] = React.useState<"all" | CompetitionApiStatus>("all")

  // Create/edit dialog. `editingId === null` means we are creating a new row.
  const [formOpen, setFormOpen] = React.useState(false)
  const [editingId, setEditingId] = React.useState<string | null>(null)
  const [editingTicketsSold, setEditingTicketsSold] = React.useState(0)
  const [statusOptions, setStatusOptions] = React.useState<CompetitionApiStatus[]>(
    statusOptionsFor(null)
  )
  const [form, setForm] = React.useState<CompetitionForm>(EMPTY_FORM)
  const [slugTouched, setSlugTouched] = React.useState(false)
  const [formError, setFormError] = React.useState<string | null>(null)
  const [submitting, setSubmitting] = React.useState(false)

  // Delete confirmation dialog target.
  const [deleteTarget, setDeleteTarget] = React.useState<Competition | null>(null)
  const [deleteError, setDeleteError] = React.useState<string | null>(null)
  // Typed-title confirmation guard, and the "has entrants" (409) blocked state.
  const [deleteConfirmText, setDeleteConfirmText] = React.useState("")
  const [deleteBlocked, setDeleteBlocked] = React.useState(false)
  const [deleting, setDeleting] = React.useState(false)

  async function load() {
    setLoading(true)
    setLoadError(null)
    try {
      const [draftRes, liveRes, closedRes, categoriesRes] = await Promise.all([
        apiGet<CompetitionsListResponse>("/apis/competition/v1/competitions?status=draft"),
        apiGet<CompetitionsListResponse>("/apis/competition/v1/competitions?status=live"),
        apiGet<CompetitionsListResponse>("/apis/competition/v1/competitions?status=closed"),
        apiGet<CategoriesListResponse>("/apis/competition/v1/categories"),
      ])
      setCompetitions([
        ...draftRes.competitions,
        ...liveRes.competitions,
        ...closedRes.competitions,
      ])
      setCategories(categoriesRes.categories)
    } catch (err) {
      setLoadError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setLoading(false)
    }
  }

  React.useEffect(() => {
    void load()
  }, [])

  const filtered = React.useMemo(() => {
    const term = query.trim().toLowerCase()
    return competitions.filter((competition) => {
      const matchesStatus = statusFilter === "all" || competition.status === statusFilter
      const matchesQuery =
        term === "" ||
        competition.title.toLowerCase().includes(term) ||
        competition.prize.toLowerCase().includes(term)
      return matchesStatus && matchesQuery
    })
  }, [competitions, query, statusFilter])

  function updateForm<K extends keyof CompetitionForm>(key: K, value: CompetitionForm[K]) {
    setForm((prev) => ({ ...prev, [key]: value }))
  }

  function handleTitleChange(value: string) {
    updateForm("title", value)
    if (!slugTouched) {
      updateForm("slug", slugify(value))
    }
  }

  function handleSlugChange(value: string) {
    setSlugTouched(true)
    updateForm("slug", value)
  }

  function openCreate() {
    setEditingId(null)
    setEditingTicketsSold(0)
    setStatusOptions(statusOptionsFor(null))
    setForm(EMPTY_FORM)
    setSlugTouched(false)
    setFormError(null)
    setFormOpen(true)
  }

  function openEdit(competition: Competition) {
    setEditingId(competition.id)
    setEditingTicketsSold(competition.tickets_sold)
    setStatusOptions(statusOptionsFor(competition.status))
    setForm(toForm(competition))
    setSlugTouched(true)
    setFormError(null)
    setFormOpen(true)
  }

  async function handleSave(event: React.FormEvent) {
    event.preventDefault()
    setFormError(null)
    setSubmitting(true)

    const slugValue = form.slug.trim() ? slugify(form.slug) : ""

    const body: CompetitionInput = {
      title: form.title.trim(),
      description: form.description,
      prize: form.prize.trim(),
      ticket_price_pence: poundsToPence(form.price || "0"),
      tickets_total: Math.max(0, Number.parseInt(form.total, 10) || 0),
      category_id: form.categoryId || null,
      status: form.status,
      starts_at: localInputToRfc3339(form.startsAt),
      ends_at: localInputToRfc3339(form.endsAt),
    }
    if (slugValue) body.slug = slugValue

    try {
      if (editingId === null) {
        await apiPost("/apis/competition/v1/admin/competitions", body)
      } else {
        await apiPut(`/apis/competition/v1/admin/competitions/${editingId}`, body)
      }
      setFormOpen(false)
      await load()
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setSubmitting(false)
    }
  }

  function closeDeleteDialog() {
    setDeleteTarget(null)
    setDeleteError(null)
    setDeleteConfirmText("")
    setDeleteBlocked(false)
  }

  async function confirmDelete() {
    if (deleteTarget === null) return
    setDeleteError(null)
    setDeleting(true)
    try {
      await apiDelete(`/apis/competition/v1/admin/competitions/${deleteTarget.id}`)
      closeDeleteDialog()
      await load()
    } catch (err) {
      // 409 = the competition has sold tickets or a draw; offer to close instead.
      if (err instanceof ApiError && err.status === 409) {
        setDeleteBlocked(true)
        setDeleteError(err.message)
      } else {
        setDeleteError(err instanceof ApiError ? err.message : "Something went wrong")
      }
    } finally {
      setDeleting(false)
    }
  }

  // "Close competition instead" — the safe alternative to deleting a
  // competition with entrants. Sends the full-field PUT with status=closed.
  async function closeInstead() {
    if (deleteTarget === null) return
    setDeleteError(null)
    setDeleting(true)
    try {
      const body: CompetitionInput = {
        title: deleteTarget.title,
        slug: deleteTarget.slug,
        description: deleteTarget.description,
        prize: deleteTarget.prize,
        ticket_price_pence: deleteTarget.ticket_price_pence,
        tickets_total: deleteTarget.tickets_total,
        category_id: deleteTarget.category_id || null,
        status: "closed",
        starts_at: deleteTarget.starts_at,
        ends_at: deleteTarget.ends_at,
      }
      await apiPut(`/apis/competition/v1/admin/competitions/${deleteTarget.id}`, body)
      closeDeleteDialog()
      await load()
    } catch (err) {
      setDeleteError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setDeleting(false)
    }
  }

  return (
    <div className="space-y-4">
      {/* Toolbar: search (left) + status filter, refresh and create action (right). */}
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
              setStatusFilter(event.target.value as "all" | CompetitionApiStatus)
            }
            className="sm:w-44"
            aria-label="Filter by status"
          >
            <option value="all">All statuses</option>
            {STATUS_FILTER_ORDER.map((status) => (
              <option key={status} value={status}>
                {STATUS_META[status].label}
              </option>
            ))}
          </Select>

          <Button
            variant="outline"
            size="icon"
            onClick={() => void load()}
            disabled={loading}
            aria-label="Refresh competitions"
          >
            <RefreshCw className={loading ? "animate-spin" : undefined} />
          </Button>

          <Button onClick={openCreate} className="shrink-0">
            <Plus data-icon="inline-start" />
            New Competition
          </Button>
        </div>
      </div>

      {loadError && (
        <p className="inline-flex items-center gap-2 rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">
          <AlertCircle className="size-4" />
          {loadError}
        </p>
      )}

      <p className="text-xs text-muted-foreground">
        Showing {formatNumber(filtered.length)} of {formatNumber(competitions.length)}{" "}
        competitions
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
              {loading ? (
                <TableRow>
                  <TableCell
                    colSpan={7}
                    className="h-32 text-center text-sm text-muted-foreground"
                  >
                    <span className="inline-flex items-center gap-2">
                      <Loader2 className="size-4 animate-spin" />
                      Loading competitions…
                    </span>
                  </TableCell>
                </TableRow>
              ) : filtered.length === 0 ? (
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
                    competition.tickets_sold,
                    competition.tickets_total
                  )
                  const meta = STATUS_META[competition.status]

                  return (
                    <TableRow key={competition.id}>
                      <TableCell>
                        <div className="font-medium">{competition.title}</div>
                        <div className="text-xs text-muted-foreground">
                          {competition.category_name ?? "—"}
                        </div>
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {competition.prize}
                      </TableCell>
                      <TableCell>
                        {competition.ticket_price_pence === 0 ? (
                          <span className="font-medium text-botb-green">FREE</span>
                        ) : (
                          formatPence(competition.ticket_price_pence)
                        )}
                      </TableCell>
                      <TableCell>
                        <div className="min-w-32 space-y-1.5">
                          <div className="text-xs text-muted-foreground">
                            {formatNumber(competition.tickets_sold)} /{" "}
                            {formatNumber(competition.tickets_total)}
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
                        {formatEndDate(competition.ends_at)}
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
        <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-xl">
          <DialogHeader>
            <DialogTitle>
              {editingId === null ? "Create New Competition" : "Edit Competition"}
            </DialogTitle>
            <DialogDescription>
              {editingId === null
                ? "Add a new competition to the platform."
                : "Update the details for this competition."}
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleSave} className="space-y-4">
            {formError && (
              <p className="inline-flex items-center gap-2 rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">
                <AlertCircle className="size-4 shrink-0" />
                {formError}
              </p>
            )}

            <div className="space-y-2">
              <Label htmlFor="competition-title">Title</Label>
              <Input
                id="competition-title"
                value={form.title}
                onChange={(event) => handleTitleChange(event.target.value)}
                placeholder="Win a Defender D350 X"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="competition-slug">Slug</Label>
              <Input
                id="competition-slug"
                value={form.slug}
                onChange={(event) => handleSlugChange(event.target.value)}
                placeholder="win-a-defender-d350-x"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="competition-description">Description</Label>
              <textarea
                id="competition-description"
                value={form.description}
                onChange={(event) => updateForm("description", event.target.value)}
                placeholder="Tell players what they're winning…"
                className="flex min-h-20 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 dark:bg-input/30"
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
                <Select
                  id="competition-category"
                  value={form.categoryId}
                  onChange={(event) => updateForm("categoryId", event.target.value)}
                >
                  <option value="">— None —</option>
                  {categories.map((category) => (
                    <option key={category.id} value={category.id}>
                      {category.name}
                    </option>
                  ))}
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="competition-status">Status</Label>
                <Select
                  id="competition-status"
                  value={form.status}
                  onChange={(event) =>
                    updateForm("status", event.target.value as CompetitionApiStatus)
                  }
                >
                  {statusOptions.map((status) => (
                    <option key={status} value={status}>
                      {STATUS_META[status].label}
                    </option>
                  ))}
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="competition-starts">Starts at</Label>
                <Input
                  id="competition-starts"
                  type="datetime-local"
                  value={form.startsAt}
                  onChange={(event) => updateForm("startsAt", event.target.value)}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="competition-ends">Ends at</Label>
                <Input
                  id="competition-ends"
                  type="datetime-local"
                  value={form.endsAt}
                  onChange={(event) => updateForm("endsAt", event.target.value)}
                />
              </div>
            </div>

            {editingId !== null && (
              <p className="text-xs text-muted-foreground">
                Tickets sold: {formatNumber(editingTicketsSold)} (read-only)
              </p>
            )}

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setFormOpen(false)}
                disabled={submitting}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={submitting}>
                {submitting ? (
                  <>
                    <Loader2 className="animate-spin" />
                    Saving…
                  </>
                ) : editingId === null ? (
                  "Create Competition"
                ) : (
                  "Save Changes"
                )}
              </Button>
            </DialogFooter>
          </form>

          {editingId !== null && <CompetitionMediaManager competitionId={editingId} />}
        </DialogContent>
      </Dialog>

      {/* Delete confirmation dialog. */}
      <Dialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) closeDeleteDialog()
        }}
      >
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>
              {deleteBlocked ? "Can’t delete this competition" : "Delete competition?"}
            </DialogTitle>
            <DialogDescription>
              {deleteBlocked
                ? "It has sold tickets or a draw, so entrants have paid and results reference it. Close it instead — closed competitions stay on record but leave the live grids."
                : deleteTarget
                  ? `“${deleteTarget.title}” and its media will be permanently removed. This can’t be undone.`
                  : "This can’t be undone."}
            </DialogDescription>
          </DialogHeader>

          {/* Type-to-confirm guard (only when a plain delete is still possible). */}
          {!deleteBlocked && deleteTarget && (
            <div className="space-y-1.5">
              <Label htmlFor="delete-confirm">
                Type <span className="font-semibold">{deleteTarget.title}</span> to confirm
              </Label>
              <Input
                id="delete-confirm"
                value={deleteConfirmText}
                onChange={(e) => setDeleteConfirmText(e.target.value)}
                autoComplete="off"
                placeholder={deleteTarget.title}
              />
            </div>
          )}

          {deleteError && !deleteBlocked && (
            <p className="inline-flex items-center gap-2 rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">
              <AlertCircle className="size-4 shrink-0" />
              {deleteError}
            </p>
          )}

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={closeDeleteDialog}
              disabled={deleting}
            >
              Cancel
            </Button>
            {deleteBlocked ? (
              <Button
                type="button"
                onClick={() => void closeInstead()}
                disabled={deleting}
                className="gap-1.5"
              >
                {deleting ? (
                  <>
                    <Loader2 className="animate-spin" />
                    Closing…
                  </>
                ) : (
                  "Close competition instead"
                )}
              </Button>
            ) : (
              <Button
                type="button"
                variant="destructive"
                onClick={() => void confirmDelete()}
                disabled={deleting || deleteConfirmText.trim() !== deleteTarget?.title}
                className="gap-1.5"
              >
                {deleting ? (
                  <>
                    <Loader2 className="animate-spin" />
                    Deleting…
                  </>
                ) : (
                  <>
                    <Trash2 data-icon="inline-start" />
                    Delete
                  </>
                )}
              </Button>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
