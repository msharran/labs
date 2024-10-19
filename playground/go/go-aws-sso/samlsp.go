package main

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
)

func main() {
	keyPair, err := tls.LoadX509KeyPair("caterpillar-api.cert", "caterpillar-api.key")
	if err != nil {
		panic(err) // TODO handle error
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err) // TODO handle error
	}

	idpMetadataURL, err := url.Parse("https://portal.sso.ap-south-1.amazonaws.com/saml/metadata/NDIxNjg5Mjg0MTU5X2lucy0zNzY2MjQ4ZjhmNzBmMjYy")
	if err != nil {
		panic(err) // TODO handle error
	}
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient, *idpMetadataURL)
	if err != nil {
		panic(err) // TODO handle error
	}

	rootURL, err := url.Parse("http://localhost:65432")
	if err != nil {
		panic(err) // TODO handle error
	}

	opts := samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
		CookieName:  "token",
		SignRequest: true,
		ForceAuthn:  true,
	}
	samlmw := &samlsp.Middleware{
		ServiceProvider: samlsp.DefaultServiceProvider(opts),
		Binding:         "",
		ResponseBinding: saml.HTTPPostBinding,
		OnError:         samlsp.DefaultOnError,
		Session:         samlsp.DefaultSessionProvider(opts),
	}
	samlmw.ServiceProvider.AuthnNameIDFormat = saml.EmailAddressNameIDFormat // very important
	samlmw.RequestTracker = samlsp.DefaultRequestTracker(opts, &samlmw.ServiceProvider)
	if opts.UseArtifactResponse {
		samlmw.ResponseBinding = saml.HTTPArtifactBinding
	}

	mux := http.NewServeMux()
	mux.Handle("/saml/", samlmw)

	loginDone := make(chan string)

	loginHandler := func(w http.ResponseWriter, r *http.Request) {
		out := bytes.Buffer{}
		email := samlsp.AttributeFromContext(r.Context(), "email")
		out.WriteString("User logged in: " + email + "\n")
		fmt.Fprintf(w, "%s", out.String())
		loginDone <- email // signal that the login is done
	}
	mux.Handle("/hello", samlmw.RequireAccount(http.HandlerFunc(loginHandler)))
	server := &http.Server{
		Addr:    "localhost:65432",
		Handler: mux,
	}

	fmt.Println("Server starting on localhost:65432")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.ListenAndServe()
	}()

	<-loginDone

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	wg.Wait()

	fmt.Println("Server gracefully stopped")
}
