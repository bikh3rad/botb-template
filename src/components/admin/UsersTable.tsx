"use client"

import * as React from "react"
import { ChevronLeft, ChevronRight, Search, Ticket } from "lucide-react"

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
import { Input } from "@/components/ui/input"
import { Select } from "@/components/ui/select"
import { adminUsers, formatDate, formatNumber, formatPence } from "@/lib/admin-data"
import type { AdminUser, UserStatus } from "@/types/admin"

/** Filter options for the status dropdown — "all" bypasses status filtering. */
type StatusFilter = "all" | UserStatus

/** Page-size choices offered by the pagination control. */
const PAGE_SIZE_OPTIONS = [10, 20, 50] as const

/** Maps a user status to its badge variant and human-readable label. */
const STATUS_META: Record<
  UserStatus,
  { label: string; variant: "success" | "warning" | "destructive" }
> = {
  active: { label: "Active", variant: "success" },
  vip: { label: "VIP", variant: "warning" },
  suspended: { label: "Suspended", variant: "destructive" },
}

/** Derive up to two uppercase initials from a full name, e.g. "Olivia Bennett" -> "OB". */
function getInitials(name: string): string {
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("")
}

export function UsersTable() {
  const [query, setQuery] = React.useState("")
  const [status, setStatus] = React.useState<StatusFilter>("all")
  const [pageSize, setPageSize] = React.useState<number>(PAGE_SIZE_OPTIONS[0])
  const [page, setPage] = React.useState(1)

  // Case-insensitive match on name OR email, combined with the status filter.
  const filtered = React.useMemo<AdminUser[]>(() => {
    const needle = query.trim().toLowerCase()
    return adminUsers.filter((user) => {
      const matchesStatus = status === "all" || user.status === status
      if (!matchesStatus) return false
      if (!needle) return true
      return (
        user.name.toLowerCase().includes(needle) ||
        user.email.toLowerCase().includes(needle)
      )
    })
  }, [query, status])

  const total = filtered.length
  const pageCount = Math.max(1, Math.ceil(total / pageSize))

  // Guard against the current page falling out of range after filtering or a
  // page-size change (e.g. searching down to a single page of results).
  const currentPage = Math.min(page, pageCount)
  const startIndex = (currentPage - 1) * pageSize
  const pageRows = filtered.slice(startIndex, startIndex + pageSize)

  // Human-friendly 1-based range for the "Showing X–Y of Z" label.
  const rangeStart = total === 0 ? 0 : startIndex + 1
  const rangeEnd = Math.min(startIndex + pageSize, total)

  // Any change to the result set or page size should return the user to page 1.
  function handleQueryChange(event: React.ChangeEvent<HTMLInputElement>) {
    setQuery(event.target.value)
    setPage(1)
  }

  function handleStatusChange(event: React.ChangeEvent<HTMLSelectElement>) {
    setStatus(event.target.value as StatusFilter)
    setPage(1)
  }

  function handlePageSizeChange(event: React.ChangeEvent<HTMLSelectElement>) {
    setPageSize(Number(event.target.value))
    setPage(1)
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
              value={query}
              onChange={handleQueryChange}
              placeholder="Search by name or email…"
              className="pl-9"
              aria-label="Search players by name or email"
            />
          </div>
          <div className="w-full sm:w-40">
            <Select
              value={status}
              onChange={handleStatusChange}
              aria-label="Filter by status"
            >
              <option value="all">All statuses</option>
              <option value="active">Active</option>
              <option value="vip">VIP</option>
              <option value="suspended">Suspended</option>
            </Select>
          </div>
        </div>
      </CardHeader>

      <CardContent className="px-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="px-6">User</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Tickets Owned</TableHead>
              <TableHead className="text-right">Total Spent</TableHead>
              <TableHead className="px-6 text-right">Joined</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {pageRows.length === 0 ? (
              <TableRow className="hover:bg-transparent">
                <TableCell colSpan={5} className="px-6 py-12 text-center">
                  <p className="font-heading text-sm font-semibold">
                    No users found
                  </p>
                  <p className="mt-1 text-sm text-muted-foreground">
                    Try a different name, email or status filter.
                  </p>
                </TableCell>
              </TableRow>
            ) : (
              pageRows.map((user) => {
                const statusMeta = STATUS_META[user.status]
                return (
                  <TableRow key={user.id}>
                    <TableCell className="px-6">
                      <div className="flex items-center gap-3">
                        <Avatar className="bg-primary/10">
                          <AvatarFallback className="bg-transparent text-botb-orange-hover dark:text-primary">
                            {getInitials(user.name)}
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
                      <Badge variant={statusMeta.variant}>{statusMeta.label}</Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <span className="inline-flex items-center justify-end gap-1.5 tabular-nums">
                        <Ticket className="size-3.5 text-muted-foreground" />
                        {formatNumber(user.ticketsOwned)}
                      </span>
                    </TableCell>
                    <TableCell className="text-right font-medium tabular-nums">
                      {formatPence(user.totalSpentPence)}
                    </TableCell>
                    <TableCell className="px-6 text-right text-muted-foreground tabular-nums">
                      {formatDate(user.joinDate)}
                    </TableCell>
                  </TableRow>
                )
              })
            )}
          </TableBody>
        </Table>
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
              value={String(pageSize)}
              onChange={handlePageSizeChange}
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

        <div className="flex items-center gap-3">
          <p className="text-sm text-muted-foreground">
            Page <span className="font-medium text-foreground">{currentPage}</span>{" "}
            of <span className="font-medium text-foreground">{pageCount}</span>
          </p>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPage((prev) => Math.max(1, prev - 1))}
              disabled={currentPage <= 1}
              aria-label="Previous page"
            >
              <ChevronLeft />
              Prev
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPage((prev) => Math.min(pageCount, prev + 1))}
              disabled={currentPage >= pageCount}
              aria-label="Next page"
            >
              Next
              <ChevronRight />
            </Button>
          </div>
        </div>
      </CardFooter>
    </Card>
  )
}
