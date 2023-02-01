package main

import "fmt"

// Create the function fuelGauge() here
func fuelGauge(fuel int) {
	fmt.Println("Remaining fuel:", fuel)
}

// Create the function calculateFuel() here
func calculateFuel(planet string) int {
	var fuel int
	switch planet {
	case "Venus":
		fuel = 300000
	case "Mercury":
		fuel = 500000
	case "Mars":
		fuel = 700000
	default:
		fuel = 2000000
	}
	return fuel
}

// Create the function greetPlanet() here
func greetPlanet(planet string) {
	a := (planet)
	fmt.Println("Welcome to", a)
}

// Create the function cantFly() here
func cantFly() {
	fmt.Println("We do not have the available fuel to fly here.")
}

func currentPlanet(planet string) string {
	fmt.Println("You are currently on planet", planet)
	return planet
}

// Create the function flyToPlanet() here
func flyToPlanet(planet string, fuel int) int {
	var fuelRemaining int
	var fuelCost int
	fuelRemaining = fuel
	fuelCost = calculateFuel(planet)
	if fuelRemaining >= fuelCost {
		greetPlanet(planet)
		fuelRemaining = fuelRemaining - fuelCost
		currentPlanet(planet)
	} else {
		cantFly()
	}
	return fuelRemaining
}

func whereWasI(last string) string {
	fmt.Println("You were previously on", last)
	return last
}

func main() {
	// Test your functions!

	// Create `planetChoice` and `fuel`
	var fuel int
	fuel = 1000000
	for fuel > 0 {
		var planetChoice string
		var planetOn string
		fmt.Println("Which planet are we going to: ")
		fmt.Scanln(&planetChoice)
		planetOn = planetChoice
		//planetChoice = "Mars"
		// And then liftoff!
		fuel = flyToPlanet(planetChoice, fuel)
		fuelGauge(fuel)
		//currentPlanet(planetOn)

		var previousPlanet string
		previousPlanet = planetOn
		whereWasI(previousPlanet)
	}

}
