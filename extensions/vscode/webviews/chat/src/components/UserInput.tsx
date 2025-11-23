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
  usage,
  session,
  chatState,
  handleInput,
  handleSelectModel,
  handleSelectAgent,
  handleSend,
}:{
  sid: string,
  usage: any,
  session: any,
  chatState: any,
  handleInput: (input: string) => void,
  handleSelectModel: (input: string) => void,
  handleSelectAgent: (input: string) => void,
  handleSend: () => void,
}) => {

  return (
    <div 
      className={cn(
      "my-auto flex flex-col gap-2 mx-2",
      "bg-slate-800/80 rounded-xl border-gray-500",
      )}
    >
      <Header sid={sid} usage={usage} session={session} chatState={chatState} />

      <textarea
        className="rounded-lg p-2 text-md m-2 bg-slate-700/80"
        rows={5}
        value={chatState.input}
        onChange={(e) => handleInput(e.target.value)}
        onKeyDown={(e) => {
          if (e.metaKey && e.key === 'Enter') {
            handleSend();
          }
        }}
        placeholder="Type a message..."
      />

      { chatState?.error && <span
        className="p-2 w-full border rounded bg-red-800 font-heavy"
      >{chatState.error}</span>}

      <div className="flex gap-2 m-3">
        <AgentSelect chatState={chatState} handleSelect={handleSelectAgent} />
        <ModelSelect chatState={chatState} handleSelect={handleSelectModel} />
        <button 
          onClick={handleSend}
          className="w-32 p-2 rounded-lg border border-lime-600 bg-slate-800 hover:bg-slate-600"
        >Send</button>
      </div>
    </div>
  )
}

// @ts-ignore
const AgentSelect = ({chatState, handleSelect}: {chatState: any, handleSelect:(s: string)=>void}) => {
  const agents: string[] = []
  for (const [key, _] of Object.entries(chatState?.agents)) {
    agents.push(key)
  }
  console.log("chat.input.agents", agents)

  // const agents = ["coding", "coding-ro", "basic", "general", "filesys"]

  return (
    <Select 
      defaultValue={chatState?.agent}
      onValueChange={(v: string) => {
        handleSelect(v)
      }}
    >
      <SelectTrigger className="flex-3">
        <SelectValue placeholder="Select an agent" />
      </SelectTrigger>
      <SelectContent>
        {agents.map((v: any) => {
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
const ModelSelect = ({chatState, handleSelect}: {chatState: any, handleSelect:(s: string)=>void}) => {

  const models = [
    "gemini-2.5-flash-lite",
    "gemini-2.5-flash",
    "gemini-2.5-pro",
    "gemini-3.0-pro-preview",
  ]

  return (
    <Select 
      defaultValue={chatState?.model}
      onValueChange={(v: string) => {
        handleSelect(v)
      }}
    >
      <SelectTrigger className="flex-2">
        <SelectValue placeholder="Select a model" />
      </SelectTrigger>
      <SelectContent>
        {models.map((v: any) => {
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