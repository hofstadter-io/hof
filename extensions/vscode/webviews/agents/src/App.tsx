import { useState, useEffect } from 'react'
import { vscodeApi } from './vscodeApi.js'
import './index.css' // We'll add some styles

// generic message with type & payload
interface Message {
  type: string;
  payload: unknown;
}

// specific message types we expect
// (These should match your Go server and extension)
class TerminalPayload {
	id: number = -1;
	name?: string;
	history?: HistoryPayload[];
}

class HistoryPayload {
	cmd?: any;
	cwd?: string;
	out?: string;
	exit?: number;
}

interface TerminalsPayload {
  terminals: string;
}

class SessionPayload {
	id?: string;
	state?: any;
	lastUpdate?: string;
}

function App() {
  const [debugValue, setDebugValue] = useState({});

  // Effect to listen for messages from the extension
  useEffect(() => {
    // The 'onMessage' helper returns a cleanup function
    const removeListener = vscodeApi.onMessage((event) => {
      console.log("debug: event handler:", event)
      const message = event.data as Message;
      
      // Handle the message based on its type
      if (message.type === 'terminalInfo') {
        // console.log("debug: got terminalInfo")
        const payload = message.payload as TerminalsPayload;
        var terminals: TerminalPayload[] = []
        if (typeof payload.terminals === "string") {
          terminals = JSON.parse(payload.terminals)
        } else {
          terminals = payload.terminals as TerminalPayload[]
        }
        // console.log("debug: terminals", terminals)
        setDebugValue((prevData) => {
          return {
            ...prevData,
            terminals,
          }
        });
      }

      if (message.type === 'sessionsList') {
        // console.log("debug: got terminalInfo")
        var sessions: SessionPayload[] = []
        if (typeof message.payload === "string") {
          sessions = JSON.parse(message.payload)
        } else {
          sessions = message.payload as SessionPayload[];
        }
        // console.log("debug: terminals", terminals)
        setDebugValue((prevData) => {
          return {
            ...prevData,
            sessions,
          }
        });
      }
    });

    // Return the cleanup function
    return removeListener;
  }, []); // Empty dependency array means this runs once

  return (
    <div className="flex">
      <pre className="text-sm">
        {JSON.stringify(debugValue, null, "  ")}
      </pre>
    </div>
  )
}

export default App