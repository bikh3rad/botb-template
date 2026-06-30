"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useCart } from "@/lib/cart";
import { CartIcon, CloseIcon } from "@/components/icons";

export default function CartPage() {
  const router = useRouter();
  const { items, total, removeItem, updateQty } = useCart();

  // The cart provider starts empty on the server and on the first client
  // render, then loads persisted items in its own effect — so rendering the
  // empty state directly stays hydration-safe.
  const isEmpty = items.length === 0;

  return (
    <div className="mx-auto max-w-[1100px] px-4 py-10">
      <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
        Your Basket
      </h1>

      {isEmpty ? (
        <div className="flex flex-col items-center justify-center gap-4 py-20 text-center">
          <CartIcon className="h-12 w-12 text-botb-muted" />
          <p className="font-jost text-lg font-semibold text-botb-text">
            Your basket is empty
          </p>
          <Link
            href="/competitions"
            className="rounded-md bg-botb-orange px-6 py-3 font-jost text-sm font-bold uppercase text-white transition-colors hover:bg-botb-orange-hover"
          >
            Browse Competitions
          </Link>
        </div>
      ) : (
        <div className="mt-8 grid grid-cols-1 gap-8 lg:grid-cols-[1fr_360px]">
          {/* Line items */}
          <div className="flex flex-col gap-4">
            {items.map((item) => (
              <div
                key={item.slug}
                className="flex gap-4 rounded-lg border border-botb-card-border bg-white p-4"
              >
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={item.image}
                  alt={item.title}
                  className="h-16 w-24 rounded object-cover"
                />

                <div className="min-w-0 flex-1">
                  <h2 className="font-jost font-semibold text-botb-text">
                    {item.title}
                  </h2>
                  {item.note ? (
                    <p className="text-xs text-botb-muted">{item.note}</p>
                  ) : null}
                  <p className="mt-1 text-sm text-botb-muted">
                    £{item.unitPrice.toFixed(2)} each
                  </p>

                  {/* Quantity stepper */}
                  <div className="mt-3 inline-flex items-center rounded-md border border-botb-card-border">
                    <button
                      type="button"
                      aria-label="Decrease quantity"
                      onClick={() => updateQty(item.slug, item.qty - 1)}
                      className="px-3 py-1 text-lg leading-none text-botb-text transition-colors hover:bg-botb-gray"
                    >
                      −
                    </button>
                    <span className="min-w-8 px-2 text-center font-jost text-sm font-semibold text-botb-text">
                      {item.qty}
                    </span>
                    <button
                      type="button"
                      aria-label="Increase quantity"
                      onClick={() => updateQty(item.slug, item.qty + 1)}
                      className="px-3 py-1 text-lg leading-none text-botb-text transition-colors hover:bg-botb-gray"
                    >
                      +
                    </button>
                  </div>
                </div>

                <div className="flex flex-col items-end justify-between">
                  <button
                    type="button"
                    aria-label={`Remove ${item.title}`}
                    onClick={() => removeItem(item.slug)}
                    className="text-botb-muted transition-colors hover:text-botb-text"
                  >
                    <CloseIcon className="h-5 w-5" />
                  </button>
                  <p className="font-jost font-bold text-botb-text">
                    £{(item.unitPrice * item.qty).toFixed(2)}
                  </p>
                </div>
              </div>
            ))}
          </div>

          {/* Order summary */}
          <aside className="lg:sticky lg:top-6 lg:self-start">
            <div className="rounded-lg border border-botb-card-border bg-white p-6">
              <h2 className="font-jost text-lg font-semibold text-botb-text">
                Order Summary
              </h2>

              <div className="mt-4 space-y-3 text-sm">
                <div className="flex justify-between text-botb-text">
                  <span>Subtotal</span>
                  <span>£{total.toFixed(2)}</span>
                </div>
                <div className="flex justify-between text-botb-text">
                  <span>Delivery</span>
                  <span className="font-semibold text-botb-orange">FREE</span>
                </div>
                <div className="border-t border-botb-card-border pt-3">
                  <div className="flex justify-between font-jost text-base font-bold text-botb-text">
                    <span>Total</span>
                    <span>£{total.toFixed(2)}</span>
                  </div>
                </div>
              </div>

              <button
                type="button"
                disabled={items.length === 0}
                onClick={() => router.push("/checkout")}
                className="mt-6 w-full rounded-md bg-botb-orange px-6 py-3 font-jost text-sm font-bold uppercase text-white transition-colors hover:bg-botb-orange-hover disabled:cursor-not-allowed disabled:opacity-50"
              >
                Proceed to Checkout
              </button>

              <Link
                href="/competitions"
                className="mt-4 block text-center text-sm text-botb-muted transition-colors hover:text-botb-text"
              >
                Continue shopping
              </Link>
            </div>
          </aside>
        </div>
      )}
    </div>
  );
}
