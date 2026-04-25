package wire

import "tiny-forum/pkg/utils"

type Helpers struct {
	TimeHelpers *utils.TimeHelpers
	// 未来可以加其他: Logger, Cache 等
}

func NewHelpers() *Helpers {
	defaultChain := utils.NewParserChain(
		utils.AbsoluteParser{},
		utils.NewRelativeParser(),
	)
	timeHelpers := &utils.TimeHelpers{
		SingleParser: utils.NewTimeParser(defaultChain),
		RangeParser:  utils.NewTimeRangeParser(defaultChain),
	}
	return &Helpers{TimeHelpers: timeHelpers}
}
