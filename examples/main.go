package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/liviudnicoara/gap"
)

// CountryResponse represents the JSON response structure from the API
type CountryResponse struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
}

// Todo represents the JSON response structure from the API
type Todo struct {
	ID int `json:"id"`
}

func CallApi[T any](url string) (*T, error) {
	// Send an HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return new(T), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return new(T), fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	// Decode the JSON response
	var response *T
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return new(T), err
	}

	return response, nil
}

// GetCommonNameByAlphaCode retrieves the common name of a country by its alpha code
func GetCommonNameByAlphaCode(alphaCode string) (string, error) {
	// Construct the URL
	url := fmt.Sprintf("https://restcountries.com/v3.1/alpha/%s", alphaCode)

	countryResponses, err := CallApi[[]CountryResponse](url)
	if err != nil {
		return "", err
	}

	if len(*countryResponses) == 0 {
		return "", fmt.Errorf("No country data found for alpha code: %s", alphaCode)
	}

	return (*countryResponses)[0].Name.Common, nil
}

// GetTodoByID retrieves the "id" from the specified URL
func GetTodoByID(id int) (int, error) {

	// Construct the URL
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", id)

	todoReponse, err := CallApi[Todo](url)
	if err != nil {
		return 0, err
	}

	return todoReponse.ID, nil
}

func main() {
	defer gap.Stop()

	alphaCodes := []string{
		"USA", "CAN", "GBR", "FRA", "GER",
		"AUS", "JPN", "CHN", "BRA", "IND",
		"RUS", "MEX", "ARG", "ITA", "ESP",
		"NLD", "BEL", "SWE", "NOR", "FIN",
	}

	countryGroup := gap.NewGroup()

	for _, c := range alphaCodes {
		code := c
		countryGroup.Do(func() (interface{}, error) {
			return GetCommonNameByAlphaCode(code)
		})
	}

	todoGroup := gap.NewGroup()

	fmt.Println("Active go routines: ", gap.Running())

	for i := 1; i < 5; i++ {
		id := i
		countryGroup.Do(func() (interface{}, error) {
			return GetTodoByID(id)
		})
	}

	fmt.Println("Getting country results")
	countryResults := todoGroup.GetResults()
	for _, r := range countryResults {
		fmt.Println(r.Result)
	}

	fmt.Println("Getting todo results")
	todoResults := countryGroup.GetResults()
	for _, r := range todoResults {
		fmt.Println(r.Result)
	}
}
