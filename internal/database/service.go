package database

import (
	"fmt"
	"github.com/jackc/pgx/v5/tracelog"
	"os"
)

func LogLevelFromEnv() (tracelog.LogLevel, error) {
	if level := os.Getenv("PGX_LOG_LEVEL"); level != "" {
		l, err := tracelog.LogLevelFromString(level)
		if err != nil {
			return tracelog.LogLevelDebug, fmt.Errorf("pgx configuration: %w", err)
		}
		return l, nil
	}
	return tracelog.LogLevelInfo, nil
}
