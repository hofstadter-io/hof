import { useState } from "react"
// import { vscodeApi } from './vscodeApi.js'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import {Prism as SyntaxHighlighter} from 'react-syntax-highlighter'
import {vscDarkPlus} from 'react-syntax-highlighter/dist/esm/styles/prism'
import { CopyToClipboard } from "react-copy-to-clipboard-ts";
import { cn } from '@/lib/utils'

import {
  Code,
} from 'lucide-react'

import CopyClipboardButton from './CopyClipboardButton'


export const Markdown = ({children}:{children: any}) => {
  return (
    <ReactMarkdown
      children={children}
      components={components}
      remarkPlugins={[[remarkGfm, {singleTilde: false}]]}
    />
  )
}

const components = {
  code(props: any) {
    const [show, setShow] = useState(true)
    const {children, className, node, ...rest} = props
    const match = /language-(\w+)/.exec(className || '') // || "\n" in children?
    return match ? (
      // code block?
      <div className={cn(
        "flex flex-col relative w-full bg-[#1e1e1e] [&>*]:bg-[#1e1e1e] veg-highlight [&>*]:veg-highlight",
        show ? "" : "max-h-64"
      )}>

        {/* actual code */}
        <SyntaxHighlighter
          {...rest}
          PreTag="div"
          children={String(children).replace(/\n$/, '')}
          language={match[1]}
          style={vscDarkPlus}
          // className="bg-stone-800"
          codeTagProps={{
            className: "not-prose bg-[#1e1e1e]"
          }}
          customStyle={{
            lineHeight: "1.4"
          }}
        />

        {/* buttons */}
        <CopyClipboardButton size={20} source={children} positioning="absolute top-[-6px] left-[-12px] p-1"/>

        <div className="flex justify-center">
          <div className="bg-fuchsia-500/20 hover:bg-fuchsia-500/70 w-12 rounded flex justify-center" onClick={() => setShow(!show)}>
            <Code size={16} strokeWidth={1.5} /> 
          </div>
        </div>

      </div>
    ) : (
      // inline?
      <code {...rest} className={className}>
        {children}
      </code>
    )
  }
}
