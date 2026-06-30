import Link from "next/link";
import { CompetitionCard } from "@/components/CompetitionCard";
import { dreamCarCard, featuredCompetitions, subscriberCard } from "@/lib/data";

function DreamCarCard() {
  const d = dreamCarCard;
  return (
    <Link href="/prizes/cars" className="group flex h-full flex-col overflow-hidden rounded-lg border border-botb-card-border bg-white shadow-[0_1.5px_8px_0_rgba(0,0,0,0.2)] transition-all duration-200 hover:-translate-y-0.5">
      {/* How-it-works strip over dark image */}
      <div className="relative">
        <img src={d.heroImage} alt="Dream Car" className="aspect-[342/150] w-full object-cover" />
        <div className="absolute inset-0 bg-black/55" />
        <div className="absolute inset-0 flex items-center justify-around px-3 text-white">
          {d.steps.map((step, i) => (
            <div key={step} className="flex items-center gap-1.5">
              <span className="flex h-5 w-5 items-center justify-center rounded-full bg-botb-orange text-[11px] font-bold">
                {i + 1}
              </span>
              <span className="font-jost text-[12px] font-medium leading-tight">{step}</span>
              {i < d.steps.length - 1 && <span className="text-botb-orange">→</span>}
            </div>
          ))}
        </div>
        <span className="botb-badge-gradient absolute left-2 top-2 z-[1] rounded px-2 py-1 font-jost text-[12px] font-medium uppercase leading-none text-white">
          {d.badge}
        </span>
      </div>

      {/* Orange title bar */}
      <div className="bg-botb-orange px-5 py-3 text-center font-jost text-[14px] font-bold uppercase text-white">
        {d.titleBar}
      </div>

      {/* Body */}
      <div className="flex flex-1 items-center gap-3 px-4 py-4">
        <img src={d.stbBadge} alt="Spot the Ball" className="w-24 shrink-0" />
        <div className="flex-1 text-center">
          <p className="text-[15px] leading-snug text-botb-text">{d.description}</p>
          <p className="mt-3 text-[11px] font-medium uppercase text-botb-muted">{d.priceLabel}</p>
          <p className="font-jost text-[22px] font-bold text-botb-text">{d.price}</p>
        </div>
      </div>
      <div className="px-4 pb-4">
        <span className="block w-full rounded bg-botb-orange px-4 py-2.5 text-center font-jost text-[15px] font-medium uppercase text-white transition-colors group-hover:bg-botb-orange-hover">
          Play Now »
        </span>
      </div>
    </Link>
  );
}

function SubscriberCard() {
  const s = subscriberCard;
  return (
    <Link
      href="/prizes/iphone-17-and-1-249-prizes"
      className="group relative flex h-full flex-col overflow-hidden rounded-lg border border-botb-card-border bg-[#1f5b3a] shadow-[0_1.5px_8px_0_rgba(0,0,0,0.2)] transition-all duration-200 hover:-translate-y-0.5"
    >
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img src={s.image} alt={s.headline} className="h-full w-full object-cover" />
      <span className="botb-badge-gradient absolute left-2 top-2 z-[1] rounded px-2 py-1 font-jost text-[12px] font-medium uppercase leading-none text-white">
        {s.badge}
      </span>
    </Link>
  );
}

export function FeaturedSection() {
  return (
    <section
      id="featured-competitions"
      className="mx-auto w-full max-w-[1360px] scroll-mt-44 px-2 md:px-5"
    >
      <h2 className="font-jost text-[26px] font-semibold uppercase text-botb-text md:text-[30px]">
        Featured Competitions
      </h2>
      <div className="mt-5 grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-4">
        <DreamCarCard />
        <CompetitionCard competition={featuredCompetitions[0]} />
        <CompetitionCard competition={featuredCompetitions[1]} />
        <CompetitionCard competition={featuredCompetitions[2]} />
        <SubscriberCard />
        <CompetitionCard competition={featuredCompetitions[3]} />
        <CompetitionCard competition={featuredCompetitions[4]} />
      </div>
    </section>
  );
}
