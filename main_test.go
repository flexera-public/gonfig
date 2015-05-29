package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGonfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gonfig Suite")
}

var _ = Describe("gonfig", func() {

	var fatalError string

	BeforeEach(func() {
		fatalError = ""
		fatalf = func(format string, args ...interface{}) {
			fatalError = fmt.Sprintf(format, args...)
			panic("done") // :( is there a better way?
		}
	})

	AfterEach(func() {
		os.Remove(*outFile)
		os.Remove(*cfgPath)
	})

	Describe("with an invalid JSON config file path", func() {
		BeforeEach(func() {
			tf, err := ioutil.TempFile("", "")
			Ω(err).ShouldNot(HaveOccurred())
			*outFile = tf.Name()
			*cfgPath = "/some/invalid/path/__"
		})

		It("generates an error message indicating so", func() {
			defer func() { recover() }() // :( is there a better way?
			run()
			Ω(fatalError).Should(ContainSubstring("no configuration file at " + *cfgPath))
		})
	})

	Describe("with an existing JSON config file", func() {

		var cfgContent string

		JustBeforeEach(func() {
			tempFile, err := ioutil.TempFile("", "")
			Ω(err).ShouldNot(HaveOccurred())
			tempFile.WriteString(cfgContent)
			*cfgPath = tempFile.Name()
		})

		Describe("with a non existing output file path", func() {
			BeforeEach(func() {
				cfgContent = `{"foo":"bar"}`
				*outFile = filepath.Join(os.TempDir(), "/new/config.json")
			})

			It("creates the output file directory", func() {
				run()
				Ω(fatalError).Should(BeEmpty())
				_, err := os.Stat(*outFile)
				Ω(err).ShouldNot(HaveOccurred())
				content, err := ioutil.ReadFile(*outFile)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(content).ShouldNot(BeEmpty())
				Ω(content).Should(ContainSubstring("Foo"))
			})
		})

		Describe("with an existing output file path", func() {
			BeforeEach(func() {
				cfgContent = `{"foo":"bar"}`
				*outFile = filepath.Join(os.TempDir(), "config.json")
			})

			It("uses the existing file", func() {
				run()
				Ω(fatalError).Should(BeEmpty())
				_, err := os.Stat(*outFile)
				Ω(err).ShouldNot(HaveOccurred())
				content, err := ioutil.ReadFile(*outFile)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(content).ShouldNot(BeEmpty())
				Ω(content).Should(ContainSubstring("Foo"))
			})
		})

		Describe("with an invalid JSON content", func() {
			BeforeEach(func() {
				cfgContent = `nope{`
			})

			It("returns a helpful error message", func() {
				defer func() { recover() }() // :( is there a better way?
				run()
				Ω(fatalError).ShouldNot(BeEmpty())
				Ω(fatalError).Should(ContainSubstring("failed to unmarshal JSON"))
			})
		})

		Describe("with a non-object JSON content", func() {
			BeforeEach(func() {
				cfgContent = `42`
			})

			It("returns a helpful error message", func() {
				defer func() { recover() }() // :( is there a better way?
				run()
				Ω(fatalError).ShouldNot(BeEmpty())
				Ω(fatalError).Should(ContainSubstring("must define an object"))
			})
		})

	})

})
