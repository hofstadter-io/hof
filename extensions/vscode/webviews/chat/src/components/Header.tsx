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
}:{
  ref?: Ref<HTMLDivElement>,
  className?: string,
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

  return (

    <div ref={ref} className={cn("flex flex-col gap-2 py-2", className)}>

      <div className="flex justify-between items-center gap-2">
        <span>{session?.state?.title || sid || "no session"}</span>
        <Menu hidden={hidden} setHidden={setHidden} refresh userInput={chatState?.userInput}/>

        <ToolTipper label="Create">
          <Plus size={16}
            aria-label="create"
            className="hover:text-green-500"
            onClick={() => {
              const state = vscodeApi.getState()
              console.log("Create!", state)
              vscodeApi.postMessage({
                type: 'session.create',
                payload: { 
                  focus: true, 
                  agent: chatState?.userInput?.agent,
                  model: chatState?.userInput?.model,
                  envName: chatState?.userInput?.environ,
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
              const state = vscodeApi.getState()
              console.log("Delete!", state)
              vscodeApi.postMessage({
                type: 'session.delete',
                payload: { sid },
              });
            }}
          />
        </ToolTipper>
      </div>

      <div className="flex justify-between items-center gap-2">
        <span className="font-thin italic text-xs">tasks & planning</span>
        <UsageInfo usage={usage} size={16}/>
      </div>

      { !hidden && <JsonObject data={{
        sid,
        usage,
        session,
        chatState,
        diff,
      }} />}
    </div>

  )
}