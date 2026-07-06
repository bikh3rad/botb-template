"use client"

import * as React from "react"
import { Check, FileText, Search } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { apiGet, apiPut, ApiError } from "@/lib/admin/client"
import { formatDateTime } from "@/lib/admin/format"
import type { ContentRow } from "@/types/admin-api"

/** Shape returned by GET /apis/competition/v1/content. */
interface ContentResponse {
  items: Record<string, string>
  rows: ContentRow[]
}

export function SiteTextsEditor() {
  const [rows, setRows] = React.useState<ContentRow[]>([])
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)
  const [query, setQuery] = React.useState("")

  // Local editable copy of each row's value, keyed by content key.
  const [drafts, setDrafts] = React.useState<Record<string, string>>({})
  // Per-row save-in-flight state.
  const [saving, setSaving] = React.useState<Record<string, boolean>>({})
  // Per-row save error, keyed by content key.
  const [saveErrors, setSaveErrors] = React.useState<Record<string, string>>({})
  // Per-row "saved" success flash, keyed by content key.
  const [saved, setSaved] = React.useState<Record<string, boolean>>({})

  async function load() {
    setLoading(true)
    setError(null)
    try {
      const data = await apiGet<ContentResponse>("/apis/competition/v1/content")
      const sorted = [...(data.rows ?? [])].sort((a, b) =>
        a.key.localeCompare(b.key)
      )
      setRows(sorted)
      setDrafts((prev) => {
        const next: Record<string, string> = {}
        for (const row of sorted) {
          // Preserve any in-progress edit that hasn't been saved yet.
          next[row.key] = prev[row.key] ?? row.value
        }
        return next
      })
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong")
    } finally {
      setLoading(false)
    }
  }

  React.useEffect(() => {
    void load()
  }, [])

  function updateDraft(key: string, value: string) {
    setDrafts((prev) => ({ ...prev, [key]: value }))
    setSaved((prev) => (prev[key] ? { ...prev, [key]: false } : prev))
  }

  function resetDraft(key: string, original: string) {
    setDrafts((prev) => ({ ...prev, [key]: original }))
    setSaveErrors((prev) => {
      if (!(key in prev)) return prev
      const next = { ...prev }
      delete next[key]
      return next
    })
    setSaved((prev) => (prev[key] ? { ...prev, [key]: false } : prev))
  }

  async function save(key: string) {
    const value = drafts[key] ?? ""
    setSaving((prev) => ({ ...prev, [key]: true }))
    setSaveErrors((prev) => {
      if (!(key in prev)) return prev
      const next = { ...prev }
      delete next[key]
      return next
    })

    try {
      await apiPut(
        `/apis/competition/v1/admin/content/${encodeURIComponent(key)}`,
        { value }
      )
      await load()
      setSaved((prev) => ({ ...prev, [key]: true }))
    } catch (err) {
      setSaveErrors((prev) => ({
        ...prev,
        [key]: err instanceof ApiError ? err.message : "Something went wrong",
      }))
    } finally {
      setSaving((prev) => ({ ...prev, [key]: false }))
    }
  }

  const filtered = React.useMemo(() => {
    const term = query.trim().toLowerCase()
    if (term === "") return rows
    return rows.filter((row) => row.key.toLowerCase().includes(term))
  }, [rows, query])

  return (
    <div className="space-y-4">
      <div className="relative w-full sm:max-w-xs">
        <Search className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          type="search"
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Search keys…"
          className="pl-9"
          aria-label="Search site texts"
        />
      </div>

      {error && <p className="text-sm text-destructive">{error}</p>}

      {loading ? (
        <Card>
          <CardContent className="py-10 text-center text-sm text-muted-foreground">
            Loading site texts…
          </CardContent>
        </Card>
      ) : filtered.length === 0 ? (
        <Card>
          <CardContent className="py-10 text-center text-sm text-muted-foreground">
            <span className="inline-flex flex-col items-center gap-2">
              <FileText className="size-6 text-muted-foreground/60" />
              {rows.length === 0
                ? "No site texts yet."
                : "No keys match your search."}
            </span>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {filtered.map((row) => {
            const draft = drafts[row.key] ?? row.value
            const isDirty = draft !== row.value
            const isSaving = saving[row.key] ?? false
            const rowError = saveErrors[row.key]
            const isSaved = saved[row.key] ?? false

            return (
              <Card key={row.key}>
                <CardHeader>
                  <CardTitle className="font-mono text-sm font-medium">
                    {row.key}
                  </CardTitle>
                  <CardDescription>
                    Updated {formatDateTime(row.updated_at)}
                  </CardDescription>
                </CardHeader>

                <CardContent>
                  <textarea
                    value={draft}
                    onChange={(event) => updateDraft(row.key, event.target.value)}
                    disabled={isSaving}
                    aria-label={`Value for ${row.key}`}
                    className="flex min-h-24 w-full rounded-lg border border-input bg-background px-3 py-2 font-mono text-xs shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 dark:bg-input/30"
                  />
                  {rowError && (
                    <p className="mt-2 text-sm text-destructive">{rowError}</p>
                  )}
                </CardContent>

                <CardFooter className="justify-end gap-2">
                  {isSaved && !isDirty && (
                    <span className="mr-auto inline-flex items-center gap-1 text-xs font-medium text-botb-green">
                      <Check className="size-3.5" />
                      Saved
                    </span>
                  )}
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => resetDraft(row.key, row.value)}
                    disabled={!isDirty || isSaving}
                  >
                    Reset
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    onClick={() => void save(row.key)}
                    disabled={!isDirty || isSaving}
                  >
                    {isSaving ? "Saving…" : "Save"}
                  </Button>
                </CardFooter>
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}
