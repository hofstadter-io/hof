"use client";

import React from "react";

import { ThemeProvider, useTheme } from "@wits/next-themes";

import { MDXProvider } from "@mdx-js/react";
import "@code-hike/mdx/dist/index.css";

export function Providers({ children }) {
  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <MDXProvider>{children}</MDXProvider>
    </ThemeProvider>
  );
}

export function ThemeChanger() {
  const { resolvedTheme, setTheme } = useTheme();

  const handle = function () {
    if (resolvedTheme === "light") {
      setTheme("dark");
    } else {
      setTheme("light");
    }
  };

  return (
    <span aria-label="Toggle dark mode" onClick={() => handle()}>
      Toggle dark mode
    </span>
  );
}
