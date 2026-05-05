import { useState } from "react"
import { CopyToClipboard } from "react-copy-to-clipboard-ts";

import {
  ClipboardCopy,
  ClipboardCheck,
} from 'lucide-react'

export default ({
  source,
  size = 24,
  strokeWidth = 1,
  positioning = "",
  copyClasses = "hover:text-sky-500",
  checkClasses = "text-green-500",
}:{
  source: string,
  size?: number,
  strokeWidth?: number,
  positioning?: string,
  copyClasses?: string,
  checkClasses?: string,
}) => {
  const [copied, setCopied] = useState(false)
  return (
    <CopyToClipboard text={source} onCopy={() => {
      setCopied(true)
      setTimeout(() => setCopied(false), 3000)
    }}>
      <span className={positioning}>
        { copied ?
          <ClipboardCheck size={size} strokeWidth={strokeWidth} className={checkClasses} /> 
        :   
          <ClipboardCopy size={size} strokeWidth={strokeWidth} className={copyClasses} /> 
        }
      </span>
    </CopyToClipboard>
  )
}

