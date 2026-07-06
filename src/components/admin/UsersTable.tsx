"use client"

import * as React from "react"
import { ChevronLeft, ChevronRight, Loader2, Pencil, Search, Ticket } from "lucide-react"

import {
  Card,
  CardContent,
  CardFooter,
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
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import { apiGet, apiPost, apiPut, ApiError } from "@/lib/admin/client"
import { formatDate, formatNumber, formatPence, initials } from "@/lib/admin/format"
import type { AdminUserRow } from "@/types/admin-api"

/** Response shape for the users listing endpoint. */
interface UsersListResponse {
  users: AdminUserRow[]
  total: number
  count: number
  limit: number
  offset: number
}

/** Page-size choices offered by the pagination control. */
const PAGE_SIZE_OPTIONS = [10, 20, 50] as const

/** Editable fields captured by the edit dialog. */
interface UserForm {
  name: string
  email: string
}

/** Build a form snapshot from an existing user for editing. */
function toForm(user: AdminUserRow): UserForm {
  return { name: user.name, email: user.email }
}

export function UsersTable() {
  const [users, setUsers] = React.useState<AdminUserRow[]>([])
  const [total, setTotal] = React.useState(0)
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)

  const [q, setQ] = React.useState("")
  const [debouncedQ, setDebouncedQ] = React.useState("")
  const [limit, setLimit] = React.useState<number>(PAGE_SIZE_OPTIONS[1])
  const [offset, setOffset] = React.useState(0)

  // Edit dialog state.
  const [editingUser, setEditingUser] = React.useState<AdminUserRow | null>(null)
  const [form, setForm] = React.useState<UserForm>({ name: "", email: "" })
  const [saving, setSaving] = React.useState(false)
  const [formError, setFormError] = React.useState<string | null>(null)

  // Suspend/activate confirmation dialog state.
  const [statusTarget, setStatusTarget] = React.useState<AdminUserRow | null>(null)
  const [statusBusy, setStatusBusy] = React.useState(false)
  const [statusError, setStatusError] = React.useState<string | null>(null)

  // Debounce the search input by ~300ms before it drives a request.
  React.useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQ(q)
      setOffset(0)
    }, 300)
    return () => clearTimeout(timer)
  }, [q])

  const load = React.useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const params = new URLSearchParams({
        q: debouncedQ,
        limit: String(limit),
        offset: String(offset),
      })
      const res = await apiGet<UsersListResponse>(
        `/apis/user/v1/admin/users?${params.toString()}`
      )
      setUsers(res.users)
      setTotal(res.total)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setLoading(false)
    }
  }, [debouncedQ, limit, offset])

  React.useEffect(() => {
    void load()
  }, [load])

  const rangeStart = total === 0 ? 0 : offset + 1
  const rangeEnd = Math.min(offset + limit, total)

  function handleQueryChange(event: React.ChangeEvent<HTMLInputElement>) {
    setQ(event.target.value)
  }

  function handleLimitChange(event: React.ChangeEvent<HTMLSelectElement>) {
    setLimit(Number(event.target.value))
    setOffset(0)
  }

  function openEdit(user: AdminUserRow) {
    setEditingUser(user)
    setForm(toForm(user))
    setFormError(null)
  }

  async function handleSave(event: React.FormEvent) {
    event.preventDefault()
    if (!editingUser) return

    setSaving(true)
    setFormError(null)
    try {
      await apiPut(`/apis/user/v1/admin/users/${editingUser.id}`, {
        name: form.name.trim(),
        email: form.email.trim(),
      })
      setEditingUser(null)
      await load()
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setSaving(false)
    }
  }

  async function confirmStatusChange() {
    if (!statusTarget) return
    setStatusBusy(true)
    setStatusError(null)
    try {
      const action = statusTarget.is_active ? "suspend" : "activate"
      await apiPost(`/apis/user/v1/admin/users/${statusTarget.id}/${action}`)
      setStatusTarget(null)
      await load()
    } catch (err) {
      setStatusError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setStatusBusy(false)
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>All players</CardTitle>
        <div className="flex flex-col gap-3 pt-2 sm:flex-row sm:items-center sm:justify-between">
          <div className="relative w-full sm:max-w-xs">
            <Search className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="search"
              value={q}
              onChange={handleQueryChange}
              placeholder="Search by name or email…"
              className="pl-9"
              aria-label="Search players by name or email"
            />
          </div>
        </div>
      </CardHeader>

      <CardContent className="px-0">
        {error ? (
          <p className="mx-6 rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">
            {error}
          </p>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="px-6">User</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Tickets Owned</TableHead>
                <TableHead className="text-right">Total Spent</TableHead>
                <TableHead className="px-6">Joined</TableHead>
                <TableHead className="px-6 text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow className="hover:bg-transparent">
                  <TableCell colSpan={6} className="px-6 py-12 text-center">
                    <span className="inline-flex items-center gap-2 text-sm text-muted-foreground">
                      <Loader2 className="size-4 animate-spin" />
                      Loading users…
                    </span>
                  </TableCell>
                </TableRow>
              ) : users.length === 0 ? (
                <TableRow className="hover:bg-transparent">
                  <TableCell colSpan={6} className="px-6 py-12 text-center">
                    <p className="font-heading text-sm font-semibold">
                      No users found
                    </p>
                    <p className="mt-1 text-sm text-muted-foreground">
                      Try a different name or email.
                    </p>
                  </TableCell>
                </TableRow>
              ) : (
                users.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell className="px-6">
                      <div className="flex items-center gap-3">
                        <Avatar className="bg-primary/10">
                          <AvatarFallback className="bg-transparent text-botb-orange-hover dark:text-primary">
                            {initials(user.name)}
                          </AvatarFallback>
                        </Avatar>
                        <div className="min-w-0">
                          <div className="font-medium">{user.name}</div>
                          <div className="truncate text-xs text-muted-foreground">
                            {user.email}
                          </div>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      {user.is_active ? (
                        <Badge variant="success">Active</Badge>
                      ) : (
                        <Badge variant="destructive">Suspended</Badge>
                      )}
                    </TableCell>
                    <TableCell className="text-right">
                      <span className="inline-flex items-center justify-end gap-1.5 tabular-nums">
                        <Ticket className="size-3.5 text-muted-foreground" />
                        {formatNumber(user.tickets_owned)}
                      </span>
                    </TableCell>
                    <TableCell className="text-right font-medium tabular-nums">
                      {formatPence(user.total_spent_pence)}
                    </TableCell>
                    <TableCell className="px-6 text-muted-foreground tabular-nums">
                      {formatDate(user.created_at)}
                    </TableCell>
                    <TableCell className="px-6">
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          onClick={() => openEdit(user)}
                          aria-label={`Edit ${user.name}`}
                        >
                          <Pencil />
                        </Button>
                        {user.is_active ? (
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => setStatusTarget(user)}
                            className="text-muted-foreground hover:text-destructive"
                          >
                            Suspend
                          </Button>
                        ) : (
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => setStatusTarget(user)}
                          >
                            Activate
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        )}
      </CardContent>

      <CardFooter className="flex flex-col gap-4 border-t border-border pt-6 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <p className="text-sm text-muted-foreground">
            Showing{" "}
            <span className="font-medium text-foreground">{rangeStart}</span>–
            <span className="font-medium text-foreground">{rangeEnd}</span> of{" "}
            <span className="font-medium text-foreground">{total}</span>
          </p>
          <div className="w-[4.5rem]">
            <Select
              value={String(limit)}
              onChange={handleLimitChange}
              aria-label="Rows per page"
            >
              {PAGE_SIZE_OPTIONS.map((size) => (
                <option key={size} value={size}>
                  {size}
                </option>
              ))}
            </Select>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setOffset((prev) => Math.max(0, prev - limit))}
            disabled={offset === 0}
            aria-label="Previous page"
          >
            <ChevronLeft />
            Prev
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setOffset((prev) => prev + limit)}
            disabled={offset + limit >= total}
            aria-label="Next page"
          >
            Next
            <ChevronRight />
          </Button>
        </div>
      </CardFooter>

      {/* Edit dialog. */}
      <Dialog
        open={editingUser !== null}
        onOpenChange={(open) => {
          if (!open) setEditingUser(null)
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit user</DialogTitle>
            <DialogDescription>
              Update this player&rsquo;s name and email address.
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleSave} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="user-name">Name</Label>
              <Input
                id="user-name"
                value={form.name}
                onChange={(event) =>
                  setForm((prev) => ({ ...prev, name: event.target.value }))
                }
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="user-email">Email</Label>
              <Input
                id="user-email"
                type="email"
                value={form.email}
                onChange={(event) =>
                  setForm((prev) => ({ ...prev, email: event.target.value }))
                }
                required
              />
            </div>

            {editingUser ? (
              <p className="text-xs text-muted-foreground">
                Tickets owned: {formatNumber(editingUser.tickets_owned)} · Total
                spent: {formatPence(editingUser.total_spent_pence)} — read-only
              </p>
            ) : null}

            {formError ? (
              <p className="rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">
                {formError}
              </p>
            ) : null}

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setEditingUser(null)}
                disabled={saving}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={saving}>
                {saving ? (
                  <>
                    <Loader2 className="animate-spin" />
                    Saving…
                  </>
                ) : (
                  "Save Changes"
                )}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Suspend / activate confirmation dialog. */}
      <Dialog
        open={statusTarget !== null}
        onOpenChange={(open) => {
          if (!open) {
            setStatusTarget(null)
            setStatusError(null)
          }
        }}
      >
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>
              {statusTarget?.is_active ? "Suspend user?" : "Activate user?"}
            </DialogTitle>
            <DialogDescription>
              {statusTarget
                ? statusTarget.is_active
                  ? `${statusTarget.name} will lose access to their account until reactivated.`
                  : `${statusTarget.name} will regain access to their account.`
                : ""}
            </DialogDescription>
          </DialogHeader>

          {statusError ? (
            <p className="rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">
              {statusError}
            </p>
          ) : null}

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setStatusTarget(null)}
              disabled={statusBusy}
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant={statusTarget?.is_active ? "destructive" : "default"}
              onClick={() => void confirmStatusChange()}
              disabled={statusBusy}
            >
              {statusBusy ? (
                <Loader2 className="animate-spin" />
              ) : statusTarget?.is_active ? (
                "Suspend"
              ) : (
                "Activate"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  )
}
