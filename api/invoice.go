package handler

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/ipsn/go-libtor"
)

func startTorAndGetDialer() (dialer *tor.Dialer, err error) {
	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	fmt.Println("Starting tor and fetching title of https://check.torproject.org, please wait a few seconds...")
	t, err := tor.Start(nil, &tor.StartConf{ProcessCreator: libtor.Creator, DebugWriter: os.Stderr})
	if err != nil {
		return nil, err
	}
	defer t.Close()
	// Wait at most a minute to start network and get
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()
	// Make connection
	dialer, err = t.Dialer(dialCtx, nil)
	if err != nil {
		return nil, err
	}
	return dialer, nil
}

func getInfo(dialer *tor.Dialer, onion, macaroon string) (info string, err error) {

	tr := &http.Transport{
		DialContext:     dialer.DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	fmt.Println("making request over tor")
	req, err := http.NewRequest("GET", onion, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Grpc-Metadata-macaroon", macaroon)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	return string(result), nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}
