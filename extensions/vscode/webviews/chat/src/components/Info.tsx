import { useState } from "react"
import { vscodeApi } from "@/vscodeApi";
import {
  Braces,
  Circle,
  CircleQuestionMark,
  CircleCheckBig,
  CircleDashed,
  Columns3,
  GitFork,
  GitPullRequestArrow,
  RefreshCcw,
  GraduationCap,
  OctagonAlert,
  PanelLeftOpen,
  PanelRightClose,
  BrainCircuit,
  EqualApproximately,
  SquareSigma,
  DatabaseBackup,
  MessageSquareMore,
} from 'lucide-react'

import JsonView from '@uiw/react-json-view';
import { vscodeTheme } from '@uiw/react-json-view/vscode';

export const EventDetails = ({evt}:{evt: any}) => {
  const [hidden, setHidden] = useState(true);
  return (
    <div className="flex flex-col">
      <div className="flex gap-2">
        <UsageInfo evt={evt} size={16} />
        <MetaInfo evt={evt} size={12} />
        <span className="text-sm font-thin ml-auto m-0 p-0">{evt.Author}</span>
        <div className="ml-auto flex justify-end items-center gap-2">
          <TimeInfo timestamp={evt.Timestamp} />
          <Braces size={16} onClick={(e) => {
            e.preventDefault()
            e.stopPropagation()
            setHidden(!hidden)
          }} />
        </div>
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
    <span className="ml-auto font-thin text-sm">
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

export const Menu = ({
  sid,
  hidden,
  setHidden,
}:{
  sid: string,
  hidden: boolean,
  setHidden: (prev: any) => any
}) => {
  return (
    <div className="ml-auto flex justify-end items-center gap-2">
      <Columns3 size={16}
        onClick={() => {
          vscodeApi.postMessage({
            type: "session.diff",
            payload: {
              sid,
            }
          })
        }}
      />
      <GitPullRequestArrow size={16} />
      <GitFork size={16} />
      <RefreshCcw size={16}
        onClick={() => {
          vscodeApi.postMessage({
            type: "session.get",
            payload: {
              sid,
            }
          })
        }}
      />
      <Braces size={16} onClick={(e) => {
        e.preventDefault()
        e.stopPropagation()
        setHidden(!hidden)
      }} />
    </div>
  )

}