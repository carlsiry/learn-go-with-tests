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

func SumAllTails(listNum ...[]int) (sums []int) {

	for _, list := range listNum {
		if len(list) == 0 {
			sums = append(sums, 0)
		} else {
			sums = append(sums, Sum(list[1:]))
		}
	}

	return
}
