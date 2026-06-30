"use client";

import Link from "next/link";
import { TicketIcon } from "@/components/icons";

// Static order reference — deterministic so server and client markup match
// (no Math.random during render, which would trigger a hydration mismatch).
const orderRef = "BOTB-2026-004521";

export default function ConfirmationPage() {

  return (
    <div className="mx-auto max-w-[1100px] px-4 py-16">
      <div className="mx-auto max-w-md rounded-lg border border-botb-card-border bg-white p-8 text-center">
        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            className="h-8 w-8 text-green-600"
          >
            <path
              d="M20 6L9 17l-5-5"
              stroke="currentColor"
              strokeWidth="2.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        </div>

        <h1 className="mt-6 font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          Order Confirmed!
        </h1>
        <p className="mt-2 text-botb-muted">
          Thanks for your entry — good luck!
        </p>

        <div className="mt-6 flex items-center justify-center gap-2 rounded-md bg-botb-gray px-4 py-3">
          <TicketIcon className="h-5 w-5 text-botb-orange" />
          <span className="font-jost text-sm font-semibold text-botb-text">
            Order reference: {orderRef}
          </span>
        </div>

        <p className="mt-4 text-sm text-botb-muted">
          Winners are revealed daily — keep an eye on your account to see if
          you&apos;ve won.
        </p>

        <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:justify-center">
          <Link
            href="/account"
            className="rounded-md bg-botb-orange px-6 py-3 font-jost text-sm font-bold uppercase text-white transition-colors hover:bg-botb-orange-hover"
          >
            View My Account
          </Link>
          <Link
            href="/competitions"
            className="rounded-md border border-botb-card-border px-6 py-3 font-jost text-sm font-bold uppercase text-botb-text transition-colors hover:bg-botb-gray"
          >
            Back to Competitions
          </Link>
        </div>
      </div>
    </div>
  );
}
