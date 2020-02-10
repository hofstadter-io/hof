package phases

import (
	"github.com/hofstadter-io/hof/pkg/context"
)

type Phase func (*context.Context) error
