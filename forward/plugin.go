package forward

import (
	"fmt"
	"github.com/kohmebot/pkg/command"
	"github.com/kohmebot/pkg/version"
	"github.com/kohmebot/plugin"
	"github.com/wdvxdr1123/ZeroBot"
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

func (p *PluginForward) Init(engine *zero.Engine, env plugin.Env) error {
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

func (p *PluginForward) Name() string {
	return "传话"
}

func (p *PluginForward) Description() string {
	return "转发私聊消息到对应群"
}

func (p *PluginForward) Commands() fmt.Stringer {
	return command.NewCommands()
}

func (p *PluginForward) Version() uint64 {
	return uint64(version.NewVersion(0, 0, 1))
}

func (p *PluginForward) OnBoot() {

}
