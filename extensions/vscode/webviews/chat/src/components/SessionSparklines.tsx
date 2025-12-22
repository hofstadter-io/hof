import React from 'react';
import Sparkline from './Sparkline';
import { processEvents } from "@/lib/utils";

type SessionSparklinesProps = {
  events: any[];
  session: any;
  chatState: any;
}

export const SessionSparklines: React.FC<SessionSparklinesProps> = ({ events, session, chatState }) => {
  const { usages } = processEvents(events);
  
  const cached: number[] = [];
  const prompt: number[] = [];
  const inputs: number[] = [];
  const thinks: number[] = [];
  const writes: number[] = [];
  const output: number[] = [];
  const totals: number[] = [];

  usages?.forEach((u) => {
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
  });

  const lines = [
    { value: 0, className: "stroke-white" },
    { value: 25000, className: "stroke-yellow-400" },
    { value: 50000, className: "stroke-orange-500" },
    { value: 100000, className: "stroke-red-500" }
  ];

  if (!events || events.length === 0) {
    return null;
  }

  return (
    <div className="flex ml-auto gap-4 h-8">
      <div className="w-64">
        <Sparkline
          lines={lines}
          series={[
            {
              title: "prompt",
              values: prompt,
              className: "stroke-amber-300 fill-amber-200/5 stroke-2"
            },
            {
              title: "cached",
              values: cached,
              className: "stroke-lime-400 fill-lime-300/20"
            },
            {
              title: "thinks",
              values: thinks,
              className: "stroke-cyan-400 fill-cyan-300/20"
            },
            {
              title: "writes",
              values: writes,
              className: "stroke-blue-400 fill-blue-300/20"
            }
          ]}
        />
      </div>

      <div className="w-64">
        <Sparkline
          lines={lines}
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
              title: "output",
              values: output,
              className: "stroke-sky-400 fill-sky-300/20"
            }
          ]}
        />
      </div>
    </div>
  );
};
