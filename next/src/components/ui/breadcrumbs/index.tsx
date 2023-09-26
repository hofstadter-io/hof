'use client'

import classnames from 'classnames';
import Link from '@/src/components/link'
import { usePathname } from 'next/navigation'

export default function Component() {	
  var path = usePathname()
	var parts = path.split("/")

	return (
		<span className="text-xl m-4">
			<Link href="/"
				className={classnames(
					"normal-case m-1",
					// "transition-all",
					"dark:hover:text-yellow-300 decoration-yellow-300",
					"no-underline hover:underline",
					"underline-offset-2 hover:underline-offset-8",
					"decoration-0 hover:decoration-2",
				)}
			>Home</Link>
			{ parts.map((part, idx) => {
				var ps = parts.slice(0, idx+1)
				var hr = ps.join("/")
				// ignore empty first element and home page
				if (idx > 0 && path !== "/") {
					return (
						<span key={idx}> / <Link
							href={hr}
							className={classnames(
								"normal-case m-1",
								// "transition-all",
								"dark:hover:text-yellow-300 decoration-yellow-300",
								"no-underline hover:underline",
								"underline-offset-2 hover:underline-offset-8",
								"decoration-0 hover:decoration-2",
							)}
						>{part}</Link></span>
					)
				}
			})}
		</span>
	)

}
