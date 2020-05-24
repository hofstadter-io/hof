# ConfigDir for Go

This library provides a cross platform means of detecting the system's
configuration directories so that your Go app can store config files in a
standard location. For Linux and other Unixes (BSD, etc.) this means using the
[Freedesktop.org XDG Base Directory Specification][1] (0.8), and for Windows
and macOS it uses their standard directories.

This is a simple no-nonsense module that just gives you the path names to do
with as you please. You can either get the bare root config path, or get a
path with any number of names suffixed onto it for vendor- or
application-specific namespacing.

For the impatient, the directories this library can return tend to be like
the following:

|         | **System-wide Configuration**                       |
|---------|-----------------------------------------------------|
| Windows | `%PROGRAMDATA%` or `C:\ProgramData`                 |
| Linux   | `$XDG_CONFIG_DIRS` or `/etc/xdg`                    |
| macOS   | `/Library/Application Support`                      |
|         | **User-level Configuration**                        |
| Windows | `%APPDATA%` or `C:\Users\%USER%\AppData\Roaming`    |
| Linux   | `$XDG_CONFIG_HOME` or `$HOME/.config`               |
| macOS   | `$HOME/Library/Application Support`                 |
|         | **User-level Cache Folder**                         |
| Windows | `%LOCALAPPDATA%` or `C:\Users\%USER%\AppData\Local` |
| Linux   | `$XDG_CACHE_HOME` or `$HOME/.cache`                 |
| macOS   | `$HOME/Library/Caches`                              |

## Quick Start

```go
// A common use case is to get a private config folder for your app to
// place its settings files into, that are specific to the local user.
configPath := configdir.LocalConfig("my-app")
err := configdir.MakePath(configPath) // Ensure it exists.
if err != nil {
    panic(err)
}

// Deal with a JSON configuration file in that folder.
configFile := filepath.Join(configPath, "settings.json")
type AppSettings struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
var settings AppSettings

// Does the file not exist?
if _, err = os.Stat(configFile); os.IsNotExist(err) {
    // Create the new config file.
    settings = AppSettings{"MyUser", "MyPassword"}
    fh, err := os.Create(configFile)
    if err != nil {
        panic(err)
    }
    defer fh.Close()

    encoder := json.NewEncoder(fh)
    encoder.Encode(&settings)
} else {
    // Load the existing file.
    fh, err := os.Open(configFile)
    if err != nil {
        panic(err)
    }
    defer fh.Close()

    decoder := json.NewDecoder(fh)
    decoder.Decode(&settings)
}
```

## Documentation

Package documentation is available at
<https://godoc.org/github.com/kirsle/configdir>

## Author

Noah Petherbridge, [@kirsle](https://github.com/kirsle)

## License

MIT

[1]: https://specifications.freedesktop.org/basedir-spec/basedir-spec-0.8.html
[2]: https://github.com/shibukawa/configdir
