'use client'

import Link from '@/src/components/link'
import { usePathname } from 'next/navigation'

export default function Component({ nodes, className }) {
  var path = usePathname()
  return (
    <ul className={className}>
      {nodes.map( node => { 
        return <Entry key={node.name} node={node} path={path}/>
      })}
    </ul>
  )
}

function Entry({ node, path }) {
  var active = (node.href === path);

  if (node.children?.length === 0) {
    return (
      <li>
        <Link href={node.href} className={active ? "active" : ""}>{node.name}</Link>
      </li>
    )
  } else {
    return (
      <li>
        <details open>
          <summary className={active ? "active" : ""}>  
            <Link href={node.href} >{node.name}</Link>
          </summary>
          <ul>
            {node.children?.map( node => { 
              return <Entry key={node.name} node={node} path={path}/>
            })}
          </ul>
        </details>
      </li>
    )
  }
}
