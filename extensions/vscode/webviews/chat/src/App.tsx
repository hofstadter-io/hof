import { useRef } from 'react'
import '@/index.css' // We'll add some styles

// test comment to see what debug looks like

import { Header } from '@/components/Header.js'
import { Events } from '@/components/Messages.js';
import { UserInput } from '@/components/UserInput.js';

import { useChat } from '@/hooks/useChat.js';

function App() {
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const {
    sid,
    pos,
    setPos,
    session,
    usage,
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
    <div className="flex flex-col p-2 gap-2 min-h-screen">
      <Header
        sid={sid}
        setPos={setPos}
        session={session}
        usage={usage}
        chatState={chatState}
        className="mx-2"
      />

      <Events
        sid={sid}
        currPos={pos}
        setPos={setPos}
        events={session?.events}
        messagesEndRef={messagesEndRef}
      />

      <UserInput
        sid={sid}
        setPos={setPos}
        usage={usage}
        session={session}
        chatState={chatState}
        handleSend={handleSend}
      />

    </div>
  )
}



export default App
