package main

import (
	"exchanger/internal/gui"
	"fmt"
	"os"

	"golang.org/x/sys/windows/registry"
)

// Function for adding an application to the autorun
func addToAutorun() error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	err = key.SetStringValue("CurrencyApp", exePath)
	if err != nil {
		return err
	}

	return nil
}

// Function to remove an application from the autorun
func removeFromAutorun() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	err = key.DeleteValue("CurrencyApp")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Check if there are any command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-install":
			// Add the application to the autorun
			if err := addToAutorun(); err != nil {
				fmt.Println("Error adding to autorun:", err)
			} else {
				fmt.Println("The program has been successfully added to the autorun.")
			}
			return // End the programme after execution
		case "-uninstall":
			// Remove the application from the autorun
			if err := removeFromAutorun(); err != nil {
				fmt.Println("Error when removing from autorun:", err)
			} else {
				fmt.Println("The program has been successfully removed from the autorun.")
			}
			return // End the programme after execution
		default:
			fmt.Println("Unknown team:", os.Args[1])
			return
		}
	}

	// If there are no arguments, open the main GUI
	gui.MainWindow()
}
