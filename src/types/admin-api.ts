// API shapes returned by the backend admin endpoints. These mirror the live
// JSON contract (snake_case) and are the source of truth for the rewritten
// admin editor UIs. Kept separate from the legacy mock types in `admin.ts`.

/** Competition lifecycle. Transitions are forward-only: draft → live → closed. */
export type CompetitionApiStatus = "draft" | "live" | "closed"

/** A media object attached to a competition (embedded in the competition read). */
export interface CompetitionMedia {
  id: string
  kind: string
  bucket: string
  object_key: string
  content_type: string
  position: number
}

/** A competition as returned by the competition service. */
export interface Competition {
  id: string
  title: string
  slug: string
  description: string
  prize: string
  ticket_price_pence: number
  tickets_total: number
  tickets_sold: number
  category_id: string | null
  category_name: string | null
  status: CompetitionApiStatus
  starts_at: string
  ends_at: string
  media: CompetitionMedia[]
}

/** Payload for creating/updating a competition (tickets_sold is read-only). */
export interface CompetitionInput {
  title: string
  slug?: string
  description: string
  prize: string
  ticket_price_pence: number
  tickets_total: number
  category_id?: string | null
  status: CompetitionApiStatus
  starts_at: string
  ends_at: string
}

/** A competition category. */
export interface Category {
  id: string
  name: string
  slug: string
  created_at: string
}

/** A registered player, as returned by the user service admin endpoints. */
export interface AdminUserRow {
  id: string
  name: string
  email: string
  tickets_owned: number
  total_spent_pence: number
  is_active: boolean
  created_at: string
}

/** A ticket owned by a user (read-only detail view). */
export interface UserTicket {
  id: string
  competition_id?: string
  competition_title?: string
  ticket_number?: number
  created_at?: string
}

/** Draw lifecycle. */
export type DrawApiStatus = "pending" | "drawn" | "void"

/** A prize draw for a competition. */
export interface Draw {
  id: string
  competition_id: string
  winner_user_id?: string | null
  winner_ticket_id?: string | null
  prize: string
  status: DrawApiStatus
  void_reason?: string | null
  drawn_at?: string | null
  created_at: string
  updated_at: string
}

/** A media object in the global library / per-owner listing. */
export interface MediaItem {
  id: string
  owner_type: string
  owner_id: string
  kind: string
  bucket: string
  object_key: string
  content_type: string
  position: number
  created_at: string
}

/** A single editable site-content row. */
export interface ContentRow {
  key: string
  value: string
  updated_at: string
}

/** Admin account role. */
export type AdminRole = "admin" | "superadmin"

/** An admin console account (superadmin-managed). */
export interface AdminAccount {
  id: string
  name: string
  email: string
  role: AdminRole
  is_active: boolean
  created_at: string
  last_login_at?: string | null
}
