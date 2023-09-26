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
		<input
			key="themeToggle"
			type="checkbox"
			className="toggle"
			checked={ resolvedTheme === "dark" }
			aria-label="Toggle dark mode"
			onClick={ () => handle() }
		/>
  )
}
