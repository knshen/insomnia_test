package model

import (
	"encoding/json"
)

type CacheEntry struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedAt int64  `json:"created_at"`
}

func (e *CacheEntry) Dump() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func NewCacheFromRawString(raw string) *CacheEntry {
	var ret CacheEntry
	json.Unmarshal([]byte(raw), &ret)
	return &ret
}
