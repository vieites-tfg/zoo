package main

import (
	"dagger/dagger/internal/dagger"
	"slices"
)

type SecMap struct {
	Keys   []string
	Values []*dagger.Secret
}

func (m SecMap) Get(key string) *dagger.Secret {
	return m.Values[slices.Index(m.Keys, key)]
}
