"use client";

import Link from "@/src/components/link";

export default function NotFound() {
  return (
    <div className="p-8">
      <h1>Not Found</h1>
      <Link href="/">Go Home</Link>
    </div>
  );
}
