import { useState, useEffect, useCallback, createContext, useContext, type ReactNode } from 'react';
import { vscodeApi } from '@/vscodeApi.js';

// Define the message types we expect
interface SidPayload {
  sid: string;
}

interface ServerMessage {
  type: string;
  payload: unknown;
}

const defaults = {
  agent: 'veggie',
  model: 'gemini-2.5-flash',
};

function processEvents(session: any) {
  const usage: any = {
    candidatesTokenCount: 0,
    promptTokenCount: 0,
    cachedContentTokenCount: 0,
    thoughtsTokenCount: 0,
    totalTokenCount: 0,
  };
  if (!!(session?.events) && session.events.length > 0) {
    session?.events?.forEach((e: any) => {
      if (e?.UsageMetadata) {
        const u = e.UsageMetadata;
        usage.candidatesTokenCount += u.candidatesTokenCount || 0;
        usage.cachedContentTokenCount += u.cachedContentTokenCount || 0;
        usage.promptTokenCount += u.promptTokenCount || 0;
        usage.thoughtsTokenCount += u.thoughtsTokenCount || 0;
        usage.totalTokenCount += u.totalTokenCount || 0;
      }
    });
  }

  return {
    usage,
  };
}

interface ChatContextType {
  sid: string;
  pos: number;
  session: any;
  usage: any;
  chatState: any;
  diff: any;
  setDiff: (diff: any) => void;
  setPos: (pos: number) => void;
  handleSend: (userInput: any) => void;
}

const ChatContext = createContext<ChatContextType | undefined>(undefined);

export function ChatProvider({ children }: { children: ReactNode }) {
  const state = vscodeApi.getState() || {};

  const [sid, setSid] = useState(state?.sid || '');
  const [pos, setPos] = useState(state?.pos || -1);
  const [session, setSession] = useState<any>(state?.session || {});
  const [diff, setDiff] = useState<any>(state?.diff || {});
  const [usage, setUsage] = useState<any>({});
  const [chatState, setChatState] = useState<any>(state?.chatState || defaults);

  // Effect to listen for messages from the extension
  useEffect(() => {
    console.log("[] effect");

    // The 'onMessage' helper returns a cleanup function
    const removeListener = vscodeApi.onMessage((event) => {
      const message = event.data as ServerMessage;

      if (message.type === 'config.info.resp') {
        setChatState((prev: any) => {
          const s = vscodeApi.getState();
          const next = {
            ...prev,
            config: message.payload,
          };
          vscodeApi.setState({
            ...s,
            chatState: next,
          });
          return next;
        });
      }

      if (message.type === 'env.info.resp') {
        setChatState((prev: any) => {
          const s = vscodeApi.getState();
          const next = {
            ...prev,
            env: message.payload,
          };
          vscodeApi.setState({
            ...s,
            chatState: next,
          });
          return next;
        });
      }

      if (message.type === 'terminal.info') {
        setChatState((prev: any) => {
          const s = vscodeApi.getState();
          const next = {
            ...prev,
            terminals: message.payload,
          };
          vscodeApi.setState({
            ...s,
            chatState: next,
          });
          return next;
        });
      }

      if (message.type === 'window.info.resp') {
        setChatState((prev: any) => {
          const s = vscodeApi.getState();
          const next = {
            ...prev,
            window: message.payload,
          };
          vscodeApi.setState({
            ...s,
            chatState: next,
          });
          return next;
        });
      }

      if (message.type === 'workspace.info.resp') {
        setChatState((prev: any) => {
          const s = vscodeApi.getState();
          const next = {
            ...prev,
            workspace: message.payload,
          };
          vscodeApi.setState({
            ...s,
            chatState: next,
          });
          return next;
        });
      }

      if (message.type === 'models.list.resp') {
        setChatState((prev: any) => {
          const s = vscodeApi.getState();
          const next = {
            ...prev,
            models: message.payload,
          };
          vscodeApi.setState({
            ...s,
            chatState: next,
          });
          return next;
        });
      }

      if (message.type === 'agents.list.resp') {
        setChatState((prev: any) => {
          const s = vscodeApi.getState();
          const next = {
            ...prev,
            agents: message.payload,
          };
          vscodeApi.setState({
            ...s,
            chatState: next,
          });
          return next;
        });
      }

      // Handle the message based on its type
      if (message.type === 'chat.event') {
        setSession((prev: any) => {
          const next = {
            ...prev,
            events: [...(prev.events || []), message.payload],
          };
          const { usage: u } = processEvents(next);
          setUsage(u);
          const s = vscodeApi.getState();
          vscodeApi.setState({
            ...s,
            session: next,
          });
          return next;
        });
      }

      if (message.type === 'chat.event.error') {
        setSession((prev: any) => {
          const next = {
            ...prev,
            error: message.payload,
          };
          const s = vscodeApi.getState();
          vscodeApi.setState({
            ...s,
            session: next,
          });
          return next;
        });
      }

      if (message.type === 'chat.loadSession') {
        const payload = message.payload as SidPayload;
        if (payload.sid) {
          // update core state
          const s = vscodeApi.getState();
          vscodeApi.setState({
            ...s,
            sid: payload.sid,
            session: { sid: payload.sid },
          });
          setSid(payload.sid);
          setSession({ sid: payload.sid });

          // request session info when we load
          vscodeApi.postMessage({
            type: 'session.get',
            payload: {
              sid: payload.sid,
            },
          });
          vscodeApi.postMessage({
            type: 'session.diff',
            payload: {
              sid: payload.sid,
            },
          });
        }
      }
    });

    // console.log("initial fetch", state?.sid, state?.sid !== "")
    vscodeApi.postMessage({ type: 'requestSync' });
    if (state?.sid !== '') {
      vscodeApi.postMessage({
        type: 'session.get',
        payload: {
          sid: state.sid,
        },
      });
      vscodeApi.postMessage({
        type: 'session.diff',
        payload: {
          sid: state.sid,
        },
      });
    }

    // Return the cleanup function
    return removeListener;
  }, []); // Empty dependency array means this runs once

  // Calculate usage when session changes
  useEffect(() => {
    const { usage: u } = processEvents(session);
    setUsage(u);
  }, [session?.events]);

  // update our listener when the sid changes
  useEffect(() => {
    console.log("sid effect");
    // The 'onMessage' helper returns a cleanup function
    const removeListener = vscodeApi.onMessage((event) => {
      const message = event.data as ServerMessage;

      if (message.type === 'session.get.resp' || message.type === 'session.info') {
        const payload = message.payload as SidPayload;
        if (payload.sid === sid) {
          const s = vscodeApi.getState();
          vscodeApi.setState({
            ...s,
            sid: payload.sid,
            session: payload,
          });
          setSession(payload as any);
          const { usage: u } = processEvents(session);
          setUsage(u);
        }
      }

      if (message.type === 'session.diff') {
        const payload = message.payload as { sid: string, pos?: number };
        if (payload.sid === sid && pos && pos > 0) {

          const s = vscodeApi.getState();
          vscodeApi.setState({
            ...s,
            pos,
          });
          setPos(pos)
        }
      }

      if (message.type === 'session.diff.resp') {
        const payload = message.payload as { sid: string };
        if (payload.sid === sid && pos && pos > 0) {
          setDiff(payload)
        }
      }

      if (message.type === 'session.delete') {
        const payload = message.payload as SidPayload;
        if (payload.sid === sid) {
          const s = vscodeApi.getState();
          vscodeApi.setState({
            ...s,
            sid: "",
            pos: -1,
            usage: {},
            session: {},
            diff: {},
          });
          setSid('');
          setPos(-1)
          setSession({});
          setUsage({});
          setDiff({});
        }
      }
    });

    // Return the cleanup function
    return removeListener;
  }, [sid]); // Empty dependency array means this runs once

  // Handler for sending a message
  const handleSend = useCallback((userInput: any) => {
    console.log("handleSend", userInput);
    //
    // send a message to the agent / model
    //
    if (userInput.text.trim() && userInput.model && userInput.agent) {
      console.log("CHAT!", userInput);

      vscodeApi.postMessage({
        type: 'chat.userMessage',
        payload: {
          sid,
          ...userInput,
        },
      });

      setSession((prev: any) => {
        const next = {
          ...prev,
          events: [...(prev?.events || []), {
            Author: "user",
            Content: {
              role: "user",
              parts: [{
                text: userInput.text,
              }],
            },
            Timestamp: new Date().toISOString(),
          }],
        };
        const s = vscodeApi.getState();
        vscodeApi.setState({
          ...s,
          session: next,
        });

        return next;
      });
    }
  }, [sid]);

  const value = {
    sid,
    pos,
    session,
    usage,
    chatState,
    diff,
    setDiff,
    setPos,
    handleSend,
  };

  return (
    <ChatContext.Provider value={value}>
      {children}
    </ChatContext.Provider>
  );
}

export function useChat() {
  const context = useContext(ChatContext);
  if (context === undefined) {
    throw new Error('useChat must be used within a ChatProvider');
  }
  return context;
}
