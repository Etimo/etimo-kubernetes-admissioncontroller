package ingresses

import (
	"regexp"

	"github.com/Etimo/etimo-kubernetes-admissioncontroller"
	v1 "k8s.io/api/admission/v1"
)

func validateCreate() admissioncontroller.AdmitFunc {
	return func(r *v1.AdmissionRequest, settings *admissioncontroller.Settings) (*admissioncontroller.Result, error) {
		dp, err := parseIngress(r.Object.Raw)
		if err != nil {
			return &admissioncontroller.Result{Msg: err.Error()}, nil
		}

		if dp.Namespace == "special" {
			return &admissioncontroller.Result{Msg: "You cannot create a deployment in `special` namespace."}, nil
		}

		ns := dp.Namespace
		pattern := `^` + regexp.QuoteMeta(ns) + `(\-[^\.]+)?` + regexp.QuoteMeta(settings.Domain) + `$`
		re, err := regexp.Compile(pattern)
		for _, s := range dp.Spec.Rules {
			if !re.MatchString(s.Host) {
				return &admissioncontroller.Result{Msg: "ingress host must be " + ns + settings.Domain + " or " + ns + "-<something>" + settings.Domain}, nil
			}
		}

		return &admissioncontroller.Result{Allowed: true}, nil
	}
}
