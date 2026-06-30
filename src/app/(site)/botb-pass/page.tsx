import Link from "next/link";

type Tier = {
  name: string;
  price: string;
  period: string;
  features: string[];
  popular: boolean;
};

const tiers: Tier[] = [
  {
    name: "Lite",
    price: "£9.99",
    period: "/mo",
    features: [
      "Extra free tickets every month",
      "Access to exclusive competitions",
      "£5 bonus credit each month",
      "Member newsletter",
    ],
    popular: false,
  },
  {
    name: "Plus",
    price: "£19.99",
    period: "/mo",
    features: [
      "More extra free tickets",
      "Exclusive member-only comps",
      "Priority winner selection",
      "£15 bonus credit each month",
      "Priority support",
    ],
    popular: true,
  },
  {
    name: "Ultimate",
    price: "£39.99",
    period: "/mo",
    features: [
      "Maximum extra free tickets",
      "All exclusive competitions",
      "Top priority winner selection",
      "£40 bonus credit each month",
      "Dedicated VIP support",
    ],
    popular: false,
  },
];

const perks = [
  "Exclusive competitions",
  "Extra free tickets",
  "Member-only prizes",
  "Priority support",
];

export default function BotbPassPage() {
  return (
    <div className="mx-auto max-w-[1360px] px-4 py-10">
      {/* Page header */}
      <header className="mb-10 text-center">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          BOTB Pass
        </h1>
        <p className="mt-3 text-[16px] text-botb-muted">
          Better odds, more entries, exclusive perks.
        </p>
      </header>

      {/* Pricing tiers */}
      <section className="mb-16">
        <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
          {tiers.map((tier) => (
            <div
              key={tier.name}
              className={`relative flex flex-col rounded-lg border bg-white p-6 shadow-sm ${
                tier.popular
                  ? "border-botb-orange shadow-md"
                  : "border-botb-card-border"
              }`}
            >
              {tier.popular && (
                <span className="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-botb-orange px-4 py-1 font-jost text-[12px] font-bold uppercase tracking-wide text-white">
                  Most Popular
                </span>
              )}

              <h2 className="font-jost text-[22px] font-bold uppercase text-botb-text">
                {tier.name}
              </h2>

              <div className="mt-3 flex items-baseline gap-1">
                <span className="font-jost text-[40px] font-bold leading-none text-botb-text">
                  {tier.price}
                </span>
                <span className="text-[15px] text-botb-muted">
                  {tier.period}
                </span>
              </div>

              <ul className="mt-6 flex-1 space-y-3">
                {tier.features.map((feature) => (
                  <li
                    key={feature}
                    className="flex items-start gap-2 text-[15px] leading-6 text-botb-text"
                  >
                    <span
                      aria-hidden
                      className="mt-0.5 font-bold text-botb-orange"
                    >
                      ✓
                    </span>
                    <span>{feature}</span>
                  </li>
                ))}
              </ul>

              <Link
                href="/competitions"
                className={
                  tier.popular
                    ? "mt-8 inline-block rounded bg-botb-orange px-6 py-2.5 text-center font-jost text-[15px] font-medium uppercase text-white hover:bg-botb-orange-hover"
                    : "mt-8 inline-block rounded border border-botb-orange px-6 py-2.5 text-center font-jost text-[15px] font-medium uppercase text-botb-orange hover:bg-botb-orange hover:text-white"
                }
              >
                Choose Plan
              </Link>
            </div>
          ))}
        </div>
      </section>

      {/* Perks strip */}
      <section className="rounded-lg bg-botb-gray p-8">
        <div className="grid grid-cols-2 gap-6 text-center md:grid-cols-4">
          {perks.map((perk) => (
            <div key={perk} className="flex flex-col items-center gap-3">
              <span className="flex h-12 w-12 items-center justify-center rounded-full bg-botb-orange font-bold text-white">
                ✓
              </span>
              <span className="font-jost text-[15px] font-medium uppercase text-botb-text">
                {perk}
              </span>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
