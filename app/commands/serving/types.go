package serving

type ccbGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ccbMember struct {
	ID         int           `json:"group_id"`
	Individual ccbIndividual `json:"individual"`
}

type ccbIndividual struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ccbServing struct {
	ID    int    `json:"id"`
	Count int    `json:"count"`
	Start string `json:"start"`
}
