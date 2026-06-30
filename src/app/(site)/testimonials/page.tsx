type Testimonial = {
  quote: string;
  name: string;
  prize: string;
  rating: number;
};

const testimonials: Testimonial[] = [
  {
    quote:
      "I genuinely didn't believe it when the team turned up at my door with the keys. The whole process was smooth and the car was exactly as advertised. Still pinching myself!",
    name: "Daniel R.",
    prize: "Won a Porsche 911 Carrera",
    rating: 5,
  },
  {
    quote:
      "Played for years for a bit of fun and never expected to actually win. £25,000 landed in my account within days — no fuss, no catch. Absolutely brilliant.",
    name: "Sophie M.",
    prize: "Won £25,000 Cash",
    rating: 5,
  },
  {
    quote:
      "The lifestyle prize changed everything for us. We put the money towards a deposit and finally got on the property ladder. Forever grateful to BOTB.",
    name: "James & Hannah T.",
    prize: "Won a £200,000 Dream Home Fund",
    rating: 5,
  },
  {
    quote:
      "Won the full tech bundle — new laptop, phone and a massive TV. Delivery was quick and everything was brand new and boxed. Couldn't ask for more.",
    name: "Aisha K.",
    prize: "Won the Ultimate Tech Bundle",
    rating: 5,
  },
  {
    quote:
      "I've recommended BOTB to all my mates. The weekly draws are exciting and when I won my Audi the customer service was second to none.",
    name: "Mark P.",
    prize: "Won an Audi RS3",
    rating: 5,
  },
  {
    quote:
      "Never won anything in my life until now. £10k cash just before Christmas — it covered all the presents and a family holiday. Made our year.",
    name: "Lucy B.",
    prize: "Won £10,000 Cash",
    rating: 5,
  },
  {
    quote:
      "The dream car came with a tank of fuel and a year's insurance sorted. They really do think of everything. A proper, trustworthy company.",
    name: "Omar S.",
    prize: "Won a BMW M4 Competition",
    rating: 5,
  },
  {
    quote:
      "From entering to winning to delivery, every step was clear and professional. The team even filmed the surprise — such a lovely touch.",
    name: "Charlotte W.",
    prize: "Won a Tesla Model 3",
    rating: 5,
  },
  {
    quote:
      "I was over the moon to win the cash alternative and clear my debts. Honestly life-changing. Thank you BOTB for being the real deal.",
    name: "Gary H.",
    prize: "Won £50,000 Cash",
    rating: 5,
  },
];

export default function TestimonialsPage() {
  return (
    <div className="mx-auto max-w-[1360px] px-4 py-10">
      <header className="mb-8">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          Testimonials
        </h1>
        <p className="mt-2 text-botb-muted">
          Don&apos;t just take our word for it.
        </p>
      </header>

      <section className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        {testimonials.map((testimonial, index) => (
          <TestimonialCard
            key={`${testimonial.name}-${index}`}
            testimonial={testimonial}
          />
        ))}
      </section>

      <section className="mt-12 rounded-lg bg-botb-gray p-8 text-center">
        <div className="flex flex-col items-center justify-center gap-6 sm:flex-row sm:gap-12">
          <div>
            <div className="font-jost text-xl font-bold text-botb-text">
              feefo{" "}
              <span className="text-botb-orange" aria-hidden="true">
                ★★★★★
              </span>{" "}
              4.8
            </div>
            <div className="text-sm text-botb-muted">Independent reviews</div>
          </div>
          <div>
            <div className="font-jost text-xl font-bold text-botb-text">
              Trustpilot{" "}
              <span className="text-botb-orange" aria-hidden="true">
                ★★★★★
              </span>
            </div>
            <div className="text-sm text-botb-muted">Rated Excellent</div>
          </div>
          <div>
            <div className="font-jost text-xl font-bold text-botb-text">
              Over 721k winners
            </div>
            <div className="text-sm text-botb-muted">and counting</div>
          </div>
        </div>
      </section>
    </div>
  );
}

function TestimonialCard({ testimonial }: { testimonial: Testimonial }) {
  return (
    <article className="flex flex-col rounded-lg border border-botb-card-border bg-white p-6 shadow-sm">
      <div
        className="text-lg text-botb-orange"
        aria-label={`${testimonial.rating} out of 5 stars`}
      >
        {"★".repeat(testimonial.rating)}
      </div>
      <p className="mt-4 flex-1 text-botb-text">{testimonial.quote}</p>
      <footer className="mt-4 border-t border-botb-card-border pt-4">
        <div className="font-jost font-semibold text-botb-text">
          {testimonial.name}
        </div>
        <div className="text-sm text-botb-muted">{testimonial.prize}</div>
      </footer>
    </article>
  );
}
