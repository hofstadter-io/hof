import React, { useEffect, useImperativeHandle, useState } from 'react'

import { cn } from '@/lib/utils'

export default (props: any) => {
  const [selectedIndex, setSelectedIndex] = useState(0)

  const selectItem = (index: any) => {
    const item = props.items[index]

    if (item) {
      props.command({ id: item })
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

  useImperativeHandle(props.ref, () => ({
    onKeyDown: ({ event }: { event: any }) => {
      // if ((event.shiftKey && event.key === 'Tab') || event.key === "ArrowUp") {
      if (event.key === "ArrowUp") {
        upHandler()
        return true
      }

      if (event.key === "ArrowDown") {
        downHandler()
        return true
      }

      // if (event.key === 'Enter' || (event.shiftKey && (event.key == " " || event.code == "Space" || event.keyCode == 32 )) ) {
      if (event.key === "Enter" || event.key === "Tab" ) {
        enterHandler()
        return true
      }

      return false
    },
  }))

  return (
    <div className={cn(
      "relative flex flex-col gap-1 p-1",
      "border rounded bg-stone-800",
    )}>
      {props.items.length ? (
        props.items.map((item: any, index: number) => (
          <button
            className={cn(
              "w-full flex gap-1 items-center",
              "bg-transparent hover:bg-stone-600",
              index === selectedIndex ?  'bg-stone-700' : '',
            )}
            key={index}
            onClick={() => selectItem(index)}
          >
            {item}
          </button>
        ))
      ) : (
        <div className="item">No result</div>
      )}
    </div>
  )
}