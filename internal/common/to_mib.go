package common

func ToMiB(bytes uint64) float64 {
	mib := float64(bytes) / 1024 / 1024
	mib = mib * 100
	mibInt := int(mib)
	return float64(mibInt) / 100
}
