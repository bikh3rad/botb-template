"use client";

import { useEffect, useState } from "react";

import { CloseIcon } from "@/components/icons";
import { cn } from "@/lib/utils";

/**
 * Fixed remaining duration used to seed the countdown. Kept as a constant so the
 * initial server render and the first client render agree (avoids hydration
 * mismatch). The live ticking only begins inside the effect on the client.
 */
const INITIAL_TIME = { h: 15, m: 33, s: 33 } as const;

type CountdownTime = {
  h: number;
  m: number;
  s: number;
};

/** Pad a number to a fixed 2-digit string for stable display. */
function pad(value: number): string {
  return value.toString().padStart(2, "0");
}

/**
 * Decrement the countdown by one second, cascading seconds -> minutes -> hours
 * and clamping at zero so the timer never goes negative.
 */
function tick(time: CountdownTime): CountdownTime {
  if (time.h === 0 && time.m === 0 && time.s === 0) {
    return time;
  }

  let { h, m, s } = time;

  if (s > 0) {
    s -= 1;
  } else {
    s = 59;
    if (m > 0) {
      m -= 1;
    } else {
      m = 59;
      h -= 1;
    }
  }

  return { h, m, s };
}

type CountdownGroupProps = {
  value: string;
  label: string;
};

/** A single big-number + tiny-label unit within the countdown. */
function CountdownGroup({ value, label }: CountdownGroupProps) {
  return (
    <div className="flex flex-col items-center">
      <span className="font-jost text-[20px] font-bold leading-none md:text-[24px]">
        {value}
      </span>
      <span className="text-[9px] uppercase tracking-wide opacity-80">
        {label}
      </span>
    </div>
  );
}

export function PromoRibbon() {
  const [dismissed, setDismissed] = useState(false);
  const [time, setTime] = useState<CountdownTime>(INITIAL_TIME);

  useEffect(() => {
    // Start ticking only on the client to keep SSR markup deterministic.
    const intervalId = setInterval(() => {
      setTime((current) => tick(current));
    }, 1000);

    return () => clearInterval(intervalId);
  }, []);

  if (dismissed) {
    return null;
  }

  return (
    <div className="botb-ribbon-gradient fixed bottom-0 left-0 right-0 z-[200000] text-white">
      <div className="mx-auto flex max-w-[1360px] items-center justify-between gap-3 px-4 py-2.5">
        {/* LEFT: headline message */}
        <p className="flex min-w-0 items-center gap-2 truncate font-jost text-[15px] font-semibold md:text-[17px]">
          <span aria-hidden="true">🏡</span>
          <span className="truncate">
            LAST HOURS: Win a £1.2M Home in Zone 1!
          </span>
        </p>

        {/* CENTER: price pill + live countdown */}
        <div className="flex shrink-0 items-center gap-3">
          <span className="hidden rounded-full bg-black px-4 py-1 text-[14px] font-semibold text-white sm:inline-block">
            Only £1
          </span>

          <div className="flex items-end gap-2">
            <CountdownGroup value={pad(time.h)} label="Hours" />
            <span className="font-jost text-[20px] font-bold leading-none md:text-[24px]">
              :
            </span>
            <CountdownGroup value={pad(time.m)} label="Minutes" />
            <span className="font-jost text-[20px] font-bold leading-none md:text-[24px]">
              :
            </span>
            <CountdownGroup value={pad(time.s)} label="Seconds" />
          </div>
        </div>

        {/* RIGHT: CTA + dismiss */}
        <div className="flex shrink-0 items-center gap-2">
          <button
            type="button"
            className={cn(
              "rounded bg-white px-6 py-2 font-jost text-[14px] font-semibold uppercase text-botb-orange",
              "hover:bg-white/90",
            )}
          >
            Enter Now
          </button>
          <button
            type="button"
            aria-label="Dismiss promotion"
            onClick={() => setDismissed(true)}
            className="flex items-center justify-center p-1 text-white hover:opacity-80"
          >
            <CloseIcon className="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  );
}
