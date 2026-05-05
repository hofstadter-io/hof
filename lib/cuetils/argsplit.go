package cuetils

// splits args at the first '%'such that everything after is for CUE
func PercentSplitArgs(orig []string) (args, cueargs []string) {
	args = orig
	for i, a := range orig {
		if a == "%" {
			args = orig[:i]
			if i+1 < len(orig) {
				cueargs = orig[i+1:]
			}
			break
		}
	}
	return args, cueargs
}
