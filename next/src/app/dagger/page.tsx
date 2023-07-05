"use client";

import Content from "./content.mdx";

export default function Page(props) {
  return (
		<div className="flex w-full">
			<div className="w-64 bg-slate-600 text-slate-50 p-8">tree</div>
			<div className="p-8">
				<Content {...props} />
			</div>
		</div>
	)
}

