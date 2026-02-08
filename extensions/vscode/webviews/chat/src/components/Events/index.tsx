import { useEffect } from "react";
import { useStickToBottomContext } from "use-stick-to-bottom";

import { cn, processEvents } from "@/lib/utils";
import { useChat } from "@/hooks/useChat";

import {
  UnknownEvent,
  StopEvent,
  UserMessage,
  ModelMessage,
} from './Messages'

export const Events = ({
  messagesEndRef
}:{
  messagesEndRef?: any,
}) => {
  const {
    pos: currPos,
    session,
  } = useChat();

  const events = session?.events;

  const { isAtBottom, scrollToBottom } = useStickToBottomContext();

  useEffect(() => {
    isAtBottom && scrollToBottom()
  }, [events])

  if (!events?.length) {
    return null
  }

  // TODO, extract this
  // console.log("Events.events", events)
  // todo, coalesce events here
  const { merged } = processEvents(events)

  // console.log("Events.merged", merged)

  return (
    <div className="grow flex flex-col mx-2 gap-1 overflow-y">
      {merged?.map((e: any, pos: number) => {
        if (!e) {
          return null
        }
        return (
          <div className={cn(pos === currPos && "bg-violet-500/30 rounded")}>
            <Event pos={pos} key={e.ID} evt={e} />
          </div>
        )
      })}
      <div ref={messagesEndRef} />
    </div>
  )
}

export const Event = ({
  pos,
  evt,
}: {
  pos: number,
  evt: any,
}) => {
  // HACK: to ignore function responses, which get merged with the call and rendered together
  var fnRespCnt: number = 0
  var fnCallCnt: number = 0
  if (evt?.Content?.parts) {
    for (const p of evt?.Content?.parts) {
      if (p.functionCall) {
        fnCallCnt++
      }
      if (p.functionResponse) {
        fnRespCnt++
      }
    }
    if (fnRespCnt > 0 && fnCallCnt === 0) {
      return null
    }
  }

  if (evt?.Author === "user") {
    return <UserMessage pos={pos} evt={evt}/>
  } else {

    if (evt?.TurnComplete && evt?.Interrupted) {
    return <StopEvent pos={pos} evt={evt}/>
    }

    // weird stop message
    if (!evt?.Content && evt?.ErrorCode === "STOP" && evt?.FinishReason === "") {
      return null
    }
    // all agent messages should have parts?
    if (!evt?.Content?.parts) {
      return <UnknownEvent pos={pos} evt={evt} msg="missing Content.parts"/>
    }

    return <ModelMessage pos={pos} evt={evt}/>
  }
}

