import { cn } from "@/lib/utils";

import { Markdown } from '@/components/Markdown'
import { EventDetails } from "@/components/Events/Details";
import { FuncCall } from "./ToolCall";


export const UnknownEvent = ({sid, pos, setPos, msg, data}: {sid: string, pos: number, setPos: any, msg?: string, data: any}) => {
  return (
    <div className="flex flex-col text-sm px-2 py-1 border-red-800 rounded gap-2">
      <div className="flex justify-between items-center gap-2 p-2">
        <div className="font-bold">unknown event</div>
        <div className="font-thin">{msg}</div>
      </div>
      <EventDetails sid={sid} pos={pos} setPos={setPos} evt={data} />
    </div>
  )
}


export const UserMessage = ({sid, pos, setPos, evt}:{sid: string, pos: number, setPos: any, evt: any}) => {
  return (
    <div className={cn(
      "ml-40 py-[1px] pl-[2px] rounded",
      "bg-linear-to-r from-sky-500/80 from-[20%] via-[#1e1e1e] via-[50%] to-[#1e1e1e]",
    )}>
    <div className={cn(
      "flex flex-col p-2 rounded",
      "bg-linear-to-r from-slate-800/90 from-[20%] via-[#1e1e1e] via-[40%] to-[#1e1e1e]",
    )}>
      {/* <CopyClipboardButton source={part.text} positioning="ml-auto"/> */}
      { evt.Content.parts.map((p: any) => <MessagePart sid={sid} pos={pos} setPos={setPos} part={p} evt={evt}/>)}
      <EventDetails sid={sid} pos={pos} setPos={setPos} evt={evt} />
    </div>
    </div>
  )
}

export const ModelMessage = ({sid, pos, setPos, evt}:{sid: string, pos: number, setPos: any, evt: any}) => {
  return (
    <div className={cn(
      "py-[1px] pl-[2px] rounded",
      "bg-linear-to-r from-lime-500/80 from-[5%] via-[#1e1e1e] via-[10%] to-[#1e1e1e]",
    )}>
    <div className={cn(
      "flex flex-col p-2 rounded",
      "bg-[#1e1e1e]"
    )}>
      { evt.Content.parts.map((p: any) => <MessagePart sid={sid} pos={pos} setPos={setPos} part={p} evt={evt}/>)}
      <EventDetails sid={sid} pos={pos} setPos={setPos} evt={evt} />
    </div>
    </div>
  )
}

const MessagePart = ({ sid, pos, setPos, part, evt }:{ sid: string, pos: number, setPos: any, part: any, evt: any }) => {
  if (part.text) { return <TextPart part={part} evt={evt}/> }
  if (part.functionCall) { return <FuncCall part={part} evt={evt}/> }
  if (part.functionResponse) { return null }

  // what is it?
  return <UnknownEvent sid={sid} pos={pos} setPos={setPos} data={part} msg="unknown part"/>
}

const TextPart = ({ part, evt }:{ part: any, evt: any }) => {
  // TODO, add copy button, size limiter (3 options)
  return (
    <div className="p-2 flex flex-col dark:prose-invert prose-sm prose-stone">
      <Markdown>{part.text}</Markdown>
    </div>
  )
}
