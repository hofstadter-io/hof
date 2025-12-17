import { useChat } from "@/hooks/useChat"

export const Welcome = () => {
  const { chatState } = useChat()
  const username = chatState?.env?.user
  
  return (
    <div className="flex flex-col m-auto border border-sky-500/50 rounded-xl px-20 py-10 gap-5">
      <div className="flex gap-4 items-baseline">
        <span className="font-thin text-5xl"
        >Welcome</span>
        { username && <span className="font-thin text-5xl text-sky-500" >{username}</span> }
        { !username && <span className="font-thin text-5xl" >to <span className=" text-lime-500">Veggie</span></span> }
      </div>
      <span className="font-thin text-xl w-full text-center"
      >How can I help today?</span>
    </div>
  )
}