package main

import "time"

// Time after which the cashe will be concidered outdated
var CasheTTL time.Duration = time.Minute * 5

type PagesCashe struct {
	CashedPages map[string]CashedPage
}

// Single cashed page
type CashedPage struct {
	Request Request
	FirstRequested time.Time
	InvalidationTime time.Time
	// In case we will extend TTL
	CustomTTL bool
	TTL time.Duration
}

// Adds request to cashed pages
func (inp *PagesCashe) AddPage(ToAdd Request) {
	if inp.CheckCasheValidity(ToAdd.URI) {
		return
	}

	inp.CashedPages[ToAdd.URI] = CashedPage{
		Request: ToAdd,
		FirstRequested: time.Now(),
		InvalidationTime: time.Now().Add(CasheTTL),
	}
}

func (inp *PagesCashe) GetPageFromCashe(query string) *Request {
	var val, exists = inp.CashedPages[query]
	if exists {
		return &val.Request
	}
	return nil
}

func (inp *PagesCashe) InvalidatePage(query string) {
	delete(inp.CashedPages, query)
}

func (inp *PagesCashe) CheckCasheValidity(query string) bool {
	var val, exists = inp.CashedPages[query]
	if !exists {
		return false
	}

	if val.InvalidationTime.After(time.Now()) {
		inp.InvalidatePage(query)
		return false
	}
	return true
}

func (inp *PagesCashe) ClearOutdatedPages() {
	for k := range inp.CashedPages {
		if !inp.CheckCasheValidity(k) {
			inp.InvalidatePage(k)
		}
	}
}
