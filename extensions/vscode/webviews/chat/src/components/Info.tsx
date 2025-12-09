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

import { JsonObject, ToolTipper } from "veg-webview-common";
import { Menu } from '@/components/SessionMenu'

import { cn } from "@/lib/utils";


export const TimeInfo = ({ timestamp, className }: { timestamp?: string, className?: string }) => {
  if (!timestamp || timestamp === "") {
    return null
  }
  const ts = new Date(timestamp)
  return (
    <div className={className}>
      {ts.toLocaleString()}
    </div>
  )
}

export const MetaInfo = ({ evt, size }: { evt?: any, size: any }) => {
  return (
    <div className="flex gap-1 align-bottom align-end">
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
    </div>
  )
}

export function UsageNumber(num?: number): string {
  if (!num) {
    return "0"
  }
  if (num < 10000) {
    return `${num}`
  }
  if (num < 500000) {
    return `${(num / 1000.0).toFixed(1)}k`
  }
  if (num < 1000000) {
    return `${(num / 1000000.0).toFixed(2)}M`
  }
  return `${(num / 1000000.0).toFixed(1)}M`
}

export const UsageInfo = ({ evt, usage, size }: { evt?: any, usage?: any, size: any }) => {
  var u = evt?.UsageMetadata || usage || {}
  return (
    <div className="flex gap-2 h-4">
      <ToolTipper side="bottom" label="cached input tokens">
        <div className="flex gap-1">
          <DatabaseBackup size={size}/>
          {UsageNumber(u.cachedContentTokenCount) || "0"}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label="normal input tokens">
        <div className="flex gap-1">
          <GraduationCap size={size}/>
          {UsageNumber(u.promptTokenCount - (u.cachedContentTokenCount || 0))}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label="thinking tokens">
        <div className="flex gap-1">
          <BrainCircuit size={size}/>
          {UsageNumber(u.thoughtsTokenCount)}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label="response tokens">
        <div className="flex gap-1">
          <MessageSquareMore size={size}/>
          {UsageNumber(u.candidatesTokenCount)}
        </div>
      </ToolTipper>

      <div>
        <EqualApproximately size={size}/>
      </div>

      <ToolTipper side="bottom" label="total input">
        <div className="flex gap-1">
          <PanelRightClose size={size}/>
          {UsageNumber(u.promptTokenCount)}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label="total output">
        <div className="flex gap-1">
          <PanelLeftOpen size={size}/>
          {UsageNumber(u.candidatesTokenCount + u.thoughtsTokenCount)}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label="total tokens">
        <div className="flex gap-1">
          <SquareSigma size={size}/>
          {UsageNumber(u.totalTokenCount)}
        </div>
      </ToolTipper>
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
    <div className="flex gap-1 h-4">
      <ToolTipper side="bottom" label={ap?.join("\n") || "nothing new"}>
        <div className={cn("flex gap-1 hover:text-green-500", ap?.length && "text-green-500")}>
          <FilePlus size={size}/>
          {ap?.length || 0}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label={mp?.join("\n") || "no edits"}>
        <div className={cn("flex gap-1 hover:text-yellow-500", mp?.length && "text-yellow-500")}>
          <FilePen size={size}/>
          {mp?.length || 0}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label={dp?.join("\n") || "did you take out the trash?"}>
        <div className={cn("flex gap-1 hover:text-red-500", dp?.length && "text-red-500")}>
          <FileX size={size}/>
          {dp?.length || 0}
        </div>
      </ToolTipper>
    </div>
  )
}