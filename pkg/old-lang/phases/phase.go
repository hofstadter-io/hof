package phases

import (
	"github.com/hofstadter-io/hof/pkg/old-lang/context"
)

type Phase func (*context.Context) error
