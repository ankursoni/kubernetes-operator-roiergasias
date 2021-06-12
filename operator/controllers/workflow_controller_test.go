package controllers

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"time"

	workflowv1 "github.com/ankursoni/kubernetes-operator-roiergasias/operator/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Workflow controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		WorkflowName      = "test-workflow"
		WorkflowNamespace = "default"
		timeout           = time.Second * 10
		duration          = time.Second * 10
		interval          = time.Millisecond * 250
	)

	Context("When updating Workflow Status", func() {
		It("Should increase Workflow Status.Active count when new Jobs are created", func() {
			By("By creating a new Workflow")
			ctx := context.Background()
			workflow := &workflowv1.Workflow{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "batch.ankursoni.github.io/v1",
					Kind:       "Workflow",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      WorkflowName,
					Namespace: WorkflowNamespace,
				},
				Spec: workflowv1.WorkflowSpec{
					WorkflowYAML: workflowv1.WorkflowYAMLSpec{
						Name: "hello-world",
						YAML: `
version: 1.0
task:
  - sequential:
    - print:
      - "Hello World!"`,
					},
					JobTemplate: batchv1beta1.JobTemplateSpec{
						Spec: batchv1.JobSpec{
							Template: v1.PodTemplateSpec{
								Spec: v1.PodSpec{
									Containers: []v1.Container{
										{
											Name:  "roiergasias",
											Image: "roiergasias",
										},
									},
									RestartPolicy: v1.RestartPolicyNever,
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, workflow)).Should(Succeed())

			workflowLookupKey := types.NamespacedName{Name: WorkflowName, Namespace: WorkflowNamespace}
			createdWorkflow := &workflowv1.Workflow{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, workflowLookupKey, createdWorkflow)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Let's make sure our Schedule string value was properly converted/handled.
			Expect(createdWorkflow.Spec.WorkflowYAML.Name).Should(Equal("hello-world"))
		})
	})
})
