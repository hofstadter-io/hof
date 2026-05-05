"use client"
import * as React from "react"

import { useCookies } from 'react-cookie';
import { useLocalStorage } from "@uidotdev/usehooks";

import * as jose from 'jose';

export type AuthrContext = {
  session: any;
  sessions: any;
  options: AuthrOptions;
  login: (value: AuthrLoginInput) => Promise<void>;
  switchAccount: (did: string) => void;
  // logout: (did: string) => void;
}

export type AuthrOptions = {
  cookieName: string;
  cookieDomain: string;
  oauthHost: string;
  xrpcHost: string;
}

export type AuthrLoginInput = {
  handle: string;
  redirect?: string;
}

export type AuthrSession = {
  cookie?: string
  did: string
  pds: string
  handle: string
  sessionId?: string
}

export type AuthrSessions = {
  current?: AuthrSession
  accounts: AuthrSession[]
}

const AuthrContext = React.createContext<AuthrContext | undefined>(undefined);

export const useAuthr = () => {
  const context = React.useContext(AuthrContext);
  if (context === undefined) {
    throw new Error("useAuthr must be used within a AuthrProvider");
  }
  return context;
};

export default AuthrContext;

export const AuthrProvider: React.FC<{ options: AuthrOptions; children: React.ReactNode }> = ({ options, children }) => {

  // hooks

  const [cookies, setCookies] = useCookies()
  const [sessions, setSessions] = useLocalStorage<AuthrSessions>("veggie/sessions", { current: undefined, accounts: [] });

  React.useEffect(() => { 
    // if we have an Authr cookie, lets see if we need to add it
    const authrCookie = cookies[options.cookieName]
    if (authrCookie) {
      const claims = jose.decodeJwt(authrCookie || "")

      // build session object
      const session: any = {
        cookie: authrCookie,
        ...claims,
      }
      // console.log("AuthrProvider.session", session)
      // console.log("AuthrProvider.sessions.pre", sessions)

      // if we have a match, overwrite
      let found = false
      for (const s in sessions.accounts) {
        if (claims.did === sessions.accounts[s].did) {
          sessions.accounts[s] = session
          found = true
          break
        }
      }
      // if no match, new account
      if (!found && session.did) {
        sessions.accounts.push(session)
      }

      // finalize
      setSessions({ current: session, accounts: sessions.accounts })
    }
  }, [cookies])

  // Build up final context object
  const contextValue: AuthrContext = {
    session: sessions.current,
    sessions,
    options,
    login: async (value: AuthrLoginInput) => {
      const b = JSON.stringify(value)
      const resp = await fetch(`${options.oauthHost}/oauth/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: b
      })
      const data: any = await resp.json()
      if (data.error) {
        // TODO, update form or page...
        alert(data.error)
        return
      }

      // this should always be the case, this is the callback to our auth server (post user approval)
      const redir = data.redirect
      window.location.href = redir

    },
    switchAccount: (did: string) => {
      const newSession = sessions.accounts.find(s => s.did === did);
      if (newSession) {
        // this should trigger changing the current session
        // @not-ts-ignore (double check this)
        setSessions({ current: newSession, accounts: sessions.accounts })
        setCookies(options.cookieName, newSession.cookie, {
          path: '/',
          domain: options.cookieDomain,
          sameSite: 'lax', 
        });
      }
    }
  }

  return (
    <AuthrContext.Provider value={contextValue}>
      {children}
    </AuthrContext.Provider>
  );
};
