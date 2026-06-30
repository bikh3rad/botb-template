# BOTB Homepage — Behavior Bible

Observed via scroll/click/hover/responsive sweep.

## Hero carousel (time-driven)
- 8 slides, auto-advance (observed slides change ~every 5-6s without interaction: Lifestyle → £1.2M Home → £2,250 Cash → Defender D350 …).
- Controls: left/right chevron arrows (desktop, vertically centered, white on translucent), and a row of dots (active dot filled/wider).
- Each slide: background webp (2496×787) + overlaid text block (badge, title, subtitle, ENTER TO WIN button). Desktop: text left-aligned. Mobile: text centered, near bottom, full-width button.
- Implementation: React state index + `setInterval` auto-advance (pause on hover optional), crossfade/slide transition.

## Header (sticky)
- `position: sticky; top: 0; z-index: 1003`. Stays on scroll; appearance constant (no shrink observed). Promo bar scrolls away with header as one block (header includes promo bar at top — actually promo bar is topmost; header sticks). Treat promo+nav as a single sticky header for fidelity.

## Category pills (sticky scroll-spy)
- `position: sticky; z-40`, white, shadow `0 1.5px 8px rgba(0,0,0,0.2)`, sits directly under header.
- **Scroll-spy**: active pill updates as you scroll through corresponding section (observed "Featured Competitions" active at top → "Ends Today" active when scrolled into that section → "Free Comps" active near bottom). Active = orange text + orange underline indicator under the item.
- Click a pill → smooth-scroll/jump to that section.
- Mobile: horizontal scrollable row (overflow-x auto), no wrap.

## Winners ticker (time-driven marquee)
- Auto-scrolling horizontal list of recent winner cards (avatar + name + prize + green "Won for £0.0X" + pink "Revealed Xh ago" pill). Continuous marquee animation. Ends with "9,700 winners in the last 24 hours" + pulsing red dot.

## Countdown timers (time-driven)
- Bottom ribbon shows live HOURS : MINUTES : SECONDS, decrementing each second (observed 15:33:33 → 15:27:42 across captures).
- Implementation: `setInterval` 1s, format HH:MM:SS, labels under each number.

## Bottom promo ribbon (fixed, dismissible)
- `position: fixed; bottom: 0; z-index: 200000`. Orange→red gradient. Content: house emoji + "LAST HOURS: Win a £1.2M Home in Zone 1!" + black "Only £1" pill + countdown + "ENTER NOW" button + ✕ close (top-right). Dismiss hides it (state).

## Card hover
- Cards are links; subtle hover (cursor pointer; minor elevation). Buttons: orange fill buttons darken slightly on hover; outline ENTER NOW fills orange on hover (typical pattern). Progress bar static.

## Responsive (breakpoints ~ Tailwind md 768px, lg 1024px)
- **Desktop (1440)**: nav full links; hero text left; pills in one row centered; cards 4-col; footer 5-col links.
- **Tablet (768)**: cards ~2-3 col; pills may scroll.
- **Mobile (390)**: header → hamburger + logo + account + cart; hero text centered/bottom, full-width ENTER TO WIN; pills horizontal-scroll; cards **2-col**; footer link groups become accordions (chevron-down); stats strip stacks.

## Global
- Native scrolling (no Lenis/Locomotive detected).
- Fonts: Jost + Roboto + Josefin Sans (Google).
- Currency localized to viewer (showed €); clone standardizes on **£**.
