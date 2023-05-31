package logger

import (
	"context"
	"fmt"

	"github.com/gookit/color"
)

// Infof logs at Info level and formats analogous to fmt.Printf.
func Infof(ctx context.Context, format string, args ...interface{}) {
	logMessage(ctx, sINFO, format, args...)
}

// Warningf logs at Info level and formats analogous to fmt.Printf.
func Warningf(ctx context.Context, format string, args ...interface{}) {
	logMessage(ctx, sWARNING, format, args...)
}

// Errorf logs at Info level and formats analogous to fmt.Printf.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logMessage(ctx, sERROR, format, args...)
}

// Criticalf logs at Info level and formats analogous to fmt.Printf.
func Criticalf(ctx context.Context, format string, args ...interface{}) {
	logMessage(ctx, sCRITICAL, format, args...)
}

func logMessage(ctx context.Context, severity string, format string, args ...interface{}) {
	//if we don't have trace ID, just log as plain text message
	if severity == sCRITICAL {
		color.Error.Println(fmt.Sprintf(format, args...))
	} else if severity == sERROR {
		color.Danger.Println(fmt.Sprintf(format, args...))
	} else if severity == sWARNING {
		color.Warn.Println(fmt.Sprintf(format, args...))
	} else {
		color.Info.Println(fmt.Sprintf(format, args...))
	}
}

// The severity of the event described in a log entry.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
const (
	sDEFAULT   = "DEFAULT"
	sDEBUG     = "DEBUG"
	sINFO      = "INFO"
	sNOTICE    = "NOTICE"
	sWARNING   = "WARNING"
	sERROR     = "ERROR"
	sCRITICAL  = "CRITICAL"
	sALERT     = "ALERT"
	sEMERGENCY = "EMERGENCY"
)
