"use client";

import { useMemo, useState } from "react";
import { allCompetitions } from "@/lib/competitions";
import { CompetitionCard } from "@/components/CompetitionCard";
import { categoryNav } from "@/lib/data";
import { cn } from "@/lib/utils";

const ALL = "all";

/** Map a category id to a human-readable label using categoryNav, falling back
 * to a prettified version of the id (e.g. "ends-soon" -> "Ends Soon"). */
function categoryLabel(id: string): string {
  const match = categoryNav.find((item) => item.targetId === id);
  if (match) return match.label;
  return id
    .split("-")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

export default function CompetitionsPage() {
  const [activeCategory, setActiveCategory] = useState<string>(ALL);

  // Distinct category ids present in the data, ordered by categoryNav where
  // possible so the filter row matches the site's navigation ordering.
  const categories = useMemo(() => {
    const present = new Set(allCompetitions.map((c) => c.category));
    const ordered = categoryNav
      .map((item) => item.targetId)
      .filter((id) => present.has(id));
    // Append any category ids that aren't represented in categoryNav.
    const extras = Array.from(present).filter((id) => !ordered.includes(id));
    return [...ordered, ...extras];
  }, []);

  const filtered = useMemo(() => {
    if (activeCategory === ALL) return allCompetitions;
    return allCompetitions.filter((c) => c.category === activeCategory);
  }, [activeCategory]);

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

        {/* Filter pills */}
        <div className="no-scrollbar mt-6 flex gap-2 overflow-x-auto">
          <FilterPill
            label="All"
            active={activeCategory === ALL}
            onClick={() => setActiveCategory(ALL)}
          />
          {categories.map((id) => (
            <FilterPill
              key={id}
              label={categoryLabel(id)}
              active={activeCategory === id}
              onClick={() => setActiveCategory(id)}
            />
          ))}
        </div>

        {/* Result count */}
        <p className="mt-4 text-[14px] text-botb-muted">
          Showing {filtered.length}{" "}
          {filtered.length === 1 ? "competition" : "competitions"}
        </p>

        {/* Grid */}
        <div className="mt-6 grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-4">
          {filtered.map((competition) => (
            <CompetitionCard key={competition.slug} competition={competition} />
          ))}
        </div>
      </div>
    </section>
  );
}

function FilterPill({
  label,
  active,
  onClick,
}: {
  label: string;
  active: boolean;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      aria-pressed={active}
      className={cn(
        "shrink-0 whitespace-nowrap rounded-full px-4 py-2 font-jost text-[14px] font-medium uppercase transition-colors",
        active
          ? "bg-botb-orange text-white"
          : "border border-botb-card-border bg-white text-botb-text hover:bg-botb-gray",
      )}
    >
      {label}
    </button>
  );
}
