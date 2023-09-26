import "./globals.css";

import classnames from "classnames";
import { Inter } from "next/font/google";
const inter = Inter({ subsets: ["latin"] });

import { Providers } from "./providers";
import Navbar from "@/src/components/ui/navbars/default";
import ThemeChanger from "@/src/components/ui/themes/changer";
import Footer from "@/src/components/ui/footers/default";
import TOC from "@/src/components/ui/sidenavs/toc";
import Breadcrumbs from "@/src/components/ui/breadcrumbs";

import siteTOC from "./menu.json";

const navbar = {
  brand: {
    name: "Documentation",
    href: "/",
  },

  links: [{
		title: "v0.6.9-beta.1",
		href: "https://github.com/hofstadter-io/hof/releases/tag/v0.6.9-beta.1",
	},{
		title: "GitHub",
		href: "https://github.com/hofstadter-io/hof",
	}, {
		title: "Chat",
		href: "#",
		items: [{
			title: "Discord",
			href: "https://discord.gg/BXwX7n6B8w",
		},{
			title: "Slack",
			href: "https://join.slack.com/t/hofstadter-io/shared_invite/zt-e5f90lmq-u695eJur0zE~AG~njNlT1A",
		}]
	}, {
		title: "hof.io",
		href: "https://hofstadter.io",
	}],

	themer: <ThemeChanger />,
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

						<div className="flex grow content-stretch min-w-screen">
							<TOC nodes={siteTOC} className="pt-4 pr-4 border-r border-slate-300 dark:border-slate-500 shadow"/>
							<div
								className="flex flex-col"
							>
								<Breadcrumbs />
								<div className="p-8 text-lg prose dark:prose-invert">{children}</div>
							</div>
						</div>

            <Footer copyright="Hofstadter, Inc" />
          </div>
        </Providers>
      </body>
    </html>
  );
}
