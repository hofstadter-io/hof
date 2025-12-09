import { useEffect, useMemo, useState } from 'react'
import { cn } from '@/lib/utils'
import { useEditor, EditorContent, EditorContext } from '@tiptap/react'
import Document from '@tiptap/extension-document'
import StarterKit from '@tiptap/starter-kit'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import { TaskItem, TaskList } from '@tiptap/extension-list'
import { Markdown } from '@tiptap/markdown'


import { all, createLowlight } from 'lowlight'
const lowlight = createLowlight(all)

import Mention from '@tiptap/extension-mention'
import suggest from './suggest'

export const ChatEditor = ({
  chatState,
  onUpdate,
  editorRef,
}:{
  chatState?: any
  onUpdate?: any
  editorRef?: any
}) => {

  const editor = useEditor({
    editorProps: {
      attributes: {
        class: 'prose-invert prose-sm flex-grow flex flex-col h-full min-w-full min-h-64',
      },
    },
    

    onUpdate,
    extensions: [
      Markdown,
      StarterKit.configure({
        codeBlock: false,
      }),
      TaskList,
      TaskItem.configure({
        nested: true,
      }),
      CodeBlockLowlight.configure({ lowlight }),
      Mention.configure({
        HTMLAttributes: {
          class: 'suggest',
        },
        suggestions: [{
          char: '@',
          ...suggest.mentioner,
          items: ({ query }: { query: string }) => {
            const options = Object.keys(chatState?.config?.agents)
            const results = options
              .filter(item => item.toLowerCase().startsWith(query.toLowerCase()))
              .slice(0, 5)
            if (!results || results.length === 0) {
              return [query]
            }
            return results
          },
        },{
          char: '>',
          ...suggest.mentioner,
          items: ({ query }: { query: string }) => {
            const options = Object.keys(chatState?.config?.environs)
            const results = options
              .filter(item => item.toLowerCase().startsWith(query.toLowerCase()))
              .slice(0, 5)
            if (!results || results.length === 0) {
              return [query]
            }
            return results
          },
        },{
          char: '#',
          ...suggest.mentioner,
          items: ({ query }: { query: string }) => {
            // const options = Object.keys(chatState?.config?.runenv)
            const options = ["term-0", "term-1", "term-2", "main.go", "pkg/runtime.go", "src/App.tsx", "src/components/Events.tsx"]
            const results = options
              .filter(item => item.toLowerCase().startsWith(query.toLowerCase()))
              .slice(0, 5)
            if (!results || results.length === 0) {
              return [query]
            }
            return results
          },
        }]
      })
    ], // define your extension array
//     content: `
//         <p><span data-type="mention" data-id="veggie">@veggie</span> Can you help me with the error in <span data-type="mention" data-id="term-1" data-mention-suggestion-char="#">#term-1</span>?
//         <p>Use: <span data-type="mention" data-id="golang" data-mention-suggestion-char=">">&gt;golang</span> and <span data-type="mention" data-id="example.go" data-mention-suggestion-char="#">#example.go</span></p>
//         <h2>fix this... <b>now!</b></h2>
//         <pre><code class="language-go">package main

// func main() {
//     fmt.Println("hallo world!")
// }</code></pre>
//         <p>– Thanks, your human friend!</p>
//         <p><span data-type="mention" data-id="code_assist">@code_assist</span> double check this fool
//       `,

  })

  // Memoize the provider value to avoid unnecessary re-renders
  const providerValue = useMemo(() => ({ editor }), [editor])

  if (editorRef) {
    // console.log("REF:", editor, editorRef)
    editorRef.current = editor
  }

  return (
    <EditorContext.Provider value={providerValue}>
      <div className="flex-grow flex flex-col">
        <EditorContent editor={editor} />
      </div>
    </EditorContext.Provider>
  )
}