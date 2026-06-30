import Link from "next/link";

type ImpactCard = {
  title: string;
  stat: string;
  description: string;
};

const impactCards: ImpactCard[] = [
  {
    title: "Total Donated",
    stat: "£2M+",
    description:
      "Across more than two decades, BOTB and its players have helped raise over £2 million for causes that matter, both at home and abroad.",
  },
  {
    title: "Community Projects",
    stat: "150+ projects",
    description:
      "We back local community initiatives near our offices and airport sites — from food banks to community centres keeping neighbourhoods thriving.",
  },
  {
    title: "Disaster Relief",
    stat: "Rapid response",
    description:
      "When crises strike, we partner with established relief charities to provide fast, meaningful support to families who need it most.",
  },
  {
    title: "Youth Sports",
    stat: "40+ clubs",
    description:
      "We fund kit, equipment and facilities for grassroots youth sports clubs, helping the next generation stay active, healthy and inspired.",
  },
];

export default function CharityPage() {
  return (
    <div className="mx-auto max-w-4xl px-4 py-10">
      <header className="mb-8">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          Charity
        </h1>
        <p className="mt-2 text-botb-muted">
          Giving back is part of who we are.
        </p>
      </header>

      <section className="space-y-4 text-botb-text">
        <p>
          At BOTB, winning isn&apos;t just about the dream cars, the cash and the
          life-changing prizes. From the very beginning, giving back to the
          communities that have supported us has been woven into the fabric of
          our business.
        </p>
        <p>
          Every week, alongside crowning new winners, we set aside a portion of
          what we do to support charitable causes — from local grassroots
          projects to international relief efforts. When our players win, good
          causes win too.
        </p>
      </section>

      <section className="mt-10 grid grid-cols-1 gap-6 sm:grid-cols-2">
        {impactCards.map((card) => (
          <article
            key={card.title}
            className="rounded-lg border border-botb-card-border bg-white p-6 shadow-sm"
          >
            <div className="font-jost text-2xl font-bold text-botb-orange">
              {card.stat}
            </div>
            <h2 className="mt-1 font-jost font-semibold text-botb-text">
              {card.title}
            </h2>
            <p className="mt-2 text-sm text-botb-muted">{card.description}</p>
          </article>
        ))}
      </section>

      <section className="mt-12 rounded-lg bg-botb-gray p-8 text-center">
        <p className="text-botb-text">
          Want to know more about our charitable partnerships or suggest a cause
          close to your heart?
        </p>
        <Link
          href="/contact-us"
          className="mt-5 inline-block rounded-md bg-botb-orange px-8 py-3 font-jost font-bold uppercase text-white transition-opacity hover:opacity-90"
        >
          Get in Touch
        </Link>
      </section>
    </div>
  );
}
