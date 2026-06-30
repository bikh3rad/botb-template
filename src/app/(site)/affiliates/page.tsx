import Link from "next/link";

type Step = {
  number: string;
  title: string;
  description: string;
};

const steps: Step[] = [
  {
    number: "1",
    title: "Sign up",
    description:
      "Join the BOTB affiliate programme in minutes. It's completely free and there are no targets to hit.",
  },
  {
    number: "2",
    title: "Share your link",
    description:
      "Get your unique referral link and share it with your audience, friends and followers wherever they are.",
  },
  {
    number: "3",
    title: "Earn commission",
    description:
      "Every time someone plays through your link, you earn. Track your performance from your dashboard.",
  },
];

const benefits: string[] = [
  "Competitive commission on every qualifying player you refer",
  "Real-time tracking and transparent reporting",
  "A library of banners, creatives and copy ready to use",
  "Dedicated affiliate support from our UK-based team",
  "Regular bonuses and incentives for top performers",
  "Reliable monthly payouts, no minimum threshold games",
];

export default function AffiliatesPage() {
  return (
    <section className="bg-white">
      <div className="mx-auto max-w-[1360px] px-4 py-10">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          Affiliate Programme
        </h1>
        <p className="mt-2 max-w-3xl text-botb-muted">
          Partner with the home of the Dream Car Competition. Refer players to
          BOTB and earn generous commission while sharing competitions people
          genuinely love. No fees, no targets — just rewards for every player
          you send our way.
        </p>

        {/* How it works */}
        <h2 className="mt-12 font-jost text-[22px] font-bold uppercase text-botb-text">
          How it works
        </h2>
        <div className="mt-6 grid grid-cols-1 gap-6 md:grid-cols-3">
          {steps.map((step) => (
            <div
              key={step.number}
              className="rounded-md border border-botb-card-border p-6"
            >
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-botb-orange font-jost text-[20px] font-bold text-white">
                {step.number}
              </div>
              <h3 className="mt-4 font-jost text-[18px] font-semibold text-botb-text">
                {step.title}
              </h3>
              <p className="mt-2 text-[14px] text-botb-muted">
                {step.description}
              </p>
            </div>
          ))}
        </div>

        {/* Benefits */}
        <h2 className="mt-12 font-jost text-[22px] font-bold uppercase text-botb-text">
          Why join
        </h2>
        <ul className="mt-6 grid grid-cols-1 gap-3 md:grid-cols-2">
          {benefits.map((benefit) => (
            <li
              key={benefit}
              className="flex items-start gap-3 rounded-md border border-botb-card-border p-4"
            >
              <span
                aria-hidden
                className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-botb-orange text-[12px] font-bold text-white"
              >
                ✓
              </span>
              <span className="text-[15px] text-botb-text">{benefit}</span>
            </li>
          ))}
        </ul>

        {/* CTA */}
        <div className="mt-12 rounded-md bg-botb-gray p-8 text-center">
          <h2 className="font-jost text-[22px] font-bold uppercase text-botb-text">
            Ready to start earning?
          </h2>
          <p className="mx-auto mt-2 max-w-xl text-botb-muted">
            Create your account today and get your unique affiliate link in
            minutes.
          </p>
          <Link
            href="/register"
            className="mt-6 inline-block rounded bg-botb-orange px-6 py-2.5 font-jost text-[15px] font-medium uppercase text-white hover:bg-botb-orange-hover"
          >
            Join now
          </Link>
        </div>
      </div>
    </section>
  );
}
