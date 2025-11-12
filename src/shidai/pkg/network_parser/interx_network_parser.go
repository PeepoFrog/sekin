package networkparser

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	interxhelper "github.com/kiracore/sekin/src/shidai/internal/interx_handler/interx_helper"
	"github.com/kiracore/sekin/src/shidai/internal/types"
	interxv2 "github.com/kiracore/sekin/src/shidai/internal/types/endpoints/interx_V2"
)

type InterxNetworkParser struct {
	mu sync.Mutex
}

func NewInterxNetworkParser() *InterxNetworkParser {
	return &InterxNetworkParser{}
}

// get nodes that are available by 11000 port
func (np *InterxNetworkParser) Scan(ctx context.Context, firstNode string, port, depth int, ignoreDepth bool) (map[string]Node, map[string]BlacklistedNode, error) {
	nodesPool := make(map[string]Node)
	blacklist := make(map[string]BlacklistedNode)
	processed := make(map[string]string)
	client := &http.Client{Timeout: 10 * time.Second} // ✅ ensure client timeout
	node, err := interxhelper.GetNetInfoV2(ctx, firstNode, port)
	if err != nil {
		return nil, nil, err
	}

	var wg sync.WaitGroup

	// 🔍 debug goroutine monitor
	go func() {
		for {
			log.Printf("[DEBUG] Goroutines active: %d", runtime.NumGoroutine())
			time.Sleep(3 * time.Second)
		}
	}()

	for _, n := range node.Peers {
		wg.Add(1)
		go np.loopFunc(ctx, &wg, client, nodesPool, blacklist, processed, n.RemoteIP, 0, depth, ignoreDepth)
	}

	wg.Wait()
	fmt.Println()
	log.Printf("\nTotal saved peers:%v\nOriginal node peer count: %v\nBlacklisted nodes(not reachable): %v\n", len(nodesPool), len(node.Peers), len(blacklist))

	return nodesPool, blacklist, nil
}

func (np *InterxNetworkParser) loopFunc(
	ctx context.Context,
	wg *sync.WaitGroup,
	client *http.Client,
	pool map[string]Node,
	blacklist map[string]BlacklistedNode,
	processed map[string]string,
	ip string,
	currentDepth, totalDepth int,
	ignoreDepth bool,
) {
	start := time.Now()
	log.Printf("[DEBUG] loopFunc start: depth=%v ip=%v", currentDepth, ip)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] panic in loopFunc(%v): %v", ip, r)
		}
		wg.Done()
		log.Printf("[DEBUG] loopFunc end: depth=%v ip=%v duration=%v", currentDepth, ip, time.Since(start))
	}()

	if !ignoreDepth && currentDepth >= totalDepth {
		log.Printf("[DEBUG] depth limit reached for %v", ip)
		return
	}

	np.mu.Lock()
	if _, exist := blacklist[ip]; exist {
		np.mu.Unlock()
		log.Printf("[DEBUG] skip %v (blacklisted)", ip)
		return
	}
	if _, exist := pool[ip]; exist {
		np.mu.Unlock()
		log.Printf("[DEBUG] skip %v (already in pool)", ip)
		return
	}
	if _, exist := processed[ip]; exist {
		np.mu.Unlock()
		log.Printf("[DEBUG] skip %v (already processed)", ip)
		return
	}
	processed[ip] = ip
	np.mu.Unlock()

	currentDepth++

	// debug timing and tracking
	log.Printf("[DEBUG] querying node %v depth=%v", ip, currentDepth)

	var nodeInfo *interxv2.NetInfo
	var status *interxv2.Status
	var errNetInfo, errStatus error

	var localWaitGroup sync.WaitGroup
	localWaitGroup.Add(2)

	go func() {
		defer localWaitGroup.Done()
		netCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		log.Printf("[DEBUG] GetNetInfoV2 start %v", ip)
		nodeInfo, errNetInfo = interxhelper.GetNetInfoV2(netCtx, ip, types.DEFAULT_INTERX_PORT)
		log.Printf("[DEBUG] GetNetInfoV2 end %v (err=%v)", ip, errNetInfo)
	}()

	go func() {
		defer localWaitGroup.Done()
		statusCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		log.Printf("[DEBUG] GetInterxStatusV2 start %v", ip)
		status, errStatus = interxhelper.GetInterxStatusV2(statusCtx, ip, types.DEFAULT_INTERX_PORT)
		log.Printf("[DEBUG] GetInterxStatusV2 end %v (err=%v)", ip, errStatus)
	}()

	localWaitGroup.Wait()
	log.Printf("[DEBUG] finished both HTTP calls for %v", ip)

	if errNetInfo != nil || errStatus != nil {
		np.mu.Lock()
		log.Printf("[DEBUG] adding <%v> to blacklist", ip)
		blacklist[ip] = BlacklistedNode{IP: ip, Error: []error{errNetInfo, errStatus}}
		cleanValue(processed, ip)
		np.mu.Unlock()
		return
	}

	np.mu.Lock()
	log.Printf("[DEBUG] adding <%v> to pool (nPeers=%v)", ip, nodeInfo.NPeers)
	npeers, err := strconv.Atoi(nodeInfo.NPeers)
	if err != nil {
		np.mu.Unlock()
		log.Printf("[WARN] failed to parse nPeers for %v", ip)
		return
	}

	node := Node{
		IP:      ip,
		NCPeers: npeers,
		ID:      status.NodeInfo.ID,
	}

	for _, nn := range nodeInfo.Peers {
		ip, port, err := extractIP(nn.NodeInfo.ListenAddr)
		if err != nil {
			continue
		}
		node.Peers = append(node.Peers, Node{IP: fmt.Sprintf("%v:%v", ip, port), ID: nn.NodeInfo.ID})
	}

	pool[ip] = node
	cleanValue(processed, ip)
	np.mu.Unlock()

	for _, p := range nodeInfo.Peers {
		wg.Add(1)
		go np.loopFunc(ctx, wg, client, pool, blacklist, processed, p.RemoteIP, currentDepth, totalDepth, ignoreDepth)

		listenAddr, _, err := extractIP(p.NodeInfo.ListenAddr)
		if err != nil {
			continue
		}
		if listenAddr != p.RemoteIP {
			log.Printf("[DEBUG] listen addr (%v) != remoteIp (%v) -> new goroutine", listenAddr, p.RemoteIP)
			wg.Add(1)
			go np.loopFunc(ctx, wg, client, pool, blacklist, processed, listenAddr, currentDepth, totalDepth, ignoreDepth)
		}
	}
}
