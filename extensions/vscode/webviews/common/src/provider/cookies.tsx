"use client"

import * as React from "react"
import { CookiesProvider } from 'react-cookie';
 
export function CookieProvider({
  children,
  ...props
}: React.ComponentProps<typeof CookiesProvider>) {
  return <CookiesProvider {...props}>{children}</CookiesProvider>
}