// SERVER-ONLY typed fetch helpers for the Go backend gateway.
//
// This module is imported ONLY by Server Components / route handlers. It must
// never reach a Client Component: it reads the internal gateway URL and (via
// server-token.ts) the server-only JWT_SECRET. The `server-only` package is not
// available in this template, so the boundary is enforced by convention.
//
// All helpers fail SOFT: on any network/parse error they log and return an
// empty result so pages still render (and so `next build` never hard-fails when
// no backend is running).
import { mintAdminToken } from "@/lib/server-token";

// ---------------------------------------------------------------------------
// Backend JSON shapes (snake_case, mirrors the Go DTOs).
// ---------------------------------------------------------------------------

/** A media object attached to a competition (populated on read). */
export interface ApiMediaRef {
  id: string;
  kind: string;
  bucket: string;
  object_key: string;
  content_type: string;
  position: number;
}

/** A competition as returned by the competition service. */
export interface ApiCompetition {
  id: string;
  title: string;
  slug: string;
  description: string;
  prize: string;
  ticket_price_pence: number;
  tickets_total: number;
  tickets_sold: number;
  status: string;
  starts_at: string;
  ends_at: string;
  created_at: string;
  updated_at: string;
  media: ApiMediaRef[];
}

interface CompetitionListResp {
  count: number;
  competitions: ApiCompetition[];
}

/** A draw row from the admin draws list. */
export interface ApiDraw {
  id: string;
  competition_id: string;
  winner_user_id?: string;
  winner_ticket_id?: string;
  prize: string;
  status: string;
  drawn_at?: string;
  created_at: string;
  updated_at: string;
}

interface DrawListResp {
  draws: ApiDraw[];
}

/** A user row from the admin users list. */
export interface ApiUser {
  id: string;
  name: string;
  email: string;
  tickets_owned: number;
  total_spent_pence: number;
  created_at: string;
}

interface UserListResp {
  users: ApiUser[];
}

/** A media row from the public media-by-owner list. */
export interface ApiMedia {
  id: string;
  owner_type: string;
  owner_id: string;
  kind: string;
  bucket: string;
  object_key: string;
  content_type: string;
  position: number;
}

interface MediaListResp {
  media: ApiMedia[];
}

/** A resolved winner ready for presentation (public-safe fields only). */
export interface WinnerFeedItem {
  name: string;
  prize: string;
  /** Relative "revealed" label derived from the draw timestamp. */
  revealed: string;
  /** Avatar URL built from the user's first media object, or a fallback. */
  image: string;
}

// ---------------------------------------------------------------------------
// Configuration.
// ---------------------------------------------------------------------------

/**
 * Server-side gateway base URL. Prefers the internal compose-network address,
 * falling back to the browser-facing public URL, then localhost.
 */
function apiBase(): string {
  return (
    process.env.API_BASE_URL_INTERNAL ||
    process.env.NEXT_PUBLIC_API_BASE_URL ||
    "http://localhost:8080"
  );
}

/** Fallback avatar used when a winner has no avatar media object. */
const AVATAR_FALLBACK = "/images/winners/kfi-1.webp";

/** Shared fetch options: soft ISR-style caching (30s) across requests. */
const REVALIDATE_SECONDS = 30;

// ---------------------------------------------------------------------------
// Media URL helpers.
// ---------------------------------------------------------------------------

/** Browser-facing MinIO base URL (public-read bucket). */
export function mediaBase(): string {
  return process.env.NEXT_PUBLIC_MEDIA_BASE_URL || "http://localhost:9000";
}

/**
 * Build an object URL from a media record: `${MEDIA_BASE}/${bucket}/${key}`.
 * Returns null when the record is missing bucket/key.
 */
export function mediaObjectUrl(
  ref: { bucket: string; object_key: string } | undefined,
): string | null {
  if (!ref || !ref.bucket || !ref.object_key) return null;
  return `${mediaBase()}/${ref.bucket}/${ref.object_key}`;
}

/** Pick the lowest-position media ref from a list (the primary hero image). */
export function primaryMedia<T extends { position: number }>(
  refs: T[] | undefined,
): T | undefined {
  if (!refs || refs.length === 0) return undefined;
  return [...refs].sort((a, b) => a.position - b.position)[0];
}

// ---------------------------------------------------------------------------
// Competition fetches.
// ---------------------------------------------------------------------------

/**
 * Fetch competitions, filtered by status (default "live" so the closed
 * winners-archive competition never appears in public grids).
 */
export async function getCompetitions(
  status = "live",
): Promise<ApiCompetition[]> {
  const url = `${apiBase()}/apis/competition/v1/competitions?status=${encodeURIComponent(
    status,
  )}`;
  try {
    const res = await fetch(url, { next: { revalidate: REVALIDATE_SECONDS } });
    if (!res.ok) {
      console.error(`getCompetitions: ${res.status} from ${url}`);
      return [];
    }
    const data = (await res.json()) as CompetitionListResp;
    return data.competitions ?? [];
  } catch (err) {
    console.error("getCompetitions failed:", err);
    return [];
  }
}

/**
 * Fetch a single live competition by slug. The backend's single-competition
 * endpoint is keyed by UUID, so we resolve the slug against the live list
 * (which is cached). Returns null when not found or on error.
 */
export async function getCompetitionBySlug(
  slug: string,
): Promise<ApiCompetition | null> {
  const competitions = await getCompetitions("live");
  return competitions.find((c) => c.slug === slug) ?? null;
}

// ---------------------------------------------------------------------------
// Winners feed (admin-authenticated joins) — SERVER ONLY.
// ---------------------------------------------------------------------------

/** Authenticated GET against an admin endpoint using a freshly minted token. */
async function adminGet<T>(path: string): Promise<T | null> {
  const url = `${apiBase()}${path}`;
  try {
    const res = await fetch(url, {
      headers: { Authorization: `Bearer ${mintAdminToken()}` },
      next: { revalidate: REVALIDATE_SECONDS },
    });
    if (!res.ok) {
      console.error(`adminGet: ${res.status} from ${url}`);
      return null;
    }
    return (await res.json()) as T;
  } catch (err) {
    console.error(`adminGet ${path} failed:`, err);
    return null;
  }
}

/** Public GET (no auth) that tolerates failure by returning null. */
async function publicGet<T>(path: string): Promise<T | null> {
  const url = `${apiBase()}${path}`;
  try {
    const res = await fetch(url, { next: { revalidate: REVALIDATE_SECONDS } });
    if (!res.ok) return null;
    return (await res.json()) as T;
  } catch {
    return null;
  }
}

/** Turn an ISO draw timestamp into a human "revealed" label. */
function revealedLabel(drawnAt: string | undefined): string {
  if (!drawnAt) return "recently";
  const then = new Date(drawnAt).getTime();
  if (Number.isNaN(then)) return "recently";
  const diffMs = Date.now() - then;
  const minutes = Math.floor(diffMs / 60_000);
  if (minutes < 1) return "just now";
  if (minutes < 60) return `${minutes} minutes ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours} hours ago`;
  const days = Math.floor(hours / 24);
  if (days === 1) return "yesterday";
  return `${days} days ago`;
}

/** Resolve a winner's avatar URL from the public media-by-owner endpoint. */
async function winnerAvatar(userId: string): Promise<string> {
  const resp = await publicGet<MediaListResp>(
    `/apis/media/v1/media?owner_type=user&owner_id=${encodeURIComponent(
      userId,
    )}`,
  );
  const first = primaryMedia(resp?.media);
  return mediaObjectUrl(first) ?? AVATAR_FALLBACK;
}

/**
 * Build the public winners feed from the admin draw + user lists (the only
 * source of winner records). Joins each DRAWN draw to its user (for the name)
 * and to the user's avatar media. Runs only in Server Components. Ordered most
 * recently drawn first. Returns [] on any failure.
 */
export async function getWinners(): Promise<WinnerFeedItem[]> {
  const [drawsResp, usersResp] = await Promise.all([
    adminGet<DrawListResp>("/apis/draw/v1/admin/draws"),
    adminGet<UserListResp>("/apis/user/v1/admin/users"),
  ]);

  const draws = (drawsResp?.draws ?? []).filter(
    (d) => d.status === "drawn" && d.winner_user_id,
  );
  if (draws.length === 0) return [];

  const usersById = new Map<string, ApiUser>(
    (usersResp?.users ?? []).map((u) => [u.id, u]),
  );

  // Most recently drawn first (mirrors the "another winner, now" ordering).
  draws.sort((a, b) => {
    const at = a.drawn_at ? Date.parse(a.drawn_at) : 0;
    const bt = b.drawn_at ? Date.parse(b.drawn_at) : 0;
    return bt - at;
  });

  const items = await Promise.all(
    draws.map(async (draw): Promise<WinnerFeedItem> => {
      const userId = draw.winner_user_id as string;
      const user = usersById.get(userId);
      return {
        name: user?.name ?? "Winner",
        prize: draw.prize,
        revealed: revealedLabel(draw.drawn_at),
        image: await winnerAvatar(userId),
      };
    }),
  );

  return items;
}
