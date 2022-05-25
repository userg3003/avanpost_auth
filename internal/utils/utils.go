package utils

import (
	"github.com/rs/zerolog"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type CallerHook struct{}

func (h CallerHook) Run(event *zerolog.Event, _ zerolog.Level, _ string) {
	if pc, _, _, ok := runtime.Caller(3); ok {
		details := runtime.FuncForPC(pc)
		name := "???"
		if ok && details != nil {
			name = details.Name()
		}
		event.Str("fn", name[strings.LastIndex(name, "/")+1:])
	}
}
