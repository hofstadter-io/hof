import Signin from '@/src/components/auth/signin'

export default function Page() {
	return (
		<div className="flex flex-col justify-start items-center w-screen">

			<div className="h-fit mt-20 py-20 px-8 prose dark:prose-inverted">
				<Signin />

			</div>
		</div>
	)
}