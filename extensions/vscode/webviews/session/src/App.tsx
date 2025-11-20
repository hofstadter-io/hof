import { useState, useEffect, useCallback } from 'react'
import { vscodeApi } from './vscodeApi.js'
import './index.css' // We'll add some styles

// generic message with type & payload
interface Message {
  type: string;
  payload: unknown;
}

class SessionPayload {
	id?: string;
	state?: any;
	lastUpdate?: string;
}

function App() {
  const [sessionsList, setSessionsList] = useState<SessionPayload[]>([]);

  // Effect to listen for messages from the extension
  useEffect(() => {
    // The 'onMessage' helper returns a cleanup function
    const removeListener = vscodeApi.onMessage((event) => {
      // console.log("session: event handler:", event)
      const message = event.data as Message;
      
      if (message.type === 'sessionsList') {
        var sessions: SessionPayload[] = []
        if (typeof message.payload === "string") {
          sessions = JSON.parse(message.payload)
        } else {
          sessions = message.payload as SessionPayload[];
        }
        setSessionsList(sessions);
      }
    });

    // Return the cleanup function
    return removeListener;
  }, []); // Empty dependency array means this runs once

  // Handler for sending a message
  // @ts-ignore
  const loadSession = useCallback(() => {
      // 1. Post message to the extension
      // This will be caught by `sidebar.ts`
      vscodeApi.postMessage({
        type: 'loadSession',
        text: JSON.stringify(sessionsList),
      });
  }, [sessionsList]);

  // @ts-ignore
  const deleteSession = useCallback((id: string) => {
      // 1. Post message to the extension
      // This will be caught by `sidebar.ts`
      vscodeApi.postMessage({
        type: 'deleteSession',
        text: JSON.stringify({
          id,
        }),
      });
  }, []);

  return (
    <div className="flex flex-col">
      {sessionsList.map((s) => {
        return (
          <Session session={s} key={s.id} />
        )
      })}
    </div>
  )
}

function Session({
  session
}:{ 
  session: any
}) {
  return (
    <div className="p-1 border">
      {session.id}
    </div>
  )
}

export default App