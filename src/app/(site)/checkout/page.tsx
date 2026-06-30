"use client";

import { type FormEvent } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useCart } from "@/lib/cart";

const INPUT_CLASS =
  "w-full rounded-md border border-botb-card-border px-3 py-2 outline-none focus:border-botb-orange";

export default function CheckoutPage() {
  const router = useRouter();
  const { items, total, clear } = useCart();

  // The cart provider starts empty on the server and on the first client
  // render, then loads persisted items in its own effect — so rendering the
  // empty state directly stays hydration-safe.
  const isEmpty = items.length === 0;

  // Mock checkout: never posts anywhere. Clears the basket and forwards to the
  // confirmation page.
  function handlePlaceOrder(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    clear();
    router.push("/checkout/confirmation");
  }

  return (
    <div className="mx-auto max-w-[1100px] px-4 py-10">
      <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
        Checkout
      </h1>

      {isEmpty ? (
        <div className="py-20 text-center">
          <p className="font-jost text-lg font-semibold text-botb-text">
            Your basket is empty
          </p>
          <Link
            href="/competitions"
            className="mt-4 inline-block rounded-md bg-botb-orange px-6 py-3 font-jost text-sm font-bold uppercase text-white transition-colors hover:bg-botb-orange-hover"
          >
            Browse Competitions
          </Link>
        </div>
      ) : (
        <form
          onSubmit={handlePlaceOrder}
          className="mt-8 grid grid-cols-1 gap-8 lg:grid-cols-[1fr_360px]"
        >
          {/* Details + payment */}
          <div className="flex flex-col gap-6">
            <section className="rounded-lg border border-botb-card-border bg-white p-6">
              <h2 className="font-jost text-lg font-semibold text-botb-text">
                Your Details
              </h2>

              <div className="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div>
                  <label className="mb-1 block text-sm text-botb-text">
                    First name
                  </label>
                  <input type="text" required className={INPUT_CLASS} />
                </div>
                <div>
                  <label className="mb-1 block text-sm text-botb-text">
                    Last name
                  </label>
                  <input type="text" required className={INPUT_CLASS} />
                </div>
                <div>
                  <label className="mb-1 block text-sm text-botb-text">
                    Email
                  </label>
                  <input type="email" required className={INPUT_CLASS} />
                </div>
                <div>
                  <label className="mb-1 block text-sm text-botb-text">
                    Phone
                  </label>
                  <input type="tel" required className={INPUT_CLASS} />
                </div>
                <div className="sm:col-span-2">
                  <label className="mb-1 block text-sm text-botb-text">
                    Address line 1
                  </label>
                  <input type="text" required className={INPUT_CLASS} />
                </div>
                <div>
                  <label className="mb-1 block text-sm text-botb-text">
                    City
                  </label>
                  <input type="text" required className={INPUT_CLASS} />
                </div>
                <div>
                  <label className="mb-1 block text-sm text-botb-text">
                    Postcode
                  </label>
                  <input type="text" required className={INPUT_CLASS} />
                </div>
              </div>
            </section>

            <section className="rounded-lg border border-botb-card-border bg-white p-6">
              <h2 className="font-jost text-lg font-semibold text-botb-text">
                Payment
              </h2>
              <p className="mt-2 text-sm text-botb-muted">
                Demo only — no real payment is taken.
              </p>

              <div className="mt-4 space-y-4 opacity-60">
                <div>
                  <label className="mb-1 block text-sm text-botb-text">
                    Card number
                  </label>
                  <input
                    type="text"
                    disabled
                    placeholder="•••• •••• •••• ••••"
                    className={INPUT_CLASS}
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="mb-1 block text-sm text-botb-text">
                      Expiry
                    </label>
                    <input
                      type="text"
                      disabled
                      placeholder="MM / YY"
                      className={INPUT_CLASS}
                    />
                  </div>
                  <div>
                    <label className="mb-1 block text-sm text-botb-text">
                      CVC
                    </label>
                    <input
                      type="text"
                      disabled
                      placeholder="•••"
                      className={INPUT_CLASS}
                    />
                  </div>
                </div>
              </div>
            </section>
          </div>

          {/* Order summary */}
          <aside className="lg:sticky lg:top-6 lg:self-start">
            <div className="rounded-lg border border-botb-card-border bg-white p-6">
              <h2 className="font-jost text-lg font-semibold text-botb-text">
                Order Summary
              </h2>

              <div className="mt-4 space-y-3">
                {items.map((item) => (
                  <div
                    key={item.slug}
                    className="flex justify-between gap-4 text-sm"
                  >
                    <span className="min-w-0 truncate text-botb-text">
                      {item.title}
                      <span className="text-botb-muted"> × {item.qty}</span>
                    </span>
                    <span className="shrink-0 text-botb-text">
                      £{(item.unitPrice * item.qty).toFixed(2)}
                    </span>
                  </div>
                ))}
              </div>

              <div className="mt-4 space-y-3 border-t border-botb-card-border pt-4 text-sm">
                <div className="flex justify-between text-botb-text">
                  <span>Delivery</span>
                  <span className="font-semibold text-botb-orange">FREE</span>
                </div>
                <div className="flex justify-between font-jost text-base font-bold text-botb-text">
                  <span>Total</span>
                  <span>£{total.toFixed(2)}</span>
                </div>
              </div>

              <button
                type="submit"
                className="mt-6 w-full rounded-md bg-botb-orange px-6 py-3 font-jost text-sm font-bold uppercase text-white transition-colors hover:bg-botb-orange-hover"
              >
                Place Order
              </button>

              <Link
                href="/cart"
                className="mt-4 block text-center text-sm text-botb-muted transition-colors hover:text-botb-text"
              >
                Back to basket
              </Link>
            </div>
          </aside>
        </form>
      )}
    </div>
  );
}
