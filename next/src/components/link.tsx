import Link from 'next/link';

export default function link({href, ...props}) { 
	return ( <Link href={href} prefetch={false} {...props} /> )
}
