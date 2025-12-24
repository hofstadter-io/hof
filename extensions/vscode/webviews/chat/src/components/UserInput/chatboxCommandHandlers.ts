import { vscodeApi } from '@/vscodeApi.js';

export const handleChatboxCommand = (text: string, sid: string, chatState: any, handlers: any) => {
  const lines = text.trim().split("\n")
  // is this a special thing
  if (lines.length === 1 && lines[0].startsWith("[@")) {
    const line = lines[0] as string
    const parts = line.split("[")
    parts.forEach((part, p) => {
      // probably the first part
      if (part === "") {
        return
      }
      // trim endings: [@ ... 
      part = part.substring(2).trim()
      const pairs = part.split(/\s+/)
      console.log("handling special:", part, pairs)
      var char = "@"
      var rest = ""
      var extra: string[] = []
      pairs.forEach((p) => {
        if (p.endsWith("]")) {
          p = p.substring(0, p.length-1)
        }
        if (p.startsWith("id")) {
          // trim endings: id="..."
          p = p.trim().substring(4, p.length-1).trim()
          rest = p
        } else if (p.startsWith("char")) {
          // trim endings: char="."
          p = p.trim().substring(6, p.length-1).trim()
          char = p[0]
        } else {
          extra.push(p)
        }
      })

      rest = rest.trim()
      console.log("parsed:", char, rest, extra)

      // / for slash commands for agents
      // ! should run a command in the environ

      switch (char) {
        case "@":
          // only send if known
          // autocomplete (substr/middle) partial
          if (Object.keys(chatState?.config?.agents).includes(rest)) {
            handlers.handleSelectAgent(rest)
          }
          if (Object.keys(chatState?.config?.models).includes(rest)) {
            handlers.handleSelectModel(rest)
          }
          break;

        case ">":
          // only send if known
          handlers.handleSelectEnviron(rest)
          break;

        case "#":
          // todo, send partial messages / events to attach, or store in userInput for when we do send a message
          break;

        case "$":
          console.log(`COMMAND$${rest}:`, extra)

          if (rest === "chat") {
            console.error("implement the chat command dummy!")
            const payload: any = { 
              focus: true, 
              agent: chatState?.userInput?.agent,
              model: chatState?.userInput?.model,
              envName: chatState?.userInput?.environ,
            }
            if (extra?.length > 0) {
              payload.environ = {
                srcUri: extra[0]
              }
              // assume remaining is title
              if (extra.length > 1) {
                payload.title = extra.splice(1).join(" ")
              }
            }

            vscodeApi.postMessage({
              type: 'session.create',
              payload,
            });

            return true
          }
          

          // only send if state
          if (rest === "state") {
            if (!extra || extra.length < 1) {
              // todo, let user know by showing a help message
              return
            }
            const key = extra[0]
            var val: any = null
            if (extra.length > 1) {
              val = extra.splice(1).join(" ")
            }

            console.log("state!", key, val)
            if (!val || val.trim() === "") {
              vscodeApi.postMessage({
                type: 'session.state.del',
                payload: {
                  sid,
                  key,
                },
              });
            } else {
              vscodeApi.postMessage({
                type: 'session.state.put',
                payload: {
                  sid,
                  key,
                  val,
                },
              });
            }
            vscodeApi.postMessage({
              type: 'session.get',
              payload: {
                sid,
              },
            });

            // we can't process any more
            return true
          }
          break;
      }
    }) // end of part loop

    // end of our special char handling, we should return true to indicate it was handled
    return true 
  }
  return false
}
