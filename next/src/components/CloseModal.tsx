"use client"

import { useRouter }  from 'next/navigation';

const CloseModal = ({}) => {
	const router = useRouter();
	return (
		<button 
			onClick={() => router.back()}
			className="border-2 border-slate-600 rounded-sm px-2 prose">
			X
		</button>
	)
}

export default CloseModal
