package catalog

import (
	"fmt"
	"os"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func BuildArtifactRef(
	version int64,
	url string,
	path string,
) (ArtifactRef, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return ArtifactRef{}, fmt.Errorf(
			"read artifact: %w",
			err,
		)
	}

	info, err := os.Stat(path)
	if err != nil {
		return ArtifactRef{}, fmt.Errorf(
			"stat artifact: %w",
			err,
		)
	}

	return ArtifactRef{
		Version: version,
		URL:     url,
		Size:    info.Size(),
		Digest:  catalogcrypto.Digest(data),
	}, nil
}
