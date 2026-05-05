import { useState } from "react";
import { JsonObject } from "veg-webview-common";
import { Menu } from '@/components/SessionMenu'
import { UsageInfo, TimeInfo } from '@/components/Info'
import CopyClipboardButton from '@/components/CopyClipboardButton'

export const EventDetails = ({pos, evt}:{pos: number, evt: any}) => {
  const [hidden, setHidden] = useState(true);

  var text = evt?.Content?.parts?.map((p: any) => p.text).join("\n\n")

  return (
    <div className="flex flex-col gap-1 mt-2">

      { !hidden && <JsonObject data={evt} /> }

      <div className="flex items-end">

        <div className="flex-grow flex flex-col gap-1 items-start">
          <span className="flex gap-2 items-center">
            <CopyClipboardButton size={16} strokeWidth={1.5} source={text}/>
            <span className="text-sky-500 font-thin text-[1.2em]">@{evt.Author}</span>
          </span>
          <UsageInfo evt={evt} size={16} />
        </div>

        <div className="flex flex-col gap-1 items-end">
          <span className="flex gap-2 font-thin items-center">
            <TimeInfo timestamp={evt.Timestamp} />
          </span>
          <Menu pos={pos} hidden={hidden} setHidden={setHidden} />
        </div>
        
      </div>
    </div>
  )
}
