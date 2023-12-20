package lib

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func lcmOfSlice(intArr []int) int {
	result := 1
	for _, num := range intArr {
		result = lcm(result, num)
	}
	return result
}
