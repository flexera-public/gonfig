package main

import (
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("gonfig example", func() {

	var build string      // Path to temporary compiled `config` binary
	var exampleDir string // Path to example directory

	BeforeEach(func() {
		var err error
		build, err = Build("github.com/rightscale/gonfig")
		Ω(err).ShouldNot(HaveOccurred())
		_, filename, _, ok := runtime.Caller(0)
		Ω(ok).Should(BeTrue())
		exampleDir = filepath.Join(path.Dir(filename), "example")
	})

	Describe("when generating Go code", func() {
		BeforeEach(func() {
			cmd := exec.Command(build, "-c", "config.json", "-o", "config.go")
			cmd.Dir = exampleDir
			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(session.Wait().Err.Contents()).Should(BeEmpty())
			Ω(session.Out.Contents()).Should(Equal([]byte("config.go\n")))
		})

		It("generates Go code that compiles", func() {
			_, err := Build("github.com/rightscale/gonfig/example")
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			//os.Remove(filepath.Join(exampleDir, "config.go"))
		})

	})

	AfterEach(func() {
		CleanupBuildArtifacts()
	})
})
