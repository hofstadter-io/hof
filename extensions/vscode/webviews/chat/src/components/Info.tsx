import { useState } from "react"
import {
  Circle,
  CircleQuestionMark,
  CircleCheckBig,
  CircleDashed,
  GraduationCap,
  OctagonAlert,
  PanelLeftOpen,
  PanelRightClose,
  BrainCircuit,
  EqualApproximately,
  SquareSigma,
  DatabaseBackup,
  MessageSquareMore,
  FilePlus,
  FileX,
  FilePen,
} from 'lucide-react'

import JsonView from '@uiw/react-json-view';
import { vscodeTheme } from '@uiw/react-json-view/vscode';

import { Tooltipped } from "@/components/Tooltipped";
import { Menu } from '@/components/SessionMenu'

import { cn } from "@/lib/utils";

export const EventDetails = ({sid, pos, setPos, evt}:{sid: string, pos: number, setPos: any, evt: any}) => {
  const [hidden, setHidden] = useState(true);
  return (
    <div className="flex flex-col gap-1 mt-2">
      <div className="flex justify-between text-xs font-thin">
        <TimeInfo timestamp={evt.Timestamp} />
        <span className="text-sm">{evt.Author}</span>
      </div>
      <div className="flex flex-col justify-between items-center">
        <UsageInfo evt={evt} size={16} />
        <Menu sid={sid} pos={pos} setPos={setPos} hidden={hidden} setHidden={setHidden} />
      </div>
      { !hidden && <JsonInfo data={evt} /> }
    </div>
  )
}


export const JsonInfo = ({ hidden, data }: { hidden?: boolean, data: any }) => {
  return (
    <div hidden={hidden} className="transition hidden:h-0 whitespace-pre-line m-2 max-h-200 overflow-y-auto text-sm">
      <JsonView value={data} style={vscodeTheme} displayDataTypes={false} indentWidth={12} shortenTextAfterLength={40} className="p-2"/>
    </div>
  )
}

export const TimeInfo = ({ timestamp }: { timestamp?: string }) => {
  if (!timestamp || timestamp === "") {
    return null
  }
  const ts = new Date(timestamp)
  return (
    <span className="">
      {ts.toLocaleString()}
    </span>
  )
}

export const MetaInfo = ({ evt, size }: { evt?: any, size: any }) => {
  return (
    <span className="flex gap-1 align-bottom align-end">
      { evt.FinishReason === "STOP" ? 
        <OctagonAlert size={size} strokeWidth={2} className="text-red-600"/>
        :
        <CircleQuestionMark size={size} strokeWidth={2} className="text-yellow-600" />
      }
      { evt.Partial ? 
        <CircleDashed size={size} strokeWidth={2} className="text-yellow-600"/>
        :
        <Circle size={size} strokeWidth={2} className="text-blue-600"/>
      }
      { evt.TurnComplete ? 
        <CircleCheckBig size={size} strokeWidth={2} className="text-green-600"/>
        :
        <Circle size={size} strokeWidth={2} className="text-yellow-600"/>
      }
    </span>
  )
}

export const UsageInfo = ({ evt, usage, size }: { evt?: any, usage?: any, size: any }) => {
  var u = evt?.UsageMetadata || usage || {}
  return (
    <div className="flex gap-2">
      <span className="flex gap-1">
        <DatabaseBackup size={size}/>
        {u.cachedContentTokenCount || "0"}
      </span>
      <span className="flex gap-1">
        <GraduationCap size={size}/>
        {u.promptTokenCount - (u.cachedContentTokenCount || 0) || "0"}
      </span>
      <span className="flex gap-1">
        <BrainCircuit size={size}/>
        {u.thoughtsTokenCount || "0"}
      </span>
      <span className="flex gap-1">
        <MessageSquareMore size={size}/>
        {u.candidatesTokenCount || "0"}
      </span>

      <span>
        <EqualApproximately size={size}/>
      </span>

      <span className="flex gap-1">
        <PanelRightClose size={size}/>
        {u.promptTokenCount || "0"}
      </span>
      <span className="flex gap-1">
        <PanelLeftOpen size={size}/>
        {u.candidatesTokenCount + u.thoughtsTokenCount || "0"}
      </span>
      <span className="flex gap-1">
        <SquareSigma size={size}/>
        {u.totalTokenCount || "0"}
      </span>
    </div>
  )
   
}

export const DiffInfo = ({
  diff,
  size = 16,
}: {
  diff: any
  size?: number
}) => {
  const ap = diff?.addpaths
  const mp = diff?.modpaths
  const dp = diff?.delpaths
  return (
    <div className="flex gap-1">
      <Tooltipped label={ap?.join("\n") || "nothing new"}>
        <div className={cn("flex gap-1 hover:text-green-500", ap?.length && "text-green-500")}>
          <FilePlus size={size}/>
          {ap?.length || 0}
        </div>
      </Tooltipped>
      <Tooltipped label={mp?.join("\n") || "no edits"}>
        <div className={cn("flex gap-1 hover:text-yellow-500", mp?.length && "text-yellow-500")}>
          <FilePen size={size}/>
          {mp?.length || 0}
        </div>
      </Tooltipped>
      <Tooltipped label={dp?.join("\n") || "did you take out the trash?"}>
        <div className={cn("flex gap-1 hover:text-red-500", dp?.length && "text-red-500")}>
          <FileX size={size}/>
          {dp?.length || 0}
        </div>
      </Tooltipped>
    </div>
  )
}