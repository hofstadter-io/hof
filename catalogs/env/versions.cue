// note how this intentionally does NOT match the dir name
package versions

// BIG version object for all the things
// 1. feeds down from subdirs (somehow)
// 2. we embed all the subdirs in a sibiling package

// N. we want the user to be able to create a custom versions value they
// can inject when they import, so they don't have to all the time
// i.e. the setup a local package with the version + embedding, and things should just work
// (make fractal on both side, here we should collect versions to here)
versions: {}
