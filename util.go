package amyssh

type StringSet map[string]struct{}

func StringsFromSet(stringSet StringSet) []string {
	strings := make([]string, 0, len(stringSet))
	for k, _ := range stringSet {
		strings = append(strings, k)
	}
	return strings
}

func SetFromList(set StringSet, list []string) StringSet {
	if set == nil {
		set = make(StringSet)
	}
	for _, v := range list {
		set[v] = struct{}{}
	}
	return set
}

func NewSet() StringSet {
	return make(StringSet)
}

func NewSetFromList(list []string) StringSet {
	set := NewSet()
	for _, v := range list {
		set[v] = struct{}{}
	}
	return set
}

func (a StringSet) Union(b StringSet) {
	for k, _ := range b {
		a[k] = struct{}{}
	}
}

func SetUnion(a StringSet, b StringSet) StringSet {
	ret := NewSet()
	for k, _ := range a {
		ret[k] = struct{}{}
	}
	for k, _ := range b {
		ret[k] = struct{}{}
	}
	return ret
}
