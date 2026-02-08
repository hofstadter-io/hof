import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function processEvents(events: any[]):{
  merged: any[],
  usages: any[],
} {
  // TODO, extract this
  // console.log("Events.events", events)
  // todo, coalesce events here
  var merged: any[] = []
  var usages: any[] = []
  events?.forEach((E1: any, e1: number) => {
    if (E1?.Author !== "user" && E1?.UsageMetadata) {
      usages.push(E1.UsageMetadata)
    }
    var didMerge: boolean = false
    // loop over earlier events
    for (var e2 = e1-1; e2 >= 0; e2--) {
      const E2 = merged[e2]
      const P2 = E2?.Content?.parts
      if (!P2) {
        continue
      }
      // if we have a matching invocation id, lets do some matching
      if (E1?.InvocationID === E2?.InvocationID) {
        const P1 = E1?.Content?.parts
        if (!P1) {
          continue
        }
        // todo, we need to loop over parts here too
        for (const p1 of P1) {
          if (!p1?.functionResponse) {
            continue
          }
          // console.log("functionCall", pi[0], ei)
          for (const pi2 in P2) {
            const p2 = P2[pi2]
            if (!p2?.functionCall || p1.functionResponse.id !== p2.functionCall.id) {
              continue
            }
            E2.Content.parts[pi2].functionResponse = p1.functionResponse
            didMerge = true
            // console.log("MATCH!", ej.InvocationID, ej, merged[j], pi[0], pj[0])
            break;
          }
        }
        merged[e2] = E2
        break;
      }
    }
    if (!didMerge) {
      merged.push(E1)
    }
  })


  return {
    merged,
    usages,
  }
}