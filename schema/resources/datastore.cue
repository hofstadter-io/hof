package resources

#Datastores: [N=string]: #Datastore & { Name: N, ... }
#Datastore: {
  Name: string

	// TODO, should we specialize this with Cue and capabilities? the string gives flexibility to modules, and as a single token, we should be good, just need to be unique, can always prefix like go/cue/hof mods
	//  def should put a regex on the fields around this directory
  Type: string
  Auth: string
  Config: {...}
}
