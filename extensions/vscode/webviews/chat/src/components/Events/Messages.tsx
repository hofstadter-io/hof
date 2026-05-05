import { cn } from "@/lib/utils";

import {
  BrainCircuit,
  CheckCircle,
  MoveRight,
  Siren,
} from "lucide-react";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"

import { Markdown, TailwindClasses } from '@/components/Markdown'
import { EventDetails } from "@/components/Events/Details";
import { LightDetails } from "@/components/Events/LightDetails";


const OPEN_BY_DEFAULT = ["fs_edit", "exec"];

const shouldOpenDetails = (evt: any) => {
  const parts = evt?.Content?.parts || [];
  return parts.some((p: any) => 
    p.functionCall && OPEN_BY_DEFAULT.includes(p.functionCall.name)
  );
};

export const UnknownEvent = ({
  pos,
  msg,
  evt,
}: {
  pos: number,
  msg?: string,
  evt: any
}) => {
  const defaultOpen = shouldOpenDetails(evt) ? "details" : undefined;
  return (
    <div className={cn(
      "ml-40 my-2 py-px pl-2px rounded",
      "bg-linear-to-r from-red-500/80 from-20% via-[#1e1e1e] via-50% to-[#1e1e1e]",
    )}>
      <div className={cn(
        "flex flex-col p-2 rounded",
        "bg-linear-to-r from-slate-800/50 from-20% via-[#1e1e1e] via-40% to-[#1e1e1e]",
      )}>
        <div className="flex flex-col">
          <div className="font-bold">unknown event</div>
          <div className="font-thin">{msg}</div>
        </div>
        <div className="mt-[-1em] w-full">
          <Accordion type="single" collapsible defaultValue={defaultOpen}>
            <AccordionItem value="details">
              <AccordionTrigger className="h-3">
                <div className="flex gap-2 font-thin items-center ml-2">
                  {/* TODO, this needs to com from the session */}
                  {evt?.Actions?.StateDelta?.currEnv && (
                    <span className="text-amber-500 font-mono text-xs">
                      ({evt.Actions.StateDelta.currEnv.split(':').pop()})
                    </span>
                  )}
                <span className="text-violet-500 font-mono text-xs">
                  [{pos}]
                </span>
              </div>
            </AccordionTrigger>
              <AccordionContent>
                <LightDetails evt={evt} />
                <EventDetails pos={pos} evt={evt}/>
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </div>
      </div>
    </div>
  )
}

export const StopEvent = ({
  pos,
  msg,
  evt,
}: {
  pos: number,
  msg?: string,
  evt: any
}) => {
  const defaultOpen = shouldOpenDetails(evt) ? "details" : undefined;
  return (
    <div className={cn(
      "ml-10 mr-20 my-2 py-px pl-0.5 rounded",
      "bg-linear-to-r from-red-500/80 from-20% via-[#1e1e1e] via-50% to-[#1e1e1e]",
    )}>
      <div className={cn(
        "flex flex-col p-2 rounded",
        "bg-linear-to-r from-slate-800/50 from-20% via-[#1e1e1e] via-40% to-[#1e1e1e]",
      )}>
        <div className="flex flex-col">
          <div className="font-thin text-lg pl-2">User Interrupt</div>
        </div>
        <div className="mt-[-1em] w-full">
          <Accordion type="single" collapsible defaultValue={defaultOpen}>
            <AccordionItem value="details">
              <AccordionTrigger className="h-3">
                <div className="flex gap-2 font-thin items-center ml-2">
                  {/* TODO, this needs to com from the session */}
                  {evt?.Actions?.StateDelta?.currEnv && (
                    <span className="text-amber-500 font-mono text-xs">
                      ({evt.Actions.StateDelta.currEnv.split(':').pop()})
                    </span>
                  )}
                <span className="text-violet-500 font-mono text-xs">
                  [{pos}]
                </span>
              </div>
            </AccordionTrigger>
              <AccordionContent>
                <LightDetails evt={evt} />
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
  const defaultOpen = shouldOpenDetails(evt) ? "details" : undefined;
  return (
    <div className={cn(
      "mx-20 my-2 py-px pl-0.5 rounded",
      "bg-linear-to-r from-sky-500/80 from-20% via-[#1e1e1e] via-50% to-[#1e1e1e]",
    )}>
      <div className={cn(
        "flex flex-col p-2 rounded",
        "bg-linear-to-r from-slate-800/90 from-20% via-[#1e1e1e] via-40% to-[#1e1e1e]",
      )}>
        <div className="flex flex-col">
          { evt?.Content?.parts && evt.Content.parts.map((p: any) => <MessagePart pos={pos} part={p} evt={evt}/>) }
          { hasStateDelta && Object.entries(evt.Actions.StateDelta).map(([key, val]) => {
              return (
                <div className="flex gap-1 items-center px-2 border-l-3 border-amber-500/80">
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
        <div className="mt-[-1em] w-full">
          <Accordion type="single" collapsible defaultValue={defaultOpen}>
            <AccordionItem value="details">
              <AccordionTrigger className="h-3">
              <div className="flex gap-2 font-thin items-center ml-2">
                {evt?.Actions?.StateDelta?.currEnv && (
                  <span className="text-amber-500 font-mono text-xs">
                    ({evt.Actions.StateDelta.currEnv.split(':').pop()})
                  </span>
                )}
                <span className="text-violet-500 font-mono text-xs">
                  [{pos}]
                </span>
              </div>
            </AccordionTrigger>
              <AccordionContent>
                <LightDetails evt={evt} />
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
  const defaultOpen = shouldOpenDetails(evt) ? "details" : undefined;
  return (
    <div className={cn(
      "mr-20 py-px pl-0.5 rounded",
      "bg-linear-to-r from-lime-500/80 from-5% via-[#1e1e1e] via-10% to-[#1e1e1e]",
    )}>
      <div className={cn(
        "flex flex-col p-2 rounded",
        "bg-[#1e1e1e]"
      )}>
        <div className="flex flex-col gap-1">
          { evt.Content.parts.map((p: any) => <MessagePart pos={pos} part={p} evt={evt}/>)}
        </div>
        <div className="mt-[-1em] w-full">
          <Accordion type="single" collapsible defaultValue={defaultOpen}>
            <AccordionItem value="details">
              <AccordionTrigger className="h-4">
              <div className="flex gap-2 font-thin items-center ml-2">
                {evt?.Actions?.StateDelta?.currEnv && (
                  <span className="text-amber-500 font-mono text-xs">
                    ({evt.Actions.StateDelta.currEnv.split(':').pop()})
                  </span>
                )}
                <span className="text-violet-500 font-mono text-xs">
                  [{pos}]
                </span>
              </div>
            </AccordionTrigger>
              <AccordionContent>
                <LightDetails evt={evt} />
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
    <div 
      className={cn(
        "pl-2 py-1 mr-8 flex flex-col",
        ...TailwindClasses,
      )}
    >
      <Markdown>{part.text}</Markdown>
    </div>
  )
}

export const FuncCall = ({ part, evt }:{ part: any, evt: any }) => {
  const fn = part.functionCall?.name as string
  const args = part.functionCall?.args
  const fnArgs = f2NameArgs[fn]
  const argVals = fnArgs?.map(a=> {
    if(a in args) {
      return args[a]
    }
  })
  const resp = part.functionResponse?.response

  const isPlanning = fn === "cache_put" && args["key"] === "planning"
  const isExec = fn === "exec"

  return (
    <div  className="mr-auto ml-1 pl-1 border-l-3 border-yellow-500 flex flex-col gap-2">
      <div  className="flex gap-2 items-baseline">
        <span className="font-heavy text-yellow-500">{fn}</span>
        <span className="font-thin">{(argVals || []).join(" ")}</span>
        { !resp && <BrainCircuit size={10} strokeWidth={1} className="text-yellow-300 animate-ping" />}
        { resp?.status === "ok" && <CheckCircle size={12} className="text-lime-500" />}
        { resp?.status === "error" && <Siren size={12} className="text-red-500" />}
      </div>
      
      { isPlanning && (
        <pre className="m-2 p-2 border border-violet-500">
          {args.value}
        </pre>
      )}
      { isExec && (
        <pre className="m-2 p-2 border border-green-500">
          {args.script}
        </pre>
      )}
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

  "exec": ["key"],

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
