package breedapi

import (
	"encoding/json"
	"fmt"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/utils"
	"net/http"
	"strings"
)

type APIResponseBreed struct {
	Name string `json:"name"`
}

const url = "https://api.thecatapi.com/v1/breeds"

func ValidateBreed(breedInput string) error {
	const op = "api.ValidateBreed"
	var apiResponse []APIResponseBreed

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return utils.ErrApiServerError
	}

	if err = json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, breed := range apiResponse {
		if strings.EqualFold(breed.Name, breedInput) {
			return nil
		}
	}

	return utils.ErrInvalidBreed
}
