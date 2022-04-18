package gen

// Seed the Range
Seed: {
  @task(gen.Seed)
  $task: "gen.Seed"

  // only set to ensure consistent output while testing 
  seed?: int // defaults to time.Now()
}

Int: {
  @task(gen.Int)
  $task: "gen.Int)
  max?: int // max value if set

  // the random val returned
  val: int
}

Str: {
  @task(gen.Str)
  $task: "gen.Str"

  // number of runes to generate
  n: int | *12

  // possible runes, defaults to [a-zA-Z]
  runes?: string
}

// the other tasks don't really have schema or input

// c: string @task(gen.CUID)  // like UUID, but for cloud
// f: float  @task(gen.Float)
// n: float  @task(gen.Norm)
// n: string @task(gen.Now)   // RFC-3339
// s: string @task(gen.Slug)  // related to CUID
// u: string @task(gen.UUID)

