import { useState } from "react"
import { Plus, Trash2 } from "lucide-react";

import { vscodeApi } from "@/vscodeApi";
import { cn } from "@/lib/utils"

import { JsonInfo, UsageInfo, DiffInfo } from '@/components/Info';
import { Menu } from '@/components/SessionMenu'
import { Tooltipped } from "@/components/Tooltipped";

export const Header = ({
  sid,
  setPos,
  diff,
  usage,
  session,
  chatState,
  className,
}:{
  sid: string,
  setPos: any,
  diff?: any,
  usage?: any,
  session?: any,
  chatState?: any,
  className?: string,
}) => {
  const [hidden, setHidden] = useState(true);

  return (

    <div className={cn("flex flex-col m-2 border-b", className)}>

      <div className="flex justify-between items-center gap-2 p-2">
        <span>{session?.state?.title || sid || "no session"}</span>
        <Menu sid={sid} setPos={() => setPos(-1)} hidden={hidden} setHidden={setHidden} refresh/>

        <Tooltipped label="Create">
          <Plus size={16}
            aria-label="create"
            className="hover:text-green-500"
            onClick={() => {
              const state = vscodeApi.getState()
              console.log("Create!", state)
              vscodeApi.postMessage({
                type: 'session.create',
                payload: { focus: true, dir: chatState?.env?.workspaceDir },
              });
            }}
          />
        </Tooltipped>
        <Tooltipped label="Delete">
          <Trash2 size={16}
            aria-label="delete"
            className="hover:text-red-500"
            onClick={() => {
              const state = vscodeApi.getState()
              console.log("Delete!", state)
              vscodeApi.postMessage({
                type: 'session.delete',
                payload: { sid },
              });
            }}
          />
        </Tooltipped>
      </div>

      <div className="flex justify-between items-center gap-2 p-2">
        <UsageInfo usage={usage} size={16}/>
        <DiffInfo diff={diff} size={16}/>
      </div>

      { !hidden && <JsonInfo data={{
        sid,
        usage,
        session,
        chatState,
        diff,
      }} />}
    </div>

  )
}