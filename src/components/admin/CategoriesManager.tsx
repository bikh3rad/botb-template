"use client"

import * as React from "react"
import { Pencil, Plus, Tags, Trash2 } from "lucide-react"

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
import { apiDelete, apiGet, apiPost, apiPut, ApiError } from "@/lib/admin/client"
import { formatDate, slugify } from "@/lib/admin/format"
import type { Category } from "@/types/admin-api"

/** Editable fields captured by the create/rename dialog. */
interface CategoryForm {
  name: string
  slug: string
}

const EMPTY_FORM: CategoryForm = { name: "", slug: "" }

export function CategoriesManager() {
  const [categories, setCategories] = React.useState<Category[]>([])
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)

  // Create/rename dialog. `editingId === null` means we are creating a new row.
  const [formOpen, setFormOpen] = React.useState(false)
  const [editingId, setEditingId] = React.useState<string | null>(null)
  const [form, setForm] = React.useState<CategoryForm>(EMPTY_FORM)
  const [slugTouched, setSlugTouched] = React.useState(false)
  const [formError, setFormError] = React.useState<string | null>(null)
  const [submitting, setSubmitting] = React.useState(false)

  // Delete confirmation dialog target + reassignment state.
  const [deleteTarget, setDeleteTarget] = React.useState<Category | null>(null)
  const [deleteError, setDeleteError] = React.useState<string | null>(null)
  const [needsReassign, setNeedsReassign] = React.useState(false)
  const [reassignTo, setReassignTo] = React.useState("")
  const [deleting, setDeleting] = React.useState(false)

  async function load() {
    setLoading(true)
    setError(null)
    try {
      const data = await apiGet<{ categories: Category[] }>(
        "/apis/competition/v1/categories"
      )
      setCategories(data.categories ?? [])
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setLoading(false)
    }
  }

  React.useEffect(() => {
    void load()
  }, [])

  function updateForm<K extends keyof CategoryForm>(
    key: K,
    value: CategoryForm[K]
  ) {
    setForm((prev) => ({ ...prev, [key]: value }))
  }

  function handleNameChange(value: string) {
    setForm((prev) => ({
      ...prev,
      name: value,
      slug: slugTouched ? prev.slug : slugify(value),
    }))
  }

  function handleSlugChange(value: string) {
    setSlugTouched(true)
    updateForm("slug", value)
  }

  function openCreate() {
    setEditingId(null)
    setForm(EMPTY_FORM)
    setSlugTouched(false)
    setFormError(null)
    setFormOpen(true)
  }

  function openEdit(category: Category) {
    setEditingId(category.id)
    setForm({ name: category.name, slug: category.slug })
    setSlugTouched(true)
    setFormError(null)
    setFormOpen(true)
  }

  async function handleSave(event: React.FormEvent) {
    event.preventDefault()
    setFormError(null)
    setSubmitting(true)

    const name = form.name.trim()
    const slug = form.slug.trim()
    const body = slug ? { name, slug } : { name }

    try {
      if (editingId === null) {
        await apiPost("/apis/competition/v1/admin/categories", body)
      } else {
        await apiPut(`/apis/competition/v1/admin/categories/${editingId}`, body)
      }
      setFormOpen(false)
      await load()
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setSubmitting(false)
    }
  }

  function openDelete(category: Category) {
    setDeleteTarget(category)
    setDeleteError(null)
    setNeedsReassign(false)
    setReassignTo("")
  }

  function closeDelete() {
    setDeleteTarget(null)
    setDeleteError(null)
    setNeedsReassign(false)
    setReassignTo("")
  }

  async function confirmDelete() {
    if (deleteTarget === null) return
    setDeleting(true)
    setDeleteError(null)

    try {
      await apiDelete(`/apis/competition/v1/admin/categories/${deleteTarget.id}`)
      closeDelete()
      await load()
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        setNeedsReassign(true)
        setDeleteError(err.message)
      } else {
        setDeleteError(err instanceof ApiError ? err.message : "Something went wrong")
      }
    } finally {
      setDeleting(false)
    }
  }

  async function confirmReassignAndDelete() {
    if (deleteTarget === null || !reassignTo) return
    setDeleting(true)
    setDeleteError(null)

    try {
      await apiDelete(
        `/apis/competition/v1/admin/categories/${deleteTarget.id}?reassign_to=${reassignTo}`
      )
      closeDelete()
      await load()
    } catch (err) {
      setDeleteError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setDeleting(false)
    }
  }

  const reassignOptions = categories.filter((c) => c.id !== deleteTarget?.id)

  return (
    <div className="space-y-4">
      {/* Toolbar: create action on the right. */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
        <Button onClick={openCreate} className="shrink-0">
          <Plus data-icon="inline-start" />
          New Category
        </Button>
      </div>

      {error && (
        <p className="text-sm text-destructive">{error}</p>
      )}

      <Card className="overflow-hidden py-0">
        <CardContent className="px-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Slug</TableHead>
                <TableHead>Created</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell
                    colSpan={4}
                    className="h-32 text-center text-sm text-muted-foreground"
                  >
                    Loading categories…
                  </TableCell>
                </TableRow>
              ) : categories.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={4}
                    className="h-32 text-center text-sm text-muted-foreground"
                  >
                    <span className="inline-flex flex-col items-center gap-2">
                      <Tags className="size-6 text-muted-foreground/60" />
                      No categories yet. Create your first one to get started.
                    </span>
                  </TableCell>
                </TableRow>
              ) : (
                categories.map((category) => (
                  <TableRow key={category.id}>
                    <TableCell className="font-medium">{category.name}</TableCell>
                    <TableCell>
                      <span className="font-mono text-xs text-muted-foreground">
                        {category.slug}
                      </span>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {formatDate(category.created_at)}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          onClick={() => openEdit(category)}
                          aria-label={`Edit ${category.name}`}
                        >
                          <Pencil />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          onClick={() => openDelete(category)}
                          className="text-muted-foreground hover:text-destructive"
                          aria-label={`Delete ${category.name}`}
                        >
                          <Trash2 />
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

      {/* Create / rename dialog — reused for both flows. */}
      <Dialog open={formOpen} onOpenChange={setFormOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editingId === null ? "Create New Category" : "Rename Category"}
            </DialogTitle>
            <DialogDescription>
              {editingId === null
                ? "Add a new competition category."
                : "Update the name and slug for this category."}
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleSave} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="category-name">Name</Label>
              <Input
                id="category-name"
                value={form.name}
                onChange={(event) => handleNameChange(event.target.value)}
                placeholder="Cars"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="category-slug">Slug</Label>
              <Input
                id="category-slug"
                value={form.slug}
                onChange={(event) => handleSlugChange(event.target.value)}
                placeholder="cars"
              />
            </div>

            {formError && <p className="text-sm text-destructive">{formError}</p>}

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setFormOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={submitting}>
                {editingId === null ? "Create Category" : "Save Changes"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Delete confirmation dialog — reveals reassignment picker on 409. */}
      <Dialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) closeDelete()
        }}
      >
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Delete category?</DialogTitle>
            <DialogDescription>
              {deleteTarget
                ? `“${deleteTarget.name}” will be permanently removed. This can’t be undone.`
                : "This can’t be undone."}
            </DialogDescription>
          </DialogHeader>

          {deleteError && (
            <p className="text-sm text-destructive">{deleteError}</p>
          )}

          {needsReassign && (
            <div className="space-y-2">
              <Label htmlFor="category-reassign">Reassign competitions to</Label>
              <Select
                id="category-reassign"
                value={reassignTo}
                onChange={(event) => setReassignTo(event.target.value)}
                aria-label="Reassign competitions to"
              >
                <option value="">Select a category…</option>
                {reassignOptions.map((category) => (
                  <option key={category.id} value={category.id}>
                    {category.name}
                  </option>
                ))}
              </Select>
            </div>
          )}

          <DialogFooter>
            <Button type="button" variant="outline" onClick={closeDelete}>
              Cancel
            </Button>
            {needsReassign ? (
              <Button
                type="button"
                variant="destructive"
                onClick={confirmReassignAndDelete}
                disabled={!reassignTo || deleting}
                className="gap-1.5"
              >
                <Trash2 data-icon="inline-start" />
                Reassign &amp; delete
              </Button>
            ) : (
              <Button
                type="button"
                variant="destructive"
                onClick={confirmDelete}
                disabled={deleting}
                className="gap-1.5"
              >
                <Trash2 data-icon="inline-start" />
                Delete
              </Button>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
