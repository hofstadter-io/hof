package phases

import (
	"github.com/hofstadter-io/hof/pkg/lang/context"
)

type Phase func (*context.Context) error
