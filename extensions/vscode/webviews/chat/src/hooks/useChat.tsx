import { useState, useEffect, useCallback, createContext, useContext, useMemo, type ReactNode } from 'react';
import { vscodeApi } from '@/vscodeApi.js';
import { prices } from '@/lib/prices';

interface ServerMessage {
  type: string;
  payload: any;
}

const DEFAULTS = {
  agent: 'coding_assist',
  model: 'gemini-3-flash',
};

interface ChatContextType {
  sid: string;
  pos: number;
  session: any;
  usage: any;
  chatState: any;
  diff: any;
  setChatState: (updater: any) => void;
  setDiff: (diff: any) => void;
  setPos: (pos: number) => void;
  handleSend: (userInput: any) => void;
}

const ChatContext = createContext<ChatContextType | undefined>(undefined);

export function ChatProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState(() => ({
    sid: '',
    pos: -1,
    session: {},
    chatState: DEFAULTS,
    diff: {},
    ...vscodeApi.getState()
  }));

  const updateState = useCallback((updates: Partial<typeof state> | ((prev: typeof state) => Partial<typeof state>)) => {
    setState((prev: any) => {
      const nextUpdates = typeof updates === 'function' ? updates(prev) : updates;
      const next = { ...prev, ...nextUpdates };
      vscodeApi.setState(next);
      return next;
    });
  }, []);

  // Derived usage stats
  const usage = useMemo(() => calculateUsage(state.session), [state.session?.events, state.session?.state?.model]);

  useEffect(() => {
    const removeListener = vscodeApi.onMessage((event) => {
      const { type, payload } = event.data as ServerMessage;

      switch (type) {
        case 'config.info.resp':
        case 'env.info.resp':
        case 'terminal.info':
        case 'window.info.resp':
        case 'workspace.info.resp':
        case 'models.list.resp':
        case 'agents.list.resp': {
          const keyMap: Record<string, string> = {
            'config.info.resp': 'config',
            'env.info.resp': 'env',
            'terminal.info': 'terminals',
            'window.info.resp': 'window',
            'workspace.info.resp': 'workspace',
            'models.list.resp': 'models',
            'agents.list.resp': 'agents',
          };
          updateState(prev => ({
            chatState: { ...prev.chatState, [keyMap[type]]: payload }
          }));
          break;
        }

        case 'chat.event':
          updateState(prev => ({
            session: {
              ...prev.session,
              events: [...(prev.session.events || []), payload]
            }
          }));
          break;

        case 'chat.event.error':
          updateState({ session: { ...state.session, error: payload } });
          break;

        case 'chat.loadSession':
          if (payload.sid) {
            updateState({ sid: payload.sid, session: { sid: payload.sid } });
            vscodeApi.postMessage({ type: 'session.get', payload: { sid: payload.sid } });
            vscodeApi.postMessage({ type: 'session.diff', payload: { sid: payload.sid } });
          }
          break;

        case 'session.get.resp':
        case 'session.info':
          updateState(prev => {
            if (payload.sid !== prev.sid) return {};
            const hasMoreEvents = (payload.events?.length || 0) > (prev.session.events?.length || 0);
            if (prev.session.sid !== prev.sid || !prev.session.events || hasMoreEvents) {
              return { session: payload };
            }
            return {};
          });
          break;

        case 'session.diff':
          if (payload.sid === state.sid && payload.pos !== undefined) {
            updateState({ pos: payload.pos });
          }
          break;

        case 'session.diff.resp':
          if (payload.sid === state.sid && state.pos > 0) {
            updateState({ diff: payload });
          }
          break;

        case 'session.delete':
          if (payload.sid === state.sid) {
            updateState({ sid: '', pos: -1, session: {}, diff: {} });
          }
          break;

        case 'requestSync':
        case 'chat.userInput':
          vscodeApi.postMessage({ type: 'chat.userInput.resp', payload: state.chatState.userInput });
          break;
      }
    });

    vscodeApi.postMessage({ type: 'requestSync' });
    if (state.sid) {
      vscodeApi.postMessage({ type: 'session.get', payload: { sid: state.sid } });
      vscodeApi.postMessage({ type: 'session.diff', payload: { sid: state.sid } });
    }

    return removeListener;
  }, [state.sid, state.pos, updateState]);


  useEffect(() => {
    const removeListener = vscodeApi.onMessage((event) => {
      const { type, payload } = event.data as ServerMessage;

      switch (type) {
        case 'requestSync':
        case 'chat.userInput':
          vscodeApi.postMessage({ type: 'chat.userInput.resp', payload: state.chatState.userInput });
          break;
      }
    });

    // vscodeApi.postMessage({ type: 'requestSync' });
    // if (state.sid) {
    //   vscodeApi.postMessage({ type: 'session.get', payload: { sid: state.sid } });
    //   vscodeApi.postMessage({ type: 'session.diff', payload: { sid: state.sid } });
    // }

    return removeListener;
  }, [state.sid, state.chatState.userInput]);


  const handleSend = useCallback((userInput: any) => {
    if (!userInput.text.trim() || !userInput.model || !userInput.agent) return;

    vscodeApi.postMessage({
      type: 'chat.userMessage',
      payload: { sid: state.sid, ...userInput },
    });

    // updateState(prev => ({
    //   session: {
    //     ...prev.session,
    //     events: [...(prev.session?.events || []), {
    //       Author: "user",
    //       Content: { role: "user", parts: [{ text: userInput.text }] },
    //       Timestamp: new Date().toISOString(),
    //     }],
    //   }
    // }));
  }, [state.sid]);

  const value = {
    ...state,
    usage,
    setChatState: (updater: any) => updateState(prev => ({
      chatState: typeof updater === 'function' ? updater(prev.chatState) : updater
    })),
    setDiff: (diff: any) => updateState({ diff }),
    setPos: (pos: number) => updateState({ pos }),
    handleSend,
  };

  return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>;
}

export function useChat() {
  const context = useContext(ChatContext);
  if (!context) throw new Error('useChat must be used within a ChatProvider');
  return context;
}

function calculateUsage(session: any) {
  const usage: any = {
    candidatesTokenCount: 0,
    promptTokenCount: 0,
    cachedContentTokenCount: 0,
    thoughtsTokenCount: 0,
    totalTokenCount: 0,
    costInput: 0,
    costCache: 0,
    costThink: 0,
    costWrite: 0,
    costOutput: 0,
    costTotal: 0,
    model: session?.state?.model,
  };

  const model = session?.state?.model;
  const price = (model && (prices as any)[model]) ? (prices as any)[model] : null;

  console.log("calc'n prices with:", price)

  session?.events?.forEach((e: any) => {
    if (e?.UsageMetadata) {
      const u = e.UsageMetadata;
      usage.candidatesTokenCount += u.candidatesTokenCount || 0;
      usage.cachedContentTokenCount += u.cachedContentTokenCount || 0;
      usage.promptTokenCount += u.promptTokenCount || 0;
      usage.thoughtsTokenCount += u.thoughtsTokenCount || 0;
      usage.totalTokenCount += u.totalTokenCount || 0;

      if (price?.input) {
        const isLong = (u.promptTokenCount || 0) > (price.input.cutoff || 0);
        const pInput = isLong ? price.input.long : price.input.short;
        const pCache = isLong ? price.cache.long : price.cache.short;
        const pOutput = isLong ? price.output.long : price.output.short;

        const uncached = (u.promptTokenCount || 0) - (u.cachedContentTokenCount || 0);
        // const output = (u.thoughtsTokenCount || 0) + (u.candidatesTokenCount || 0);

        usage.costInput += (uncached / 1000000) * pInput;
        usage.costCache += ((u.cachedContentTokenCount || 0) / 1000000) * pCache;
        // const costOut = (output / 1000000) * pOutput;
        usage.costThink += ((u.thoughtsTokenCount || 0) / 1000000) * pOutput;
        usage.costWrite += ((u.candidatesTokenCount || 0) / 1000000) * pOutput;
      }
    }
    usage.costOutput = usage.costThink + usage.costWrite;
    usage.costTotal = usage.costInput + usage.costCache + usage.costOutput;
  });

  return usage;
}
