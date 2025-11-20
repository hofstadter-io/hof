import { useState, useEffect } from 'react'
import { vscodeApi } from './vscodeApi.js'
import './index.css' // We'll add some styles

// generic message with type & payload
interface Message {
  type: string;
  payload: any;
}

// // specific message types we expect
// // (These should match your Go server and extension)
// class TerminalPayload {
// 	id: number = -1;
// 	name?: string;
// 	history?: HistoryPayload[];
// }

// class HistoryPayload {
// 	cmd?: any;
// 	cwd?: string;
// 	out?: string;
// 	exit?: number;
// }

// interface TerminalsPayload {
//   terminals: string;
// }

// class SessionPayload {
// 	id?: string;
// 	state?: any;
// 	lastUpdate?: string;
// }

function App() {
  const state = vscodeApi.getState()
  const [debugValue, setDebugValue] = useState(state || {});

  // Effect to listen for messages from the extension
  useEffect(() => {
    // The 'onMessage' helper returns a cleanup function
    const removeListener = vscodeApi.onMessage((event) => {
      const message = event.data as Message;

      var payload: any = message.payload;
      if (typeof message.payload === "string") {
        payload = JSON.parse(message.payload);
      }

      const setPairs: Record<string,string> = {
        "terminal.info": "terminals",
        "models.list.resp": "models",
        "agents.list.resp": "agents",
        "chat.loadSession": "sid",
        "session.info": "session",
        "session.list": "sessions",
      }

      const pair = setPairs[message.type]
      if (!!pair && pair !== "") {
        setDebugValue((prev: any) => {
          const next = {
            ...prev,
          }
          next[pair] = payload
          vscodeApi.setState(next)
          return next
        });
      } else {
        // legacy "off-by-one"
        if (message.type === 'terminalInfo') {
          setDebugValue((prevData: any) => {
            const next = {
              ...prevData,
              terminals: payload?.terminals, // old nested format
            }
            vscodeApi.setState(next)
            return next
          });
        }
      }
    });

    vscodeApi.postMessage({
      type: 'requestSync',
      payload: {
        id: state.sid,
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