package services

import (
	"errors";
	"library_managment/models"
)

// LibraryManager interface defines the methods for managing the library.
type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

// Library implements the LibraryManager interface.
type Library struct {
	Name string
	Books   map[int]models.Book
	Members map[int]models.Member
	nextBookID   int
	nextMemberID int
}

// NewLibrary creates and returns a new Library instance.
func NewLibrary() *Library {
	return &Library{
		Name:         "My Library",
		Books:        make(map[int]models.Book),
		Members:      make(map[int]models.Member),
		nextBookID:   1,
		nextMemberID: 1,
	}
}

// AddBook adds a new book to the library.
func (l *Library) AddBook(book models.Book) {
	book.ID = l.nextBookID
	book.Status = "Available"
	l.Books[book.ID] = book
	l.nextBookID++
}

// RemoveBook removes a book from the library by its ID.
func (l *Library) RemoveBook(bookID int) error {
	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}
	if book.Status == "Borrowed" {
		return errors.New("cannot remove a borrowed book, please return it first")
	}
	delete(l.Books, bookID)
	return nil
}

// BorrowBook allows a member to borrow a book if it is available.
func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, bookExists := l.Books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}
	if book.Status == "Borrowed" {
		return errors.New("book is already borrowed")
	}

	member, memberExists := l.Members[memberID]
	if !memberExists {
		return errors.New("member not found")
	}

	book.Status = "Borrowed"
	l.Books[bookID] = book // Update the book status in the library's book map

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member // Update the member's borrowed books

	return nil
}

// ReturnBook allows a member to return a borrowed book.
func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, bookExists := l.Books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}
	if book.Status == "Available" {
		return errors.New("book is already available, cannot return")
	}

	member, memberExists := l.Members[memberID]
	if !memberExists {
		return errors.New("member not found")
	}

	foundAndRemoved := false
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			foundAndRemoved = true
			break
		}
	}

	if !foundAndRemoved {
		return errors.New("member did not borrow this book")
	}

	book.Status = "Available"
	l.Books[bookID] = book // Update the book status in the library's book map
	l.Members[memberID] = member // Update the member's borrowed books

	return nil
}

// ListAvailableBooks lists all available books in the library.
func (l *Library) ListAvailableBooks() []models.Book {
	var availableBooks []models.Book
	for _, book := range l.Books {
		if book.Status == "Available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

// ListBorrowedBooks lists all books borrowed by a specific member.
func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, exists := l.Members[memberID]
	if !exists {
		return nil // Or return an error, depending on desired behavior
	}
	return member.BorrowedBooks
}

// AddMember is a helper function to add a member (not part of the interface, but useful for testing/setup)
func (l *Library) AddMember(member models.Member) {
	member.ID = l.nextMemberID
	l.Members[member.ID] = member
	l.nextMemberID++
}