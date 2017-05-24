package emitter

import (
	"encoding/json"
	"fmt"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/monkeyherder/salus/checks"
	"sync"
)

const TAG string = "emitter"

type Target string

const (
	HealthMonitor = Target("hm")
	Director      = Target("director")
)

type Topic string

const (
	Check  = Topic("check")
	Metric = Topic("metric")
)

type Client interface {
	Publish(subject string, bytes []byte) error
}

type Emitter struct {
	Logger           boshlog.Logger
	Client           Client
	checkInfoChannel chan checks.CheckInfo
	done             chan struct{}
	startOnce        sync.Once
	shutdownOnce     sync.Once
}

func (emitter *Emitter) Start() {
	emitter.startOnce.Do(func() {
		emitter.checkInfoChannel = make(chan checks.CheckInfo)
		emitter.done = make(chan struct{})

		go func(checkChannel chan checks.CheckInfo, done chan struct{}) {
			for {
				select {
				case c := <-checkChannel:
					emitter.send(HealthMonitor, Check, c)
				case <-done:
					close(checkChannel)
					return
				}
			}
		}(emitter.checkInfoChannel, emitter.done)
	})

}

func (emitter *Emitter) Shutdown() {
	emitter.shutdownOnce.Do(func() {
		close(emitter.done)
	})
}

func (e *Emitter) EmitCheck() chan<- checks.CheckInfo {
	return e.checkInfoChannel
}

func (e *Emitter) send(target Target, topic Topic, message interface{}) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return bosherr.WrapErrorf(err, "Marshalling message (target=%s, topic=%s): %#v", target, topic, message)
	}

	e.Logger.Info(TAG, "Sending %s message '%s'", target, topic)
	e.Logger.DebugWithDetails(TAG, "Message Payload", string(bytes))

	//TODO: how should the uuid be determined?
	uuid := "uuid"
	subject := fmt.Sprintf("%s.salus.%s.%s", target, topic, uuid)
	return e.Client.Publish(subject, bytes)
}
