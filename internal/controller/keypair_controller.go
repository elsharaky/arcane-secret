/*
Copyright 2025 elsharaky.

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

package controller

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/elsharaky/arcane-secret/api/v1alpha1"
	"github.com/elsharaky/arcane-secret/internal/utils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeyPairReconciler reconciles a KeyPair object
type KeyPairReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=api.arcanesecret.io,resources=keypairs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=api.arcanesecret.io,resources=keypairs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=api.arcanesecret.io,resources=keypairs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KeyPair object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *KeyPairReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	keyPair := &apiv1alpha1.KeyPair{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Get(ctx, req.NamespacedName, keyPair); err != nil {
		log.Error(err, "unable to fetch KeyPair")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	generatedKeyPair, err := utils.GenerateKeyPair(keyPair.Spec.Algorithm, keyPair.Spec.Size, keyPair.Spec.SSHFormat)
	if err != nil {
		log.Error(err, "unable to generate key pair")
		return ctrl.Result{}, err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      keyPair.Name,
			Namespace: keyPair.Namespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"privateKey": generatedKeyPair.PrivateKey,
			"publicKey":  generatedKeyPair.PublicKey,
		},
	}

	if err := ctrl.SetControllerReference(keyPair, secret, r.Scheme); err != nil {
		log.Error(err, "unable to set controller reference")
		return ctrl.Result{}, err
	}

	foundSecret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, foundSecret)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info("creating secret", "secret.Namespace", secret.Namespace, "secret.Name", secret.Name)
		err = r.Create(ctx, secret)
		if err != nil {
			log.Error(err, "unable to create secret")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "unable to get secret")
		return ctrl.Result{}, err
	}

	keyPair.Status.Conditions = []metav1.Condition{
		{
			Type:    "Generated",
			Status:  "True",
			Reason:  "KeyPairGenerated",
			Message: "Key pair generated successfully",
			LastTransitionTime: metav1.Time{
				Time: time.Now(),
			},
		},
	}

	if generatedKeyPair.Fingerprint != nil {
		keyPair.Status.Fingerprint = *generatedKeyPair.Fingerprint
	}

	if err := r.Status().Update(ctx, keyPair); err != nil {
		log.Error(err, "unable to update KeyPair status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeyPairReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.KeyPair{}).
		Named("keypair").
		Complete(r)
}
