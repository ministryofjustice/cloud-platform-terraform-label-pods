package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ministryofjustice/cloud-platform-label-pods/init_app"
)

type Server struct {
	http.Server
	OnHandshakeFailure func(*tls.Conn)
}

type listener struct {
	net.Listener
	onHandshakeFailure func(*tls.Conn)
}

func main() {
	ginMode := init_app.InitEnvVars()

	r := init_app.InitGin(ginMode)

	server := &Server{
		Server: http.Server{
			Addr:         ":3000",
			Handler:      r,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  5 * time.Second,
		},
		OnHandshakeFailure: onHandshakeFailure,
	}

	// to run this locally provide a self signed cert
	err := server.ListenAndServeTLS("/app/certs/tls.crt", "/app/certs/tls.key")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	addr := s.Addr
	if addr == "" {
		addr = ":https"
	}

	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer tcpListener.Close()

	cfg := &tls.Config{}
	if s.TLSConfig != nil {
		cfg = s.TLSConfig.Clone()
	}

	cfg.Certificates = make([]tls.Certificate, 1)
	cfg.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(tcpListener, cfg)

	return s.Serve(&listener{
		Listener:           tlsListener,
		onHandshakeFailure: s.OnHandshakeFailure,
	})
}

func (l *listener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return c, err
	}

	// tls.Conn tracks the handshake state so that multiple calls to
	// Handshake() are no-ops.
	// We need to modify how handshake (client bad certificate) errors are handled
	// Currently they cause the pod to hang and not accept new connections
	// A pod restart solves the problem
	if err := c.(*tls.Conn).Handshake(); err != nil {
		l.onHandshakeFailure(c.(*tls.Conn))
	}

	return c, nil
}

func onHandshakeFailure(c *tls.Conn) {
	fmt.Printf("http: TLS handshake error from %s killing pod", c.RemoteAddr().String())
	c.Close()
	os.Exit(1)
}
