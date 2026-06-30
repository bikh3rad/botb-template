"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { CartIcon, ChevronRightIcon } from "@/components/icons";

interface CarPrize {
  id: string;
  title: string;
  sub: string;
  cash: string;
  price: string;
  was?: string;
  image: string;
  discount?: string;
  doubleUp?: boolean;
}

/** Mock prize catalogue mirroring BOTB's Dream Car selection page. */
const CAR_PRIZES: CarPrize[] = [
  {
    id: "9211-money",
    title: "£250,000 Cash Prize",
    sub: "Tax-free cash, straight to your account",
    cash: "£250,000",
    price: "£5.00",
    was: "£9.25",
    discount: "45% OFF",
    image: "/images/comps/9211-money.webp",
  },
  {
    id: "13558",
    title: "Double Up - Aston Martin Vantage & McLaren",
    sub: "Two supercars, one winner",
    cash: "£195,000",
    price: "£4.63",
    image: "/images/comps/13558.webp",
    doubleUp: true,
  },
  {
    id: "13740",
    title: "Double Up - Lotus Exige + Bentley",
    sub: "Track day legend meets grand tourer",
    cash: "£168,000",
    price: "£4.63",
    image: "/images/comps/13740.webp",
    doubleUp: true,
  },
  {
    id: "13736",
    title: "Jaguar F-Type + Land Rover Defender 90",
    sub: "Road and rough, covered",
    cash: "£142,500",
    price: "£4.34",
    image: "/images/hero/slide-13736.webp",
  },
  {
    id: "13734",
    title: "Autotrail Expedition + MG S5",
    sub: "Adventure camper plus a runabout",
    cash: "£118,000",
    price: "£3.99",
    image: "/images/hero/slide-13734.webp",
  },
  {
    id: "8941-1000cash",
    title: "£104,720 Cash Prize",
    sub: "Take the lot in cash",
    cash: "£104,720",
    price: "£3.50",
    image: "/images/comps/8941-1000cash.webp",
  },
  {
    id: "13728",
    title: "Porsche 911 Carrera",
    sub: "The everyday icon",
    cash: "£96,400",
    price: "£2.99",
    image: "/images/hero/slide-13728.webp",
  },
  {
    id: "13740-mini",
    title: "Mini Cooper S + £10,000 Cash",
    sub: "City classic with spending money",
    cash: "£38,500",
    price: "£2.10",
    image: "/images/comps/13740.webp",
  },
  {
    id: "renault-5",
    title: "Renault 5 Iconic",
    sub: "Retro reborn, all electric",
    cash: "£26,995",
    price: "£2.10",
    image: "/images/comps/13558.webp",
  },
];

type SortKey = "featured" | "price-asc" | "price-desc";

/** Parse a "£4.63" style string into a comparable number. */
function toNumber(value: string): number {
  return Number(value.replace(/[^0-9.]/g, "")) || 0;
}

function StepIndicator() {
  return (
    <div className="flex flex-wrap items-center justify-center gap-2 text-sm font-jost font-semibold text-botb-text sm:gap-3">
      <span className="flex h-9 items-center gap-2 rounded-full bg-botb-orange px-3 text-white">
        <span className="grid h-6 w-6 place-items-center rounded-full bg-white text-xs font-bold text-botb-orange">
          1/3
        </span>
        Select Your Prizes
      </span>
      <ChevronRightIcon className="h-4 w-4 text-botb-muted" />
      <span className="text-botb-muted">Play the game</span>
      <ChevronRightIcon className="h-4 w-4 text-botb-muted" />
      <span className="flex items-center gap-1.5 text-botb-muted">
        <CarIcon className="h-5 w-5" />
        Win your dream car
      </span>
    </div>
  );
}

function CarIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" className={className} aria-hidden>
      <path
        d="M3 13l1.5-4.5A2 2 0 0 1 6.4 7h11.2a2 2 0 0 1 1.9 1.5L21 13m-18 0v4a1 1 0 0 0 1 1h1a1 1 0 0 0 1-1v-1h12v1a1 1 0 0 0 1 1h1a1 1 0 0 0 1-1v-4m-18 0h18"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <circle cx="7" cy="16" r="1" fill="currentColor" />
      <circle cx="17" cy="16" r="1" fill="currentColor" />
    </svg>
  );
}

function CountdownBox({ value, label }: { value: string; label: string }) {
  return (
    <div className="flex flex-col items-center">
      <span className="grid h-12 w-12 place-items-center rounded-md bg-botb-text font-jost text-xl font-bold text-white sm:h-14 sm:w-14 sm:text-2xl">
        {value}
      </span>
      <span className="mt-1 text-[11px] uppercase tracking-wide text-botb-muted">
        {label}
      </span>
    </div>
  );
}

function SearchIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" className={className} aria-hidden>
      <circle cx="11" cy="11" r="7" stroke="currentColor" strokeWidth="2" />
      <path d="M20 20l-3-3" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
    </svg>
  );
}

function SortIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" className={className} aria-hidden>
      <path
        d="M3 6h18M6 12h12M10 18h4"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
      />
    </svg>
  );
}

export default function DreamCarsPage() {
  const [query, setQuery] = useState("");
  const [sort, setSort] = useState<SortKey>("featured");

  const visiblePrizes = useMemo(() => {
    const term = query.trim().toLowerCase();
    const filtered = term
      ? CAR_PRIZES.filter(
          (p) =>
            p.title.toLowerCase().includes(term) ||
            p.sub.toLowerCase().includes(term),
        )
      : CAR_PRIZES;

    if (sort === "featured") return filtered;
    const sorted = [...filtered].sort(
      (a, b) => toNumber(a.price) - toNumber(b.price),
    );
    return sort === "price-asc" ? sorted : sorted.reverse();
  }, [query, sort]);

  return (
    <div className="container mx-auto max-w-6xl px-4 py-8">
      <div className="mb-6">
        <StepIndicator />
      </div>

      <h1 className="text-center font-jost text-[28px] font-bold uppercase leading-tight text-botb-orange md:text-[40px]">
        Dream Car Competition - All Cars Brand New
      </h1>

      <div className="mt-4 flex flex-col items-center gap-2">
        <span className="font-jost text-sm font-semibold uppercase text-botb-muted">
          Ends in
        </span>
        <div className="flex items-start gap-2 sm:gap-3">
          <CountdownBox value="00" label="Days" />
          <span className="pt-3 font-jost text-xl font-bold text-botb-muted sm:pt-4">:</span>
          <CountdownBox value="00" label="Hrs" />
          <span className="pt-3 font-jost text-xl font-bold text-botb-muted sm:pt-4">:</span>
          <CountdownBox value="00" label="Mins" />
          <span className="pt-3 font-jost text-xl font-bold text-botb-muted sm:pt-4">:</span>
          <CountdownBox value="00" label="Secs" />
        </div>
      </div>

      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src="/images/hero/slide-13736.webp"
        alt="Dream car competition banner"
        className="mt-6 h-[200px] w-full rounded-lg object-cover md:h-[280px]"
      />

      <div className="mt-6 flex flex-col gap-3 sm:flex-row sm:items-center">
        <div className="relative flex-1">
          <SearchIcon className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-botb-muted" />
          <input
            type="search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search dream cars..."
            className="h-11 w-full rounded-md border border-botb-card-border bg-white pl-10 pr-4 text-sm text-botb-text outline-none focus:border-botb-orange"
          />
        </div>
        <div className="flex items-center gap-2">
          <select
            value={sort}
            onChange={(e) => setSort(e.target.value as SortKey)}
            className="h-11 rounded-md border border-botb-card-border bg-white px-3 text-sm text-botb-text outline-none focus:border-botb-orange"
            aria-label="Sort prizes"
          >
            <option value="featured">Featured</option>
            <option value="price-asc">Price: Low to High</option>
            <option value="price-desc">Price: High to Low</option>
          </select>
          <button
            type="button"
            className="grid h-11 w-11 place-items-center rounded-md border border-botb-card-border bg-white text-botb-text transition-colors hover:bg-botb-gray"
            aria-label="Filter"
          >
            <SortIcon className="h-5 w-5" />
          </button>
        </div>
      </div>

      <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {visiblePrizes.map((prize) => (
          <Link
            key={prize.id}
            href={`/play/${prize.id}`}
            className="group relative flex flex-col overflow-hidden rounded-lg border border-botb-card-border bg-white shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md"
          >
            {prize.doubleUp && (
              <span className="absolute left-3 top-3 z-10 rounded-full bg-botb-badge px-2.5 py-0.5 text-xs font-bold uppercase text-white">
                Double Ups!
              </span>
            )}
            {prize.discount && (
              <span className="absolute right-3 top-3 z-10 rounded-full bg-botb-green px-2.5 py-0.5 text-xs font-bold uppercase text-white">
                {prize.discount}
              </span>
            )}

            <div className="bg-botb-gray p-4">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src={prize.image}
                alt={prize.title}
                className="h-40 w-full object-contain transition-transform duration-200 group-hover:scale-105"
              />
            </div>

            <div className="flex flex-1 flex-col p-4">
              <h2 className="font-jost text-base font-semibold leading-tight text-botb-text">
                {prize.title}
              </h2>
              <p className="mt-1 text-sm text-botb-muted">
                Or take {prize.cash} cash
              </p>

              <div className="mt-auto flex items-end justify-between pt-4">
                <div className="flex flex-col">
                  {prize.was && (
                    <span className="text-xs text-botb-muted line-through">
                      was {prize.was}
                    </span>
                  )}
                  <span className="font-jost text-xl font-bold text-botb-text">
                    {prize.price}
                  </span>
                </div>
                <span className="inline-flex items-center gap-1.5 rounded-md bg-botb-orange px-3 py-2 text-sm font-semibold text-white transition-colors group-hover:bg-botb-orange-hover">
                  <CartIcon className="h-4 w-4" />
                  Play
                </span>
              </div>
            </div>
          </Link>
        ))}
      </div>

      {visiblePrizes.length === 0 && (
        <p className="mt-12 text-center text-botb-muted">
          No dream cars match your search.
        </p>
      )}
    </div>
  );
}
