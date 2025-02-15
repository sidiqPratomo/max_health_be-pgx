package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sidiqPratomo/max-health-backend/entity"
)

func GetUnofficialDelivery(origin int64, destination int64, weight int64, courier string) ([]entity.CourierOption, error) {
	services := []entity.CourierOption{}
	url := "https://api.rajaongkir.com/starter/cost"

	param := fmt.Sprintf("origin=%d&destination=%d&weight=%d&courier=%s", origin, destination, weight, courier)
	payload := strings.NewReader(param)
	req, _ := http.NewRequest("POST", url, payload)

	key := os.Getenv("RAJA_ONGKIR_API_KEY")
	req.Header.Add("key", key)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return services, nil
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return services, nil
	}

	var raw map[string]map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	result := raw["rajaongkir"]["results"]
	if result == nil {
		return services, nil
	}

	choice := (((result.([]interface{})[0]).(map[string]interface{})["costs"]).([]interface{}))
	var costs []map[string]interface{}
	for _, cost := range choice {
		temp := ((cost.(map[string]interface{})["cost"]).([]interface{})[0]).(map[string]interface{})
		costs = append(costs, temp)
	}

	for i := range costs {
		var service entity.CourierOption
		service.Price = costs[i]["value"].(float64)
		service.Etd = costs[i]["etd"].(string)
		services = append(services, service)
	}
	return services, nil
}
