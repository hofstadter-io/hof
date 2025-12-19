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
    case "fs_edit":
    case "fs_del":
    case "fs_list":
      return (
        <div className="text-red-500">
          <div>Filesystem error ({args.path}):</div>
          <pre className="whitespace-pre-wrap">{typeof error === 'string' ? error : JSON.stringify(error, null, 2)}</pre>
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
          <pre className="whitespace-pre-wrap">{typeof error === 'string' ? error : JSON.stringify(error, null, 2)}</pre>
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

  switch (name) {
    case "cache_put":
    case "cache_write":
    case "cache_edit":
    case "cache_del":
    case "cache_remove":
      return <div>Cache dummy: {args.key}</div>;

    case "fs_read":
    case "fs_write":
    case "fs_edit":
    case "fs_del":
    case "fs_list":
      return <div>Filesystem dummy: {args.path}</div>;

    case "fs_grep":
    case "fs_glob":
      return <div>Search dummy: {args.path} {args.regexp || args.glob}</div>;

    case "exec":
      return <div>Exec dummy: {args.script?.slice(0, 50)}...</div>;

    // Legacy / Others
    case "read_file":
    case "write_file":
    case "read_dir":
    case "tree_dir":
    case "cache_file":
    case "cache_dir":
      return <div>Legacy dummy: {args.path}</div>;

    default:
      return <pre className="whitespace-pre-wrap">{JSON.stringify(args, null, 2)}</pre>;
  }
}
