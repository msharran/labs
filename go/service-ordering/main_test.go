package main

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestOrderServicesValidCases(t *testing.T) {
	tests := []struct {
		name     string
		services map[string][]string
		want     [][]string
	}{
		{
			name: "example from question",
			services: map[string][]string{
				"api":      {"db", "cache"},
				"worker":   {"db", "queue"},
				"frontend": {"api"},
				"db":       {},
				"cache":    {},
				"queue":    {},
				"migrate":  {"db"},
			},
			want: [][]string{
				{"cache", "db", "queue"},
				{"api", "migrate", "worker"},
				{"frontend"},
			},
		},
		{
			name: "disconnected independent service starts immediately",
			services: map[string][]string{
				"web":    {"api"},
				"api":    {"db"},
				"db":     {},
				"logger": {},
			},
			want: [][]string{
				{"db", "logger"},
				{"api"},
				{"web"},
			},
		},
		{
			name: "single service with nil dependency list",
			services: map[string][]string{
				"api": nil,
			},
			want: [][]string{{"api"}},
		},
		{
			name: "all independent services sorted alphabetically",
			services: map[string][]string{
				"worker":   {},
				"api":      {},
				"frontend": {},
				"db":       {},
				"cache":    {},
			},
			want: [][]string{{"api", "cache", "db", "frontend", "worker"}},
		},
		{
			name: "linear dependency chain",
			services: map[string][]string{
				"frontend": {"web"},
				"web":      {"api"},
				"api":      {"db"},
				"db":       {},
			},
			want: [][]string{
				{"db"},
				{"api"},
				{"web"},
				{"frontend"},
			},
		},
		{
			name: "fan out then fan in",
			services: map[string][]string{
				"build":       {},
				"lint":        {},
				"unit":        {"build"},
				"integration": {"build"},
				"package":     {"unit", "integration", "lint"},
				"deploy":      {"package"},
			},
			want: [][]string{
				{"build", "lint"},
				{"integration", "unit"},
				{"package"},
				{"deploy"},
			},
		},
		{
			name: "multiple disconnected chains merge later",
			services: map[string][]string{
				"a": {},
				"b": {"a"},
				"c": {},
				"d": {"c"},
				"e": {"b", "d"},
				"f": {},
			},
			want: [][]string{
				{"a", "c", "f"},
				{"b", "d"},
				{"e"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := OrderServices(cloneServices(tc.services))
			if err != nil {
				t.Fatalf("OrderServices returned unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("OrderServices() = %#v, want %#v", got, tc.want)
			}
			assertValidStartupStages(t, tc.services, got)
		})
	}
}

func TestOrderServicesRejectsInvalidGraphs(t *testing.T) {
	tests := []struct {
		name     string
		services map[string][]string
	}{
		{
			name: "missing dependency",
			services: map[string][]string{
				"api": {"db"},
			},
		},
		{
			name: "missing dependency among otherwise valid services",
			services: map[string][]string{
				"api":   {"db", "cache"},
				"cache": {},
			},
		},
		{
			name: "two service cycle",
			services: map[string][]string{
				"api":    {"worker"},
				"worker": {"api"},
			},
		},
		{
			name: "self cycle",
			services: map[string][]string{
				"api": {"api"},
			},
		},
		{
			name: "larger cycle",
			services: map[string][]string{
				"api":   {"db"},
				"db":    {"queue"},
				"queue": {"api"},
				"cache": {},
			},
		},
		{
			name: "cycle is rejected even with a valid component",
			services: map[string][]string{
				"db":  {},
				"api": {"db"},
				"a":   {"b"},
				"b":   {"c"},
				"c":   {"a"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := OrderServices(cloneServices(tc.services))
			if err == nil {
				t.Fatalf("OrderServices() error = nil, want non-nil error; got stages %#v", got)
			}
		})
	}
}

func TestOrderServicesLargeWideDAG(t *testing.T) {
	const n = 100

	services := make(map[string][]string, 2*n+1)
	roots := make([]string, 0, n)
	mids := make([]string, 0, n)

	for i := n - 1; i >= 0; i-- {
		root := fmt.Sprintf("root-%03d", i)
		mid := fmt.Sprintf("mid-%03d", i)

		services[root] = nil
		services[mid] = []string{root}

		roots = append(roots, root)
		mids = append(mids, mid)
	}

	services["release"] = append([]string(nil), mids...)

	sort.Strings(roots)
	sort.Strings(mids)
	want := [][]string{roots, mids, {"release"}}

	got, err := OrderServices(cloneServices(services))
	if err != nil {
		t.Fatalf("OrderServices returned unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("OrderServices() = %#v, want %#v", got, want)
	}
	assertValidStartupStages(t, services, got)
}

func cloneServices(services map[string][]string) map[string][]string {
	clone := make(map[string][]string, len(services))
	for service, deps := range services {
		clone[service] = append([]string(nil), deps...)
	}
	return clone
}

func assertValidStartupStages(t *testing.T, services map[string][]string, stages [][]string) {
	t.Helper()

	seenStage := make(map[string]int, len(services))
	for stageIndex, stage := range stages {
		if !sort.StringsAreSorted(stage) {
			t.Fatalf("stage %d is not sorted alphabetically: %#v", stageIndex, stage)
		}

		for _, service := range stage {
			if _, exists := services[service]; !exists {
				t.Fatalf("stage %d contains undefined service %q", stageIndex, service)
			}
			if previousStage, exists := seenStage[service]; exists {
				t.Fatalf("service %q appears in both stage %d and stage %d", service, previousStage, stageIndex)
			}
			seenStage[service] = stageIndex
		}
	}

	if len(seenStage) != len(services) {
		missing := make([]string, 0)
		for service := range services {
			if _, exists := seenStage[service]; !exists {
				missing = append(missing, service)
			}
		}
		sort.Strings(missing)
		t.Fatalf("missing services in output: %#v", missing)
	}

	for service, deps := range services {
		serviceStage := seenStage[service]
		for _, dep := range deps {
			depStage, exists := seenStage[dep]
			if !exists {
				t.Fatalf("service %q depends on %q, but dependency is absent from output", service, dep)
			}
			if depStage >= serviceStage {
				t.Fatalf("service %q is in stage %d but depends on %q in stage %d; stages: %#v", service, serviceStage, dep, depStage, stages)
			}
		}
	}
}
