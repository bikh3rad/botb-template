import { HeroCarousel } from "@/components/HeroCarousel";
import { CategoryNav } from "@/components/CategoryNav";
import { WinnersTicker } from "@/components/WinnersTicker";
import { FeaturedSection } from "@/components/FeaturedSection";
import { CompetitionSection } from "@/components/CompetitionSection";
import { TrustBand } from "@/components/TrustBand";
import { JustLaunched } from "@/components/JustLaunched";
import { competitionSections } from "@/lib/data";

export default function Home() {
  return (
    <>
      <HeroCarousel />
      <CategoryNav />

      {/* Competition area (gray) */}
      <div className="bg-botb-gray pb-12 pt-2">
        <WinnersTicker />
        <div className="flex flex-col gap-10">
          <FeaturedSection />
          {competitionSections.map((section) => (
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
