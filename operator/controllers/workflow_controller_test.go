package controllers

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"time"

	workflowv1 "github.com/ankursoni/kubernetes-operator-roiergasias/operator/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("test workflow controller", func() {
	const (
		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)
	var (
		workflowName      string
		workflowNamespace string
		workflowYAMLName  string
		workflowYAMLText  string
	)

	Context("when running a single node, single task workflow", func() {
		BeforeEach(func() {
			workflowName = "test-workflow"
			workflowNamespace = "default"
			workflowYAMLName = "test-yaml"
			workflowYAMLText = `
version: 0.1
task:
  - sequential:
    - print:
      - "Hello World!"`
		})
		It("should increase workflow Status.ActiveJobs count when new workflow is created", func() {
			By("by creating a new Workflow")
			ctx := context.Background()
			testWorkflow := &workflowv1.Workflow{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "batch.ankursoni.github.io/v1",
					Kind:       "Workflow",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      workflowName,
					Namespace: workflowNamespace,
				},
				Spec: workflowv1.WorkflowSpec{
					WorkflowYAML: workflowv1.WorkflowYAMLSpec{
						Name: workflowYAMLName,
						YAML: workflowYAMLText,
					},
					JobTemplate: batchv1beta1.JobTemplateSpec{
						Spec: batchv1.JobSpec{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "test-container",
											Image: "test-image",
										},
									},
									RestartPolicy: corev1.RestartPolicyNever,
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, testWorkflow)).Should(Succeed())

			workflowLookupKey := types.NamespacedName{Name: workflowName, Namespace: workflowNamespace}
			createdWorkflow := &workflowv1.Workflow{}
			By("by checking the test workflow is created")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, workflowLookupKey, createdWorkflow)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdWorkflow.Spec.WorkflowYAML).ShouldNot(BeNil())
			Expect(createdWorkflow.Spec.WorkflowYAML.Name).Should(Equal(workflowYAMLName))

			By("by checking the created workflow has one active job count")
			Eventually(func() (int, error) {
				err := k8sClient.Get(ctx, workflowLookupKey, createdWorkflow)
				if err != nil {
					return -1, err
				}
				return len(createdWorkflow.Status.ActiveJobs), nil
			}, duration, interval).Should(Equal(1))

			configMapLookupKey := types.NamespacedName{Name: fmt.Sprintf("%s-%s", workflowName, workflowYAMLName),
				Namespace: workflowNamespace}
			configMapDataKey := fmt.Sprintf("%s.yaml", workflowYAMLName)
			createdConfigMap := &corev1.ConfigMap{}
			By("by checking the test workflow has one configmap created")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, configMapLookupKey, createdConfigMap)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(len(createdConfigMap.Data)).Should(Equal(1))
			Expect(string(createdConfigMap.Data[configMapDataKey])).Should(ContainSubstring("Hello World!"))
		})
	})

	Context("when running a multi node, multi task workflow", func() {
		BeforeEach(func() {
			workflowName = "test-multi-workflow"
			workflowNamespace = "default"
			workflowYAMLName = "test-multi-yaml"
			workflowYAMLText = `
version: 0.1
task:
  - node: "node1"
    sequential:
    - print:
      - "Hello World!"
  - node: "node2"
    sequential:
    - print:
      - "Hi Universe!"`
		})
		It("should increase workflow Status.ActiveJobs count when new workflow is created", func() {
			By("by creating a new Workflow")
			ctx := context.Background()
			testWorkflow := &workflowv1.Workflow{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "batch.ankursoni.github.io/v1",
					Kind:       "Workflow",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      workflowName,
					Namespace: workflowNamespace,
				},
				Spec: workflowv1.WorkflowSpec{
					WorkflowYAML: workflowv1.WorkflowYAMLSpec{
						Name: workflowYAMLName,
						YAML: workflowYAMLText,
					},
					JobTemplate: batchv1beta1.JobTemplateSpec{
						Spec: batchv1.JobSpec{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "test-container",
											Image: "test-image",
										},
									},
									RestartPolicy: corev1.RestartPolicyNever,
									Volumes: []corev1.Volume{
										{
											Name: "yaml",
											VolumeSource: corev1.VolumeSource{
												EmptyDir: &corev1.EmptyDirVolumeSource{},
											},
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, testWorkflow)).Should(Succeed())

			workflowLookupKey := types.NamespacedName{Name: workflowName, Namespace: workflowNamespace}
			createdWorkflow := &workflowv1.Workflow{}
			By("by checking the test workflow is created")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, workflowLookupKey, createdWorkflow)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdWorkflow.Spec.WorkflowYAML).ShouldNot(BeNil())
			Expect(createdWorkflow.Spec.WorkflowYAML.Name).Should(Equal(workflowYAMLName))

			By("by checking the created workflow has one active job count")
			Eventually(func() (int, error) {
				err := k8sClient.Get(ctx, workflowLookupKey, createdWorkflow)
				if err != nil {
					return -1, err
				}
				return len(createdWorkflow.Status.ActiveJobs), nil
			}, duration, interval).Should(Equal(1))

			configMapLookupKey := types.NamespacedName{Name: fmt.Sprintf("%s-%s-%d-%s", workflowName, workflowYAMLName, 1, "node1"),
				Namespace: workflowNamespace}
			configMapDataKey := fmt.Sprintf("%s.yaml", workflowYAMLName)
			createdConfigMap := &corev1.ConfigMap{}
			By("by checking the test workflow has one configmap created for first split")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, configMapLookupKey, createdConfigMap)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(len(createdConfigMap.Data)).Should(Equal(1))
			Expect(string(createdConfigMap.Data[configMapDataKey])).Should(ContainSubstring("Hello World!"))
			Expect(string(createdConfigMap.Data[configMapDataKey])).ShouldNot(ContainSubstring("Hi Universe!"))
		})
	})
})
