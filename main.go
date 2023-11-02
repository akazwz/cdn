package main

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/quic-go/quic-go/http3"

	"cdn/api"
	"cdn/cache"
	"cdn/cert"
	"cdn/global"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func cacheAndProxy(w http.ResponseWriter, r *http.Request) {
	origin, ok := global.HostTargetMap.Load(r.Host)
	if !ok {
		http.Error(w, "host not found", http.StatusNotFound)
		return
	}
	// method host path
	cacheKey := fmt.Sprintf("%s %s %s", r.Method, r.Host, r.URL.Path)
	cachePath := path.Join("caches", hashString(cacheKey))
	// add cached path
	global.CacheKeyPath.LoadOrStore(cacheKey, cachePath)
	file, err := os.Open(cachePath)
	cached := false
	if err == nil {
		defer func() {
			_ = file.Close()
		}()
		func() {
			resp, err := cache.ParseCache(file)
			// cache hit
			if err == nil {
				for k, v := range resp.Header {
					w.Header().Set(k, strings.Join(v, ""))
				}
				w.Header().Set("Date", time.Now().Format(http.TimeFormat))
				w.Header().Set("X-Cache", "HIT")
				w.WriteHeader(resp.StatusCode)
				_, err = io.Copy(w, resp.Body)
				if err != nil {
					return
				}
				cached = true
				return
			}
		}()
	}
	if cached {
		return
	}
	// cache miss origin
	target, err := url.Parse(origin.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	pxy := httputil.NewSingleHostReverseProxy(target)
	pxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
	originalDirector := pxy.Director
	pxy.Director = func(req *http.Request) {
		originalDirector(req)
		// change host
		req.Host = target.Host
		req.Header.Set("Host", target.Host)
	}
	pxy.ModifyResponse = func(resp *http.Response) error {
		return cache.MakeCache(cachePath, resp)
	}
	pxy.ServeHTTP(w, r)
}

func corsMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow all origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// allow all method
		w.Header().Set("Access-Control-Allow-Methods", "*")
		// allow all header
		w.Header().Set("Access-Control-Allow-Headers", "*")
		// allow credentials
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// preflight cache 1 hour
		w.Header().Set("Access-Control-Max-Age", "3600")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func authMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if authHeader != os.Getenv("AUTH_TOKEN") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "user", "admin"))
		next.ServeHTTP(w, r)
	})
}

func spaHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets") || strings.HasPrefix(r.URL.Path, "vite.svg") {
			http.FileServer(http.Dir("web/dist")).ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, "web/dist/index.html")
	})
}

func main() {
	fmt.Println("simple cdn server")
	r := mux.NewRouter()

	a := r.PathPrefix("/api").Subrouter()
	a.Use(corsMid)
	a.Use(authMid)
	a.HandleFunc("/me", api.Me).Methods(http.MethodGet, http.MethodOptions)
	a.HandleFunc("/login", api.Login).Methods(http.MethodPost, http.MethodOptions)
	a.HandleFunc("/hosts", api.AddHosts).Methods(http.MethodPost, http.MethodOptions)
	a.HandleFunc("/hosts", api.DeleteHosts).Methods(http.MethodDelete, http.MethodOptions)
	a.HandleFunc("/hosts", api.UpdateHosts).Methods(http.MethodPut, http.MethodOptions)
	a.HandleFunc("/hosts", api.GetHosts).Methods(http.MethodGet, http.MethodOptions)
	a.HandleFunc("/cached", api.GetCached).Methods(http.MethodGet, http.MethodOptions)
	a.HandleFunc("/cached", api.PurgeCached).Methods(http.MethodDelete, http.MethodOptions)

	// health check
	r.HandleFunc("/health", api.Health).Methods(http.MethodGet)
	// dash site
	r.PathPrefix("/").Handler(spaHandler()).Host(fmt.Sprintf("dash.%s", os.Getenv("DOMAIN")))
	// proxy site
	r.PathPrefix("/").HandlerFunc(cacheAndProxy)

	certificate, err := cert.GenerateCert("certificates", fmt.Sprintf("*.%s", os.Getenv("DOMAIN")))
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{*certificate}}
	server := &http3.Server{
		TLSConfig: tlsConfig,
		Handler:   r,
	}
	// handle http2 for browsers that don't yet support HTTP/3 and add QUIC endpoint headers
	go func() {
		s := &http.Server{
			TLSConfig: tlsConfig,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err = server.SetQuicHeaders(w.Header())
				if err != nil {
					log.Fatal(err)
				}
				server.Handler.ServeHTTP(w, r)
			}),
		}
		log.Fatal(s.ListenAndServeTLS("", ""))
	}()
	// start http3 server
	log.Fatal(server.ListenAndServe())
}

func hashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
