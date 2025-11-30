
export const Welcome = ({
  username
}:{
  username: string
}) => {
  return (
    <div className="flex flex-col m-auto border border-sky-500/50 rounded-xl px-20 py-10 gap-5">
      <div className="flex gap-4 items-baseline">
        <span className="font-thin text-5xl"
        >Welcome</span>
        <span className="font-thin text-5xl text-sky-500"
        >{username}</span>
      </div>
      <span className="font-thin text-xl w-full text-center"
      >How can I help today?</span>
    </div>
  )
}