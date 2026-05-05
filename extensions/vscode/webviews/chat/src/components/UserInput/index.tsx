import { useEffect, useRef, useState } from 'react'
import { BotMessageSquare } from 'lucide-react'

import { vscodeApi } from '@/vscodeApi.js'

import { useChat } from '@/hooks/useChat'
import { cn } from "@/lib/utils"

import { Header } from "../Header"
import { ChatStatePills } from '../ChatStatePills';
import { SessionSparklines } from '../SessionSparklines';
import { handleChatboxCommand } from './chatboxCommandHandlers';
import { ChatEditor } from './editor'

export const UserInput = () => {
  const {
    sid,
    session,
    chatState,
    setChatState,
    handleSend,
  } = useChat();

  const s = vscodeApi.getState()
  // console.log("chat.state", s)
  const [userInput, setUserInput] = useState<any>({ 
    sid: s?.userInput?.sid || "",
    agent: s?.userInput?.agent || "",
    model: s?.userInput?.model || "",
    environ: s?.userInput?.environ || "",
    text:  s?.userInput?.text || "",
  })
  const [working, setWorking] = useState<boolean>(false)
  const editorRef = useRef<any>(null);

  useEffect(() => {
    const extra = session?.state ? {
      agent: session.state?.agent,
      model: session.state?.model,
      environ: session.state?.envName,
    } : {
      agent: userInput.agent,
      model: userInput.model,
      environ: userInput.environ,
    }

    // construct next state
    const next = {
      ...userInput,
      sid: session.sid,
      ...extra,
    };

    // has the machine finished?

    // update state in react
    setUserInput(next);    
    setChatState((c: any) => ({ ...c, userInput: next }));
    // persist state to vscode
    const s = vscodeApi.getState();
    vscodeApi.setState({ ...s, userInput: next });

  }, [session?.sid, session?.state?.model, session?.state?.agent, session?.state?.environ]);

  useEffect(() =>{
    if (session.events && session.events.length > 1) {
      const last = session.events[session.events.length-1] as any
      if (!!last) {
        console.log("considering last message...", last, last["TurnComplete"])
        if (last["TurnComplete"] === true) {
          setWorking(false)
        }
      }
    }
  }, [session?.events])

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
    setUserInput((prev: any) => {
      if (input === "none") {
        input = ""
      }
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
    setWorking(true);
    setUserInput({
      ...userInput,
      text: "",
    })
    console.log("doSend.clear", editorRef.current)
    editorRef?.current?.commands.clearContent()
  }

  const doStop = () => {
    console.log("cancel!", sid)
    vscodeApi.postMessage({
      type: 'session.cancel',
      payload: { sid },
    });
  }

  return (
    <div 
      className={cn(
      "h-full flex flex-col gap-2 py-2",
      // "bg-slate-800/80 border-gray-500",
      )}
    >
      <Header userInput={userInput}/>

      <div className="flex flex-col lg:flex-row gap-2">
        <div className="max-w-150">
          <SessionSparklines events={session?.events} />
        </div>
        <div className="flex gap-1 items-center">
          <ChatStatePills userInput={userInput} session={session} />
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
          className={cn(
            "z-50 absolute rounded-4xl right-2 top-2 p-2",
            "text-white/50 hover:text-white",
            working ? "bg-red-500/50 hover:bg-red-500" : "bg-sky-500/50 hover:bg-sky-500",
          )}
          onClick={() => working ? doStop() : doSend()}
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
