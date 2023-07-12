"use client";

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
    <nav className="navbar z-40 sticky top-0 drop-shadow-xl">
      <div className="flex-1">
        <Link href={brand.href} className="btn btn-ghost normal-case text-xl">
          {brand.name}
        </Link>
      </div>
      <div className="flex-1">
        { links.map( item => (
          <Link key={ item.title } href={ item.href } className="btn btn-ghost normal-case text-lg">
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
