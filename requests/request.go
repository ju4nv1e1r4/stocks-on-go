package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type APIResponsePrice struct {
	Meta       MetaData       `json:"meta"`
	PricesBody PricesBodyData `json:"body"`
}

type APIResponseNews struct {
        NewsBody []NewsBodyData   `json:"body"`
}

type MetaData struct {
	Version        string `json:"version"`
	Status         int    `json:"status"`
	Symbol         string `json:"symbol"`
	ProcessedTime  string `json:"processedTime"`
}

type PricesBodyData struct {
	CurrentPrice struct {
		Raw       float64 `json:"raw"`
		Fmt       string `json:"fmt"`
	} `json:"currentPrice"`
	RecommendationKey string `json:"recommendationKey"`
}

type NewsBodyData struct {
	Time   string `json:"time"`
	Ago    string `json:"ago"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Text   string `json:"text"`
	Source string `json:"source"`
}

func Start() *cli.App {
	app := cli.NewApp()
	app.Name = "Stock News"
	app.Usage = "Show last news about a company"

	flags := []cli.Flag{
			cli.StringFlag{
					Name: "ticker",
					Value: "AAPL",
			},
	}

	app.Commands = []cli.Command{
	{
		Name: "news",
		Usage: "Show last news about a company",
		Flags: flags,
		Action: News,
	},
	{
		Name: "price",
		Usage: "Show price information",
		Flags: flags,
		Action: Price,
	},
}

	return app
}

func Price(c *cli.Context) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("yellow")
	s.Start()

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	apiKey := viper.GetString("API_KEY")
	ticker := c.String("ticker")

	url := fmt.Sprintf("https://yahoo-finance15.p.rapidapi.com/api/v1/markets/stock/modules?ticker=%s&module=financial-data", ticker)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 
	}

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", "yahoo-finance15.p.rapidapi.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 
	}

	time.Sleep(2*time.Second)

	var data APIResponsePrice
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return
	}

	s.Stop()

	ticker = data.Meta.Symbol
	recommendation := data.PricesBody.RecommendationKey
	price := data.PricesBody.CurrentPrice.Raw

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Ticker", "Last Price", "Recommendation"})
	table.SetBorder(false)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_CENTER})
	table.Append([]string{
		ticker,
		fmt.Sprintf("%.2f", price),
		recommendation,
	})
	table.Render()
}

func News(c *cli.Context)  {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("yellow")
	s.Start()

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
			panic(err)
	}

	apiKey := viper.GetString("API_KEY")

	ticker := c.String("ticker")
	url := fmt.Sprintf("https://yahoo-finance15.p.rapidapi.com/api/v2/markets/news?tickers=%s&type=ALL", ticker)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
			fmt.Println("Erro ao criar requisição:", err)
			return
	}

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", "yahoo-finance15.p.rapidapi.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
			fmt.Println("Erro na requisição:", err)
			return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
			fmt.Println("Erro ao ler o corpo da resposta:", err)
			return
	}

	time.Sleep(2*time.Second)

	var data APIResponseNews
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return
	}

	s.Stop()

	colors := color.New(color.FgYellow).Add(color.Bold)
	colors.Println("## News ##")

	for i := 0; i < len(data.NewsBody) && i < 5; i++ {
			news := data.NewsBody[i]
			style := color.New(color.Bold)
			style.Println("\nTitle:", news.Title)
			style.Println("URL :", news.URL)
			style.Println("Source:", news.Source)
			style.Println("Resume:", news.Text)
			style.Println("Time:", news.Time)
			style.Println("Publish at:", news.Ago)
	}
}
