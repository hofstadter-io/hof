import { useState } from "react"

import { cn } from "@/lib/utils"

import { Menu, JsonInfo, UsageInfo } from './Info';

export const Header = ({
  sid,
  setPos,
  usage,
  session,
  chatState,
  className,
}:{
  sid: string,
  setPos: any,
  usage: any,
  session: any,
  chatState: any,
  className?: string,
}) => {
  const [hidden, setHidden] = useState(true);

  return (

    <div className={cn("flex flex-col m-2 border-b", className)}>

      <div className="flex justify-between items-center gap-2 p-2">
        <span>{session?.state?.title || sid}</span>
        <Menu sid={sid} setPos={() => setPos(-1)} hidden={hidden} setHidden={setHidden} />
      </div>

      <div className="flex justify-between items-center gap-2 p-2">
        <UsageInfo usage={usage} size={16}/>
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

