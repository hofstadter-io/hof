package kv

// in memory cue.Value storage
Mem: {
  @task(kv.Mem)
  $task: "kv.Mem"

  // key to store under 
  key: string

  // if specified, value to store
  // otherwise the value in memory as filled in
  val?: _

  // delete the key & value
  delete: bool | *false

  // boolean for if the value was loaded
  loaded: bool
}

// redis, etcd

// obj Stores (here?)
// mongo, couch, elastic
// s3 / gcs
