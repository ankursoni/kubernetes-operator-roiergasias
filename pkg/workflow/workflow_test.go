package workflow_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/lib"
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/mocks"
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/workflow"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("test workflow for hello world", func() {
	Context("given single task and single step workflow", func() {
		var (
			logger         *zap.Logger
			text           string
			t              *mocks.MockITasks
			tw             *mocks.MockITaskWorkflow
			mockController *gomock.Controller
		)

		BeforeEach(func() {
			var err error
			logger, err = lib.NewZapLogger(true)
			if err != nil {
				log.Fatalln(fmt.Errorf("error creating new zap logger: %w", err))
				return
			}
			text = `
version: 0.1
environment:
  - welcome: "Welcome to the demo workflow!"
task:
  - sequential:
      - print:
          - "Hello World!"`
			mockController = gomock.NewController(GinkgoT())
			t = mocks.NewMockITasks(mockController)
			tw = mocks.NewMockITaskWorkflow(mockController)
		})

		It("create new workflow and run with task mock", func() {
			t.EXPECT().NewTask(gomock.Eq("sequential"), gomock.Any(), gomock.Eq("")).Return(tw)
			w, err := workflow.NewWorkflows(t, logger).NewWorkflowFromText(text)
			Expect(err).ToNot(HaveOccurred())
			Expect(w).ToNot(BeNil())

			tw.EXPECT().Run()
			err = w.Run()
			Expect(err).ToNot(HaveOccurred())
		})
		It("create new workflow and run", func() {
			w, err := workflow.NewWorkflows(nil, logger).NewWorkflowFromText(text)
			Expect(err).ToNot(HaveOccurred())
			Expect(w).ToNot(BeNil())

			rescueStdout := os.Stdout
			rf, wf, _ := os.Pipe()
			os.Stdout = wf

			err = w.Run()

			wf.Close()
			out, _ := ioutil.ReadAll(rf)
			os.Stdout = rescueStdout

			Expect(err).ToNot(HaveOccurred())
			Expect(string(out)).To(Equal("Hello World!\n"))
		})

		AfterEach(func() {
			mockController.Finish()
		})
	})

	Context("given 2 sequential tasks with node specifications", func() {
		var (
			logger *zap.Logger
			text   string
		)

		BeforeEach(func() {
			var err error
			logger, err = lib.NewZapLogger(true)
			if err != nil {
				log.Fatalln(fmt.Errorf("error creating new zap logger: %w", err))
				return
			}
			text = `
version: 0.1
environment:
  - welcome: "Welcome to the demo workflow!"
task:
  - node: node1
    sequential:
      - print:
          - "Hello World!"
      - environment:
          - greeting: "Warm greetings!"
  - node: node2
    sequential:
      - print:
          - "Hi Universe!"
      - print:
          - "{{env:welcome}}"
      - execute:
          - "echo {{env:greeting}}"`
		})

		It("create new workflow, call split nodes and run", func() {
			w, err := workflow.NewWorkflows(nil, logger).NewWorkflowFromText(text)
			Expect(err).ToNot(HaveOccurred())
			Expect(w).ToNot(BeNil())

			splitNodesList := w.SplitNodes()

			Expect(splitNodesList).ToNot(BeEmpty())
			Expect(len(splitNodesList)).To(Equal(2))

			Expect(splitNodesList[0].Node).To(Equal("node1"))
			Expect(splitNodesList[0].TaskList[0]["node"].(string)).To(Equal("node1"))
			Expect(len(splitNodesList[0].TaskList[0]["sequential"].([]interface{}))).To(Equal(2))
			Expect(len(splitNodesList[0].EnvironmentList)).To(Equal(1))

			Expect(splitNodesList[1].Node).To(Equal("node2"))
			Expect(splitNodesList[1].TaskList[0]["node"].(string)).To(Equal("node2"))
			Expect(len(splitNodesList[1].TaskList[0]["sequential"].([]interface{}))).To(Equal(3))
			Expect(len(splitNodesList[1].EnvironmentList)).To(Equal(2))

			rescueStdout := os.Stdout
			rf, wf, _ := os.Pipe()
			os.Stdout = wf

			err = w.Run()

			wf.Close()
			out, _ := ioutil.ReadAll(rf)
			os.Stdout = rescueStdout

			Expect(err).ToNot(HaveOccurred())
			Expect(string(out)).To(Equal("Hello World!\nHi Universe!\nWelcome to the demo workflow!\nWarm greetings!\n"))
		})
	})

	Context("given no task", func() {
		var (
			logger *zap.Logger
			text   string
		)

		BeforeEach(func() {
			var err error
			logger, err = lib.NewZapLogger(true)
			if err != nil {
				log.Fatalln(fmt.Errorf("error creating new zap logger: %w", err))
				return
			}
			text = `
version: 0.1
environment:
  - welcome: "Welcome to the demo workflow!"
task:`
		})

		It("create new workflow, call split nodes and run", func() {
			w, err := workflow.NewWorkflows(nil, logger).NewWorkflowFromText(text)
			Expect(err).ToNot(HaveOccurred())
			Expect(w).ToNot(BeNil())

			err = w.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("validation error: no task list found"))
		})
	})

	Context("given invalid task type", func() {
		var (
			logger *zap.Logger
			text   string
		)

		BeforeEach(func() {
			var err error
			logger, err = lib.NewZapLogger(true)
			if err != nil {
				log.Fatalln(fmt.Errorf("error creating new zap logger: %w", err))
				return
			}
			text = `
version: 0.1
environment:
  - welcome: "Welcome to the demo workflow!"
task:
  - sequential:
      - print:
          - "Hello World!"
  - invalid:
      - print:
          - "{{env:welcome}}"`
		})

		It("create new workflow, call split nodes and run", func() {
			w, err := workflow.NewWorkflows(nil, logger).NewWorkflowFromText(text)
			Expect(err).ToNot(HaveOccurred())
			Expect(w).ToNot(BeNil())

			err = w.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("validation error: invalid task type"))
		})
	})

	Context("given invalid step type", func() {
		var (
			logger *zap.Logger
			text   string
		)

		BeforeEach(func() {
			var err error
			logger, err = lib.NewZapLogger(true)
			if err != nil {
				log.Fatalln(fmt.Errorf("error creating new zap logger: %w", err))
				return
			}
			text = `
version: 0.1
environment:
  - welcome: "Welcome to the demo workflow!"
task:
  - sequential:
      - print:
          - "Hello World!"
  - sequential:
      - invalid:
          - "{{env:welcome}}"`
		})

		It("create new workflow, call split nodes and run", func() {
			w, err := workflow.NewWorkflows(nil, logger).NewWorkflowFromText(text)
			Expect(err).ToNot(HaveOccurred())
			Expect(w).ToNot(BeNil())

			err = w.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("validation error: invalid step type"))
		})
	})

	Context("given invalid step list in a valid task", func() {
		var (
			logger *zap.Logger
			text   string
		)

		BeforeEach(func() {
			var err error
			logger, err = lib.NewZapLogger(true)
			if err != nil {
				log.Fatalln(fmt.Errorf("error creating new zap logger: %w", err))
				return
			}
			text = `
version: 0.1
environment:
  - welcome: "Welcome to the demo workflow!"
task:
  - sequential:
      print:
        - "Hello World!"
  - sequential:
      - print:
          - "{{env:welcome}}"`
		})

		It("create new workflow, call split nodes and run", func() {
			w, err := workflow.NewWorkflows(nil, logger).NewWorkflowFromText(text)
			Expect(err).ToNot(HaveOccurred())
			Expect(w).ToNot(BeNil())

			err = w.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("validation error: invalid step list"))
		})
	})
})
