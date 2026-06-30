import type { CSSProperties } from "react";

import { winners, winnersCount } from "@/lib/data";
import type { Winner } from "@/types";
import { cn } from "@/lib/utils";

/**
 * Single winner card rendered inline within the marquee track. Kept as a local
 * component so the winners list can be rendered twice (for the seamless loop)
 * without duplicating markup.
 */
function WinnerCard({ winner }: { winner: Winner }) {
  return (
    <div className="relative mx-2 flex min-w-[230px] shrink-0 items-center gap-3 rounded-lg border border-[#e3e3e3] bg-white px-3 py-2 shadow-sm">
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={winner.image}
        alt={`${winner.name} winner`}
        className="h-12 w-12 rounded object-cover"
        loading="lazy"
      />
      <div className="min-w-0">
        <p className="font-jost text-[14px] font-semibold text-botb-text">
          {winner.name}
        </p>
        <p className="text-[13px] text-botb-muted">{winner.prize}</p>
        <p className="text-[12px] font-medium text-green-600">
          Won for {winner.wonFor}
        </p>
      </div>
      <span className="botb-badge-gradient absolute right-1 top-1 rounded-full px-2 py-0.5 text-[10px] text-white">
        Revealed {winner.revealed}
      </span>
    </div>
  );
}

/**
 * Auto-scrolling horizontal ticker of recent winners. Pure CSS marquee — the
 * winners list is rendered twice back-to-back so the `.animate-marquee`
 * translateX(0 → -50%) loop is seamless. Hovering the viewport pauses it.
 */
export function WinnersTicker() {
  // Custom property consumed by the `.animate-marquee` keyframe duration var.
  const trackStyle = { "--marquee-duration": "40s" } as CSSProperties;

  return (
    <section className="mx-auto my-5 max-w-[1360px] px-2 md:px-5 lg:my-6">
      <div className="flex items-center">
        {/* Left label — the brand's handwritten "Another winner. Now." mark */}
        <div className="shrink-0 pr-4">
          <p className="font-jost text-[18px] font-bold italic leading-tight text-botb-orange">
            Another winner.
          </p>
          <p className="font-jost text-[18px] font-bold italic leading-tight text-botb-text">
            Now.
          </p>
        </div>

        {/* Marquee viewport */}
        <div className="marquee-pause min-w-0 flex-1 overflow-hidden">
          <div className={cn("flex w-max animate-marquee")} style={trackStyle}>
            {[...winners, ...winners].map((winner, index) => (
              <WinnerCard key={`${winner.name}-${index}`} winner={winner} />
            ))}
          </div>
        </div>

        {/* Far-right live counter */}
        <div className="ml-3 flex shrink-0 items-center gap-2">
          <span className="h-2 w-2 animate-pulse rounded-full bg-red-500" />
          <div className="leading-tight">
            <p className="text-[14px] font-bold text-botb-orange">
              {winnersCount}
            </p>
            <p className="text-[12px] text-botb-muted">in the last 24 hours</p>
          </div>
        </div>
      </div>
    </section>
  );
}
