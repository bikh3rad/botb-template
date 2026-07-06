import type { Metadata } from "next"

import { CategoriesManager } from "@/components/admin/CategoriesManager"

export const metadata: Metadata = {
  title: "Categories | Admin Console",
  description: "Manage competition categories.",
}

export default function CategoriesPage() {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="font-heading text-2xl font-semibold tracking-tight">
          Categories
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Manage competition categories.
        </p>
      </div>

      <CategoriesManager />
    </div>
  )
}
