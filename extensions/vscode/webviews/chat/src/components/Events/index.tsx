import { useEffect } from "react";
import { useStickToBottomContext } from "use-stick-to-bottom";

import { cn } from "@/lib/utils";
import { useChat } from "@/hooks/useChat";

import {
  UnknownEvent,
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
  var merged: any[] = []
  events.forEach((E1: any, e1: number) => {
    // loop over earlier events
    for (var e2 = e1-1; e2 >= 0; e2--) {
      const E2 = merged[e2]
      const P2 = E2?.Content?.parts
      if (!P2) {
        continue
      }
      // if we have a matching invocation id, lets do some matching
      if (E1.InvocationID === E2.InvocationID) {
        const P1 = E1?.Content?.parts
        if (!P1) {
          continue
        }
        // todo, we need to loop over parts here too
        for (const p1 of P1) {
          if (!p1?.functionResponse) {
            continue
          }
          // console.log("functionCall", pi[0], ei)
          for (const pi2 in P2) {
            const p2 = P2[pi2]
            if (!p2?.functionCall || p1.functionResponse.id !== p2.functionCall.id) {
              continue
            }
            E2.Content.parts[pi2].functionResponse = p1.functionResponse
            // console.log("MATCH!", ej.InvocationID, ej, merged[j], pi[0], pj[0])
            break;
          }
        }
        merged[e2] = E2
        break;
      }
    }
    merged.push(E1)
  })

  // console.log("Events.merged", merged)

  return (
    <div className="flex-grow flex flex-col mx-2 gap-1 overflow-y">
      {merged?.map((e: any, pos: number) => {
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

