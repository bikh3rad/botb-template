"use client";

import { useEffect, useState } from "react";
import { categoryNav } from "@/lib/data";
import { cn } from "@/lib/utils";

/**
 * Sticky horizontal category bar shown directly under the site header.
 * Smooth-scrolls to homepage sections and highlights the active section
 * via an IntersectionObserver scroll-spy.
 */
export function CategoryNav() {
  // Default active = first item so the bar never renders with nothing highlighted.
  const [activeId, setActiveId] = useState<string>(categoryNav[0]?.targetId ?? "");

  useEffect(() => {
    // Collect the section elements that actually exist in the DOM.
    const sections = categoryNav
      .map((item) => document.getElementById(item.targetId))
      .filter((el): el is HTMLElement => el !== null);

    if (sections.length === 0) return;

    const observer = new IntersectionObserver(
      (entries) => {
        // Pick the most recently intersecting section near the upper-middle
        // of the viewport so the active item tracks the user's scroll position.
        for (const entry of entries) {
          if (entry.isIntersecting) {
            setActiveId(entry.target.id);
          }
        }
      },
      // Shrinks the observer root to a thin band around the upper-middle of
      // the viewport, so a section becomes "active" as it crosses that line.
      { rootMargin: "-45% 0px -50% 0px", threshold: 0 }
    );

    sections.forEach((section) => observer.observe(section));

    return () => observer.disconnect();
  }, []);

  const handleClick = (targetId: string) => {
    const target = document.getElementById(targetId);
    if (!target) return;
    // Sticky offset is handled by each section's own scroll-mt utility.
    target.scrollIntoView({ behavior: "smooth", block: "start" });
  };

  return (
    <nav className="sticky top-[100px] z-40 w-full bg-white shadow-[0_1.5px_8px_0_rgba(0,0,0,0.2)]">
      <div className="mx-auto max-w-[1360px] px-4">
        <ul className="no-scrollbar flex flex-nowrap justify-start gap-6 overflow-x-auto lg:justify-center lg:gap-8">
          {categoryNav.map((item) => {
            const isActive = item.targetId === activeId;

            return (
              <li key={item.targetId} className="flex-none">
                <button
                  type="button"
                  onClick={() => handleClick(item.targetId)}
                  className="flex items-center py-3.5"
                  aria-current={isActive ? "true" : undefined}
                >
                  <span
                    className={cn(
                      "-mb-[2px] flex items-center gap-2 border-b-2 border-transparent pb-3.5",
                      isActive && "border-botb-orange"
                    )}
                  >
                    {/* Decorative pill icon — label conveys meaning. */}
                    <img
                      src={item.icon}
                      alt=""
                      aria-hidden="true"
                      className="h-6 w-6 object-contain"
                    />
                    <span
                      className={cn(
                        "whitespace-nowrap font-roboto text-[15px] font-medium",
                        isActive ? "text-botb-orange" : "text-botb-text"
                      )}
                    >
                      {item.label}
                    </span>
                  </span>
                </button>
              </li>
            );
          })}
        </ul>
      </div>
    </nav>
  );
}
