// Downloads all BOTB homepage assets into public/images/, preserving a meaningful structure.
// Usage: node scripts/download-assets.mjs
import { mkdir, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";

const ROOT = new URL("../public/images/", import.meta.url).pathname;

// Map each remote URL to a local path (relative to public/images/).
const ASSETS = [
  // --- Hero slider (desktop) ---
  ["https://cdn.botb.com/media/xlbd5se4/botb-13344-slide-desktop-2496-787-3.webp", "hero/slide-13344.webp"],
  ["https://cdn.botb.com/media/2lunnupg/botb-13728-slider-desktop-2496-787.webp", "hero/slide-13728.webp"],
  ["https://cdn.botb.com/media/j4jpq2wn/botb-13736-slide-desktop-2496-787.webp", "hero/slide-13736.webp"],
  ["https://cdn.botb.com/media/5zmnkd4t/botb-13551-slide-desktop-2496-787.webp", "hero/slide-13551.webp"],
  ["https://cdn.botb.com/media/g4xanpav/botb-13982-slider-desktop-2496-787.webp", "hero/slide-13982.webp"],
  ["https://cdn.botb.com/media/01wfxefq/botb-13548-botb-13548-slide-desktop-2496-787.webp", "hero/slide-13548.webp"],
  ["https://cdn.botb.com/media/mpndkht5/botb-13589-slider-desktop-2496-787.webp", "hero/slide-13589.webp"],
  ["https://cdn.botb.com/media/f10lcyvh/botb-13734-slider-desktop-2496-787.webp", "hero/slide-13734.webp"],

  // --- As seen on (press logos) ---
  ["https://www.botb.com/assets/images/as-seen-logos/sunmasthead.png", "as-seen/sun.png"],
  ["https://www.botb.com/assets/images/as-seen-logos/the-mirror.png", "as-seen/mirror.png"],
  ["https://www.botb.com/assets/images/as-seen-logos/masthead-dailystar.png", "as-seen/daily-star.png"],
  ["https://www.botb.com/assets/images/as-seen-logos/express.png", "as-seen/express.png"],
  ["https://www.botb.com/assets/images/as-seen-logos/daily_mail.png", "as-seen/daily-mail.png"],
  ["https://www.botb.com/assets/images/as-seen-logos/all4_logo.png", "as-seen/all4.png"],
  ["https://www.botb.com/assets/images/as-seen-logos/sky_group.png", "as-seen/sky.png"],

  // --- Category pill icons ---
  ["https://cdn.botb.com/media/bxjnajmj/featuredicon.png", "pills/featured.png"],
  ["https://cdn.botb.com/media/i1kniy5e/endstodayicon.png", "pills/ends-today.png"],
  ["https://cdn.botb.com/media/zloopari/endstomorrowicon.png", "pills/ends-tomorrow.png"],
  ["https://cdn.botb.com/media/zpxb2xen/instantwinsicon.png", "pills/instant-wins.png"],
  ["https://cdn.botb.com/media/dyxe3by5/endssoonicon.png", "pills/ends-soon.png"],
  ["https://cdn.botb.com/media/5lcmdivu/freecompetitionicon.png", "pills/free-comps.png"],
  ["https://cdn.botb.com/media/kfqglore/image-989-1.png", "pills/pass-exclusives.png"],

  // --- Winners ticker title + thumbs ---
  ["https://www.botb.com/assets/images/win/win-widget-title.svg", "winners/another-winner.svg"],
  ["https://cdn.botb.com/media/etbbvzif/kfi.webp", "winners/kfi-1.webp"],
  ["https://cdn.botb.com/media/cjqdqaaj/kfi.webp", "winners/kfi-2.webp"],
  ["https://cdn.botb.com/media/vhbnhvpq/kfi.webp", "winners/kfi-3.webp"],
  ["https://cdn.botb.com/media/nunfxyfy/kfi.webp", "winners/kfi-4.webp"],
  ["https://cdn.botb.com/media/reud0fvu/kfi.webp", "winners/kfi-5.webp"],
  ["https://cdn.botb.com/media/1bebnhnx/kfi-5.webp", "winners/kfi-6.webp"],
  ["https://cdn.botb.com/media/dykammax/kfi.webp", "winners/kfi-7.webp"],
  ["https://cdn.botb.com/media/iwuf12il/kfi.webp", "winners/kfi-8.webp"],
  ["https://cdn.botb.com/media/mdehgztw/kfi.webp", "winners/kfi-9.webp"],
  ["https://cdn.botb.com/media/o25dkcbb/kfi.webp", "winners/kfi-10.webp"],
  ["https://cdn.botb.com/media/t5cdl0lp/kfi.webp", "winners/kfi-11.webp"],
  ["https://cdn.botb.com/media/mccldb52/kfi.webp", "winners/kfi-12.webp"],

  // --- Competition card images ---
  ["https://cdn.botb.com/media/deklnbjo/botb-8934-thumbnail-iphone-17.webp", "comps/iphone-17.webp"],
  ["https://cdn.botb.com/media/4eincllc/botb-13643-thumbnail.png", "comps/13643.png"],
  ["https://cdn.botb.com/media/dj0ljvke/botb-10106-thumbnail-toshiba.webp", "comps/toshiba.webp"],
  ["https://cdn.botb.com/media/bzlbvgld/botb-13345-comp-thumb-v1-750.png", "comps/13345-750.png"],
  ["https://cdn.botb.com/media/x3djygid/botb-13345-comp-thumb-v1-1000.png", "comps/13345-1000.png"],
  ["https://cdn.botb.com/media/ozioeuoy/botb-13732-competition-card-659x201-x-2.webp", "comps/13732-wide.webp"],
  ["https://www.botb.com/assets/images/stb-competition-cta-badge.png", "comps/stb-badge.png"],
  ["https://cdn.botb.com/media/f1ybr14p/botb-13344-competition-card-342x212.webp", "comps/13344.webp"],
  ["https://cdn.botb.com/media/fohfobpt/botb-13548-botb-13548-competition-card-342x212.webp", "comps/13548.webp"],
  ["https://cdn.botb.com/media/ejvhhhb3/botb-13558-competition-card-342x212.webp", "comps/13558.webp"],
  ["https://cdn.botb.com/media/uwike3g3/botb-13551-competition-card-342x212.webp", "comps/13551.webp"],
  ["https://cdn.botb.com/media/t4gmakep/botb-13740-competition-card-342x212.webp", "comps/13740.webp"],
  ["https://cdn.botb.com/media/u25cg1er/botb-13883-botb-13883-competition-card-342x212.webp", "comps/13883.webp"],
  ["https://cdn.botb.com/media/grhj0mgy/botb-13982-competition-card-684x424.webp", "comps/13982.webp"],
  ["https://cdn.botb.com/media/hvrnsnbm/botb-11196-comp-card_1.webp", "comps/11196.webp"],
  ["https://www.botb.com/assets/images/subscription-competitions/exclusive-badge.webp", "comps/exclusive-badge.webp"],
  ["https://cdn.botb.com/media/xwie15kw/botb-13831-competition-card-342x212-v1-1000.webp", "comps/13831-v1.webp"],
  ["https://cdn.botb.com/media/vyhhg5yt/botb-13095-comp-card_.webp", "comps/13095.webp"],
  ["https://cdn.botb.com/media/ot3nygab/botb-11121-comp-card.webp", "comps/11121.webp"],
  ["https://cdn.botb.com/media/hohmf3am/botb-13763-comp-card_.webp", "comps/13763.webp"],
  ["https://cdn.botb.com/media/djnbheiu/botb-9098-competition_card_template-5g-gold-bar.webp", "comps/9098-gold.webp"],
  ["https://cdn.botb.com/media/1imjywte/botb-10287-comp-card-2250.webp", "comps/10287-2250.webp"],
  ["https://cdn.botb.com/media/dvkfyqko/botb-9275-comp-card-iphone-17-pro-max.webp", "comps/9275-iphone.webp"],
  ["https://cdn.botb.com/media/w4kbypud/botb-13767-competition-card-342x212.webp", "comps/13767.webp"],
  ["https://cdn.botb.com/media/n2ceh0eo/botb-11668-comp-card_.webp", "comps/11668.webp"],
  ["https://cdn.botb.com/media/ak0hoxyz/botb-10106-comp-card-toshiba.webp", "comps/10106-toshiba.webp"],
  ["https://cdn.botb.com/media/snvjgwu1/botb-8941-competition_card_template-nintendo.webp", "comps/8941-nintendo.webp"],
  ["https://cdn.botb.com/media/lxkhrpl1/botb-11402-comp-card.webp", "comps/11402.webp"],
  ["https://cdn.botb.com/media/deplhctd/botb-9211-comp-card-money.webp", "comps/9211-money.webp"],
  ["https://cdn.botb.com/media/ai2a1obh/botb-8941-competition_card_template-1-000-cash.webp", "comps/8941-1000cash.webp"],
  ["https://cdn.botb.com/media/ysfjbslw/botb-12931-comp-card.webp", "comps/12931.webp"],
  ["https://cdn.botb.com/media/wlrkqe4d/botb-13932-comp-card_.webp", "comps/13932.webp"],
  ["https://cdn.botb.com/media/kf3ft5sc/botb-11364-comp-card_.webp", "comps/11364.webp"],
  ["https://cdn.botb.com/media/u5qon3gl/botb-10873-comp-card.webp", "comps/10873.webp"],
  ["https://cdn.botb.com/media/gpwdgwbp/botb-10168-comp-card_.webp", "comps/10168.webp"],
  ["https://cdn.botb.com/media/45xfr4gy/botb-13732-competition-card-342x212.webp", "comps/13732.webp"],
  ["https://cdn.botb.com/media/5qrft215/botb-11858-competition-card-684x424-copy.webp", "comps/11858.webp"],
  ["https://cdn.botb.com/media/wh4kulmy/botb-13936-comp-card_.webp", "comps/13936.webp"],
  ["https://cdn.botb.com/media/5sshkvxm/botb-10415-comp-card.webp", "comps/10415.webp"],
  ["https://cdn.botb.com/media/gn3fgkym/botb-13734-competition-card-342x212.webp", "comps/13734.webp"],

  // --- Featured "Just Launched" image ---
  ["https://cdn.botb.com/media/bwcjclci/homepage-ft-image-desktop.webp", "misc/just-launched.webp"],
  ["https://www.botb.com/assets/images/win/fallback-subscribers.webp", "comps/fallback-subscribers.webp"],
  ["https://www.botb.com/assets/images/win/fallback-easy-wins.webp", "comps/fallback-easy-wins.webp"],
  ["https://cdn.botb.com/media/acodhhhn/29-500.png", "comps/29-500.png"],

  // --- Footer ---
  ["https://cdn.botb.com/media/4vojnlll/footer-est99-icon.png", "footer/est99.png"],
  ["https://cdn.botb.com/media/getonzu2/footer-tabler-icon.png", "footer/gift.png"],
  ["https://cdn.botb.com/media/ilnn03ry/footer-cup-icon.png", "footer/cup.png"],
  ["https://cdn.botb.com/media/gq2hjjbm/feefo-footer.png", "footer/feefo.png"],
  ["https://cdn.botb.com/media/bsbdo2dr/trustpilot-footer.png", "footer/trustpilot.png"],
  ["https://www.botb.com/assets/images/apple-store-badge2.png", "footer/app-store.png"],
  ["https://www.botb.com/assets/images/google-play-store-badge.png", "footer/google-play.png"],
  ["https://www.botb.com/assets/images/socials/fb-logo.png", "footer/fb.png"],
  ["https://www.botb.com/assets/images/socials/ig-logo.png", "footer/ig.png"],
  ["https://www.botb.com/assets/images/socials/yt-logo.png", "footer/yt.png"],
  ["https://www.botb.com/assets/images/socials/tt-logo.png", "footer/tt.png"],

  // --- Misc UI ---
  ["https://www.botb.com/assets/images/chevron.png", "ui/chevron.png"],
  ["https://www.botb.com/assets/images/tabler-icon-chevron-down.svg", "ui/chevron-down.svg"],
  ["https://www.botb.com/assets/images/modal-close-icon.png", "ui/close.png"],
  ["https://www.botb.com/assets/images/icon-minus-black.png", "ui/minus.png"],
  ["https://www.botb.com/assets/images/icon-plus-black.png", "ui/plus.png"],

  // --- OG / social share ---
  ["https://www.botb.com/assets/images/botb-wheel-social.png", "../seo/og-image.png"],
];

async function download([url, rel]) {
  const dest = join(ROOT, rel);
  await mkdir(dirname(dest), { recursive: true });
  try {
    const res = await fetch(url, { headers: { "User-Agent": "Mozilla/5.0", Referer: "https://www.botb.com/" } });
    if (!res.ok) return { url, rel, ok: false, status: res.status };
    const buf = Buffer.from(await res.arrayBuffer());
    await writeFile(dest, buf);
    return { url, rel, ok: true, bytes: buf.length };
  } catch (e) {
    return { url, rel, ok: false, error: String(e) };
  }
}

async function run() {
  const results = [];
  for (let i = 0; i < ASSETS.length; i += 4) {
    const batch = ASSETS.slice(i, i + 4);
    results.push(...(await Promise.all(batch.map(download))));
  }
  const ok = results.filter((r) => r.ok);
  const fail = results.filter((r) => !r.ok);
  console.log(`Downloaded ${ok.length}/${results.length} assets`);
  if (fail.length) console.log("FAILED:\n" + fail.map((f) => `  ${f.status || f.error} ${f.url}`).join("\n"));
}

run();
