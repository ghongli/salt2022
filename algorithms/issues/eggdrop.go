package issues

var (
	memo = make(map[int]int)
)

func superEggDrop(k, n int) int {
	return dp(k, n)
}

func dp(k, n int) int {
	if _, ok := memo[n*100+k]; !ok {
		var ans int
		switch {
		case n == 0:
			ans = 0
		case k == 1:
			ans = n
		default:
			lo, hi := 1, n
			for lo+1 < hi {
				x := (lo + hi) / 2
				t1 := dp(k-1, x-1)
				t2 := dp(k, n-x)

				switch {
				case t1 < t2:
					lo = x
				case t1 > t2:
					hi = x
				default:
					hi = x
					lo = hi
				}
			}
			ans = 1 + min(max(dp(k-1, lo-1), dp(k, n-lo)), max(dp(k-1, hi-1), dp(k, n-hi)))
		}
		memo[n*100+k] = ans
	}

	return memo[n*100+k]
}

func min(lo, hi int) int {
	if lo > hi {
		return hi
	}
	return lo
}

func max(lo, hi int) int {
	if hi < lo {
		return lo
	}
	return hi
}
