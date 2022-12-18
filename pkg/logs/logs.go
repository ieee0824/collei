package logs

import "github.com/ieee0824/collei/pkg/aggregator"

type Log map[string]any
type Tag string
type Logs map[Tag]aggregator.Aggregator[Log]
