package templates

import (
	"strings"
)

//
// Template delimiters
//
//   these are for advanced usage, you shouldn't have to modify them normally

type Config struct {

	// "global" Template parameters
	TemplateSystem string  // which system ['text/template'(default), 'mustache']

  // Alt and Swap Delims,
	//   becuase the defaulttemplate systems use `{{` and `}}`
	//   and you may choose to use other delimiters, but the lookup system is still based on the template system
  //   and if you want to preserve those, we need three sets of delimiters
  AltDelims  bool
  SwapDelims bool

  // The default delimiters
  // You should change these when using alternative style like jinjas {% ... %}
  // They also need to be different when using the swap system
  LHS2_D string
  RHS2_D string
  LHS3_D string
  RHS3_D string

  // These are the same as the default becuase
  // the current template systems require these.
  //   So these should really never change or be overriden until there is a new template system
  //     supporting setting the delimiters dynamicalldelimiters dynamicallyy
  LHS2_S string
  RHS2_S string
  LHS3_S string
  RHS3_S string

  // The temporary delims to replace swap with while also swapping
  // the defaults you set to the swap that is required by the current templet systems
	// You need this when you are double templating a file and the top-level system is not the default
  LHS2_T string
  RHS2_T string
  LHS3_T string
  RHS3_T string
}

func (D *Config) SwitchBefore(content string) string {

	// Multi switch with temporary
	if D.AltDelims && D.SwapDelims {
		// Replace the swap or secondary with temp (this is the default for the template system)
		content = strings.ReplaceAll(content, D.LHS3_S, D.LHS3_T)
		content = strings.ReplaceAll(content, D.RHS3_S, D.RHS3_T)
		content = strings.ReplaceAll(content, D.LHS2_S, D.LHS2_T)
		content = strings.ReplaceAll(content, D.RHS2_S, D.RHS2_T)
	}

	if D.AltDelims {
		// Switch Swap for Default, which if you only set default, will work
		// do triple first, douvle second
		content = strings.ReplaceAll(content, D.LHS3_D, D.LHS3_S)
		content = strings.ReplaceAll(content, D.RHS3_D, D.RHS3_S)
		content = strings.ReplaceAll(content, D.LHS2_D, D.LHS2_S)
		content = strings.ReplaceAll(content, D.RHS2_D, D.RHS2_S)
	}

	return content
}

func (D *Config) SwitchAfter(content string) string {

	// Multi switch undo
	if D.AltDelims && D.SwapDelims {
		// Undo the default to temp swap
		content = strings.ReplaceAll(content, D.LHS3_T, D.LHS3_S)
		content = strings.ReplaceAll(content, D.RHS3_T, D.RHS3_S)
		content = strings.ReplaceAll(content, D.LHS2_T, D.LHS2_S)
		content = strings.ReplaceAll(content, D.RHS2_T, D.RHS2_S)

	}

	// shouldn't have to do anything since we rendered by now
	//   and there should have been only one template system
	//   or  another that is not effected by renering

	return content
}

// Override this Delim's dot defaults with 'delim' values
func (D *Config) OverrideDotDefaults(delim *Config) {

	if D.TemplateSystem == "." {
		D.TemplateSystem = delim.TemplateSystem
	}

  D.AltDelims = D.AltDelims || delim.AltDelims
  D.SwapDelims = D.SwapDelims || delim.SwapDelims

	if D.LHS2_D == "." {
		D.LHS2_D = delim.LHS2_D
	}
	if D.RHS2_D == "." {
		D.RHS2_D = delim.RHS2_D
	}
	if D.LHS3_D == "." {
		D.LHS3_D = delim.LHS3_D
	}
	if D.RHS3_D == "." {
		D.RHS3_D = delim.RHS3_D
	}

	if D.LHS2_S == "." {
		D.LHS2_S = delim.LHS2_S
	}
	if D.RHS2_S == "." {
		D.RHS2_S = delim.RHS2_S
	}
	if D.LHS3_S == "." {
		D.LHS3_S = delim.LHS3_S
	}
	if D.RHS3_S == "." {
		D.RHS3_S = delim.RHS3_S
	}

	if D.LHS2_T == "." {
		D.LHS2_T = delim.LHS2_T
	}
	if D.RHS2_T == "." {
		D.RHS2_T = delim.RHS2_T
	}
	if D.LHS3_T == "." {
		D.LHS3_T = delim.LHS3_T
	}
	if D.RHS3_T == "." {
		D.RHS3_T = delim.RHS3_T
	}

}
