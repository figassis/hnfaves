package util

import (
	"errors"
	"fmt"
	_ "net/http/pprof"
	"os"
	"time"

	uuid "github.com/gofrs/uuid"
	"github.com/lithammer/shortuuid"
)

const (
	Epoch = "1970-01-01 00:00:00"
)

var (
	timeFormats = []string{time.RFC3339, time.RFC1123Z, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006-01-02", "1/2/06", "1/2/06 15:05", "1_2_06"}
)

type (
	NullLogger struct{}
)

func ShortID() string {
	return shortuuid.New()
}

func ShortIDS(number int) (result []string) {
	for index := 0; index < number; index++ {
		result = append(result, shortuuid.New())
	}
	return
}

func GenerateUUID() (result string, err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return
	}
	result = newUuid.String()
	return
}

func GenerateUUIDS(number int) (result []string, err error) {
	for index := 0; index < number; index++ {
		newUuid, err2 := uuid.NewV4()
		if err = err2; err != nil {
			return
		}
		result = append(result, newUuid.String())
	}

	return
}

func CheckPath(path string) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("%s does not exist", path))
	}
	return nil
}

func (NullLogger) Print(...interface{}) {}
