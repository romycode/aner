package pkg

type Downloader interface {
	GetEpisodes() ([]string, error)
	DownloadEpisodes(episodes []string)
}
