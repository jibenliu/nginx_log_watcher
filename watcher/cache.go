package watcher

import (
	"fmt"
	"sync"
	"time"
)

type errorType string

const (
	IpErrType  errorType = "IP"
	UrlErrType errorType = "URL"
)

type Cache struct {
	hCache map[string]*hostCache
	rw     sync.RWMutex
}

func (c *Cache) getHostCache(hostname string) *hostCache {
	cache.rw.RLock()
	hMap := cache.hCache
	cache.rw.RUnlock()
	hCache, ok := hMap[hostname]
	if !ok {
		c.rw.Lock()
		hCache = newHostCache(hostname)
		cache.hCache[hostname] = hCache
		c.rw.Unlock()
		return hCache
	} else {
		return hCache
	}
}

type hostCache struct {
	hostname    string
	counter     map[errorType]*errorCount
	pErrorCount int //当前类型错误总数
	cErrorCount int //连续错误次数
	rw          sync.RWMutex
}

func (h *hostCache) cleanHostCError() {
	h.rw.Lock()
	defer h.rw.Unlock()
	h.cErrorCount = 0
}

func (h *hostCache) inCreTErrorCount() {
	h.rw.Lock()
	defer h.rw.Unlock()
	h.pErrorCount++
}
func (h *hostCache) inCreCErrorCount() {
	h.rw.Lock()
	defer h.rw.Unlock()
	h.cErrorCount++
}

func (h *hostCache) inCreErrorCounter(eType errorType, tag string) {
	h.rw.Lock()
	defer h.rw.Unlock()
	counter := h.counter[eType]
	counter.NodeIncrement(tag)
}

func (h *hostCache) getPErrorCounter() int {
	h.rw.RLock()
	defer h.rw.RUnlock()
	return h.pErrorCount
}
func (h *hostCache) getCErrorCounter() int {
	h.rw.RLock()
	defer h.rw.RUnlock()
	return h.cErrorCount
}

func (h *hostCache) getErrorByTypeCounter(eType errorType, tag string) int {
	h.rw.RLock()
	defer h.rw.RUnlock()
	counter := h.counter[eType]
	return counter.GetCount(tag)
}

func (h *hostCache) getTopTagError(eType errorType, top int) map[string]int {
	h.rw.RLock()
	defer h.rw.RUnlock()
	counter := h.counter[eType]
	return counter.GetTopN(top)
}

type errorCount struct {
	mapper   map[string]*countNode //实现O(1)的读取速度
	rw       sync.RWMutex
	nodeTail *countNode //存储顶点指针
}

type countNode struct {
	count int
	tag   string
	rw    sync.RWMutex
	pre   *countNode
	next  *countNode
}

func (i *errorCount) NodeIncrement(tag string) {
	i.rw.Lock()
	defer i.rw.Unlock()
	node, ok := i.mapper[tag]
	if !ok {
		node := &countNode{
			tag:   tag,
			count: 1,
		}
		tail := i.nodeTail
		i.nodeTail = node
		if tail != nil {
			tail.next = node
			node.pre = tail
		}
		i.mapper[tag] = node //追加到next节点上
	} else {
		i.mapper[tag].rw.Lock()
		i.mapper[tag].count++
		i.mapper[tag].rw.Unlock()
		for node.pre != nil && node.count > node.pre.count {
			node, node.pre = node.pre, node
		}
	}
}

func (i *errorCount) GetCount(tag string) int {
	i.rw.RLock()
	old, ok := i.mapper[tag]
	defer i.rw.RUnlock()
	if ok {
		return old.count
	}
	return 0
}

func (i *errorCount) GetTopN(in int) map[string]int {
	var ret = make(map[string]int)
	var index = 1
	tail := i.nodeTail
	for tail != nil && index <= in {
		ret[i.nodeTail.tag] = i.nodeTail.count
		tail = i.nodeTail.next
		index++
	}
	return ret
}

func newErrorCount() *errorCount {
	return &errorCount{
		rw:       sync.RWMutex{},
		mapper:   make(map[string]*countNode),
		nodeTail: nil,
	}
}

func newHostCache(hostname string) *hostCache {
	return &hostCache{
		hostname: hostname,
		rw:       sync.RWMutex{},
		counter: map[errorType]*errorCount{
			IpErrType:  newErrorCount(),
			UrlErrType: newErrorCount(),
		},
		pErrorCount: 0,
		cErrorCount: 0,
	}
}

var cache *Cache

func InitCache() {
	c := &Cache{
		rw:     sync.RWMutex{},
		hCache: make(map[string]*hostCache),
	}
	cache = c
}

func TickerFlushCache() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case t := <-ticker.C:
				Debug("flush ticker triggered at" + fmt.Sprintf(t.Format("2006-01-02 15:04:05 +08:00")))
				cache = &Cache{
					rw:     sync.RWMutex{},
					hCache: make(map[string]*hostCache),
				}
			}
		}
	}()
}
