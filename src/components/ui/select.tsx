import * as React from "react"
import { ChevronDown } from "lucide-react"

import { cn } from "@/lib/utils"

/**
 * Lightweight styled wrapper over the native <select> element. Keeps behaviour
 * accessible and dependency-free while matching the design system. Swap for a
 * richer popup control later without changing call sites.
 */
function Select({ className, children, ...props }: React.ComponentProps<"select">) {
  return (
    <div data-slot="select" className="relative inline-flex w-full">
      <select
        className={cn(
          "h-9 w-full appearance-none rounded-lg border border-input bg-background py-1 pr-9 pl-3 text-sm shadow-xs transition-[color,box-shadow] outline-none",
          "focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50",
          "disabled:pointer-events-none disabled:opacity-50",
          "dark:bg-input/30",
          className
        )}
        {...props}
      >
        {children}
      </select>
      <ChevronDown className="pointer-events-none absolute top-1/2 right-3 size-4 -translate-y-1/2 text-muted-foreground" />
    </div>
  )
}

export { Select }
