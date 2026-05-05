import { useChat } from "@/hooks/useChat"

export const Welcome = ({
  title = "Welcome",
  subtitle = "How can I help today?",
}:{
  title?: string
  subtitle?: string
}) => {
  const { chatState } = useChat()
  const username = chatState?.env?.user
  
  return (
    <div className="flex flex-col m-auto border border-sky-500/50 rounded-xl px-20 py-10 gap-5 justify-stretch text-center">
      <div className="sm:flex gap-4 item-center sm:items-baseline w-full">
        <p className="font-thin text-4xl sm:text-5xl mb-4 sm:mb-0">{title}</p>
        { username && <p className="font-thin text-5xl text-sky-500" >{username}</p> }
        { !username && <p className="font-thin text-5xl" >to <span className=" text-lime-500">Veggie</span></p> }
      </div>
      <span className="font-thin text-xl w-full">{subtitle}</span>
    </div>
  )
}