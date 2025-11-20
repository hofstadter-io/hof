import { useState } from "react"
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import {Prism as SyntaxHighlighter} from 'react-syntax-highlighter'
import {vscDarkPlus} from 'react-syntax-highlighter/dist/esm/styles/prism'
import { CopyToClipboard } from "react-copy-to-clipboard-ts";
import { cn } from './util'

import {
  Code,
  ClipboardCopy,
  ClipboardCheck,
  Circle,
  CircleQuestionMark,
  CircleCheckBig,
  CircleDashed,
  OctagonAlert,
  Megaphone,
  BrainCircuit,
  EqualApproximately,
  Coins,
  DatabaseBackup,
} from 'lucide-react'

const Markdown = ({children}:{children: any}) => {
  return (
    <ReactMarkdown
      children={children}
      components={components}
      remarkPlugins={[[remarkGfm, {singleTilde: false}]]}
    />
  )
}

const components = {
  code(props: any) {
    const [show, setShow] = useState(true)
    const {children, className, node, ...rest} = props
    const match = /language-(\w+)/.exec(className || '') // || "\n" in children?
    return match ? (
      // code block?
      <div className={cn(
        "flex flex-col relative",
        show ? "" : "max-h-64"
      )}>
        <CopyButton source={children} />
        <SyntaxHighlighter
          {...rest}
          PreTag="div"
          children={String(children).replace(/\n$/, '')}
          language={match[1]}
          style={vscDarkPlus}
          // className="bg-stone-800"
          codeTagProps={{
            className: "not-prose bg-[#1e1e1e]"
          }}
          customStyle={{
            lineHeight: "1.3"
          }}
        />
        <div className="bg-slate-700 hover:bg-slate-500 w-full flex justify-center" onClick={() => setShow(!show)}>
          <Code size={20} strokeWidth={2} /> 
        </div>
      </div>
    ) : (
      // inline?
      <code {...rest} className={className}>
        {children}
      </code>
    )
  }
}

const CopyButton = ({source}:{source: string}) => {
  const [copied, setCopied] = useState(false)
  const common = "p-1"
  return (
    <CopyToClipboard text={source} onCopy={() => {
      setCopied(true)
      setTimeout(() => setCopied(false), 3000)
    }}>
      <span className="absolute z-50 top-2 right-1">
        { copied ?
          <ClipboardCheck size={32} strokeWidth={1} className={cn(common, "text-green-500")} /> 
        :   
          <ClipboardCopy size={32} strokeWidth={1} className={cn(common, "hover:text-blue-300")} /> 
        }
      </span>
    </CopyToClipboard>
  )
}


export const Event = ({evt}: {evt: any}) => {
  // console.log("session event:", evt)
  const msg = evt.Content.parts[0]

  if (evt.Content.role === "user") {
    if (msg.text) { return <UserMessage evt={evt}/> }
    if (msg.functionResponse) { return <FunctionResp evt={evt}/> }
  }

  if (evt.Content.role === "model") {
    // if (msg.text && evt.FinishReason === "STOP") { return <ModelMessage evt={evt}/> }
    // if (msg.text && evt.FinishReason === "") { return <ThinkingMessage evt={evt}/> }
    if (msg.functionCall) { return <FunctionCall evt={evt}/> }
    if (msg.text) { return <ModelMessage evt={evt}/> }
  }

  return (
    <div className="flex flex-col text-sm my-2 mx-4 px-2 py-1 bg-red-800">
      <div className="font-bold">unknown event type</div>
      <pre>
        {JSON.stringify(evt, null, "  ")}
      </pre>
    </div>
  )
}

export const TimeInfo = ({ evt }: { evt?: any }) => {
  if (!evt.Timestamp) {
    return null
  }
  const ts = evt.Timestamp
  const t = new Date(ts)
  return (
    <span className="ml-auto font-thin text-sm">
      {t.toLocaleString()}
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
        <CircleQuestionMark size={size}/>
        {u.candidatesTokenCount || "0"}
      </span>
      <span className="flex gap-1">
        <Megaphone size={size}/>
        {u.promptTokenCount || "0"}
      </span>
      <span className="flex gap-1">
        <BrainCircuit size={size}/>
        {u.thoughtsTokenCount || "0"}
      </span>
      <span>
        <EqualApproximately size={size}/>
      </span>
      <span className="flex gap-1">
        {u.totalTokenCount || "0"}
        <Coins size={size}/>
      </span>
      <span className="flex gap-1">
        {u.cachedContentTokenCount || "0"}
        <DatabaseBackup size={size}/>
      </span>
    </div>
  )
   
}

export const UserMessage = ({evt}:{evt: any}) => {
  return (
    <div className="ml-16 p-2 rounded-lg bg-slate-800/80 flex flex-col">
      <div className="p-2 dark:prose-invert prose-sm prose-stone">
        <Markdown>{evt.Content.parts.map((p: any) => p.text).join("\n\n")}</Markdown>
      </div>
      <TimeInfo evt={evt} />
    </div>
  )
}

export const ModelMessage = ({evt}:{evt: any}) => {
  return (
    <div className="mr-16 my-2 p-3 rounded-lg bg-stone-800/80 flex flex-col gap-1">
      <div className="p-3 dark:prose-invert prose-sm prose-stone">
        <Markdown>{evt.Content.parts.map((p: any) => p.text).join("\n\n")}</Markdown>
      </div>
      <span className="text-sm font-thin ml-auto m-0 p-0">{evt.Author}</span>
      <div className="flex gap-4">
        <UsageInfo evt={evt} size={16} />
        <MetaInfo evt={evt} size={12} />
        <TimeInfo evt={evt} />
      </div>
    </div>
  )
}

// @ts-ignore
export const ThinkingMessage = ({evt}:{evt: any}) => {
  const [hidden, setHidden] = useState(true);
  return (
    <div className="mr-16 my-2 p-3 rounded-lg bg-purple-800/80 flex flex-col gap-1">
      <div className="flex justify-between" onClick={() => setHidden(!hidden)}>
        <span hidden={!hidden} className="font-italic">thinking...</span>
        <div className="flex justify-between">
          <UsageInfo evt={evt} size={14} />
          <TimeInfo evt={evt} />
        </div>
      </div>
      <div hidden={hidden} className="p-3 bg-purple-900/80 dark:prose-invert prose-sm prose-stone">
        <Markdown>{evt.Content.parts.map((p: any) => p.text).join("\n\n")}</Markdown>
      </div>
    </div>
  )
}

export const FunctionCall = ({evt}:{evt: any}) => {
  const msg = evt.Content.parts[0]
  return (
    <div className="flex justify-between text-xs mt-2 mr-16 py-1 px-2 bg-green-700/80 gap-2 rounded-t-lg">
      <UsageInfo evt={evt} size={14} />
      <div className="flex gap-2 overflow-x-auto">
        <div className="font-bold">{msg.functionCall.name}</div>
        <div className="monospace text-xs">{JSON.stringify(msg.functionCall.args)}</div>
      </div>
    </div>
  )
}

export const FunctionResp = ({evt}:{evt: any}) => {
  const [hidden, setHidden] = useState(true);
  const msg = evt.Content.parts[0]
  const r = msg.functionResponse.response;
  if (r.error && r.error !== "") {
    return (
      <div className="flex flex-col text-xs mb-2 mr-16 py-1 px-2 bg-red-700/50 rounded-b-lg">
        <div className="flex justify-between" onClick={() => setHidden(!hidden)}>
          <div className="font-bold">{msg.functionResponse.name}</div>
          <TimeInfo evt={evt} />
        </div>
        { !hidden && 
          <div hidden={hidden} className="transition hidden:h-0 whitespace-pre-line ml-2 max-h-200 overflow-y-auto">
            {r.error}
          </div>
        }
      </div>
    )
  }
  return (
    <div className="flex flex-col text-xs mb-2 mr-16 py-1 px-2 bg-green-900/80 rounded-b-lg">
      <div className="flex justify-between" onClick={() => setHidden(!hidden)}>
        <div className="font-bold">{msg.functionResponse.name}</div>
        <TimeInfo evt={evt} />
      </div>
      { !hidden && 
        <div hidden={hidden} className="transition hidden:h-0 whitespace-pre-line ml-2 max-h-200 overflow-y-auto">
          {msg.functionResponse.response.listing}
        </div>
      }
    </div>
  )
}