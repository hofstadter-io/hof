import { useState, useEffect } from 'react'
import { vscodeApi } from './vscodeApi.js'
import './index.css' // We'll add some styles
import JsonView from '@uiw/react-json-view';
import { vscodeTheme } from '@uiw/react-json-view/vscode';


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

const setPairs: Record<string,string> = {
  "terminal.info": "terminals",
  "models.list.resp": "models",
  "agents.list.resp": "agents",
  "chat.loadSession": "sid",
  "session.info": "session",
  "session.list": "sessions",
  "env.info.resp": "env",
  "window.info.resp": "window",
  "workspace.info.resp": "workspace",
}

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

      //
      // Check for a bunch of informational messages we want to capture
      //
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
        return
      }

      //
      // handle other message types
      //

      // legacy message type (rename)
      if (message.type === 'terminalInfo') {
        setDebugValue((prevData: any) => {
          const next = {
            ...prevData,
            terminals: payload?.terminals, // old nested format
          }
          vscodeApi.setState(next)
          return next
        });
        return
      }

    });

    if (state?.sid) {
      vscodeApi.postMessage({
        type: 'requestSync',
        payload: {
          sid: state.sid,
        }
      });
    }

    // Return the cleanup function
    return removeListener;
  }, []); // Empty dependency array means this runs once

  useEffect(() => {
    const removeListener = vscodeApi.onMessage((event) => {
      const message = event.data as Message;

      var payload: any = message.payload;
      if (typeof message.payload === "string") {
        payload = JSON.parse(message.payload);
      }

      // check if our current session has been deleted
      console.log("debug.sessionDelete?", message, debugValue)
      if (message.type === 'session.delete' && payload.sid === debugValue.session.sid) {
        setDebugValue((prevData: any) => {
          const next = {
            ...prevData,
            sid: null,
            session: null,
          }
          vscodeApi.setState(next)
          return next
        });
        return
      }
    })

    // Return the cleanup function
    return removeListener;

  }, [debugValue.sid, debugValue.session])

  return (
    <div className="flex">
      <pre className="text-sm">
        <JsonView value={debugValue} style={vscodeTheme} displayDataTypes={false} indentWidth={12} shortenTextAfterLength={120}/>
      </pre>
    </div>
  )
}

export default App