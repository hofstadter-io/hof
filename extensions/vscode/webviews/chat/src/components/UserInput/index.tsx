import { useRef, useState } from 'react'
import { vscodeApi } from '@/vscodeApi.js'

import { Badge } from "@/components/ui/badge"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import { cn } from "@/lib/utils"

import { Header } from "../Header"
import { ChatEditor } from './editor'
import { AtSign, AudioLines, Bot, BotMessageSquare, Boxes, Drama, Forward, Send } from 'lucide-react'

export const UserInput = ({
  sid,
  setPos,
  usage,
  session,
  chatState,
  diff,
  handleSend,
}:{
  sid: string,
  setPos: any,
  usage: any,
  session: any,
  chatState: any,
  diff?: any,
  handleSend: (userInput: any) => void,
}) => {

  const [userInput, setUserInput] = useState<any>({ 
    agent: chatState?.agent || "veggie",
    model: chatState?.model || "default",
    environ: chatState?.environ || "debian",
    text:  chatState?.input || "",
  })

  const editorRef = useRef<any>(null);
  const inputReady: boolean = (userInput?.text as string).startsWith("/") ||
                              (userInput?.agent !== "" && 
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
    setUserInput((prev: any) => {
      const s = vscodeApi.getState()
      const next = {
        ...prev,
        model: input,
      }
      vscodeApi.setState({
        ...s,
        userInput: next,
        chatState: {
          ...s.chatState,
          model: input,
        },
      })
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
        chatState: {
          ...s.chatState,
          agent: input,
        },
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
        chatState: {
          ...s.chatState,
          environ: input,
        },
      })
      return next
    })
  }

  const doSend = () => {
    handleSend(userInput);
    setUserInput({
      ...userInput,
      text: "",
    })
    console.log("doSend.clear", editorRef.current)
    editorRef?.current?.commands.clearContent()
  }


  return (
    <div 
      className={cn(
      "h-full flex flex-col gap-3",
      // "bg-slate-800/80 border-gray-500",
      )}
    >
      <Header
        sid={sid}
        setPos={setPos}
        diff={diff}
        usage={usage}
        session={session}
        chatState={chatState}
        // className="border-t pt-2"
      />

      {/* <div className="flex gap-1">
        <Badge className="text-violet-300 bg-violet-700/30 "><TerminalSquare size={12}/>term-1</Badge>
        <Badge className="text-violet-300 bg-violet-700/30 "><FileBracesCorner size={12}/>example.go</Badge>
      </div> */}

      <div className="flex gap-1">
        { userInput?.agent && <Badge className="text-sky-300 bg-sky-700/30 "><Bot size={12}/>{userInput?.agent}</Badge>}
        { userInput?.model && <Badge className="text-sky-300 bg-sky-700/30 "><Drama size={12}/>{userInput?.model}</Badge>}
        { userInput?.environ && <Badge className="text-lime-400 bg-lime-800/30 "><Boxes size={12}/>{userInput?.environ}</Badge>}
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
          chatState={chatState}  
          onUpdate={handleInputUpdate}
          editorRef={editorRef}
        />
      </div>

      <div className="flex gap-2 m-3">
        <AgentSelect agent={userInput?.agent} agents={chatState?.config?.agents} handleSelect={handleSelectAgent} />
        <ModelSelect model={userInput?.model} models={chatState?.config?.models} handleSelect={handleSelectModel} />
        <EnvironSelect environ={userInput?.environ} environs={chatState?.config?.environs} handleSelect={handleSelectEnviron} />
      </div>

    </div>
  )
}

// @ts-ignore
const AgentSelect = ({agent, agents, handleSelect}: {agent: string, agents: any, handleSelect:(s: string)=>void}) => {
  const as: string[] = []
  if (agents) {
    for (const [key, _] of Object.entries(agents)) {
      as.push(key)
    }
  }
  // console.log("chat.input.agents", agents)

  // const agents = ["coding", "coding-ro", "basic", "general", "filesys"]

  return (
    <Select 
      defaultValue={agent}
      onValueChange={(v: string) => {
        handleSelect(v)
      }}
    >
      <SelectTrigger className="flex-3">
        <SelectValue placeholder="Select an agent" />
      </SelectTrigger>
      <SelectContent>
        {as.map((v: any) => {
          const val = `${v}`
          return (
            <SelectItem key={val} value={val}
              className="ml-2 p-1 text-sm font-thin text-gray-800"
            >{val}</SelectItem>
          )
        })}
      </SelectContent>
    </Select>
  )
}

// @ts-ignore
const ModelSelect = ({model, models, handleSelect}: {model: string, models: any, handleSelect:(s: string)=>void}) => {

  const ms: string[] = []
  if (models) {
    for (const [key, _] of Object.entries(models)) {
      ms.push(key)
    }
  }

  return (
    <Select 
      defaultValue={model}
      onValueChange={(v: string) => {
        handleSelect(v)
      }}
    >
      <SelectTrigger className="flex-2">
        <SelectValue placeholder="Select a model" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value={"default"}
          className="ml-2 p-1 text-sm font-thin text-gray-800"
        >agent-default</SelectItem>
        {ms.map((v: any) => {
          const val = `${v}`
          return (
            <SelectItem key={val} value={val}
              className="ml-2 p-1 text-sm font-thin text-gray-800"
            >{val}</SelectItem>
          )
        })}
      </SelectContent>
    </Select>
  )
}

const EnvironSelect = ({environ, environs, handleSelect}: {environ: string, environs: any, handleSelect:(s: string)=>void}) => {

  const ms: string[] = []
  if (environs) {
    for (const [key, _] of Object.entries(environs)) {
      ms.push(key)
    }
  }

  return (
    <Select 
      defaultValue={environ}
      onValueChange={(v: string) => {
        handleSelect(v)
      }}
    >
      <SelectTrigger className="flex-2">
        <SelectValue placeholder="Select a environ" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value={"default"}
          className="ml-2 p-1 text-sm font-thin text-gray-800"
        >agent-default</SelectItem>
        {ms.map((v: any) => {
          const val = `${v}`
          return (
            <SelectItem key={val} value={val}
              className="ml-2 p-1 text-sm font-thin text-gray-800"
            >{val}</SelectItem>
          )
        })}
      </SelectContent>
    </Select>
  )
}