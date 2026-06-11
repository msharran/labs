package main

import (
	"fmt"
	"slices"
)

// {
// 			name: "example from question",
// 			services: map[string][]string{
// 				"api":      {"db", "cache"},
// 				"worker":   {"db", "queue"},
// 				"frontend": {"api"},
// 				"db":       {},
// 				"cache":    {},
// 				"queue":    {},
// 				"migrate":  {"db"},
// 			},
// 			want: [][]string{
// 				{"cache", "db", "queue"},
// 				{"api", "migrate", "worker"},
// 				{"frontend"},
// 			},
// 		},
//
//
//

type serviceGraph struct {
	services map[string][]string
	edges    map[string][]string
	pending  map[string]int
	deployed map[string]bool
	ordered  [][]string
}

func computeGraph(services map[string][]string) (serviceGraph, error) {
	g := serviceGraph{
		services: services,
		edges:    make(map[string][]string),
		pending:  make(map[string]int),
		deployed: make(map[string]bool),
		ordered:  make([][]string, 0),
	}

	for k, vals := range services {
		if _, ok := g.edges[k]; !ok {
			// add k as a node, but it unlocks none as of yet
			g.edges[k] = []string{}
		}
		for _, v := range vals {
			// v unlocks k
			if _, ok := services[v]; !ok {
				return serviceGraph{}, fmt.Errorf("missing dependency %v", v)
			}
			// linking v -> k
			// check k -> v or v -> v exists and fail
			if k == v {
				return serviceGraph{}, fmt.Errorf("service %v depends on itself", k)
			}
			if slices.Contains(services[v], k) {
				return serviceGraph{}, fmt.Errorf("cycle detected between %v and %v", k, v)
			}
			g.edges[v] = append(g.edges[v], k)
		}

		g.pending[k] = len(vals)
	}
	return g, nil
}

func (g *serviceGraph) removeDeployedServices() {
	for k, count := range g.pending {
		if count == 0 && g.deployed[k] {
			delete(g.pending, k)
		}
	}
}

func (g *serviceGraph) refreshServiceDependencies() {
	for svc := range g.pending {
		for _, dep := range g.services[svc] {
			if g.deployed[dep] && g.pending[svc] > 0 {
				count := g.pending[svc]
				g.pending[svc] = count - 1
			}
		}
	}
}

func OrderServices(services map[string][]string) ([][]string, error) {
	fmt.Printf("services::::::::::::::::%v\n", services)

	// get depends_on for every service
	// start with 0 deps service
	// minus deployed deps from service depends_on list
	// now repeat again till everything is 0

	g, err := computeGraph(services)
	if err != nil {
		return nil, err
	}

	for len(g.pending) > 0 {
		currOrdered := []string{}
		for s, count := range g.pending {
			if count == 0 {
				g.deployed[s] = true
				currOrdered = append(currOrdered, s)
			}
		}
		slices.Sort(currOrdered)
		fmt.Println("deployed:", currOrdered)
		g.removeDeployedServices()
		g.refreshServiceDependencies()
		g.ordered = append(g.ordered, currOrdered)

		fmt.Println("pending services:", g.pending)
	}

	return g.ordered, nil
}

func main() {

}
