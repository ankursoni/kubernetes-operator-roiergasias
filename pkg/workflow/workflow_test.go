package workflow_test

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/mocks"
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/workflow"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("test workflow for hello world", func() {
	Context("given single task and single step workflow from text", func() {
		var (
			text           string
			t              *mocks.MockITasks
			tw             *mocks.MockITaskWorkflow
			mockController *gomock.Controller
		)

		BeforeEach(func() {
			text = `
version: 1.0
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
			w, err := workflow.NewWorkflows(t).NewWorkflowFromText(text)
			Expect(err).ToNot(HaveOccurred())
			Expect(w).ToNot(BeNil())

			tw.EXPECT().Run()
			err = w.Run()
			Expect(err).ToNot(HaveOccurred())
		})
		It("create new workflow and run", func() {
			w, err := workflow.NewWorkflows(nil).NewWorkflowFromText(text)
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
			text string
		)

		BeforeEach(func() {
			text = `
version: 1.0
environment:
  - welcome: "Welcome to the demo workflow!"
task:
  - node: node1
    sequential:
      - print:
          - "Hello World!"
        set-environment:
          - greeting: "Warm greetings!"
  - node: node2
    sequential:
      - print:
          - "{{env:welcome}}"
      - execute:
          - "echo {{env:greeting}}"`
		})

		It("create new workflow, call split nodes and run", func() {
			w, err := workflow.NewWorkflows(nil).NewWorkflowFromText(text)
			Expect(err).ToNot(HaveOccurred())
			Expect(w).ToNot(BeNil())

			splitNodesList := w.SplitNodes()

			Expect(splitNodesList).ToNot(BeEmpty())
			Expect(len(splitNodesList)).To(Equal(2))

			Expect(splitNodesList[0].Node).To(Equal("node1"))
			Expect(splitNodesList[0].TaskList[0]["node"].(string)).To(Equal("node1"))
			Expect(len(splitNodesList[0].TaskList[0]["sequential"].([]interface{}))).To(Equal(1))
			Expect(len(splitNodesList[0].EnvironmentList)).To(Equal(1))

			Expect(splitNodesList[1].Node).To(Equal("node2"))
			Expect(splitNodesList[1].TaskList[0]["node"].(string)).To(Equal("node2"))
			Expect(len(splitNodesList[1].TaskList[0]["sequential"].([]interface{}))).To(Equal(2))
			Expect(len(splitNodesList[1].EnvironmentList)).To(Equal(2))

			rescueStdout := os.Stdout
			rf, wf, _ := os.Pipe()
			os.Stdout = wf

			err = w.Run()

			wf.Close()
			out, _ := ioutil.ReadAll(rf)
			os.Stdout = rescueStdout

			Expect(err).ToNot(HaveOccurred())
			Expect(string(out)).To(Equal("Hello World!\nWelcome to the demo workflow!\nWarm greetings!\n"))
		})
	})
})
