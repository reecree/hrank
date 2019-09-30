package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// Complete the roadsAndLibraries function below.
func roadsAndLibraries(n int32, c_lib int32, c_road int32, cities [][]int32) int64 {
	var (
		noRoadCost      = int64(c_lib) * int64(n)
		connections     = map[int32][]int32{}
		foundCities     = map[int32]struct{}{}
		clusters, roads int64
	)
	// Create "graph"
	for _, row := range cities {
		connections[row[0]] = append(connections[row[0]], row[1])
		connections[row[1]] = append(connections[row[1]], row[0])
	}

	for i := int32(1); i <= n; i++ {
		if _, ok := foundCities[i]; ok {
			// If we've already discovered city, it was a part of a previous
			// cluster so skip
			continue
		}
		clusters++
		foundCities[i] = struct{}{}
		// Iterate through this cluster ignoring cycles via 'visitedCity'
		var (
			visitedCity = map[int32]struct{}{i: struct{}{}}
			checkCities = Stack{i}
			checkCity   int32
		)
		for len(checkCities) > 0 {
			checkCities, checkCity = checkCities.Pop()
			for _, city := range connections[checkCity] {
				if _, ok := visitedCity[city]; ok {
					// Skip cycles
					continue
				}
				//fmt.Printf("Road from %d to %d\n", checkCity, city)
				foundCities[city] = struct{}{}
				visitedCity[city] = struct{}{}
				checkCities = checkCities.Push(city)
				roads++
			}
		}
	}
	roadCost := clusters*int64(c_lib) + roads*int64(c_road)
	//fmt.Printf("NoRoad: %d. Road: %d\nClusters: %d Roads: %d\n", noRoadCost, roadCost, clusters, roads)
	return min(noRoadCost, roadCost)
}

type Stack []int32

func (s Stack) Push(v int32) Stack {
	return append(s, v)
}

func (s Stack) Pop() (Stack, int32) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 1024*1024)

	qTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	q := int32(qTemp)

	for qItr := 0; qItr < int(q); qItr++ {
		nmC_libC_road := strings.Split(readLine(reader), " ")

		nTemp, err := strconv.ParseInt(nmC_libC_road[0], 10, 64)
		checkError(err)
		n := int32(nTemp)

		mTemp, err := strconv.ParseInt(nmC_libC_road[1], 10, 64)
		checkError(err)
		m := int32(mTemp)

		c_libTemp, err := strconv.ParseInt(nmC_libC_road[2], 10, 64)
		checkError(err)
		c_lib := int32(c_libTemp)

		c_roadTemp, err := strconv.ParseInt(nmC_libC_road[3], 10, 64)
		checkError(err)
		c_road := int32(c_roadTemp)

		var cities [][]int32
		for i := 0; i < int(m); i++ {
			citiesRowTemp := strings.Split(readLine(reader), " ")

			var citiesRow []int32
			for _, citiesRowItem := range citiesRowTemp {
				citiesItemTemp, err := strconv.ParseInt(citiesRowItem, 10, 64)
				checkError(err)
				citiesItem := int32(citiesItemTemp)
				citiesRow = append(citiesRow, citiesItem)
			}

			if len(citiesRow) != 2 {
				panic("Bad input")
			}

			cities = append(cities, citiesRow)
		}

		result := roadsAndLibraries(n, c_lib, c_road, cities)

		fmt.Fprintf(writer, "%d\n", result)
	}

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
