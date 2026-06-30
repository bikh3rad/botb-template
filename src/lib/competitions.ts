import { competitionSections, featuredCompetitions } from "@/lib/data";
import type { Competition } from "@/types";

export interface CompetitionEntry extends Competition {
  slug: string;
  /** category section id this comp primarily belongs to */
  category: string;
  /** cash alternative shown on detail page */
  cashAlternative?: string;
  /** headline value (e.g. prize worth) */
  value?: string;
}

function slugify(title: string): string {
  return title
    .toLowerCase()
    .replace(/&/g, "and")
    .replace(/[£$€%+]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "")
    .slice(0, 60);
}

// Flatten all competitions across featured + sections, de-duplicated by slug.
const map = new Map<string, CompetitionEntry>();

function add(c: Competition, category: string) {
  const slug = slugify(c.title);
  // ensure uniqueness across categories
  if (map.has(slug) && map.get(slug)!.category !== category) {
    // keep first occurrence; skip duplicates (same prize appears in multiple sections)
    return;
  }
  if (map.has(slug)) return;
  map.set(slug, { ...c, slug, category });
}

featuredCompetitions.forEach((c) => add(c, "featured-competitions"));
competitionSections.forEach((s) => s.competitions.forEach((c) => add(c, s.id)));

export const allCompetitions: CompetitionEntry[] = Array.from(map.values());

export function getCompetitionBySlug(slug: string): CompetitionEntry | undefined {
  return allCompetitions.find((c) => c.slug === slug);
}

export function competitionHref(c: Competition): string {
  return `/prizes/${slugify(c.title)}`;
}

export { slugify };
