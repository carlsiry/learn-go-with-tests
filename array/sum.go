package array

func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}

	return sum
}

func SumAll(listNum ...[]int) (sums []int) {

	for _, list := range listNum {
		sums = append(sums, Sum(list))
	}

	return
}
