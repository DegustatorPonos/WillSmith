package geminiprotocol

import (
	globalstate "WillSmith/GlobalState"
	"time"
)

// Time after which the cache will be concidered outdated
var CacheTTL time.Duration

type PagesCache struct {
	CachedPages map[string]CachedPage
}

func InitCache() {
	CacheTTL = time.Minute * time.Duration(globalstate.CurrentSettings.CacheTTL)
}

// Single cached page
type CachedPage struct {
	Request Request
	FirstRequested time.Time
	InvalidationTime time.Time
	// In case we will extend TTL
	CustomTTL bool
	TTL time.Duration
}

// Adds request to cached pages
func (inp *PagesCache) AddPage(ToAdd Request) {
	if inp.CheckCacheValidity(ToAdd.URI) {
		return
	}

	inp.CachedPages[ToAdd.URI] = CachedPage{
		Request: ToAdd,
		FirstRequested: time.Now(),
		InvalidationTime: time.Now().Add(CacheTTL),
	}
}

func (inp *PagesCache) GetPageFromCache(query string) *Request {
	var val, exists = inp.CachedPages[query]
	if exists {
		return &val.Request
	}
	return nil
}

func (inp *PagesCache) InvalidatePage(query string) {
	delete(inp.CachedPages, query)
}

func (inp *PagesCache) CheckCacheValidity(query string) bool {
	var val, exists = inp.CachedPages[query]
	if !exists {
		return false
	}

	if val.InvalidationTime.After(time.Now()) {
		inp.InvalidatePage(query)
		return false
	}
	return true
}

func (inp *PagesCache) ClearOutdatedPages() {
	for k := range inp.CachedPages {
		if !inp.CheckCacheValidity(k) {
			inp.InvalidatePage(k)
		}
	}
}
