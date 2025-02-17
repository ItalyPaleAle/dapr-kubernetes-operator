package operator

import (
	"testing"

	"k8s.io/apimachinery/pkg/types"

	"github.com/rs/xid"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	. "github.com/dapr-sandbox/dapr-kubernetes-operator/test/support"
	. "github.com/onsi/gomega"

	daprCP "github.com/dapr-sandbox/dapr-kubernetes-operator/internal/controller/operator"
	daprAc "github.com/dapr-sandbox/dapr-kubernetes-operator/pkg/client/operator/applyconfiguration/operator/v1alpha1"
	daprTC "github.com/dapr-sandbox/dapr-kubernetes-operator/test/e2e/common"
)

func TestDaprDeploy(t *testing.T) {
	test := With(t)

	instance := test.NewDaprControlPlane(daprAc.DaprControlPlaneSpec().
		WithValues(nil),
	)

	test.Eventually(Deployment(test, "dapr-operator", instance.Namespace), TestTimeoutLong).Should(
		WithTransform(ConditionStatus(appsv1.DeploymentAvailable), Equal(corev1.ConditionTrue)))
	test.Eventually(Deployment(test, "dapr-sentry", instance.Namespace), TestTimeoutLong).Should(
		WithTransform(ConditionStatus(appsv1.DeploymentAvailable), Equal(corev1.ConditionTrue)))
	test.Eventually(Deployment(test, "dapr-sidecar-injector", instance.Namespace), TestTimeoutLong).Should(
		WithTransform(ConditionStatus(appsv1.DeploymentAvailable), Equal(corev1.ConditionTrue)))

	//
	// Dapr Application
	//

	daprTC.ValidateDaprApp(test, instance.Namespace)
}

func TestDaprDeployWrongCR(t *testing.T) {
	test := With(t)

	instance := test.NewNamespacedNameDaprControlPlane(
		types.NamespacedName{
			Name:      xid.New().String(),
			Namespace: daprCP.DaprControlPlaneNamespaceDefault,
		},
		daprAc.DaprControlPlaneSpec().
			WithValues(nil),
	)

	test.Eventually(ControlPlane(test, instance), TestTimeoutLong).Should(
		WithTransform(ConditionStatus(daprCP.DaprConditionReconciled), Equal(corev1.ConditionFalse)))
	test.Eventually(ControlPlane(test, instance), TestTimeoutLong).Should(
		WithTransform(ConditionReason(daprCP.DaprConditionReconciled), Equal(daprCP.DaprConditionReasonUnsupportedConfiguration)))

}
