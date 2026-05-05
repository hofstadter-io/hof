import { useRef, useEffect, useMemo } from 'react'
import { cn } from '@/lib/utils'

import { vscodeApi } from '@/vscodeApi.js'
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
import { useChat } from '@/hooks/useChat'

export const ChatEditor = ({
  userInput,
  handlers,
  editorRef,
}: {
  userInput?: any
  handlers?: any
  editorRef?: any
}) => {
  const {
    chatState,
  } = useChat();

  const chatStateRef = useRef(chatState)
  useEffect(() => {
    chatStateRef.current = chatState
  }, [chatState])

  const editor = useEditor({
    editorProps: {
      attributes: {
        class: cn("flex-grow flex flex-col h-full min-w-full min-h-48", ...TailwindClasses),
      },
    },

    onFocus: () => {
      vscodeApi.postMessage({ type: 'requestSync' })
    },

    // need to parse markdown here
    // content: userInput.text || "",
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

      Mention.configure({
        HTMLAttributes: {
          class: 'suggest',
        },
        suggestion: {
          items: ({ query }: { query: string }) => {
            const options = Object.keys(chatStateRef.current?.config?.agents || {}).concat(Object.keys(chatStateRef.current?.config?.models || {}))
            return options
              .filter(item => item.toLowerCase().includes(query.toLowerCase()))
              .slice(0, 10)
          },
          ...suggest.mentioner,
        }
      }),
      Mention.extend({ name: 'environ' }).configure({
        HTMLAttributes: {
          class: 'suggest',
        },
        suggestion: {
          char: '>',
          items: ({ query }: { query: string }) => {
            const options = ["none"].concat(Object.values(chatStateRef.current?.config?.environs || {}).map((e: any) => e.name) || [])
            return options
              .filter(item => item.toLowerCase().includes(query.toLowerCase()))
              .slice(0, 10)
          },
          ...suggest.mentioner,
        }
      }),
      Mention.extend({ name: 'context' }).configure({
        HTMLAttributes: {
          class: 'suggest',
        },
        suggestion: {
          char: '#',
          items: ({ query }: { query: string }) => {
            const options: string[] = ["fresh"]
            if (chatStateRef.current?.terminals?.terminals) {
              chatStateRef.current.terminals.terminals.forEach((t: any) => options.push(`term-${t.id}`))
            }
            if (chatStateRef.current?.window?.tabgroups) {
              chatStateRef.current.window.tabgroups.forEach((tg: any) => {
                tg.tabs.forEach((t: any) => options.push(t.label))
              })
            }
            return options
              .filter(item => item.toLowerCase().includes(query.toLowerCase()))
              .slice(0, 10)
          },
          ...suggest.mentioner,
        }
      }),
      Mention.extend({ name: 'veg' }).configure({
        HTMLAttributes: {
          class: 'suggest',
        },
        suggestion: {
          char: '$',
          items: ({ query }: { query: string }) => {
            const options = ["state", "chat"]
            const results = options
              .filter(item => item.toLowerCase().includes(query.toLowerCase()))
              .slice(0, 5)
            if (!results || results.length === 0) {
              return [query]
            }
            return results
          },
          ...suggest.mentioner,
        }
      }),
      Mention.extend({ name: 'hist' }).configure({
        HTMLAttributes: {
          class: 'suggest',
        },
        suggestion: {
          char: '%',
          items: ({ query }: { query: string }) => {
            const options = ["rewind", "clone", "thread", "compact", "slice"]
            const results = options
              .filter(item => item.toLowerCase().includes(query.toLowerCase()))
              .slice(0, 5)
            if (!results || results.length === 0) {
              return [query]
            }
            return results
          },
          ...suggest.mentioner,
        }
      })
    ], // define your extension array
  }, [])

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