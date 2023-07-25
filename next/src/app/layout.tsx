import "./globals.css";

import classnames from "classnames";
import { Inter } from "next/font/google";
const inter = Inter({ subsets: ["latin"] });

import { Providers } from "./providers";
import Navbar from "@/src/components/ui/navbars/default";
import Footer from "@/src/components/ui/footers/default";

import ThemeChanger from "@/src/components/ui/themes/changer";

import PrefMenu from "@/src/components/ui/navbars/pref-menu";

const navbar = {
  brand: {
    name: "Documentation",
    href: "/",
  },

  links: [
    {
      title: "Dagger",
      href: "/dagger",
    },
    {
      title: "Editor",
      href: "/editor",
    },
    {
      title: "About",
      href: "/about",
    },
  ],

  dropdown: <PrefMenu ThemeChanger={ThemeChanger} />,
};

export const metadata = {
  title: "Hof Docs",
  description: "Hofstadter Documentation Site",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={classnames(
          inter.className,
          "min-h-screen antialiased",
          "transition-all"
        )}
      >
        <Providers>
          <div className="flex flex-col min-h-screen">
            <Navbar {...navbar} />

						<div className="flex grow content-stretch min-w-screen">{children}</div>

            <Footer copyright="Hofstadter, Inc" />
          </div>
        </Providers>
      </body>
    </html>
  );
}
