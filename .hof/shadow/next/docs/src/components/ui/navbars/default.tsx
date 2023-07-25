"use client";

import classnames from 'classnames';
import React from "react";
import Link from "@/src/components/link";

export default function Navbar({
	brand,
	links,
	search,
	dropdown,
}:{
	brand: any,
	links: any,
	search: React.Node,
	dropdown: React.Node,
}) {

  return (
    <nav className="navbar z-40 sticky top-0 drop-shadow-xl gap-x-16">
      <div className="flex-none">
        <Link href={brand.href} className="btn btn-ghost normal-case text-xl">
          {brand.name}
        </Link>
      </div>
      <div className="flex-1 justify-content-start">
        { links.map( item => (
          <Link key={ item.title } href={ item.href }
						className={classnames(
							"btn btn-ghost",
							"normal-case text-lg",
							"transition-all",
							"hover:text-yellow-300 decoration-yellow-300",
							"no-underline hover:underline",
							"underline-offset-2 hover:underline-offset-8",
							"decoration-0 hover:decoration-2",
						)}
					>
            { item.title }
          </Link>
        ))}
      </div>
      <div className="flex-none gap-2">
				{ search }
				{ dropdown }
      </div>
    </nav>
  );
}