/*
Package configdir provides a cross platform means of detecting the system's
configuration directories.

This makes it easy to program your Go app to store its configuration files in
a standard location relevant to the host operating system. For Linux and some
other Unixes (like BSD) this means following the Freedesktop.org XDG Base
Directory Specification 0.8, and for Windows and macOS it uses their standard
directories.

This is a simple no-nonsense module that just gives you the paths, with
optional components tacked on the end for vendor- or app-specific namespacing.
It also provides a convenience function to call `os.MkdirAll()` on the paths to
ensure they exist and are ready for you to read and write files in.

Standard Global Configuration Paths

  * Linux: $XDG_CONFIG_DIRS or "/etc/xdg"
  * Windows: %PROGRAMDATA% or "C:\\ProgramData"
  * macOS: /Library/Application Support

Standard User-Specific Configuration Paths

  * Linux: $XDG_CONFIG_HOME or "$HOME/.config"
  * Windows: %APPDATA% or "C:\\Users\\%USER%\\AppData\\Roaming"
  * macOS: $HOME/Library/Application Support

Standard User-Specific Cache Paths

  * Linux: $XDG_CACHE_HOME or "$HOME/.cache"
  * Windows: %LOCALAPPDATA% or "C:\\Users\\%USER%\\AppData\\Local"
  * macOS: $HOME/Library/Caches

[1]: https://specifications.freedesktop.org/basedir-spec/basedir-spec-0.8.html
*/
package configdir
