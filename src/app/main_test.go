package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"

	"app/utils/logger"
)

func init() {
	if os.Getenv("LOG_LEVEL") == "" {
		godotenv.Load(".env")
	}
	logLv := logger.Error
	envLogLv := os.Getenv("LOG_LEVEL")
	if envLogLv != "" {
		n, _ := strconv.Atoi(envLogLv)
		logLv = logger.LogLv(n)
	}
	logger.LogInitialize(logLv, 25)
}

func TestHandlerSuccess(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{
			name:       "normal",
			args:       args{key: "type", value: "main_summary:221309_hamamatsu_covid19_patients,main_summary:221309_hamamatsu_covid19_patients_summary,patients:221309_hamamatsu_covid19_patients,patients_summary:221309_hamamatsu_covid19_patients,inspection_persons:221309_hamamatsu_covid19_test_people,contacts:221309_hamamatsu_covid19_call_center"},
			statusCode: 200,
		},
	}
	for _, tt := range tests {
		req := events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{tt.args.key: tt.args.value},
		}
		t.Run(tt.name, func(t *testing.T) {
			ret, err := handler(req)
			if err != nil {
				t.Errorf(tt.name)
			} else if ret.StatusCode != tt.statusCode {
				t.Errorf(tt.name)
			}
		})
	}
}
