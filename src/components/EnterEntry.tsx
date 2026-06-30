"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useCart } from "@/lib/cart";
import { CloseIcon } from "@/components/icons";

interface EnterEntryProps {
  slug: string;
  /** Display title for the basket line and modal header. */
  title: string;
  image: string;
  /** Price string from the data, e.g. "£1.25". */
  price: string;
}

const MIN_QTY = 1;
const MAX_QTY = 1000;

/** Quick-select shortcuts shown in the entry modal. */
const QUICK_SELECT: { qty: number; label?: string }[] = [
  { qty: 10 },
  { qty: 20 },
  { qty: 50 },
  { qty: 1000, label: "MOST CHANCES" },
];

function clampQty(value: number): number {
  if (Number.isNaN(value)) return MIN_QTY;
  return Math.min(MAX_QTY, Math.max(MIN_QTY, Math.round(value)));
}

export function EnterEntry({ slug, title, image, price }: EnterEntryProps) {
  const router = useRouter();
  const { addItem } = useCart();

  // Derive a numeric unit price from the display string (e.g. "£1.25" -> 1.25).
  const unit = parseFloat(price.replace(/[^0-9.]/g, "")) || 0;

  const [open, setOpen] = useState(false);
  const [qty, setQty] = useState(MIN_QTY);
  const [tab, setTab] = useState<"online" | "postal">("online");

  const total = (unit * qty).toFixed(2);

  function handleAddToBasket() {
    addItem({ slug, title, image, unitPrice: unit, qty });
    setOpen(false);
    router.push("/cart");
  }

  return (
    <>
      {/* CTA — full width, becomes a sticky bottom bar on small screens */}
      <div className="sticky bottom-0 z-20 -mx-4 border-t border-botb-card-border bg-white/95 px-4 py-3 backdrop-blur md:static md:mx-0 md:border-0 md:bg-transparent md:p-0 md:backdrop-blur-none">
        <button
          type="button"
          onClick={() => setOpen(true)}
          className="w-full rounded bg-botb-orange px-6 py-3.5 font-jost text-[16px] font-semibold uppercase tracking-wide text-white transition-colors hover:bg-botb-orange-hover"
        >
          Enter Now »
        </button>
      </div>

      {open && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
          role="dialog"
          aria-modal="true"
          aria-label={`Enter ${title}`}
          onClick={() => setOpen(false)}
        >
          <div
            className="relative w-full max-w-md rounded-lg bg-white p-6 shadow-lg"
            onClick={(e) => e.stopPropagation()}
          >
            <button
              type="button"
              onClick={() => setOpen(false)}
              aria-label="Close"
              className="absolute right-3 top-3 rounded-full p-1 text-botb-muted transition-colors hover:bg-botb-gray hover:text-botb-text"
            >
              <CloseIcon className="h-5 w-5" />
            </button>

            {/* Entry-method tabs */}
            <div className="mb-4 flex gap-6 border-b border-botb-card-border pr-8">
              <button
                type="button"
                onClick={() => setTab("online")}
                className={`-mb-px border-b-2 pb-2 font-jost text-[14px] font-semibold uppercase transition-colors ${
                  tab === "online"
                    ? "border-botb-orange text-botb-orange"
                    : "border-transparent text-botb-muted hover:text-botb-text"
                }`}
              >
                Online entry
              </button>
              <button
                type="button"
                onClick={() => setTab("postal")}
                className={`-mb-px border-b-2 pb-2 font-jost text-[14px] font-semibold uppercase transition-colors ${
                  tab === "postal"
                    ? "border-botb-orange text-botb-orange"
                    : "border-transparent text-botb-muted hover:text-botb-text"
                }`}
              >
                Free postal entry
              </button>
            </div>

            {/* Prize summary */}
            <div className="mb-5 flex items-center gap-3">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src={image}
                alt={title}
                className="h-16 w-16 flex-none rounded-md object-cover"
              />
              <div className="min-w-0">
                <p className="truncate font-jost text-[15px] font-bold text-botb-text">
                  {title}
                </p>
                <p className="text-[13px] text-botb-muted">{price} per ticket</p>
              </div>
            </div>

            {/* Quantity */}
            <div className="mb-4">
              <div className="mb-2 flex items-center justify-between">
                <span className="text-[13px] font-medium uppercase tracking-wide text-botb-muted">
                  Quantity
                </span>
                <span className="font-jost text-[20px] font-bold text-botb-text">
                  {qty}
                </span>
              </div>

              <div className="flex items-center gap-3">
                <button
                  type="button"
                  onClick={() => setQty((q) => clampQty(q - 1))}
                  aria-label="Decrease quantity"
                  className="flex h-9 w-9 flex-none items-center justify-center rounded-full border border-botb-card-border text-[20px] leading-none text-botb-text transition-colors hover:bg-botb-gray"
                >
                  −
                </button>
                <input
                  type="range"
                  min={MIN_QTY}
                  max={MAX_QTY}
                  value={qty}
                  onChange={(e) => setQty(clampQty(Number(e.target.value)))}
                  aria-label="Quantity"
                  className="h-2 flex-1 cursor-pointer accent-botb-orange"
                />
                <button
                  type="button"
                  onClick={() => setQty((q) => clampQty(q + 1))}
                  aria-label="Increase quantity"
                  className="flex h-9 w-9 flex-none items-center justify-center rounded-full border border-botb-card-border text-[20px] leading-none text-botb-text transition-colors hover:bg-botb-gray"
                >
                  +
                </button>
              </div>
            </div>

            {/* Quick select */}
            <div className="mb-5">
              <p className="mb-2 text-[13px] font-medium uppercase tracking-wide text-botb-muted">
                Quick Select
              </p>
              <div className="grid grid-cols-4 gap-2">
                {QUICK_SELECT.map(({ qty: q, label }) => (
                  <button
                    key={q}
                    type="button"
                    onClick={() => setQty(clampQty(q))}
                    className={`flex flex-col items-center justify-center rounded border px-2 py-2 font-jost text-[14px] font-bold transition-colors ${
                      qty === q
                        ? "border-botb-orange bg-botb-orange text-white"
                        : "border-botb-card-border bg-white text-botb-text hover:border-botb-orange hover:text-botb-orange"
                    }`}
                  >
                    {q}
                    {label && (
                      <span className="mt-0.5 text-[8px] font-semibold uppercase leading-none">
                        {label}
                      </span>
                    )}
                  </button>
                ))}
              </div>
            </div>

            {/* Total + add to basket */}
            <div className="mb-4 flex items-center justify-between border-t border-botb-card-border pt-4">
              <span className="text-[14px] font-medium uppercase text-botb-muted">
                Total
              </span>
              <span className="font-jost text-[22px] font-bold text-botb-text">
                £{total}
              </span>
            </div>

            <button
              type="button"
              onClick={handleAddToBasket}
              className="w-full rounded bg-botb-orange px-6 py-3.5 font-jost text-[16px] font-semibold uppercase tracking-wide text-white transition-colors hover:bg-botb-orange-hover"
            >
              Add to Basket »
            </button>
          </div>
        </div>
      )}
    </>
  );
}
