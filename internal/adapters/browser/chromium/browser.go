package chromium

import "github.com/vo1dFl0w/market-parser/internal/repository"

type Browser struct {
	browserRepo repository.BrowserRepository
}

func NewBrowser(browserRepo repository.BrowserRepository) *Browser {
	return &Browser{browserRepo: browserRepo}
}

func (b *Browser) Chromium() repository.BrowserRepository {
	return b.browserRepo
}
