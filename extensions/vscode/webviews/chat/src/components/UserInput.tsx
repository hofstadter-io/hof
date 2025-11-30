import { useState } from 'react'
import { vscodeApi } from '@/vscodeApi.js'

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import { cn } from "@/lib/utils"

import { Header } from "./Header"

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
    agent: chatState?.agent || "general_assistant",
    model: chatState?.model || "default",
    text:  chatState?.input || "",
  })

  const inputReady: boolean = (userInput?.text as string).startsWith("/") ||
                              (userInput?.agent !== "" && 
                               userInput?.model !== "" &&
                               userInput?.text  !== "" )

  const handleInput = (input: string) => {
    // console.log("handleInput")
    setUserInput((prev: any) => {
      const next = {
        ...prev,
        text: input,
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

  const doSend = (input: any) => {
    handleSend(input);
    setUserInput({
      ...input,
      text: "",
    })
  }


  return (
    <div 
      className={cn(
      "flex flex-col gap-2 mx-2",
      "bg-slate-800/80 rounded-xl border-gray-500",
      )}
    >
      <Header
        sid={sid}
        setPos={setPos}
        diff={diff}
        usage={usage}
        session={session}
        chatState={chatState}
      />

      <textarea
        className="rounded-lg p-2 text-md m-2 bg-slate-700/80"
        rows={5}
        value={userInput.text}
        onChange={(e) => {
          // console.log("textarea", e)
          handleInput(e.target.value)
        }}
        onKeyDown={(e) => {
          if (e.metaKey && e.key === 'Enter') {
            doSend(userInput);
          }
        }}
        placeholder="Type a message..."
      />

      { userInput?.error && <span
        className="p-2 w-full border rounded bg-red-800 font-heavy"
      >{userInput.error}</span>}

      <div className="flex gap-2 m-3">
        <AgentSelect agent={userInput?.agent} agents={chatState?.config?.agents} handleSelect={handleSelectAgent} />
        <ModelSelect model={userInput?.model} models={chatState?.config?.models} handleSelect={handleSelectModel} />
        <button 
          disabled={!inputReady}
          onClick={() => doSend(userInput)}
          className={cn(
            "w-32 p-2 rounded-lg border",
            inputReady
            ? "border-2 border-sky-500 bg-slate-700 hover:bg-slate-600"
            : "border-gray-600 bg-slate-800",
          )}
        >Send</button>
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
      defaultValue={model || "gemini-2.5-flash"}
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