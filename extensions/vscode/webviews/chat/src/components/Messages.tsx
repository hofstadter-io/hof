import { useState } from "react"

import { Markdown } from './Markdown'
import { EventDetails, JsonInfo, UsageInfo, TimeInfo } from "./Info";
import { FunctionCall, FunctionResp } from './FunctionCalls'

export const Events = ({
  events,
  messagesEndRef
}:{
  events: any[],
  messagesEndRef?: any,
}) => {
  console.log("Events", events, messagesEndRef)
  if (!events?.length) {
    return null
  }
  return (
    <div className="flex-grow mx-2 overflow-y-auto">
      {events?.map((e: any) => {
        if (!(e?.Content)) {
          return null
        }
        return (
          <Event key={e.ID} evt={e}/>
        )
      })}
      <div ref={messagesEndRef} />
    </div>
  )
}


export const Event = ({evt}: {evt: any}) => {
  // console.log("session event:", evt)
  const msg = evt.Content.parts[0]

  if (evt.Content.role === "user") {
    if (msg.text) { return <UserMessage evt={evt}/> }
    if (msg.functionResponse) { return <FunctionResp evt={evt}/> }
  }

  if (evt.Content.role === "model") {
    // if (msg.text && evt.FinishReason === "STOP") { return <ModelMessage evt={evt}/> }
    // if (msg.text && evt.FinishReason === "") { return <ThinkingMessage evt={evt}/> }
    if (msg.functionCall) { return <FunctionCall evt={evt}/> }
    if (msg.text) { return <ModelMessage evt={evt}/> }
  }

  return (
    <div className="flex flex-col text-sm my-2 mx-4 px-2 py-1 bg-red-800">
      <div className="font-bold">unknown event type</div>
      <JsonInfo data={evt} />
    </div>
  )
}

export const UserMessage = ({evt}:{evt: any}) => {
  return (
    <div className="ml-16 p-2 rounded-lg bg-slate-800/80 flex flex-col">
      <div className="p-2 dark:prose-invert prose-sm prose-stone">
        <Markdown>{evt.Content.parts.map((p: any) => p.text).join("\n\n")}</Markdown>
      </div>
      <EventDetails evt={evt} />
    </div>
  )
}

export const ModelMessage = ({evt}:{evt: any}) => {
  return (
    <div className="mr-16 my-2 p-3 rounded-lg bg-stone-800/80 flex flex-col gap-1">
      <div className="p-3 dark:prose-invert prose-sm prose-stone">
        <Markdown>{evt.Content.parts.map((p: any) => p.text).join("\n\n")}</Markdown>
      </div>
      <EventDetails evt={evt} />
    </div>
  )
}

// @ts-ignore
export const ThinkingMessage = ({evt}:{evt: any}) => {
  const [hidden, setHidden] = useState(true);
  return (
    <div className="mr-16 my-2 p-3 rounded-lg bg-purple-800/80 flex flex-col gap-1">
      <div className="flex justify-between" onClick={() => setHidden(!hidden)}>
        <span hidden={!hidden} className="font-italic">thinking...</span>
        <div className="flex justify-between">
          <UsageInfo evt={evt} size={14} />
          <TimeInfo timestamp={evt.Timestamp} />
        </div>
      </div>
      <div hidden={hidden} className="p-3 bg-purple-900/80 dark:prose-invert prose-sm prose-stone">
        <Markdown>{evt.Content.parts.map((p: any) => p.text).join("\n\n")}</Markdown>
      </div>
    </div>
  )
}
