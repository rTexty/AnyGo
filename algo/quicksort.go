package main



func Find(sl []int, target int) int {
	low := 0
	high := len(sl)
	for low <= high{
		mid := low + (high- low )/ 2
		guess := sl[mid]
		if guess == target{
			return mid
		}
		if guess > target {
			high = mid -1
		} else{
			low = mid +1
		}
	}
	return -1
	
}

