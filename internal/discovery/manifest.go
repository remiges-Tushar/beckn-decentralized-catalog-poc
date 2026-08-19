package discovery

import (
	"encoding/json"
	"fmt"
)

const manifestPath = "/.well-known/dedi.json"

func (c *Crawler) FetchManifest(
	baseURL string,
) (Manifest, error) {

	url := baseURL + manifestPath

	body, _, err := c.Fetcher.Get(url)
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest

	if err := json.Unmarshal(body, &manifest); err != nil {
		return Manifest{}, fmt.Errorf(
			"parse manifest: %w",
			err,
		)
	}

	if err := ValidateManifest(manifest); err != nil {
		return Manifest{}, fmt.Errorf(
			"invalid manifest: %w",
			err,
		)
	}

	return manifest, nil
}

func ValidateManifest(
	manifest Manifest,
) error {

	if manifest.DeDiVersion == "" {
		return fmt.Errorf("missing dedi_version")
	}

	if manifest.Type == "" {
		return fmt.Errorf("missing type")
	}

	if manifest.Domain == "" {
		return fmt.Errorf("missing domain")
	}

	if len(manifest.Keys) == 0 {
		return fmt.Errorf("manifest contains no keys")
	}

	if len(manifest.Files) == 0 {
		return fmt.Errorf("manifest contains no files")
	}

	return nil
}

func (c *Crawler) FetchAndVerifyManifest(
	baseURL string,
) (Manifest, error) {

	manifest, err := c.FetchManifest(
		baseURL,
	)
	if err != nil {
		return Manifest{}, err
	}

	_, err = VerifyManifest(
		manifest,
	)
	if err != nil {
		return Manifest{}, fmt.Errorf(
			"manifest verification failed: %w",
			err,
		)
	}

	return manifest, nil
}
