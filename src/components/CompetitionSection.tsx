import { CompetitionCard } from "@/components/CompetitionCard";
import type { CompetitionSection as Section } from "@/types";

export function CompetitionSection({ section }: { section: Section }) {
  return (
    <section
      id={section.id}
      className="mx-auto w-full max-w-[1360px] scroll-mt-44 px-2 md:px-5"
    >
      <h2 className="font-jost text-[26px] font-semibold uppercase text-botb-text md:text-[30px]">
        {section.heading}
      </h2>
      {section.subtitle && (
        <p className="mt-1 text-[15px] text-botb-muted">{section.subtitle}</p>
      )}
      <div className="mt-5 grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-4">
        {section.competitions.map((competition, i) => (
          <CompetitionCard key={`${section.id}-${i}`} competition={competition} />
        ))}
      </div>
    </section>
  );
}
