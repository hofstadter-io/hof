'use client';

import Signup from '@/src/components/auth/signup'
import CloseModal from '@/src/components/CloseModal'

export default function Modal() {
	console.log("MODAL")
	return (
			<div className="z-50 fixed inset-0">
				<div className="container flex items-center h-full max-w-lg mx-auto">

				<div className="relative w-full h-fit py-20 px-8 rounded-lg prose bg-slate-300">
					<div className="absolute top-4 right-4">
						<CloseModal />
					</div>
					<Signup />
				</div>

				</div>
			</div>
	)
}