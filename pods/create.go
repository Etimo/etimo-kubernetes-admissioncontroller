package pods

import (
	"strings"

	"github.com/douglasmakey/admissioncontroller"

	v1 "k8s.io/api/admission/v1"
)

func validateCreate() admissioncontroller.AdmitFunc {
	return func(r *v1.AdmissionRequest, settings *admissioncontroller.Settings) (*admissioncontroller.Result, error) {
		pod, err := parsePod(r.Object.Raw)
		if err != nil {
			return &admissioncontroller.Result{Msg: err.Error()}, nil
		}

		for _, c := range pod.Spec.Containers {
			if strings.HasSuffix(c.Image, ":latest") {
				return &admissioncontroller.Result{Msg: "You cannot use the tag 'latest' in a container."}, nil
			}
		}

		return &admissioncontroller.Result{Allowed: true}, nil
	}
}
