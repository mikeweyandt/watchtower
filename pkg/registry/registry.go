package registry

import (
	"context"

	ref "github.com/distribution/reference"
	"github.com/mikeweyandt/watchtower/pkg/registry/helpers"
	watchtowerTypes "github.com/mikeweyandt/watchtower/pkg/types"
	log "github.com/sirupsen/logrus"
)

// GetPullAuth returns the encoded registry credentials to use when pulling the
// given image, or an empty string when no credentials are configured.
//
// The caller turns this into the SDK's pull options, which keeps the Docker
// client dependency confined to the pkg/container package.
func GetPullAuth(imageName string) (string, error) {
	auth, err := EncodedAuth(imageName)
	log.Debugf("Got image name: %s", imageName)
	if err != nil {
		return "", err
	}

	// CREDENTIAL: Uncomment to log docker config auth
	// log.Tracef("Got auth value: %s", auth)

	return auth, nil
}

// DefaultAuthHandler will be invoked if an AuthConfig is rejected
// It could be used to return a new value for the "X-Registry-Auth" authentication header,
// but there's no point trying again with the same value as used in AuthConfig
func DefaultAuthHandler(context.Context) (string, error) {
	log.Debug("Authentication request was rejected. Trying again without authentication")
	return "", nil
}

// WarnOnAPIConsumption will return true if the registry is known-expected
// to respond well to HTTP HEAD in checking the container digest -- or if there
// are problems parsing the container hostname.
// Will return false if behavior for container is unknown.
func WarnOnAPIConsumption(container watchtowerTypes.Container) bool {

	normalizedRef, err := ref.ParseNormalizedNamed(container.ImageName())
	if err != nil {
		return true
	}

	containerHost, err := helpers.GetRegistryAddress(normalizedRef.Name())
	if err != nil {
		return true
	}

	if containerHost == helpers.DefaultRegistryHost || containerHost == "ghcr.io" {
		return true
	}

	return false
}
