package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Etimo/etimo-kubernetes-admissioncontroller"
	"github.com/Etimo/etimo-kubernetes-admissioncontroller/http"

	log "k8s.io/klog/v2"
)

var (
	tlscert, tlskey, port, domain string
)

func main() {
	flag.StringVar(&tlscert, "tlscert", "", "Path to the TLS certificate")
	flag.StringVar(&tlskey, "tlskey", "", "Path to the TLS key")
	flag.StringVar(&port, "port", "8443", "The port to listen")
	flag.StringVar(&domain, "domain", ".staging.etimo-test.live", "The domain used by the cluster")
	flag.Parse()

	settings := admissioncontroller.Settings{
		Domain: domain,
	}
	server := http.NewServer(port, &settings)
	go func() {
		if tlscert == "" {
			server.ListenAndServe()
		} else if err := server.ListenAndServeTLS(tlscert, tlskey); err != nil {
			log.Errorf("Failed to listen and serve: %v", err)
		}
	}()

	log.Infof("Server running in port: %s", port)

	// listen shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Infof("Shutdown gracefully...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}
}
