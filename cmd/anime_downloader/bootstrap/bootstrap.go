package bootstrap

import (
	"github.com/romycode/anime-downloader/pkg/errors"
	"github.com/romycode/anime-downloader/pkg/logger"
	"github.com/romycode/anime-downloader/pkg/storage"
	"github.com/romycode/anime-downloader/pkg/web"
)

func WarmUp(path string) (*errors.ErrorHandler, *web.URLExtractor, storage.Storage) {
	var il = logger.NewInfoLogger()
	var wl = logger.NewWarnLogger()
	var el = logger.NewErrorLogger()

	var l = logger.NewStdoutLogger(il, wl, el)
	var eh = errors.NewErrorHandler(l)

	var e = web.NewURLExtractor(eh)

	localStorage := storage.NewLocalStorage(path, eh)

	return eh, e, localStorage
}
