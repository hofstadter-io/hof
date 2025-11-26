import { useState } from "react"

import { Markdown } from './Markdown'
import { EventDetails, JsonInfo, UsageInfo, TimeInfo, Menu } from "./Info";

export const Events = ({
  events,
  messagesEndRef
}:{
  events: any[],
  messagesEndRef?: any,
}) => {
  if (!events?.length) {
    return null
  }
  console.log("Events", events)
  // todo, coalesce events here

  return (
    <div className="flex-grow mx-2 overflow-y-auto">
      {events?.map((e: any) => {
        // console.log("event.loop", e)
        return (
          <Event key={e.ID} evt={e}/>
        )
      })}
      <div ref={messagesEndRef} />
    </div>
  )
}


export const Event = ({evt}: {evt: any}) => {
  if (evt?.Content?.role === "user") {
    return <UserMessage evt={evt}/>
  }
  if (evt?.Content?.role === "model") {
    return <ModelMessage evt={evt}/>
  }
  return <UnknownEvent data={evt} msg="missing Content.role"/>
}

const UnknownEvent = ({msg, data}: {msg?: string, data: any}) => {
  return (
    <div className="flex flex-col text-sm my-2 mx-4 px-2 py-1 bg-red-800">
      <div className="font-bold">unknown event</div>
      <div className="font-thin">{msg}</div>
      <JsonInfo data={data} />
    </div>
  )
}


const UserMessage = ({evt}:{evt: any}) => {
  if (!evt?.Content?.parts) {
    return <UnknownEvent data={evt} msg="missing Content.parts"/>
  }
  return (
    <div className="ml-16 p-2 rounded-lg bg-slate-800/80 flex flex-col">
      { evt.Content.parts.map((p: any) => <MessagePart part={p} evt={evt}/>)}
      <EventDetails evt={evt} />
    </div>
  )
}

const ModelMessage = ({evt}:{evt: any}) => {
  if (!evt?.Content?.parts) {
    return <UnknownEvent data={evt} msg="missing Content.parts"/>
  }
  return (
    <div className="mr-16 my-2 p-3 rounded-lg bg-stone-800/80 flex flex-col gap-1">
      { evt.Content.parts.map((p: any, n: any) => <MessagePart key={`${evt.id}`} part={p} evt={evt}/>)}
      <EventDetails evt={evt} />
    </div>
  )
}

const MessagePart = ({ part, evt }:{ part: any, evt: any }) => {
  if (part.text) { return <TextPart part={part} evt={evt}/> }
  if (part.functionCall) { return <FuncCall part={part} evt={evt}/> }
  if (part.functionResponse) { return <FuncResp part={part} evt={evt}/> }

  // what is it?
  return <UnknownEvent data={part} msg="unknown part"/>
}

const TextPart = ({ part }:{ part: any, evt: any }) => {
  // TODO, add copy button, size limiter (3 options)
  return (
    <div className="p-2 dark:prose-invert prose-sm prose-stone">
      <Markdown>{part.text}</Markdown>
    </div>
  )
}

const NameArgTitle = ({name, args}:{name: string, args: string[]}) => {
  return (
    <div className="flex gap-2 overflow-x">
      <div className="font-bold">{name}</div>
      <div className="monospace text-xs">{args.join("")}</div>
    </div>
  )
}

const FuncTitle = ({name, args}:{name: string, args: any}) => {

  const f2NameArgs: Record<string,string[]> = {
    "cache_write": ["key"],
    "cache_put": ["key"],
    "cache_edit": ["key"],
    "cache_del": ["key"],
    "cache_remove": ["key"],

    "fs_read": ["path"],
    "fs_list": ["path"],
    "fs_grep": ["path", "regexp"],
    "fs_write": ["path"],
    "fs_edit": ["path"],
    "fs_del": ["path"],

    // legacy
    "read_file": ["path"],
    "read_dir":  ["path"],
    "tree_dir": ["path"],
    "write_file": ["path"],
    "cache_glob": ["path", "regexp"],
    "cache_grep": ["path", "regexp"],
    "cache_file": ["path"],
    "cache_dir": ["path"],
  }

  const fnArgs = f2NameArgs[name]

  return (
    <div>
    { args && fnArgs ?
      <NameArgTitle name={name} args={fnArgs.map(a=> args[a] || a)}/>
      :
      <NameArgTitle name={name} args={["???"]}/>
    }
    </div>
  )
}

const FuncCall = ({ part }:{ part: any, evt: any }) => {
  const fn = part.functionCall.name as string
  return <div  className="mr-auto">
      <FuncTitle name={fn} args={part.functionCall.args} />
    </div>

}
const FuncResp = ({ part }:{ part: any, evt: any }) => {
  console.log()
  const fn = part.functionResponse.name as string
  return (
    <div  className="mr-auto">
      <FuncTitle name={fn} args={part.functionResponse.response} />
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
