"use client"

import React, { useState } from 'react';
import { signIn } from 'next-auth/react'

interface LoginFormProps extends React.HTMLAttributes<HTMLDivElement> {
	provider: string
}

const LoginForm = ({ provider, ...props } : LoginFormProps) => {

	const [isLoading, setIsLoading] = useState(false);

	const loginWith = async () => {
		setIsLoading(true)

		try {
			await signIn(provider)
		} catch (error) {
			console.error(error)
		} finally {
		  setIsLoading(false)
		}
	}

	return (
		<div {...props}>
			<button
				onClick={loginWith}
			>
				{provider}
			</button>
		</div>
	)

}

export default LoginForm
