'use client'

import classnames from "classnames";
import Link from '@/src/components/link'
import { usePathname } from 'next/navigation'

export default function Component({ nodes, className }) {
  var path = usePathname()
  return (
    <ul className={className}>
      {nodes.map( node => { 
        return <Entry key={node.name} node={node} path={path} indent={0}/>
      })}
    </ul>
  )
}

function Entry({ node, path, indent }) {
  var active = node.href !== "/" && path.startsWith(node.href);
	if (node.href === path) {
		active = true
	}
	// var active = true;

	var margin = "ml-4 mb-1";

  if (node.children?.length > 0) {
    return (
      <li className={margin}>
        <details open={active}>
          <summary>  
            <Link
							href={node.href}
							className={classnames(
								{"dark:text-yellow-300" : active},
								{"font-bold" : active},
								"no-underline hover:underline",
								"decoration-yellow-400 dark:decoration-yellow-300",
								"decoration-2",
							)}
						>{node.name}</Link>
          </summary>
          <ul>
            {node.children?.map( node => { 
              return <Entry key={node.name} node={node} path={path} indent={indent + 2}/>
            })}
          </ul>
        </details>
      </li>
    )
  } else {
    return (
      <li className={margin}>
        <Link
					href={node.href}
					className={classnames(
						{"dark:text-yellow-300" : active},
						{"font-bold" : active},
						"pl-4",
						"no-underline hover:underline",
						"decoration-yellow-400 dark:decoration-yellow-300",
						"decoration-2",
					)}
				>{node.name}</Link>
      </li>
    )
  }
}
