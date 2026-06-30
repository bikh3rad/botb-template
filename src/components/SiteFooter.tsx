"use client";

/* eslint-disable @next/next/no-img-element */

import { useState } from "react";
import { footerColumns, footerDescription } from "@/lib/data";
import type { FooterColumn } from "@/types";
import { ChevronDownIcon } from "@/components/icons";
import { cn } from "@/lib/utils";

/** Small icon + two-line stat used in the stats strip. */
function StatItem({ icon, value, label }: { icon: string; value: string; label: string }) {
  return (
    <div className="flex items-center gap-3">
      <img src={icon} alt="" className="h-9 w-9 shrink-0 object-contain" />
      <div className="leading-tight">
        <div className="text-[14px] font-semibold text-botb-secondary">{value}</div>
        <div className="text-[13px] text-botb-muted">{label}</div>
      </div>
    </div>
  );
}

/**
 * A single footer link column. When `collapsible` is true the title becomes a
 * toggle button on mobile (links hidden by default); on md+ the links are
 * always visible regardless of the collapsed state.
 */
function FooterLinkColumn({ column }: { column: FooterColumn }) {
  const [open, setOpen] = useState(false);
  const collapsible = column.collapsible ?? false;

  return (
    <div>
      {collapsible ? (
        <button
          type="button"
          onClick={() => setOpen((prev) => !prev)}
          aria-expanded={open}
          className="mb-3 flex w-full items-center justify-between md:cursor-default md:pointer-events-none"
        >
          <span className="font-jost text-[16px] font-semibold text-botb-text">{column.title}</span>
          <ChevronDownIcon
            className={cn(
              "h-4 w-4 text-botb-muted transition-transform md:hidden",
              open && "rotate-180",
            )}
          />
        </button>
      ) : (
        <h3 className="mb-3 font-jost text-[16px] font-semibold text-botb-text">{column.title}</h3>
      )}

      <ul
        className={cn(
          "space-y-2",
          // Collapsible columns are hidden on mobile unless toggled open; always shown on md+.
          collapsible && !open && "hidden md:block",
        )}
      >
        {column.links.map((link) => (
          <li key={link.label}>
            <a
              href={link.href}
              className="text-[14px] text-botb-muted transition-colors hover:text-botb-orange"
            >
              {link.label}
            </a>
          </li>
        ))}
      </ul>
    </div>
  );
}

export function SiteFooter() {
  return (
    <footer>
      {/* 1) STATS STRIP */}
      <div className="mx-auto w-full max-w-[1360px] px-4 py-4">
        <div className="flex flex-wrap items-center gap-6 rounded bg-botb-gray px-4 py-4">
          <StatItem icon="/images/footer/est99.png" value="26 years" label="UK's No.1" />
          <StatItem icon="/images/footer/gift.png" value="Over £160 million" label="in prizes won" />
          <StatItem icon="/images/footer/cup.png" value="Over 721k" label="guaranteed winners" />

          <div className="ml-auto hidden items-center gap-4 md:flex">
            <img src="/images/footer/feefo.png" alt="Feefo rating" className="h-10 w-auto object-contain" />
            <span className="h-8 border-l border-botb-card-border" aria-hidden="true" />
            <img
              src="/images/footer/trustpilot.png"
              alt="Trustpilot rating"
              className="h-10 w-auto object-contain"
            />
          </div>
        </div>
      </div>

      {/* 2) MAIN FOOTER */}
      <div className="mx-auto w-full max-w-[1360px] bg-white px-4 py-10">
        {/* Top area: description + app/social */}
        <div className="grid gap-10 md:grid-cols-[2fr_1fr]">
          {/* LEFT: heading + description */}
          <div>
            <h2 className="font-jost text-[18px] font-semibold text-botb-text">
              WEEKLY COMPETITION WINNERS - GUARANTEED
            </h2>
            <div className="mt-4 space-y-3">
              {footerDescription.map((paragraph, index) => (
                <p key={index} className="text-[14px] leading-relaxed text-botb-muted">
                  {paragraph}
                </p>
              ))}
              <p className="text-[14px] leading-relaxed text-botb-muted">
                Winvia Entertainment PLC is part of{" "}
                <a href="#" className="font-bold text-botb-text underline">
                  Winvia Entertainment Group
                </a>
                .
              </p>
            </div>
          </div>

          {/* RIGHT: app downloads + social */}
          <div>
            <h3 className="font-jost text-[16px] font-semibold text-botb-text">Download App&apos;s</h3>
            <div className="mt-3 flex flex-wrap items-center gap-3">
              <a href="#" aria-label="Download on the App Store">
                <img src="/images/footer/app-store.png" alt="Download on the App Store" className="h-10 w-auto" />
              </a>
              <a href="#" aria-label="Get it on Google Play">
                <img src="/images/footer/google-play.png" alt="Get it on Google Play" className="h-10 w-auto" />
              </a>
            </div>

            <div className="mt-4 flex items-center gap-2 text-[15px]">
              <span className="text-botb-orange" aria-hidden="true">
                ★★★★★
              </span>
              <span className="font-semibold text-botb-text">4.8</span>
            </div>

            <p className="mt-4 text-[14px] font-semibold text-botb-text">Follow on:</p>
            <div className="mt-2 flex items-center gap-3">
              <a href="#" aria-label="Facebook">
                <img src="/images/footer/fb.png" alt="Facebook" className="h-9 w-9" />
              </a>
              <a href="#" aria-label="Instagram">
                <img src="/images/footer/ig.png" alt="Instagram" className="h-9 w-9" />
              </a>
              <a href="#" aria-label="YouTube">
                <img src="/images/footer/yt.png" alt="YouTube" className="h-9 w-9" />
              </a>
              <a href="#" aria-label="TikTok">
                <img src="/images/footer/tt.png" alt="TikTok" className="h-9 w-9" />
              </a>
            </div>
          </div>
        </div>

        {/* LINK COLUMNS */}
        <div className="mt-10 grid grid-cols-2 gap-6 border-t border-botb-card-border pt-8 sm:grid-cols-3 lg:grid-cols-6">
          {footerColumns.map((column) => (
            <FooterLinkColumn key={column.title} column={column} />
          ))}
        </div>
      </div>
    </footer>
  );
}
