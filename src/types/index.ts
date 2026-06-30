// Content types for the BOTB homepage clone.

export interface HeroSlide {
  badge: string; // e.g. "ENDS SOON", "SUPERCAR"
  title: string; // e.g. "A LUXURY £1.2M HOME IN ZONE 1"
  subtitle: string;
  image: string; // /images/hero/*.webp
}

export interface CategoryNavItem {
  label: string;
  icon: string; // /images/pills/*.png
  targetId: string; // anchor id of the section
}

export interface Winner {
  name: string;
  prize: string;
  wonFor: string; // e.g. "£0.02"
  revealed: string; // e.g. "7 hours ago"
  image: string; // /images/winners/*.webp
}

/** Colored title-bar variant for a competition card. */
export type TitleBarColor =
  | "orange"
  | "teal"
  | "purple"
  | "red"
  | "green"
  | "blue"
  | "pass";

export interface Competition {
  badge: string; // e.g. "ENDS TONIGHT"
  title: string; // colored title bar text
  titleColor: TitleBarColor;
  description: string;
  /** Label above the price: "TICKET PRICE" or "STARTING FROM". */
  priceLabel?: string;
  price?: string; // e.g. "£1.25"
  /** Progress stat: { value: "92.5% sold" | "27.6k Won", percent, tickets } */
  stat?: {
    label: string; // "92.5% sold" or "27.6k Won"
    percent: number; // 0-100 progress fill
    tickets?: string; // "2.5M / 2.7M" or "170,326 Left"
  };
  cta: "ENTER NOW" | "PLAY NOW" | "DETAILS";
  /** Show the small orange cart button next to the CTA. */
  showCart?: boolean;
  image: string; // /images/comps/*.webp
  /** Wide card spanning 2 columns. */
  wide?: boolean;
}

export interface CompetitionSection {
  id: string;
  heading: string;
  subtitle?: string;
  competitions: Competition[];
}

export interface FooterLink {
  label: string;
  href: string;
}

export interface FooterColumn {
  title: string;
  links: FooterLink[];
  collapsible?: boolean; // accordion on mobile
}
