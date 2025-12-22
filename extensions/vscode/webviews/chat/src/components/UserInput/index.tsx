import { useRef, useState } from 'react'

import { vscodeApi } from '@/vscodeApi.js'

import { Badge } from "@/components/ui/badge"
import { cn, processEvents } from "@/lib/utils"

import { Header } from "../Header"
import { ChatEditor } from './editor'
import { BotMessageSquare } from 'lucide-react'
import { useChat } from '@/hooks/useChat'
import { SessionSparklines } from '../SessionSparklines';
import { ChatStatePills } from '../ChatStatePills';
import { handleChatboxCommand } from '@/lib/chatboxCommandHandlers';

export const UserInput = () => {
  const {
    sid,
    setPos,
    usage,
    session,
    chatState,
    setChatState,
    diff,
    handleSend,
  } = useChat();

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
      
      // Update chatState so it's available for the Header Menu
      setChatState((c: any) => ({ ...c, userInput: next }));

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
      
      // Update chatState
      setChatState((c: any) => ({ ...c, userInput: next }));

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
      
      // Update chatState
      setChatState((c: any) => ({ ...c, userInput: next }));

      return next
    })
  }

  const handleSelectEnviron = (input: string) => {
    if (input === "none") {
      input = ""
    }
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

      // Update chatState
      setChatState((c: any) => ({ ...c, userInput: next }));

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

    const handled = handleChatboxCommand(text, sid, chatState, {
      handleSelectAgent,
      handleSelectModel,
      handleSelectEnviron,
    });

    if (handled) {
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

  return (
    <div 
      className={cn(
      "h-full flex flex-col gap-2",
      // "bg-slate-800/80 border-gray-500",
      )}
    >
      <Header />

      <div className="flex gap-1 items-center">
        <ChatStatePills userInput={userInput} session={session} />
        <SessionSparklines events={session?.events} session={session} chatState={chatState} />
      </div>

      { userInput?.error && <span
        className="p-2 w-full border rounded bg-red-800 font-heavy"
      >{userInput.error}</span>}

      <div
        className="flex flex-col relative" 
        // need to capture all keyboard events for special case
        // tiptap isn't letting us do CMD + Enter to send easily
        onKeyDown={(evt: any)=>{
          // TODO, need to make sure this does not make it into the text
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
