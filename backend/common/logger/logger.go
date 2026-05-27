package logger

import "go.uber.org/zap"

func New(serviceName string) (*zap.Logger, error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return log.With(zap.String("service", serviceName)), nil
}
