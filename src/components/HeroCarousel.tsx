"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { heroSlides, heroStats, asSeenOn } from "@/lib/data";
import { ChevronLeftIcon, ChevronRightIcon, CartIcon } from "@/components/icons";
import { cn } from "@/lib/utils";

/** Auto-advance interval for the hero carousel, in milliseconds. */
const AUTO_ADVANCE_MS = 5500;

export function HeroCarousel() {
  const [index, setIndex] = useState(0);
  const slideCount = heroSlides.length;

  // Hold the active interval id so manual navigation can reset the timer,
  // preventing a slide from flipping immediately after a user click.
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const clearTimer = useCallback(() => {
    if (timerRef.current !== null) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }
  }, []);

  const startTimer = useCallback(() => {
    clearTimer();
    timerRef.current = setInterval(() => {
      setIndex((current) => (current + 1) % slideCount);
    }, AUTO_ADVANCE_MS);
  }, [clearTimer, slideCount]);

  // Drive the auto-advance loop and tear it down on unmount.
  useEffect(() => {
    startTimer();
    return clearTimer;
  }, [startTimer, clearTimer]);

  // Manual navigation restarts the timer so the next auto-advance is a full cycle away.
  const goTo = useCallback(
    (next: number) => {
      setIndex(((next % slideCount) + slideCount) % slideCount);
      startTimer();
    },
    [slideCount, startTimer],
  );

  const goPrev = useCallback(() => goTo(index - 1), [goTo, index]);
  const goNext = useCallback(() => goTo(index + 1), [goTo, index]);

  const activeSlide = heroSlides[index];

  return (
    <section className="relative bg-black">
      {/* Slide area */}
      <div className="relative w-full overflow-hidden">
        <div className="relative h-[420px] w-full md:h-[520px]">
          {/* Full-bleed slide image */}
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            key={activeSlide.image}
            src={activeSlide.image}
            alt={activeSlide.title}
            className="absolute inset-0 h-full w-full object-cover object-center"
          />

          {/* Legibility gradient (stronger on desktop where text sits left) */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/10 to-transparent md:bg-gradient-to-r md:from-black/70 md:via-black/20 md:to-transparent" />

          {/* Overlay text */}
          <div className="absolute inset-0">
            <div className="mx-auto flex h-full max-w-[1360px] flex-col justify-end px-4 pb-16 text-center md:justify-center md:pb-0 md:text-left">
              <div className="max-w-xl md:max-w-2xl">
                <span className="inline-block rounded border border-white/60 px-3 py-1 font-jost text-[12px] uppercase tracking-wide text-white">
                  {activeSlide.badge}
                </span>
                <h2 className="mt-3 font-jost text-[36px] font-bold uppercase leading-tight text-white md:text-[44px]">
                  {activeSlide.title}
                </h2>
                <p className="mt-1 text-[16px] text-white/90">{activeSlide.subtitle}</p>

                <div className="mt-5 flex items-center justify-center gap-2 md:justify-start">
                  <button
                    type="button"
                    className="flex-1 rounded bg-botb-orange px-8 py-2.5 font-jost text-[16px] uppercase text-white transition-colors hover:bg-botb-orange-hover md:flex-none"
                  >
                    Enter to win &raquo;
                  </button>
                  <button
                    type="button"
                    aria-label="Add to basket"
                    className="rounded bg-botb-orange p-2.5 transition-colors hover:bg-botb-orange-hover"
                  >
                    <CartIcon className="h-5 w-5 text-white" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          {/* Prev / next arrows (desktop only) */}
          <button
            type="button"
            onClick={goPrev}
            aria-label="Previous slide"
            className="absolute left-3 top-1/2 hidden -translate-y-1/2 text-white/80 transition-colors hover:text-white md:block"
          >
            <ChevronLeftIcon className="h-8 w-8" />
          </button>
          <button
            type="button"
            onClick={goNext}
            aria-label="Next slide"
            className="absolute right-3 top-1/2 hidden -translate-y-1/2 text-white/80 transition-colors hover:text-white md:block"
          >
            <ChevronRightIcon className="h-8 w-8" />
          </button>

          {/* Dots */}
          <div className="absolute bottom-4 left-1/2 flex -translate-x-1/2 items-center gap-1.5">
            {heroSlides.map((slide, dotIndex) => (
              <button
                key={slide.image}
                type="button"
                onClick={() => goTo(dotIndex)}
                aria-label={`Go to slide ${dotIndex + 1}`}
                aria-current={dotIndex === index}
                className={cn(
                  "h-2 rounded-full transition-all",
                  dotIndex === index ? "w-4 bg-white" : "w-2 bg-white/50",
                )}
              />
            ))}
          </div>
        </div>
      </div>

      {/* Stats band */}
      <div className="mx-auto max-w-[1360px] px-4 py-3">
        <div className="flex items-stretch justify-center gap-3 md:justify-start md:gap-6">
          {heroStats.map((stat, statIndex) => (
            <div
              key={stat.label}
              className={cn(
                "flex items-center gap-2 pl-3 md:pl-6",
                statIndex > 0 && "border-l border-white/20",
              )}
            >
              {stat.est ? (
                <span className="flex h-9 w-9 flex-col items-center justify-center rounded-full border border-white/40 text-center text-[8px] leading-tight text-white/70">
                  <span>Est</span>
                  <span>99</span>
                </span>
              ) : (
                <span className="flex h-9 w-9 items-center justify-center rounded-full border border-white/40 text-white/70" aria-hidden="true">
                  {stat.label === "Winners" ? (
                    <svg viewBox="0 0 24 24" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="1.6">
                      <path d="M8 4h8v4a4 4 0 0 1-8 0V4z" strokeLinecap="round" strokeLinejoin="round" />
                      <path d="M8 5H5v1a3 3 0 0 0 3 3M16 5h3v1a3 3 0 0 1-3 3M10 14h4M9 20h6M12 14v6" strokeLinecap="round" strokeLinejoin="round" />
                    </svg>
                  ) : (
                    <svg viewBox="0 0 24 24" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="1.6">
                      <path d="M4 11h16v9H4z M12 7v13 M4 7h16v4H4z M12 7C9 7 7 4 9 3s3 4 3 4 1-3 3-3-1 3-3 3" strokeLinecap="round" strokeLinejoin="round" />
                    </svg>
                  )}
                </span>
              )}
              <span className="flex flex-col">
                <span className="font-jost text-[18px] font-bold text-white">{stat.value}</span>
                <span className="text-[12px] text-white/60">{stat.label}</span>
              </span>
            </div>
          ))}
        </div>
      </div>

      {/* As seen on band */}
      <div className="mx-auto max-w-[1360px] px-4 py-3">
        <div className="flex flex-wrap items-center justify-center gap-x-6 gap-y-3 md:justify-between">
          <span className="text-[13px] font-medium uppercase tracking-wide text-white/70">
            As seen on
          </span>
          {asSeenOn.map((press) => (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              key={press.name}
              src={press.image}
              alt={press.name}
              className="h-7 w-auto object-contain opacity-90 md:h-8"
            />
          ))}
        </div>
      </div>
    </section>
  );
}
