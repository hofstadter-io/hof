import { vscodeApi } from "@/vscodeApi";
import {
  Braces,
  Columns3,
  RefreshCcw,
  SquareTerminal,
  GitPullRequestCreateArrow,
  GitGraph,
  ListTree,
} from 'lucide-react'

import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"

export const Menu = ({
  sid,
  pos,
  setPos,
  hidden,
  refresh,
  openFS,
  // checkpoint,
  setHidden,
}:{
  sid: string,
  pos?: number,
  setPos: any,
  hidden: boolean,
  refresh?: boolean,
  openFS?: boolean,
  checkpoint?: boolean,
  setHidden: (prev: any) => any
}) => {

  const Tooltipped = ({ label, children }: { label: string, children: any}) => {
    return (
      <Tooltip>
        <TooltipTrigger>{children}</TooltipTrigger>
        <TooltipContent>
          <p>{label}</p>
        </TooltipContent>
      </Tooltip>
    )
  }
  return (
    <div className="ml-auto flex justify-end items-center gap-2">
      {openFS && <Tooltipped label="diff">
      <Columns3 size={16}
        aria-label="diff"
        className="hover:text-violet-500"
        onClick={() => {
          setPos(pos)
          vscodeApi.postMessage({
            type: "session.diff",
            payload: {
              sid,
              pos,
              show: true,
              // todo, add start here
            }
          })
        }}
      />
      </Tooltipped>}

      <Tooltipped label="browse">
      <ListTree size={16}
        aria-label="diff"
        className="hover:text-violet-500"
        onClick={() => {
          setPos(pos)
          vscodeApi.postMessage({
            type: "session.fs.open",
            payload: {
              sid,
              pos,
              show: true,
              // todo, add start here
            }
          })
        }}
      />
      </Tooltipped>

      <Tooltipped label="terminal">
      <div className="hover:text-green-500">
        <SquareTerminal size={16}
          aria-label="terminal"
          onClick={(e) => {
            vscodeApi.postMessage({
              type: "session.term",
              payload: {
                sid,
                pos,
              }
            })
          }}
        />
      </div>
      </Tooltipped>

      <Tooltipped label="merge">
      <GitPullRequestCreateArrow size={16}
        aria-label="merge"
        className="hover:text-yellow-500"
        onClick={() => {
          vscodeApi.postMessage({
            type: "session.merge",
            payload: {
              sid,
              pos,
            }
          })
        }}
      />
      </Tooltipped>

      <Tooltipped label="fork">
      <GitGraph size={16}
        aria-label="fork"
        className="hover:text-sky-500"
        onClick={() => {
          vscodeApi.postMessage({
            type: "session.fork",
            payload: {
              from: sid,
              pos,
            }
          })
        }}
      />
      </Tooltipped>

      <Tooltipped label="refresh">
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
      </Tooltipped>

      <Tooltipped label="details">
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
      </Tooltipped>

    </div>
  )

}