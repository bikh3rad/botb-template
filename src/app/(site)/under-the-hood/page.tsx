import Link from "next/link";

type BlogPost = {
  title: string;
  excerpt: string;
  date: string;
  category: string;
  image: string;
};

const posts: BlogPost[] = [
  {
    title: "Meet the winner who drove away in a brand-new supercar",
    excerpt:
      "We caught up with last week's Dream Car winner to hear what it felt like when our team turned up on the doorstep.",
    date: "26 June 2026",
    category: "Winner Story",
    image: "/images/comps/10168.webp",
  },
  {
    title: "Behind the scenes: how we film every single winner reveal",
    excerpt:
      "From the surprise knock to the keys in hand, here's how the BOTB team captures those unforgettable moments.",
    date: "21 June 2026",
    category: "Behind the Scenes",
    image: "/images/comps/10415.webp",
  },
  {
    title: "Cash, gadgets and getaways — inside our Lifestyle competitions",
    excerpt:
      "Cars aren't the only thing up for grabs. Take a look at the huge range of Lifestyle prizes won every week.",
    date: "18 June 2026",
    category: "Prizes",
    image: "/images/comps/10873.webp",
  },
  {
    title: "Five spot-the-ball tips from our most successful players",
    excerpt:
      "Want to sharpen your eye? Our community shares the habits that keep them coming back to the winners' circle.",
    date: "14 June 2026",
    category: "Tips & Tricks",
    image: "/images/comps/11121.webp",
  },
  {
    title: "How BOTB players have raised millions for charity",
    excerpt:
      "Every entry helps. Here's a look at the causes our community has supported through the BOTB Foundation.",
    date: "9 June 2026",
    category: "Charity",
    image: "/images/comps/11196.webp",
  },
  {
    title: "A weekend with the latest electric Dream Car winner",
    excerpt:
      "Range anxiety? Not a chance. Our newest winner takes their silent supercar on its very first road trip.",
    date: "4 June 2026",
    category: "Winner Story",
    image: "/images/comps/11364.webp",
  },
  {
    title: "New competitions just launched — here's what to watch",
    excerpt:
      "Fresh prizes drop every week. We round up the most exciting new competitions you won't want to miss.",
    date: "30 May 2026",
    category: "News",
    image: "/images/comps/11402.webp",
  },
  {
    title: "The story behind our London Gatwick HQ",
    excerpt:
      "Ever wondered where the magic happens? Take a tour of the home of BOTB and the team that makes it tick.",
    date: "24 May 2026",
    category: "Behind the Scenes",
    image: "/images/comps/11668.webp",
  },
  {
    title: "From sceptic to superfan: one player's winning journey",
    excerpt:
      "He almost didn't enter. Now he's part of the family. Read how a single ticket changed everything.",
    date: "19 May 2026",
    category: "Winner Story",
    image: "/images/comps/11858.webp",
  },
];

export default function UnderTheHoodPage() {
  return (
    <section className="bg-white">
      <div className="mx-auto max-w-[1360px] px-4 py-10">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          Under the Hood
        </h1>
        <p className="mt-2 text-botb-muted">
          News, winner stories &amp; behind the scenes.
        </p>

        <div className="mt-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {posts.map((post) => (
            <Link
              key={post.title}
              href="/under-the-hood"
              className="group flex flex-col overflow-hidden rounded-md border border-botb-card-border transition-shadow hover:shadow-md"
            >
              <div className="aspect-[16/10] overflow-hidden bg-botb-gray">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={post.image}
                  alt={post.title}
                  className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
                />
              </div>
              <div className="flex flex-1 flex-col p-5">
                <span className="self-start rounded-full bg-botb-gray px-3 py-1 text-[12px] font-medium uppercase text-botb-secondary">
                  {post.category}
                </span>
                <h2 className="mt-3 font-jost text-[18px] font-semibold text-botb-text group-hover:text-botb-orange">
                  {post.title}
                </h2>
                <p className="mt-2 flex-1 text-[14px] text-botb-muted">
                  {post.excerpt}
                </p>
                <p className="mt-4 text-[13px] text-botb-muted">{post.date}</p>
              </div>
            </Link>
          ))}
        </div>
      </div>
    </section>
  );
}
