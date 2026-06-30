# BOTB Homepage — Page Topology

Target: https://www.botb.com/ (Best of the Best — UK online car competitions).
Tech: Angular SPA + Tailwind utility classes. Total page height ≈ 12,700px desktop @1440.

## Section order (top → bottom)

| # | Section | Component | Interaction | Notes |
|---|---------|-----------|-------------|-------|
| 0 | **Top promo bar** | part of Header | static (dismissible NEW pill) | Black bar, full width: `🎁 FINAL DAYS: New Subscribers get 20 EXTRA Home Tickets!` with red **NEW** pill |
| 1 | **Nav bar** | `SiteHeader` | sticky (z-1003) | White. Desktop: COMPETITIONS · WINNERS · [BOTB logo center] · LOG IN · SIGN UP (orange) · cart. Mobile: hamburger · logo · account · cart |
| 2 | **Hero carousel** | `HeroCarousel` | **time-driven** auto-rotate + arrows + dots | Black bg, 8 full-bleed slides (2496×787 webp). Overlaid: ENDS-badge, title, subtitle, ENTER TO WIN button + cart btn, dots, prev/next arrows. Below slide: stats row (26 Years/721k+/£160M+) + **AS SEEN ON** press logos |
| 3 | **Category pills** | `CategoryNav` | **sticky scroll-spy** (z-40) | White + shadow. 7 items w/ icons: Featured Competitions, Ends Today, Ends Tomorrow, Instant Wins, Ends Soon, Free Comps, Pass Exclusives. Active = orange text + underline. Mobile: horizontal scroll. Clicking jumps to section |
| 4 | **Winners ticker** | `WinnersTicker` | **time-driven** auto-scroll marquee | "Another winner. Now." graphic + cards of recent winners (avatar, name, prize, "Won for £0.02", "Revealed Xh ago" pink pill) + "9,700 winners in the last 24 hours" |
| 5 | **Competition sections** | `CompetitionSection` ×N | static grid (cards have hover) | Gray bg `#EFEFEF`. Each: heading + subtitle + responsive card grid. Sections: **Featured Competitions** (incl. Spot-the-Ball hero card + 1-2-3 how-it-works), **Ends Today**, **Ends Tomorrow**, **Instant Wins**, **Ends Soon** (largest), **Free Comps** |
| 6 | **Trust band** | `TrustBand` | static | Gray. `EST. 1999 — £160M+ IN PRIZES` / `Guaranteed winners every week` + DREAM CAR (orange) + FIRST VISIT? (outline) buttons |
| 7 | **Just Launched band** | `JustLaunched` | static | White card w/ image left, text right: `JUST LAUNCHED!` / `Win cars, bikes, tech or cash!` / ENTER NOW |
| 8 | **Stats strip** | part of Footer | static | 26 years UK's No.1 · Over £160 million in prizes won · Over 721k guaranteed winners · feefo ★★★★★ · Trustpilot |
| 9 | **Footer** | `SiteFooter` | static (accordions on mobile) | Description paragraphs (Winvia Entertainment PLC) + Download Apps badges + 4.8 rating + social icons (FB/IG/YT/TT) + 5 link columns |
| 10 | **Fixed bottom ribbon** | `PromoRibbon` | **time-driven** countdown, dismissible | Orange gradient, z-200000, fixed bottom. `🏡 LAST HOURS: Win a £1.2M Home in Zone 1!` · Only £1 · HOURS:MINUTES:SECONDS live countdown · ENTER NOW · ✕ |

## The Competition Card (central reusable unit)
Appears across all sections 5. Structure:
- White card, `rounded-lg` (8px), border `0.67px #C3C3C3`, shadow `0 1.5px 8px rgba(0,0,0,0.2)`, `overflow-hidden`, ~318px wide in 4-col grid.
- **Image** (342×212 webp) top, `rounded-t-lg`, with:
  - **ENDS badge** top-left (e.g. ENDS TONIGHT/TOMORROW/SUNDAY) — Jost 12px/500 uppercase white on pink/red pill.
- **Colored title bar** (full-width strip below image, e.g. orange/teal/purple/red/green): bold uppercase white title.
- **Body** (centered, padded): description (~16px #333), `TICKET PRICE`/`STARTING FROM` label (Roboto 11px/500 uppercase #898994) + **price** (bold ~24px), progress row (`% sold` orange 14px/600 + tickets `2.5M / 2.7M` with ticket icon), **ENTER NOW** outline button (+ orange cart icon button) OR **PLAY NOW** filled orange button.

Card variants: standard, wide/feature (684×424 image, spans 2 cols), Spot-the-Ball game card (orange PLAY NOW + STB badge), subscriber promo card.

## Layout
- Page scroll container = body (no smooth-scroll library detected — native scroll).
- Content max-width: **1360px** (`max-w-[1360px]`), padded `px-2 md:px-5`.
- Card grid: 4 cols desktop → 2 cols mobile (~768px breakpoint). gap ~ 16-20px.
- z-index layers: ribbon 200000 > header 1003 > pills 40 > card overlays 1.
- Currency: live site geo-shows €; **clone uses £ (GBP)** — the canonical UK brand currency (hero/meta use 35p/50p/£).
