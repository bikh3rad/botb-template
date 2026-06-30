import Link from "next/link";

interface PrizeCollection {
  title: string;
  count: number;
  image: string;
}

const collections: PrizeCollection[] = [
  { title: "Dream Cars", count: 24, image: "/images/comps/13558.webp" },
  { title: "Cash Prizes", count: 18, image: "/images/comps/9211-money.webp" },
  { title: "Luxury Watches", count: 9, image: "/images/comps/13344.webp" },
  { title: "Tech & Gadgets", count: 31, image: "/images/comps/13767.webp" },
  { title: "Houses", count: 4, image: "/images/comps/13982.webp" },
  { title: "Lifestyle", count: 15, image: "/images/comps/11364.webp" },
];

export default function PrizeCollectionsPage() {
  return (
    <div className="mx-auto max-w-[1360px] px-4 py-10">
      <header className="mb-8">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          Prize Collections
        </h1>
        <p className="mt-2 text-botb-muted">Browse our prize categories.</p>
      </header>

      <section className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {collections.map((collection) => (
          <Link
            key={collection.title}
            href="/competitions"
            className="group relative block overflow-hidden rounded-lg border border-botb-card-border shadow-sm transition-transform duration-200 hover:-translate-y-1 hover:shadow-md"
          >
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={collection.image}
              alt={collection.title}
              className="aspect-[4/3] w-full object-cover transition-transform duration-300 group-hover:scale-105"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/75 via-black/20 to-transparent" />
            <div className="absolute bottom-0 left-0 p-5">
              <h2 className="font-jost text-xl font-bold text-white">
                {collection.title}
              </h2>
              <p className="mt-1 text-sm font-medium text-white/80">
                {collection.count} prizes
              </p>
            </div>
          </Link>
        ))}
      </section>
    </div>
  );
}
