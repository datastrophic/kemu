package test

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"testing"

	"datastrophic.io/kemu/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"sigs.k8s.io/e2e-framework/support/kind"
)

func TestKEMU(t *testing.T) {
	format.MaxLength = 0
	RegisterFailHandler(Fail)
	RunSpecs(t, "running KEMU test suite")
}

var _ = BeforeSuite(func() {
	logTextHandler := slog.NewTextHandler(GinkgoWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(logTextHandler)
	slog.SetDefault(logger)

	cmd := exec.Command("go", "build", "-cover", "-o", ".run/kemu", "main.go")
	_, err := utils.Run(cmd)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "failed to build the kemu binary")
})

var _ = AfterSuite(func() {
	knownClusters := []string{"e2e-cli-test-cluster", "it-simple", "it-with-kind-config"}
	for _, cluster := range knownClusters {
		err := kind.NewProvider().SetDefaults().WithName(cluster).Destroy(context.Background())
		Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("failed to destroy cluster %s", cluster))
	}
})
