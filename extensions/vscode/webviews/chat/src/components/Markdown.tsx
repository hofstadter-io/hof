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
  ClipboardCopy,
  ClipboardCheck,
} from 'lucide-react'


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
        "flex flex-col relative",
        show ? "" : "max-h-64"
      )}>
        <CopyButton source={children} />
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
            lineHeight: "1.3"
          }}
        />
        <div className="bg-slate-700 hover:bg-slate-500 w-full flex justify-center" onClick={() => setShow(!show)}>
          <Code size={20} strokeWidth={2} /> 
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

export const CopyButton = ({source}:{source: string}) => {
  const [copied, setCopied] = useState(false)
  const common = "p-1"
  return (
    <CopyToClipboard text={source} onCopy={() => {
      setCopied(true)
      setTimeout(() => setCopied(false), 3000)
    }}>
      <span className="absolute z-50 top-2 right-1">
        { copied ?
          <ClipboardCheck size={32} strokeWidth={1} className={cn(common, "text-green-500")} /> 
        :   
          <ClipboardCopy size={32} strokeWidth={1} className={cn(common, "hover:text-sky-500")} /> 
        }
      </span>
    </CopyToClipboard>
  )
}

