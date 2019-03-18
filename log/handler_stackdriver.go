package log

import (
	"encoding/json"
	"fmt"

	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
)

// LabelMap defines stackdriver labels struct.
type LabelMap map[string]string

// LogMsg defines json object of log message.
type LogMsg struct {
	Msg   string   `json:"msg"`
	Label LabelMap `json:"label"`
}

// StackdriverHandler is a log handler of stackdrvier.
type StackdriverHandler struct {
	client *logging.Client
	logger *logging.Logger
}

// NewStackDriverHandler returns a new handler.
func NewStackDriverHandler(projectID, logName string) (*StackdriverHandler, error) {
	client, err := logging.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}
	return &StackdriverHandler{
		client:    client,
		logger:    client.Logger(logName),
		OmitEmpty: true,
	}, nil
}

// Log send to stackdriver.
func (h *StackdriverHandler) Log(r *Record) error {
	e := logging.Entry{
		Payload:   r.Msg,
		Timestamp: r.Time,
		Severity:  asSeverity(r.Lvl),
	}
	e.Labels = LabelMap{"tag": "test"}
	parseLabel(&e)
	h.logger.Log(e)
}

func parseLabel(e *logging.Entry) {
	var log LogMsg
	payload := []byte(e.Payload.(string))
	err := json.Unmarshal(payload, &log)
	if err != nil {
		fmt.Println("%v", err)
		return
	}

	e.Payload = log.Msg
	e.Labels = log.Label
}

func asSeverity(l Lvl) logging.Severity {
	switch l {
	case LvlDebug:
		return logging.Debug
	case LvlCrit:
		return logging.Emergency
	case LvlError:
		return logging.Error
	case LvlInfo:
		return logging.Info
	case LvlWarn:
		return logging.Warning
	default:
		return logging.Notice
	}
}

// Close is delegated to the client (if any)
func (h *StackdriverHandler) Close() error {
	if h.client == nil {
		return nil
	}
	return h.client.Close()
}
