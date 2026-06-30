type Milestone = {
  year: string;
  text: string;
};

const milestones: Milestone[] = [
  {
    year: "1999",
    text: "BOTB is founded with a single dream-car stand at London Gatwick Airport, inviting travellers to spot the ball and drive away in their dream car.",
  },
  {
    year: "2000s",
    text: "The instantly recognisable airport stands roll out across major UK terminals, turning a clever idea into a household name for jet-setting players.",
  },
  {
    year: "2014",
    text: "BOTB lists on the London Stock Exchange's AIM market, cementing its reputation as a trusted, transparent and growing British business.",
  },
  {
    year: "2020",
    text: "A surge in online play transforms the business — players from across the country join weekly draws from their phones and laptops.",
  },
  {
    year: "2023",
    text: "BOTB passes a remarkable milestone, having handed over more than £160M in prizes to lucky winners throughout its history.",
  },
  {
    year: "2025",
    text: "Celebrating 26 years and over 721k guaranteed winners — and still giving away dream cars, cash and lifestyle prizes every single week.",
  },
];

const locations = [
  "London Gatwick",
  "Heathrow",
  "Manchester",
  "Birmingham",
  "Online HQ",
];

export default function HistoryLocationsPage() {
  return (
    <div className="mx-auto max-w-4xl px-4 py-10">
      <header className="mb-10">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          Our History
        </h1>
        <p className="mt-2 text-botb-muted">
          Est. 1999 — 26 years of real winners.
        </p>
      </header>

      <section className="relative">
        <ol className="relative ml-3 border-l-2 border-botb-card-border">
          {milestones.map((milestone) => (
            <li key={milestone.year} className="relative mb-10 pl-8 last:mb-0">
              <span
                className="absolute -left-[9px] top-1.5 h-4 w-4 rounded-full border-2 border-white bg-botb-orange"
                aria-hidden="true"
              />
              <div className="font-jost text-xl font-bold text-botb-orange">
                {milestone.year}
              </div>
              <p className="mt-1 text-botb-text">{milestone.text}</p>
            </li>
          ))}
        </ol>
      </section>

      <section className="mt-16">
        <h2 className="font-jost text-2xl font-bold uppercase text-botb-text">
          Locations
        </h2>
        <p className="mt-2 max-w-2xl text-botb-muted">
          From our first airport stand to a thriving online operation, BOTB has
          grown a proud UK-wide presence. Our team operates across major British
          airports and our digital HQ keeps the dream alive for players
          everywhere.
        </p>
        <div className="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
          {locations.map((location) => (
            <div
              key={location}
              className="rounded-lg border border-botb-card-border bg-botb-gray p-4 text-center font-jost font-semibold text-botb-text"
            >
              {location}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
