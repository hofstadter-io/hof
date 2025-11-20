import { useState, useEffect, useCallback } from 'react'
import { vscodeApi } from './vscodeApi.js'
import './index.css' // We'll add some styles

// Define the message types we expect
// (These should match your Go server and extension)
interface ConfigResponsePayload {
  responseText: string;
}

interface ServerMessage {
  type: string;
  payload: unknown;
}

function App() {
  const [configValue, setConfigValue] = useState({});

  // Effect to listen for messages from the extension
  useEffect(() => {
    // The 'onMessage' helper returns a cleanup function
    const removeListener = vscodeApi.onMessage((event) => {
      console.log("App event:", event)
      const message = event.data as ServerMessage;
      
      // Handle the message based on its type
      if (message.type === 'configResponse') {
        const payload = message.payload as ConfigResponsePayload;
        const data = JSON.parse(payload.responseText)
        setConfigValue(data);
      }
      // You could add other 'else if' blocks here
      // for different message types from the server.
    });

    // Return the cleanup function
    return removeListener;
  }, []); // Empty dependency array means this runs once

  // Handler for sending a message
  // @ts-ignore
  const handleSend = useCallback(() => {
      // 1. Post message to the extension
      // This will be caught by `sidebar.ts`
      vscodeApi.postMessage({
        type: 'sendConfig',
        text: JSON.stringify(configValue),
      });
  }, [configValue]);

  return (
    <div className="flex">
      <pre className="">
        {JSON.stringify(configValue, null, "  ")}
      </pre>
    </div>
  )
}

export default App