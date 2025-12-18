import { cn } from "@/lib/utils";
import { BrainCircuit, CheckCircle, Shell, Siren } from "lucide-react";

export const FuncCall = ({ part }:{ part: any, evt: any }) => {
  const fn = part.functionCall?.name as string
  const args = part.functionCall?.args
  const fnArgs = f2NameArgs[fn]
  const argVals = fnArgs?.map(a=> {
    if(a in args) {
      return args[a]
    }
  })
  const resp = part.functionResponse?.response

  const isPlanning = fn === "cache_put" && args["key"] === "planning"
  const isExec = fn === "exec"

  return (
    <div  className="mr-auto ml-1 pl-1 border-l-3 border-yellow-500 flex flex-col gap-2">
      <div  className="flex gap-2 items-baseline">
        <span className="font-heavy text-yellow-500">{fn}</span>
        <span className="font-thin">{(argVals || []).join(" ")}</span>
        { !resp && <BrainCircuit size={10} strokeWidth={1} className="text-yellow-300 animate-ping" />}
        { resp?.status === "ok" && <CheckCircle size={12} className="text-lime-500" />}
        { resp?.status === "error" && <Siren size={12} className="text-red-500" />}
      </div>
      { isPlanning && (
        <pre className="m-2 p-2 border border-violet-500">
          {args.value}
        </pre>
      )}
      { isExec && (
        <pre className="m-2 p-2 border border-green-500">
          {args.script}
        </pre>
      )}
    </div>
  )
}

// export const FuncResp = ({ part }:{ part: any, evt: any }) => {
//   // console.log("FuncResp.part", part)

//   const fn = part.functionResponse?.name as string
//   const resp = part.functionResponse?.response

//   const err = resp.error
//   // console.log("FuncResp.prep", fn, resp, err, !err)
//   var argVals: any[] = []
//   if (!err) {
//     const fnArgs = f2NameArgs[fn]
//     argVals = fnArgs?.map(a=> {
//       if(a in resp) {
//         return resp[a]
//       }
//       // return a
//     })
//   }

//   // console.log("FuncResp.render", fn, argVals, err)
//   return (
//     <div className="flex flex-col gap-2">
//       <div  className={cn("flex gap-2 items-baseline")}>
//         <span className="font-heavy">{fn}</span>
//         <span className="font-thin">{(argVals || []).join(" ")}</span>
//         { err && <span className="text-red-400">Error</span> }
//       </div>
//       { err && 
//         <div className="m-2 p-3 border border-red-400 max-h-32">
//           <pre className="overflow-auto">{err}</pre>
//         </div>
//       }
//     </div>
//   )
// }

const f2NameArgs: Record<string,string[]> = {
  "cache_write": ["key"],
  "cache_put": ["key"],
  "cache_edit": ["key"],
  "cache_del": ["key"],
  "cache_remove": ["key"],

  "fs_read": ["path"],
  "fs_list": ["path"],
  "fs_grep": ["path", "regexp"],
  "fs_write": ["path"],
  "fs_edit": ["path"],
  "fs_del": ["path"],

  "exec": ["key"],

  // legacy
  "read_file": ["path"],
  "read_dir":  ["path"],
  "tree_dir": ["path"],
  "write_file": ["path"],
  "cache_glob": ["path", "regexp"],
  "cache_grep": ["path", "regexp"],
  "cache_file": ["path"],
  "cache_dir": ["path"],
}