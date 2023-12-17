package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const API_KEY = "fca_live_6RfbtFLrVplZfxCnNKFLP3mcTiC5hVwBqrKdVwTu"
const BASE_URL = "https://api.freecurrencyapi.com/v1/latest?apikey=" + API_KEY

var currencies = [...]string{"EUR", "USD", "BRL", "CAD", "AUD", "CNY"}

func clearConsole() {
	runCmd := func(name string, args ...string) {
		cmd := exec.Command(name, args...)
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	switch runtime.GOOS {
		case "linux", "darwin":
			runCmd("clear")
		case "windows":
			runCmd("cmd", "/c", "cls")
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
			if b == a {
					return true
			}
	}
	return false
}

func getAmount() float64 {
	var amount float64 = 0

	for amount <= 0 {
		fmt.Print("Insira o valor que deseja converter: ")
		fmt.Scanln(&amount)

		if amount <= 0 {
			fmt.Println("O valor inserido deve ser maior que zero!")
		}
	}

	return amount
}

func getBaseCurrency() string {
	var baseCurrency string = ""

	for !stringInSlice(baseCurrency, currencies[:]) {
		fmt.Print("Insira a moeda base: ")
		fmt.Scanln(&baseCurrency)
		baseCurrency = strings.ToUpper(baseCurrency)

		if !stringInSlice(baseCurrency, currencies[:]) {
			fmt.Println("Moeda invalida. Tente novamente")
		}
	}

	return baseCurrency
}

func convertCurrency(baseCurrency string) map[string]float64 {
	formattedTargetCurrencies := strings.Join(currencies[:], ",")

	url := fmt.Sprintf("%s&base_currency=%s&currencies=%s", BASE_URL, baseCurrency, formattedTargetCurrencies)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	data := result["data"].(map[string]interface{})
	delete(data, baseCurrency)

	resultMap := make(map[string]float64)
	for currency, value := range data {
		resultMap[currency] = value.(float64)
	}

	return resultMap
}

func printResult(data map[string]float64, amount float64) {
	for currency, value := range data {
		fmt.Printf("%s: %.2f\n", currency, value*amount)
	}
}

func main() {
	for {
		amount := getAmount()
		baseCurrency := getBaseCurrency()

		clearConsole()

		fmt.Printf("Convertendo %.2f %s:\n\n", amount, baseCurrency)

		data := convertCurrency(baseCurrency)
		if data == nil {
			continue
		}

		printResult(data, float64(amount))

		var input string
		fmt.Print("\nDeseja realizar uma nova conversão? (s/n): ")
		fmt.Scanln(&input)

		if input != "s" {
			break
		} else {
			clearConsole()
		}
	}
}