package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gocolly/colly/v2"

	"syutuba/convert"
	"syutuba/log"
)

// ---- struct

// POSTされてくるJSONデータ構造体
type Request struct {
	Action string `json:"action"`
	RaceId string `json:"raceid"`
}

// 競走馬情報構造体
type HorseInformation struct {
	Id           int    `json:"id"`
	HorseName    string `json:"horsename"`
	StallionName string `json:"stallionname"`
	BnsName      string `json:"bnsname"`
	JockeyName   string `json:"jockeyname"`
}

// レース情報構造体
type RaceInformation struct {
	RaceText  string             `json:"racetext"`
	HorseData []HorseInformation `json:"horsedata"`
}

// レスポンス構造体
type CheckEntriesResponse struct {
	//	RaceText  string             `json:"racetext"`
	//	HorseData []HorseInformation `json:"horsedata"`
	RaceData []RaceInformation `json:"racedata"`
}

// ---- Global Variable

// ---- Package Global Variable

var checkStallion = []string{
	"パイロ",
	"ホッコータルマエ",
	"マクフィ",
	"グレーターロンドン",
}

// ---- public function ----

// ---- private function

// POSTされてくるJSONデータを構造体に変換する
func convertPostDataToStruct(inputs string) (*Request, error) {
	var req Request
	err := json.Unmarshal([]byte(inputs), &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// 1レースのurlを引数として、該当した馬の情報を返す
func checkOneRace(url string) RaceInformation {

	// Instantiate default collector
	c := colly.NewCollector()
	// Extract title element is Sample
	/*
		c.OnHTML("title", func(e *colly.HTMLElement) {
			fmt.Println("Title:", e.Text)
		})
	*/
	// Before making a request print "Visiting URL: https://XXX" is Sample
	/*
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting URL:", r.URL.String())
		})
	*/

	var raceInformationText string
	// レース詳細
	c.OnHTML(".RaceData01", func(e *colly.HTMLElement) {
		raceInformationText = e.ChildText("span")
	})

	c.OnHTML(".RaceName", func(e *colly.HTMLElement) {
		raceInformationText = fmt.Sprintf("%s:%s", e.Text, raceInformationText)
	})
	c.OnHTML(".RaceNum", func(e *colly.HTMLElement) {
		raceInformationText = fmt.Sprintf("%s:%s", e.Text, raceInformationText)
		raceInformationText = strings.ReplaceAll(raceInformationText, " ", "")
		raceInformationText = strings.ReplaceAll(raceInformationText, "\n", "")
		raceInformationText, _ = convert.EucjpToUtf8(raceInformationText)
	})

	var retValue RaceInformation
	var applicable []HorseInformation
	horseNumber := 0
	// Extract class="Horse_Info"
	c.OnHTML(".Horse_Info", func(e *colly.HTMLElement) {

		// Extract class="Horse01 fc" element
		stallionName := e.DOM.Find(".Horse01").Text()
		stallionName = strings.Trim(stallionName, "\n")
		stallionName = strings.Trim(stallionName, " ")
		stallionName, _ = convert.EucjpToUtf8(stallionName) // チェックする種牡馬名は先にUTF8化
		horseName := e.DOM.Find(".Horse02").Text()
		horseName = strings.Trim(horseName, "\n")
		horseName = strings.Trim(horseName, " ")
		bmsName := e.DOM.Find(".Horse04").Text()
		bmsName = strings.Trim(bmsName, "\n")
		bmsName = strings.Trim(bmsName, " ")
		if horseName != "" && len(horseName) > 0 {
			horseNumber++ // 馬番はインクリメント
			// 該当種牡馬の産駒かをチェック
			for _, check := range checkStallion {
				if stallionName == check {
					var single HorseInformation
					single.Id = horseNumber
					single.StallionName = stallionName
					single.HorseName, _ = convert.EucjpToUtf8(horseName)
					single.BnsName, _ = convert.EucjpToUtf8(bmsName)
					applicable = append(applicable, single)
					break
				}
			}
		}
		/*
			OnHTML(".Horse01") の場合の取得
			syuboba = e.Text
			syuboba_utf8, _ = convert.EucjpToUtf8(syuboba)
			fmt.Println(syuboba_utf8)
					name_utf8, _ = convert.EucjpToUtf8(stallion_name)
			fmt.Print(name_utf8, len(name_utf8))
		*/
		/*	colly サンプル
			// Extract href
			link, _ := e.DOM.Find("a[href]").Attr("href")
			fmt.Println(link)

			article := articleInfo{
				Title:  title,
				URL:    link,
			}

			articles = append(articles, article)
		*/
	})

	i := 0
	c.OnHTML(".Jockey", func(e *colly.HTMLElement) {
		i++
		jockeyName := e.ChildText("a")
		if jockeyName != "" && len(jockeyName) > 0 {
			// 馬番が該当馬になっているか
			for _, data := range applicable {
				if data.Id == i {
					data.JockeyName, _ = convert.EucjpToUtf8(jockeyName)
				}
			}
		}
	})

	// Start scraping on https://XXX
	c.Visit(url)

	retValue.RaceText = raceInformationText
	retValue.HorseData = applicable
	return retValue
}

func checkJraEntries(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	jsonLogger := log.GetInstance()
	jsonLogger.Info("checkJraEntries")

	req, errPost := convertPostDataToStruct(request.Body)
	//var req Request
	//var errPost error
	if errPost != nil {
		slog.Error("Convert Post Failed")
		jsonLogger.Error("Convert Post Failed", slog.String("error", errPost.Error()))
		return events.APIGatewayProxyResponse{
			Body:       "NG",
			StatusCode: 500,
		}, nil
	}

	//	var baseUrl = "https://race.netkeiba.com/race/shutuba_past.html?race_id=2024080504"
	var baseUrl = "https://race.netkeiba.com/race/shutuba_past.html?race_id="
	var url string
	var response CheckEntriesResponse
	// 12レースループ処理
	for i := 1; i <= 12; i++ {

		if i < 10 {
			url = fmt.Sprintf("%s%s0%d", baseUrl, req.RaceId, i)
		} else {
			url = fmt.Sprintf("%s%s%d", baseUrl, req.RaceId, i)
		}
		slog.Info(url)
		var single RaceInformation
		single = checkOneRace(url)
		// 該当馬いなければSkip
		if len(single.HorseData) < 1 {
			continue
		}

		response.RaceData = append(response.RaceData, single)
	}

	// json化
	var jsonBytes []byte
	var errMarshal error
	jsonBytes, errMarshal = json.Marshal(response)
	if errMarshal != nil {
		return events.APIGatewayProxyResponse{
			Body:       "NG",
			StatusCode: 500,
		}, nil
	}

	// 返り値としてレスポンスを返す
	return events.APIGatewayProxyResponse{
		Body: string(jsonBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
			// CORS対応
			/*
				"Access-Control-Allow-Headers":     "*",                       // CORS対応
				"Access-Control-Allow-Origin":      "http://localhost:5173/",  // CORS対応
				"Access-Control-Allow-Methods":     "GET, POST, PUT, OPTIONS", // CORS対応
				"Access-Control-Allow-Credentials": "true",                    // CORS対応
			*/
			// CORS対応　ここまで
		},
		StatusCode: 200,
	}, nil
}

// ---- main
func main() {
	lambda.Start(checkJraEntries)
}