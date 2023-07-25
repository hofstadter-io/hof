'use client'

import dynamic from 'next/dynamic'

const MDXEditor = dynamic(
  () => import('@mdxeditor/editor').then((mod) => mod.MDXEditor), 
  { ssr: false }
)

export default function Page() {
  return (
		<div className="p-8">
			<MDXEditor
				markdown="Hallo *world!!*"
				onChange={(markdown) => console.log(markdown)}
				contentEditableClassName="editor-prose"
			/>
		</div>
  );
}

