package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"

    "gopkg.in/telebot.v3"
)

type ExchangeResponse struct {
    Amount float64            `json:"amount"`
    Base   string             `json:"base"`
    Rates  map[string]float64 `json:"rates"`
}

func main() {
    bot, err := telebot.NewBot(telebot.Settings{
        Token: "7781307048:AAE_Vi4ro4o_b3ZUx2HEE9VNEfnDHc3hsaQ",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Обработчик команды /start
    bot.Handle("/start", func(c telebot.Context) error {
        return c.Send("Привет! Я бот для конвертации валют. 🤑\n\n" +
            "Используй команду:\n/convert <сумма> <из валюты> to <в валюту>\n" +
            "Пример: /convert 100 USD to EUR")
    })

    // Обработчик команды /convert
    bot.Handle("/convert", func(c telebot.Context) error {
        args := c.Args()
        if len(args) < 4 || strings.ToLower(args[len(args)-2]) != "to" {
            return c.Send("❌ Неправильный формат. Используй:\n/convert <сумма> <из валюты> to <в валюту>\nПример: /convert 100 USD to EUR")
        }

        // Парсим аргументы
        amountStr := args[0]
        fromCurrency := strings.ToUpper(args[1])
        toCurrency := strings.ToUpper(args[len(args)-1])

        // Конвертируем сумму в число
        amount, err := strconv.ParseFloat(amountStr, 64)
        if err != nil || amount <= 0 {
            return c.Send("❌ Сумма должна быть положительным числом. Пример: /convert 100 USD to RUB")
        }

        // Запрос к API
        url := fmt.Sprintf("https://api.frankfurter.app/latest?from=%s&to=%s", fromCurrency, toCurrency)
        resp, err := http.Get(url)
        if err != nil {
            return c.Send("❌ Ошибка связи с сервером валют.")
        }
        defer resp.Body.Close()

        // Парсим ответ API
        var exchangeResp ExchangeResponse
        if err := json.NewDecoder(resp.Body).Decode(&exchangeResp); err != nil {
            return c.Send("❌ Ошибка обработки данных.")
        }

        // Проверяем, есть ли нужная валюта в ответе
        rate, ok := exchangeResp.Rates[toCurrency]
        if !ok {
            return c.Send("❌ Неверные коды валют (например, USD, EUR, RUB).")
        }

        // Конвертация
        convertedAmount := amount * rate

        // Форматируем ответ
        resultMsg := fmt.Sprintf("💸 %.2f %s = *%.2f %s*", amount, fromCurrency, convertedAmount, toCurrency)
        return c.Send(resultMsg, telebot.ModeMarkdown)
    })

    // Запуск бота
    bot.Start()
}

