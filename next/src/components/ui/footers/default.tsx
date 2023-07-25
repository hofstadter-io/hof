import classnames from 'classnames';

export default function Footer({ copyright }: { copyright: string }) {
	const today = new Date()
	const year = today.getFullYear();
  return (
    <footer className={classnames(
		  "footer flex w-100 min-w-screen h-12",
			"justify-center items-center",
			"hover:text-yellow-300",
			)}
		>
      <span>© {year} {copyright}</span>
    </footer>
  );
}
