package delivery

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/xmtp/example-notification-server-go/pkg/options"
	"go.uber.org/zap"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type EmailDelivery struct {
	ctx    context.Context
	opts   options.EmailOptions
	logger *zap.Logger
}

func NewEmailDelivery(ctx context.Context, logger *zap.Logger, opts options.EmailOptions) (*EmailDelivery, error) {
	return &EmailDelivery{
		ctx:    ctx,
		opts:   opts,
		logger: logger,
	}, nil
}

func (e *EmailDelivery) Send(ctx context.Context, token, topic, message string) error {

	var body = map[string]interface{}{
		"token":     token,
		"topic":     topic,
		"timestamp": time.Now().Unix(),
		"message":   message,
	}

	bz, _ := json.Marshal(body)
	request, err := http.NewRequest("POST", os.Getenv("MSG_HANDLER_ADDRESS")+"/v1/notifications/new-message", bytes.NewBuffer(bz))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Basic "+basicAuth(os.Getenv("MSG_HANDLER_USER"), os.Getenv("MSG_HANDLER_PASS")))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		e.logger.Error("error sending email", zap.Error(err))
		return err
	}
	defer response.Body.Close()

	return nil
}
