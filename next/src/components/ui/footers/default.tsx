export default function Footer({ copyright }: { copyright: string }) {
	const today = new Date()
	const year = today.getFullYear();
  return (
    <footer className="footer flex w-100 min-w-screen h-12 justify-center items-center">
      <span>© {year} {copyright}</span>
    </footer>
  );
}
