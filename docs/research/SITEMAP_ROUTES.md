# BOTB — Full Site Routes & Workflow

Routes mirror botb.com. Shared shell: SiteHeader + (CategoryNav on home) + SiteFooter + PromoRibbon.

## Pages
| Route | Page | Notes |
|-------|------|-------|
| `/` | Home | done |
| `/competitions` | All competitions listing | grid of all comps + category filter |
| `/prizes/[slug]` | Competition detail | gallery tabs (Tour/Photos/Floorplan/Location), title, cash-alt pill, price, progress, spec grid, sticky ENTER NOW → EnterModal (quantity→basket) |
| `/prizes/cars` | Dream Car selection | steps 1/3 Select Prizes → 2/3 Play game → 3/3 Win. Prize grid (cars/cash, "Double Up"), search/sort. Leads to Spot-the-Ball |
| `/play/[slug]` | Spot-the-Ball game | crosshair placement on action image, then add to basket (step 2/3) |
| `/winners` | Previous Winners | grid of winners w/ filters |
| `/collections/prize-collections` | Prize Collections | |
| `/how-to-play` | How to Play | steps explainer |
| `/botb-pass` | BOTB Pass | subscription tiers |
| `/testimonials` | Testimonials | |
| `/about/history-locations` | History | |
| `/about/charity` | Charity | |
| `/contact-us` | Contact Us | form |
| `/under-the-hood` | Blog | post grid |
| `/account/perks-for-playing` | Perks | |
| `/affiliates` | Affiliate Programme | |
| `/site-map` | Sitemap | link index |
| `/terms` `/privacy` `/cookies` `/complaints` `/mindful-play` | Legal | text content |
| `/login` `/register` | Auth | mock |
| `/cart` | Basket | line items, totals |
| `/checkout` | Checkout | mock details + payment (no real payment) |
| `/checkout/confirmation` | Order confirmation | |
| `/account` | My Account | dashboard (mock) |

## Workflow (mock state via CartContext + localStorage)
1. Home/Competitions → click card → `/prizes/[slug]`.
2. ENTER NOW → EnterModal: Online entry / Free postal entry tabs, quantity (1–1000 slider + quick-select 10/20/50/1000), live total, ADD TO BASKET → cart context.
3. Dream Car: `/prizes/cars` select a car prize → `/play/[slug]` Spot-the-Ball → ADD TO BASKET.
4. Header cart badge shows count → `/cart` → `/checkout` → `/checkout/confirmation`.
5. `/login` or `/register` (mock auth, no backend) → `/account`.

## Detail page data (from /prizes/2house sample)
- Countdown sub-bar "Ends in DD:HH:MM:SS".
- "OR TAKE €X CASH" pill, big price, progress (% sold + tickets X / Y).
- Spec grid: House Value, Monthly Rental, Availability, Cash Alternative (for house); generic for others.
- "What this place offers": bedrooms/bathrooms/sqft (house) — generic specs for cars/tech.
- Sticky bottom ENTER NOW.

Out of scope: real backend, payments, real auth (all mock).
