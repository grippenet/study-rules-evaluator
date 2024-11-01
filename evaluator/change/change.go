package change

const(
	Equal = 0
	Changed = 1
	Added = 2
	Deleted = 3
)

func CompareMap(initialFlags map[string]string, flags map[string]string) map[string]int {
	events := make(map[string]int)

	for name, initialValue := range initialFlags {
		flagValue, ok := flags[name]
		if(!ok) {
			events[name] = Deleted
		} else {
			if(initialValue != flagValue) {
				events[name] = Changed
			} else {
				events[name] = Equal
			}
		}
	}

	for name, _ := range flags {
		_, has := events[name]
		if has {
			continue
		}
		events[name] = Added
	}
	return events
}

