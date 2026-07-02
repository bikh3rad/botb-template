import { HeroCarousel } from "@/components/HeroCarousel";
import { CategoryNav } from "@/components/CategoryNav";
import { WinnersTicker } from "@/components/WinnersTicker";
import { FeaturedSection } from "@/components/FeaturedSection";
import { CompetitionSection } from "@/components/CompetitionSection";
import { TrustBand } from "@/components/TrustBand";
import { JustLaunched } from "@/components/JustLaunched";
import { getCompetitions, getWinners } from "@/lib/api";
import { buildHomeView } from "@/lib/presentation";

// Rendered at request time so `next build` never needs a live backend and the
// grids always reflect current competition/winner data.
export const dynamic = "force-dynamic";

export default async function Home() {
  // Fetch live competitions + the winners feed in parallel, then group them into
  // the homepage sections using the slug-keyed presentation map.
  const [competitions, winners] = await Promise.all([
    getCompetitions("live"),
    getWinners(),
  ]);
  const view = buildHomeView(competitions, winners);

  return (
    <>
      <HeroCarousel />
      <CategoryNav />

      {/* Competition area (gray) */}
      <div className="bg-botb-gray pb-12 pt-2">
        <WinnersTicker winners={view.winners} winnersCount={view.winnersCount} />
        <div className="flex flex-col gap-10">
          <FeaturedSection featured={view.featured} />
          {view.sections.map((section) => (
            <CompetitionSection key={section.id} section={section} />
          ))}
          {/* Pass Exclusives anchor (category nav target) */}
          <div id="pass-exclusives" className="scroll-mt-44" />
        </div>
      </div>

      <TrustBand />
      <JustLaunched />
    </>
  );
}
