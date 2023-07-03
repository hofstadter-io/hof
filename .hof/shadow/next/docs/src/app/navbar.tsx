"use client";

import React from "react";
import Link from "@/src/components/link";

import { ThemeChanger } from "./providers";

const content = {
  brand: "Supacode",
  href: "/",

  menu: [
    {
      title: "Docs",
      href: "/docs",
    },
    {
      title: "Essays",
      href: "/essays",
    },
    {
      title: "Posts",
      href: "/posts",
    },
    {
      title: "User",
      href: "/user",
    },
  ],
};

export default function Navbar() {
  return (
    <nav className="navbar z-40 sticky top-0 drop-shadow-xl">
      <div className="flex-1">
        <Link href={content.href} className="btn btn-ghost normal-case text-xl">
          {content.brand}
        </Link>
      </div>
      <div className="flex-1">
        {content.menu.map((item) => (
          <Link
            key={item.title}
            href={item.href}
            className="btn btn-ghost normal-case text-lg"
          >
            {item.title}
          </Link>
        ))}
      </div>
      <div className="flex-none gap-2"></div>
    </nav>
  );
}
