import { randomUUID } from "node:crypto"
import { mkdir, writeFile } from "node:fs/promises"
import path from "node:path"

import { NextResponse, type NextRequest } from "next/server"

import { isAuthorizedAdmin } from "@/lib/admin-auth"

const UPLOAD_DIR = path.join(process.cwd(), "public", "uploads")
const MAX_BYTES = 25 * 1024 * 1024 // 25 MB

// Allowed MIME types mapped to the extension we persist them under.
const ALLOWED_TYPES = new Map<string, string>([
  ["image/jpeg", "jpg"],
  ["image/png", "png"],
  ["image/webp", "webp"],
  ["image/gif", "gif"],
  ["image/avif", "avif"],
  ["video/mp4", "mp4"],
  ["video/webm", "webm"],
])

/**
 * POST /admin/api/media
 *
 * Guarded media upload. Accepts a `multipart/form-data` body with a single
 * `file` field and stores it under `public/uploads/` with a generated name.
 * Requires a valid admin token — see `isAuthorizedAdmin`.
 */
export async function POST(request: NextRequest) {
  if (!isAuthorizedAdmin(request)) {
    return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
  }

  let form: FormData
  try {
    form = await request.formData()
  } catch {
    return NextResponse.json(
      { error: "Expected multipart/form-data body" },
      { status: 400 },
    )
  }

  const file = form.get("file")
  if (!(file instanceof File)) {
    return NextResponse.json(
      { error: "Missing 'file' field" },
      { status: 400 },
    )
  }

  const ext = ALLOWED_TYPES.get(file.type)
  if (!ext) {
    return NextResponse.json(
      { error: `Unsupported media type: ${file.type || "unknown"}` },
      { status: 415 },
    )
  }

  if (file.size === 0) {
    return NextResponse.json({ error: "File is empty" }, { status: 400 })
  }
  if (file.size > MAX_BYTES) {
    return NextResponse.json(
      { error: "File exceeds the 25 MB limit" },
      { status: 413 },
    )
  }

  // Generate the stored name ourselves — never trust the client filename,
  // which could contain path-traversal segments.
  const filename = `${randomUUID()}.${ext}`
  const bytes = Buffer.from(await file.arrayBuffer())

  await mkdir(UPLOAD_DIR, { recursive: true })
  await writeFile(path.join(UPLOAD_DIR, filename), bytes)

  return NextResponse.json(
    { url: `/uploads/${filename}`, filename, size: file.size, type: file.type },
    { status: 201 },
  )
}
