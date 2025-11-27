import { useState } from "react"

import { cn } from "@/lib/utils";

import { Markdown } from './Markdown'
import { EventDetails, JsonInfo, Menu } from "./Info";
import { BadgeQuestionMark, Check, X } from "lucide-react";

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
  console.log("Events", events)
  // todo, coalesce events here

  return (
    <div className="flex-grow flex flex-col mx-2 gap-4 overflow-y">
      {events?.map((e: any, pos: number) => {
        console.log("event.loop", pos, currPos)
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
    return <UserMessage sid={sid} pos={pos} setPos={setPos} evt={evt}/>
  }
  if (evt?.Content?.role === "model") {
    return <ModelMessage sid={sid} pos={pos} setPos={setPos} evt={evt}/>
  }
  return <UnknownEvent sid={sid} pos={pos} setPos={setPos} data={evt} msg="missing Content.role"/>
}

const UnknownEvent = ({sid, pos, setPos, msg, data}: {sid: string, pos: number, setPos: any, msg?: string, data: any}) => {
  return (
    <div className="flex flex-col text-sm m-2 px-2 py-1 border-red-800 rounded gap-2">
      <div className="flex justify-between items-center gap-2 p-2">
        <div className="font-bold">unknown event</div>
        <div className="font-thin">{msg}</div>
      </div>
      <EventDetails sid={sid} pos={pos} setPos={setPos} evt={data} />
    </div>
  )
}


const UserMessage = ({sid, pos, setPos, evt}:{sid: string, pos: number, setPos: any, evt: any}) => {
  return (
    <div className={cn("ml-16 rounded-lg border", "border-sky-300")}>
      <div className="flex flex-col m-2 p-2 gap-2">
        { evt.Content.parts.map((p: any) => <MessagePart sid={sid} pos={pos} setPos={setPos} part={p} evt={evt}/>)}
        <EventDetails sid={sid} pos={pos} setPos={setPos} evt={evt} />
      </div>
    </div>
  )
}

const ModelMessage = ({sid, pos, setPos, evt}:{sid: string, pos: number, setPos: any, evt: any}) => {
  return (
    <div className="mr-16 rounded-lg border border-green-300">
      <div className="flex flex-col m-2 p-2 gap-2">
        { evt.Content.parts.map((p: any) => <MessagePart sid={sid} pos={pos} setPos={setPos} part={p} evt={evt}/>)}
        <EventDetails sid={sid} pos={pos} setPos={setPos} evt={evt} />
      </div>
    </div>
  )
}

const MessagePart = ({ sid, pos, setPos, part, evt }:{ sid: string, pos: number, setPos: any, part: any, evt: any }) => {
  if (part.text) { return <TextPart part={part} evt={evt}/> }
  if (part.functionCall) { return <FuncCall part={part} evt={evt}/> }
  if (part.functionResponse) { return <FuncResp part={part} evt={evt}/> }

  // what is it?
  return <UnknownEvent sid={sid} pos={pos} setPos={setPos} data={part} msg="unknown part"/>
}

const TextPart = ({ part }:{ part: any, evt: any }) => {
  // TODO, add copy button, size limiter (3 options)
  return (
    <div className="p-4 dark:prose-invert prose-sm prose-stone">
      <Markdown>{part.text}</Markdown>
    </div>
  )
}

const FuncCall = ({ part }:{ part: any, evt: any }) => {
  const fn = part.functionCall?.name as string
  const args = part.functionCall?.args
  const fnArgs = f2NameArgs[fn]
  const argVals = fnArgs?.map(a=> args[a] || a)

  return (
    <div  className="mr-auto flex gap-2 items-baseline">
      <span className="font-heavy">{fn}</span>
      <span className="font-thin">{argVals.join(" ")}</span>
    </div>
  )
}
const FuncResp = ({ part }:{ part: any, evt: any }) => {
  const fn = part.functionResponse?.name as string
  const resp = part.functionResponse?.response
  const fnArgs = f2NameArgs[fn]
  const argVals = fnArgs?.map(a=> resp[a] || a)
  // const status = resp.status
  const err = resp.error

  return (
    <div className="mr-auto flex flex-col gap-2">
      <div  className="flex gap-2 items-baseline">
        {/* <RespStatus size={14} status={status} /> */}
        <span className="font-heavy">{fn}</span>
        <span className="font-thin">{argVals.join(" ")}</span>
      </div>
      { err && (
        <div className="m-1 p-3 mx-auto rounded border border-red-600 font-thin">
          {err}
        </div>
      )}
    </div>
  )
}

const RespStatus = ({status, size}:{status:string, size: number}) => {
  switch (status) {
    case "ok":
      return <Check size={size} color="green" />
      break;

    case "error":
      return <X size={size} color="red" />
      break;

    default:
      return <BadgeQuestionMark size={size} />
  }
}

const NameArgTitle = ({name, args}:{name: string, args: string[]}) => {
  return (
    <div className="flex gap-2 overflow-x items-baseline">
      <div className="font-bold">{name}</div>
      <div className="monospace text-xs">{args.join("")}</div>
    </div>
  )
}

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
