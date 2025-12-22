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

type SparklineProps = {
    title?: string;
    series: SparklineSeries[];
    lines?: SparklineLine[];
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
  props.series.forEach((d) => {
    // Don't render the chart for less than 2 points
    if (!d?.values || d.values.length < 2) {
      console.log("bad values for series", d)
      return null
    }
    min = Math.min(...d.values, min);
    max = Math.max(...d.values, max);
  })
  // small buffer on points to reduce clipping
  max += 1000
  // clamp values
  if (props.dims?.minMin) min = Math.max(min, props.dims.minMin);
  if (props.dims?.maxMax) max = Math.min(max, props.dims.maxMax);

  const dims = { min, max, width, height };

  return (
    <div ref={containerRef} style={{ width: '100%', height: '100%' }}>
      <svg width={width} height={height} className="overflow-hidden">
        {/* <!-- a transparent glow that takes on the colour of the object it's applied to --> */}
        <filter id="glow">
            <feGaussianBlur stdDeviation="1.5" result="coloredBlur"/>
            <feMerge>
                <feMergeNode in="coloredBlur"/>
                <feMergeNode in="SourceGraphic"/>
            </feMerge>
        </filter>
        { props.lines?.map(l => <Dashline line={l} dims={dims} />)}
        { props.series.map(s => <Polyline series={s} dims={dims} />)}
      </svg>
    </div>
  );
};

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