import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"

export const Tooltipped = ({ label, children }: { label: string, children: any}) => {
  return (
    <Tooltip>
      <TooltipTrigger>{children}</TooltipTrigger>
      <TooltipContent>
        <pre>{label}</pre>
      </TooltipContent>
    </Tooltip>
  )
}