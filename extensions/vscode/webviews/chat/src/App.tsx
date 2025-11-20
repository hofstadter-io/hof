import { useState, useEffect, useCallback, useRef } from 'react'
import { vscodeApi } from './vscodeApi.js'
import './index.css' // We'll add some styles

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import {
  Event,
  UsageInfo,
} from './Messages';

// Define the message types we expect
// (These should match your Go server and extension)
interface IdPayload {
  id: string
}

interface ServerMessage {
  type: string;
  payload: unknown;
}

function App() {
  const state = vscodeApi.getState() || {}
  const [sid, setSid] = useState(state?.sid || '');
  const [session, setSession] = useState<any>(state?.session || {})
  const [chatState, setChatState] = useState<any>(state?.chatState || {})
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Scroll to bottom when chat log changes
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [session?.events]);

  // update our listener when the sid changes
  useEffect(() => {
    // The 'onMessage' helper returns a cleanup function
    const removeListener = vscodeApi.onMessage((event) => {
      // console.log("Chat.sid event:", event)
      const message = event.data as ServerMessage;
      // console.log("Chat.sid message:", message)
      
      if (message.type === "session.get.resp" || message.type === "session.info") {
        console.log("Chat session.info", sid, message.payload)
        const payload = message.payload as IdPayload;
        if (payload.id === sid) {
          vscodeApi.setState({
            ...state,
            session: payload,
          })
          setSession(message.payload as any)
        }
      }
    });

    // Return the cleanup function
    return removeListener;
  }, [sid]); // Empty dependency array means this runs once

  // Effect to listen for messages from the extension
  useEffect(() => {
    // The 'onMessage' helper returns a cleanup function
    const removeListener = vscodeApi.onMessage((event) => {
      // console.log("Chat message:", event)
      const message = event.data as ServerMessage;

      if (message.type === 'agents.list.resp') {
        setChatState((prev: any) => {
          const next = {
            ...prev,
            agents: message.payload
          }
          vscodeApi.setState({
            ...state,
            chatState: next,
          })
          return next
        })
      }
      
      // Handle the message based on its type
      if (message.type === 'chat.event') {
        // console.log("Chat chat.event", message.payload)
        setSession((prev: any) => {
          const next = {
            ...prev,
            events: [...(prev.events), message.payload],
          }
          vscodeApi.setState({
            ...state,
            session: next,
          })
          return next
        });
      }

      if (message.type === "chat.loadSession") {
        const payload = message.payload as IdPayload;

        // update core state
        vscodeApi.setState({
          ...state,
          sid: payload.id
        })
        setSid(payload.id)

        // request session info when we load
        vscodeApi.postMessage({
          type: 'session.get',
          payload: {
            id: payload.id,
          }
        });
      }
    });

    // console.log("initial fetch", state?.sid, state?.sid !== "")
    if (state?.sid !== "") {
      vscodeApi.postMessage({
        type: 'session.get',
        payload: {
          id: state.sid,
        }
      });
    }

    // scroll immediately if we already have history
    messagesEndRef.current?.scrollIntoView({ behavior: "instant" });

    // Return the cleanup function
    return removeListener;
  }, []); // Empty dependency array means this runs once

  const handleInput = (input: string) => {
    setChatState((prev: any) => {
      const next = {
        ...prev,
        input,
        error: undefined,
      }
      vscodeApi.setState({
        ...state,
        chatState: next,
      })
      return next
    })
  }

  const handleSelect = (input: string) => {
    setChatState((prev: any) => {
      const next = {
        ...prev,
        agent: input,
      }
      vscodeApi.setState({
        ...state,
        chatState: next,
      })
      return next
    })
  }

  // Handler for sending a message
  const handleSend = useCallback(() => {
    if (chatState.input.trim()) {
      console.log("CHAT!", chatState)

      const input = chatState.input as string
      if (input.startsWith("/")) {
        const parts = input.split(/[,\s]/)
        const cmd = parts[0]
        const args = parts.splice(1)
        console.log("CMD!", cmd, args)

        if (cmd === "/state") {
          if (args.length === 1) {
            vscodeApi.postMessage({
              type: 'session.state.get',
              payload: {
                sid,
                key: args[0],
              }
            });

          }
          if (args.length > 1) {
            vscodeApi.postMessage({
              type: 'session.state.put',
              payload: {
                sid,
                key: args[0],
                val: args[1],
              }
            });

            setSession((prev: any) => {
              const next = {
                ...prev,
                state: {
                  ...prev?.state
                }
              }
              next.state[parts[0]] = parts[1]
              return next
            });
          }
          return
        }

        setChatState((prev: any) => {
          return {
            ...prev,
            error: "unknown command: " + cmd
          }
        })
        return
      }


      // 1. Post message to the extension
      // This will be caught by `sidebar.ts`
      vscodeApi.postMessage({
        type: 'chat.userMessage',
        payload: {
          text: chatState.input,
          sid,
          agent: chatState.agent || "basic-fast"
        }
      });

      setSession((prev: any) => {
        return {
          ...prev,
          events: [...(prev?.events || []), {
            Content: {
              role: "user",
              parts: [{
                text: chatState.input
              }]
            },
            Timestamp: new Date().toISOString()
          }],
        }
      });

      // 3. Clear the input
      handleInput('');
    }
  }, [chatState.input, chatState.agent]);

  // console.log("chat.session", session)
  // console.log("chat.chatState", chatState)

  const usage: any = {
    candidatesTokenCount: 0,
    promptTokenCount: 0,
    cachedContentTokenCount: 0,
    thoughtsTokenCount: 0,
    totalTokenCount: 0,
  }
  if (!!(session?.events) && session.events.length > 0) {
    session?.events?.forEach( (e: any) => {
      if (e?.UsageMetadata) {
        const u = e.UsageMetadata
        usage.candidatesTokenCount += u.candidatesTokenCount || 0
        usage.cachedContentTokenCount += u.cachedContentTokenCount || 0
        usage.promptTokenCount += u.promptTokenCount || 0
        usage.thoughtsTokenCount += u.thoughtsTokenCount || 0
        usage.totalTokenCount += u.totalTokenCount || 0
      }
    })
  }

  return (
    <div className="flex flex-col p-2">
      <Header sid={sid} usage={usage} />

      <Events
        session={session}
        messagesEndRef={messagesEndRef}
      />

      <UserInputs 
        sid={sid}
        usage={usage}
        chatState={chatState}
        handleInput={handleInput}
        handleSelect={handleSelect}
        handleSend={handleSend}
      />

    </div>
  )
}

const Header = ({
  sid,
  usage,
}:{
  sid: string,
  usage: any,
}) => {
  return (
    <div className="m-4 p-3 border-b text-md font-thin flex gap-2 justify-between">
      <span>{sid}</span>
      <UsageInfo usage={usage} size={16}/>
    </div>
  )
}

const Events = ({
  session,
  messagesEndRef
}:{
  session: any,
  messagesEndRef: any,
}) => {
  if (!session?.events?.length) {
    return null
  }
  return (
    <div className="flex-grow mx-2 overflow-y-auto">
      {session?.events?.map((e: any) => {
        if (!(e?.Content)) {
          return null
        }
        return (
          <Event key={e.ID} evt={e}/>
        )
      })}
      <div ref={messagesEndRef} />
    </div>
  )
}

const UserInputs = ({
  sid,
  usage,
  chatState,
  handleInput,
  handleSelect,
  handleSend,
}:{
  sid: string,
  usage: any,
  chatState: any,
  handleInput: any,
  handleSelect: any,
  handleSend: any,
}) => {
  // console.log("chat.input.chatState", chatState)

  return (
    <div className="flex flex-col bg-slate-800/80 mx-3 px-3 gap-2 mt-2 rounded-lg">
      <div className="py-3 px-1 text-md flex gap-2 justify-between">
        <span>{sid}</span>
        <UsageInfo usage={usage} size={16}/>
      </div>
      <textarea
        className="w-full rounded p-2 text-md mt-2 bg-slate-700/80"
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
      <div className="flex gap-2 justify-evenly">
        <AgentSelect chatState={chatState} handleSelect={handleSelect} />
        <button 
          onClick={handleSend}
          className="w-40 p-2 rounded border border-white bg-slate-800 hover:bg-slate-600"
        >Send</button>
      </div>
    </div>
  )
}

const AgentSelect = ({chatState, handleSelect}: {chatState: any, handleSelect:(s: string)=>void}) => {
  const agents: Record<string,any> = {}
  // for (const [key, _] of Object.entries(chatState?.agents)) {
  //   const parts = key.split('-')
  //   const k = parts[0]
  //   var curr: any[] = agents[k]
  //   if (!curr) {
  //     curr = []
  //   }
  //   if (parts.length > 1) {
  //     const v = parts.splice(1).join(" ")
  //     curr.push(v)
  //   } 
  //   agents[k] = curr
  // }
  agents["coding"] = ["lite", "fast", "norm", "hard"]
  agents["basic"] = ["lite", "fast", "norm", "hard"]
  agents["general"] = ["lite", "fast", "norm", "hard"]
  agents["filesys"] = ["lite", "fast", "norm", "hard"]
  // console.log("chat.input.agents", agents)

  return (
    <Select 
      defaultValue={chatState.agent || "general-fast"}
      onValueChange={(v: string) => {
        handleSelect(v)
      }}
    >
      <SelectTrigger className="w-40">
        <SelectValue placeholder="Select an agent" />
      </SelectTrigger>
      <SelectContent>
        {[...Object.entries(agents)].map(([key, value]: [string, any]) => {
          if (value.length > 0) {
            return (
              <SelectGroup>
                <SelectLabel
                  className="text-lg text-gray-900"
                >{key}</SelectLabel>

                {value.map((v: any) => {
                  const val = `${key}-${v}`
                  return (
                    <SelectItem key={val} value={val}
                      className="ml-2 p-1 text-sm font-thin text-gray-800"
                    >{val}</SelectItem>
                  )
                })}
              </SelectGroup>
            )
          } else {
            return (
              <SelectItem key={key} value={key}
                className="ml-2 p-1 text-sm font-thin text-gray-800"
              >{key}</SelectItem>
            )
          }
        })}
      </SelectContent>
    </Select>
  )
}


export default App