package cluster

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	httpd "github.com/open-falcon/falcon-plus/modules/rrdlite/http"
	"github.com/open-falcon/falcon-plus/modules/rrdlite/g"
	"github.com/toolkits/file"
)

const numAttempts int = 3
const attemptInterval time.Duration = 5 * time.Second

// Join attempts to join the cluster at one of the addresses given in joinAddr.
// It walks through joinAddr in order, and sets the node ID and Raft address of
// the joining node as id addr respectively. It returns the endpoint
// successfully used to join the cluster.
func Join(joinAddr []string, id, addr string, meta map[string]string, skip bool) (string, error) {
	var err error
	var j string
	logFile := file.MustOpenLogFile(g.Config().LogPath)
	logger := log.New(logFile, "[cluster-join] ", log.LstdFlags|log.Lshortfile)
	for i := 0; i < numAttempts; i++ {
		for _, a := range joinAddr {
			j, err = join(a, id, addr, meta, skip)
			if err == nil {
				// Success!
				return j, nil
			}
		}
		logger.Printf("failed to join cluster at %s, sleeping %s before retry", joinAddr, attemptInterval)
		time.Sleep(attemptInterval)
	}
	logger.Printf("failed to join cluster at %s, after %d attempts", joinAddr, numAttempts)
	return "", err
}

func join(joinAddr string, id, addr string, meta map[string]string, skip bool) (string, error) {
	if id == "" {
		return "", fmt.Errorf("node ID not set")
	}

	// Join using IP address, as that is what Hashicorp Raft works in.
	resv, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return "", err
	}

	// Check for protocol scheme, and insert default if necessary.
	fullAddr := httpd.NormalizeAddr(fmt.Sprintf("%s/join", joinAddr))

	// Enable skipVerify as requested.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skip},
	}
	client := &http.Client{Transport: tr}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	for {
		b, err := json.Marshal(map[string]interface{}{
			"id":   id,
			"addr": resv.String(),
			"meta": meta,
		})
		if err != nil {
			return "", err
		}

		// Attempt to join.
		resp, err := client.Post(fullAddr, "application-type/json", bytes.NewReader(b))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		switch resp.StatusCode {
		case http.StatusOK:
			return fullAddr, nil
		case http.StatusMovedPermanently:
			fullAddr = resp.Header.Get("location")
			if fullAddr == "" {
				return "", fmt.Errorf("failed to join, invalid redirect received")
			}
			continue
		default:
			return "", fmt.Errorf("failed to join, node returned: %s: (%s)", resp.Status, string(b))
		}
	}
}
