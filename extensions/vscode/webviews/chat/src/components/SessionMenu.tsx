import { vscodeApi } from "@/vscodeApi";
import { useChat } from "@/hooks/useChat";
import {
  Braces,
  FileDiff,
  GitGraph,
  GitPullRequestCreateArrow,
  ListTree,
  RefreshCcw,
  SquareTerminal,
} from 'lucide-react'

import { ToolTipper } from 'veg-webview-common'

export const Menu = ({
  pos,
  hidden,
  refresh,
  setHidden,
}:{
  pos?: number,
  hidden: boolean,
  refresh?: boolean,
  setHidden: (prev: any) => any
}) => {
  const { sid, session, setPos, pos: currentPos, chatState } = useChat();
  const effectivePos = pos ?? currentPos;

  return (
    <div className="ml-auto flex justify-end items-center gap-2">

      <ToolTipper label="browse">
      <ListTree size={16}
        aria-label="diff"
        className="hover:text-sky-500"
        onClick={() => {
          // currently, this is only enabled on the headers (not per message/event)
          console.log("Browse session:", session)
          if (!session) {
            return
          }
          vscodeApi.postMessage({
            type: "filesys.openEnviron",
            payload: {
              session,
            }
          })
        }}
      />
      </ToolTipper>

      <ToolTipper label="terminal">
      <div className="hover:text-green-500">
        <SquareTerminal size={16}
          aria-label="terminal"
          onClick={(e) => {
            console.log("Terminal session:", session)
            vscodeApi.postMessage({
              type: "session.term.open",
              payload: {
                sid,
                pos: effectivePos,
              }
            })
          }}
        />
      </div>
      </ToolTipper>

      <ToolTipper label="diff">
      <FileDiff size={16}
        aria-label="diff"
        className="hover:text-yellow-500"
        onClick={() => {
          if (pos !== undefined) {
            setPos(pos)
          }
          vscodeApi.postMessage({
            type: "session.diff",
            payload: {
              sid,
              pos: effectivePos,
              show: true,
              currEnv: session?.state?.currEnv,
            }
          })
        }}
      />
      </ToolTipper>

      <ToolTipper label="merge">
      <GitPullRequestCreateArrow size={16}
        aria-label="merge"
        className="hover:text-yellow-500"
        onClick={() => {
          vscodeApi.postMessage({
            type: "session.merge",
            payload: {
              sid,
              pos: effectivePos,
              currEnv: session?.state?.currEnv,
            }
          })
        }}
      />
      </ToolTipper>

      <ToolTipper label="fork">
      <GitGraph size={16}
        aria-label="fork"
        className="hover:text-sky-500"
        onClick={() => {
          // vscodeApi.postMessage({
          //   type: "session.fork",
          //   payload: {
          //     from: sid,
          //     pos,
          //   }
          // })
        }}
      />
      </ToolTipper>

      <ToolTipper label="refresh">
      { refresh && <RefreshCcw size={16}
        aria-label="refresh"
        className="hover:text-sky-500"
        onClick={() => {
          vscodeApi.postMessage({
            type: "session.get",
            payload: {
              sid,
            }
          })
        }}
      /> }
      </ToolTipper>

      <ToolTipper label="details">
      <div className="hover:text-sky-500">
        <Braces size={16}
          aria-label="details"
          className="hover:text-sky-500" 
          onClick={(e) => {
            e.preventDefault()
            e.stopPropagation()
            setHidden(!hidden)
          }}
        />
      </div>
      </ToolTipper>

    </div>
  )

}