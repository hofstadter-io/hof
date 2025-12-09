import JsonView from '@uiw/react-json-view';
import { vscodeTheme } from '@uiw/react-json-view/vscode';

export const JsonObject = ({ hidden, data }: { hidden?: boolean, data: any }) => {
  if (hidden) return null

  return (
    <div className="transition whitespace-pre-line m-2 max-h-200 overflow-y-auto text-sm">
      <JsonView 
        value={data}
        style={vscodeTheme}
        className="p-2"

        displayDataTypes={false}
        indentWidth={12}
        shortenTextAfterLength={80}

        // collapsed={2}
        // @ts-ignore
        shouldExpandNodeInitially={(isExpanded, { value, keys, level }) => {
          if (level > 2) {
            return false
          }
          return true
          // if (keys.length > 0 && keys[0] == "object") {
          //   return false
          // }
          return isExpanded
        }}
      />
    </div>
  )
}
