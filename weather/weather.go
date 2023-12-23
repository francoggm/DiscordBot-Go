package weather

import (
	"discord-bot/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

const URL string = "https://api.openweathermap.org/data/2.5/weather?"

type WeatherData struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

func makeRequest(url string) (*WeatherData, int, error) {
	client := http.Client{Timeout: 5 * time.Second}

	response, err := client.Get(url)
	if err != nil {
		return nil, response.StatusCode, err
	}

	var data WeatherData

	body, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal([]byte(body), &data)

	fmt.Printf("Response=%s : StatusCode=%d", string(body), response.StatusCode)

	return &data, response.StatusCode, nil
}

func GetWeather(msg string) *discordgo.MessageSend {
	r, _ := regexp.Compile(`\d+`)

	zip := r.FindString(msg)
	if zip == "" {
		return &discordgo.MessageSend{
			Content: "Invalid ZIP code!",
		}
	}

	cfg := config.GetConfig()

	requestURL := fmt.Sprintf("%szip=%s&lang=pt_br&appid=%s", URL, zip, cfg.OpenWeatherToken)

	data, statusCode, err := makeRequest(requestURL)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "Error getting weather, please try again!",
		}
	}

	if statusCode == http.StatusOK {
		return &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{{
				Type:        discordgo.EmbedTypeRich,
				Title:       "Current Weather",
				Description: "Weather for " + data.Name,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Conditions",
						Value:  data.Weather[0].Description,
						Inline: true,
					},
					{
						Name:   "Temperature",
						Value:  strconv.FormatFloat(data.Main.Temp, 'f', 2, 64) + "F",
						Inline: true,
					},
					{
						Name:   "Humidty",
						Value:  strconv.Itoa(data.Main.Humidity) + "%",
						Inline: true,
					},
					{
						Name:   "Wind",
						Value:  strconv.FormatFloat(data.Wind.Speed, 'f', 2, 64) + " mph",
						Inline: true,
					},
				},
			}},
		}
	} else {
		return &discordgo.MessageSend{
			Content: "Error getting weather, please try again!",
		}
	}
}
