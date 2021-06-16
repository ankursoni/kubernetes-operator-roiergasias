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
	workflowv1 "github.com/ankursoni/kubernetes-operator-roiergasias/operator/api/v1"
	wf "github.com/ankursoni/kubernetes-operator-roiergasias/pkg/workflow"
	"gopkg.in/yaml.v3"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ref "k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// WorkflowReconciler reconciles a Workflow object
type WorkflowReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=workflowv1.ankursoni.github.io,resources=workflows,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=workflowv1.ankursoni.github.io,resources=workflows/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=workflowv1.ankursoni.github.io,resources=workflows/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configMaps,verbs=get;list;watch;create;update;patch;delete
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
	workflow := workflowv1.Workflow{}
	if err := r.Get(ctx, req.NamespacedName, &workflow); err != nil {
		logger.Error(err, "unable to fetch workflow")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// get reference to child jobs
	childJobs := batchv1.JobList{}
	if err := r.List(ctx, &childJobs, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		logger.Error(err, "unable to list child jobs")
		return ctrl.Result{}, err
	}

	// find the active list of jobs
	activeJobs := []batchv1.Job{}
	successfulJobs := []batchv1.Job{}
	failedJobs := []batchv1.Job{}
	// check all child jobs and update workflow status based on job status
	for _, job := range childJobs.Items {
		_, finishedType := r.isJobFinished(&job)
		switch finishedType {
		case "":
			activeJobs = append(activeJobs, job)
		case batchv1.JobComplete:
			successfulJobs = append(successfulJobs, job)
		case batchv1.JobFailed:
			failedJobs = append(failedJobs, job)
		}
	}
	activeJobRefs := []corev1.ObjectReference{}
	for _, activeJob := range activeJobs {
		jobRef, err := ref.GetReference(r.Scheme, &activeJob)
		if err != nil {
			logger.Error(err, "unable to make reference to active job", "job", activeJob)
			continue
		}
		activeJobRefs = append(activeJobRefs, *jobRef)
	}
	successfulJobRefs := []corev1.ObjectReference{}
	for _, successfulJob := range successfulJobs {
		jobRef, err := ref.GetReference(r.Scheme, &successfulJob)
		if err != nil {
			logger.Error(err, "unable to make reference to successful job", "job", successfulJob)
			continue
		}
		successfulJobRefs = append(successfulJobRefs, *jobRef)
	}
	failedJobRefs := []corev1.ObjectReference{}
	for _, failedJob := range failedJobs {
		jobRef, err := ref.GetReference(r.Scheme, &failedJob)
		if err != nil {
			logger.Error(err, "unable to make reference to failed job", "job", failedJob)
			continue
		}
		failedJobRefs = append(failedJobRefs, *jobRef)
	}
	logger.V(1).Info("job count", "active jobs", len(activeJobs),
		"successful jobs", len(successfulJobs), "failed jobs", len(failedJobs), "total jobs", len(childJobs.Items))

	if err := r.Get(ctx, req.NamespacedName, &workflow); err != nil {
		logger.Error(err, "unable to fetch workflow")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	workflow.Status.ActiveJobs = activeJobRefs
	workflow.Status.SuccessfulJobs = successfulJobRefs
	workflow.Status.FailedJobs = failedJobRefs
	if err := r.Status().Update(ctx, &workflow); err != nil {
		logger.Error(err, "unable to update workflow status")
		return ctrl.Result{}, err
	}

	// validate spec workflow yaml
	if workflow.Spec.WorkflowYAML.Name == "" || workflow.Spec.WorkflowYAML.YAML == "" {
		err := fmt.Errorf("empty spec.workflowYAML name or yaml")
		logger.Error(err, "workflowYAML not proper in the spec")
		return ctrl.Result{}, err
	}

	// check if split execution is needed in spec workflow yaml
	wfYAML, err := wf.NewWorkflows(nil).NewWorkflowFromText(workflow.Spec.WorkflowYAML.YAML)
	if err != nil {
		err = fmt.Errorf("invalid spec.workflowYAML yaml with error %w", err)
		logger.Error(err, "workflowYAML not proper in the spec")
		return ctrl.Result{}, err
	}
	splitWfList := wfYAML.SplitNodes()
	if len(splitWfList) == 0 { // construct single configMap + job k8s resource
		if len(activeJobs)+len(successfulJobs)+len(failedJobs) > 0 {
			return ctrl.Result{}, nil
		}

		configMap, err := r.constructConfigMapForWorkflow(&workflow, workflow.Spec.WorkflowYAML.YAML, "")
		if err != nil {
			logger.Error(err, "unable to construct configMap from template")
			return ctrl.Result{}, nil
		}
		if err := r.Create(ctx, configMap); err != nil {
			logger.Error(err, "unable to create configMap for workflow", "configMap", configMap)
			return ctrl.Result{}, err
		}
		logger.V(1).Info("created configMap for workflow", "configMap", configMap)

		job, err := r.constructJobForWorkflow(&workflow, configMap, "")
		if err != nil {
			logger.Error(err, "unable to construct job from template")
			return ctrl.Result{}, nil
		}
		if err := r.Create(ctx, job); err != nil {
			logger.Error(err, "unable to create job for workflow", "job", job)
			return ctrl.Result{}, err
		}
		logger.V(1).Info("created job for workflow run", "job", job)
	} else { // else if the workflow is split then construct multiple configMaps + jobs k8s resources
		for k, splitWf := range splitWfList {
			if len(activeJobs)+len(successfulJobs)+len(failedJobs) > 0 {
				if len(activeJobs) > 0 || k > len(successfulJobs) {
					return ctrl.Result{}, nil
				}
				if k < len(successfulJobs) {
					continue
				}
			}
			yamlBytes, _ := yaml.Marshal(splitWf)
			configMap, err := r.constructConfigMapForWorkflow(&workflow, string(yamlBytes), splitWf.Node)
			if err != nil {
				logger.Error(err, "unable to construct configMap from split node template")
				return ctrl.Result{}, nil
			}
			if err := r.Create(ctx, configMap); err != nil {
				logger.Error(err, "unable to create configMap for workflow", "configMap", configMap)
				return ctrl.Result{}, err
			}
			logger.V(1).Info("created configMap for workflow", "configMap", configMap)

			job, err := r.constructJobForWorkflow(&workflow, configMap, splitWf.Node)
			if err != nil {
				logger.Error(err, "unable to construct job from split node template")
				return ctrl.Result{}, nil
			}
			if err := r.Create(ctx, job); err != nil {
				logger.Error(err, "unable to create job for workflow", "job", job)
				return ctrl.Result{}, err
			}
			logger.V(1).Info("created job for workflow run", "job", job)
			break
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

func (r *WorkflowReconciler) isJobFinished(job *batchv1.Job) (bool, batchv1.JobConditionType) {
	for _, c := range job.Status.Conditions {
		if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == corev1.ConditionTrue {
			return true, c.Type
		}
	}
	return false, ""
}

// +kubebuilder:docs-gen:collapse=isJobFinished

func (r *WorkflowReconciler) constructConfigMapForWorkflow(workflow *workflowv1.Workflow, yamlText string, node string) (
	*corev1.ConfigMap, error) {
	name := ""
	if node == "" {
		name = fmt.Sprintf("%s-%s", workflow.Name, workflow.Spec.WorkflowYAML.Name)
	} else {
		name = fmt.Sprintf("%s-%s-%s", workflow.Name, workflow.Spec.WorkflowYAML.Name, node)
	}
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

func (r *WorkflowReconciler) constructJobForWorkflow(workflow *workflowv1.Workflow, configMap *corev1.ConfigMap,
	node string) (*batchv1.Job, error) {
	name := ""
	if node == "" {
		name = workflow.Name
	} else {
		name = fmt.Sprintf("%s-%s", workflow.Name, node)
	}
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
	if node != "" {
		job.Annotations["node"] = node
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
