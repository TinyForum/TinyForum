package wire

import (
	"tiny-forum/pkg/timeutil"
)

type Helpers struct {
	TimeHelpers *timeutil.TimeHelpers
	// 未来可以加其他: Logger, Cache 等
}

func NewHelpers() *Helpers {
	defaultChain := timeutil.NewParserTimeChain(
		timeutil.AbsoluteParser{},
		timeutil.NewRelativeParser(),
	)
	timeHelpers := &timeutil.TimeHelpers{
		SingleParser: timeutil.NewTimeParser(defaultChain),
		RangeParser:  timeutil.NewTimeRangeParser(defaultChain),
	}
	return &Helpers{TimeHelpers: timeHelpers}
}
