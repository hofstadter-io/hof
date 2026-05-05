import { useState } from "react"
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import {Prism as SyntaxHighlighter} from 'react-syntax-highlighter'
import {vscDarkPlus} from 'react-syntax-highlighter/dist/esm/styles/prism'
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

export const TailwindClasses: string[] = [
  "prose prose-sm prose-invert",
  "prose-sm font-thin",
  "max-w-full",

  "prose-h1:my-[.5em]",
  "prose-h2:my-[.5em]",
  "prose-h3:my-[.5em]",

  "prose-pre:bg-[#1e1e1e] prose-pre:rounded prose-pre:border prose-pre:border-[#7e7e7e]",
  "prose-p:my-[.5em]",
  "prose-hr:my-[1em]",

  "prose-ol:list-decimal",
  "prose-ul:list-disc",
  "prose-li:mt-0",
]

const components = {
  code(props: any) {
    const [show, setShow] = useState(true)
    const {children, className, node, ...rest} = props
    const match = /language-(\w+)/.exec(className || '') // || "\n" in children?
    const lang = match ? match[1] : "text"

    const inline = typeof children === "string" && !String(children).includes("\n")

    return inline ? (
      <code {...rest} className={className}>
        {children}
      </code>
    ) : (
      // code block
      <div className={cn(
        "flex flex-col relative w-full bg-[#1e1e1e] [&>*]:bg-[#1e1e1e] veg-highlight [&>*]:veg-highlight",
        show ? "" : "max-h-64"
      )}>

        {/* actual code */}
        <SyntaxHighlighter
          {...rest}
          PreTag="div"
          children={String(children).replace(/\n$/, '')}
          language={lang}
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
    )
  }
}
