package logging

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"
)

type Logger struct {
    logger *log.Logger
    file   *os.File
}

func NewLogger(logDir string) (*Logger, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

    logFile := filepath.Join(logDir, fmt.Sprintf("backup_%s.log", time.Now().Format("2006-01-02")))
    file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }

    logger := log.New(file, "", log.LstdFlags)
    return &Logger{
        logger: logger,
        file:   file,
    }, nil
}

func (l *Logger) Info(format string, v ...interface{}) {
    l.logger.Printf("INFO: "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
    l.logger.Printf("ERROR: "+format, v...)
}

func (l *Logger) Close() error {
    return l.file.Close()
}