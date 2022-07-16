package cache

func GetFileInfoCacheKey(name string) string {
	return "f:" + name
}
