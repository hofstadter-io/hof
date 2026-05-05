import React, { forwardRef, useEffect, useImperativeHandle, useState } from 'react'
import { cn } from '@/lib/utils'

export const EmojiList = forwardRef((props: any, ref: any) => {
  const [selectedIndex, setSelectedIndex] = useState(0)

  const selectItem = (index: number) => {
    const item = props.items[index]

    if (item) {
      props.command({ name: item.name })
    }
  }

  const upHandler = () => {
    setSelectedIndex((selectedIndex + props.items.length - 1) % props.items.length)
  }

  const downHandler = () => {
    setSelectedIndex((selectedIndex + 1) % props.items.length)
  }

  const enterHandler = () => {
    selectItem(selectedIndex)
  }

  useEffect(() => setSelectedIndex(0), [props.items])

  useImperativeHandle(ref, () => {
    return {
      onKeyDown: (x: any) => {
        if (x.event.key === 'ArrowUp') {
          upHandler()
          return true
        }

        if (x.event.key === 'ArrowDown') {
          downHandler()
          return true
        }

        if (x.event.key === 'Enter' || x.event.key === 'Tab') {
          enterHandler()
          return true
        }

        return false
      },
    }
  }, [upHandler, downHandler, enterHandler])

  return (
    <div className={cn(
      "relative flex flex-col gap-1 p-1",
      "border rounded bg-stone-800",
      "max-h-40 overflow-y-auto"
    )}>
      {props.items.length ? (
        props.items.map((item: any, index: number) => (
          <button
            className={cn(
              "w-full flex gap-1 items-center",
              "bg-transparent hover:bg-stone-600 [&>*]:size-4",
              index === selectedIndex ?  'bg-stone-700' : '',
            )}
            key={index}
            onClick={() => selectItem(index)}
          >
            {item.fallbackImage ? <img src={item.fallbackImage} /> : item.emoji}:{item.name}:
          </button>
        ))
      ) : (
        <div className="item">No result</div>
      )}
    </div>
  )
})