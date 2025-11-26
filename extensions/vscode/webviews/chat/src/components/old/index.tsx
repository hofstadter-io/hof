import { useState } from "react"
import { cn } from "@/lib/utils";
import {Prism as SyntaxHighlighter} from 'react-syntax-highlighter'
import {vscDarkPlus} from 'react-syntax-highlighter/dist/esm/styles/prism'

import { CopyButton } from '../Markdown'

import { JsonInfo, UsageInfo, TimeInfo } from "../Info";

import {
  Braces,
  FileDiff,
} from 'lucide-react'

// change this to just the part?
export const FunctionCall = ({evt}:{evt: any}) => {
  const [hidden, setHidden] = useState(true);
  console.log("FunctionCall", evt)
  const msg = evt.Content.parts[0]
  const fn = msg.functionCall.name as string

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

  const fnArgs = f2NameArgs[fn]
  console.log("funcCall:", fn, msg, fnArgs)

  return (
    <div className="flex flex-col text-xs mt-2 mr-16 py-1 px-2 bg-green-700/80 rounded-t">
        <div className="flex gap-2 items-baseline">
        <UsageInfo evt={evt} size={12} />
        <Braces size={12} onClick={(e) => {
          e.preventDefault()
          e.stopPropagation()
          setHidden(!hidden)
        }} />
        <div className="ml-auto">
        { fnArgs ?
          <NameArgCall name={fn} args={fnArgs}/>
          :
          <NameArgCall name={fn} args={[]}/>
        }
        </div>
      </div>
      <JsonInfo hidden={hidden} data={evt} />
    </div>
  )
}

const NameArgCall = ({name, args}:{name: string, args: string[]}) => {
  return (
    <div className="flex gap-2 overflow-x">
      <div className="font-bold">{name}</div>
      <div className="monospace text-xs">{args.join("")}</div>
    </div>
  )
}

export const FunctionResp = ({evt}:{evt: any}) => {
  const [hiddenDetails, setHiddenDetails] = useState(true);
  const [hiddenJson, setHiddenJson] = useState(true);
  // console.log("FunctionResp", evt)
  const msg = evt.Content.parts[0]
  const r = msg.functionResponse.response;

  if (r?.error && r?.error !== "") {
    return <FunctionRespError evt={evt} />
  }

  const fn = msg.functionResponse.name

  return (
    <div className="flex flex-col text-xs mb-2 mr-16 py-1 px-2 bg-green-900/80 rounded-b">
      <div className="flex justify-between" onClick={() => setHiddenDetails(!hiddenDetails)}>
        <div className="flex gap-2 items-baseline">
          <div className="font-bold">{fn}</div>
          <Braces size={14} onClick={(e) => {
            e.preventDefault()
            e.stopPropagation()
            setHiddenJson(!hiddenJson)
          }} />
          { fn === "write_file" && <FileDiff size={14} onClick={(e) => {
            e.preventDefault()
            e.stopPropagation()
            console.log("diff.write_file")
            // todo, send some event 
          }}/>}

        </div>
        <TimeInfo timestamp={evt.Timestamp} />
      </div>
      { !hiddenDetails && 
        <div className="">
          { fn === "read_file" && <ReadFileResp evt={evt}/> }
          { fn === "read_dir" && <ReadDirResp evt={evt}/> }
          { fn === "tree_dir" && <TreeDirResp evt={evt}/> }
          { fn === "write_file" && <WriteFileResp evt={evt}/> }
        </div>
      }
      <JsonInfo hidden={hiddenJson} data={evt} />
    </div>
  )
}

const FunctionRespError = ({evt}:{evt: any}) => {
  console.log("FunctionRespError")
  const [hidden, setHidden] = useState(true);
  const msg = evt.Content.parts[0]

  return (
    <div className="flex flex-col text-xs mb-2 mr-16 py-1 px-2 bg-red-700/50 rounded-b-lg">
      <div className="flex justify-between" onClick={() => setHidden(!hidden)}>
        <div className="font-bold">{msg.functionResponse.name}</div>
        <TimeInfo timestamp={evt.Timestamp} />
      </div>
      <JsonInfo hidden={hidden} data={evt} />
    </div>
  )
}

export const ReadFileRespHighlight = ({evt}:{evt: any}) => {
  console.log("ReadFileResp", evt)

  const msg = evt.Content.parts[0]
  const source = msg.response.content

  return (
      <div className={cn(
        "flex flex-col relative",
        "max-h-64",
      )}>
        <CopyButton source={source} />
        <SyntaxHighlighter
          PreTag="div"
          children={String(source).replace(/\n$/, '')}
          language={"tsx"}
          style={vscDarkPlus}
          // className="bg-stone-800"
          codeTagProps={{
            className: "not-prose bg-[#1e1e1e]"
          }}
          customStyle={{
            lineHeight: "1.3"
          }}
        />
      </div>
  )
}

export const ReadFileResp = ({evt}:{evt: any}) => {
  console.log("ReadFileResp", evt)

  const msg = evt.Content.parts[0]
  const source = msg.functionResponse?.response?.content

  return (
    <pre className="m-2 p-2 max-h-64 overflow-auto bg-slate-800 text-xs">
      {source}
    </pre>
  )
}

export const ReadDirResp = ({evt}:{evt: any}) => {
  console.log("ReadDirResp", evt)
  const msg = evt.Content.parts[0]

  return (
    <pre className="m-2 p-2 max-h-64 overflow-auto bg-slate-800 text-xs">
      {msg.functionResponse.response.listing}
    </pre>
  )
}

export const TreeDirResp = ({evt}:{evt: any}) => {
  console.log("TreeDirResp", evt)
  const msg = evt.Content.parts[0]

  return (
    <pre className="m-2 p-2 max-h-64 overflow-auto bg-slate-800 text-xs">
      {msg.functionResponse.response.listing}
    </pre>
  )
}

export const WriteFileResp = ({evt}:{evt: any}) => {
  console.log("WriteFileResp", evt)
  const msg = evt.Content.parts[0]
  const fp = msg.functionResponse?.response?.path
  const content = evt.Actions.StateDelta.fs[fp]

  return (
    <pre className="m-2 p-2 max-h-64 overflow-auto bg-slate-800 text-xs">
      {content}
    </pre>
  )
}

// const showDiff = () => {
//   vscodeApi.getState()
// }


