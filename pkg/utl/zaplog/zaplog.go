package zaplog

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log represents zerolog logger
type (
	Logger struct {
		logger  *zap.Logger
		version string
		// sugar   *zap.SugaredLogger
	}
)

var log *Logger

// New instantiates new zero logger
func Initialize(version string) (err error) {

	cfg := zap.NewProductionConfig()
	if os.Getenv("ENVIRONMENT") != "production" {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))
	if err != nil {
		return
	}

	log = &Logger{logger: logger, version: version}
	return nil
}

// Log logs using zerolog
func Log(msg string, params map[string]interface{}) (err error) {
	if log == nil {
		fmt.Println(errors.New("Logging not configured"))
		return
	}

	if params == nil {
		params = make(map[string]interface{})
	}

	delete(params, "time")
	var fields []zap.Field
	fields = append(fields, zap.String("version", log.version))
	for name, param := range params {
		fields = append(fields, zap.Any(name, param))
		if err2, ok := param.(error); ok && err2 != nil {
			err = err2
		}
	}

	if err != nil {
		params["error"] = err
		log.logger.Error(msg, fields...)
		return
	}

	log.logger.Info(msg, fields...)
	return nil
}

func ZLog(msg interface{}) (err error) {
	if log == nil {
		err = errors.New("Logging no configured")
		errors.New("Logging no configured")
		fmt.Println(err)
		panic(err)
	}
	// now := time.Now().Format("2006-01-02 15:04:05")
	switch msgType := msg.(type) {
	case nil:
		return
	case error:
		log.logger.Error(msgType.Error())
		return msgType
	case string:
		log.logger.Info(msgType)
	default:
		if JSON(msg) != "" {
			log.logger.Info(JSON(msg))
		}
	}
	return nil
}

func JSON(value interface{}) string {
	bytes, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		return ""
	}
	return string(bytes)
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}
