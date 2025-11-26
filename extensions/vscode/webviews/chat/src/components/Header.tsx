import { useState } from "react"
import {
  Braces,
} from 'lucide-react'

import { cn } from "@/lib/utils"


import { JsonInfo, UsageInfo } from './Info';

export const Header = ({
  sid,
  usage,
  session,
  chatState,
  className,
}:{
  sid: string,
  usage: any,
  session: any,
  chatState: any,
  className?: string,
}) => {
  const [hidden, setHidden] = useState(true);

  return (

    <div className={cn("flex flex-col p-2", className)}>
      <div className="text-md flex gap-2 justify-between border-b px-4 py-2">
        <span>{session?.state?.title || sid}</span>
        <div className="ml-auto flex justify-end items-center gap-2">
          <UsageInfo usage={usage} size={16}/>
          <Braces size={16} onClick={(e) => {
            e.preventDefault()
            e.stopPropagation()
            setHidden(!hidden)
          }} />
        </div>
      </div>
      { !hidden && <JsonInfo data={{
        sid,
        usage,
        session,
        chatState,
      }} />}
    </div>

  )
}

