import { useState, useEffect, useCallback } from 'react'
import { vscodeApi } from './vscodeApi.js'
import './index.css' // We'll add some styles

import { Button } from 'veg-webview-common'

// Define the message types we expect
// (These should match your Go server and extension)
interface ManageResponsePayload {
  responseText: string;
}

interface ServerMessage {
  type: string;
  payload: unknown;
}

function App() {
  const state = vscodeApi.getState()
  console.log("manage state:", state)
  const [value, setValue] = useState({});

  // Effect to listen for messages from the extension
  useEffect(() => {
    // The 'onMessage' helper returns a cleanup function
    const removeListener = vscodeApi.onMessage((event) => {
      console.log("App event:", event)
      const message = event.data as ServerMessage;
      
      // Handle the message based on its type
      if (message.type === 'manageResponse') {
        const payload = message.payload as ManageResponsePayload;
        const data = JSON.parse(payload.responseText)
        setValue(data);
      }
      // You could add other 'else if' blocks here
      // for different message types from the server.
    });

    // Return the cleanup function
    return removeListener;
  }, []); // Empty dependency array means this runs once

  return (
    <div className="flex flex-col p-2 gap-2 min-h-screen">
      <div className="rounded border m-2 p-2 font-thin text-xl">
        <Button>veggie!</Button>
        <div className="border m-1 p-1 gap-1 flex flex-col md:flex-row">
          <div className="border p-1 w-full">A1</div>
          <div className="border p-1 w-full">A2</div>
        </div>
      </div>
    </div>
  )
}

export default App