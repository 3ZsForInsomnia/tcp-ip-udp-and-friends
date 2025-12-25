package logger

import (
	"context"
	"slices"
	"sync"
)

// MockLogger is a logger implementation for testing that captures log calls
type MockLogger struct {
	mu       sync.Mutex
	Warnings []string
	Errors   []string
	Infos    []string
}

// NewMockLogger creates a new mock logger for testing
func NewMockLogger() *MockLogger {
	return &MockLogger{
		Warnings: []string{},
		Errors:   []string{},
		Infos:    []string{},
	}
}

func PrepTest() *context.Context {
	ctx := NewMockLogger().WithLogger(context.Background())
	return &ctx
}

// WithLogger adds the mock logger to a context
func (m *MockLogger) WithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerContextKey, m)
}

// Warn captures warning messages
func (m *MockLogger) Warn(msg string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Warnings = append(m.Warnings, msg)
	return msg
}

// Error captures error messages
func (m *MockLogger) Error(msg string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Errors = append(m.Errors, msg)
	return msg
}

// Info captures info messages
func (m *MockLogger) Info(msg string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Infos = append(m.Infos, msg)
	return msg
}

// Reset clears all captured log messages
func (m *MockLogger) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Warnings = []string{}
	m.Errors = []string{}
	m.Infos = []string{}
}

// HasWarning checks if a specific warning message was logged
func (m *MockLogger) HasWarning(msg string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return slices.Contains(m.Warnings, msg)
}

// HasError checks if a specific error message was logged
func (m *MockLogger) HasError(msg string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return slices.Contains(m.Errors, msg)
}

// HasInfo checks if a specific info message was logged
func (m *MockLogger) HasInfo(msg string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return slices.Contains(m.Infos, msg)
}

// WarningCount returns the number of warnings logged
func (m *MockLogger) WarningCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Warnings)
}

// ErrorCount returns the number of errors logged
func (m *MockLogger) ErrorCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Errors)
}

// InfoCount returns the number of info messages logged
func (m *MockLogger) InfoCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Infos)
}
