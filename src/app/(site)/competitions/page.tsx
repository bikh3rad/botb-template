import { CompetitionsBrowser } from "@/components/CompetitionsBrowser";
import { getCompetitions } from "@/lib/api";
import { buildCompetitionList } from "@/lib/presentation";

// Rendered at request time (backend not required at build).
export const dynamic = "force-dynamic";

export default async function CompetitionsPage() {
  const competitions = await getCompetitions("live");
  const items = buildCompetitionList(competitions);

  return (
    <section className="bg-botb-gray">
      <div className="mx-auto max-w-[1360px] px-4 py-8">
        {/* Heading */}
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          All Competitions
        </h1>
        <p className="mt-1 text-botb-muted">
          Choose your tickets and check out — real winners, every day.
        </p>

        <CompetitionsBrowser items={items} />
      </div>
    </section>
  );
}
