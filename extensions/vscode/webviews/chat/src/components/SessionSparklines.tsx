import React from 'react';
import Sparkline from './Sparkline';
import { processEvents } from "@/lib/utils";
import { UsageNumber } from './Info';

type SessionSparklinesProps = {
  events: any[];
}

export const SessionSparklines: React.FC<SessionSparklinesProps> = ({ events }) => {
  if (!events || events.length < 2) {
    return null;
  }

  const { merged } = processEvents(events);
  
  const cached: number[] = [];
  const prompt: number[] = [];
  const inputs: number[] = [];
  const thinks: number[] = [];
  const writes: number[] = [];
  const output: number[] = [];
  const totals: number[] = [];

  const verts: any[] = [];
  const ticks: any[] = [];
  const meta: any[] = [];
  let idx = 0;

  merged?.forEach((e, originalIndex) => {
    const u = e?.UsageMetadata
    const c = u?.cachedContentTokenCount || 0;
    const p = u?.promptTokenCount || 0;
    const i = p - c;
    const t = u?.thoughtsTokenCount || 0;
    const w = u?.candidatesTokenCount || 0;
    const o = t + w;
    const T = u?.totalTokenCount || 0;
    cached.push(c);
    prompt.push(p);
    inputs.push(i);
    thinks.push(t);
    writes.push(w);
    output.push(o);
    totals.push(T);

    const isUser = e.Author === "user";
    const hasUsage = e.Author !== "user" && e.UsageMetadata;
    const parts = e.Content?.parts || [];
    const hasFunc = parts.some((p: any) => p.functionCall);

    if (isUser) {
        verts.push({ index: idx, className: "stroke-sky-300/50" });
    }
    if (!isUser && parts.some((p: any) => p.text)) {
        verts.push({ index: idx, className: "stroke-lime-400/50" });
    }
    if (hasFunc) {
        const isError = parts.some((p: any) => 
            p.error || 
            p.functionCall?.error || 
            p.functionResponse?.error ||
            p.functionResponse?.response?.error ||
            p.functionResponse?.status === "error" ||
            p.functionResponse?.response?.status === "error"
        );

        ticks.push({ index: idx, className: isError ? "stroke-red-500" : "stroke-yellow-400" });
    }

    if (hasUsage) {
        let title = "Step";
        let titleColor = "";
        
        if (hasFunc) {
            const isError = parts.some((p: any) => 
                p.error || 
                p.functionCall?.error || 
                p.functionResponse?.error ||
                p.functionResponse?.response?.error ||
                p.functionResponse?.status === "error" ||
                p.functionResponse?.response?.status === "error"
            );
            if (isError) {
                titleColor = "text-red-500";
            }

            const part = parts.find((p: any) => p.functionCall);
            if (part?.functionCall?.name) {
                title = part.functionCall.name;
            } else {
                title = "Tool Use";
            }
        } else if (parts.some((p: any) => p.text)) {
             const part = parts.find((p: any) => p.text);
             if (part?.text) {
                 title = part.text.trim();
             } else {
                 title = "Model";
             }
        }
        
        meta.push({ index: originalIndex, title, className: titleColor });
    }
    idx++;
  });

  if (totals.length < 2) {
    return null
  }

  // console.log("sparkline input:", {
  //   verts,
  //   ticks,
  //   meta,
  // })

  const lines = [
    { value: 0, className: "stroke-white" },
    { value: 25000, className: "stroke-yellow-400" },
    { value: 50000, className: "stroke-orange-500" },
    { value: 100000, className: "stroke-red-500" }
  ];

  const allMetrics = [
    { title: "cached", values: cached, className: "text-lime-400" },
    { title: "inputs", values: inputs, className: "text-amber-200" },
    { title: "prompt", values: prompt, className: "text-amber-300" },
    { title: "thinks", values: thinks, className: "text-cyan-400" },
    { title: "writes", values: writes, className: "text-sky-400" },
    { title: "output", values: output, className: "text-blue-400" },
    { title: "totals", values: totals, className: "text-fuchsia-400" },
  ];

  return (
    <div className="flex ml-auto gap-4 h-12">
      <div className="w-full">
        <Sparkline
          lines={lines}
          verts={verts}
          ticks={ticks}
          meta={meta}
          tooltipData={allMetrics}
          formatter={UsageNumber}
          series={[
            {
              title: "totals",
              values: totals,
              className: "stroke-fuchsia-400 fill-fuchsia-300/5 stroke-2"
            },
            {
              title: "prompt",
              values: prompt,
              className: "stroke-amber-300 fill-amber-200/5"
            },
            {
              title: "cached",
              values: cached,
              className: "stroke-lime-400 fill-lime-300/10"
            },
            {
              title: "output",
              values: output,
              className: "stroke-blue-400"
            },
            {
              title: "thinks",
              values: thinks,
              className: "stroke-cyan-400"
            },
            {
              title: "writes",
              values: writes,
              className: "stroke-sky-400 fill-sky-300/20"
            }
          ]}
        />
      </div>
    </div>
  );
};
