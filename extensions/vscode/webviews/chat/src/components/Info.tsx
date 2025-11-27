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
  SquareTerminal,
} from 'lucide-react'

import JsonView from '@uiw/react-json-view';
import { vscodeTheme } from '@uiw/react-json-view/vscode';

export const EventDetails = ({sid, pos, setPos, evt}:{sid: string, pos: number, setPos: any, evt: any}) => {
  const [hidden, setHidden] = useState(true);
  return (
    <div className="flex flex-col gap-1 mt-2">
      <div className="flex justify-between text-xs font-thin">
        <TimeInfo timestamp={evt.Timestamp} />
        <span className="text-sm">{evt.Author}</span>
      </div>
      <div className="flex justify-between items-center">
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

export const Menu = ({
  sid,
  pos,
  setPos,
  hidden,
  refresh,
  // checkpoint,
  setHidden,
}:{
  sid: string,
  pos?: number,
  setPos: any,
  hidden: boolean,
  refresh?: boolean,
  checkpoint?: boolean,
  setHidden: (prev: any) => any
}) => {
  return (
    <div className="ml-auto flex justify-end items-center gap-2">
      <Columns3 size={16}
        aria-label="diff"
        className="hover:text-violet-500"
        onClick={() => {
          setPos(pos)
          vscodeApi.postMessage({
            type: "session.diff",
            payload: {
              sid,
              pos,
            }
          })
        }}
      />
      <GitPullRequestArrow size={16}
        aria-label="merge"
        className="hover:text-yellow-500"
        onClick={() => {
          vscodeApi.postMessage({
            type: "session.diff",
            payload: {
              sid,
              pos,
            }
          })
        }}
      />
      <GitFork size={16}
        aria-label="fork"
        className="hover:text-sky-500"
        onClick={() => {
          vscodeApi.postMessage({
            type: "session.create",
            payload: {
              from: sid,
              pos,
            }
          })
        }}
      />
      { refresh && <RefreshCcw size={16}
        aria-label="refresh"
        className="hover:text-sky-500"
        onClick={() => {
          vscodeApi.postMessage({
            type: "session.get",
            payload: {
              sid,
            }
          })
        }}
      /> }
      <div className="hover:text-green-500">
        <SquareTerminal size={16}
          aria-label="details"
          onClick={(e) => {
            e.preventDefault()
            e.stopPropagation()
            setHidden(!hidden)
          }}
        />
      </div>
      <div className="hover:text-sky-500">
        <Braces size={16}
          aria-label="details"
          className="hover:text-sky-500" 
          onClick={(e) => {
            e.preventDefault()
            e.stopPropagation()
            setHidden(!hidden)
          }}
        />
      </div>
    </div>
  )

}