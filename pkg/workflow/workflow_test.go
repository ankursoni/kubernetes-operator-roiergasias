package workflow_test

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/mocks"
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/workflow"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
task:
  - sequential:
    - print:
      - "Hello"`
			mockController = gomock.NewController(GinkgoT())
			t = mocks.NewMockITasks(mockController)
			tw = mocks.NewMockITaskWorkflow(mockController)
		})

		It("create new workflow with no error", func() {
			w, err := workflow.NewWorkflows(nil).NewWorkflowFromText(text)
			Expect(err).ToNot(HaveOccurred())
			Expect(w).ToNot(BeNil())
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

		AfterEach(func() {
			mockController.Finish()
		})
	})
})
