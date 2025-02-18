// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceBuilder(t *testing.T) {
	for _, test := range []string{"default", "all_set", "none_set"} {
		t.Run(test, func(t *testing.T) {
			cfg := loadResourceAttributesConfig(t, test)
			rb := NewResourceBuilder(cfg)
			rb.SetK8sDeploymentName("k8s.deployment.name-val")
			rb.SetK8sDeploymentUID("k8s.deployment.uid-val")
			rb.SetK8sNamespaceName("k8s.namespace.name-val")
			rb.SetOpencensusResourcetype("opencensus.resourcetype-val")

			res := rb.Emit()
			assert.Equal(t, 0, rb.Emit().Attributes().Len()) // Second call should return 0

			switch test {
			case "default":
				assert.Equal(t, 4, res.Attributes().Len())
			case "all_set":
				assert.Equal(t, 4, res.Attributes().Len())
			case "none_set":
				assert.Equal(t, 0, res.Attributes().Len())
				return
			default:
				assert.Failf(t, "unexpected test case: %s", test)
			}

			val, ok := res.Attributes().Get("k8s.deployment.name")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "k8s.deployment.name-val", val.Str())
			}
			val, ok = res.Attributes().Get("k8s.deployment.uid")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "k8s.deployment.uid-val", val.Str())
			}
			val, ok = res.Attributes().Get("k8s.namespace.name")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "k8s.namespace.name-val", val.Str())
			}
			val, ok = res.Attributes().Get("opencensus.resourcetype")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "opencensus.resourcetype-val", val.Str())
			}
		})
	}
}
