package main

import (
	"container/heap"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
)

// main concept is to build the mapping relationship between k and state(day, planet, fuel)
// use a min-heap to maintain and compute the next state with the minimum k
// every time we want to do something, there are 3 options:
// 1. wait 1 day
// 2. refuel 1 day
// 3. jump to the neighbor
// So after every time we compute, the heap will have 2 + number_of_neighbors new states pushed in
// we keep doing this until we reach the goal state (planet == Endor && day <= countdown)
// The first time we reach the goal state, it must be with the minimum k, which will be the answer we want

// for decoding JSON configs
type Route struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	TravelTime  int    `json:"travelTime"`
}

type MFConfig struct {
	Autonomy int     `json:"autonomy"`
	Routes   []Route `json:"routes"`
}

type Hunter struct {
	Planet string `json:"planet"`
	Day    int    `json:"day"`
}

type EmpireConfig struct {
	Countdown     int      `json:"countdown"`
	BountyHunters []Hunter `json:"bounty_hunters"`
}

type C3PO struct {
	autonomy int
	graph    map[string][]edge
}

type edge struct {
	to   string
	cost int // travelTime (days)
}

func buildGraph(routes []Route) (map[string][]edge, error) {
	g := make(map[string][]edge)
	for _, r := range routes {
		if r.Origin == "" || r.Destination == "" || r.TravelTime <= 0 {
			return nil, errors.New("invalid route in millennium-falcon.json")
		}
		g[r.Origin] = append(g[r.Origin], edge{to: r.Destination, cost: r.TravelTime})
		g[r.Destination] = append(g[r.Destination], edge{to: r.Origin, cost: r.TravelTime})
	}
	return g, nil
}

func NewC3PO(millenniumFalconJsonFile string) (*C3PO, error) {
	raw, err := os.ReadFile(millenniumFalconJsonFile)
	if err != nil {
		return nil, fmt.Errorf("read mf file: %w", err)
	}
	var cfg MFConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("parse mf json: %w", err)
	}
	g, err := buildGraph(cfg.Routes)
	if err != nil {
		return nil, err
	}
	return &C3PO{
		autonomy: cfg.Autonomy,
		graph:    g,
	}, nil
}

func (c *C3PO) GiveMeTheOdds(empireJsonFile string) (float64, error) {
	// 1) read and load empire.json
	raw, err := os.ReadFile(empireJsonFile)
	if err != nil {
		return 0, fmt.Errorf("read empire file: %w", err)
	}
	var emp EmpireConfig
	if err := json.Unmarshal(raw, &emp); err != nil {
		return 0, fmt.Errorf("parse empire json: %w", err)
	}

	// 2) build bounty hunters : hunters[planet][day] = true
	hunters := make(map[string]map[int]bool)
	for _, h := range emp.BountyHunters {
		if h.Planet == "" || h.Day < 0 {
			return 0, errors.New("invalid bounty hunter record")
		}
		if hunters[h.Planet] == nil {
			hunters[h.Planet] = make(map[int]bool)
		}
		hunters[h.Planet][h.Day] = true
	}
	// planet + day = composite primary key

	// 3) run the core code to get the minimum k
	const start = "Tatooine"
	const goal = "Endor"
	minK := c.minCaptureAttempts(start, goal, c.autonomy, emp.Countdown, hunters)

	// 4) 0.9 ^ minK
	if math.IsInf(minK, 1) {
		return 0, nil // never reach the destination
	}
	if minK == 0 {
		return 1, nil // never encounter a hunter
	}
	return math.Pow(0.9, minK), nil
}

type state struct {
	day    int
	planet string
	fuel   int
}

type pqNode struct {
	k  int
	st state
}

// mini-heap, using heap in golang is really annoying
type minHeap []*pqNode

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].k < h[j].k }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(*pqNode)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func (c *C3PO) minCaptureAttempts(start, goal string, autonomy, deadline int, hunters map[string]map[int]bool) float64 {
	// state.init
	startSt := state{
		day:    0,
		planet: start,
		fuel:   autonomy,
	}

	// min-heap, sorted by k
	pq := &minHeap{}
	heap.Init(pq)
	heap.Push(pq, &pqNode{k: 0, st: startSt})

	// record the best k for each state
	best := make(map[state]int)
	best[startSt] = 0

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*pqNode)
		k := cur.k
		st := cur.st

		// discard if the current k is not the best k, It might be the old one we pushed long before
		if bk, ok := best[st]; !ok || k != bk {
			continue
		}

		// we made it, the first time we cost the least k and get to the destination
		if st.planet == goal && st.day <= deadline {
			return float64(k)
		}

		// 1) wait 1 day
		if st.day+1 <= deadline {
			ns := state{day: st.day + 1, planet: st.planet, fuel: st.fuel}
			if old, ok := best[ns]; !ok || k < old {
				best[ns] = k
				heap.Push(pq, &pqNode{k: k, st: ns}) //push the new one and continue to compute the current one
			}
		}

		// 2) refuel 1 day (if there are hunters on the refuel day, k+1)
		if st.day+1 <= deadline {
			nk := k
			if hunters[st.planet][st.day+1] {
				nk++
			}
			ns := state{day: st.day + 1, planet: st.planet, fuel: autonomy}
			if old, ok := best[ns]; !ok || nk < old {
				best[ns] = nk
				heap.Push(pq, &pqNode{k: nk, st: ns})
			}
		}

		// 3) jump to the neighbor（fuel - cost；k+1 if there is a hunter）
		for _, e := range c.graph[st.planet] {
			if e.cost > st.fuel {
				continue
			}
			arrival := st.day + e.cost
			if arrival > deadline {
				continue
			}
			nk := k
			if hunters[e.to][arrival] {
				nk++
			}
			ns := state{
				day:    arrival,
				planet: e.to,
				fuel:   st.fuel - e.cost,
			}
			if old, ok := best[ns]; !ok || nk < old {
				best[ns] = nk
				heap.Push(pq, &pqNode{k: nk, st: ns})
			}
		}
	}

	return math.Inf(1)
}
