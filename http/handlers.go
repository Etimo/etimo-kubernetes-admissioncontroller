package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	admissioncontroller "github.com/Etimo/etimo-kubernetes-admissioncontroller"
	v1 "k8s.io/api/admission/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	log "k8s.io/klog/v2"
)

// admissionHandler represents the HTTP handler for an admission webhook
type admissionHandler struct {
	decoder  runtime.Decoder
	settings *admissioncontroller.Settings
}

// newAdmissionHandler returns an instance of AdmissionHandler
func newAdmissionHandler(settings *admissioncontroller.Settings) *admissionHandler {
	return &admissionHandler{
		settings: settings,
		decoder:  serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer(),
	}
}

// Serve returns a http.HandlerFunc for an admission webhook
func (h *admissionHandler) Serve(hook admissioncontroller.Hook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprint("invalid method only POST requests are allowed"), http.StatusMethodNotAllowed)
			return
		}

		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			http.Error(w, fmt.Sprint("only content type 'application/json' is supported"), http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not read request body: %v", err), http.StatusBadRequest)
			return
		}

		var review v1.AdmissionReview
		if _, _, err := h.decoder.Decode(body, nil, &review); err != nil {
			http.Error(w, fmt.Sprintf("could not deserialize request: %v", err), http.StatusBadRequest)
			return
		}

		if review.Request == nil {
			http.Error(w, "malformed admission review: request is nil", http.StatusBadRequest)
			return
		}

		result, err := hook.Execute(review.Request, h.settings)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		admissionResponse := v1.AdmissionReview{
			TypeMeta: meta.TypeMeta{
				APIVersion: "admission.k8s.io/v1",
				Kind:       "AdmissionReview",
			},
			Response: &v1.AdmissionResponse{
				UID:     review.Request.UID,
				Allowed: result.Allowed,
				Result:  &meta.Status{Message: result.Msg},
			},
		}

		// set the patch operations for mutating admission
		if len(result.PatchOps) > 0 {
			patchBytes, err := json.Marshal(result.PatchOps)
			if err != nil {
				log.Error(err)
				http.Error(w, fmt.Sprintf("could not marshal JSON patch: %v", err), http.StatusInternalServerError)
			}
			admissionResponse.Response.Patch = patchBytes
		}

		res, err := json.Marshal(admissionResponse)
		if err != nil {
			log.Error(err)
			http.Error(w, fmt.Sprintf("could not marshal response: %v", err), http.StatusInternalServerError)
			return
		}

		log.Infof("Webhook [%s - %s] - Allowed: %t", r.URL.Path, review.Request.Operation, result.Allowed)
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

func healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
