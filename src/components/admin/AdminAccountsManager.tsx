"use client"

import * as React from "react"
import { useRouter } from "next/navigation"
import { Pencil, Plus, ShieldAlert, UserCog } from "lucide-react"

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
import { apiGet, apiPost, apiPut, ApiError } from "@/lib/admin/client"
import { useAdminAuth } from "@/lib/admin/auth-context"
import { formatDateTime } from "@/lib/admin/format"
import type { AdminAccount, AdminRole } from "@/types/admin-api"

/** Render order + label for the role select and badges. */
const ROLE_ORDER: AdminRole[] = ["admin", "superadmin"]

type BadgeVariant = React.ComponentProps<typeof Badge>["variant"]

const ROLE_BADGE: Record<AdminRole, BadgeVariant> = {
  superadmin: "default",
  admin: "secondary",
}

/** Fields captured by the create-account dialog. */
interface CreateForm {
  name: string
  email: string
  password: string
  role: AdminRole
}

const EMPTY_CREATE: CreateForm = {
  name: "",
  email: "",
  password: "",
  role: "admin",
}

/** Fields captured by the edit-account dialog. Password is optional. */
interface EditForm {
  name: string
  email: string
  role: AdminRole
  status: "active" | "disabled"
  password: string
}

const EMPTY_EDIT: EditForm = {
  name: "",
  email: "",
  role: "admin",
  status: "active",
  password: "",
}

function toEditForm(account: AdminAccount): EditForm {
  return {
    name: account.name,
    email: account.email,
    role: account.role,
    status: account.is_active ? "active" : "disabled",
    password: "",
  }
}

export function AdminAccountsManager() {
  const { isSuperadmin, loading: authLoading } = useAdminAuth()
  const router = useRouter()

  React.useEffect(() => {
    if (!authLoading && !isSuperadmin) router.replace("/admin")
  }, [authLoading, isSuperadmin, router])

  const [accounts, setAccounts] = React.useState<AdminAccount[]>([])
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)

  // Create dialog.
  const [createOpen, setCreateOpen] = React.useState(false)
  const [createForm, setCreateForm] = React.useState<CreateForm>(EMPTY_CREATE)
  const [createError, setCreateError] = React.useState<string | null>(null)
  const [creating, setCreating] = React.useState(false)

  // Edit dialog. `editTarget === null` means it is closed.
  const [editTarget, setEditTarget] = React.useState<AdminAccount | null>(null)
  const [editForm, setEditForm] = React.useState<EditForm>(EMPTY_EDIT)
  const [editError, setEditError] = React.useState<string | null>(null)
  const [saving, setSaving] = React.useState(false)

  async function load() {
    setLoading(true)
    setError(null)
    try {
      const data = await apiGet<{ accounts: AdminAccount[] }>(
        "/apis/adminauth/v1/admin/accounts"
      )
      setAccounts(data.accounts ?? [])
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setLoading(false)
    }
  }

  React.useEffect(() => {
    if (isSuperadmin) void load()
  }, [isSuperadmin])

  function updateCreateForm<K extends keyof CreateForm>(
    key: K,
    value: CreateForm[K]
  ) {
    setCreateForm((prev) => ({ ...prev, [key]: value }))
  }

  function updateEditForm<K extends keyof EditForm>(key: K, value: EditForm[K]) {
    setEditForm((prev) => ({ ...prev, [key]: value }))
  }

  function openCreate() {
    setCreateForm(EMPTY_CREATE)
    setCreateError(null)
    setCreateOpen(true)
  }

  function openEdit(account: AdminAccount) {
    setEditTarget(account)
    setEditForm(toEditForm(account))
    setEditError(null)
  }

  function closeEdit() {
    setEditTarget(null)
    setEditError(null)
  }

  async function handleCreate(event: React.FormEvent) {
    event.preventDefault()
    setCreateError(null)
    setCreating(true)

    try {
      await apiPost("/apis/adminauth/v1/admin/accounts", {
        name: createForm.name.trim(),
        email: createForm.email.trim(),
        password: createForm.password,
        role: createForm.role,
      })
      setCreateOpen(false)
      await load()
    } catch (err) {
      setCreateError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setCreating(false)
    }
  }

  async function handleEdit(event: React.FormEvent) {
    event.preventDefault()
    if (editTarget === null) return
    setEditError(null)
    setSaving(true)

    const body: {
      name: string
      email: string
      role: AdminRole
      is_active: boolean
      password?: string
    } = {
      name: editForm.name.trim(),
      email: editForm.email.trim(),
      role: editForm.role,
      is_active: editForm.status === "active",
    }
    if (editForm.password.trim()) {
      body.password = editForm.password
    }

    try {
      await apiPut(`/apis/adminauth/v1/admin/accounts/${editTarget.id}`, body)
      closeEdit()
      await load()
    } catch (err) {
      setEditError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setSaving(false)
    }
  }

  if (authLoading) {
    return <p className="text-sm text-muted-foreground">Loading…</p>
  }

  if (!isSuperadmin) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center gap-3 py-12 text-center">
          <ShieldAlert className="size-8 text-muted-foreground/60" />
          <div>
            <p className="font-medium">Forbidden</p>
            <p className="text-sm text-muted-foreground">
              You do not have permission to view this page.
            </p>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-4">
      {/* Toolbar: create action on the right. */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
        <Button onClick={openCreate} className="shrink-0">
          <Plus data-icon="inline-start" />
          New Account
        </Button>
      </div>

      {error && <p className="text-sm text-destructive">{error}</p>}

      <Card className="overflow-hidden py-0">
        <CardContent className="px-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Email</TableHead>
                <TableHead>Role</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Last login</TableHead>
                <TableHead>Created</TableHead>
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
                    Loading accounts…
                  </TableCell>
                </TableRow>
              ) : accounts.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={7}
                    className="h-32 text-center text-sm text-muted-foreground"
                  >
                    <span className="inline-flex flex-col items-center gap-2">
                      <UserCog className="size-6 text-muted-foreground/60" />
                      No admin accounts yet.
                    </span>
                  </TableCell>
                </TableRow>
              ) : (
                accounts.map((account) => (
                  <TableRow key={account.id}>
                    <TableCell className="font-medium">{account.name}</TableCell>
                    <TableCell className="text-muted-foreground">
                      {account.email}
                    </TableCell>
                    <TableCell>
                      <Badge variant={ROLE_BADGE[account.role]} className="capitalize">
                        {account.role}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {account.is_active ? (
                        <Badge variant="success">Active</Badge>
                      ) : (
                        <Badge variant="muted">Disabled</Badge>
                      )}
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {formatDateTime(account.last_login_at)}
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {formatDateTime(account.created_at)}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          onClick={() => openEdit(account)}
                          aria-label={`Edit ${account.name}`}
                        >
                          <Pencil />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Create dialog. */}
      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create New Account</DialogTitle>
            <DialogDescription>
              Add a new admin console account.
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleCreate} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="account-name">Name</Label>
              <Input
                id="account-name"
                value={createForm.name}
                onChange={(event) => updateCreateForm("name", event.target.value)}
                placeholder="Ada Lovelace"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="account-email">Email</Label>
              <Input
                id="account-email"
                type="email"
                value={createForm.email}
                onChange={(event) => updateCreateForm("email", event.target.value)}
                placeholder="ada@example.com"
                required
              />
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="account-password">Password</Label>
                <Input
                  id="account-password"
                  type="password"
                  value={createForm.password}
                  onChange={(event) =>
                    updateCreateForm("password", event.target.value)
                  }
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="account-role">Role</Label>
                <Select
                  id="account-role"
                  value={createForm.role}
                  onChange={(event) =>
                    updateCreateForm("role", event.target.value as AdminRole)
                  }
                >
                  {ROLE_ORDER.map((role) => (
                    <option key={role} value={role} className="capitalize">
                      {role}
                    </option>
                  ))}
                </Select>
              </div>
            </div>

            {createError && (
              <p className="text-sm text-destructive">{createError}</p>
            )}

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setCreateOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={creating}>
                Create Account
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Edit dialog. */}
      <Dialog
        open={editTarget !== null}
        onOpenChange={(open) => {
          if (!open) closeEdit()
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Account</DialogTitle>
            <DialogDescription>
              Update the details for this admin account.
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleEdit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="edit-account-name">Name</Label>
              <Input
                id="edit-account-name"
                value={editForm.name}
                onChange={(event) => updateEditForm("name", event.target.value)}
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="edit-account-email">Email</Label>
              <Input
                id="edit-account-email"
                type="email"
                value={editForm.email}
                onChange={(event) => updateEditForm("email", event.target.value)}
                required
              />
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="edit-account-role">Role</Label>
                <Select
                  id="edit-account-role"
                  value={editForm.role}
                  onChange={(event) =>
                    updateEditForm("role", event.target.value as AdminRole)
                  }
                >
                  {ROLE_ORDER.map((role) => (
                    <option key={role} value={role} className="capitalize">
                      {role}
                    </option>
                  ))}
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="edit-account-status">Status</Label>
                <Select
                  id="edit-account-status"
                  value={editForm.status}
                  onChange={(event) =>
                    updateEditForm(
                      "status",
                      event.target.value as "active" | "disabled"
                    )
                  }
                >
                  <option value="active">Active</option>
                  <option value="disabled">Disabled</option>
                </Select>
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="edit-account-password">New password</Label>
              <Input
                id="edit-account-password"
                type="password"
                value={editForm.password}
                onChange={(event) =>
                  updateEditForm("password", event.target.value)
                }
                placeholder="Leave blank to keep current"
              />
            </div>

            {editError && <p className="text-sm text-destructive">{editError}</p>}

            <DialogFooter>
              <Button type="button" variant="outline" onClick={closeEdit}>
                Cancel
              </Button>
              <Button type="submit" disabled={saving}>
                Save Changes
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
