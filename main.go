package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"strings"
)

type Weather struct {

	PinpointLocations []PinpointLocations	`json:"pinpointLocations"`
	Forecasts []Forecasts	`json:"forecasts"`
}

type PinpointLocations struct {
	CityLink string	`json:"link"`
	CityName string `json:"name"`
}

type Forecasts struct {
	DateLabel string	`json:"dateLabel"`
	Telop string	`json:"telop"`
	DispDay string	`json:"date"`
	Temperature Temperature	`json:"temperature"`
	Image Image	`json:"image"`
}

type Temperature struct {
	MinTemperature MinTemperature    `json:"min"`
	MaxTemperature MaxTemperature	 `json:"max"`
}

type MinTemperature struct {
	MinCelsius string	`json:"celsius"`
}

type MaxTemperature struct {
	MaxCelsius string	`json:"celsius"`
}

type Image struct {
	Icon string	`json:"url"`
}

func main(){

	// ============= 天気情報取得 =============
	getUrl := "http://weather.livedoor.com/forecast/webservice/json/v1?city=130010"

	req, _ := http.NewRequest("GET", getUrl, nil)
	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	// ============= 天気情報取得 =============

	// ============= jsonパース =============
	byteArray, _ := ioutil.ReadAll(resp.Body)
	data := new(Weather)
	if err := json.Unmarshal(byteArray, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return
	}
	// ============= jsonパース =============

	// ============= 各都市の詳細Link取得 =============
	cityName := ""
	cityLink := ""
	cityName2 := ""
	cityLink2 := ""

	for _, pinpointLocation := range data.PinpointLocations {

		if pinpointLocation.CityName == "千代田区" {
			cityName = pinpointLocation.CityName
			cityLink = pinpointLocation.CityLink
		}

		if pinpointLocation.CityName == "東村山市" {
			cityName2 = pinpointLocation.CityName
			cityLink2 = pinpointLocation.CityLink
		}
	}
	// ============= 各都市の詳細Link取得 =============

	// ============= 当日、翌日の天気取得 =============
	todayDate := ""
	todayTelop := ""
	todayMinTemperature := ""
	todayMaxTemperature := ""
	todayIcon := ""
	tomorrowDate := ""
	tomorrowTelop := ""
	tomorrowIcon := ""

	for _, forecast := range data.Forecasts {

		if forecast.DateLabel == "今日" {

			todayDate = forecast.DispDay
			todayTelop = forecast.Telop
			todayMinTemperature = forecast.Temperature.MinTemperature.MinCelsius
			todayMaxTemperature = forecast.Temperature.MaxTemperature.MaxCelsius
			todayIcon = forecast.Image.Icon
		}

		if forecast.DateLabel == "明日" {

			tomorrowDate = forecast.DispDay
			tomorrowTelop = forecast.Telop
			tomorrowIcon = forecast.Image.Icon
		}
	}
	// ============= 各都市の詳細Link取得 =============

	// ============= Slack表示文字列 =============
	display := "*☆東京都の天気☆*\n" +
		"*" + todayDate + "*\n" +
		"天気： " + todayTelop + "\n" +
		"最低気温： " + todayMinTemperature + "\n" +
		"最高気温： " + todayMaxTemperature + "\n" +
		todayIcon + "\n" +
		"*" + tomorrowDate + "*\n" +
		"天気： " + tomorrowTelop + "\n" +
		tomorrowIcon + "\n" +
		"\n 各地域週間情報 \n" +
		cityName + "： " + cityLink + "\n" +
		cityName2 + "： " + cityLink2 + "\n"
	// ============= Slack表示文字列 =============

	// ============= Slack通知先設定 =============
	postUrl := "https://slack.com/api/chat.postMessage"
	values := url.Values{}
	values.Add("token", "xxxx")
	values.Add("channel", "xxxx")
	values.Add("text", display)
	// ============= Slack通知先設定 =============

	// ============= Slack通知実行 =============
	req2, _ := http.NewRequest("POST", postUrl, strings.NewReader(values.Encode()))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp2, _ := client.Do(req2)
	body2, _ := ioutil.ReadAll(resp2.Body)
	defer resp.Body.Close()
	// ============= Slack通知実行 =============

	fmt.Printf(string(body2))
	fmt.Printf("\n成功！\n")
}