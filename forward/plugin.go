package forward

import (
	"github.com/kohmebot/pkg/chain"
	"github.com/kohmebot/plugin/v2"
	"github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"time"
)

type PluginForward struct {
	env  plugin.Env
	conf Config
	tmp  *timerMap
}

func NewPlugin() plugin.Plugin {
	return new(PluginForward)
}

func (p *PluginForward) OnInit(engine plugin.Engine, env plugin.Env) error {
	p.env = env
	err := env.GetConf(&p.conf)
	if err != nil {
		return err
	}
	p.tmp = newTimerMap(time.Duration(p.conf.StopDur) * time.Second)
	p.SetOnStart(engine)
	p.SetOnMsg(engine)
	return nil
}

func (p *PluginForward) OnHelp(ctx *zero.Ctx) {
	if !p.env.SuperUser().Rule()(ctx) {
		return
	}

	var msg chain.MessageChain

	msg.Split(
		message.Text("forward 插件所有命令"),
		message.Text("forward：开始传话"),
	)

	ctx.Send(msg)

}

func (p *PluginForward) Name() string {
	return "forward"
}

func (p *PluginForward) Version() string {
	return "v0.1.0"
}

func (p *PluginForward) OnBoot() {

}
