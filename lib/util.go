package amyssh

func StringsFromSet(stringSet map[string]struct{}) []string {
	strings := make([]string, 0, len(stringSet))
	for k, _ := range stringSet {
		strings = append(strings, k)
	}
	return strings
}

func SetFromList(set map[string]struct{}, list []string) map[string]struct{} {
	if set == nil {
		set = make(map[string]struct{})
	}
	for _, v := range list {
		set[v] = struct{}{}
	}
	return set
}
