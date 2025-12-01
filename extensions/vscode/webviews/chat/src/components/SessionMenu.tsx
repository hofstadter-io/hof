import { vscodeApi } from "@/vscodeApi";
import {
  Braces,
  FileDiff,
  GitGraph,
  GitPullRequestCreateArrow,
  ListTree,
  RefreshCcw,
  SquareTerminal,
} from 'lucide-react'

import { Tooltipped } from "@/components/Tooltipped";

export const Menu = ({
  sid,
  pos,
  setPos,
  hidden,
  refresh,
  // openFS,
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

  return (
    <div className="ml-auto flex justify-end items-center gap-2">

      <Tooltipped label="browse">
      <ListTree size={16}
        aria-label="diff"
        className="hover:text-sky-500"
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
              type: "session.term.open",
              payload: {
                sid,
                pos,
              }
            })
          }}
        />
      </div>
      </Tooltipped>

      <Tooltipped label="diff">
      <FileDiff size={16}
        aria-label="diff"
        className="hover:text-yellow-500"
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