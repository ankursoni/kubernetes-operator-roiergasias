/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ref "k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	workflowv1 "github.com/ankursoni/kubernetes-operator-roiergasias/operator/api/v1"
	wf "github.com/ankursoni/kubernetes-operator-roiergasias/pkg/workflow"
)

// WorkflowReconciler reconciles a Workflow object
type WorkflowReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=workflowv1.ankursoni.github.io,resources=workflows,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=workflowv1.ankursoni.github.io,resources=workflows/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=workflowv1.ankursoni.github.io,resources=workflows/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Workflow object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *WorkflowReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// get reference to workflow
	var workflow workflowv1.Workflow
	if err := r.Get(ctx, req.NamespacedName, &workflow); err != nil {
		logger.Error(err, "unable to fetch workflow")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// get reference to workflow jobs
	var workflowJobs batchv1.JobList
	if err := r.List(ctx, &workflowJobs, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		logger.Error(err, "unable to get workflow job")
		return ctrl.Result{}, err
	}

	// check if workflow jobs and update workflow status based on job status
	if workflowJobs.Items != nil && len(workflowJobs.Items) > 0 {
		workflowJob := workflowJobs.Items[0]
		jobRef, err := ref.GetReference(r.Scheme, &workflowJob)
		if err != nil {
			logger.Error(err, "unable to make reference to job", "job", workflowJob)
			return ctrl.Result{}, err
		}
		workflow.Status.Job = *jobRef

		if err := r.Status().Update(ctx, &workflow); err != nil {
			logger.Error(err, "unable to update workflow status")
		}
		return ctrl.Result{}, err
	}

	// validate spec workflow yaml
	if workflow.Spec.WorkflowYAML.Name == "" || workflow.Spec.WorkflowYAML.YAML == "" {
		err := fmt.Errorf("empty spec.workflowYAML name or yaml")
		logger.Error(err, "workflowYAM not proper in the spec")
		return ctrl.Result{}, err
	}

	// check if split execution is needed in spec workflow yaml
	wfYAML := wf.NewWorkflowFromText(workflow.Spec.WorkflowYAML.YAML)
	splitWfList := wfYAML.SplitNodes()

	// if the workflow is split then construct multiple configMap + job k8s resources
	if splitWfList != nil && len(splitWfList) > 0 {
		for k := range splitWfList {
			constructConfigMapForSplitWorkflow := func(workflow *workflowv1.Workflow, node string, yamlText string) (
				*corev1.ConfigMap, error) {
				name := fmt.Sprintf("%s-%s-%s", workflow.Name, workflow.Spec.WorkflowYAML.Name, node)
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Labels:      make(map[string]string),
						Annotations: make(map[string]string),
						Name:        name,
						Namespace:   workflow.Namespace,
					},
					Data: map[string]string{fmt.Sprintf("%s.yaml", workflow.Spec.WorkflowYAML.Name): yamlText},
				}
				if err := ctrl.SetControllerReference(workflow, configMap, r.Scheme); err != nil {
					return nil, err
				}
				return configMap, nil
			}
			// +kubebuilder:docs-gen:collapse=constructConfigMapForWorkflow

			splitWf := &splitWfList[k]
			yamlBytes, _ := yaml.Marshal(splitWf)
			configMap, err := constructConfigMapForSplitWorkflow(&workflow, splitWf.Node, string(yamlBytes))
			if err != nil {
				logger.Error(err, "unable to construct configMap from split node")
				return ctrl.Result{}, nil
			}

			existingConfigMap := &corev1.ConfigMap{}
			err = r.Get(ctx, client.ObjectKey{Name: configMap.Name, Namespace: configMap.Namespace}, existingConfigMap)
			if err != nil || existingConfigMap.Name == "" {
				if err := r.Create(ctx, configMap); err != nil {
					logger.Error(err, "unable to create configMap for workflow", "configMap", configMap)
					return ctrl.Result{}, err
				}
				logger.V(1).Info("created configMap for workflow run", "configMap", configMap)
			} else {
				if err := r.Update(ctx, configMap); err != nil {
					logger.Error(err, "unable to update configMap for workflow", "configMap", configMap)
					return ctrl.Result{}, err
				}
				logger.V(1).Info("updated configMap for workflow run", "configMap", configMap)
			}

			constructJobForWorkflow := func(workflow *workflowv1.Workflow, node string) (*batchv1.Job, error) {
				name := workflow.Name
				job := &batchv1.Job{
					ObjectMeta: metav1.ObjectMeta{
						Labels:      make(map[string]string),
						Annotations: make(map[string]string),
						Name:        name,
						Namespace:   workflow.Namespace,
					},
					Spec: *workflow.Spec.JobTemplate.Spec.DeepCopy(),
				}
				for k, v := range workflow.Spec.JobTemplate.Annotations {
					job.Annotations[k] = v
				}
				job.Annotations["node"] = node
				for k, v := range workflow.Spec.JobTemplate.Labels {
					job.Labels[k] = v
				}
				// find and remove any volumes with name "yaml"
				for i := 0; i < len(job.Spec.Template.Spec.Volumes); i++ {
					volume := job.Spec.Template.Spec.Volumes[i]
					if volume.Name == "yaml" {
						job.Spec.Template.Spec.Volumes = append(job.Spec.Template.Spec.Volumes[:i], job.Spec.Template.Spec.Volumes[i+1:]...)
						break
					}
				}
				// create volume with name "yaml"
				job.Spec.Template.Spec.Volumes = append(job.Spec.Template.Spec.Volumes,
					corev1.Volume{
						Name: "yaml",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: configMap.Name},
							},
						},
					},
				)
				if err := ctrl.SetControllerReference(workflow, job, r.Scheme); err != nil {
					return nil, err
				}
				return job, nil
			}
			// +kubebuilder:docs-gen:collapse=constructJobForWorkflow

			job, err := constructJobForWorkflow(&workflow, splitWf.Node)
			if err != nil {
				logger.Error(err, "unable to construct job from template")
				return ctrl.Result{}, nil
			}

			existingJob := &batchv1.Job{}
			err = r.Get(ctx, client.ObjectKey{Name: job.Name, Namespace: job.Namespace}, existingJob)
			if err != nil || existingJob.Name == "" {
				if err := r.Create(ctx, job); err != nil {
					logger.Error(err, "unable to create job for workflow", "job", job)
					return ctrl.Result{}, err
				}
				logger.V(1).Info("created job for workflow run", "job", job)
			} else {
				if err := r.Update(ctx, job); err != nil {
					logger.Error(err, "unable to update job for workflow", "job", job)
					return ctrl.Result{}, err
				}
				logger.V(1).Info("updated job for workflow run", "job", job)
			}
		}
	} else { // construct single configMap + job k8s resource
		constructConfigMapForWorkflow := func(workflow *workflowv1.Workflow) (*corev1.ConfigMap, error) {
			name := fmt.Sprintf("%s-%s", workflow.Name, workflow.Spec.WorkflowYAML.Name)
			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      make(map[string]string),
					Annotations: make(map[string]string),
					Name:        name,
					Namespace:   workflow.Namespace,
				},
				Data: map[string]string{fmt.Sprintf("%s.yaml", workflow.Spec.WorkflowYAML.Name): workflow.Spec.WorkflowYAML.YAML},
			}
			if err := ctrl.SetControllerReference(workflow, configMap, r.Scheme); err != nil {
				return nil, err
			}
			return configMap, nil
		}
		// +kubebuilder:docs-gen:collapse=constructConfigMapForWorkflow

		configMap, err := constructConfigMapForWorkflow(&workflow)
		if err != nil {
			logger.Error(err, "unable to construct configMap from template")
			return ctrl.Result{}, nil
		}

		existingConfigMap := &corev1.ConfigMap{}
		err = r.Get(ctx, client.ObjectKey{Name: configMap.Name, Namespace: configMap.Namespace}, existingConfigMap)
		if err != nil || existingConfigMap.Name == "" {
			if err := r.Create(ctx, configMap); err != nil {
				logger.Error(err, "unable to create configMap for workflow", "configMap", configMap)
				return ctrl.Result{}, err
			}
			logger.V(1).Info("created configMap for workflow run", "configMap", configMap)
		} else {
			if err := r.Update(ctx, configMap); err != nil {
				logger.Error(err, "unable to update configMap for workflow", "configMap", configMap)
				return ctrl.Result{}, err
			}
			logger.V(1).Info("updated configMap for workflow run", "configMap", configMap)
		}

		constructJobForWorkflow := func(workflow *workflowv1.Workflow) (*batchv1.Job, error) {
			name := workflow.Name
			job := &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      make(map[string]string),
					Annotations: make(map[string]string),
					Name:        name,
					Namespace:   workflow.Namespace,
				},
				Spec: *workflow.Spec.JobTemplate.Spec.DeepCopy(),
			}
			for k, v := range workflow.Spec.JobTemplate.Annotations {
				job.Annotations[k] = v
			}
			for k, v := range workflow.Spec.JobTemplate.Labels {
				job.Labels[k] = v
			}
			// find and remove any volumes with name "yaml"
			for i := 0; i < len(job.Spec.Template.Spec.Volumes); i++ {
				volume := job.Spec.Template.Spec.Volumes[i]
				if volume.Name == "yaml" {
					job.Spec.Template.Spec.Volumes = append(job.Spec.Template.Spec.Volumes[:i], job.Spec.Template.Spec.Volumes[i+1:]...)
					break
				}
			}
			// create volume with name "yaml"
			job.Spec.Template.Spec.Volumes = append(job.Spec.Template.Spec.Volumes,
				corev1.Volume{
					Name: "yaml",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: configMap.Name},
						},
					},
				},
			)
			if err := ctrl.SetControllerReference(workflow, job, r.Scheme); err != nil {
				return nil, err
			}
			return job, nil
		}
		// +kubebuilder:docs-gen:collapse=constructJobForWorkflow

		job, err := constructJobForWorkflow(&workflow)
		if err != nil {
			logger.Error(err, "unable to construct job from template")
			return ctrl.Result{}, nil
		}

		existingJob := &batchv1.Job{}
		err = r.Get(ctx, client.ObjectKey{Name: job.Name, Namespace: job.Namespace}, existingJob)
		if err != nil || existingJob.Name == "" {
			if err := r.Create(ctx, job); err != nil {
				logger.Error(err, "unable to create job for workflow", "job", job)
				return ctrl.Result{}, err
			}
			logger.V(1).Info("created job for workflow run", "job", job)
		} else {
			if err := r.Update(ctx, job); err != nil {
				logger.Error(err, "unable to update job for workflow", "job", job)
				return ctrl.Result{}, err
			}
			logger.V(1).Info("updated job for workflow run", "job", job)
		}
	}
	return ctrl.Result{}, nil
}

var (
	jobOwnerKey = ".metadata.controller"
	apiGVStr    = workflowv1.GroupVersion.String()
)

// SetupWithManager sets up the controller with the Manager.
func (r *WorkflowReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &batchv1.Job{}, jobOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner...
		job := rawObj.(*batchv1.Job)
		owner := metav1.GetControllerOf(job)
		if owner == nil {
			return nil
		}
		// ...make sure it's a Workflow...
		if owner.APIVersion != apiGVStr || owner.Kind != "Workflow" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&workflowv1.Workflow{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
