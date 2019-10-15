package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/zserge/lorca"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"time"
)

const remoteHostOrIp = "120.163.249.4"

//https://github.com/golang/go/issues/28168

func main2() {
	proxy := &httputil.ReverseProxy{
		Transport: roundTripper(rt),
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = remoteHostOrIp
			req.Header.Set("user-agent", getMacAddrMd5()) // <--- I set it here first
		},
	}
	log.Fatal(http.ListenAndServe("127.0.0.1:8888", proxy))
}

func rt(req *http.Request) (*http.Response, error) {
	log.Printf("request received. url=%s", req.URL)
	req.Header.Set("Host", "dev.tech.mojotv.cn") // <--- I set it here as well
	defer log.Printf("request complete. url=%s", req.URL)

	return http.DefaultTransport.RoundTrip(req)
}

// roundTripper makes func signature a http.RoundTripper
type roundTripper func(*http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func mustRunProxy() net.Addr {
	proxy := &httputil.ReverseProxy{
		Transport: roundTripper(rt),
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = remoteHostOrIp
			req.Header.Set("user-agent", getMacAddrSha256()) // <--- I set it here first
		},
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, proxy)
	return ln.Addr()
}

func main() {
	log.Println(getMacAddrMd5())
	// Create and bind Go object to the UI

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	proxy := &httputil.ReverseProxy{
		Transport: roundTripper(rt),
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = remoteHostOrIp
			req.Header.Set("user-agent", getMacAddrMd5()) // <--- I set it here first
		},
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	firstUrl := fmt.Sprintf("http://%s/#/login", ln.Addr())

	go func() {
		time.Sleep(time.Second * 1)
		browserOpen(firstUrl)
	}()

	log.Fatal("run proxy failed: ", http.Serve(ln, proxy))
}

func browserOpen(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func mainGuiLora() {
	log.Println(getMacAddrMd5())
	// Create and bind Go object to the UI

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	proxy := &httputil.ReverseProxy{
		Transport: roundTripper(rt),
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = remoteHostOrIp
			req.Header.Set("user-agent", getMacAddrMd5()) // <--- I set it here first
		},
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go func() {
		log.Fatal("run proxy failed: ", http.Serve(ln, proxy))
	}()
	time.Sleep(time.Second * 2)
	firstUrl := fmt.Sprintf("http://%s/#/login", ln.Addr())
	ui, err := lorca.New(firstUrl, "", 1280, 960)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Println("exiting...")
}

func getMacAddrMd5() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return ""
	}
	var macAddrs []string
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	str := strings.Join(macAddrs, "_")
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func getMacAddrSha256() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return ""
	}
	var macAddrs []string
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	str := strings.Join(macAddrs, "_")
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
