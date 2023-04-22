package utils

func Int64SliceContainEle(list []int64, v int64) bool {
	for _, i := range list {
		if i == v {
			return true
		}
	}

	return false
}
