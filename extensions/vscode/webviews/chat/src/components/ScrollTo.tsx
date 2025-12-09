import { ArrowUp, ArrowDown } from "lucide-react";
import React, { useState, useEffect, type Ref } from "react";
import { useStickToBottomContext } from 'use-stick-to-bottom';

export function ScrollTo({
  target,
}:{
  target: Ref<HTMLDivElement>,
}) {
  const { isAtBottom, scrollToBottom } = useStickToBottomContext();
  function scrollToTop() {
    // @ts-ignore
    target?.current?.scrollIntoView({
      behavior: "smooth"
    });
  }

  return (
    !isAtBottom && (
      <div className="absolute right-2 bottom-2 flex flex-col gap-2 text-white/30">
        <ArrowUp
          onClick={() => scrollToTop()}
          size={24}
          strokeWidth={1}
          className="rounded-4xl bg-fuchsia-500/30 hover:text-white hover:bg-fuchsia-500"
          aria-label="Scroll to top"
        />
        <ArrowDown
          onClick={() => scrollToBottom()}
          size={24}
          strokeWidth={1}
          className="rounded-4xl bg-fuchsia-500/30 hover:text-white hover:bg-fuchsia-500"
          aria-label="Scroll to bottom"
        />
      </div>
    )
  );
}
