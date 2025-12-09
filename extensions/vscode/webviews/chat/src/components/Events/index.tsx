import { cn } from "@/lib/utils";

import {
  UnknownEvent,
  UserMessage,
  ModelMessage,
} from './Messages'

export const Events = ({
  sid,
  currPos,
  setPos,
  events,
  messagesEndRef
}:{
  sid: string,
  currPos: number,
  setPos: any,
  events: any[],
  messagesEndRef?: any,
}) => {
  if (!events?.length) {
    return null
  }

  // TODO, extract this
  // console.log("Events.events", events)
  // todo, coalesce events here
  var merged: any[] = []
  events.forEach((E1, e1) => {
    // loop over earlier events
    for (var e2 = e1-1; e2 >= 0; e2--) {
      const E2 = merged[e2]
      // if we have a matching invocation id, lets do some matching
      if (E1.InvocationID === E2.InvocationID) {
        const P1 = E1?.Content?.parts
        const P2 = E2?.Content?.parts
        // todo, we need to loop over parts here too
        for (const p1 of P1) {
          if (!p1?.functionResponse) {
            continue
          }
          // console.log("functionCall", pi[0], ei)
          for (const pi2 in P2) {
            const p2 = P2[pi2]
            if (!p2.functionCall || p1.id !== p2.id) {
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
    <div className="flex-grow flex flex-col mx-2 gap-3 overflow-y">
      {merged?.map((e: any, pos: number) => {
        return (
          <div className={cn(pos === currPos && "bg-violet-500/30 rounded")}>
            <Event sid={sid} pos={pos} setPos={setPos} key={e.ID} evt={e}/>
          </div>
        )
      })}
      <div ref={messagesEndRef} />
    </div>
  )
}

export const Event = ({sid, pos, setPos, evt}: {sid: string, pos: number, setPos: any, evt: any}) => {
  if (!evt?.Content?.parts) {
    return <UnknownEvent sid={sid} pos={pos} setPos={setPos} data={evt} msg="missing Content.parts"/>
  }
  if (evt?.Content?.role === "user") {

    var found = false
    evt.Content.parts.forEach((p: any) => {
      if (p.text) {
        found = true
      }
    })
    if (!found) {
      return null
    }

    return <UserMessage sid={sid} pos={pos} setPos={setPos} evt={evt}/>
  }
  if (evt?.Content?.role === "model") {
    return <ModelMessage sid={sid} pos={pos} setPos={setPos} evt={evt}/>
  }
  return <UnknownEvent sid={sid} pos={pos} setPos={setPos} data={evt} msg="missing Content.role"/>
}

