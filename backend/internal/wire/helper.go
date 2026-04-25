package wire

import "tiny-forum/pkg/utils"

func NewTimeHelpers() *utils.TimeHelpers {
	// 使用默认的解析链（绝对时间 + 相对时间）
	defaultChain := utils.NewParserChain(
		utils.AbsoluteParser{},
		utils.NewRelativeParser(),
	)
	return &utils.TimeHelpers{
		SingleParser: utils.NewTimeParser(defaultChain),
		RangeParser:  utils.NewTimeRangeParser(defaultChain),
	}
}
