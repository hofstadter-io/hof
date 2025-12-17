import { useRef } from 'react'
import '@/index.css' // We'll add some styles

import { vscodeApi } from '@/vscodeApi.js'
import { StickToBottom, useStickToBottomContext } from 'use-stick-to-bottom';

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable"

import { Header } from '@/components/Header'
import { Welcome } from '@/components/Welcome'
import { Events } from '@/components/Events';
import { UserInput } from '@/components/UserInput';

import { useChat } from '@/hooks/useChat';

import {
  ScrollTo,
} from '@/components/ScrollTo'

function App() {
  const state = vscodeApi.getState()
  console.log("chat state:", state)
  const headerRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const {
    sid,
    pos,
    setPos,
    session,
    usage,
    diff,
    chatState,
    handleSend,
  } = useChat(messagesEndRef as any);


  //
  // IMPORTANT, we should coalesce events intelligently, but basically first too
  //   1. group streaming partial responses, don't duplicate (are the events making it this far (yet?)?)
  //   2. more coherent tool components with spinners and buttons to perform actions
  //   3. Snapshots and time travel
  //
  // separately, but related, how do we do custom events like /state update. Those might not be committed
  // ... or are we manually doing that and making other sessions dirty and/or non-reproducible? (we might be ok, and it's more not having snapshots for app: / user: values)

  return (
    <ResizablePanelGroup direction="vertical" className="min-h-screen p-4 bg-[#1e1e1e]">
      <ResizablePanel>
        <StickToBottom className="flex flex-col p-2 gap-2 h-full relative" resize="smooth" initial="smooth">
          <StickToBottom.Content className="flex-grow h-full flex flex-col gap-4">

            <Header
              ref={headerRef}
              sid={sid}
              setPos={setPos}
              session={session}
              usage={usage}
              chatState={chatState}
              diff={diff}
              className="border-b-2 border-fuchsia-700/70 pb-2"
            />

            {session?.events?.length > 0 ?
              <Events
                sid={sid}
                currPos={pos}
                setPos={setPos}
                session={session}
                events={session?.events}
                messagesEndRef={messagesEndRef}
              />
              : <Welcome username={chatState?.env?.user} />
            }

          </StickToBottom.Content>

          <ScrollTo target={headerRef} />

          {/* This component uses `useStickToBottomContext` to scroll to bottom when the user enters a message */}
          {/* <ChatBox /> */}
        </StickToBottom>

      </ResizablePanel>
      <ResizableHandle className="pt-[3px] rounded-xl bg-fuchsia-500/20 hover:bg-fuchsia-500/70"/>
      <ResizablePanel defaultSize={20}>
        <div className="h-full overflow-y-auto">
          <UserInput
            sid={sid}
            setPos={setPos}
            usage={usage}
            session={session}
            chatState={chatState}
            diff={diff}
            handleSend={handleSend}
          />
        </div>
      </ResizablePanel>

    </ResizablePanelGroup>
  )
}



export default App
