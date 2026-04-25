package wire

import "tiny-forum/pkg/utils"

type TimeHelpers struct {
	SingleParser *utils.TimeParser
	RangeParser  *utils.TimeRangeParser
}

func NewTimeHelpers() *TimeHelpers {
	// 使用默认的解析链（绝对时间 + 相对时间）
	defaultChain := utils.NewParserChain(
		utils.AbsoluteParser{},
		utils.NewRelativeParser(),
	)
	return &TimeHelpers{
		SingleParser: utils.NewTimeParser(defaultChain),
		RangeParser:  utils.NewTimeRangeParser(defaultChain),
	}
}
