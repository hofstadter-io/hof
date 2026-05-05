import { useState, type Ref } from "react"
import { Plus, Trash2 } from "lucide-react";

import { vscodeApi } from "@/vscodeApi";
import { cn } from "@/lib/utils"

import { UsageInfo, DiffInfo } from '@/components/Info';
import { Menu } from '@/components/SessionMenu'
import { JsonObject, ToolTipper } from 'veg-webview-common'
import { useChat } from "@/hooks/useChat";

export const Header = ({
  ref,
  className,
  userInput,
}:{
  ref?: Ref<HTMLDivElement>,
  className?: string,
  userInput?: any
}) => {
  const {
    sid,
    setPos,
    diff,
    usage,
    session,
    chatState,
  } = useChat();

  const [hidden, setHidden] = useState(true);

  const realInput = userInput || chatState?.userInput

  return (

    <div ref={ref} className={cn("flex flex-col gap-2", className)}>

      <div className="flex justify-between items-center gap-2">
        <span>{session?.state?.title || sid || "no session"}</span>
        <Menu hidden={hidden} setHidden={setHidden} refresh userInput={realInput}/>

        <ToolTipper label="Create">
          <Plus size={16}
            aria-label="create"
            className="hover:text-green-500"
            onClick={() => {
              console.log("Create!", realInput)
              vscodeApi.postMessage({
                type: 'session.create',
                payload: { 
                  focus: true, 
                  agent: realInput?.agent,
                  model: realInput?.model,
                  envName: realInput?.environ,
                },
              });
            }}
          />
        </ToolTipper>
        <ToolTipper label="Delete">
          <Trash2 size={16}
            aria-label="delete"
            className="hover:text-red-500"
            onClick={() => {
              console.log("Delete!", sid)
              vscodeApi.postMessage({
                type: 'session.delete',
                payload: { sid },
              });
            }}
          />
        </ToolTipper>
      </div>

      <div className="flex flex-col gap-2">
        <span className="font-thin italic text-xs">tasks & planning</span>
        <UsageInfo usage={usage} size={16}/>
      </div>

      { !hidden && <JsonObject data={{
        sid,
        userInput,
        usage,
        session,
        chatState,
        diff,
      }} />}
    </div>

  )
}