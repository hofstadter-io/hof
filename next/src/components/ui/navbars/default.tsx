"use client";

import classnames from 'classnames';
import React from "react";
import Link from "@/src/components/link";

export default function Navbar({
	brand,
	links,
	search,
	themer,
}:{
	brand: any,
	links: any,
	search: React.Node,
	themer: React.Node,
}) {
  return (
    <nav className="navbar z-40 sticky top-0 drop-shadow-xl gap-x-16">
      <div className="navbar-start">
        <Link
					href={brand.href}
					className={classnames(
						"group",
						"btn btn-ghost",
						"normal-case text-xl",
						"transition-all",
						"hover:text-yellow-300 decoration-yellow-300",
						"no-underline hover:underline",
						"underline-offset-2 hover:underline-offset-8",
						"decoration-0 hover:decoration-2",
					)}
				>
          <i className="text-yellow-300 p-0 m-0 -mr-2 group-hover:text-slate-800">_</i>{brand.name}
        </Link>
      </div>
      <div className="navbar-end hidden md:flex mr-2">
        { links.map( (item,idx) => {
					if ( item.items !== undefined ) {
						return (
							<details key={idx} className="dropdown">
							<summary
								className={classnames(
									"btn btn-ghost",
									"normal-case text-lg",
									"transition-all",
									"border-0 hover:border-0",
									"hover:text-yellow-300 decoration-yellow-300",
									"no-underline hover:underline",
									"underline-offset-2 hover:underline-offset-8",
									"decoration-0 hover:decoration-2",
								)}
							>{item.title}</summary>
							<ul className="shadow p-2 dropdown-content z-[1] rounded-box border-0">
								{ item.items.map( it => {
									return (<li key={it.title}><NavLink item={it} /></li>)
							 })}
							</ul>
							</details>
						)
					} else {
						return <NavLink key={idx} item={item} />
					}
					})	
        }
				{ themer }
      </div>
			<div className="dropdown md:hidden navbar-end">
				<label tabIndex={0} className="btn btn-ghost float-right">
					<svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h8m-8 6h16" /></svg>
				</label>
				<DropdownLinks links={links} themer={themer} />
			</div>
    </nav>
  );
}

function DropdownLinks({
	links,
	themer,
}:{
	links: any,
	themer: React.Node,
}) {
	return (
		<ul tabIndex={0} className="dropdown-content mt-3 z-[1] p-2 shadow rounded-box w-52">
			{ links.map( (item,idx) => {
				if ( item.items !== undefined ) {
					return (
						<li key={idx}>
							<NavLink item={item} />
							<ul className="pl-4">
								{ item.items.map( it => {
									return (<li key={it.href}><NavLink item={it} /></li>)
							 })}
							</ul>
						</li>
					)
				} else {
					return (<li><NavLink key={idx} item={item} /></li>)
				}
			})}
			<li key="theme-changer" className="pl-4 p-2 text-lg flex items-center align-center"><b className="mr-4">theme:</b> { themer }</li>
		</ul>
	)
}

function NavLink({ item }: { item: any }) {
	return (
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
	)
}
