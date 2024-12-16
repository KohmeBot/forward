package forward

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"sync"
	"time"
)

type timerMap struct {
	mu  sync.Mutex
	m   map[int64]*time.Timer
	dur time.Duration
}

func newTimerMap(dur time.Duration) *timerMap {
	return &timerMap{m: make(map[int64]*time.Timer), dur: dur}
}

func (m *timerMap) Refresh(uid int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t, ok := m.m[uid]; ok {
		t.Reset(m.dur)
	}
}

func (m *timerMap) Has(uid int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.m[uid]
	return ok
}

func (m *timerMap) Start(uid int64, ctx *zero.Ctx) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.m[uid]
	if ok {
		return false
	}
	t := time.NewTimer(m.dur)
	m.m[uid] = t
	go func() {
		<-t.C
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.m, uid)

		ctx.Send(message.Text("你太久没说话啦，退出传话啦！"))

	}()
	return true
}
