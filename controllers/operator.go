package controllers

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"cd 
)

// NamespaceReconciler reconciles a Namespace object
type NamespaceReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	AWSClient *AWSManager // Mocked AWS SDK wrapper
}

const cleanupFinalizer = "aws.cleanup.finalizer"

// Reconcile handles the lifecycle of the Namespace and its external AWS resources
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	ns := &corev1.Namespace{}

	// Fetch the Namespace instance [cite: 75]
	if err := r.Get(ctx, req.NamespacedName, ns); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// 1. Handle Deletion (Cleanup Logic) [cite: 76, 77]
	if !ns.DeletionTimestamp.IsZero() {
		if containsString(ns.Finalizers, cleanupFinalizer) {
			l.Info("Namespace is being deleted. Cleaning up external AWS resources.", "Namespace", ns.Name)

			// Call AWS SDK to delete associated RDS [cite: 79, 80]
			if err := r.AWSClient.DeleteRDS(ns.Name); err != nil {
				l.Error(err, "Failed to delete external RDS. Retrying with backoff.")
				// Returning error triggers exponential backoff via controller-runtime
				return ctrl.Result{}, err
			}

			// Remove finalizer so K8s can finish namespace deletion [cite: 85, 86]
			ns.Finalizers = removeString(ns.Finalizers, cleanupFinalizer)
			if err := r.Update(ctx, ns); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// 2. Ensure Finalizer is present (Safe Lifecycle Management) [cite: 170, 173]
	if !containsString(ns.Finalizers, cleanupFinalizer) {
		ns.Finalizers = append(ns.Finalizers, cleanupFinalizer)
		if err := r.Update(ctx, ns); err != nil {
			return ctrl.Result{}, err
		}
	}

	// 3. Day 2 Logic: Drift Detection & IRSA Provisioning [cite: 111, 163]
	// This function would check if IAM roles or security tags still exist in AWS
	if err := r.ensureExternalResources(ctx, ns); err != nil {
		l.Error(err, "Drift detected or provisioning failed. Requeuing.")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, err
	}

	return ctrl.Result{RequeueAfter: 1 * time.Hour}, nil
}

// ensureExternalResources handles logic like IRSA and Security Tags [cite: 149, 163, 164]
func (r *NamespaceReconciler) ensureExternalResources(ctx context.Context, ns *corev1.Namespace) error {
	// Provision IAM Role for Service Account (IRSA) [cite: 163]
	if err := r.AWSClient.EnsureIAMRole(ns.Name); err != nil {
		return err
	}
	// Inject mandatory labels for security compliance [cite: 160]
	if ns.Labels["enterprise-security"] != "compliant" {
		ns.Labels["enterprise-security"] = "compliant"
		return r.Update(ctx, ns)
	}
	return nil
}

// Helper functions for string slice management
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) []string {
	result := []string{}
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		Complete(r)
}
