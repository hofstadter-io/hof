import { useRef, useState } from 'react'
import { Sparklines, SparklinesLine, SparklinesSpots } from 'react-sparklines';

import { vscodeApi } from '@/vscodeApi.js'

import { Badge } from "@/components/ui/badge"
import { cn, processEvents } from "@/lib/utils"

import { Header } from "../Header"
import { ChatEditor } from './editor'
import { Bot, BotMessageSquare, Drama, FileCode, Megaphone, ScrollText, TerminalSquare } from 'lucide-react'
import { useChat } from '@/hooks/useChat'
import { ToolTipper } from 'veg-webview-common'

export const UserInput = () => {
  const {
    sid,
    setPos,
    usage,
    session,
    chatState,
    diff,
    handleSend,
  } = useChat();

  const cacheKeys = Object.keys(session?.state || {}).filter(k => k.startsWith("cache:"));
  const fileKeys = Object.keys(session?.state || {}).filter(k => k.startsWith("files:"));
  const agentmdKeys = Object.keys(session?.state || {}).filter(k => k.startsWith("agentmd:"));

  const s = vscodeApi.getState()
  // console.log("chat.state", s)
  const [userInput, setUserInput] = useState<any>({ 
    agent: s?.userInput?.agent || "",
    model: s?.userInput?.model || "",
    environ: s?.userInput?.environ || "",
    text:  s?.userInput?.text || "",
  })

  const editorRef = useRef<any>(null);
  const inputReady: boolean = (userInput?.agent !== "" && 
                               userInput?.model !== "" &&
                               userInput?.text  !== "" )

  const handleInputUpdate = ({ editor }:{ editor: any }) => {
    // console.log("handleInputUpdate", editor)
    const markdown = editor.getMarkdown()
    // console.log(markdown)

    setUserInput((prev: any) => {
      const next = {
        ...prev,
        text: markdown,
      }
      const s = vscodeApi.getState()
      vscodeApi.setState({
        ...s,
        userInput: next,
      })
      return next
    })
  }

  const handleSelectModel = (input: string) => {
    console.log("setting model:", input)
    setUserInput((prev: any) => {
      const next = {
        ...prev,
        model: input,
      }
      const s = vscodeApi.getState()
      const n = {
        ...s,
        userInput: next,
      }
      console.log("setting userInput.model:", input, prev, next, s, n)
      vscodeApi.setState(n)
      return next
    })
  }

  const handleSelectAgent = (input: string) => {
    setUserInput((prev: any) => {
      const s = vscodeApi.getState()
      const next = {
        ...prev,
        agent: input,
      }
      vscodeApi.setState({
        ...s,
        userInput: next,
      })
      return next
    })
  }

  const handleSelectEnviron = (input: string) => {
    setUserInput((prev: any) => {
      const s = vscodeApi.getState()
      const next = {
        ...prev,
        environ: input,
      }
      vscodeApi.setState({
        ...s,
        userInput: next,
      })
      return next
    })
  }

  const doSend = () => {
    console.log("doSend", userInput)
    // what is the userInput?
    // single line, starting with special char, or maybe a few?
    // we are quickly going towards a parser here...
    const text = (userInput.text as string).trim()
    if (text.length === 0) {
      // no input
      return
    }
    const lines = text.split("\n")
    // is this a special thing
    if (lines.length === 1 && lines[0].startsWith("[@")) {
      const line = lines[0] as string
      const parts = line.split("[")
      parts.forEach((part, p) => {
        // probably the first part
        if (part === "") {
          return
        }
        // trim endings: [@ ... 
        part = part.substring(2).trim()
        const pairs = part.split(/\s+/)
        console.log("handling special:", part, pairs)
        var char = "@"
        var rest = ""
        var extra: string[] = []
        pairs.forEach((p) => {
          if (p.endsWith("]")) {
            p = p.substring(0, p.length-1)
          }
          if (p.startsWith("id")) {
            // trim endings: id="..."
            p = p.trim().substring(4, p.length-1).trim()
            rest = p
          } else if (p.startsWith("char")) {
            // trim endings: char="."
            p = p.trim().substring(6, p.length-1).trim()
            char = p[0]
          } else {
            extra.push(p)
          }
        })

        rest = rest.trim()
        console.log("parsed:", char, rest, extra)

        switch (char) {
          case "@":
            if (Object.keys(chatState?.config?.agents).includes(rest)) {
              handleSelectAgent(rest)
            }
            if (Object.keys(chatState?.config?.models).includes(rest)) {
              handleSelectModel(rest)
            }
            break;

          case ">":
            handleSelectEnviron(rest)
            break;

          case "#":
            break;

          case "$":
            console.log(`$${rest}:`, extra)
            if (rest === "state") {
              if (!extra || extra.length < 1) {
                // todo, let user know by showing a help message
                return
              }
              const key = extra[0]
              var val: any = null
              if (extra.length > 1) {
                val = extra.splice(1).join(" ")
              }

              console.log("state!", key, val)
              if (!val || val.trim() === "") {
                vscodeApi.postMessage({
                  type: 'session.state.del',
                  payload: {
                    sid,
                    key,
                  },
                });
              } else {
                vscodeApi.postMessage({
                  type: 'session.state.put',
                  payload: {
                    sid,
                    key,
                    val,
                  },
                });
              }
              vscodeApi.postMessage({
                type: 'session.get',
                payload: {
                  sid,
                },
              });


              // we can't process any more
              return
            }
            break;
        }
      }) // end of part loop

      // end of our special char handling, we should return and not send a message
      editorRef?.current?.commands.clearContent()
      return 
    }

    // otherwise assume a message for the agent
    // TODO, facet detection and settings updates from longer text
    // a user could want to delegate different parts of the task to different agents in their message

    handleSend(userInput);
    setUserInput({
      ...userInput,
      text: "",
    })
    console.log("doSend.clear", editorRef.current)
    editorRef?.current?.commands.clearContent()
  }

  const { usages } = processEvents(session?.events)
  var cached: any[] = []
  var prompt: any[] = []
  var inputs: any[] = []
  var thinks: any[] = []
  var writes: any[] = []
  var output: any[] = []
  var totals: any[] = []

  usages?.forEach((u) => {
    const c = u?.cachedContentTokenCount || 0;
    const p = u?.promptTokenCount || 0;
    const i = p - c
    const t = u?.thoughtsTokenCount || 0;
    const w = u?.candidatesTokenCount || 0;
    const o = t + w
    const T = u?.totalTokenCount || 0;
    cached.push(c)
    prompt.push(p)
    inputs.push(i)
    thinks.push(t)
    writes.push(w)
    output.push(o)
    totals.push(T)
  })


  return (
    <div 
      className={cn(
      "h-full flex flex-col gap-2",
      // "bg-slate-800/80 border-gray-500",
      )}
    >
      <Header />

      <div className="flex gap-1 items-center">

        {/* User Settings */}
        { userInput?.agent && <Badge className="text-sky-300 bg-sky-600/50 "><Bot size={12}/>{userInput?.agent}</Badge>}
        { userInput?.model && <Badge className="text-sky-300 bg-sky-600/50 "><Drama size={12}/>{userInput?.model}</Badge>}
        { userInput?.environ && <Badge className="text-lime-300/80 bg-lime-600/50 "><TerminalSquare size={12}/>{userInput?.environ}</Badge>}

        {/* Context Info */}
        {cacheKeys.length > 0 && (
          <ToolTipper label={cacheKeys.map(k => k.split(':').slice(2).join(':')).join('\n')}>
            <Badge className="text-fuchsia-200/80 bg-fuchsia-600/50 flex gap-1 items-center px-2">
              <ScrollText size={12}/>
              <span>{cacheKeys.length}</span>
            </Badge>
          </ToolTipper>
        )}
        {fileKeys.length > 0 && (
          <ToolTipper label={fileKeys.map(k => k.split(':').slice(2).join(':')).join('\n')}>
            <Badge className="text-violet-300 bg-violet-600/50 flex gap-1 items-center px-2">
              <FileCode size={12}/>
              <span>{fileKeys.length}</span>
            </Badge>
          </ToolTipper>
        )}
        {agentmdKeys.length > 0 && (
          <ToolTipper label={agentmdKeys.map(k => k.split(':').slice(2).join(':')).join('\n')}>
            <Badge className="text-amber-200 bg-amber-600/50 flex gap-1 items-center px-2">
              <Megaphone size={12}/>
              <span>{agentmdKeys.length}</span>
            </Badge>
          </ToolTipper>
        )}

        {/* Token Usage */}
        <div className="w-100 ml-4 px-2 h-6 flex relative rounded border-b border-dashed border-gray-400">
          <div className="absolute top-0 left-0 h-3 w-full border-t border-dashed border-red-400 z-20">
          </div>
          <div className="absolute top-0 left-0 h-3 w-full border-b border-dashed border-amber-400 z-20">
          </div>
          <div className="absolute top-[-3px] left-0 h-6 w-50 z-30">
            <Sparklines data={cached} width={140} height={20} min={0} max={100000}>
              <SparklinesLine style={{ stroke: "oklch(84.1% 0.238 128.85)", fill: "oklch(84.1% 0.238 128.85)" }} />
            </Sparklines>
          </div>
          <div className="absolute top-[-3px] left-0 h-6 w-50 z-30">
            <Sparklines data={prompt} width={140} height={20} min={0} max={100000}>
              <SparklinesLine style={{ stroke: "oklch(87.9% 0.169 91.605)", fill: "oklch(87.9% 0.169 91.605)" }} />
            </Sparklines>
          </div>

          <div className="absolute top-[-3px] left-50 h-6 w-50">
            <Sparklines data={output} width={140} height={20} min={0} max={100000}>
              <SparklinesLine style={{ stroke: "oklch(74.6% 0.16 232.661)", fill: "oklch(74.6% 0.16 232.661)" }} />
            </Sparklines>
          </div>
          <div className="absolute top-[-3px] left-50 h-6 w-50">
            <Sparklines data={totals} width={140} height={20} min={0} max={100000}>
              <SparklinesLine style={{ stroke: "oklch(74% 0.238 322.16)", fill: "oklch(74% 0.238 322.16)" }} />
            </Sparklines>
          </div>
        </div>

      </div>

      { userInput?.error && <span
        className="p-2 w-full border rounded bg-red-800 font-heavy"
      >{userInput.error}</span>}

      <div
        className="flex flex-col relative" 
        // need to capture all keyboard events for special case
        // tiptap isn't letting us do CMD + Enter to send easily
        onKeyDown={(evt: any)=>{
          if (evt.metaKey && evt.key === "Enter") {
            console.log("Send", evt)
            doSend()
          }
        }}
      >

        <BotMessageSquare
          size={32}
          strokeWidth={1.5}
          className="z-50 absolute rounded-4xl right-2 top-2 p-2 text-white/50 bg-sky-500/50 hover:text-white hover:bg-sky-500"
          onClick={() => doSend()}
        />

        <ChatEditor
          userInput={userInput}
          handlers={{
            handleInputUpdate,
            handleSelectAgent,
            handleSelectModel,
            handleSelectEnviron,
          }}
          editorRef={editorRef}
        />
      </div>
    </div>
  )
}
