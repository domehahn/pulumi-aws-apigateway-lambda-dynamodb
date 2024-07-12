package error

import (
	"github.com/pkg/errors"
	"log"
)

func HandlingError(message string, details ...interface{}) {
	log.Fatalln(errors.Errorf(message, details...))
}
