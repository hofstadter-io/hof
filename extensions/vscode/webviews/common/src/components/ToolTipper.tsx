import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "./ui/tooltip"

export const ToolTipper = ({ label, side="top", children }: { label: string, side?: any, children: any}) => {
  return (
    <Tooltip>
      <TooltipTrigger>{children}</TooltipTrigger>
      <TooltipContent side={side}>
        <pre>{label}</pre>
      </TooltipContent>
    </Tooltip>
  )
}