package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


func main() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Which task do you want to run? \n 1. Word Frequency Counter \n 2.Palindrome Checker \n")
	input.Scan()
	choice := input.Text()
	switch choice {
	case "1":
		fmt.Println("Running Word Frequency Counter...")
		WordFrequencyCounter()
	case "2":
		fmt.Println("Running Palindrome Checker...")
		PalindromeChecker()
	default:
		fmt.Println("Invalid choice. Please select 1 or 2.")

	}
}
func WordFrequencyCounter() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Println("Please provide the text for word frequency counting:")
	input.Scan()
	text := input.Text()

	wordSeparate := strings.Fields(text)
	wordCount := make(map[string]int)
	punctuation := ".,!?:;'\"-_( )[]{}<>/\\|@#$%^&*~"
	for _, word := range wordSeparate {
		word = strings.Trim(word, punctuation)
		wordCount[word]++
	}
	fmt.Println("Word Frequency Count:")
	for word, count := range wordCount {
		fmt.Printf("word: %v, count: %v\n", word, count)


	}
}
func PalindromeChecker() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Println("Please provide the text to check for palindrome:")
	input.Scan()
	text := input.Text()
	var textSlice []string
	punctuation := ".,!?:;'\"-_( )[]{}<>/\\|@#$%^&*~"
	// removing punctuation and converting to lowercase
	for _, char := range text {
		newchar := strings.ToLower(string(char))
		if !strings.ContainsRune(punctuation, rune(newchar[0])) {
			textSlice = append(textSlice, newchar)
		}
	}
	// reversing the text slice
	reversedText := make([]string, len(textSlice))
	j := len(textSlice) - 1
	for i:=0; i < len(textSlice); i++ {
		reversedText[i] = textSlice[j]
		j--
	}
	// checking palindrome
	if strings.Join(textSlice, "") == strings.Join(reversedText, "") {
		fmt.Println("The text is a palindrome.")
	} else {
		fmt.Println("The text is not a palindrome.")
	}
	fmt.Printf("word1: %v, word2: %v\n", strings.Join(textSlice, ""), strings.Join(reversedText, ""))
}