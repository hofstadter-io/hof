import Link from '@/src/components/link'

import Login from '@/src/components/auth/login'

export default function Signin() {
	return (
		<div className="flex flex-col items-center">
			<h1>Sign up to Supacode</h1>

			<Login provider="google" className="m-2 p-2 px-8 border border-slate-800"/>

			<Link href="/signin" className="mt-8">Sign in to existing account</Link>
		</div>
	)
}