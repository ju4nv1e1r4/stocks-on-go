package main

import (
        "fmt"
        "net/http"
        "io"
        "github.com/spf13/viper"
)

func Run()  {
        // Chamar as funções para obter informações de preço e notícias
        Price()
        News()
}

func Price()  {
        viper.SetConfigFile(".env")
        err := viper.ReadInConfig()
        if err != nil {
                panic(err)
        }

        apiKey := viper.GetString("API_KEY")

        url := "https://yahoo-finance15.p.rapidapi.com/api/v1/markets/stock/modules?ticker=AAPL&module=financial-data"

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

        fmt.Println("Preço:")
        fmt.Println(string(body))
}

func News()  {
        viper.SetConfigFile(".env")
        err := viper.ReadInConfig()
        if err != nil {
                panic(err)
        }

        apiKey := viper.GetString("API_KEY")

        url := "https://yahoo-finance15.p.rapidapi.com/api/v2/markets/news?tickers=AAPL&type=ALL"

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
