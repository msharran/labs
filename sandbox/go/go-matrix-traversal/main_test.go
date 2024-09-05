package main

import (
	"reflect"
	"testing"
)

func Test_isAfter(t *testing.T) {
	type args struct {
		nearbyHouse houseAddress
		i           int
		j           int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "afterMe",
			args: args{
				nearbyHouse: houseAddress{i: 2, j: 0},
				i:           0,
				j:           0,
			},
			want: true,
		},
		{
			name: "beforeMe",
			args: args{
				nearbyHouse: houseAddress{i: 2, j: 0},
				i:           3,
				j:           2,
			},
			want: false,
		},
		{
			name: "AtStart",
			args: args{
				nearbyHouse: houseAddress{i: 0, j: 0},
				i:           0,
				j:           0,
			},
			want: false,
		},
		{
			name: "DiagonallyAfter",
			args: args{
				nearbyHouse: houseAddress{i: 2, j: 2},
				i:           1,
				j:           1,
			},
			want: true,
		},
		{
			name: "DiagonallyBefore",
			args: args{
				nearbyHouse: houseAddress{i: 2, j: 2},
				i:           3,
				j:           3,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAfter(tt.args.nearbyHouse, tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("isAfter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findNearbyActivePeers(t *testing.T) {
	type args struct {
		houses [][]int
		i      int
		j      int
	}
	tests := []struct {
		name string
		args args
		want []houseAddress
	}{
		{
			name: "NoPeers",
			args: args{
				houses: [][]int{
					{1, 0, 1, 0, 1},
					{0, 0, 0, 0, 1},
					{0, 1, 0, 1, 0},
					{1, 0, 0, 0, 1},
					{1, 0, 0, 0, 1},
				},
				i: 0,
				j: 0,
			},
			want: nil,
		},
		{
			name: "HasPeers",
			args: args{
				houses: [][]int{
					{1, 0, 1, 0, 1},
					{0, 0, 0, 0, 1},
					{0, 1, 0, 1, 0},
					{1, 0, 0, 0, 1},
					{1, 0, 0, 0, 1},
				},
				i: 0,
				j: 1,
			},
			want: []houseAddress{
				{
					i: 0,
					j: 2,
				},
				{
					i: 0,
					j: 0,
				},
			},
		},
		{
			name: "HasPeers",
			args: args{
				houses: [][]int{
					{1, 0, 1, 0, 1},
					{0, 0, 0, 0, 1},
					{0, 1, 0, 1, 0},
					{1, 0, 0, 0, 1},
					{1, 0, 0, 0, 1},
				},
				i: 2,
				j: 3,
			},
			want: []houseAddress{
				{
					i: 1,
					j: 4,
				},
				{
					i: 3,
					j: 4,
				},
			},
		},
		{
			name: "HasPeers",
			args: args{
				houses: [][]int{
					{1, 0, 1, 0, 1},
					{0, 0, 0, 0, 1},
					{0, 1, 0, 1, 0},
					{1, 0, 0, 0, 1},
					{1, 0, 0, 0, 1},
				},
				i: 4,
				j: 4,
			},
			want: []houseAddress{
				{
					i: 3,
					j: 4,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findNearbyActivePeers(tt.args.houses, tt.args.i, tt.args.j); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findNearbyActivePeers() = %v, want %v", got, tt.want)
			}
		})
	}
}
