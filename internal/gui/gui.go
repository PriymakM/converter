package gui

import (
	"exchanger/internal/exchanger"
	"exchanger/internal/logical"
	"exchanger/internal/logical/rates"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/getlantern/systray"
)

var w fyne.Window
var currencyFrom, currencyTo *widget.Select
var amountFrom, amountTo *widget.Entry
var ratesLabel *widget.Label
var a fyne.App

func onReady() {
	// Tray icon
	iconData, err := ioutil.ReadFile("../image/convertimg.png")
	if err != nil {
		log.Fatal("Не вдалося завантажити іконку:", err)
	}
	systray.SetIcon(iconData)
	systray.SetTitle("Currency Converter")
	systray.SetTooltip("Currency Converter in system tray")

	// Tray menu items
	todayRatesMenuItem := systray.AddMenuItem("Rates for today", "Exchange rates for today")
	monthRatesMenuItem := systray.AddMenuItem("Rates for this month", "This month's exchange rate")
	lastMonthRatesMenuItem := systray.AddMenuItem("Rates for last month", "Last month's exchange rate")
	quitMenuItem := systray.AddMenuItem("Exit", "Exit application")

	// Processing clicks on menu items
	go func() {
		for {
			select {
			case <-todayRatesMenuItem.ClickedCh:
				showTodayRates()
			case <-monthRatesMenuItem.ClickedCh:
				showMonthRates()
			case <-lastMonthRatesMenuItem.ClickedCh:
				showLastMonthRates()
			case <-quitMenuItem.ClickedCh:
				systray.Quit()
				a.Quit()
			}
		}
	}()
}

func onExit() {
	// Exit actions
	fmt.Println("Processing exit!")
}

func MainWindow() {
	// Запуск трея в отдельной горутине
	go func() {
		systray.Run(onReady, onExit)
	}()

	// Запуск Fyne приложения в главной горутине
	a = app.New()
	w = a.NewWindow("Currency Converter")
	w.Resize(fyne.NewSize(800, 800))

	// Завантажуємо іконку з файлу як []byte
	iconData, err := ioutil.ReadFile("../image/convertimgBig.png")
	if err != nil {
		log.Fatal("Could not load the icon:", err)
	}

	// Встановлюємо іконку для вікна
	w.SetIcon(fyne.NewStaticResource("icon.png", iconData))

	// Элементы интерфейса
	currencyFrom = widget.NewSelect([]string{"USD", "EUR", "UAH", "GBP"}, nil)
	currencyTo = widget.NewSelect([]string{"USD", "EUR", "UAH", "GBP"}, nil)
	amountFrom = widget.NewEntry()
	amountFrom.SetPlaceHolder("Enter the amount")
	amountTo = widget.NewEntry()
	amountTo.SetPlaceHolder("Conversion result")
	amountTo.SetReadOnly(true)
	ratesLabel = widget.NewLabel("")

	// Кнопки
	convertButton := widget.NewButton("Convert", func() {
		convertCurrency()
	})
	todayRatesButton := widget.NewButton("Exchange rates for today", func() {
		showTodayRates()
	})
	monthRatesButton := widget.NewButton("This months exchange rate", func() {
		showMonthRates()
	})
	lastMonthRatesButton := widget.NewButton("Last month's exchange rate", func() {
		showLastMonthRates()
	})
	quitButton := widget.NewButton("Exit", func() {
		a.Quit()
	})

	// Добавляем элементы в окно
	w.SetContent(container.NewVBox(
		widget.NewLabel("Convert currency"),
		container.NewGridWithColumns(2, currencyFrom, amountFrom),
		container.NewGridWithColumns(2, currencyTo, amountTo),
		convertButton,
		todayRatesButton,
		container.NewGridWithColumns(2, monthRatesButton, lastMonthRatesButton),
		quitButton,
		ratesLabel,
	))

	// Показываем основное окно приложения
	w.ShowAndRun()
}

func convertCurrency() {
	fromCurrency := currencyFrom.Selected
	toCurrency := currencyTo.Selected

	if fromCurrency == "" || toCurrency == "" {
		amountTo.SetText("Please select a currency.")
		return
	}

	amountText := amountFrom.Text
	amount, err := strconv.ParseFloat(amountText, 64)
	if err != nil {
		amountTo.SetText("Error: enter a number")
		return
	}

	exchanger.GetCurrencyRates()
	convertAmount, err := logical.ConvertCurrency(fromCurrency, toCurrency, amount)
	if err != nil {
		amountTo.SetText("Error converting currency")
		return
	}

	amountTo.SetText(strconv.FormatFloat(convertAmount, 'f', 2, 64))
}

func showTodayRates() {
	ratesMap, err := rates.LoadCurrentRatesFromFile()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ratesText := ""
	for currency, rate := range ratesMap {
		ratesText += fmt.Sprintf("%s: %.2f\n", currency, rate)
	}
	ratesLabel.SetText(ratesText)
}

func showMonthRates() {
	ratesMap, err := rates.LoadRatesFromFile("current")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ratesText := ""
	for currency, rate := range ratesMap {
		ratesText += fmt.Sprintf("%s: %.2f\n", currency, rate)
	}
	ratesLabel.SetText(ratesText)
}

func showLastMonthRates() {
	ratesMap, err := rates.LoadRatesFromFile("last")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ratesText := ""
	for currency, rate := range ratesMap {
		ratesText += fmt.Sprintf("%s: %.2f\n", currency, rate)
	}
	ratesLabel.SetText(ratesText)
}
