import { type ComponentProps } from "react";

export const Button = (props: ComponentProps<'button'>) => {
  return (
    <button {...props} />
  )
}