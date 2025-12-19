import { useRef, useState } from 'react';
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

import { 
  ToolTipper,
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuLabel,
} from 'veg-webview-common'

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

  const [termOpen, setTermOpen] = useState(false);
  const [mergeOpen, setMergeOpen] = useState(false);
  const termTimer = useRef<any>(null);
  const mergeTimer = useRef<any>(null);

  const startTermTimer = () => {
    termTimer.current = setTimeout(() => {
      setTermOpen(true);
    }, 2000);
  };
  const stopTermTimer = () => {
    if (termTimer.current) {
      clearTimeout(termTimer.current);
      termTimer.current = null;
    }
  };

  const startMergeTimer = () => {
    mergeTimer.current = setTimeout(() => {
      setMergeOpen(true);
    }, 2000);
  };
  const stopMergeTimer = () => {
    if (mergeTimer.current) {
      clearTimeout(mergeTimer.current);
      mergeTimer.current = null;
    }
  };


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

      <DropdownMenu open={termOpen} onOpenChange={setTermOpen}>
        <DropdownMenuTrigger asChild>
          <div 
            onMouseEnter={startTermTimer}
            onMouseLeave={stopTermTimer}
            onPointerDown={(e) => e.preventDefault()}
            className="hover:text-green-500 h-4"
          >
            <ToolTipper label="terminal">
              <SquareTerminal size={16}
                aria-label="terminal"
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  vscodeApi.postMessage({
                    type: "session.term.open",
                    payload: {
                      sid,
                      pos: effectivePos,
                    }
                  })
                }}
              />
            </ToolTipper>
          </div>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="bg-[#0e0e0e] text-white border-0 font-mono text-xs font-light p-1">
          <DropdownMenuLabel className="px-1 pb-1 pt-0 mb-1 border-b-1 text-xs">Terminals</DropdownMenuLabel>
          <DropdownMenuItem 
            className="px-2 py-[2px] text-xs focus:bg-[#2e2e2e] hover:bg-[#2e2e2e] focus:text-white hover-text-white"
            onClick={() => {
            vscodeApi.postMessage({
              type: "session.term.open",
              payload: {
                sid,
                pos: effectivePos,
              }
            })
            }}
          >
            fresh
          </DropdownMenuItem>
          {chatState?.terminals?.terminals?.map((t: any) => (
            <DropdownMenuItem key={t.id}
              className="px-2 py-[2px] text-xs focus:bg-[#2e2e2e] hover:bg-[#2e2e2e] focus:text-white hover-text-white"
              onClick={() => {
              vscodeApi.postMessage({
                type: "session.term.open",
                payload: {
                  sid,
                  pos: effectivePos,
                  termId: t.id,
                }
              })
              }}
            >
              {t.name} ({t.id})
            </DropdownMenuItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>

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

      <DropdownMenu open={mergeOpen} onOpenChange={setMergeOpen}>
        <DropdownMenuTrigger asChild>
          <div 
            onMouseEnter={startMergeTimer}
            onMouseLeave={stopMergeTimer}
            onPointerDown={(e) => e.preventDefault()}
            className="hover:text-yellow-500 h-4"
          >
            <ToolTipper label="merge">
              <GitPullRequestCreateArrow size={16}
                aria-label="merge"
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
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
          </div>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="bg-[#0e0e0e] text-white border-0 font-mono text-xs font-light p-1">
          <DropdownMenuLabel className="px-1 pb-1 pt-0 mb-1 border-b-1 text-xs">merge into</DropdownMenuLabel>
          {session?.state?.initEnv?.srcUri && (
            <DropdownMenuItem
              className="px-2 py-[2px] text-xs focus:bg-[#2e2e2e] hover:bg-[#2e2e2e] focus:text-white hover-text-white"
              onClick={() => {
              vscodeApi.postMessage({
                type: "session.merge",
                payload: {
                  sid,
                  pos: effectivePos,
                  currEnv: session?.state?.currEnv,
                  dest: session.state.initEnv.srcUri,
                }
              })
              }}
            >
              {session.state.initEnv.srcUri}
            </DropdownMenuItem>
          )}
          <DropdownMenuItem
            className="px-2 py-[2px] text-xs focus:bg-[#2e2e2e] hover:bg-[#2e2e2e] focus:text-white hover-text-white"
            onClick={() => {
            vscodeApi.postMessage({
              type: "session.merge",
              payload: {
                sid,
                pos: effectivePos,
                currEnv: session?.state?.currEnv,
                forceInput: true,
              }
            })
            }}
          >
            input...
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

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