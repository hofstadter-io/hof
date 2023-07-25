import Hero from "@/src/components/ui/home/hero"
import Core from "@/src/components/ui/home/core"

export default function Page() {
  return (
		<div className="prose dark:prose-invert min-w-full p-4 flex flex-col">
			<Hero />
			<Core />
		</div>
  );
}
