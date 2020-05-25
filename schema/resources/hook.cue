package resources

#Hooks: [N=string]: #Hook & { Name: N, ... }
#Hook: {
	Name: string
	Evt: string
	Cmds: [...string]
}
