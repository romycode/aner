package errors

import (
	"log"

	"github.com/romycode/anime-downloader/pkg/logger"
)

type ErrorHandler struct {
	logger *log.Logger
}

func NewErrorHandler(logger logger.Logger) *ErrorHandler {
	return &ErrorHandler{}
}

func (eh ErrorHandler) HandleError(err error) {
		eh.logger.Fatalln(err)
}
