"use client"

import * as React from "react"
import { ChevronDown, ChevronUp, ImageOff, Loader2, Trash2, Upload } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { apiDelete, apiGet, apiPut, apiUpload, mediaUrl, ApiError } from "@/lib/admin/client"
import type { MediaItem } from "@/types/admin-api"

interface CompetitionMediaManagerProps {
  competitionId: string
}

/** Response shape for the media listing endpoint. */
interface MediaListResponse {
  media: MediaItem[]
}

/** Manages the media gallery (upload / reorder / delete) for one competition. */
export function CompetitionMediaManager({
  competitionId,
}: CompetitionMediaManagerProps) {
  const [items, setItems] = React.useState<MediaItem[]>([])
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)

  const [file, setFile] = React.useState<File | null>(null)
  const [uploading, setUploading] = React.useState(false)

  const [busyId, setBusyId] = React.useState<string | null>(null)
  const [deleteTarget, setDeleteTarget] = React.useState<MediaItem | null>(null)

  const fileInputRef = React.useRef<HTMLInputElement>(null)

  async function load() {
    setLoading(true)
    setError(null)
    try {
      const res = await apiGet<MediaListResponse>(
        `/apis/media/v1/admin/media?owner_type=competition&owner_id=${competitionId}`
      )
      const sorted = [...res.media].sort((a, b) => a.position - b.position)
      setItems(sorted)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setLoading(false)
    }
  }

  React.useEffect(() => {
    void load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [competitionId])

  async function handleUpload() {
    if (!file) return
    setUploading(true)
    setError(null)
    try {
      const nextPosition =
        items.length === 0 ? 0 : Math.max(...items.map((item) => item.position)) + 1
      const form = new FormData()
      form.append("file", file)
      form.append("owner_type", "competition")
      form.append("owner_id", competitionId)
      form.append("position", String(nextPosition))
      await apiUpload("/apis/media/v1/admin/uploads", form)
      setFile(null)
      if (fileInputRef.current) fileInputRef.current.value = ""
      await load()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setUploading(false)
    }
  }

  async function confirmDelete() {
    if (!deleteTarget) return
    setBusyId(deleteTarget.id)
    setError(null)
    try {
      await apiDelete(`/apis/media/v1/admin/media/${deleteTarget.id}`)
      setDeleteTarget(null)
      await load()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setBusyId(null)
    }
  }

  async function move(index: number, direction: -1 | 1) {
    const target = items[index]
    const neighbour = items[index + direction]
    if (!target || !neighbour) return

    setBusyId(target.id)
    setError(null)
    try {
      await Promise.all([
        apiPut(`/apis/media/v1/admin/media/${target.id}`, {
          position: neighbour.position,
        }),
        apiPut(`/apis/media/v1/admin/media/${neighbour.id}`, {
          position: target.position,
        }),
      ])
      await load()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setBusyId(null)
    }
  }

  return (
    <div className="space-y-3 rounded-lg border border-border p-4">
      <div className="flex items-center justify-between gap-2">
        <h4 className="text-sm font-medium">Media</h4>
      </div>

      {error && (
        <p className="rounded-md bg-destructive/10 px-3 py-2 text-xs text-destructive">
          {error}
        </p>
      )}

      {loading ? (
        <p className="inline-flex items-center gap-2 text-sm text-muted-foreground">
          <Loader2 className="size-4 animate-spin" />
          Loading media…
        </p>
      ) : items.length === 0 ? (
        <p className="text-sm text-muted-foreground">No media uploaded yet.</p>
      ) : (
        <div className="flex flex-wrap gap-3">
          {items.map((item, index) => {
            const url = mediaUrl({ bucket: item.bucket, object_key: item.object_key })
            const isBusy = busyId === item.id
            return (
              <div
                key={item.id}
                className="relative flex flex-col items-center gap-1.5 rounded-md border border-border p-2"
              >
                {url ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    src={url}
                    alt=""
                    className="h-20 w-20 rounded-md border border-border object-cover"
                  />
                ) : (
                  <div className="flex h-20 w-20 items-center justify-center rounded-md border border-border bg-muted text-muted-foreground">
                    <ImageOff className="size-5" />
                  </div>
                )}
                <span className="text-xs text-muted-foreground">#{item.position}</span>
                <div className="flex items-center gap-1">
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon-xs"
                    disabled={index === 0 || isBusy}
                    onClick={() => void move(index, -1)}
                    aria-label="Move earlier"
                  >
                    <ChevronUp />
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon-xs"
                    disabled={index === items.length - 1 || isBusy}
                    onClick={() => void move(index, 1)}
                    aria-label="Move later"
                  >
                    <ChevronDown />
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon-xs"
                    className="text-muted-foreground hover:text-destructive"
                    disabled={isBusy}
                    onClick={() => setDeleteTarget(item)}
                    aria-label="Delete media"
                  >
                    <Trash2 />
                  </Button>
                </div>
              </div>
            )
          })}
        </div>
      )}

      <div className="flex flex-wrap items-center gap-2 border-t border-border pt-3">
        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          onChange={(event) => setFile(event.target.files?.[0] ?? null)}
          className="flex-1 text-sm text-muted-foreground file:mr-3 file:rounded-md file:border-0 file:bg-muted file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-foreground"
        />
        <Button
          type="button"
          size="sm"
          disabled={!file || uploading}
          onClick={() => void handleUpload()}
        >
          {uploading ? (
            <>
              <Loader2 className="animate-spin" />
              Uploading…
            </>
          ) : (
            <>
              <Upload data-icon="inline-start" />
              Upload
            </>
          )}
        </Button>
      </div>

      <Dialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null)
        }}
      >
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Delete this image?</DialogTitle>
            <DialogDescription>
              This media item will be permanently removed. This can’t be undone.
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
              onClick={() => void confirmDelete()}
              disabled={busyId === deleteTarget?.id}
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
