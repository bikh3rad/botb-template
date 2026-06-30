import Link from "next/link";
import { winners } from "@/lib/data";
import type { Winner } from "@/types";

const extraWinners: Winner[] = [
  { name: "Aisha K.", prize: "Rolex Submariner", wonFor: "£0.05", revealed: "2 days ago", image: "/images/winners/kfi-1.webp" },
  { name: "Liam D.", prize: "£5,000 Cash", wonFor: "£0.03", revealed: "3 days ago", image: "/images/winners/kfi-3.webp" },
  { name: "Chloe F.", prize: "PlayStation 6 Bundle", wonFor: "£0.02", revealed: "3 days ago", image: "/images/winners/kfi-5.webp" },
  { name: "Raj P.", prize: "Audi RS3 Carbon Black", wonFor: "£0.06", revealed: "4 days ago", image: "/images/winners/kfi-7.webp" },
];

const allWinners: Winner[] = [...winners, ...extraWinners];

const stats = [
  { value: "26 Years", label: "UK's No.1" },
  { value: "£160M+", label: "in prizes won" },
  { value: "721k+", label: "guaranteed winners" },
];

export default function WinnersPage() {
  return (
    <div className="mx-auto max-w-[1360px] px-4 py-10">
      <header className="mb-8">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          Previous Winners
        </h1>
        <p className="mt-2 text-botb-muted">
          Over 721k guaranteed winners — and counting. Real winners, every day.
        </p>
      </header>

      <section className="mb-10 rounded-lg bg-botb-gray p-6">
        <div className="grid grid-cols-3 gap-4 text-center">
          {stats.map((stat) => (
            <div key={stat.label}>
              <div className="font-jost text-xl font-bold text-botb-text md:text-2xl">
                {stat.value}
              </div>
              <div className="text-sm text-botb-muted">{stat.label}</div>
            </div>
          ))}
        </div>
      </section>

      <section className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
        {allWinners.map((winner, index) => (
          <WinnerCardItem key={`${winner.name}-${index}`} winner={winner} />
        ))}
      </section>

      <section className="mt-12 rounded-lg bg-botb-gray p-8 text-center">
        <h2 className="font-jost text-2xl font-bold uppercase text-botb-text">
          Could you be next?
        </h2>
        <p className="mt-2 text-botb-muted">
          Join thousands of winners and grab your chance today.
        </p>
        <Link
          href="/competitions"
          className="mt-5 inline-block rounded-md bg-botb-orange px-8 py-3 font-jost font-bold uppercase text-white transition-opacity hover:opacity-90"
        >
          Enter a Competition
        </Link>
      </section>
    </div>
  );
}

function WinnerCardItem({ winner }: { winner: Winner }) {
  return (
    <article className="overflow-hidden rounded-lg border border-botb-card-border bg-white shadow-sm transition-transform duration-200 hover:-translate-y-1 hover:shadow-md">
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={winner.image}
        alt={`${winner.name} won ${winner.prize}`}
        className="aspect-square w-full rounded-t-lg object-cover"
      />
      <div className="p-3">
        <h3 className="font-jost font-semibold text-botb-text">{winner.name}</h3>
        <p className="text-sm text-botb-muted">{winner.prize}</p>
        <p className="mt-1 text-xs font-medium text-green-600">
          Won for {winner.wonFor}
        </p>
        <p className="mt-1 text-xs text-botb-muted">Revealed {winner.revealed}</p>
      </div>
    </article>
  );
}
