package util

import (
	"errors"

	"github.com/google/go-github/github"
)

func GetReleaseAsset(assets []github.ReleaseAsset, filename string) (github.ReleaseAsset, error) {
	for _, asset := range assets {
		if filename == *asset.Name {
			return asset, nil
		}
	}

	return github.ReleaseAsset{}, errors.New("util: release asset not found")
}
