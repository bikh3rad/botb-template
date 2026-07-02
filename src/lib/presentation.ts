// Presentation scaffolding for API-driven competitions.
//
// This is CSS-LIKE LAYOUT CONFIG, not data: it maps each backend competition
// (keyed by its slug = slugify(Title)) to the purely-visual attributes the
// backend schema does not model — accent colour, homepage section membership,
// per-placement badge text, CTA style, and layout flags. All real CONTENT
// (title, description, price, ticket counts, image) comes from the API.
//
// Slugs here MUST match backend/internal/seeddata/seeddata.go (the seed uses
// the same slugify rule), so the map keys line up with what the gateway serves.
import {
  mediaObjectUrl,
  primaryMedia,
  type ApiCompetition,
  type WinnerFeedItem,
} from "@/lib/api";
import type { Competition, TitleBarColor } from "@/types";

/** A competition card ready for <CompetitionCard>, carrying its slug for keys. */
export interface CardView extends Competition {
  slug: string;
}

/** A homepage competition row. */
export interface HomeSectionView {
  id: string;
  heading: string;
  subtitle?: string;
  cards: CardView[];
}

/** The bespoke "Dream Car" promo cell (price/description pulled from the API). */
export interface DreamCarView {
  steps: string[];
  titleBar: string;
  badge: string;
  heroImage: string;
  stbBadge: string;
  description: string;
  priceLabel: string;
  price: string;
  href: string;
}

/** The static new-subscriber promo cell (pure presentation). */
export interface SubscriberView {
  badge: string;
  image: string;
  headline: string;
  href: string;
}

/** The featured grid: 5 cards plus the two bespoke promo cells. */
export interface FeaturedView {
  cards: CardView[];
  dreamCar: DreamCarView;
  subscriber: SubscriberView;
}

/** Everything the homepage renders below the hero. */
export interface HomeView {
  featured: FeaturedView;
  sections: HomeSectionView[];
  winners: WinnerFeedItem[];
  winnersCount: string;
}

type Cta = Competition["cta"];
type StatStyle = "sold" | "won" | "none";

/** Intrinsic per-slug card style (accent, stat rendering, CTA, price label). */
interface SlugStyle {
  color: TitleBarColor;
  statStyle: StatStyle;
  cta: Cta;
  showCart: boolean;
  /** Present => show a price block with this label; absent => hide price. */
  priceLabel?: string;
}

/** A card's placement within a section (references a slug + per-spot overrides). */
interface CardEntry {
  slug: string;
  badge: string;
  wide?: boolean;
  /** Overrides the intrinsic showCart for this specific placement. */
  showCart?: boolean;
}

// ---------------------------------------------------------------------------
// Static promo cells + marketing copy (presentation only).
// ---------------------------------------------------------------------------

/** Fallback hero image when a competition has no media object. */
const FALLBACK_IMAGE = "/images/comps/fallback-subscribers.webp";

/** Slug of the competition that backs the bespoke Dream Car promo cell. */
const DREAM_CAR_SLUG = "dream-car";

/** Live-counter copy shown by the winners ticker (marketing, not winner data). */
export const winnersCount = "9,700 winners";

const DREAM_CAR_LAYOUT = {
  steps: ["Select prize", "Play the game", "Win your dream car"],
  titleBar: "DREAM CAR",
  badge: "ENDS SUNDAY",
  heroImage: "/images/comps/13732-wide.webp",
  stbBadge: "/images/comps/stb-badge.png",
  priceLabel: "STARTING FROM",
  href: "/prizes/cars",
};

const SUBSCRIBER_CARD: SubscriberView = {
  badge: "ENDS TONIGHT",
  image: "/images/comps/fallback-subscribers.webp",
  headline: "NEW SUBSCRIBERS GET 20 EXTRA LONDON HOME TICKETS",
  href: "/prizes/iphone-17-and-1249-prizes",
};

// ---------------------------------------------------------------------------
// Per-slug intrinsic styles.
// ---------------------------------------------------------------------------

const G = (extra?: Partial<SlugStyle>): SlugStyle => ({
  color: "green",
  statStyle: "sold",
  cta: "ENTER NOW",
  showCart: true,
  priceLabel: "TICKET PRICE",
  ...extra,
});

const SLUG_STYLE: Record<string, SlugStyle> = {
  // Featured hero prizes.
  "1-2m-home-in-zone-1": G({ color: "teal" }),
  "500k-instant-wins": G({ color: "purple", statStyle: "won" }),
  "audi-r8-for-21p": G({ color: "red" }),
  "2-5m-instant-wins": G({ color: "purple", statStyle: "won" }),
  "evoque-mini-for-9p": G({ color: "red" }),
  // Instant wins.
  "500k-midweek-instant-wins": G({ color: "purple", statStyle: "won" }),
  // "Details"-style card (no price/stat), always purple.
  "iphone-17-and-1249-prizes": {
    color: "purple",
    statStyle: "none",
    cta: "DETAILS",
    showCart: false,
  },
  // Dream Car (also backs the bespoke promo cell).
  "dream-car": {
    color: "orange",
    statStyle: "none",
    cta: "PLAY NOW",
    showCart: false,
    priceLabel: "STARTING FROM",
  },
  // Lifestyle (orange accent, standard priced card).
  "lifestyle-competition": G({ color: "orange" }),
  // Free comps (no price, no cart, no stat).
  "free-world-cup-tickets": {
    color: "blue",
    statStyle: "none",
    cta: "ENTER NOW",
    showCart: false,
  },
  "free-250-cash": {
    color: "blue",
    statStyle: "none",
    cta: "ENTER NOW",
    showCart: false,
  },
  "free-mystery-prize": {
    color: "blue",
    statStyle: "none",
    cta: "ENTER NOW",
    showCart: false,
  },
};

/** Resolve a slug's style, deriving a sensible default for unmapped slugs. */
function styleFor(comp: ApiCompetition): SlugStyle {
  const mapped = SLUG_STYLE[comp.slug];
  if (mapped) return mapped;
  // Derived fallback: free comps hide price; everything else is a green card.
  if (comp.ticket_price_pence === 0) {
    return { color: "blue", statStyle: "none", cta: "ENTER NOW", showCart: false };
  }
  return G();
}

// ---------------------------------------------------------------------------
// Homepage layout (which slugs appear where, with which badge).
// ---------------------------------------------------------------------------

const FEATURED_ENTRIES: CardEntry[] = [
  { slug: "1-2m-home-in-zone-1", badge: "ENDS TONIGHT" },
  { slug: "500k-instant-wins", badge: "ENDS TOMORROW", showCart: false },
  { slug: "audi-r8-for-21p", badge: "ENDS FRIDAY" },
  { slug: "2-5m-instant-wins", badge: "ENDS SUNDAY" },
  { slug: "evoque-mini-for-9p", badge: "ENDS IN 11 DAYS" },
];

interface SectionLayout {
  id: string;
  heading: string;
  subtitle?: string;
  entries: CardEntry[];
}

const SECTION_LAYOUT: SectionLayout[] = [
  {
    id: "ends-today",
    heading: "ENDS TODAY",
    subtitle: "Your last chance to enter, don't miss out!",
    entries: [
      { slug: "iphone-17-and-1249-prizes", badge: "ENDS TONIGHT" },
      { slug: "1-2m-home-in-zone-1", badge: "ENDS TONIGHT" },
      { slug: "1k-house-tickets", badge: "ENDS TONIGHT" },
      { slug: "500-house-tickets", badge: "ENDS TONIGHT" },
      { slug: "rattan-dining-set", badge: "ENDS TONIGHT" },
      { slug: "mystery-cash-prize", badge: "ENDS TONIGHT" },
      { slug: "ninja-autobarista-pro", badge: "ENDS TONIGHT" },
      { slug: "1000-house-tickets-1p", badge: "ENDS TONIGHT" },
    ],
  },
  {
    id: "ends-tomorrow",
    heading: "ENDS TOMORROW",
    subtitle: "Last chance to enter these competitions!",
    entries: [
      { slug: "iphone-17-and-1249-prizes", badge: "ENDS TOMORROW" },
      { slug: "5g-gold-bar", badge: "ENDS TOMORROW" },
      { slug: "2250-cash", badge: "ENDS TOMORROW" },
      { slug: "iphone-17-pro-max", badge: "ENDS TOMORROW" },
      { slug: "500-house-tickets", badge: "ENDS TOMORROW" },
      { slug: "500k-instant-wins", badge: "ENDS TOMORROW" },
    ],
  },
  {
    id: "instant-wins",
    heading: "INSTANT WINS",
    subtitle: "Play now to win instantly!",
    entries: [
      { slug: "500k-instant-wins", badge: "ENDS TOMORROW" },
      { slug: "2-5m-instant-wins", badge: "ENDS SUNDAY" },
      { slug: "500k-midweek-instant-wins", badge: "ENDS WEDNESDAY" },
    ],
  },
  {
    id: "ends-soon",
    heading: "ENDS SOON",
    subtitle: "These competitions won't be around for long.",
    entries: [
      { slug: "iphone-17-and-1249-prizes", badge: "ENDS THURSDAY" },
      { slug: "samsung-galaxy-book6", badge: "ENDS THURSDAY" },
      { slug: "macbook-neo", badge: "ENDS THURSDAY" },
      { slug: "toshiba-tv", badge: "ENDS THURSDAY" },
      { slug: "500-house-tickets", badge: "ENDS THURSDAY" },
      { slug: "1k-house-tickets", badge: "ENDS THURSDAY" },
      { slug: "audi-r8-for-21p", badge: "ENDS FRIDAY" },
      { slug: "nintendo-bundle", badge: "ENDS FRIDAY" },
      { slug: "mystery-lifestyle", badge: "ENDS FRIDAY" },
      { slug: "1250-cash", badge: "ENDS FRIDAY" },
      { slug: "5000-cash", badge: "ENDS SATURDAY" },
      { slug: "1k-amazon-voucher", badge: "ENDS SATURDAY" },
      { slug: "puremate-ac", badge: "ENDS SATURDAY" },
      { slug: "1000-house-tickets-1p", badge: "ENDS SATURDAY" },
      { slug: "ultimate-botb-pass", badge: "ENDS SUNDAY" },
      { slug: "mystery-tech-prize", badge: "ENDS SUNDAY" },
      { slug: "1750-cash", badge: "ENDS SUNDAY" },
      { slug: "dream-car", badge: "ENDS SUNDAY" },
      { slug: "10g-gold-bar", badge: "ENDS MONDAY", wide: true },
      { slug: "shark-fan-bundle", badge: "ENDS MONDAY" },
      { slug: "800-cash", badge: "ENDS MONDAY" },
      { slug: "lifestyle-competition", badge: "ENDS MONDAY", wide: true },
      { slug: "evoque-mini-for-9p", badge: "ENDS IN 11 DAYS" },
    ],
  },
  {
    id: "free-comps",
    heading: "FREE COMPS",
    subtitle: "Enter our free competitions — no purchase necessary!",
    entries: [
      { slug: "free-world-cup-tickets", badge: "ENDS TONIGHT" },
      { slug: "free-250-cash", badge: "ENDS FRIDAY" },
      { slug: "free-mystery-prize", badge: "ENDS SUNDAY" },
    ],
  },
];

/**
 * Category (section) a competition primarily belongs to, used by the
 * /competitions filter. Matches the categoryNav target ids.
 */
export function primaryCategory(slug: string): string {
  if (FEATURED_ENTRIES.some((e) => e.slug === slug)) return "featured-competitions";
  for (const section of SECTION_LAYOUT) {
    if (section.entries.some((e) => e.slug === slug)) return section.id;
  }
  return "ends-soon";
}

// ---------------------------------------------------------------------------
// Value formatting.
// ---------------------------------------------------------------------------

function fmtPrice(pence: number): string {
  return `£${(pence / 100).toFixed(2)}`;
}

/** Compact large counts: 2,700,000 -> "2.7M"; 12400 -> "12,400". */
function compact(n: number): string {
  if (n >= 1_000_000) {
    const millions = Math.round((n / 1_000_000) * 10) / 10;
    return `${millions}M`;
  }
  return n.toLocaleString("en-GB");
}

/** One-decimal percentage with trailing ".0" stripped: 92.5 -> "92.5", 62 -> "62". */
function trimDecimal(n: number): string {
  return (Math.round(n * 10) / 10).toString();
}

function buildStat(
  comp: ApiCompetition,
  style: StatStyle,
): Competition["stat"] | undefined {
  if (style === "none") return undefined;
  const total = comp.tickets_total;
  const sold = comp.tickets_sold;
  const pctRaw = total > 0 ? (sold / total) * 100 : 0;
  const percent = Math.round(pctRaw);
  if (style === "won") {
    return {
      label: `${compact(sold)} Won`,
      percent,
      tickets: `${compact(Math.max(0, total - sold))} Left`,
    };
  }
  return {
    label: `${trimDecimal(pctRaw)}% sold`,
    percent,
    tickets: `${compact(sold)} / ${compact(total)}`,
  };
}

/** Resolve a competition's hero image URL (falling back to a placeholder). */
export function competitionImage(comp: ApiCompetition): string {
  return mediaObjectUrl(primaryMedia(comp.media)) ?? FALLBACK_IMAGE;
}

// ---------------------------------------------------------------------------
// Mappers: ApiCompetition -> CardView.
// ---------------------------------------------------------------------------

/**
 * Map a backend competition into the exact `Competition` props the existing
 * <CompetitionCard> expects, applying the given placement (badge + overrides).
 */
export function toCardView(comp: ApiCompetition, entry: CardEntry): CardView {
  const style = styleFor(comp);
  const showPrice = style.priceLabel !== undefined;
  return {
    slug: comp.slug,
    badge: entry.badge,
    title: comp.title,
    titleColor: style.color,
    description: comp.description,
    priceLabel: showPrice ? style.priceLabel : undefined,
    price: showPrice ? fmtPrice(comp.ticket_price_pence) : undefined,
    stat: buildStat(comp, style.statStyle),
    cta: style.cta,
    showCart: entry.showCart ?? style.showCart,
    image: competitionImage(comp),
    wide: entry.wide,
  };
}

/** Derive a hero badge for a detail page from the competition's end time. */
function deriveBadge(comp: ApiCompetition): string {
  const hours = (Date.parse(comp.ends_at) - Date.now()) / 3_600_000;
  if (Number.isNaN(hours)) return "OPEN FOR ENTRIES";
  if (hours <= 12) return "ENDS TONIGHT";
  if (hours <= 36) return "ENDS TOMORROW";
  if (hours <= 24 * 7) return "ENDS THIS WEEK";
  return "OPEN FOR ENTRIES";
}

/**
 * Map a competition into a detail-page card view. Badge is derived from the end
 * time (detail pages have a single canonical badge, unlike the multi-placement
 * homepage). The prize string doubles as the "Prize Value" spec.
 */
export function toDetailView(
  comp: ApiCompetition,
): CardView & { value: string } {
  const card = toCardView(comp, { slug: comp.slug, badge: deriveBadge(comp) });
  return { ...card, value: comp.prize };
}

// ---------------------------------------------------------------------------
// Homepage assembly.
// ---------------------------------------------------------------------------

/** Build the featured + section views from the live competition list. */
export function buildHomeView(
  competitions: ApiCompetition[],
  winners: WinnerFeedItem[],
): HomeView {
  const bySlug = new Map(competitions.map((c) => [c.slug, c]));

  const mapEntries = (entries: CardEntry[]): CardView[] =>
    entries
      .map((entry) => {
        const comp = bySlug.get(entry.slug);
        return comp ? toCardView(comp, entry) : null;
      })
      .filter((c): c is CardView => c !== null);

  const dreamCarComp = bySlug.get(DREAM_CAR_SLUG);
  const dreamCar: DreamCarView = {
    ...DREAM_CAR_LAYOUT,
    description: dreamCarComp?.description ?? "Win a car PLUS gold!",
    price: dreamCarComp
      ? fmtPrice(dreamCarComp.ticket_price_pence)
      : "£1.13",
  };

  return {
    featured: {
      cards: mapEntries(FEATURED_ENTRIES),
      dreamCar,
      subscriber: SUBSCRIBER_CARD,
    },
    sections: SECTION_LAYOUT.map((section) => ({
      id: section.id,
      heading: section.heading,
      subtitle: section.subtitle,
      cards: mapEntries(section.entries),
    })),
    winners,
    winnersCount,
  };
}

/**
 * Build the flat, filterable list for the /competitions page: every live
 * competition once, tagged with its primary category.
 */
export function buildCompetitionList(
  competitions: ApiCompetition[],
): { card: CardView; category: string }[] {
  return competitions.map((comp) => ({
    card: toCardView(comp, {
      slug: comp.slug,
      badge: deriveBadge(comp),
    }),
    category: primaryCategory(comp.slug),
  }));
}
