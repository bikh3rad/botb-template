import type {
  HeroSlide,
  CategoryNavItem,
  Winner,
  CompetitionSection,
  FooterColumn,
} from "@/types";

export const heroSlides: HeroSlide[] = [
  {
    badge: "APP EXCLUSIVE",
    title: "WIN FREE TECH IN APP!",
    subtitle: "Free Apple Airpods, Nintendo Switch 2s & Ninja Creamis",
    image: "/images/hero/slide-13734.webp",
  },
  {
    badge: "ENDS SOON",
    title: "A LUXURY £1.2M HOME IN ZONE 1",
    subtitle: "Last 7 days to win this London home for just £1",
    image: "/images/hero/slide-13344.webp",
  },
  {
    badge: "DREAM CAR COMPETITION",
    title: "THE GOLDEN BOOT!",
    subtitle: "6 cars with gold, who will win the coveted Golden Boot!",
    image: "/images/hero/slide-13728.webp",
  },
  {
    badge: "SUPERCAR",
    title: "WIN A DEFENDER D350 X!",
    subtitle: "The ultimate defender for only 20p!",
    image: "/images/hero/slide-13736.webp",
  },
  {
    badge: "PRIZE EVERY TIME",
    title: "SUMMER FESTIVAL INSTANT WINS!",
    subtitle: "A prize every time to kick-start your summer!",
    image: "/images/hero/slide-13551.webp",
  },
  {
    badge: "OPEN FOR ENTRIES!",
    title: "LIFESTYLE COMPETITION",
    subtitle: "Win cars, cash, tech, watches, and so much more!",
    image: "/images/hero/slide-13982.webp",
  },
  {
    badge: "THE CLOCK'S TICKING",
    title: "WIMBLEDON WONDERS INSTANT WINS!",
    subtitle: "Game, set, match! Epic prizes for £1.19!",
    image: "/images/hero/slide-13548.webp",
  },
  {
    badge: "YOU'VE STILL GOT TIME",
    title: "WIN AN AUDI RS3 CARBON BLACK",
    subtitle: "Drive away in the ultimate hot hatch for 6p!",
    image: "/images/hero/slide-13589.webp",
  },
];

export const heroStats = [
  { value: "26 Years", label: "UK's No.1", est: true },
  { value: "Over 721k+", label: "Winners" },
  { value: "£160M+", label: "in Prizes Won" },
];

export const asSeenOn = [
  { name: "The Sun", image: "/images/as-seen/sun.png" },
  { name: "Mirror", image: "/images/as-seen/mirror.png" },
  { name: "Daily Star", image: "/images/as-seen/daily-star.png" },
  { name: "Express", image: "/images/as-seen/express.png" },
  { name: "Daily Mail", image: "/images/as-seen/daily-mail.png" },
  { name: "Channel 4", image: "/images/as-seen/all4.png" },
  { name: "Sky", image: "/images/as-seen/sky.png" },
];

export const categoryNav: CategoryNavItem[] = [
  { label: "Featured Competitions", icon: "/images/pills/featured.png", targetId: "featured-competitions" },
  { label: "Ends Today", icon: "/images/pills/ends-today.png", targetId: "ends-today" },
  { label: "Ends Tomorrow", icon: "/images/pills/ends-tomorrow.png", targetId: "ends-tomorrow" },
  { label: "Instant Wins", icon: "/images/pills/instant-wins.png", targetId: "instant-wins" },
  { label: "Ends Soon", icon: "/images/pills/ends-soon.png", targetId: "ends-soon" },
  { label: "Free Comps", icon: "/images/pills/free-comps.png", targetId: "free-comps" },
  { label: "Pass Exclusives", icon: "/images/pills/pass-exclusives.png", targetId: "pass-exclusives" },
];

export const winners: Winner[] = [
  { name: "Mark H.", prize: "£500 Cash", wonFor: "£0.02", revealed: "7 hours ago", image: "/images/winners/kfi-1.webp" },
  { name: "Daniel J.", prize: "Oura Ring 5 + Case", wonFor: "£0.02", revealed: "today", image: "/images/winners/kfi-2.webp" },
  { name: "Barrie M.", prize: "Toshiba Fire TV", wonFor: "£0.01", revealed: "today", image: "/images/winners/kfi-3.webp" },
  { name: "Jaisa O.", prize: "750 House Tickets", wonFor: "£0.05", revealed: "today", image: "/images/winners/kfi-4.webp" },
  { name: "Steven P.", prize: "£1,000 Cash", wonFor: "£0.03", revealed: "today", image: "/images/winners/kfi-5.webp" },
  { name: "Laura B.", prize: "iPhone 17 Pro", wonFor: "£0.04", revealed: "yesterday", image: "/images/winners/kfi-6.webp" },
  { name: "Connor M.", prize: "£2,250 Cash", wonFor: "£0.02", revealed: "yesterday", image: "/images/winners/kfi-7.webp" },
  { name: "Priya S.", prize: "MacBook Neo", wonFor: "£0.05", revealed: "yesterday", image: "/images/winners/kfi-8.webp" },
  { name: "Gary T.", prize: "Nintendo Bundle", wonFor: "£0.03", revealed: "yesterday", image: "/images/winners/kfi-9.webp" },
  { name: "Emma W.", prize: "£800 Cash", wonFor: "£0.01", revealed: "yesterday", image: "/images/winners/kfi-10.webp" },
  { name: "Tom R.", prize: "10g Gold Bar", wonFor: "£0.02", revealed: "2 days ago", image: "/images/winners/kfi-11.webp" },
  { name: "Sofia L.", prize: "Mystery Tech Prize", wonFor: "£0.04", revealed: "2 days ago", image: "/images/winners/kfi-12.webp" },
];

export const winnersCount = "9,700 winners";

/** Cards shown in the FEATURED COMPETITIONS section (excludes the Spot-the-Ball hero card). */
export const featuredCompetitions: import("@/types").Competition[] = [
  { badge: "ENDS TONIGHT", title: "£1.2M HOME IN ZONE 1", titleColor: "teal", description: "Time's Nearly Up! Win This Home in Central London for Just £1", priceLabel: "TICKET PRICE", price: "£1.25", stat: { label: "92.5% sold", percent: 93, tickets: "2.5M / 2.7M" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13344.webp" },
  { badge: "ENDS TOMORROW", title: "£500K+ INSTANT WINS", titleColor: "purple", description: "The clock's ticking! Wimbledon Wonders - instant prizes for just £1.19!", priceLabel: "TICKET PRICE", price: "£1.49", stat: { label: "27.62k Won", percent: 60, tickets: "170,326 Left" }, cta: "ENTER NOW", image: "/images/comps/13548.webp" },
  { badge: "ENDS FRIDAY", title: "AUDI R8 FOR 21P!", titleColor: "red", description: "Bag this German-engineered supercar icon worth £76k!", priceLabel: "TICKET PRICE", price: "£0.27", stat: { label: "62% sold", percent: 62, tickets: "12,400 / 20,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13558.webp" },
  { badge: "ENDS SUNDAY", title: "£2.5M+ INSTANT WINS", titleColor: "purple", description: "BOTB's summer festival is here, grab entry for just £1.49!", priceLabel: "TICKET PRICE", price: "£1.87", stat: { label: "169.7k Won", percent: 70, tickets: "1.2M Left" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13551.webp" },
  { badge: "ENDS IN 11 DAYS", title: "EVOQUE + MINI FOR 9P!", titleColor: "red", description: "Bag a Range Rover and Mini worth £52k!", priceLabel: "TICKET PRICE", price: "£0.12", stat: { label: "18% sold", percent: 18, tickets: "3,600 / 20,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13740.webp" },
];

/** The Spot-the-Ball "Dream Car" hero card (first cell of Featured). */
export const dreamCarCard = {
  steps: ["Select prize", "Play the game", "Win your dream car"],
  titleBar: "DREAM CAR",
  description: "Final Hours to Enter - Win a car PLUS gold!",
  priceLabel: "STARTING FROM",
  price: "£1.13",
  badge: "ENDS SUNDAY",
  heroImage: "/images/comps/13732-wide.webp",
  stbBadge: "/images/comps/stb-badge.png",
};

/** New-subscriber promo card in Featured. */
export const subscriberCard = {
  badge: "ENDS TONIGHT",
  image: "/images/comps/fallback-subscribers.webp",
  headline: "NEW SUBSCRIBERS GET 20 EXTRA LONDON HOME TICKETS",
};

export const competitionSections: CompetitionSection[] = [
  {
    id: "ends-today",
    heading: "ENDS TODAY",
    subtitle: "Your last chance to enter, don't miss out!",
    competitions: [
      { badge: "ENDS TONIGHT", title: "IPHONE 17 & 1,249+ PRIZES!", titleColor: "purple", description: "1,250 prizes up for grabs – including an iPhone 17 & £500!", cta: "DETAILS", image: "/images/comps/11196.webp" },
      { badge: "ENDS TONIGHT", title: "£1.2M HOME IN ZONE 1", titleColor: "teal", description: "Time's Nearly Up! Win This Home in Central London for Just £1", priceLabel: "TICKET PRICE", price: "£1.25", stat: { label: "93% sold", percent: 93, tickets: "2.5M / 2.7M" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13344.webp" },
      { badge: "ENDS TONIGHT", title: "1K HOUSE TICKETS!", titleColor: "green", description: "Want to boost your odds to win a house for 25p?", priceLabel: "TICKET PRICE", price: "£0.32", stat: { label: "48% sold", percent: 48, tickets: "3,924 / 8,499" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13831-v1.webp" },
      { badge: "ENDS TONIGHT", title: "500 HOUSE TICKETS!", titleColor: "green", description: "Low odds to bag 500 house tickets!", priceLabel: "TICKET PRICE", price: "£0.62", stat: { label: "41% sold", percent: 41, tickets: "772 / 1,899" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13831-v1.webp" },
      { badge: "ENDS TONIGHT", title: "RATTAN DINING SET", titleColor: "green", description: "Transform your garden for 5p!", priceLabel: "TICKET PRICE", price: "£0.07", stat: { label: "47% sold", percent: 47, tickets: "4,512 / 9,600" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13095.webp" },
      { badge: "ENDS TONIGHT", title: "MYSTERY CASH PRIZE", titleColor: "green", description: "Unlock Today's Mystery Cash Prize!", priceLabel: "TICKET PRICE", price: "£0.24", stat: { label: "51% sold", percent: 51, tickets: "5,100 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/11121.webp" },
      { badge: "ENDS TONIGHT", title: "NINJA AUTOBARISTA PRO", titleColor: "green", description: "Café-quality coffee at home for just 5p!", priceLabel: "TICKET PRICE", price: "£0.07", stat: { label: "31% sold", percent: 31, tickets: "3,100 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13763.webp" },
      { badge: "ENDS TONIGHT", title: "1000 HOUSE TICKETS 1P", titleColor: "green", description: "1,000 House Tickets for just 1p!", priceLabel: "TICKET PRICE", price: "£0.02", stat: { label: "89% sold", percent: 89, tickets: "8,900 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13831-v1.webp" },
    ],
  },
  {
    id: "ends-tomorrow",
    heading: "ENDS TOMORROW",
    subtitle: "Last chance to enter these competitions!",
    competitions: [
      { badge: "ENDS TOMORROW", title: "IPHONE 17 & 1,249+ PRIZES!", titleColor: "purple", description: "1,250 prizes up for grabs – including an iPhone 17 & £500!", cta: "DETAILS", image: "/images/comps/11196.webp" },
      { badge: "ENDS TOMORROW", title: "5G GOLD BAR", titleColor: "green", description: "Own a 5g gold bar for just 5p!", priceLabel: "TICKET PRICE", price: "£0.07", stat: { label: "43% sold", percent: 43, tickets: "4,300 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/9098-gold.webp" },
      { badge: "ENDS TOMORROW", title: "£2,250 CASH", titleColor: "green", description: "Just 5p for £2.25K!", priceLabel: "TICKET PRICE", price: "£0.07", stat: { label: "49% sold", percent: 49, tickets: "5,766 / 12,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/10287-2250.webp" },
      { badge: "ENDS TOMORROW", title: "IPHONE 17 PRO MAX", titleColor: "green", description: "29p for the best iPhone on the market!", priceLabel: "TICKET PRICE", price: "£0.37", stat: { label: "29.5% sold", percent: 30, tickets: "2,950 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/9275-iphone.webp" },
      { badge: "ENDS TOMORROW", title: "500 HOUSE TICKETS!", titleColor: "green", description: "Low odds to bag 500 house tickets!", priceLabel: "TICKET PRICE", price: "£0.62", stat: { label: "28.5% sold", percent: 29, tickets: "540 / 1,899" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13831-v1.webp" },
      { badge: "ENDS TOMORROW", title: "£500K+ INSTANT WINS", titleColor: "purple", description: "The clock's ticking! Wimbledon Wonders - instant prizes for just £1.19!", priceLabel: "TICKET PRICE", price: "£1.19", stat: { label: "27.6k Won", percent: 60, tickets: "170,326 Left" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13548.webp" },
    ],
  },
  {
    id: "instant-wins",
    heading: "INSTANT WINS",
    subtitle: "Play now to win instantly!",
    competitions: [
      { badge: "ENDS TOMORROW", title: "£500K+ INSTANT WINS", titleColor: "purple", description: "The clock's ticking! Wimbledon Wonders - instant prizes for just £1.19!", priceLabel: "TICKET PRICE", price: "£1.19", stat: { label: "27.6k Won", percent: 60, tickets: "170,326 Left" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13548.webp" },
      { badge: "ENDS SUNDAY", title: "£2.5M+ INSTANT WINS", titleColor: "purple", description: "BOTB's summer festival is here, grab entry for just £1.49!", priceLabel: "TICKET PRICE", price: "£1.49", stat: { label: "169.7k Won", percent: 70, tickets: "1.2M Left" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13551.webp" },
      { badge: "ENDS WEDNESDAY", title: "£500K+ INSTANT WINS", titleColor: "purple", description: "Win big instantly with Summer Festival prizes for 99p!", priceLabel: "TICKET PRICE", price: "£0.99", stat: { label: "42.1k Won", percent: 45, tickets: "320,500 Left" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13734.webp" },
    ],
  },
  {
    id: "ends-soon",
    heading: "ENDS SOON",
    subtitle: "These competitions won't be around for long.",
    competitions: [
      { badge: "ENDS THURSDAY", title: "IPHONE 17 & 1,249+ PRIZES!", titleColor: "purple", description: "1,250 prizes up for grabs – including an iPhone 17 & £500!", cta: "DETAILS", image: "/images/comps/11196.webp" },
      { badge: "ENDS THURSDAY", title: "SAMSUNG GALAXY BOOK6", titleColor: "green", description: "Redefine productivity with the Galaxy Book6!", priceLabel: "TICKET PRICE", price: "£0.05", stat: { label: "10% sold", percent: 10, tickets: "1,000 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13767.webp" },
      { badge: "ENDS THURSDAY", title: "MACBOOK NEO", titleColor: "green", description: "Win Apple's latest tech for 4p!", priceLabel: "TICKET PRICE", price: "£0.05", stat: { label: "28.5% sold", percent: 29, tickets: "2,850 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/11668.webp" },
      { badge: "ENDS THURSDAY", title: "TOSHIBA TV", titleColor: "green", description: 'Upgrade your TV with the Toshiba 70" Ultra QLED!', priceLabel: "TICKET PRICE", price: "£0.02", stat: { label: "39.5% sold", percent: 40, tickets: "3,950 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/10106-toshiba.webp" },
      { badge: "ENDS THURSDAY", title: "750 HOUSE TICKETS!", titleColor: "green", description: "750 entries for our House Competition?", priceLabel: "TICKET PRICE", price: "£0.02", stat: { label: "24% sold", percent: 24, tickets: "456 / 1,899" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13831-v1.webp" },
      { badge: "ENDS THURSDAY", title: "1K HOUSE TICKETS!", titleColor: "green", description: "Want to boost your odds to win a house for 25p?", priceLabel: "TICKET PRICE", price: "£0.32", stat: { label: "9.5% sold", percent: 10, tickets: "807 / 8,499" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13831-v1.webp" },
      { badge: "ENDS FRIDAY", title: "AUDI R8 FOR 21P!", titleColor: "red", description: "Bag this German-engineered supercar icon worth £76k!", priceLabel: "TICKET PRICE", price: "£0.21", stat: { label: "62% sold", percent: 62, tickets: "12,400 / 20,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13558.webp" },
      { badge: "ENDS FRIDAY", title: "NINTENDO BUNDLE", titleColor: "green", description: "Win this Nintendo Switch 2 & Mario Kart Bundle!", priceLabel: "TICKET PRICE", price: "£0.37", stat: { label: "19.5% sold", percent: 20, tickets: "1,950 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/8941-nintendo.webp" },
      { badge: "ENDS FRIDAY", title: "MYSTERY LIFESTYLE", titleColor: "green", description: "A mystery lifestyle prize could be yours!", priceLabel: "TICKET PRICE", price: "£0.25", stat: { label: "29.5% sold", percent: 30, tickets: "2,950 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/11402.webp" },
      { badge: "ENDS FRIDAY", title: "£1,250 CASH", titleColor: "green", description: "1p for a bundle of tax-free cash!", priceLabel: "TICKET PRICE", price: "£0.02", stat: { label: "35.5% sold", percent: 36, tickets: "3,550 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/9211-money.webp" },
      { badge: "ENDS SATURDAY", title: "£5,000 CASH", titleColor: "green", description: "Just 3p could land you £5,000 tax-free!", priceLabel: "TICKET PRICE", price: "£0.05", stat: { label: "15.5% sold", percent: 16, tickets: "1,550 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/8941-1000cash.webp" },
      { badge: "ENDS SATURDAY", title: "£1K AMAZON VOUCHER", titleColor: "green", description: "For 10p this £1,000 Amazon Voucher could be yours!", priceLabel: "TICKET PRICE", price: "£0.13", stat: { label: "8% sold", percent: 8, tickets: "800 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/12931.webp" },
      { badge: "ENDS SATURDAY", title: "PUREMATE AC", titleColor: "green", description: "Beat the heat with the PureMate AC for 5p!", priceLabel: "TICKET PRICE", price: "£0.07", stat: { label: "19% sold", percent: 19, tickets: "1,900 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13932.webp" },
      { badge: "ENDS SATURDAY", title: "1,000 HOUSE TICKETS FOR 1P!", titleColor: "green", description: "1p to win your dream home?!", priceLabel: "TICKET PRICE", price: "£0.02", stat: { label: "9.5% sold", percent: 10, tickets: "950 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13831-v1.webp" },
      { badge: "ENDS SUNDAY", title: "ULTIMATE BOTB PASS", titleColor: "green", description: "Better your odds with a year of our Ultimate BOTB Pass!", priceLabel: "TICKET PRICE", price: "£0.13", stat: { label: "14% sold", percent: 14, tickets: "1,400 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/11364.webp" },
      { badge: "ENDS SUNDAY", title: "MYSTERY TECH PRIZE", titleColor: "green", description: "A mystery tech prize could be yours!", priceLabel: "TICKET PRICE", price: "£0.19", stat: { label: "7% sold", percent: 7, tickets: "700 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/10873.webp" },
      { badge: "ENDS SUNDAY", title: "£1,750 CASH", titleColor: "green", description: "Enter now for your chance to win £1,750 Cash!", priceLabel: "TICKET PRICE", price: "£0.11", stat: { label: "10.5% sold", percent: 11, tickets: "1,050 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/10168.webp" },
      { badge: "ENDS SUNDAY", title: "DREAM CAR", titleColor: "orange", description: "Final Hours to Enter - Win a car PLUS gold!", priceLabel: "STARTING FROM", price: "£1.13", cta: "PLAY NOW", image: "/images/comps/13732.webp" },
      { badge: "ENDS MONDAY", title: "10G GOLD BAR", titleColor: "green", description: "Luxury Gold, Crazy Price – Just 1p!", priceLabel: "TICKET PRICE", price: "£0.02", stat: { label: "2% sold", percent: 2, tickets: "200 / 10,000" }, cta: "ENTER NOW", showCart: true, wide: true, image: "/images/comps/11858.webp" },
      { badge: "ENDS MONDAY", title: "SHARK FAN BUNDLE", titleColor: "green", description: "Stay cool this summer – Shark Fan Bundle for 3p", priceLabel: "TICKET PRICE", price: "£0.04", stat: { label: "0.5% sold", percent: 1, tickets: "50 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13936.webp" },
      { badge: "ENDS MONDAY", title: "£800 CASH", titleColor: "green", description: "Bank balance looking low? Win £800 for 1p", priceLabel: "TICKET PRICE", price: "£0.02", stat: { label: "3% sold", percent: 3, tickets: "300 / 10,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/10415.webp" },
      { badge: "ENDS MONDAY", title: "LIFESTYLE COMPETITION", titleColor: "orange", description: "Open For Entries! Your chance to win cars, watches, tech & more!", priceLabel: "TICKET PRICE", price: "£0.44", stat: { label: "5% sold", percent: 5, tickets: "500 / 10,000" }, cta: "ENTER NOW", showCart: true, wide: true, image: "/images/comps/13982.webp" },
      { badge: "ENDS IN 11 DAYS", title: "EVOQUE + MINI FOR 9P!", titleColor: "red", description: "Bag a Range Rover and Mini worth £52k!", priceLabel: "TICKET PRICE", price: "£0.12", stat: { label: "18% sold", percent: 18, tickets: "3,600 / 20,000" }, cta: "ENTER NOW", showCart: true, image: "/images/comps/13740.webp" },
    ],
  },
  {
    id: "free-comps",
    heading: "FREE COMPS",
    subtitle: "Enter our free competitions — no purchase necessary!",
    competitions: [
      { badge: "ENDS TONIGHT", title: "FREE WORLD CUP TICKETS!", titleColor: "blue", description: "Win World Cup tickets for free — simply enter!", cta: "ENTER NOW", image: "/images/comps/13883.webp" },
      { badge: "ENDS FRIDAY", title: "FREE £250 CASH", titleColor: "blue", description: "Enter for free for your chance to win £250 cash!", cta: "ENTER NOW", image: "/images/comps/9211-money.webp" },
      { badge: "ENDS SUNDAY", title: "FREE MYSTERY PRIZE", titleColor: "blue", description: "A free mystery prize could be yours this week!", cta: "ENTER NOW", image: "/images/comps/10873.webp" },
    ],
  },
];

export const footerColumns: FooterColumn[] = [
  {
    title: "All competitons",
    links: [
      { label: "Dream Car", href: "/prizes/cars" },
      { label: "£1.2M Home in Zone 1", href: "/prizes/1-2m-home-in-zone-1" },
      { label: "£500k+ Instant Wins", href: "/prizes/500k-instant-wins" },
      { label: "Audi R8 for 21p!", href: "/prizes/audi-r8-for-21p" },
      { label: "£2.5M+ Instant Wins", href: "/prizes/2-5m-instant-wins" },
      { label: "Evoque + Mini for 9p!", href: "/prizes/evoque-mini-for-9p" },
      { label: "Free World Cup Tickets!", href: "/prizes/free-world-cup-tickets" },
      { label: "Lifestyle Competition", href: "/prizes/lifestyle-competition" },
    ],
  },
  {
    title: "Categories",
    collapsible: true,
    links: [
      { label: "Featured Competitions", href: "/#featured-competitions" },
      { label: "Ends Today", href: "/#ends-today" },
      { label: "Ends Tomorrow", href: "/#ends-tomorrow" },
      { label: "Instant Wins", href: "/#instant-wins" },
      { label: "Ends Soon", href: "/#ends-soon" },
      { label: "Free Comps", href: "/#free-comps" },
      { label: "All Competitions", href: "/competitions" },
    ],
  },
  {
    title: "Winners",
    links: [
      { label: "Previous Winners", href: "/winners" },
      { label: "Prize Collections", href: "/collections/prize-collections" },
    ],
  },
  {
    title: "About BOTB",
    links: [
      { label: "Testimonials", href: "/testimonials" },
      { label: "History", href: "/about/history-locations" },
      { label: "How to Play", href: "/how-to-play" },
      { label: "Blog", href: "/under-the-hood" },
      { label: "Contact Us", href: "/contact-us" },
      { label: "Charity", href: "/about/charity" },
      { label: "Sitemap", href: "/site-map" },
    ],
  },
  {
    title: "Under the Hood",
    links: [
      { label: "Terms of Play", href: "/terms" },
      { label: "Privacy Policy", href: "/privacy" },
      { label: "Mindful Play", href: "/mindful-play" },
      { label: "Cookie Policy", href: "/cookies" },
      { label: "Complaints Policy", href: "/complaints" },
      { label: "Modern Slavery", href: "/terms" },
    ],
  },
  {
    title: "More from BOTB",
    links: [
      { label: "Perks For Playing", href: "/account/perks-for-playing" },
      { label: "iOS App", href: "/how-to-play" },
      { label: "Android App", href: "/how-to-play" },
      { label: "Affiliate Programme", href: "/affiliates" },
      { label: "BOTB Pass", href: "/botb-pass" },
    ],
  },
];

export const footerDescription = [
  "Winvia Entertainment PLC (formerly known as Best of the Best Limited) operates skilled prize competitions resulting in the allocation of prizes in accordance with the Terms and Conditions of the website.",
  "Win a brand-new car, take the cash alternative, or win Competition Credit in the BOTB Dream Car Competition. There are over 150 new car prizes to choose from, and the closest person in the skilled Spot the Ball game wins the car or a life-changing amount of cash!",
  "And don't miss out on our Instant Win, Lifestyle, Luxury car, House and Free Competitions to win cars, motorbikes, holidays, watches, tech, cash, and more life-changing prizes. Simply choose your tickets and check out.",
];
