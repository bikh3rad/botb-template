import Link from "next/link";
import { notFound } from "next/navigation";
import { getCompetitions, getCompetitionBySlug } from "@/lib/api";
import { competitionImage, toDetailView } from "@/lib/presentation";
import type { TitleBarColor } from "@/types";
import { TicketIcon } from "@/components/icons";
import { EnterEntry } from "@/components/EnterEntry";

const TITLE_BAR_BG: Record<TitleBarColor, string> = {
  orange: "bg-botb-orange",
  teal: "bg-botb-teal",
  purple: "bg-botb-purple",
  red: "bg-botb-red",
  green: "bg-botb-green",
  blue: "bg-botb-blue",
  pass: "bg-botb-purple",
};

/** Visual-only gallery tabs mirroring botb.com detail pages. */
const GALLERY_TABS = ["Tour", "Photos", "Floorplan", "Location"] as const;

// Rendered at request time: slugs live in the backend, so we resolve on demand
// rather than pre-generating params (which would require a backend at build).
export const dynamic = "force-dynamic";

export default async function Page({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const competition = await getCompetitionBySlug(slug);

  if (!competition) {
    notFound();
  }

  const c = toDetailView(competition);
  const displayTitle = c.description;
  const price = c.price ?? "£0.00";
  const priceLabel = c.priceLabel ?? "TICKET PRICE";

  // Generic spec grid derived from the competition data.
  const specs: { label: string; value: string }[] = [
    { label: "Prize Value", value: c.value },
    { label: "Cash Alternative", value: "Available" },
    { label: "Draw", value: "Live on Facebook" },
    { label: "Tickets from", value: price },
  ];

  // A handful of other live competitions for the bottom strip.
  const others = await getCompetitions("live");
  const more = others
    .filter((x) => x.slug !== competition.slug)
    .slice(0, 4)
    .map((x) => ({
      slug: x.slug,
      description: x.description,
      image: competitionImage(x),
    }));

  return (
    <>
      {/* Countdown sub-bar */}
      <div className="bg-black py-2.5 text-center font-jost text-[14px] font-medium text-white">
        Ends in{" "}
        <span className="text-botb-orange">00d : 14h : 22m : 04s</span>
      </div>

      <div className="mx-auto max-w-[1100px] px-4 py-8">
        {/* Hero image + badge */}
        <div className="relative">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src={c.image}
            alt={displayTitle}
            className="aspect-[16/9] w-full rounded-lg object-cover"
          />
          <span className="botb-badge-gradient absolute left-3 top-3 rounded px-2.5 py-1 font-jost text-[12px] font-medium uppercase leading-none text-white">
            {c.badge}
          </span>
        </div>

        {/* Gallery tabs (visual only) */}
        <div className="mt-4 flex gap-6 border-b border-botb-card-border">
          {GALLERY_TABS.map((tab, i) => (
            <span
              key={tab}
              className={`-mb-px border-b-2 pb-2 font-jost text-[14px] font-semibold uppercase ${
                i === 0
                  ? "border-botb-orange text-botb-orange"
                  : "border-transparent text-botb-muted"
              }`}
            >
              {tab}
            </span>
          ))}
        </div>

        {/* Category label + title */}
        <p
          className={`mt-6 inline-block rounded px-2.5 py-1 font-jost text-[12px] font-bold uppercase leading-none text-white ${TITLE_BAR_BG[c.titleColor]}`}
        >
          {c.title}
        </p>
        <h1 className="mt-3 font-jost text-[24px] font-bold text-botb-text md:text-[30px]">
          {displayTitle}
        </h1>

        {/* Prize pill */}
        <p className="mt-3 inline-flex rounded-full border border-botb-orange px-4 py-1 font-jost text-[13px] font-semibold uppercase text-botb-orange">
          {c.value}
        </p>

        {/* Two-column area: details on the left, price/entry on the right */}
        <div className="mt-8 grid grid-cols-1 gap-8 md:grid-cols-[1fr_320px]">
          {/* Left column */}
          <div className="order-2 md:order-1">
            {/* Spec grid */}
            <div className="grid grid-cols-2 gap-4">
              {specs.map((s) => (
                <div
                  key={s.label}
                  className="flex items-start gap-3 rounded-lg border border-botb-card-border bg-white p-4"
                >
                  <TicketIcon className="mt-0.5 h-5 w-5 flex-none text-botb-orange" />
                  <div className="min-w-0">
                    <p className="text-[11px] font-medium uppercase tracking-wide text-botb-muted">
                      {s.label}
                    </p>
                    <p className="mt-0.5 font-jost text-[15px] font-bold text-botb-text">
                      {s.value}
                    </p>
                  </div>
                </div>
              ))}
            </div>

            {/* Details prose */}
            <div className="mt-8">
              <h2 className="font-jost text-[20px] font-bold text-botb-text">
                Details
              </h2>
              <div className="mt-3 space-y-4 text-[15px] leading-relaxed text-botb-text">
                <p>
                  Win {displayTitle.toLowerCase()} for a fraction of its value.
                  This {c.title.toLowerCase()} prize is one of our most
                  sought-after competitions — enter now for your chance to drive
                  away a winner.
                </p>
                <p>
                  Every ticket costs just {price}, and the more you buy the
                  greater your chances. Prefer the cash? Take an equivalent cash
                  alternative instead — the choice is always yours.
                </p>
                <p>
                  Winners are drawn live on Facebook, so you can watch the moment
                  unfold in real time. Good luck!
                </p>
              </div>
            </div>
          </div>

          {/* Right column: price + entry */}
          <aside className="order-1 md:order-2">
            <div className="rounded-lg border border-botb-card-border bg-white p-5">
              <p className="text-[11px] font-medium uppercase tracking-wide text-botb-muted">
                {priceLabel}
              </p>
              <p className="mt-1 font-jost text-[28px] font-bold text-botb-text">
                {price}
              </p>

              {/* Progress */}
              {c.stat && (
                <div className="mt-4">
                  <div className="flex items-center justify-between text-[13px]">
                    <span className="font-semibold text-botb-orange">
                      {c.stat.label}
                    </span>
                    {c.stat.tickets && (
                      <span className="flex items-center gap-1 text-botb-muted">
                        <TicketIcon className="h-3.5 w-3.5 text-botb-muted" />
                        {c.stat.tickets}
                      </span>
                    )}
                  </div>
                  <div className="mt-1.5 h-2 w-full overflow-hidden rounded-full bg-[#e5e5e5]">
                    <div
                      className="h-full rounded-full bg-botb-orange"
                      style={{
                        width: `${Math.min(100, Math.max(0, c.stat.percent))}%`,
                      }}
                    />
                  </div>
                </div>
              )}

              <div className="mt-5">
                <EnterEntry
                  slug={c.slug}
                  title={displayTitle}
                  image={c.image}
                  price={price}
                />
              </div>
            </div>
          </aside>
        </div>

        {/* More competitions */}
        <section className="mt-12">
          <h2 className="font-jost text-[20px] font-bold text-botb-text">
            More competitions
          </h2>
          <div className="mt-4 grid grid-cols-2 gap-4 md:grid-cols-4">
            {more.map((m) => (
              <Link
                key={m.slug}
                href={`/prizes/${m.slug}`}
                className="group flex flex-col overflow-hidden rounded-lg border border-botb-card-border bg-white transition-all duration-200 hover:-translate-y-0.5 hover:shadow-[0_6px_18px_0_rgba(0,0,0,0.18)]"
              >
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={m.image}
                  alt={m.description}
                  className="aspect-[342/212] w-full object-cover"
                />
                <p className="px-3 py-3 font-jost text-[14px] font-semibold leading-snug text-botb-text group-hover:text-botb-orange">
                  {m.description}
                </p>
              </Link>
            ))}
          </div>
        </section>
      </div>
    </>
  );
}
