package registry_test

import (
	"os"

	"github.com/mikeweyandt/watchtower/internal/actions/mocks"
	unit "github.com/mikeweyandt/watchtower/pkg/registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"
)

var _ = Describe("Registry", func() {
	Describe("WarnOnAPIConsumption", func() {
		When("Given a container with an image from ghcr.io", func() {
			It("should want to warn", func() {
				Expect(testContainerWithImage("ghcr.io/containrrr/watchtower")).To(BeTrue())
			})
		})
		When("Given a container with an image implicitly from dockerhub", func() {
			It("should want to warn", func() {
				Expect(testContainerWithImage("docker:latest")).To(BeTrue())
			})
		})
		When("Given a container with an image explicitly from dockerhub", func() {
			It("should want to warn", func() {
				Expect(testContainerWithImage("index.docker.io/docker:latest")).To(BeTrue())
				Expect(testContainerWithImage("docker.io/docker:latest")).To(BeTrue())
			})
		})
		When("Given a container with an image from some other registry", func() {
			It("should not want to warn", func() {
				Expect(testContainerWithImage("docker.fsf.org/docker:latest")).To(BeFalse())
				Expect(testContainerWithImage("altavista.com/docker:latest")).To(BeFalse())
				Expect(testContainerWithImage("gitlab.com/docker:latest")).To(BeFalse())
			})
		})
	})

	Describe("GetPullAuth", func() {
		When("registry credentials are set in the environment", func() {
			It("should return the encoded credentials", func() {
				// The same base64 that trust_test.go asserts for these creds.
				const expected = "eyJ1c2VybmFtZSI6ImNvbnRhaW5ycnItdXNlciIsInBhc3N3b3JkIjoiY29udGFpbnJyci1wYXNzIn0="

				Expect(os.Setenv("REPO_USER", "containrrr-user")).To(Succeed())
				Expect(os.Setenv("REPO_PASS", "containrrr-pass")).To(Succeed())
				defer func() {
					_ = os.Unsetenv("REPO_USER")
					_ = os.Unsetenv("REPO_PASS")
				}()

				auth, err := unit.GetPullAuth("containrrr/watchtower")
				Expect(err).NotTo(HaveOccurred())
				Expect(auth).To(Equal(expected))
			})
		})

		When("no credentials are configured", func() {
			It("should return an empty string", func() {
				_ = os.Unsetenv("REPO_USER")
				_ = os.Unsetenv("REPO_PASS")

				// Point the docker config at an empty dir so no stored
				// credentials are found for the referenced registry.
				configDir, err := os.MkdirTemp("", "watchtower-docker-config")
				Expect(err).NotTo(HaveOccurred())
				defer func() { _ = os.RemoveAll(configDir) }()

				Expect(os.Setenv("DOCKER_CONFIG", configDir)).To(Succeed())
				defer func() { _ = os.Unsetenv("DOCKER_CONFIG") }()

				auth, err := unit.GetPullAuth("docker.io/library/nginx")
				Expect(err).NotTo(HaveOccurred())
				Expect(auth).To(Equal(""))
			})
		})
	})
})

func testContainerWithImage(imageName string) bool {
	container := mocks.CreateMockContainer("", "", imageName, time.Now())
	return unit.WarnOnAPIConsumption(container)
}
