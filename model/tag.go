package model

type Tag struct {
	Name string `bun:"name"`
}

func Take(from, target []Tag) []Tag {
	result := from
	for _, i := range target {
		list := make([]Tag, 0)
		for _, j := range result {
			if i != j {
				list = append(list, j)
			}
		}
		result = list
	}
	return result
}
