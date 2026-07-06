import type { Metadata } from "next"

import { SiteTextsEditor } from "@/components/admin/SiteTextsEditor"

export const metadata: Metadata = {
  title: "Site Texts | Admin Console",
  description: "Edit editable copy and content strings shown on the site.",
}

export default function SiteTextsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="font-heading text-2xl font-semibold tracking-tight">
          Site Texts
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Edit editable copy and content strings shown on the site.
        </p>
      </div>

      <SiteTextsEditor />
    </div>
  )
}
