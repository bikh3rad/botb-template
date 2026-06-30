# BOTB Clone — Component Specifications

Specs that were the contract for each builder. Source of truth: extracted CSS in
`DESIGN_TOKENS.md`, behaviors in `BEHAVIORS.md`, content in `../sections.json` and `src/lib/data.ts`.

## CompetitionCard (`src/components/CompetitionCard.tsx`) — shared linchpin
- White card, `rounded-lg` (8px), border `1px #C3C3C3`, shadow `0 1.5px 8px rgba(0,0,0,0.2)`, overflow-hidden, hover lift.
- Image `aspect-[342/212]` object-cover; **ENDS badge** top-left = pink→orange gradient (`.botb-badge-gradient`), Jost 12px/500 uppercase white, rounded px-2 py-1.
- **Colored title bar**: Jost 14px/700 uppercase white, padding 12px 20px. Colors: orange #FF8200, teal #32BAB5, purple #562373, red #AE1F25, green #3FA63F, blue #3772FF.
- Body (centered): description 15px #333; price block (label Roboto 11px/500 uppercase #898994 + value Jost 22px/700); progress (`% sold` orange 14px/600 + tickets w/ ticket icon, bar track #e5e5e5 fill orange 6px).
- CTA variants: `ENTER NOW` outline (orange border, fills on hover) + optional orange cart square; `PLAY NOW` filled orange; `DETAILS` link + filled orange ENTER NOW.
- `wide` → `sm:col-span-2` with 684×424 image.

## SiteHeader (`SiteHeader.tsx`) — sticky `top-0 z-[1003]`, client
- Black promo bar (red NEW pill + "🎁 FINAL DAYS…"). White nav: desktop 3-zone (COMPETITIONS/WINNERS · centered Logo · LOG IN/SIGN UP orange/cart+badge); mobile hamburger + logo + account/cart + dropdown menu. Logo SVG already contains the tagline.

## HeroCarousel (`HeroCarousel.tsx`) — black bg, time-driven, client
- 8 full-bleed slides (2496×787), overlay badge/title (Jost bold uppercase)/subtitle/ENTER TO WIN + cart. Prev/next chevrons (md+), 8 dots (active wider/white). Auto-advance 5.5s, reset on manual nav. Bottom band: stats row (26 Years/721k+/£160M+ with icons + dividers) + AS SEEN ON press logos.

## CategoryNav (`CategoryNav.tsx`) — sticky `top-[100px] z-40`, client
- 7 pills (icon + label) white bar + shadow. Active = orange text + orange underline. Scroll-spy via IntersectionObserver (`rootMargin -45% 0 -50%`); click → smooth scrollIntoView. Mobile: `overflow-x-auto no-scrollbar`.

## WinnersTicker (`WinnersTicker.tsx`) — server, CSS marquee
- "Another winner. Now." mark + horizontal marquee (`.animate-marquee`, list duplicated, pause on hover) of winner cards (thumb + name + prize + green "Won for £…" + pink "Revealed …" pill) + "9,700 winners / in the last 24 hours" with pulsing dot.

## FeaturedSection (`FeaturedSection.tsx`) — server
- Heading + grid: DreamCar (Spot-the-Ball) hero card with 1-2-3 steps strip over dark image + orange DREAM CAR bar + STB graphic + STARTING FROM + PLAY NOW; SubscriberCard (full-image); 5 standard CompetitionCards.

## CompetitionSection (`CompetitionSection.tsx`) — server
- `id` anchor + `scroll-mt-44`, heading (Jost 26–30px uppercase) + subtitle, responsive grid `grid-cols-2 lg:grid-cols-4 gap-4` of CompetitionCards.

## TrustBand / JustLaunched — server, static
- TrustBand: gray bg, "EST. 1999 — £160M+ IN PRIZES" + "Guaranteed winners every week" + DREAM CAR / FIRST VISIT? buttons.
- JustLaunched: image left + "JUST LAUNCHED! / Win cars, bikes, tech or cash!" + ENTER NOW.

## SiteFooter (`SiteFooter.tsx`) — client (mobile accordions)
- Stats strip (Est99/gift/cup + feefo/trustpilot) + description (Winvia Entertainment PLC) + Download Apps badges + 4.8 rating + socials + 6 link columns (collapsible on mobile).

## PromoRibbon (`PromoRibbon.tsx`) — fixed `bottom-0 z-[200000]`, client, dismissible
- Orange→red gradient. "🏡 LAST HOURS… · Only £1 pill · HH:MM:SS live countdown (setInterval, client-only to avoid hydration mismatch) · ENTER NOW · ✕".
