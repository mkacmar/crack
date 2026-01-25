package rule

var registry = make(map[string]Rule)

func Register(r Rule) {
	registry[r.ID()] = r
}

func Get(id string) Rule {
	return registry[id]
}
