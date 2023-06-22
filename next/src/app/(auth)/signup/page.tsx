import Signup from '@/src/components/auth/signup'

export default function Page() {
	return (
		<div className="flex flex-col justify-start items-center w-screen">

			<div className="h-fit mt-20 py-20 px-8 prose dark:prose-inverted">
				<Signup />

			</div>
		</div>
	)
}
