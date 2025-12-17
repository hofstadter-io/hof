import { computePosition } from '@floating-ui/dom'
import { ReactRenderer } from '@tiptap/react'

import { EmojiList } from './EmojiList'

export default {
  items: ({ editor, query }: { editor: any, query: string}) => {
    return editor.storage.emoji.emojis
      .filter(({ shortcodes, tags }: { shortcodes: any[], tags: any[] }) => {
        return (
          shortcodes.find(shortcode => shortcode.startsWith(query.toLowerCase())) ||
          tags.find(tag => tag.startsWith(query.toLowerCase()))
        )
      })
      // .slice(0, 10)
  },

  allowSpaces: false,

  render: () => {
    let component: any

    function repositionComponent(clientRect: any) {
      if (!component || !component.element) {
        return
      }

      const virtualElement = {
        getBoundingClientRect() {
          return clientRect
        },
      }

      computePosition(virtualElement, component.element, {
        placement: 'bottom-start',
      }).then(pos => {
        Object.assign(component.element.style, {
          left: `${pos.x}px`,
          top: `${pos.y}px`,
          position: pos.strategy === 'fixed' ? 'fixed' : 'absolute',
        })
      })
    }

    return {
      onStart: (props: any) => {
        component = new ReactRenderer(EmojiList, {
          props,
          editor: props.editor,
        })

        document.body.appendChild(component.element)
        repositionComponent(props.clientRect())
      },

      onUpdate(props: any) {
        component.updateProps(props)
        repositionComponent(props.clientRect())
      },

      onKeyDown(props: any) {
        if (props.event.key === 'Escape') {
          document.body.removeChild(component.element)
          component.destroy()

          return true
        }

        return component.ref?.onKeyDown(props)
      },

      onExit() {
        if (document.body.contains(component.element)) {
          document.body.removeChild(component.element)
        }
        component.destroy()
      },
    }
  },
}