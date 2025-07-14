package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
)

func gradecalculator() {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter your name: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)
    fmt.Printf("Hello, %s!\n", name)

    fmt.Print("Enter number of subjects: ")
    var numSubjects int
    fmt.Scanln(&numSubjects)

	grade_file := map[string]int{}
	for i := 0; i < numSubjects; i++ {
		fmt.Printf("Enter name of subject %d: ", i+1)
		subject, _ := reader.ReadString('\n')
		subject = strings.TrimSpace(subject)
		fmt.Printf("Enter marks for %s: ", subject)
		var marks int
		fmt.Scanln(&marks)
		if marks < 0 || marks > 100 {	
			fmt.Printf("Invalid marks for %s. Please enter a value between 0 and 100.\n", subject)
			i-- // Decrement i to repeat this iteration
			continue
		} else {
			grade_file[subject] = marks
		}

	}
	Total := 0.0
	for _, marks := range grade_file {
		Total += float64(marks)
	}
	average := Total / float64(numSubjects)
	fmt.Printf("________________________________\n")
	fmt.Printf("\nGrade Report for %s:\n", name)
	fmt.Printf("Student Name: %s\n", name)
	
	for subject, marks := range grade_file {
		fmt.Printf("Subject: %s, Marks: %d\n", subject, marks)
	}

	fmt.Printf("Total Marks: %.2f\n", Total)
	fmt.Printf("Average Marks: %.2f\n", average)
	fmt.Printf("________________________________\n")
}

func main() {
	gradecalculator()
}