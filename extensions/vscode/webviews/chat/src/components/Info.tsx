import { useState } from "react"
import {
  Circle,
  CircleQuestionMark,
  CircleCheckBig,
  CircleDashed,
  OctagonAlert,
  PanelLeftOpen,
  PanelRightClose,
  BrainCircuit,
  EqualApproximately,
  SquareSigma,
  FilePlus,
  FileX,
  FilePen,
  BotMessageSquare,
  BookMarked,
  NotebookTabs,
} from 'lucide-react'

import { ToolTipper } from "veg-webview-common";

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

export function PercentNumber(part: number, total: number) {
  if (!part || part < 0 || !total || total < 1) {
    return null
  }
  const p = Math.round(100 * (1 - ((total - part) / total)));
  return <span className="align-super text-[.5em]">{p}</span>
}

export const UsageInfo = ({ evt, usage, size }: { evt?: any, usage?: any, size: any }) => {
  var u = evt?.UsageMetadata || usage || {}

  const uncachedInput = u.promptTokenCount - (u.cachedContentTokenCount || 0)
  const totalOutput = (u.candidatesTokenCount || 0) + (u.thoughtsTokenCount || 0)

  return (
    <div className="flex gap-2 h-4">
      <ToolTipper side="bottom" label="cached input tokens">
        <div className="flex  text-lime-400">
          <BookMarked size={size}  className="mr-1"/>
          {UsageNumber(u.cachedContentTokenCount) || "0"}
          {PercentNumber(u.cachedContentTokenCount, u.promptTokenCount)}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label="normal input tokens">
        <div className="flex  text-amber-200">
          <NotebookTabs size={size} className="mr-1"/>
          {UsageNumber(uncachedInput)}
          {PercentNumber(uncachedInput, u.promptTokenCount)}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label="thinking tokens">
        <div className="flex  text-cyan-300">
          <BrainCircuit size={size} className="mr-1"/>
          {UsageNumber(u.thoughtsTokenCount)}
          {PercentNumber(u.thoughtsTokenCount, totalOutput)}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label="response tokens">
        <div className="flex  text-sky-400">
          <BotMessageSquare size={size} className="mr-1"/>
          {UsageNumber(u.candidatesTokenCount)}
          {PercentNumber(u.candidatesTokenCount, totalOutput)}
        </div>
      </ToolTipper>

      <div className="text-fuchsia-400">
        <EqualApproximately size={size}/>
      </div>

      <ToolTipper side="bottom" label="total input">
        <div className="flex  text-amber-300">
          <PanelRightClose size={size} className="mr-1"/>
          {UsageNumber(u.promptTokenCount)}
          {PercentNumber(u.promptTokenCount, u.totalTokenCount)}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label="total output">
        <div className="flex  text-blue-400">
          <PanelLeftOpen size={size} className="mr-1"/>
          {UsageNumber(totalOutput)}
          {PercentNumber(totalOutput, u.totalTokenCount)}
        </div>
      </ToolTipper>
      <ToolTipper side="bottom" label="total tokens">
        <div className="flex  text-fuchsia-400">
          <SquareSigma size={size} className="mr-1"/>
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