package rule

var registry = make(map[string]Rule)

func Register(r Rule) {
	registry[r.ID()] = r
}

func Get(id string) Rule {
	return registry[id]
}

func GetAll() []Rule {
	rules := make([]Rule, 0, len(registry))
	for _, r := range registry {
		rules = append(rules, r)
	}
	return rules
}
