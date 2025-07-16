package controllers

import (
	"bufio"
	"fmt"
	"library_managment/models"
	"library_managment/services"
	"os"
	"strconv"
	"strings"
)
// import (
// 	"bufio"
// 	"fmt"
// 	"library_management/models"
// 	"library_management/models"
// 	"library_management/services"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// LibraryController handles console input and invokes service methods.
type LibraryController struct {
	Service services.LibraryManager
	Reader  *bufio.Reader
}

// NewLibraryController creates and returns a new LibraryController instance.
func NewLibraryController(service services.LibraryManager) *LibraryController {
	return &LibraryController{
		Service: service,
		Reader:  bufio.NewReader(os.Stdin),
	}
}

// Run starts the console-based library management system.
func (lc *LibraryController) Run() {
	fmt.Println("Welcome to the Console-Based Library Management System!")
	for {
		lc.displayMenu()
		choice := lc.getInput("Enter your choice: ")
		switch choice {
		case "1":
			lc.addBook()
		case "2":
			lc.removeBook()
		case "3":
			lc.borrowBook()
		case "4":
			lc.returnBook()
		case "5":
			lc.listAvailableBooks()
		case "6":
			lc.listBorrowedBooks()
		case "7":
			fmt.Println("Exiting Library Management System. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
		fmt.Println()
	}
}

func (lc *LibraryController) displayMenu() {
	fmt.Println("--- Menu ---")
	fmt.Println("1. Add Book")
	fmt.Println("2. Remove Book")
	fmt.Println("3. Borrow Book")
	fmt.Println("4. Return Book")
	fmt.Println("5. List Available Books")
	fmt.Println("6. List Borrowed Books by Member")
	fmt.Println("7. Exit")
	fmt.Print("----------------\n")
}

func (lc *LibraryController) getInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := lc.Reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (lc *LibraryController) getIntInput(prompt string) (int, error) {
	inputStr := lc.getInput(prompt)
	val, err := strconv.Atoi(inputStr)
	if err != nil {
		return 0, fmt.Errorf("invalid input. Please enter a number: %w", err)
	}
	return val, nil
}

func (lc *LibraryController) addBook() {
	title := lc.getInput("Enter book title: ")
	author := lc.getInput("Enter book author: ")
	lc.Service.AddBook(models.Book{Title: title, Author: author})
	fmt.Println("Book added successfully!")
}

func (lc *LibraryController) removeBook() {
	bookID, err := lc.getIntInput("Enter book ID to remove: ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = lc.Service.RemoveBook(bookID)
	if err != nil {
		fmt.Println("Error removing book:", err)
		return
	}
	fmt.Println("Book removed successfully!")
}

func (lc *LibraryController) borrowBook() {
	bookID, err := lc.getIntInput("Enter book ID to borrow: ")
	if err != nil {
		fmt.Println(err)
		return
	}
	memberID, err := lc.getIntInput("Enter member ID: ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = lc.Service.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error borrowing book:", err)
		return
	}
	fmt.Println("Book borrowed successfully!")
}

func (lc *LibraryController) returnBook() {
	bookID, err := lc.getIntInput("Enter book ID to return: ")
	if err != nil {
		fmt.Println(err)
		return
	}
	memberID, err := lc.getIntInput("Enter member ID: ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = lc.Service.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error returning book:", err)
		return
	}
	fmt.Println("Book returned successfully!")
}

func (lc *LibraryController) listAvailableBooks() {
	books := lc.Service.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}
	fmt.Println("--- Available Books ---")
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}
}

func (lc *LibraryController) listBorrowedBooks() {
	memberID, err := lc.getIntInput("Enter member ID: ")
	if err != nil {
		fmt.Println(err)
		return
	}
	books := lc.Service.ListBorrowedBooks(memberID)
	if len(books) == 0 {
		fmt.Println("This member has not borrowed any books, or member ID is incorrect.")
		return
	}
	fmt.Printf("--- Books Borrowed by Member ID %d ---\n", memberID)
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}
}
