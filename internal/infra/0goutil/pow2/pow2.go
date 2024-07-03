package pow2

func RoundToPowerOfTwo(num int) int {
	if (num & (num - 1)) == 0 {
		return num
	}
	var bitPosition = 0
	for num != 0 {
		num >>= 1
		bitPosition++
	}
	return 1 << bitPosition
}
