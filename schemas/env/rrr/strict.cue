package env

#ImageLike: #Container | #HostImage | #DockerBuild
#DirLike:   #Dir | #HostDir | #GitRepo
#FileLike:  #File | #HostFile | #CuefigSBOM

#Container: {
	from: string | #ImageLike
}

#Service: {
	source: #ImageLike
}

#File: {
	source: #DirLike | #ImageLike
}

#Dir: {
	source: #DirLike | #ImageLike
}

#Patch: {
	source: #DirLike | #ImageLike
}

Dir: {
	source: #DirLike | #ImageLike
}

File: {
	content: string | #FileLike | #DirLike | #ImageLike
}

Secret: {
	content: string | #FileLike
}
