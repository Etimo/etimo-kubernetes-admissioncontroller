package http

import (
	"fmt"
	"net/http"

	"github.com/douglasmakey/admissioncontroller"
	"github.com/douglasmakey/admissioncontroller/ingresses"
)

// NewServer creates and return a http.Server
func NewServer(port string, settings *admissioncontroller.Settings) *http.Server {
	// Instances hooks
	// podsValidation := pods.NewValidationHook()
	// podsMutation := pods.NewMutationHook()
	// deploymentValidation := deployments.NewValidationHook()
	ingressValidation := ingresses.NewValidationHook()

	// Routers
	ah := newAdmissionHandler(settings)
	mux := http.NewServeMux()
	mux.Handle("/healthz", healthz())
	// mux.Handle("/validate/pods", ah.Serve(podsValidation))
	// mux.Handle("/mutate/pods", ah.Serve(podsMutation))
	// mux.Handle("/validate/deployments", ah.Serve(deploymentValidation))
	mux.Handle("/validate/ingresses", ah.Serve(ingressValidation))

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
}
