package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"

	"cdn/global"
)

func Login(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["token"] = os.Getenv("AUTH_TOKEN")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func Me(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	w.Header().Set("Content-Type", "application/json")
	data := make(map[string]string)
	data["user"] = user.(string)
	_ = json.NewEncoder(w).Encode(data)
}

func Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func AddHosts(w http.ResponseWriter, r *http.Request) {
	type AddHostsReq struct {
		Host   string `json:"host" form:"host"`
		Origin string `json:"origin" form:"origin"`
	}
	var req AddHostsReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, ok := global.HostTargetMap.Load(req.Host)
	if ok {
		http.Error(w, "host already exists", http.StatusBadRequest)
		return
	}
	global.HostTargetMap.LoadOrStore(req.Host, req.Origin)
	_, _ = w.Write([]byte("ok"))
}

func DeleteHosts(w http.ResponseWriter, r *http.Request) {
	type DeleteHostsReq struct {
		Host string `json:"host" form:"host"`
	}
	var req DeleteHostsReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(req.Host)
	global.HostTargetMap.Delete(req.Host)
	_, _ = w.Write([]byte("ok"))
}

func UpdateHosts(w http.ResponseWriter, r *http.Request) {
	type UpdateHostsReq struct {
		Host   string `json:"host" form:"host"`
		Origin string `json:"origin" form:"origin"`
	}
	var req UpdateHostsReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, ok := global.HostTargetMap.Load(req.Host)
	if !ok {
		http.Error(w, "host not found", http.StatusBadRequest)
		return
	}
	global.HostTargetMap.Store(req.Host, req.Origin)
	_, _ = w.Write([]byte("ok"))
}

type HostOriginData struct {
	Host   string `json:"host" form:"host"`
	Origin string `json:"origin" form:"origin"`
}

func GetHosts(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := make([]*HostOriginData, 0)
	global.HostTargetMap.Range(func(key, value interface{}) bool {
		data = append(data, &HostOriginData{
			Host:   key.(string),
			Origin: value.(string),
		})
		return true
	})
	_ = json.NewEncoder(w).Encode(data)
}

type CacheKeyPathData struct {
	CacheKey  string `json:"cache_key"`
	CachePath string `json:"cache_path"`
}

func GetCached(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := make([]*CacheKeyPathData, 0)
	global.CacheKeyPath.Range(func(key, value interface{}) bool {
		data = append(data, &CacheKeyPathData{
			CacheKey:  key.(string),
			CachePath: value.(string),
		})
		return true
	})
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].CacheKey < data[j].CacheKey
	})
	_ = json.NewEncoder(w).Encode(data)
}

func PurgeCached(w http.ResponseWriter, r *http.Request) {
	type PurgeCachedReq struct {
		CacheKey string `json:"cache_key"`
	}
	var req PurgeCachedReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cachePath, ok := global.CacheKeyPath.LoadAndDelete(req.CacheKey)
	if !ok {
		http.Error(w, "cache not found", http.StatusBadRequest)
		return
	}
	err = os.Remove(cachePath.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
