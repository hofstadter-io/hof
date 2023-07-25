'use client';

import React from 'react';
import Image from 'next/image'
import Link from '@/src/components/link';
import { redirect } from "next/navigation";
import { signOut } from 'next-auth/react'
import { useSession } from "next-auth/react";

export default function UserMenu ({ ThemeChanger }) {
	const { data: session, status } = useSession()
	/// console.log("session:", session)

	if (session?.user) {
		const signout = () => {
			signOut({
				callbackUrl: "/",
			})
		}

		return (
			<div className="dropdown dropdown-end">
				<label tabIndex={0} className="btn btn-ghost btn-circle avatar">
					<div className="w-8 h-8">
						<Image fill src={session?.user?.image} alt="profile picture" className="rounded-full"/>
					</div>
				</label>
				<ul tabIndex={0} className="mt-3 p-2 shadow menu menu-sm dropdown-content bg-base-100 rounded-box w-52">
					<li>
						<a href="/profile" className="justify-between">
							Profile
							<span className="badge badge-accent dark:badge-primary">New</span>
						</a>
					</li>
					<li><a>Settings</a></li>
					<li><ThemeChanger /></li>
					<li><div onClick={signout}>Sign Out</div></li>
				</ul>
			</div>
		)
	} else {
		return (
			<div className="text-white bg-emerald-600 rounded py-2 px-4 mx-2">
				 <Link href="/signin">Sign In</Link>
			</div>
		)
	}
}