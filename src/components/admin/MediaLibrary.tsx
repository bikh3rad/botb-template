"use client"

import * as React from "react"
import Link from "next/link"
import {
  ChevronLeft,
  ChevronRight,
  Image as ImageIcon,
  Loader2,
  Trash2,
} from "lucide-react"

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
import { apiDelete, apiGet, ApiError, mediaUrl } from "@/lib/admin/client"
import { formatDateTime, formatNumber } from "@/lib/admin/format"
import type { MediaItem } from "@/types/admin-api"

/** Response shape for the global media listing endpoint. */
interface MediaListResponse {
  media: MediaItem[]
  total: number
  count: number
  limit: number
  offset: number
}

const DEFAULT_LIMIT = 50

/** Paged grid of every media object across all owners, with delete support. */
export function MediaLibrary() {
  const [items, setItems] = React.useState<MediaItem[]>([])
  const [total, setTotal] = React.useState(0)
  const [limit] = React.useState(DEFAULT_LIMIT)
  const [offset, setOffset] = React.useState(0)
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)

  const [deleteTarget, setDeleteTarget] = React.useState<MediaItem | null>(null)
  const [deleteError, setDeleteError] = React.useState<string | null>(null)
  const [deleting, setDeleting] = React.useState(false)

  async function load() {
    setLoading(true)
    setError(null)
    try {
      const data = await apiGet<MediaListResponse>(
        `/apis/media/v1/admin/media?limit=${limit}&offset=${offset}`
      )
      setItems(data.media ?? [])
      setTotal(data.total ?? 0)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setLoading(false)
    }
  }

  React.useEffect(() => {
    void load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [offset, limit])

  function openDelete(item: MediaItem) {
    setDeleteTarget(item)
    setDeleteError(null)
  }

  function closeDelete() {
    setDeleteTarget(null)
    setDeleteError(null)
  }

  async function confirmDelete() {
    if (deleteTarget === null) return
    setDeleting(true)
    setDeleteError(null)
    try {
      await apiDelete(`/apis/media/v1/admin/media/${deleteTarget.id}`)
      closeDelete()
      await load()
    } catch (err) {
      setDeleteError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setDeleting(false)
    }
  }

  const rangeStart = total === 0 ? 0 : offset + 1
  const rangeEnd = Math.min(offset + limit, total)
  const canPrev = offset > 0
  const canNext = offset + limit < total

  return (
    <div className="space-y-4">
      {error && <p className="text-sm text-destructive">{error}</p>}

      {loading ? (
        <p className="inline-flex items-center gap-2 py-12 text-sm text-muted-foreground">
          <Loader2 className="size-4 animate-spin" />
          Loading media…
        </p>
      ) : items.length === 0 ? (
        <div className="flex flex-col items-center gap-2 py-16 text-sm text-muted-foreground">
          <ImageIcon className="size-6 text-muted-foreground/60" />
          No media uploaded yet.
        </div>
      ) : (
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6">
          {items.map((item) => {
            const url = mediaUrl({ bucket: item.bucket, object_key: item.object_key })
            return (
              <div
                key={item.id}
                className="group relative overflow-hidden rounded-lg border border-border bg-card"
              >
                {url ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    src={url}
                    alt=""
                    className="aspect-square w-full object-cover"
                  />
                ) : (
                  <div className="flex aspect-square w-full items-center justify-center bg-muted text-muted-foreground">
                    <ImageIcon className="size-6" />
                  </div>
                )}

                <Button
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  onClick={() => openDelete(item)}
                  className="absolute top-1.5 right-1.5 bg-background/80 text-muted-foreground backdrop-blur-sm hover:text-destructive"
                  aria-label="Delete media"
                >
                  <Trash2 />
                </Button>

                <div className="space-y-1 p-2.5">
                  <div className="flex items-center justify-between gap-2">
                    <Badge variant="secondary">{item.owner_type}</Badge>
                    {item.owner_type === "competition" ? (
                      <Link
                        href="/admin/competitions"
                        title={item.owner_id}
                        className="truncate text-xs text-muted-foreground hover:text-foreground hover:underline"
                      >
                        {item.owner_id}
                      </Link>
                    ) : (
                      <span
                        title={item.owner_id}
                        className="truncate text-xs text-muted-foreground"
                      >
                        {item.owner_id}
                      </span>
                    )}
                  </div>
                  <p className="truncate text-xs text-muted-foreground">
                    {item.content_type}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {formatDateTime(item.created_at)}
                  </p>
                </div>
              </div>
            )
          })}
        </div>
      )}

      <div className="flex items-center justify-between gap-3 border-t border-border pt-3">
        <p className="text-xs text-muted-foreground">
          Showing {formatNumber(rangeStart)}–{formatNumber(rangeEnd)} of{" "}
          {formatNumber(total)}
        </p>
        <div className="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            disabled={!canPrev}
            onClick={() => setOffset((prev) => Math.max(0, prev - limit))}
          >
            <ChevronLeft data-icon="inline-start" />
            Prev
          </Button>
          <Button
            type="button"
            variant="outline"
            size="sm"
            disabled={!canNext}
            onClick={() => setOffset((prev) => prev + limit)}
          >
            Next
            <ChevronRight />
          </Button>
        </div>
      </div>

      {/* Delete confirmation dialog. */}
      <Dialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) closeDelete()
        }}
      >
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Delete this media item?</DialogTitle>
            <DialogDescription>
              This file will be permanently removed from storage. This can’t be
              undone.
            </DialogDescription>
          </DialogHeader>

          {deleteError && <p className="text-sm text-destructive">{deleteError}</p>}

          <DialogFooter>
            <Button type="button" variant="outline" onClick={closeDelete}>
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={() => void confirmDelete()}
              disabled={deleting}
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
