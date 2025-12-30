@extern(embed)
package patches

import "strings"

_diffs: _ @embed(glob=*.diff,type=text)

for path, text in _diffs {
	let p = strings.TrimSuffix(path, ".diff")
	(p): text
}
