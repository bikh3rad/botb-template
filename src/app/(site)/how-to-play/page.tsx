import Image from "next/image";
import Link from "next/link";

const steps = [
  {
    number: 1,
    title: "Select your prize",
    description: "Choose from cars, cash, tech, houses & more.",
  },
  {
    number: 2,
    title: "Play the game / pick tickets",
    description:
      "For Dream Car play Spot the Ball; for others choose your tickets.",
  },
  {
    number: 3,
    title: "Win!",
    description:
      "Winners revealed every day. Take the prize or the cash alternative.",
  },
];

const faqs = [
  {
    question: "How often are winners drawn?",
    answer:
      "Winners are revealed every single day. Once a competition closes, the result is announced and the lucky winner is contacted directly.",
  },
  {
    question: "What is Spot the Ball?",
    answer:
      "Spot the Ball is our skill-based game for Dream Car competitions. You place a crosshair where you think the centre of the ball is. The closest guess to the judge's chosen spot wins.",
  },
  {
    question: "Can I take cash instead of the prize?",
    answer:
      "Yes. Every prize comes with a cash alternative, so if you'd rather have the money than the car, tech or house, the choice is yours.",
  },
  {
    question: "How do I know if I've won?",
    answer:
      "We contact every winner directly by phone and email, and winners are featured on our site and social channels.",
  },
  {
    question: "How many tickets can I buy?",
    answer:
      "You can buy as many tickets as you like, up to the limit set for each individual competition. More tickets means more chances to win.",
  },
];

export default function HowToPlayPage() {
  return (
    <div className="mx-auto max-w-[1360px] px-4 py-10">
      {/* Page header */}
      <header className="mb-10 max-w-3xl">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          How to Play
        </h1>
        <p className="mt-4 text-[16px] leading-7 text-botb-muted">
          Winning with BOTB couldn&apos;t be simpler. Pick the prize you&apos;d
          love to win, play the game or choose your tickets, and you could be
          our next lucky winner. Here&apos;s everything you need to know to get
          started.
        </p>
      </header>

      {/* 3-step section */}
      <section className="mb-16">
        <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
          {steps.map((step) => (
            <div
              key={step.number}
              className="flex flex-col items-start rounded-lg border border-botb-card-border bg-white p-6 shadow-sm"
            >
              <span className="flex h-14 w-14 items-center justify-center rounded-full bg-botb-orange font-jost text-[26px] font-bold text-white">
                {step.number}
              </span>
              <h2 className="mt-5 font-jost text-[20px] font-bold uppercase text-botb-text">
                {step.title}
              </h2>
              <p className="mt-2 text-[15px] leading-6 text-botb-muted">
                {step.description}
              </p>
            </div>
          ))}
        </div>
      </section>

      {/* Spot the Ball explained */}
      <section className="mb-16">
        <div className="grid grid-cols-1 items-center gap-8 rounded-lg border border-botb-card-border bg-botb-gray p-6 md:grid-cols-2 md:p-8">
          <div>
            <h2 className="font-jost text-[24px] font-bold uppercase text-botb-text md:text-[28px]">
              Spot the Ball explained
            </h2>
            <p className="mt-4 text-[15px] leading-7 text-botb-muted">
              For our Dream Car competitions you play Spot the Ball. We take a
              real match photo and remove the ball — your job is to place the
              crosshair exactly where you think the centre of the ball is.
            </p>
            <p className="mt-3 text-[15px] leading-7 text-botb-muted">
              Our independent judge decides the winning spot based on the
              players&apos; body language and eye line. The entry closest to
              that spot drives away with the car (or takes the cash
              alternative). It&apos;s a game of skill, not luck.
            </p>
          </div>
          <div className="overflow-hidden rounded-lg border border-botb-card-border bg-white">
            <Image
              src="/images/comps/13732-wide.webp"
              alt="Spot the Ball — place the crosshair where you think the ball's centre is"
              width={1200}
              height={675}
              className="h-auto w-full"
            />
          </div>
        </div>
      </section>

      {/* FAQ */}
      <section className="mb-16">
        <h2 className="mb-6 font-jost text-[24px] font-bold uppercase text-botb-text md:text-[28px]">
          Frequently asked questions
        </h2>
        <div className="divide-y divide-botb-card-border overflow-hidden rounded-lg border border-botb-card-border bg-white">
          {faqs.map((faq) => (
            <div key={faq.question} className="p-6">
              <h3 className="font-jost text-[18px] font-bold text-botb-text">
                {faq.question}
              </h3>
              <p className="mt-2 text-[15px] leading-7 text-botb-muted">
                {faq.answer}
              </p>
            </div>
          ))}
        </div>
      </section>

      {/* Closing CTA */}
      <section className="text-center">
        <h2 className="font-jost text-[24px] font-bold uppercase text-botb-text md:text-[28px]">
          Ready to play?
        </h2>
        <p className="mx-auto mt-3 max-w-xl text-[15px] leading-7 text-botb-muted">
          Browse our live competitions and pick the prize you&apos;d love to
          win.
        </p>
        <Link
          href="/competitions"
          className="mt-6 inline-block rounded bg-botb-orange px-6 py-2.5 font-jost text-[15px] font-medium uppercase text-white hover:bg-botb-orange-hover"
        >
          View Competitions
        </Link>
      </section>
    </div>
  );
}
