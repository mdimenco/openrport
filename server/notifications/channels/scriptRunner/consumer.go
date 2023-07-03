package scriptRunner

import (
	"context"
	"encoding/json"
	"time"

	"github.com/realvnc-labs/rport/server/notifications"
	"github.com/realvnc-labs/rport/share/logger"
)

const ScriptTimeout = time.Second * 20

type consumer struct {
	l *logger.Logger
}

//nolint:revive
func NewConsumer(l *logger.Logger) *consumer {
	return &consumer{
		l: l,
	}
}

func (c consumer) Process(details notifications.NotificationDetails) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), ScriptTimeout)
	defer cancelFunc()

	var content interface{} = map[string]interface{}{}
	var err error

	switch details.Data.ContentType {
	case notifications.ContentTypeTextJSON:
		err = json.Unmarshal([]byte(details.Data.Content), &content)
		if err != nil {
			return err
		}
	default:
		content = details.Data.Content
	}

	tmp := map[string]interface{}{
		"recipients": details.Data.Recipients,
		"data":       content,
	}

	data, err := json.Marshal(&tmp)
	if err != nil {
		return err
	}

	c.l.Debugf("running script: %s: with data: %s", details.Data.Target, string(data))

	err = RunCancelableScript(ctx, details.Data.Target, string(data))
	if err != nil {
		c.l.Debugf("failed running script: %s: with err: ", details.Data.Target, err)
		return err
	}

	return nil
}

func (c consumer) Target() notifications.Target {
	return notifications.TargetScript
}
