package st

Mask: {
  @task(st.Mask)
  val: _
  mask: _
  out: _
}

Pick: {
  @task(st.Pick)
  val: _
  pick: _
  out: _
}

Insert: {
  @task(st.Insert)
  val: _
  ins: _
  out: _
}

Replace: {
  @task(st.Pick)
  val: _
  repl: _
  out: _
}

Upsert: {
  @task(st.Upsert)
  val: _
  up: _
  out: _
}

Diff: {
  @task(st.Diff)
  orig: _
  next: _
  diff: _
}

Patch: {
  @task(st.Patch)
  orig: _
  patch: _
  next: _
}
