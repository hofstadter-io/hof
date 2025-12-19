import { ToolTipper } from "veg-webview-common";

export const LightDetails = ({ evt }: { evt: any }) => {
  const parts = evt?.Content?.parts || [];
  const calls = parts.filter((p: any) => p.functionCall);
  
  if (calls.length === 0) return null;

  return (
    <div className="flex flex-col gap-2 mb-2">
      {calls.map((p: any, i: number) => (
        <div key={i} className="p-2 border border-dashed border-muted-foreground/50 text-xs">
          <div className="font-bold text-yellow-500 mb-1 flex items-center gap-2">
            <ToolTipper label={JSON.stringify(p, null, 2)} >
              {p.functionCall.name}
            </ToolTipper>
          </div>
          {p.functionResponse?.response.status === "error" ? (
            <ToolSpecificError part={p} evt={evt} />
          ) : (
            <ToolSpecificDetails part={p} evt={evt} />
          )}
        </div>
      ))}
    </div>
  )
}

const ToolSpecificError = ({ part, evt }: { part: any, evt: any }) => {
  const name = part.functionCall.name;
  const args = part.functionCall.args;
  const error = part.functionResponse.response.error;

  switch (name) {
    case "cache_put":
    case "cache_write":
    case "cache_edit":
    case "cache_del":
    case "cache_remove":
      return (
        <div className="text-red-500">
          <div>Cache error ({args.key}):</div>
          <pre className="whitespace-pre-wrap">{typeof error === 'string' ? error : JSON.stringify(error, null, 2)}</pre>
        </div>
      );

    case "fs_read":
    case "fs_write":
    case "fs_del":
    case "fs_list":
      return (
        <div className="text-red-500">
          <div>Filesystem error ({args.path}):</div>
          <pre className="whitespace-pre-wrap">{typeof error === 'string' ? error : JSON.stringify(error, null, 2)}</pre>
        </div>
      );

    case "fs_edit":
      return (
        <div className="text-red-500">
          <div>Filesystem error ({args.path}):</div>
          <pre className="whitespace-pre-wrap mb-2">{typeof error === 'string' ? error : JSON.stringify(error, null, 2)}</pre>
          <FsEditDetails path={args.path} edits={args.edits} />
        </div>
      );

    case "fs_grep":
    case "fs_glob":
      return (
        <div className="text-red-500">
          <div>Search error ({args.path}):</div>
          <pre className="whitespace-pre-wrap">{typeof error === 'string' ? error : JSON.stringify(error, null, 2)}</pre>
        </div>
      );

    case "exec":
      return (
        <div className="text-red-500">
          <div>Exec error:</div>
          <pre className="whitespace-pre-wrap mb-2">{typeof error === 'string' ? error : JSON.stringify(error, null, 2)}</pre>
          <ExecDetails args={args} response={part.functionResponse.response} />
        </div>
      );

    case "read_file":
    case "write_file":
    case "read_dir":
    case "tree_dir":
    case "cache_file":
    case "cache_dir":
      return (
        <div className="text-red-500">
          <div>Legacy error ({args.path}):</div>
          <pre className="whitespace-pre-wrap">{typeof error === 'string' ? error : JSON.stringify(error, null, 2)}</pre>
        </div>
      );

    default:
      return (
        <div className="text-red-500">
          <div className="font-semibold underline">Error in {name}:</div>
          <pre className="whitespace-pre-wrap">{typeof error === 'string' ? error : JSON.stringify(error, null, 2)}</pre>
        </div>
      );
  }
}

const ToolSpecificDetails = ({ part, evt }: { part: any, evt: any }) => {
  const name = part.functionCall.name;
  const args = part.functionCall.args;
  const response = part.functionResponse?.response;
  const output = response?.output;

  switch (name) {
    case "cache_put":
    case "cache_write":
    case "cache_edit":
    case "cache_del":
    case "cache_remove":
      return <div>Cache: {args.key}</div>;

    case "fs_read":
    case "fs_write":
    case "fs_del":
      return <div>Filesystem: {args.path}</div>;

    case "fs_list":
      return (
        <div>
          <div>Filesystem: {args.path}</div>
          {output && <pre className="mt-1 text-[10px] font-mono whitespace-pre-wrap max-h-40 overflow-auto border-t border-muted-foreground/30 pt-1">{output}</pre>}
        </div>
      );

    case "fs_edit":
      return <FsEditDetails path={args.path} edits={args.edits} />;

    case "fs_grep":
    case "fs_glob":
      return (
        <div>
          <div>Search: {args.path || args.glob} {args.regexp || ""}</div>
          {output && <pre className="mt-1 text-[10px] font-mono whitespace-pre-wrap max-h-40 overflow-auto border-t border-muted-foreground/30 pt-1">{output}</pre>}
        </div>
      );

    case "exec":
      return <ExecDetails args={args} response={response} />;

    // Legacy / Others
    case "read_file":
    case "write_file":
    case "read_dir":
    case "tree_dir":
    case "cache_file":
    case "cache_dir":
      return <div>Legacy: {args.path}</div>;

    default:
      return <pre className="whitespace-pre-wrap">Default: {JSON.stringify(args, null, 2)}</pre>;
  }
}

const FsEditDetails = ({ path, edits }: { path: string, edits: any[] }) => {
  const items = edits || [];
  return (
    <div className="flex flex-col gap-1">
      {items.map((edit: any, i: number) => (
        <div key={i} className="flex flex-col">
          <div className="font-semibold text-yellow-600 mb-1">
            ({edit.count}) {path}
          </div>
          <div className="grid grid-cols-2 gap-2 text-[10px] font-mono whitespace-pre-wrap border border-muted-foreground/50 p-1">
            <div className="text-red-400 border-r border-dashed border-muted-foreground/50 pr-2">{edit.old}</div>
            <div className="text-green-400 pl-2">{edit.new}</div>
          </div>
        </div>
      ))}
    </div>
  );
}

const ExecDetails = ({ args, response }: { args: any, response: any }) => {
  return (
    <div>
      <div className="font-mono text-[10px] mb-1">
        Exec: {args.script?.slice(0, 100)}...
        {response?.exitCode !== undefined && (
          <span className={`ml-2 px-1 rounded ${response.exitCode === 0 ? 'bg-green-500/20 text-green-500' : 'bg-red-500/20 text-red-500'}`}>
            exit: {response.exitCode}
          </span>
        )}
      </div>
      {response?.stdout && (
        <div className="mt-1">
          <div className="text-[9px] opacity-50 uppercase">stdout</div>
          <pre className="text-[10px] font-mono whitespace-pre-wrap max-h-40 overflow-auto border border-muted-foreground/30 p-1">{response.stdout}</pre>
        </div>
      )}
      {response?.stderr && (
        <div className="mt-1">
          <div className="text-[9px] opacity-50 uppercase text-yellow-500">stderr</div>
          <pre className="text-[10px] font-mono whitespace-pre-wrap max-h-40 overflow-auto border border-yellow-500/30 p-1 text-yellow-500/80">{response.stderr}</pre>
        </div>
      )}
    </div>
  );
}
