// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dclcontroller

import (
	"context"
	"fmt"
	"strings"
	"time"

	corekccv1alpha1 "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/core/v1alpha1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/jitter"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/lifecyclehandler"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/metrics"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/predicate"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/ratelimiter"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl"
	dclclientconfig "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/clientconfig"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/conversion"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/extension"
	dclcontainer "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/extension/container"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/kcclite"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/livestate"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/schema/dclschemaloader"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/execution"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/k8s"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/lease/leasable"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/lease/leaser"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/resourceoverrides"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/servicemapping/servicemappingloader"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/util"

	mmdcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	dclunstruct "github.com/GoogleCloudPlatform/declarative-resource-client-library/unstructured"
	"github.com/go-logr/logr"
	"github.com/nasa9084/go-openapi"
	corev1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var logger = log.Log

// The DCL lifecycles define what operations are acceptable during Apply.
// KCC doesn't allow delete-and-recreate operation to enforce modifications on immutable fields.
var LifecycleParams = []mmdcl.ApplyOption{
	mmdcl.WithLifecycleParam(mmdcl.BlockDestruction),
}

type Reconciler struct {
	lifecyclehandler.LifecycleHandler
	metrics.ReconcilerMetrics
	resourceLeaser *leaser.ResourceLeaser
	mgr            manager.Manager
	crd            *apiextensions.CustomResourceDefinition
	jsonSchema     *apiextensions.JSONSchemaProps
	gvk            schema.GroupVersionKind
	logger         logr.Logger
	// DCL related fields
	schema    *openapi.Schema
	dclConfig *mmdcl.Config
	converter *conversion.Converter
	// TF related fields
	serviceMappingLoader *servicemappingloader.ServiceMappingLoader
}

func Add(mgr manager.Manager, crd *apiextensions.CustomResourceDefinition, converter *conversion.Converter,
	dclConfig *mmdcl.Config, serviceMappingLoader *servicemappingloader.ServiceMappingLoader) error {
	kind := crd.Spec.Names.Kind
	apiVersion := k8s.GetAPIVersionFromCRD(crd)
	controllerName := fmt.Sprintf("%v-controller", strings.ToLower(kind))
	r, err := NewReconciler(mgr, crd, converter, dclConfig, serviceMappingLoader)
	if err != nil {
		return err
	}
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": apiVersion,
		},
	}
	_, err = builder.
		ControllerManagedBy(mgr).
		Named(controllerName).
		WithOptions(controller.Options{MaxConcurrentReconciles: k8s.ControllerMaxConcurrentReconciles, RateLimiter: ratelimiter.NewRateLimiter()}).
		Watches(&source.Kind{Type: obj}, &handler.EnqueueRequestForObject{}, builder.OnlyMetadata).
		For(obj, builder.OnlyMetadata).
		WithEventFilter(predicate.UnderlyingResourceOutOfSyncPredicate{}).
		Build(r)
	if err != nil {
		return fmt.Errorf("error creating new controller: %v", err)
	}
	logger.Info("Registered dcl controller", "kind", kind, "apiVersion", apiVersion)
	return nil
}

func NewReconciler(mgr manager.Manager, crd *apiextensions.CustomResourceDefinition, converter *conversion.Converter, dclConfig *mmdcl.Config, serviceMappingLoader *servicemappingloader.ServiceMappingLoader) (*Reconciler, error) {
	controllerName := fmt.Sprintf("%v-controller", strings.ToLower(crd.Spec.Names.Kind))
	gvk := schema.GroupVersionKind{
		Group:   crd.Spec.Group,
		Version: k8s.GetVersionFromCRD(crd),
		Kind:    crd.Spec.Names.Kind,
	}
	dclSchema, err := dclschemaloader.GetDCLSchemaForGVK(gvk, converter.MetadataLoader, converter.SchemaLoader)
	if err != nil {
		return nil, err
	}

	return &Reconciler{
		LifecycleHandler: lifecyclehandler.LifecycleHandler{
			Client:   mgr.GetClient(),
			Recorder: mgr.GetEventRecorderFor(controllerName),
		},
		ReconcilerMetrics: metrics.ReconcilerMetrics{
			ResourceNameLabel: metrics.ResourceNameLabel,
		},
		mgr:                  mgr,
		crd:                  crd,
		jsonSchema:           k8s.GetOpenAPIV3SchemaFromCRD(crd),
		gvk:                  gvk,
		resourceLeaser:       leaser.NewResourceLeaser(nil, nil, mgr.GetClient()),
		schema:               dclSchema,
		logger:               logger.WithName(controllerName),
		dclConfig:            dclclientconfig.CopyAndModifyForKind(dclConfig, gvk.Kind),
		converter:            converter,
		serviceMappingLoader: serviceMappingLoader,
	}, nil
}

func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, err error) {
	r.logger.Info("starting reconcile", "resource", req.NamespacedName)
	startTime := time.Now()
	r.RecordReconcileWorkers(ctx, r.gvk)
	defer r.AfterReconcile()
	defer r.RecordReconcileMetrics(ctx, r.gvk, req.Namespace, req.Name, startTime, &err)

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(r.gvk)

	if err := r.Get(ctx, req.NamespacedName, u); err != nil {
		if apierrors.IsNotFound(err) {
			r.logger.Info("resource not found in API server; finishing reconcile", "resource", req.NamespacedName)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	u, err = k8s.TriggerManagedFieldsMetadata(ctx, r.Client, u)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("error triggering Server-Side Apply (SSA) metadata: %w", err)
	}

	resource, err := dcl.NewResource(u, r.schema)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not parse resource %s: %w", req.NamespacedName.String(), err)
	}
	if err := r.applyChangesForBackwardsCompatibility(ctx, resource); err != nil {
		return reconcile.Result{}, fmt.Errorf("error applying changes to resource '%v' for backwards compatibility: %v", k8s.GetNamespacedName(resource), err)
	}
	// Apply pre-actuation transformation.
	if err := resourceoverrides.Handler.PreActuationTransform(&resource.Resource); err != nil {
		return reconcile.Result{}, r.HandlePreActuationTransformFailed(ctx, &resource.Resource, fmt.Errorf("error applying pre-actuation transformation to resource '%v': %w", req.NamespacedName.String(), err))
	}
	requeue, err := r.sync(ctx, resource)
	if err != nil {
		return reconcile.Result{}, err
	}
	if requeue {
		return reconcile.Result{Requeue: true}, nil
	}
	jitteredPeriod := jitter.GenerateJitteredReenqueuePeriod()
	r.logger.Info("successfully finished reconcile", "resource", resource.GetNamespacedName(), "time to next reconciliation", jitteredPeriod)
	return reconcile.Result{RequeueAfter: jitteredPeriod}, nil
}

func (r *Reconciler) obtainResourceLeaseIfNecessary(ctx context.Context, resource *dcl.Resource, liveLabels map[string]string) error {
	conflictPolicy, err := k8s.GetManagementConflictPreventionAnnotationValue(resource)
	if err != nil {
		return err
	}
	if conflictPolicy != k8s.ManagementConflictPreventionPolicyResource {
		return nil
	}
	ok, err := leasable.DCLSchemaSupportsLeasing(resource.Schema)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("kind '%v' does not support usage of %v of '%v'", resource.GroupVersionKind(),
			k8s.ManagementConflictPreventionPolicyAnnotation, conflictPolicy)
	}
	// Use SoftObtain instead of Obtain so that obtaining the lease ONLY changes the 'labels' value on the local krmResource and does not write the results
	// to GCP. The reason to do that is to reduce the number of writes to GCP and therefore improve performance and reduce errors.
	// The labels are written to GCP by the main sync(...) function because the changes to the labels show up in the diff.
	if err := r.resourceLeaser.SoftObtain(ctx, &resource.Resource, liveLabels); err != nil {
		return r.HandleObtainLeaseFailed(ctx, &resource.Resource, fmt.Errorf("error obtaining lease on '%v': %w",
			resource.GetNamespacedName(), err))
	}
	return nil
}

func (r *Reconciler) sync(ctx context.Context, resource *dcl.Resource) (requeue bool, err error) {
	// isolate any panics to only this function
	defer execution.RecoverWithInternalError(&err)

	dclConfig := dclclientconfig.SetUserAgentWithBlueprintAttribution(r.dclConfig, resource)
	if !resource.GetDeletionTimestamp().IsZero() {
		return r.finalizeResourceDeletion(ctx, resource, dclConfig)
	}

	liveLite, err := livestate.FetchLiveState(ctx, resource, dclConfig, r.converter, r.serviceMappingLoader, r.Client)
	if err != nil {
		if unwrappedErr, ok := lifecyclehandler.CausedByUnresolvableDeps(err); ok {
			r.logger.Info(unwrappedErr.Error(), "resource", resource.GetNamespacedName())
			return true, r.HandleUnresolvableDeps(ctx, &resource.Resource, unwrappedErr)
		}
		return false, r.HandleUpdateFailed(ctx, &resource.Resource, fmt.Errorf("error fetching live state: %w", err))
	}

	ok, err := resource.HasServerGeneratedIDAndConfigured()
	if err != nil {
		return false, r.HandleUpdateFailed(ctx, &resource.Resource, err)
	}
	if liveLite == nil && ok {
		// GCP resource with server-generated ID had been created once already,
		// but no longer exists. Don't "recreate" the resource, since
		// "recreating" resources with server-generated IDs technically creates
		// a brand new resource instead.
		return false, r.HandleUpdateFailed(ctx, &resource.Resource,
			fmt.Errorf("underlying resource with server-generated Id %s no longer exists and can't be recreated without creating a brand new resource with a different identifer", resource.Spec[k8s.ResourceIDFieldName]))
	}

	// attempt to obtain the resource lease
	var liveLabels map[string]string
	if liveLite != nil {
		liveLabels = liveLite.GetLabels()
	} else {
		liveLabels = make(map[string]string, 0)
	}
	if err = r.obtainResourceLeaseIfNecessary(ctx, resource, liveLabels); err != nil {
		return false, r.HandleObtainLeaseFailed(ctx, &resource.Resource, err)
	}
	// construct the trimmed desired state by only preserving k8s-managed fields
	desired, err := r.constructDesiredStateWithManagedFields(resource)
	if err != nil {
		return false, r.HandleUpdateFailed(ctx, &resource.Resource, err)
	}
	// KCC Full to KCC Lite
	lite, secretVersions, err := kcclite.ToKCCLiteAndSecretVersions(desired, r.converter.MetadataLoader, r.converter.SchemaLoader, r.serviceMappingLoader, r.Client)
	if err != nil {
		if unwrappedErr, ok := lifecyclehandler.CausedByUnresolvableDeps(err); ok {
			r.logger.Info(unwrappedErr.Error(), "resource", resource.GetNamespacedName())
			return true, r.HandleUnresolvableDeps(ctx, &resource.Resource, unwrappedErr)
		}
		return false, r.HandleUpdateFailed(ctx, &resource.Resource, fmt.Errorf("error converting the desired state to a KCC lite resource: %w", err))
	}
	// KCC Lite to DCL resource
	dclResource, err := r.converter.KRMObjectToDCLObject(lite)
	if err != nil {
		return false, r.HandleUpdateFailed(ctx, &resource.Resource, fmt.Errorf("error expanding resource configuration: %w", err))
	}
	// Construct the state hint to check the diffs for mutable-but-unreadable fields.
	// Note that the state hint will be ignored if the resource doesn't support
	// it so we always append it to the life cycle params.
	stateHint, err := r.getStateHint(liveLite)
	if err != nil {
		return false, r.HandleUpdateFailed(ctx, &resource.Resource, fmt.Errorf("error generating the state hint: %w", err))
	}
	stateHintApplyOption := dclunstruct.WithStateHint(stateHint)
	hasDiff, err := dclunstruct.HasDiff(ctx, dclConfig, dclResource, stateHintApplyOption)
	if err != nil {
		return false, r.HandleUpdateFailed(ctx, &resource.Resource, fmt.Errorf("error generating the diffs from desired state: %w", err))
	}

	// ensure the finalizers before apply
	if err := r.EnsureFinalizers(ctx, resource.Original, &resource.Resource, k8s.ControllerFinalizerName, k8s.DeletionDefenderFinalizerName); err != nil {
		return false, err
	}
	// check if there are diffs between the desired state and the underlying resource
	if !hasDiff {
		r.logger.Info("resource is already up to date", "resource", resource.GetNamespacedName())
		return r.updateSpecAndStatusWithLiveState(ctx, liveLite, resource, secretVersions)
	}
	// create or update the underlying resource
	r.logger.Info("creating/updating underlying resource", "resource", resource.GetNamespacedName())
	if err := r.HandleUpdating(ctx, &resource.Resource); err != nil {
		return false, err
	}
	lifecycleParams := append(LifecycleParams, stateHintApplyOption)
	newState, err := dclunstruct.Apply(ctx, dclConfig, dclResource, lifecycleParams...)
	if err != nil {
		r.logger.Error(err, "error applying desired state", "resource", resource.GetNamespacedName())
		return false, r.HandleUpdateFailed(ctx, &resource.Resource, fmt.Errorf("error applying desired state: %w", err))
	}
	// update k8s api server with the new state
	newLite, err := r.converter.DCLObjectToKRMObject(newState)
	if err != nil {
		return false, r.HandleUpdateFailed(ctx, &resource.Resource, err)
	}
	return r.updateSpecAndStatusWithLiveState(ctx, newLite, resource, secretVersions)
}

func (r *Reconciler) applyChangesForBackwardsCompatibility(ctx context.Context, resource *dcl.Resource) error {
	// Ensure the resource has a hierarchical reference. This is done to be
	// backwards compatible with resources created before the webhook for
	// defaulting hierarchical references was added.
	gvk := resource.GroupVersionKind()
	containers, err := dclcontainer.GetContainersForGVK(gvk, r.converter.MetadataLoader, r.converter.SchemaLoader)
	if err != nil {
		return fmt.Errorf("error getting containers supported by GroupVersionKind %v: %v", gvk, err)
	}
	hierarchicalRefs, err := dcl.GetHierarchicalReferencesForGVK(gvk, r.converter.MetadataLoader, r.converter.SchemaLoader)
	if err != nil {
		return fmt.Errorf("error getting hierarchical references supported by GroupVersionKind %v: %v", gvk, err)
	}
	if err := k8s.EnsureHierarchicalReference(ctx, &resource.Resource, hierarchicalRefs, containers, r.Client); err != nil {
		return fmt.Errorf("error ensuring resource '%v' has a hierarchical reference: %v", k8s.GetNamespacedName(resource), err)
	}

	// Ensure the resource has a state-into-spec annotation.
	// This is done to be backwards compatible with resources
	// created before the webhook for defaulting the annotation was added.
	if err := k8s.EnsureSpecIntoSateAnnotation(&resource.Resource); err != nil {
		return fmt.Errorf("error ensuring resource '%v' has a '%v' annotation: %v", k8s.GetNamespacedName(resource), k8s.StateIntoSpecAnnotation, err)
	}
	return nil
}

func (r *Reconciler) handleDeleted(ctx context.Context, resource *dcl.Resource) error {
	if err := resourceoverrides.Handler.PostActuationTransform(resource.Original, &resource.Resource); err != nil {
		return r.HandlePostActuationTransformFailed(ctx, &resource.Resource, fmt.Errorf("error applying post-actuation transformation to resource '%v': %w", resource.GetNamespacedName(), err))
	}
	return r.HandleDeleted(ctx, &resource.Resource)
}

func (r *Reconciler) updateSpecAndStatusWithLiveState(ctx context.Context, liveLite *unstructured.Unstructured, resource *dcl.Resource, secretVersions map[string]string) (requeue bool, err error) {
	newSpec, newStatus, err := kcclite.ResolveSpecAndStatus(liveLite, resource, r.converter.MetadataLoader)
	if err != nil {
		return false, r.HandleUpdateFailed(ctx, &resource.Resource, fmt.Errorf("error resolving the live state: %w", err))
	}
	resource.Spec, resource.Status = newSpec, newStatus

	if err := updateMutableButUnreadableFieldsAnnotationFor(resource); err != nil {
		return false, err
	}
	if err := updateObservedSecretVersionsAnnotationFor(resource, secretVersions); err != nil {
		return false, err
	}

	if err := resourceoverrides.Handler.PostActuationTransform(resource.Original, &resource.Resource); err != nil {
		return false, r.HandlePostActuationTransformFailed(ctx, &resource.Resource, fmt.Errorf("error applying post-actuation transformation to resource '%v': %w", resource.GetNamespacedName(), err))
	}

	if !k8s.IsSpecOrStatusUpdateRequired(&resource.Resource, resource.Original) &&
		!k8s.IsAnnotationsUpdateRequired(&resource.Resource, resource.Original) &&
		k8s.ReadyConditionMatches(&resource.Resource, corev1.ConditionTrue, k8s.UpToDate, k8s.UpToDateMessage) {
		return false, nil
	}

	return false, r.HandleUpToDate(ctx, &resource.Resource)
}

func (r *Reconciler) constructDesiredStateWithManagedFields(original *dcl.Resource) (*dcl.Resource, error) {
	gvk := original.GroupVersionKind()
	hierarchicalRefs, err := dcl.GetHierarchicalReferencesForGVK(gvk, r.converter.MetadataLoader, r.converter.SchemaLoader)
	if err != nil {
		return nil, fmt.Errorf("error getting hierarchical references supported by GroupVersionKind %v: %v", gvk, err)
	}
	trimmed, err := k8s.ConstructTrimmedSpecWithManagedFields(&original.Resource, r.jsonSchema, hierarchicalRefs)
	if err != nil {
		return nil, err
	}
	u, err := original.MarshalAsUnstructured()
	if err != nil {
		return nil, nil
	}
	u.Object["spec"] = trimmed
	res := &dcl.Resource{
		Schema: original.Schema,
	}
	if err := util.Marshal(u, res); err != nil {
		return nil, err
	}
	if val, ok := original.Spec[k8s.ResourceIDFieldName]; ok {
		if res.Spec == nil {
			res.Spec = make(map[string]interface{})
		}
		res.Spec[k8s.ResourceIDFieldName] = val
	}
	return res, nil
}

// finalizeResourceDeletion reacts to the KCC resource object deletion from k8s cluster.
// It performs the following operations:
// 1) checks on the finalizers before any operation
// 2) checks the deletion policy and determines whether to abandon the underlying resource
// 3) checks if the resource is orphaned by its parent
// 4) deletes the underlying resources if it owns the resource lease
func (r *Reconciler) finalizeResourceDeletion(ctx context.Context, resource *dcl.Resource, dclConfig *mmdcl.Config) (requeue bool, err error) {
	r.logger.Info("finalizing resource deletion", "resource", resource.GetNamespacedName())
	if !k8s.HasFinalizer(resource, k8s.ControllerFinalizerName) {
		r.logger.Info("no controller finalizer is present; no finalization necessary",
			"resource", resource.GetNamespacedName())
		return false, nil
	}
	if k8s.HasFinalizer(resource, k8s.DeletionDefenderFinalizerName) {
		r.logger.Info("deletion defender has not yet been finalized; requeuing", "resource", resource.GetNamespacedName())
		return true, nil
	}
	if err := r.HandleDeleting(ctx, &resource.Resource); err != nil {
		return false, err
	}
	// abandon the resource if deletion policy is abandon
	if k8s.HasAbandonAnnotation(resource) {
		r.logger.Info("deletion policy set to abandon; abandoning underlying resource", "resource", resource.GetNamespacedName())
		return false, r.handleDeleted(ctx, resource)
	}
	// check if the resource is orphaned by its parent
	orphaned, parent, err := r.isOrphaned(ctx, resource)
	if err != nil {
		return false, err
	}
	if orphaned {
		r.logger.Info("resource has been orphaned; no API call necessary", "resource", resource.GetNamespacedName())
		return false, r.handleDeleted(ctx, resource)
	}
	if parent != nil && !k8s.IsResourceReady(parent) {
		// If this resource has a parent and is not orphaned, ensure its parent
		// is ready before attempting deletion.
		return true, r.HandleUnresolvableDeps(
			ctx, &resource.Resource, k8s.NewReferenceNotReadyErrorForResource(parent))
	}

	// check if the underlying resource exists
	liveLite, err := livestate.FetchLiveState(ctx, resource, r.dclConfig, r.converter, r.serviceMappingLoader, r.Client)
	if err != nil {
		return false, r.HandleDeleteFailed(ctx, &resource.Resource, fmt.Errorf("error fetching live state: %w", err))
	}
	if liveLite == nil {
		r.logger.Info("underlying resource does not exist; no API call necessary", "resource", k8s.GetNamespacedName(resource))
		return false, r.handleDeleted(ctx, resource)
	}
	// attempt to obtain the resource lease and delete the underlying resource
	if err = r.obtainResourceLeaseIfNecessary(ctx, resource, liveLite.GetLabels()); err != nil {
		return false, r.HandleObtainLeaseFailed(ctx, &resource.Resource, err)
	}
	lite, err := kcclite.ToKCCLiteBestEffort(resource, r.converter.MetadataLoader, r.converter.SchemaLoader, r.serviceMappingLoader, r.Client)
	if err != nil {
		return false, fmt.Errorf("error converting KCC full to KCC lite: %w", err)
	}
	dclResource, err := r.converter.KRMObjectToDCLObject(lite)
	if err != nil {
		return false, fmt.Errorf("error converting KCC lite to dcl resource: %w", err)
	}
	r.logger.Info("deleting underlying resource", "resource", resource.GetNamespacedName())
	if err := dclunstruct.Delete(ctx, dclConfig, dclResource); err != nil {
		if dcl.IsNoSuchMethodError(err) {
			r.logger.Info("underlying resource cannot be deleted since there is no delete API; only clean up the kubernetes resource object", "resource", k8s.GetNamespacedName(resource))
			return false, r.handleDeleted(ctx, resource)
		}
		return false, r.HandleDeleteFailed(ctx, &resource.Resource, fmt.Errorf("error deleting the resource %v: %w", resource.GetNamespacedName(), err))
	}
	return false, r.handleDeleted(ctx, resource)
}

// isOrphaned returns whether the resource has been orphaned (i.e. its parent
// Kubernetes resource has already been deleted). Note:
// * A resource with no parent will always return false.
// * A resource whose parent is an external resource will always return false.
// * Hierarchical resources are also considered parents.
// It is assumed that parent and hierarchical references are always at the top
// level.
func (r *Reconciler) isOrphaned(ctx context.Context, resource *dcl.Resource) (orphaned bool, parent *k8s.Resource, err error) {
	gvk := resource.GroupVersionKind()
	resourceMetadata, found := r.converter.MetadataLoader.GetResourceWithGVK(gvk)
	if !found {
		return false, nil, fmt.Errorf("ServiceMetadata for resource with GroupVersionKind %v not found", gvk)
	}
	parentConfigs := make([]corekccv1alpha1.TypeConfig, 0)
	for k, s := range resource.Schema.Properties {
		if s.ReadOnly {
			continue
		}
		if !extension.IsReferenceField(s) {
			continue
		}
		if s.Type != "string" {
			continue
		}

		// TODO(b/186159460): Delete this if-block once all resources support
		// hierarchical references.
		if dcl.IsContainerField([]string{k}) && !resourceMetadata.SupportsHierarchicalReferences {
			continue
		}

		if dcl.IsMultiTypeParentReferenceField([]string{k}) {
			typeConfigs, err := dcl.GetReferenceTypeConfigs(s, r.converter.MetadataLoader)
			if err != nil {
				return false, nil, fmt.Errorf("error getting reference type configs for DCL field '%v': %w", k, err)
			}
			for _, tc := range typeConfigs {
				parentConfigs = append(parentConfigs, tc)
			}
			continue
		}

		refFieldName, err := extension.GetReferenceFieldName([]string{k}, s)
		if err != nil {
			return false, nil, err
		}
		typeConfigs, err := dcl.GetReferenceTypeConfigs(s, r.converter.MetadataLoader)
		if err != nil {
			return false, nil, fmt.Errorf("error getting reference type configs for DCL field '%v': %w", k, err)
		}
		// It is assumed that parent references and single-type hierarchical
		// references only support one resource type.
		if len(typeConfigs) != 1 {
			continue
		}
		tc := typeConfigs[0]
		if !tc.Parent {
			continue
		}
		tc.Key = refFieldName
		parentConfigs = append(parentConfigs, tc)
	}
	return lifecyclehandler.IsOrphaned(&resource.Resource, parentConfigs, r.Client)
}

// getStateHint returns a state hint based on the given resource live state. A
// state hint is a DCL unstructurd object that represents the live state of the
// resource.
//
// This function returns the given live state (if it's not nil) as a DCL
// unstructured object. Otherwise, it returns nil.
func (r *Reconciler) getStateHint(liveLite *unstructured.Unstructured) (*dclunstruct.Resource, error) {
	if liveLite == nil {
		return nil, nil
	}

	stateHint, err := r.converter.KRMObjectToDCLObject(liveLite)
	if err != nil {
		return nil, err
	}
	return stateHint, nil
}

func updateMutableButUnreadableFieldsAnnotationFor(resource *dcl.Resource) error {
	hasMutableButUnreadableFields, err := resource.HasMutableButUnreadableFields()
	if err != nil {
		return fmt.Errorf("error checking if resource has mutable-but-unreadable fields: %w", err)
	}

	// The annotation should only be set for resources with mutable-but-unreadable fields.
	if !hasMutableButUnreadableFields {
		k8s.RemoveAnnotation(k8s.MutableButUnreadableFieldsAnnotation, resource)
		return nil
	}

	annotationVal, err := dcl.MutableButUnreadableFieldsAnnotationFor(resource)
	if err != nil {
		return fmt.Errorf("error constructing value for %v: %w", k8s.MutableButUnreadableFieldsAnnotation, err)
	}
	k8s.SetAnnotation(k8s.MutableButUnreadableFieldsAnnotation, annotationVal, resource)
	return nil
}

func updateObservedSecretVersionsAnnotationFor(resource *dcl.Resource, secretVersions map[string]string) error {
	hasSensitiveFields, err := extension.HasSensitiveFields(resource.Schema)
	if err != nil {
		return fmt.Errorf("error checking if resource has sensitive fields: %w", err)
	}

	return k8s.UpdateOrRemoveObservedSecretVersionsAnnotation(&resource.Resource, secretVersions, hasSensitiveFields)
}