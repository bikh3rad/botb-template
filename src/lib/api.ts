// SERVER-ONLY typed fetch helpers for the Go backend gateway.
//
// This module is imported ONLY by Server Components. It reads the internal
// gateway URL and calls PUBLIC endpoints — it holds no credentials.
//
// All helpers fail SOFT: on any network/parse error they log and return an
// empty result so pages still render (and so `next build` never hard-fails when
// no backend is running).
//
// This module no longer mints or holds any admin token: the winners feed reads
// the PUBLIC /apis/draw/v1/winners endpoint, and the admin panel authenticates
// through the adminauth service (see src/lib/admin/*). The old token-minting
// leak (a committed signing secret used to forge admin tokens) is gone.

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

/**
 * Browser-facing media base. Defaults to the same-origin "/media" path, which
 * the Next rewrite proxies to MinIO — so object URLs follow whatever server
 * IP/host the site is served from (no baked-in localhost). Set
 * NEXT_PUBLIC_MEDIA_BASE_URL to a full URL only to point the browser straight
 * at a public MinIO/CDN host instead.
 */
export function mediaBase(): string {
  return process.env.NEXT_PUBLIC_MEDIA_BASE_URL || "/media";
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
// Winners feed — SERVER ONLY, PUBLIC (no admin token).
// ---------------------------------------------------------------------------

/** A winner row from the public winners feed. */
interface ApiWinner {
  draw_id: string;
  prize: string;
  drawn_at?: string;
  winner_user_id: string;
  winner_name: string;
}

interface WinnerListResp {
  winners: ApiWinner[];
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
 * Build the public winners feed from the PUBLIC winners endpoint (the draw
 * service joins drawn draws to winner names server-side). Resolves each
 * winner's avatar from the public media-by-owner endpoint. Runs only in Server
 * Components. Ordered most recently drawn first. Returns [] on any failure.
 * No admin token is involved — winners are public data.
 */
export async function getWinners(): Promise<WinnerFeedItem[]> {
  const resp = await publicGet<WinnerListResp>("/apis/draw/v1/winners?limit=24");
  const winners = resp?.winners ?? [];
  if (winners.length === 0) return [];

  const items = await Promise.all(
    winners.map(async (winner): Promise<WinnerFeedItem> => {
      return {
        name: winner.winner_name || "Winner",
        prize: winner.prize,
        revealed: revealedLabel(winner.drawn_at),
        image: await winnerAvatar(winner.winner_user_id),
      };
    }),
  );

  return items;
}
