package requests

import (
        "fmt"
	"io"
	"net/http"
        "encoding/json"
        
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type FinancialData struct {
        CurrentPrice struct {
                Raw  float64 `json:"raw"`
                Fmt string  `json:"fmt"`
        } `json:"currentPrice"`
        RecommendationKey string `json:"recommendationKey"`
}

type PriceResponse struct {
        FinancialData FinancialData `json:"financialData"`
}

type NewsResponse struct {
        Link    string `json:"link"`
        PubDate string `json:"pubDate"`
        Title   string `json:"title"`
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

func Price(c *cli.Context) (PriceResponse, error) {
        viper.SetConfigFile(".env")
        err := viper.ReadInConfig()
        if err != nil {
                return PriceResponse{}, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
        }

        apiKey := viper.GetString("API_KEY")

        ticker := c.String("ticker")

        url := fmt.Sprintf("https://yahoo-finance15.p.rapidapi.com/api/v1/markets/stock/modules?ticker=%s&module=financial-data", ticker)

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
                return PriceResponse{}, fmt.Errorf("erro ao criar requisição: %w", err)
        }

        req.Header.Add("x-rapidapi-key", apiKey)
        req.Header.Add("x-rapidapi-host", "yahoo-finance15.p.rapidapi.com")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
                return PriceResponse{}, fmt.Errorf("erro na requisição: %w", err)
        }
        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil {
                return PriceResponse{}, fmt.Errorf("erro ao ler o corpo da resposta: %w", err)
        }

        var priceResponse PriceResponse
        err = json.Unmarshal(body, &priceResponse)
        if err != nil {
                return PriceResponse{}, fmt.Errorf("erro ao deserializar JSON: %w", err)
        }

        return PriceResponse{}, nil
        
}

func News(c *cli.Context)  {
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

        fmt.Println("Notícias:")
        fmt.Println(string(body))
}
