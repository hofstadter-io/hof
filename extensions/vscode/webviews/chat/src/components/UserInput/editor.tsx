import { useEffect, useMemo, useState } from 'react'
import { cn } from '@/lib/utils'

// tiptap
import { useEditor, EditorContent, EditorContext } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import { Table, TableRow, TableCell, TableHeader } from '@tiptap/extension-table'
import { TaskItem, TaskList } from '@tiptap/extension-list'
import { Markdown } from '@tiptap/markdown'
import Emoji, { gitHubEmojis } from '@tiptap/extension-emoji'
import Mention from '@tiptap/extension-mention'

import { all, createLowlight } from 'lowlight'
const lowlight = createLowlight(all)

import { TailwindClasses } from '../Markdown'

import suggest from './suggest'
import emojiSuggest from './suggest/emojiSuggest'

export const ChatEditor = ({
  userInput,
  chatState,
  handlers,
  editorRef,
}:{
  userInput?: any
  chatState?: any
  handlers?: any
  editorRef?: any
}) => {
  const editor = useEditor({
    editorProps: {
      attributes: {
        class: cn( "flex-grow flex flex-col h-full min-w-full min-h-48", ...TailwindClasses),
      },
    },
    

    onUpdate: handlers.handleInputUpdate,
    autofocus: true,
    extensions: [
      // package deal
      StarterKit.configure({
        codeBlock: false,
      }),

      // markdown
      Table,
      TableRow,
      TableCell,
      TableHeader,
      TaskList,
      TaskItem.configure({
        nested: true,
      }),
      Markdown.configure({
        markedOptions: { gfm: true },
      }),
      CodeBlockLowlight.configure({ lowlight }),
      Emoji.configure({
        emojis: gitHubEmojis,
        suggestion: emojiSuggest,
      }),

      // mentions
      Mention.configure({
        HTMLAttributes: {
          class: 'suggest',
        },
        suggestions: [{
          // agents / models
          char: '@',
          ...suggest.mentioner,
          items: ({ query }: { query: string }) => {
            const options = Object.keys(chatState?.config?.agents).concat(Object.keys(chatState?.config?.models))
            const results = options
              .filter(item => item.toLowerCase().startsWith(query.toLowerCase()))
              .slice(0, 5)
            if (!results || results.length === 0) {
              return [query]
            }
            return results
          },
        },{
          // environs (fs/exe)
          char: '>',
          ...suggest.mentioner,
          items: ({ query }: { query: string }) => {
            // todo, sessions here too?
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
          // context
          char: '#',
          ...suggest.mentioner,
          items: ({ query }: { query: string }) => {
            const options = ["term-0", "term-1", "term-2", "main.go", "pkg/runtime.go", "src/App.tsx", "src/components/Events.tsx"]
            const results = options
              .filter(item => item.toLowerCase().startsWith(query.toLowerCase()))
              .slice(0, 5)
            if (!results || results.length === 0) {
              return [query]
            }
            return results
          },
        },{
          // veg commands
          char: '$',
          // allowSpaces: true,
          ...suggest.mentioner,
          items: ({ query }: { query: string }) => {
            const options = ["state"]
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