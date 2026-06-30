"use client";

import { useState } from "react";
import Link from "next/link";
import { Logo, CartIcon, MenuIcon, UserIcon, CloseIcon } from "@/components/icons";
import { cn } from "@/lib/utils";
import { useCart } from "@/lib/cart";

const NAV_LINKS = [
  { label: "COMPETITIONS", href: "/competitions" },
  { label: "WINNERS", href: "/winners" },
] as const;

const MOBILE_LINKS = [
  { label: "COMPETITIONS", href: "/competitions" },
  { label: "WINNERS", href: "/winners" },
  { label: "HOW TO PLAY", href: "/how-to-play" },
  { label: "BOTB PASS", href: "/botb-pass" },
  { label: "LOG IN", href: "/login" },
  { label: "SIGN UP", href: "/register" },
] as const;

/** Small cart icon with a live count badge, reused across desktop and mobile bars. */
function CartButton() {
  const { count } = useCart();
  return (
    <Link href="/cart" className="relative inline-flex" aria-label="Cart">
      <CartIcon className="h-7 w-7 text-botb-secondary" />
      <span className="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-600 text-[10px] text-white">
        {count}
      </span>
    </Link>
  );
}

/** Stacked promo + navigation header. */
export function SiteHeader() {
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <header className="sticky top-0 z-[1003]">
      {/* Top promo bar */}
      <div className="flex h-10 w-full items-center justify-center gap-2 bg-black px-4">
        <span className="rounded-full bg-red-600 px-2 py-0.5 text-[11px] font-bold uppercase text-white">
          NEW
        </span>
        <span className="text-[13px] text-white">
          🎁 FINAL DAYS: New Subscribers get 20 EXTRA Home Tickets!
        </span>
      </div>

      {/* Nav bar */}
      <div className="w-full border-b border-botb-card-border bg-white">
        {/* Desktop */}
        <div className="relative mx-auto hidden h-[60px] max-w-[1360px] items-center px-4 lg:flex">
          {/* Left zone */}
          <nav className="flex items-center gap-8">
            {NAV_LINKS.map((link) => (
              <Link
                key={link.label}
                href={link.href}
                className="font-roboto text-[16px] font-medium uppercase text-botb-secondary transition-colors hover:text-botb-orange"
              >
                {link.label}
              </Link>
            ))}
          </nav>

          {/* Center zone */}
          <Link
            href="/"
            className="absolute left-1/2 top-1/2 flex -translate-x-1/2 -translate-y-1/2 flex-col items-center"
          >
            <Logo className="w-[109px] h-auto" />
          </Link>

          {/* Right zone */}
          <div className="ml-auto flex items-center gap-5">
            <Link
              href="/login"
              className="text-[15px] font-medium uppercase text-botb-secondary transition-colors hover:text-botb-orange"
            >
              LOG IN
            </Link>
            <Link
              href="/register"
              className="rounded border-2 border-botb-orange bg-botb-orange px-5 py-1.5 font-jost text-[14px] font-medium uppercase text-white transition-colors hover:bg-botb-orange-hover"
            >
              SIGN UP
            </Link>
            <CartButton />
          </div>
        </div>

        {/* Mobile */}
        <div className="relative flex h-[60px] items-center px-4 lg:hidden">
          {/* Left: hamburger */}
          <button
            type="button"
            aria-label={menuOpen ? "Close menu" : "Open menu"}
            aria-expanded={menuOpen}
            onClick={() => setMenuOpen((open) => !open)}
            className="inline-flex"
          >
            {menuOpen ? (
              <CloseIcon className="h-6 w-6 text-botb-text" />
            ) : (
              <MenuIcon className="h-6 w-6 text-botb-text" />
            )}
          </button>

          {/* Center: logo */}
          <Link href="/" className="absolute left-1/2 flex -translate-x-1/2 flex-col items-center">
            <Logo className="w-[80px]" />
          </Link>

          {/* Right: account + cart */}
          <div className="ml-auto flex items-center gap-4">
            <Link href="/account" aria-label="Account" className="inline-flex">
              <UserIcon className="h-6 w-6 text-botb-secondary" />
            </Link>
            <CartButton />
          </div>
        </div>

        {/* Mobile menu */}
        {menuOpen && (
          <nav className="border-t border-botb-card-border bg-white px-4 lg:hidden">
            {MOBILE_LINKS.map((link) => (
              <Link
                key={link.label}
                href={link.href}
                onClick={() => setMenuOpen(false)}
                className={cn(
                  "block border-b border-botb-card-border py-3 text-[15px] font-medium uppercase text-botb-secondary transition-colors hover:text-botb-orange",
                )}
              >
                {link.label}
              </Link>
            ))}
          </nav>
        )}
      </div>
    </header>
  );
}
