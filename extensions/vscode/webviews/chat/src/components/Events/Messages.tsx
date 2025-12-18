import { cn } from "@/lib/utils";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"

import { Markdown } from '@/components/Markdown'
import { EventDetails } from "@/components/Events/Details";
import { FuncCall } from "./ToolCall";
import { ArrowRight, MoveRight } from "lucide-react";


export const UnknownEvent = ({
  pos,
  msg,
  evt,
}: {
  pos: number,
  msg?: string,
  evt: any
}) => {
  return (
    <div className={cn(
      "ml-40 my-2 py-[1px] pl-[2px] rounded",
      "bg-linear-to-r from-red-500/80 from-[20%] via-[#1e1e1e] via-[50%] to-[#1e1e1e]",
    )}>
      <div className={cn(
        "flex flex-col p-2 rounded",
        "bg-linear-to-r from-slate-800/50 from-[20%] via-[#1e1e1e] via-[40%] to-[#1e1e1e]",
      )}>
        <div className="flex flex-col">
          <div className="font-bold">unknown event</div>
          <div className="font-thin">{msg}</div>
        </div>
        <div className="mt-[-12px] mr-auto">
          <Accordion type="single" collapsible>
            <AccordionItem value="details">
              <AccordionTrigger className="h-3"></AccordionTrigger>
              <AccordionContent>
                <EventDetails pos={pos} evt={evt}/>
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </div>
      </div>
    </div>
  )
}


export const UserMessage = ({pos, evt}:{pos: number, evt: any}) => {
  const hasStateDelta = evt?.Actions?.StateDelta && Object.keys(evt?.Actions?.StateDelta).length > 0
  return (
    <div className={cn(
      "ml-40 my-2 py-[1px] pl-[2px] rounded",
      "bg-linear-to-r from-sky-500/80 from-[20%] via-[#1e1e1e] via-[50%] to-[#1e1e1e]",
    )}>
      <div className={cn(
        "flex flex-col p-2 rounded",
        "bg-linear-to-r from-slate-800/90 from-[20%] via-[#1e1e1e] via-[40%] to-[#1e1e1e]",
      )}>
        <div className="flex flex-col">
          { evt?.Content?.parts && evt.Content.parts.map((p: any) => <MessagePart pos={pos} part={p} evt={evt}/>) }
          { hasStateDelta && Object.entries(evt.Actions.StateDelta).map(([key, val]) => {
              return (
                <div className="flex gap-1 items-center px-2 border-l-3 border-red-500">
                  <span className="font-bold">
                    ${key}
                  </span>
                  <MoveRight size={16} />
                  <span>
                    {val as string}
                  </span>
                </div>
              )
            })
          }
        </div>
        <div className="mt-[-12px] w-full">
          <Accordion type="single" collapsible>
            <AccordionItem value="details">
              <AccordionTrigger className="h-3"></AccordionTrigger>
              <AccordionContent>
                <EventDetails pos={pos} evt={evt} />
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </div>
      </div>
    </div>
  )
}

export const ModelMessage = ({pos, evt}:{pos: number, evt: any}) => {
  return (
    <div className={cn(
      "py-[1px] pl-[2px] rounded",
      "bg-linear-to-r from-lime-500/80 from-[5%] via-[#1e1e1e] via-[10%] to-[#1e1e1e]",
    )}>
      <div className={cn(
        "flex flex-col p-2 rounded",
        "bg-[#1e1e1e]"
      )}>
        <div className="flex flex-col gap-1">
          { evt.Content.parts.map((p: any) => <MessagePart pos={pos} part={p} evt={evt}/>)}
        </div>
        <div className="mt-[-12px] w-full">
          <Accordion type="single" collapsible>
            <AccordionItem value="details">
              <AccordionTrigger className="h-4"></AccordionTrigger>
              <AccordionContent>
                <EventDetails pos={pos} evt={evt} />
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </div>
      </div>
    </div>
  )
}

const MessagePart = ({ pos, part, evt }:{ pos: number, part: any, evt: any }) => {
  if (part.text) { return <TextPart part={part} evt={evt}/> }
  if (part.functionCall) { return <FuncCall part={part} evt={evt}/> }
  if (part.functionResponse) { return null }

  // what is it?
  return <UnknownEvent pos={pos} evt={evt} msg="unknown part"/>
}

const TextPart = ({ part, evt }:{ part: any, evt: any }) => {
  // TODO, add copy button, size limiter (3 options)
  return (
    <div className="pl-2 py-1 mr-8 flex flex-col dark:prose-invert prose-sm prose-stone">
      <Markdown>{part.text}</Markdown>
    </div>
  )
}
