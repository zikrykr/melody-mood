package payload

type (
	GenerateRecommendationsReq struct {
		Personality string `json:"personality"`
		Genre       string `json:"genre"`
		Occasion    string `json:"occasion"`
	}
)
