"use client";

import React from "react";

import { ThemeProvider } from "@wits/next-themes";

import { MDXProvider } from "@mdx-js/react";
import "@code-hike/mdx/dist/index.css";

export function Providers({ children }) {
  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <MDXProvider>{children}</MDXProvider>
    </ThemeProvider>
  );
}
