package servers

import (
	"gocache/caches"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

// HTTPServer 是 HTTP 服务器结构
type HTTPServer struct {
	// cache 是底层存储的结构
	cache *caches.Cache
}

// NewHTTPServer 返回一个关于 cache 的新 HTTP 服务器
func NewHTTPServer(cache *caches.Cache) *HTTPServer {
	return &HTTPServer{cache: cache}
}

func (hs *HTTPServer) Run(address string) error {
	return http.ListenAndServe(address, hs.routerHandler())
}

// routerHandler 返回路由处理器给 http 包中注册用
func (hs *HTTPServer) routerHandler() http.Handler {
	// httprouter.New() 创建一个 http 路由组件，包括各种请求方法的路由
	// GET 请求方法就用于缓存的查询，PUT 请求就用于缓存的新建，DELETE 请求用于缓存的删除
	// key 都从 url 上获取，value 从请求体中获取
	router := httprouter.New()
	router.GET("/cache/:key", hs.getHandler)
	router.PUT("/cache/:key", hs.setHandler)
	router.DELETE("/cache/:key", hs.deleteHandler)
	router.GET("/status", hs.statusHandler)
	return router
}

// getHandler 获取缓存数据
func (hs *HTTPServer) getHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	value, ok := hs.cache.Get(key)
	if !ok {
		// 如果缓存中找不到数据，就返回 404 状态码
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write(value)
}

// setHandler 保存缓存数据
func (hs *HTTPServer) setHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	// value 从请求体中读取，整个请求体都被当作 value
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// 如果读取请求体失败，就返回500
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hs.cache.Set(key, value)
}

// deleteHandler 用于删除缓存数据
func (hs *HTTPServer) deleteHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	hs.cache.Delete(key)
}

// statusHandler 用户获取缓存键值对的个数
func (hs *HTTPServer) statusHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// 将个数编码成 JSON 字符串
	status, err := json.Marshal(map[string]interface{}{
		"count": hs.cache.Count(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(status)
}