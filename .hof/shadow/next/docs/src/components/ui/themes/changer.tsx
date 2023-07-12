'use client';

import { useTheme } from '@wits/next-themes';

export default function ThemeChanger () {
  const { resolvedTheme, setTheme } = useTheme()

	const handle = function() {
		if (resolvedTheme === "light") {
			setTheme("dark")
		} else {
			setTheme("light")
		}
	}

  return (
		<span
			aria-label="Toggle dark mode"
			onClick={ () => handle() }
		>
			Toggle dark mode
		</span>
  )
}