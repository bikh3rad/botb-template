# BOTB Design Tokens (extracted via getComputedStyle)

## Colors
| Token | Value | Usage |
|-------|-------|-------|
| **Brand orange** | `#FF8200` (rgb 255,130,0) | Primary CTA bg, links, active states, % sold, logo |
| Orange hover | `#E67500` approx (darken ~6%) | button hover |
| Body text | `#333333` | default text |
| Nav/secondary text | `#696971` (rgb 105,105,113) | nav links, cart stroke |
| Muted label | `#898994` (rgb 137,137,148) | TICKET PRICE label, captions |
| Page white | `#FFFFFF` | header, cards |
| Competition area bg | `#EFEFEF` (rgb 239,239,239) | gray section background |
| Card border | `#C3C3C3` @ 0.67px | card outline |
| Black promo bar | `#000000` | top promo bar, hero bg |
| Badge pink/red | ~`#E4002B`→`#E91E63` gradient | ENDS badge pills (sampled from screenshots) |
| Title-bar colors | orange `#FF8200`, teal `#1FA8A0`, purple `#6B2FA0`, red `#D32F2F`, green `#4CAF50` | per-competition colored title strips (approx from screenshots) |

## Typography
- **Jost** — headings, buttons, badges, prices, titles (Google Font). Weights 400/500/600/700.
- **Roboto** — body copy, nav links, labels, descriptions (Google Font). 400/500/600.
- **Josefin Sans** — secondary accents (Google Font).
- Body default: Roboto 16px / line-height 24px / `#333`.

### Measured samples
| Element | font | size | weight | transform | color |
|---|---|---|---|---|---|
| Nav link | Roboto | 16px | 500 | uppercase | #696971 |
| SIGN UP btn | Jost | 14px | 500 | uppercase | #fff on #FF8200, radius 4px, pad 4px 20px, border 2px #FF8200 |
| ENTER TO WIN (hero) | Jost | 16px | 400 | uppercase | #fff on #FF8200, radius 5px, pad 9px 45px 8px 30px |
| ENTER NOW (card, outline) | Roboto | 16px | 400 | uppercase | #FF8200 on #fff, radius 4px, border ~0.67px #FF8200 |
| Card ENDS badge | Jost | 12px | 500 | uppercase | #fff, radius 4px, pad 2px 8px |
| TICKET PRICE label | Roboto | 11px | 500 | uppercase | #898994 |
| % sold | Roboto | 14px | 600 | none | #FF8200 |
| Section heading (FEATURED COMPETITIONS) | Jost | ~28-30px | 600-700 | uppercase | #333 |
| Hero title (LIFESTYLE COMPETITION) | Jost | ~40px+ | 700 | uppercase | #fff |

## Radius / Shadow / Spacing
- Radius: buttons 4-5px, cards `rounded-lg` 8px.
- Card shadow: `0 1.5px 8px rgba(0,0,0,0.2)`.
- Pills bar shadow: `0 1.5px 8px rgba(0,0,0,0.2)`.
- Content max-width `1360px`, padding `px-2 md:px-5`.
- Card grid gap ~16-20px; 4 col desktop / 2 col mobile.

## Map to shadcn tokens (globals.css)
- `--primary` → `#FF8200` (oklch); `--primary-foreground` → white.
- `--background` → white; `--foreground` → `#333`.
- `--muted-foreground` → `#898994`; secondary text `#696971`.
- Custom: `--botb-gray-bg: #EFEFEF`, `--botb-card-border: #C3C3C3`, `--botb-badge: #E4002B`.

## Fonts to load (next/font/google)
- `Jost` (weights 400,500,600,700) → `--font-jost`
- `Roboto` (weights 400,500,600,700) → `--font-roboto` (body default)
- `Josefin_Sans` (weights 400,600,700) → `--font-josefin`
