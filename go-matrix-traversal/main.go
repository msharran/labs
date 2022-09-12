package main

import "fmt"

func main() {
	houses := [][]int{
		{1, 0, 1, 0, 1},
		{0, 0, 0, 0, 1},
		{0, 1, 0, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 1, 0, 0},
	}

	online := findOnline(houses)
	peer2Peer := findPeerToPeer(houses)

	fmt.Println("online downloads:", online)
	fmt.Println("peer to peer downloads:", peer2Peer)
	fmt.Println("cost optimization: ", (float64(peer2Peer)/float64(online))*100, "%")
}

func findOnline(houses [][]int) int {
	var count int
	for _, row := range houses {
		for _, currentHouse := range row {
			if currentHouse == 1 {
				count++
			}
		}
	}
	return count
}

func findPeerToPeer(houses [][]int) int {
	peer2peerNetworks := make(map[houseAddress][]houseAddress)

	for i, row := range houses {
		for j, currentHouse := range row {

			// go to next houseAddress if current houseAddress is inactive
			if currentHouse == 0 {
				continue
			}

			nearbyActiveHouses := findNearbyActivePeers(houses, i, j)

			// if there are nearbyActiveHouse active peers, check if is present in downloaded keys and values
			// if any nearbyActiveHouse is present, add the current house to that network and break the loop
			// It searches in the ascending order of nearby house addresses
			downloadedHouseAddr := getAvailableNearbyNetwork(nearbyActiveHouses, i, j, peer2peerNetworks)
			if downloadedHouseAddr == nil {
				peer2peerNetworks[houseAddress{i: i, j: j}] = []houseAddress{}
				continue
			}

			peer2peerNetworks[*downloadedHouseAddr] = append(peer2peerNetworks[*downloadedHouseAddr], houseAddress{i: i, j: j})
		}
	}
	return len(peer2peerNetworks)
}

func getAvailableNearbyNetwork(nearbyActiveHouses []houseAddress, i int, j int, peer2peerNetworks map[houseAddress][]houseAddress) *houseAddress {
	for _, nearbyActiveHouse := range nearbyActiveHouses {
		// only the houses before the current house would have the downloaded files. We can skip the houses after us
		// those will be checked in the following iterations
		if isAfter(nearbyActiveHouse, i, j) {
			continue
		}

		for downloadedHouseAddr, connectedHouseAddrs := range peer2peerNetworks {
			if nearbyActiveHouse == downloadedHouseAddr || containsHouse(connectedHouseAddrs, nearbyActiveHouse) {
				return &downloadedHouseAddr
			}
		}
	}
	return nil
}

func isAfter(nearbyHouse houseAddress, i int, j int) bool {
	return (nearbyHouse.i > i) || (nearbyHouse.i == i && nearbyHouse.j > j)
}

type houseAddress struct {
	i int
	j int
}

func findNearbyActivePeers(houses [][]int, i, j int) []houseAddress {
	var activePeers []houseAddress

	// circle from topLeft to left
	for _, index := range [][]int{{i - 1, j - 1}, {i - 1, j}, {i - 1, j + 1}, {i, j + 1}, {i + 1, j + 1}, {i + 1, j}, {i + 1, j - 1}, {i, j - 1}} {
		if index[0] >= 0 && index[1] >= 0 && index[0] < len(houses) && index[1] < len(houses) && houses[index[0]][index[1]] == 1 {
			activePeers = append(activePeers, houseAddress{i: index[0], j: index[1]})
		}
	}

	return activePeers
}

func containsHouse(houses []houseAddress, house houseAddress) bool {
	for _, h := range houses {
		if house == h {
			return true
		}
	}
	return false
}
