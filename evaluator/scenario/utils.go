package scenario

import(
	"sort"
	"fmt"
)

func sortedMapKeys[T any](flags map[string]T) []string {
	keys := make([]string, 0, len(flags))
	for key, _ := range flags {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}


func PrintSortedMap(flags map[string]string) {
	keys := sortedMapKeys(flags)
	for _, key := range keys {
		value, _ := flags[key]
		fmt.Printf(" - %s = '%s'\n", key, value)
	}
}

func errAsString(err error) string {
	if(err == nil) {
		return ""
	}
	return fmt.Sprintf("%s", err)
}
