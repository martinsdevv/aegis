package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type Req struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestProxyHandler(t *testing.T) {
	var (
		mu            sync.Mutex
		lastPath      string
		lastRawQuery  string
		lastBodyBytes []byte
	)

	upstreamMux := http.NewServeMux()

	upstreamMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		lastPath = r.URL.Path
		lastRawQuery = r.URL.RawQuery
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("root"))
	})

	upstreamMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		lastPath = r.URL.Path
		lastRawQuery = r.URL.RawQuery
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})

	upstreamMux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		mu.Lock()
		lastPath = r.URL.Path
		lastRawQuery = r.URL.RawQuery
		lastBodyBytes = body
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	upstreamSrv := httptest.NewServer(upstreamMux)
	defer upstreamSrv.Close()

	prx := NewDynamicProxy()

	proxySrv := httptest.NewServer(HandleProxy(prx))
	defer proxySrv.Close()

	client := proxySrv.Client()

	t.Run("GET /proxy/ping rewrites to /ping", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, proxySrv.URL+"/proxy/ping", nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}

		mu.Lock()
		gotPath := lastPath
		gotQuery := lastRawQuery
		mu.Unlock()

		if gotPath != "/ping" {
			t.Fatalf("expected upstream path /ping, got %q", gotPath)
		}
		if gotQuery != "" {
			t.Fatalf("expected empty query, got %q", gotQuery)
		}
	})

	t.Run("GET /proxy/ping?x=1&y=2 preserves query", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, proxySrv.URL+"/proxy/ping?x=1&y=2", nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		mu.Lock()
		gotPath := lastPath
		gotQuery := lastRawQuery
		mu.Unlock()

		if gotPath != "/ping" {
			t.Fatalf("expected upstream path /ping, got %q", gotPath)
		}
		if gotQuery != "x=1&y=2" && gotQuery != "y=2&x=1" {
			t.Fatalf("expected query to contain x=1&y=2, got %q", gotQuery)
		}
	})

	t.Run("POST /proxy/echo passes body unchanged", func(t *testing.T) {
		data := Req{Name: "Gabriel", Age: 21}
		b, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, proxySrv.URL+"/proxy/echo", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		respBody, _ := io.ReadAll(res.Body)

		mu.Lock()
		got := lastBodyBytes
		gotPath := lastPath
		mu.Unlock()

		if gotPath != "/echo" {
			t.Fatalf("expected upstream path /echo, got %q", gotPath)
		}
		if !bytes.Equal(got, b) {
			t.Fatalf("upstream body differs.\nwant: %s\ngot:  %s", string(b), string(got))
		}

		if !bytes.Equal(respBody, b) {
			t.Fatalf("response differs.\nwant: %s\ngot:  %s", string(b), string(respBody))
		}
	})

	t.Run("GET /proxy rewrites to /", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, proxySrv.URL+"/proxy", nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		mu.Lock()
		gotPath := lastPath
		mu.Unlock()

		if gotPath != "/" {
			t.Fatalf("expected upstream path /, got %q", gotPath)
		}
		if string(body) != "root" {
			t.Fatalf("expected body root, got %q", string(body))
		}
	})
}
