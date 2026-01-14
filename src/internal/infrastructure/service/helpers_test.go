package service

import (
	"suscord/internal/domain/logger"

	"github.com/stretchr/testify/mock"
)

type mockLogger struct {
	mock.Mock
}

func (l *mockLogger) Debug(msg string, fields ...logger.Field) { l.Called(msg, fields) }
func (l *mockLogger) Info(msg string, fields ...logger.Field)  { l.Called(msg, fields) }
func (l *mockLogger) Warn(msg string, fields ...logger.Field)  { l.Called(msg, fields) }
func (l *mockLogger) Error(msg string, fields ...logger.Field) { l.Called(msg, fields) }
func (l *mockLogger) Err(err error, fields ...logger.Field)    { l.Called(err, fields) }

