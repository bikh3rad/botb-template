import type { Metadata } from "next"

import { MediaLibrary } from "@/components/admin/MediaLibrary"

export const metadata: Metadata = {
  title: "Media Library | Admin Console",
  description: "Browse and manage every uploaded media file across the platform.",
}

export default function MediaPage() {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="font-heading text-2xl font-semibold tracking-tight">
          Media Library
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Browse and manage every uploaded media file across the platform.
        </p>
      </div>

      <MediaLibrary />
    </div>
  )
}
