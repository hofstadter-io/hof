import { cn } from '@/lib/utils';
import { PoundSterling } from 'lucide-react';
import React, { useRef, useEffect, useState } from 'react';

type SparklineSeries = {
    values: number[];
    title?: string;
    className?: string;
}

type SparklineLine = {
    value: number;
    title?: string;
    className?: string;
}

type SparklineVert = {
    index: number;
    className?: string;
}

type SparklineProps = {
    title?: string;
    series: SparklineSeries[];
    lines?: SparklineLine[];
    verts?: SparklineVert[];
    ticks?: SparklineVert[];
    tooltipData?: SparklineSeries[];
    formatter?: (val: number) => string;
    meta?: { index: number, title: string, className?: string }[];
    dims?: {
      // start values
      max?: number;
      min?: number;
      // clamp values
      maxMax?: number;
      minMin?: number;
    }
}

const Sparkline: React.FC<SparklineProps> = (props) => {
  // Reference to the container element that is changeable
  const containerRef = useRef<HTMLDivElement>(null);
  // Current dimensions for the chart
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
  
  // Tooltip state
  const [hoverIndex, setHoverIndex] = useState<number | null>(null);

  // Create effect to run when component mounts
  useEffect(() => {
    const container = containerRef.current;

    // Do nothing if there is no element
    if (!container) return;

    const updateDimensions = () => {
      setDimensions({
        width: container.clientWidth,
        height: container.clientHeight,
      });
    };

    // Set initial dimensions
    updateDimensions();

    // Create the observer and subscribe to container's changes
    const resizeObserver = new ResizeObserver(updateDimensions);
    resizeObserver.observe(container);

    return () => {
      // Disconnect when the component unmounts
      resizeObserver.disconnect();
    };
  }, []);

  // Construct the points string with the simple mapping
  // Scales are [0, length] -> [0, width] and [min, max] -> [0, height]
  const { width, height } = dimensions;


  // Calculate the min and max values of the data, across all the data
  var min = props.dims?.min || 0;
  var max = props.dims?.max || 0;
  var dataLength = 0;
  var badData = false
  props.series.forEach((d) => {
    // Don't render the chart for less than 2 points
    if (!d?.values || d.values.length < 2) {
      console.log("bad values for series", d)
      return null
    }
    min = Math.min(...d.values, min);
    max = Math.max(...d.values, max);
    dataLength = Math.max(dataLength, d.values.length);
  })
  // small buffer on points to reduce clipping
  max *= 1.1
  // clamp values
  if (props.dims?.minMin) min = Math.max(min, props.dims.minMin);
  if (props.dims?.maxMax) max = Math.min(max, props.dims.maxMax);

  const dims = { min, max, width, height, dataLength };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (dataLength < 2) return;
    
    // Calculate index from mouse position
    const rect = e.currentTarget.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const index = Math.round((x / width) * (dataLength - 1));
    
    // If index changed
    if (index !== hoverIndex) {
        setHoverIndex(index);
    }
  };

  const handleMouseLeave = () => {
      setHoverIndex(null);
  };

  const tooltipValues = props.tooltipData || props.series;
  const currentMeta = props.meta?.[hoverIndex || 0];
  const tooltipTitle = currentMeta 
      ? `[${currentMeta.index}] ${currentMeta.title}`
      : `[${hoverIndex}] Step`;

  return (
    <div 
        ref={containerRef} 
        style={{ width: '100%', height: '100%', position: 'relative' }}
        onMouseMove={handleMouseMove}
        onMouseLeave={handleMouseLeave}
    >
      {hoverIndex !== null && (
          <div className="absolute z-50 text-white font-mono px-3 py-2 rounded-md shadow-md text-[10px] border border-white/10 whitespace-nowrap pointer-events-none backdrop-blur-sm" 
               style={{ 
                   backgroundColor: 'rgba(9, 9, 11, 0.7)',
                   left: `${(hoverIndex / (dataLength - 1)) * 100}%`, 
                   top: '100%', 
                   marginTop: '4px',
                   transform: hoverIndex > dataLength / 2 ? 'translateX(-100%)' : 'translateX(0)',
               }}>
              <div className={cn("font-bold border-b border-white/20 pb-1 mb-1 w-[100px] truncate", currentMeta?.className)}>{tooltipTitle}</div>
              {tooltipValues.map((s, i) => (
                  <div key={i} className="flex justify-between gap-4">
                      <span className={s.className?.split(' ')[0] || ''}>{s.title}</span>
                      <span>{props.formatter ? props.formatter(s.values[hoverIndex]) : s.values[hoverIndex]}</span>
                  </div>
              ))}
          </div>
      )}
      <svg width={width} height={height} className="overflow-hidden">
        {/* <!-- a transparent glow that takes on the colour of the object it's applied to --> */}
        <filter id="glow">
            <feGaussianBlur stdDeviation="1.5" result="coloredBlur"/>
            <feMerge>
                <feMergeNode in="coloredBlur"/>
                <feMergeNode in="SourceGraphic"/>
            </feMerge>
        </filter>
        { props.lines?.map((l, i) => <Dashline key={i} line={l} dims={dims} />)}
        { props.verts?.map((v, i) => <VertLine key={i} vert={v} dims={dims} />)}
        { hoverIndex !== null && (
            <line
                x1={(hoverIndex / (dataLength - 1)) * width}
                x2={(hoverIndex / (dataLength - 1)) * width}
                y1={0}
                y2={height}
                className="stroke-white/50 stroke-1"
            />
        )}
        { props.series.map((s, i) => <Polyline key={i} series={s} dims={dims} />)}
        { props.ticks?.map((v, i) => <TickLine key={i} vert={v} dims={dims} />)}
      </svg>
    </div>
  );
};

type VertLineProps = {
  vert: SparklineVert;
  dims: {
    width: number;
    height: number;
    dataLength: number;
  };
}

const VertLine: React.FC<VertLineProps> = ({ vert, dims }) => {
    const x = (vert.index / (dims.dataLength - 1)) * dims.width;
    return (
        <line
            x1={x} x2={x}
            y1={0} y2={dims.height}
            strokeDasharray="2 2"
            className={cn(
                "fill-none stroke-1",
                vert.className,
            )}
        />
    )
}

const TickLine: React.FC<VertLineProps> = ({ vert, dims }) => {
    const x = (vert.index / (dims.dataLength - 1)) * dims.width;
    return (
        <line
            x1={x} x2={x}
            y1={0} y2={4}
            className={cn(
                "fill-none stroke-2",
                vert.className,
            )}
        />
    )
}

type PolylineProps = {
  series: SparklineSeries;
  dims: {
    width: number;
    height: number;
    min: number;
    max: number;
  };
}

const Polyline: React.FC<PolylineProps> = ({ series, dims }) => {
  // series points
  var points = series.values.map((value, index) => {
    const x = (index / (series.values.length - 1)) * dims.width;
    const y = ((value - dims.min) / (dims.max - dims.min)) * dims.height;

    var p = `${x},${dims.height - y}`
    // extra points to hide polygon for better filling
    if (index === 0) {
      p = `-42,${dims.height+42} -42,${dims.height - y} ${p}`
    }
    if (index === series.values.length-1) {
      p = `${p} ${dims.width+42},${dims.height - y} ${dims.width+42},${dims.height+42}`
    }
    return p;
  }).join(' ')

  // cap off the ends to make a polygon
  points = `` + points + ``

  return (
    <polygon
      className={cn(
        "fill-none stroke-1 stroke-current",
        series.className,
      )}
      filter="url(#glow)"
      points={points}
      strokeLinejoin="round"
    />
  )
}

type DashlineProps = {
  line: SparklineLine;
  dims: {
    width: number;
    height: number;
    min: number;
    max: number;
  };
}

const Dashline: React.FC<DashlineProps> = ({ line, dims }) => {
  const yo = ((line.value - dims.min) / (dims.max - dims.min)) * dims.height;
  // need to inverse here, polyline does this for us
  const y = dims.height - yo
  return (
    <line
      x1="0" x2={dims.width}
      y1={y} y2={y}
      strokeDasharray="2 3"
      className={cn(
        "fill-none stroke-1",
        line.className,
      )}
    />
  )
}

export default Sparkline;