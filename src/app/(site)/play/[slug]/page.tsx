"use client";

import { use, useState, type MouseEvent } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useCart } from "@/lib/cart";
import { ChevronLeftIcon } from "@/components/icons";

const ACTION_IMAGE = "/images/comps/13732-wide.webp";
const UNIT_PRICE = 1.13;
const MIN_QTY = 1;
const MAX_QTY = 10;

interface Guess {
  /** Horizontal position as a percentage of the image width. */
  x: number;
  /** Vertical position as a percentage of the image height. */
  y: number;
}

function StepIndicator() {
  return (
    <div className="flex flex-wrap items-center justify-center gap-2 text-sm font-jost font-semibold sm:gap-3">
      <span className="text-botb-muted">1/3 Select Your Prizes</span>
      <span className="text-botb-muted">→</span>
      <span className="flex h-9 items-center gap-2 rounded-full bg-botb-orange px-3 text-white">
        <span className="grid h-6 w-6 place-items-center rounded-full bg-white text-xs font-bold text-botb-orange">
          2/3
        </span>
        Play the game
      </span>
      <span className="text-botb-muted">→</span>
      <span className="text-botb-muted">Win your dream car</span>
    </div>
  );
}

function Crosshair({ guess }: { guess: Guess }) {
  return (
    <span
      className="pointer-events-none absolute z-10 -translate-x-1/2 -translate-y-1/2"
      style={{ left: `${guess.x}%`, top: `${guess.y}%` }}
      aria-hidden
    >
      <span className="relative block h-10 w-10">
        <span className="absolute left-1/2 top-1/2 h-10 w-10 -translate-x-1/2 -translate-y-1/2 rounded-full border-2 border-botb-orange" />
        <span className="absolute left-1/2 top-0 h-10 w-0.5 -translate-x-1/2 bg-botb-orange" />
        <span className="absolute left-0 top-1/2 h-0.5 w-10 -translate-y-1/2 bg-botb-orange" />
        <span className="absolute left-1/2 top-1/2 h-2 w-2 -translate-x-1/2 -translate-y-1/2 rounded-full bg-botb-orange" />
      </span>
    </span>
  );
}

export default function SpotTheBallPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = use(params);
  const router = useRouter();
  const { addItem } = useCart();

  const [guess, setGuess] = useState<Guess | null>(null);
  const [qty, setQty] = useState(1);

  const handlePlaceGuess = (event: MouseEvent<HTMLButtonElement>) => {
    const rect = event.currentTarget.getBoundingClientRect();
    const x = ((event.clientX - rect.left) / rect.width) * 100;
    const y = ((event.clientY - rect.top) / rect.height) * 100;
    setGuess({ x, y });
  };

  const decQty = () => setQty((q) => Math.max(MIN_QTY, q - 1));
  const incQty = () => setQty((q) => Math.min(MAX_QTY, q + 1));

  const handleAddToBasket = () => {
    if (!guess) return;
    addItem({
      slug,
      title: "Dream Car — Spot the Ball",
      image: ACTION_IMAGE,
      unitPrice: UNIT_PRICE,
      qty,
      note: "Spot the Ball",
    });
    router.push("/cart");
  };

  return (
    <div className="container mx-auto max-w-4xl px-4 py-8">
      <Link
        href="/prizes/cars"
        className="mb-4 inline-flex items-center gap-1 text-sm font-medium text-botb-muted transition-colors hover:text-botb-text"
      >
        <ChevronLeftIcon className="h-4 w-4" />
        Back to prizes
      </Link>

      <div className="mb-6">
        <StepIndicator />
      </div>

      <h1 className="text-center font-jost text-[28px] font-bold uppercase leading-tight text-botb-orange md:text-[40px]">
        Spot the Ball
      </h1>
      <p className="mx-auto mt-2 max-w-xl text-center text-sm text-botb-muted">
        Click where you think the centre of the ball is. The closest guess wins!
      </p>

      <div className="mt-6 overflow-hidden rounded-lg border border-botb-card-border bg-botb-gray">
        <button
          type="button"
          onClick={handlePlaceGuess}
          className="relative block w-full cursor-crosshair"
          aria-label="Place your guess on the image"
        >
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src={ACTION_IMAGE}
            alt="Spot the ball action shot"
            className="w-full select-none"
            draggable={false}
          />
          {guess && <Crosshair guess={guess} />}
        </button>
      </div>

      {guess && (
        <p className="mt-3 text-center text-sm font-semibold text-botb-green">
          Your guess placed! Adjust by clicking again.
        </p>
      )}

      <div className="mt-6 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
        <div className="flex items-center gap-3">
          <span className="text-sm font-medium text-botb-text">Plays</span>
          <div className="flex items-center rounded-md border border-botb-card-border bg-white">
            <button
              type="button"
              onClick={decQty}
              disabled={qty <= MIN_QTY}
              className="grid h-10 w-10 place-items-center text-lg font-bold text-botb-text transition-colors hover:bg-botb-gray disabled:opacity-40"
              aria-label="Decrease plays"
            >
              −
            </button>
            <span className="grid h-10 w-12 place-items-center font-jost text-lg font-bold text-botb-text">
              {qty}
            </span>
            <button
              type="button"
              onClick={incQty}
              disabled={qty >= MAX_QTY}
              className="grid h-10 w-10 place-items-center text-lg font-bold text-botb-text transition-colors hover:bg-botb-gray disabled:opacity-40"
              aria-label="Increase plays"
            >
              +
            </button>
          </div>
        </div>

        <button
          type="button"
          onClick={handleAddToBasket}
          disabled={!guess}
          className="inline-flex h-12 items-center justify-center rounded-md bg-botb-orange px-8 font-jost text-base font-bold uppercase text-white transition-colors hover:bg-botb-orange-hover disabled:cursor-not-allowed disabled:opacity-50"
        >
          Add to Basket »
        </button>
      </div>

      <p className="mt-4 text-center text-xs text-botb-muted">
        Total: £{(UNIT_PRICE * qty).toFixed(2)} for {qty}{" "}
        {qty === 1 ? "play" : "plays"}
        {!guess && " · place your guess to continue"}
      </p>

      <p className="mt-6 text-center text-sm text-botb-muted">
        New to Spot the Ball?{" "}
        <Link
          href="/how-to-play"
          className="font-semibold text-botb-orange hover:underline"
        >
          See how it works
        </Link>
      </p>
    </div>
  );
}
