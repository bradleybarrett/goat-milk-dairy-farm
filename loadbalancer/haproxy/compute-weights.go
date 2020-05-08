package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	/* Get input from command line flags. */
	inputStringPtr := flag.String("input", "", "Json string with the map of version to count:weight pairs")
	flag.Parse()

	//inputStr := '{"version=0.0.0.1": "2:1","version=0.0.0.2": "1:0","version=0.0.0.3": "3:1"}'
	inputStr := *inputStringPtr

	/* Parse map from json input string. */
	versionToCountWeightMap := make(map[string]string)
	err := json.Unmarshal([]byte(inputStr), &versionToCountWeightMap)
	handleError(err, inputStr)

	/* Populate variables used to compute the instance weight for each version. */
	versionToCount := make(map[string]int)
	versionToWeight := make(map[string]int)
	totalWeight := 0
	for version, countWeight := range versionToCountWeightMap {
		countWeightArray := strings.Split(countWeight, ":")

		count, err := strconv.Atoi(countWeightArray[0])
		handleError(err, inputStr)
		weight, err := strconv.Atoi(countWeightArray[1])
		handleError(err, inputStr)

		versionToCount[version] = count
		versionToWeight[version] = weight
		totalWeight += weight
	}

	/* Compute the instance weight for each version scaled to a certain value and add formatted string to result slice. */
	var scale int = 256 /* ha-proxy supports weights in range [0, 256] */
	versionToInstanceWeight := []string{}
	for version := range versionToCount {

		/* Allocated weight will be zero if there is no totalWeight to allocate among the instances. */
		allocatedWeight := 0
		if totalWeight > 0 {
			allocatedWeight = scale * versionToWeight[version] / totalWeight
		}

		/* Instance weight will be zero if there are no instances of the service version. */
		instanceWeight := 0
		if versionToCount[version] > 0 {
			instanceWeight = allocatedWeight / versionToCount[version]
		}

		versionToInstanceWeight = append(versionToInstanceWeight, fmt.Sprintf("%s:%d", version, instanceWeight))
	}

	/* Create the result string from instance weight slice and return the result to Stdout. */
	/* "version=1:2500,version=2:0,version=3:1666" */
	resultString := strings.Join(versionToInstanceWeight, ",")
	fmt.Fprintln(os.Stdout, fmt.Sprintf("%s", resultString))

	os.Exit(0)
}

func handleError(err error, input string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("err: %s, input: %s", err, input))
		os.Exit(1)
	}
}

/*
	"version=1": "2:1", "version=2":"1:0", "version=3": "3:1"

	weight_1 = 1
	count_1 = 2
	weight_2 = 0
	count_2 = 1
	weight_3 = 1
	count_3 = 3

	totalWeight = weight_1 + weight_2 + weight_3 = 1 + 0 + 1 = 2

	scale = 10000

	allocated_weight_1 = scale * weight_1 / totalWeight = 10000 * 1 / 2 = 5000
	instance_weight_1 = allocated_weight_1 / count_1 = 5000 / 2 = 2500

	allocated_weight_2 = scale * weight_2 / totalWeight = 10000 * 0 / 2 = 0
	instance_weight_2 = allocated_weight_2 / count_2 = 0 / 1 = 0

	allocated_weight_3 = scale * weight_3 / totalWeight = 10000 * 1 / 2 = 5000
	instance_weight_3 = allocated_weight_3 / count_3 = 5000 / 3 = 1666

	"version=1": "2500", "version=2": "0", "version=3": "1666"

	Make sure you catch edge case where all weights are 0 => computed weight should be zero in this case.

	Best solution to the below problems is to recommend that users provide weights that sum to a number less than 10000
		Effectively representing (100.00 as 10000). So user would provide v1=9999, v2=1 to route 99.99% to v1 and 0.01% to v2.
		Most users should just use weights from 0 to 100.

	have a zero numerator problem when weight > 0 AND (scale * weight) < totalWeight
		Should log this when it occurs and set the computed weight to smallest non-zero value (1).
	have a zero numerator problem when weight > 0 AND (scale * weight / totalWeight) < count
		worst case: weight = 1, totalWeight = 1, scale = 10000...
			so... (10000 * 1 / 1) < count => 10000 < count
			Should log this when it occurs that set computed weight to smallest non-zero value (1).
			"Computed weight is 0 because allocated-scaled-weight for version is less than instance count. Using computed weight = 1 for version."
*/
