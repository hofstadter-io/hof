package st

Mask: {
	@task(st.Mask)
	$task: "st.Mask"
	val:   _
	mask:  _
	out:   _
}

Pick: {
	@task(st.Pick)
	$task: "st.Pick"
	val:   _
	pick:  _
	out:   _
}

Insert: {
	@task(st.Insert)
	$task:  "st.Insert"
	val:    _
	insert: _
	out:    _
}

Replace: {
	@task(st.Replace)
	$task:   "st.Replace"
	val:     _
	replace: _
	out:     _
}

Upsert: {
	@task(st.Upsert)
	$task:  "st.Upsert"
	val:    _
	upsert: _
	out:    _
}

Diff: {
	@task(st.Diff)
	$task: "st.Diff"
	orig:  _
	patch: diff
	next:  _
	diff:  _
}

Patch: {
	@task(st.Patch)
	$task: "st.Patch"
	orig:  _
	patch: _
	next:  _
	diff:  patch
}
