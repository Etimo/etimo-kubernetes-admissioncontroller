package ingresses

import (
	"encoding/json"

	"github.com/douglasmakey/admissioncontroller"

	v1 "k8s.io/api/networking/v1"
)

// NewValidationHook creates a new instance of deployment validation hook
func NewValidationHook() admissioncontroller.Hook {
	return admissioncontroller.Hook{
		Create: validateCreate(),
	}
}

func parseIngress(object []byte) (*v1.Ingress, error) {
	var dp v1.Ingress
	if err := json.Unmarshal(object, &dp); err != nil {
		return nil, err
	}

	return &dp, nil
}
