"use client";

import { useMemo, useState } from "react";
import { CompetitionCard } from "@/components/CompetitionCard";
import { categoryNav } from "@/lib/data";
import type { CardView } from "@/lib/presentation";
import { cn } from "@/lib/utils";

const ALL = "all";

/** One competition plus the category it filters under. */
export interface BrowserItem {
  card: CardView;
  category: string;
}

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

/**
 * Client-side category filter over the API-driven competition list. The data is
 * fetched by the Server Component parent and passed in as props — this component
 * only owns the interactive filtering.
 */
export function CompetitionsBrowser({ items }: { items: BrowserItem[] }) {
  const [activeCategory, setActiveCategory] = useState<string>(ALL);

  // Distinct category ids present in the data, ordered by categoryNav where
  // possible so the filter row matches the site's navigation ordering.
  const categories = useMemo(() => {
    const present = new Set(items.map((i) => i.category));
    const ordered = categoryNav
      .map((item) => item.targetId)
      .filter((id) => present.has(id));
    const extras = Array.from(present).filter((id) => !ordered.includes(id));
    return [...ordered, ...extras];
  }, [items]);

  const filtered = useMemo(() => {
    if (activeCategory === ALL) return items;
    return items.filter((i) => i.category === activeCategory);
  }, [activeCategory, items]);

  return (
    <>
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
        {filtered.map((item) => (
          <CompetitionCard key={item.card.slug} competition={item.card} />
        ))}
      </div>
    </>
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
