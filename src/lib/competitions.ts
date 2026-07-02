// Pure, client-safe helpers for competition slugs + links.
//
// This module intentionally holds NO data — competition content is served by
// the backend (see src/lib/api.ts). It is safe to import from Client
// Components. `slugify` mirrors the backend's slug rule so a title maps to the
// same slug the gateway stores.

/** Turn a competition title into its URL-safe slug (matches the backend). */
export function slugify(title: string): string {
  return title
    .toLowerCase()
    .replace(/&/g, "and")
    .replace(/[£$€%+]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "")
    .slice(0, 60);
}

/** Link to a competition's detail page. Accepts anything with a slug or title. */
export function competitionHref(c: { slug?: string; title: string }): string {
  const slug = c.slug ?? slugify(c.title);
  return `/prizes/${slug}`;
}
