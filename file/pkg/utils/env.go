package utils


func GetenvOrDefault(env, def string) string {
	if env == "" {
		return def
	}
	return env
}
