type Perk = {
  icon: string;
  title: string;
  description: string;
};

const perks: Perk[] = [
  {
    icon: "🎟️",
    title: "Free tickets",
    description:
      "Collect free entries the more you play. Every active week unlocks bonus tickets towards your favourite competitions.",
  },
  {
    icon: "🎂",
    title: "Birthday bonus",
    description:
      "Celebrate with us. Players receive a special birthday treat and extra entries during their birthday month.",
  },
  {
    icon: "🏆",
    title: "Loyalty rewards",
    description:
      "The longer you're with BOTB, the more you earn. Climb the loyalty tiers for ever-growing rewards.",
  },
  {
    icon: "⭐",
    title: "Exclusive comps",
    description:
      "Unlock members-only competitions with better odds and prizes reserved just for our loyal players.",
  },
  {
    icon: "🚀",
    title: "Early access",
    description:
      "Be first in line. Get early access to brand-new Dream Cars and Lifestyle prizes before anyone else.",
  },
  {
    icon: "🤝",
    title: "Referral bonus",
    description:
      "Invite friends and you both win. Earn bonus entries every time a friend joins and plays with BOTB.",
  },
];

export default function PerksForPlayingPage() {
  return (
    <section className="bg-white">
      <div className="mx-auto max-w-[1360px] px-4 py-10">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          Perks for Playing
        </h1>
        <p className="mt-2 max-w-3xl text-botb-muted">
          Playing BOTB comes with more than the chance to win. From free tickets
          to exclusive competitions, here&apos;s everything our players enjoy
          just for taking part.
        </p>

        <div className="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {perks.map((perk) => (
            <div
              key={perk.title}
              className="rounded-md border border-botb-card-border p-6"
            >
              <div
                aria-hidden
                className="flex h-12 w-12 items-center justify-center rounded-full bg-botb-gray text-[24px]"
              >
                {perk.icon}
              </div>
              <h2 className="mt-4 font-jost text-[18px] font-semibold text-botb-text">
                {perk.title}
              </h2>
              <p className="mt-2 text-[14px] text-botb-muted">
                {perk.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
