"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { UserIcon, TicketIcon } from "@/components/icons";

// Mock dashboard figures — this is a portfolio clone with no backend.
const STATS = [
  { label: "Active Entries", value: "12", Icon: TicketIcon },
  { label: "Tickets Owned", value: "348", Icon: TicketIcon },
  { label: "Competitions Won", value: "3", Icon: UserIcon },
  { label: "BOTB Credit", value: "£25.00", Icon: UserIcon },
];

const QUICK_LINKS = [
  {
    label: "My Competitions",
    description: "View the competitions you have entered.",
    href: "/competitions",
  },
  {
    label: "Previous Winners",
    description: "See who has driven away a dream car.",
    href: "/winners",
  },
  {
    label: "Perks for Playing",
    description: "Unlock rewards just for taking part.",
    href: "/account/perks-for-playing",
  },
  {
    label: "BOTB Pass",
    description: "Manage your subscription and benefits.",
    href: "/botb-pass",
  },
];

const RECENT_ACTIVITY = [
  {
    title: "Entered the Dream Car Competition",
    detail: "Porsche 911 GT3 — 5 tickets",
    time: "2 hours ago",
  },
  {
    title: "You won a Lifestyle prize!",
    detail: "£250 cash + AirPods Pro",
    time: "Yesterday",
  },
  {
    title: "Topped up your BOTB Credit",
    detail: "+ £25.00 added to your balance",
    time: "3 days ago",
  },
  {
    title: "Entered the Midweek Car Competition",
    detail: "Audi RS6 Avant — 2 tickets",
    time: "Last week",
  },
];

export default function AccountPage() {
  const router = useRouter();

  // Mock-only: there is no session to clear, so we just route home.
  function handleLogout() {
    router.push("/");
  }

  return (
    <section className="bg-white">
      <div className="mx-auto max-w-[1360px] px-4 py-10">
        <div className="flex flex-col gap-1">
          <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
            My Account
          </h1>
          <p className="text-botb-muted">Welcome back, Player!</p>
        </div>

        {/* Stat cards */}
        <div className="mt-8 grid grid-cols-2 gap-4 md:grid-cols-4">
          {STATS.map(({ label, value, Icon }) => (
            <div
              key={label}
              className="rounded-md border border-botb-card-border p-4"
            >
              <div className="flex items-center gap-2 text-botb-muted">
                <Icon className="h-5 w-5" />
                <span className="text-[13px]">{label}</span>
              </div>
              <p className="mt-2 font-jost text-[24px] font-bold text-botb-text">
                {value}
              </p>
            </div>
          ))}
        </div>

        {/* Quick links */}
        <div className="mt-10">
          <h2 className="font-jost text-[20px] font-semibold uppercase text-botb-text">
            Quick links
          </h2>
          <div className="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {QUICK_LINKS.map(({ label, description, href }) => (
              <Link
                key={label}
                href={href}
                className="rounded-md border border-botb-card-border p-4 transition-colors hover:border-botb-orange"
              >
                <p className="font-jost text-[16px] font-semibold text-botb-text">
                  {label}
                </p>
                <p className="mt-1 text-[14px] text-botb-muted">{description}</p>
              </Link>
            ))}
          </div>
        </div>

        {/* Recent activity */}
        <div className="mt-10">
          <h2 className="font-jost text-[20px] font-semibold uppercase text-botb-text">
            Recent activity
          </h2>
          <ul className="mt-4 divide-y divide-botb-card-border rounded-md border border-botb-card-border">
            {RECENT_ACTIVITY.map((activity) => (
              <li
                key={activity.title}
                className="flex items-center justify-between gap-4 p-4"
              >
                <div>
                  <p className="text-[15px] font-medium text-botb-text">
                    {activity.title}
                  </p>
                  <p className="text-[14px] text-botb-muted">
                    {activity.detail}
                  </p>
                </div>
                <span className="shrink-0 text-[13px] text-botb-muted">
                  {activity.time}
                </span>
              </li>
            ))}
          </ul>
        </div>

        {/* Account actions */}
        <div className="mt-10 flex flex-col items-start gap-4 sm:flex-row sm:items-center">
          <button
            type="button"
            onClick={handleLogout}
            className="rounded bg-botb-orange px-6 py-2.5 font-jost text-[15px] font-medium uppercase text-white hover:bg-botb-orange-hover"
          >
            Log out
          </button>
          <Link href="/" className="text-[14px] text-botb-muted hover:text-botb-orange">
            Back to home
          </Link>
        </div>
      </div>
    </section>
  );
}
