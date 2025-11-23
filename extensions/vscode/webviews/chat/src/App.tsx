import { useState, useEffect, useCallback, useRef } from 'react'
import { vscodeApi } from './vscodeApi.js'
import './index.css' // We'll add some styles

// test comment to see what debug looks like

import { Header } from './components/Header.js'
import { Events } from './components/Messages.js';
import { UserInput } from './components/UserInput.js';

// Define the message types we expect
// (These should match your Go server and extension)
interface SidPayload {
  sid: string
}

interface ServerMessage {
  type: string;
  payload: unknown;
}

function App() {

  // TODO, this stuff should go into a custom hook/provider
  const state = vscodeApi.getState() || {}
  const [sid, setSid] = useState(state?.sid || '');
  const [session, setSession] = useState<any>(state?.session || {})
  const [chatState, setChatState] = useState<any>(state?.chatState || {})
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Scroll to bottom when chat log changes
  useEffect(() => {
    // todo, add configuration
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [session?.events]);

  // update our listener when the sid changes
  useEffect(() => {
    // The 'onMessage' helper returns a cleanup function
    const removeListener = vscodeApi.onMessage((event) => {
      const message = event.data as ServerMessage;
      
      if (message.type === "session.get.resp" || message.type === "session.info") {
        const payload = message.payload as SidPayload;
        if (payload.sid === sid) {
          vscodeApi.setState({
            ...state,
            session: payload,
          })
          setSession(message.payload as any)
        }
      }

      if (message.type === "session.delete") {
        const payload = message.payload as SidPayload;
        if (payload.sid === sid) {
          vscodeApi.setState({})
          setSid('')
          setSession({})
          setChatState({})
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
        const payload = message.payload as SidPayload;
        if (payload.sid) {
          // update core state
          vscodeApi.setState({
            ...state,
            sid: payload.sid
          })
          setSid(payload.sid)

          // request session info when we load
          vscodeApi.postMessage({
            type: 'session.get',
            payload: {
              sid: payload.sid,
            }
          });
        }
      }

    });

    // console.log("initial fetch", state?.sid, state?.sid !== "")
    if (state?.sid !== "") {
      vscodeApi.postMessage({
        type: 'session.get',
        payload: {
          sid: state.sid,
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

  const handleSelectModel = (input: string) => {
    setChatState((prev: any) => {
      const next = {
        ...prev,
        model: input,
      }
      vscodeApi.setState({
        ...state,
        chatState: next,
      })
      return next
    })
  }

  const handleSelectAgent = (input: string) => {
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
          // no arg means delete
          if (args.length === 1) {
            // add an empty element so we don't need extra logic below
            args.push("")
          }
          if (args.length > 1) {
            const rest = args.splice(1).join(" ")
            vscodeApi.postMessage({
              type: 'session.state.put',
              payload: {
                sid,
                key: args[0],
                val: rest, // overly simple way to do this
              }
            });

            setSession((prev: any) => {
              const next = {
                ...prev,
                state: {
                  ...prev?.state
                }
              }
              next.state[args[0]] = rest
              return next
            });

            // clear input
            handleInput('')

            // make sure listeners have updated conent
            vscodeApi.postMessage({
              type: 'session.get',
              payload: {
                sid,
              }
            });
            vscodeApi.postMessage({
              type: 'session.getList',
              payload: {}
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

      //
      // Otherwise, send a message
      //

      const agent = chatState?.agent
      const model = chatState?.model

      vscodeApi.postMessage({
        type: 'chat.userMessage',
        payload: {
          text: chatState.input,
          sid,
          agent,
          model,
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

  //
  // IMPORTANT, we should coalesce events intelligently, but basically first too
  //   1. group streaming partial responses, don't duplicate (are the events making it this far (yet?)?)
  //   2. more coherent tool components with spinners and buttons to perform actions
  //   3. Snapshots and time travel
  //
  // separately, but related, how do we do custom events like /state update. Those might not be committed
  // ... or are we manually doing that and making other sessions dirty and/or non-reproducible? (we might be ok, and it's more not having snapshots for app: / user: values)

  // maybe put this on the chatState as read-only?
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
    <div className="flex flex-col p-2 gap-2 min-h-screen">
      <Header
        sid={sid}
        session={session}
        usage={usage}
        chatState={chatState}
        className="mx-2"
      />

      <Events
        session={session}
        messagesEndRef={messagesEndRef}
      />

      <UserInput
        sid={sid}
        usage={usage}
        session={session}
        chatState={chatState}
        handleInput={handleInput}
        handleSelectModel={handleSelectModel}
        handleSelectAgent={handleSelectAgent}
        handleSend={handleSend}
      />

    </div>
  )
}



export default App