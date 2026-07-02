import type { HeroSlide, CategoryNavItem, FooterColumn } from "@/types";

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
