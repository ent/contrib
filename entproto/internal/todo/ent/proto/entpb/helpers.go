package entpb

func toRepeatedInt(repeatedInt64 []int64) []int {
	rs := make([]int, len(repeatedInt64))

	for i, v := range repeatedInt64 {
		rs[i] = int(v)
	}

	return rs
}

func toRepeatedInt64(repeatedInt []int) []int64 {
	rs := make([]int64, len(repeatedInt))

	for i, v := range repeatedInt {
		rs[i] = int64(v)
	}

	return rs
}
