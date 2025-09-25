package test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"datastrophic.io/kemu/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"sigs.k8s.io/e2e-framework/support/kind"
)

// rootProjectDir is configured at the suite setup time and is used
// by all tests to avoid relative path inconsistencies.
var rootProjectDir string

func TestKEMU(t *testing.T) {
	format.MaxLength = 0
	RegisterFailHandler(Fail)
	RunSpecs(t, "running KEMU test suite")
}

// Using SynchronizedBeforeSuite to set rootProjectDir.
var _ = SynchronizedBeforeSuite(func() []byte {
	//runs *only* on process #1
	rootDir, err := utils.GetProjectDir()

	Expect(err).NotTo(HaveOccurred())
	Expect(rootDir).NotTo(BeEmpty())

	return []byte(rootDir)
}, func(data []byte) {
	//runs on *all* processes
	rootProjectDir = string(data)

	logTextHandler := slog.NewTextHandler(GinkgoWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(logTextHandler)
	slog.SetDefault(logger)
})

var _ = SynchronizedAfterSuite(func() {
	//runs on *all* processes, noop
}, func() {
	//runs *only* on process #1
	knownClusters := []string{"it-simple", "it-already-exists", "it-with-kind-config", "it-with-addons", "it-with-kwok-nodes", "it-with-full-config", "e2e-cli-test-cluster"}
	for _, cluster := range knownClusters {
		err := kind.NewProvider().SetDefaults().WithName(cluster).Destroy(context.Background())
		Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("failed to destroy cluster %s", cluster))
	}
})
