import Link from "next/link";
import { cn } from "@/lib/utils";
import { CartIcon, TicketIcon } from "@/components/icons";
import { competitionHref } from "@/lib/competitions";
import type { Competition, TitleBarColor } from "@/types";

const TITLE_BAR_BG: Record<TitleBarColor, string> = {
  orange: "bg-botb-orange",
  teal: "bg-botb-teal",
  purple: "bg-botb-purple",
  red: "bg-botb-red",
  green: "bg-botb-green",
  blue: "bg-botb-blue",
  pass: "bg-botb-purple",
};

export function CompetitionCard({ competition }: { competition: Competition }) {
  const c = competition;
  return (
    <Link
      href={competitionHref(c)}
      className={cn(
        "group flex h-full flex-col overflow-hidden rounded-lg border border-botb-card-border bg-white text-center shadow-[0_1.5px_8px_0_rgba(0,0,0,0.2)] transition-all duration-200 hover:-translate-y-0.5 hover:shadow-[0_6px_18px_0_rgba(0,0,0,0.18)]",
        c.wide && "sm:col-span-2",
      )}
    >
      {/* Image + badge */}
      <div className="relative">
        <img
          src={c.image}
          alt={c.title}
          className={cn(
            "w-full object-cover",
            c.wide ? "aspect-[684/424] sm:aspect-[684/212]" : "aspect-[342/212]",
          )}
        />
        <span className="botb-badge-gradient absolute left-2 top-2 z-[1] rounded font-jost text-[11px] font-medium uppercase leading-none text-white px-2 py-1 lg:text-[12px]">
          {c.badge}
        </span>
      </div>

      {/* Colored title bar */}
      <div
        className={cn(
          "px-5 py-3 font-jost text-[14px] font-bold uppercase leading-tight text-white",
          TITLE_BAR_BG[c.titleColor],
        )}
      >
        {c.title}
      </div>

      {/* Body */}
      <div className="flex flex-1 flex-col px-4 py-4">
        <p className="mx-auto max-w-[90%] text-[15px] leading-snug text-botb-text">
          {c.description}
        </p>

        {/* Price */}
        {c.price && (
          <div className="mt-4">
            {c.priceLabel && (
              <p className="text-[11px] font-medium uppercase tracking-wide text-botb-muted">
                {c.priceLabel}
              </p>
            )}
            <p className="mt-1 font-jost text-[22px] font-bold text-botb-text">{c.price}</p>
          </div>
        )}

        {/* Progress */}
        {c.stat && (
          <div className="mt-3">
            <div className="flex items-center justify-between text-[13px]">
              <span className="font-semibold text-botb-orange">{c.stat.label}</span>
              {c.stat.tickets && (
                <span className="flex items-center gap-1 text-botb-muted">
                  <TicketIcon className="h-3.5 w-3.5 text-botb-muted" />
                  {c.stat.tickets}
                </span>
              )}
            </div>
            <div className="mt-1.5 h-1.5 w-full overflow-hidden rounded-full bg-[#e5e5e5]">
              <div
                className="h-full rounded-full bg-botb-orange"
                style={{ width: `${Math.min(100, Math.max(0, c.stat.percent))}%` }}
              />
            </div>
          </div>
        )}

        {/* Spacer pushes CTA to bottom */}
        <div className="flex-1" />

        {/* CTA (whole card is a link to the detail page) */}
        <div className="mt-4">
          {c.cta === "PLAY NOW" && (
            <span className="block w-full rounded bg-botb-orange px-4 py-2.5 font-jost text-[15px] font-medium uppercase text-white transition-colors group-hover:bg-botb-orange-hover">
              Play Now »
            </span>
          )}

          {c.cta === "DETAILS" && (
            <div className="flex flex-col gap-3">
              <span className="font-jost text-[13px] font-semibold uppercase text-botb-text group-hover:text-botb-orange">
                Details »
              </span>
              <span className="block w-full rounded bg-botb-orange px-4 py-2.5 font-jost text-[15px] font-medium uppercase text-white transition-colors group-hover:bg-botb-orange-hover">
                Enter Now
              </span>
            </div>
          )}

          {c.cta === "ENTER NOW" && (
            <div className="flex items-stretch gap-2">
              <span className="flex-1 rounded border border-botb-orange bg-white px-4 py-2 font-roboto text-[15px] uppercase text-botb-orange transition-colors group-hover:bg-botb-orange group-hover:text-white">
                Enter Now
              </span>
              {c.showCart && (
                <span
                  aria-hidden
                  className="flex w-11 shrink-0 items-center justify-center rounded bg-botb-orange text-white transition-colors group-hover:bg-botb-orange-hover"
                >
                  <CartIcon className="h-5 w-5" />
                </span>
              )}
            </div>
          )}
        </div>
      </div>
    </Link>
  );
}
