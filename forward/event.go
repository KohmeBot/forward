package forward

import (
	"fmt"
	"github.com/kohmebot/pkg/chain"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"regexp"
	"strconv"
	"strings"
)

func (p *PluginForward) SetOnStart(engine *zero.Engine) {

	engine.OnCommandGroup([]string{"forward", "fw", "f"}, zero.OnlyPrivate, p.env.SuperUser().Rule()).
		Handle(func(ctx *zero.Ctx) {
			sender := ctx.Event.Sender.ID
			if p.tmp.Has(sender) {
				ctx.Send(message.Text("已经开始传话了哦！接下来说的话我都会送达的！"))
				return
			}
			p.tmp.Start(sender, ctx)
			var msg chain.MessageChain
			var groups []string
			for group := range p.env.Groups().RangeGroup {
				groups = append(groups, strconv.FormatInt(group, 10))
			}

			msg.Split(
				message.Text("现在开始传话啦！"),
				message.Text(fmt.Sprintf("接下来每说的一句话，将转发给 %s", strings.Join(groups, ","))),
				message.Text(fmt.Sprintf("在%d秒内如果没有收到新的消息，将会退出传话！", p.conf.StopDur)),
			)

			ctx.Send(msg)

			msg = chain.MessageChain{}
			msg.Split(
				message.Text("特别提醒！如果有艾特需求则按照以下规则"),
				message.Text("艾特全体: 输入 @全体"),
				//message.Text("艾特某人: 输入 @QQ号"),
			)

			ctx.Send(msg)

		})
}

func (p *PluginForward) SetOnMsg(engine *zero.Engine) {
	engine.OnMessage(zero.OnlyPrivate, p.env.SuperUser().Rule()).Handle(func(ctx *zero.Ctx) {
		sender := ctx.Event.Sender.ID
		if !p.tmp.Has(sender) {
			return
		}
		p.tmp.Refresh(sender)
		conv := convChain(ctx.Event.Message)
		for group := range p.env.Groups().RangeGroup {
			ctx.SendGroupMessage(group, conv)
		}
		ctx.Send(message.Text("已传话！"))
	})
}

func convChain(msgs message.Message) message.Message {
	cMsg := chain.MessageChain{}
	for _, msg := range msgs {
		if msg.Type == "text" {
			text := msg.Data["text"]
			res := splitAtAllToMsg(strings.Split(text, "@全体"))
			cMsg.SplitEmpty(res...)
		}
	}

	return message.Message(cMsg)
}

func splitAtAllToMsg(strs []string) []message.MessageSegment {
	var msgs []message.MessageSegment
	if len(strs) == 1 {
		msgs = append(msgs, message.Text(strs[0]))
		return msgs
	}
	if len(strs) > 1 && strs[0] == "" {
		msgs = append(msgs, message.AtAll())
		strs = strs[1:]
	}
	for idx, s := range strs {
		if idx == len(strs)-1 {
			if len(s) > 0 {
				msgs = append(msgs, message.Text(s))
			}

			break
		}
		if len(s) > 0 {
			msgs = append(msgs, message.Text(s))
			msgs = append(msgs, message.AtAll())
		}
	}

	return msgs
}

func parseAtInfo(input string) ([]string, []int64) {
	// TODO
	// 使用正则表达式匹配 @ 后的号码
	atRegex := regexp.MustCompile(`@(\d+)`)
	matches := atRegex.FindAllStringSubmatchIndex(input, -1)

	var texts []string
	var ats []int64
	var lastIndex int

	for _, match := range matches {
		// match[0] 是 "@数字" 的开始位置，match[1] 是结束位置
		// match[2] 是数字部分的开始位置，match[3] 是结束位置
		start, end, numStart, numEnd := match[0], match[1], match[2], match[3]

		// 提取 @ 前的文本
		if start > lastIndex {
			texts = append(texts, input[lastIndex:start])
		} else if start == 0 && lastIndex == 0 {
			texts = append(texts, "")
		}

		// 提取数字部分并转换为 int64
		atNumber, err := strconv.ParseInt(input[numStart:numEnd], 10, 64)
		if err != nil {
			return nil, nil
		}
		ats = append(ats, atNumber)

		lastIndex = end
	}

	// 提取剩余的文本
	if lastIndex < len(input) {
		texts = append(texts, input[lastIndex:])
	}

	return texts, ats
}
