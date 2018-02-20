package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"bytes"

	"github.com/RomanDerkach/homework/storage"
)

func BenchmarkFilterBooks(b *testing.B) {
	fmt.Print("1111")
	// defer profile.Start(profile.CPUProfile).Stop()
	store, err := storage.NewJSONStorage("test_data/test_storage.json")
	if err != nil {
		b.Fatal("cant create storage")
	}
	// filterStr := `
	// {
	//     "title": "title4"
	// }`

	handler, err := NewHandler(store)
	filter, err := json.Marshal(storage.BookFilter{Title: "title4"})
	// filter, err := json.Marshal(filterStr)

	// bookFilter := storage.BookFilter{}
	// err = json.NewDecoder(strings.NewReader(filterStr)).Decode(&bookFilter)
	// if err != nil {
	// 	fmt.Print("bli")
	// 	b.Fatal("s")
	// }

	// fmt.Print(bookFilter)
	// return

	if err != nil {
		b.Fatal("cant marshal filter")
	}
	reqw := httptest.NewRequest("POST", "/books/helpmewithbooks", bytes.NewReader(filter))

	if err != nil {
		b.Fatal("cant create request")
	}
	for i := 0; i < 2; i++ {
		rw := httptest.NewRecorder()
		handler.booksHandlerByIDPOST(rw, reqw)
	}

}

// func testEqBooks(t *testing.T, a, b storage.Books) bool {
// 	t.Helper()
// 	if a == nil && b == nil {
// 		return true
// 	}

// 	if a == nil || b == nil {
// 		return false
// 	}

// 	if len(a) != len(b) {
// 		return false
// 	}

// 	if &a == &b {
// 		return true
// 	}

// 	return true
// }

// func Test_indexByID(t *testing.T) {
// 	type args struct {
// 		id    string
// 		books storage.Books
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		want        int
// 		expectedErr error
// 	}{
// 		{"found", args{"C97376B9-6C2E-41E5-9DBE-2E82C0EF114B", storage.GetBooks()}, 1, nil},
// 		{"notfound", args{"11111111-1111-1111-1111-111111111111", storage.GetBooks()}, 0, ErrNotFound},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			got, err := indexByID(tt.args.id, tt.args.books)
// 			if err != tt.expectedErr {
// 				t.Errorf("indexByID() error = %v, wantErr %v", err, tt.expectedErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("indexByID() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_booksHandler(t *testing.T) {
// 	expBooks := storage.GetBooks()
// 	newBook := storage.Book{
// 		Title:  "test",
// 		Genres: []string{"erotic"},
// 		Pages:  111,
// 		Price:  33.33,
// 	}
// 	book, err := json.Marshal(newBook)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	type args struct {
// 		r *http.Request
// 	}
// 	tests := []struct {
// 		name     string
// 		status   int
// 		expBooks []storage.Book
// 		args     args
// 	}{
// 		{"testget", http.StatusOK, expBooks, args{httptest.NewRequest("GET", "/books", nil)}},
// 		{"testpostbad", http.StatusBadRequest, nil, args{httptest.NewRequest("POST", "/books", nil)}},
// 		{"testpostgood", http.StatusOK, nil, args{httptest.NewRequest("POST", "/books", bytes.NewReader(book))}},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rr := httptest.NewRecorder()
// 			booksHandler(rr, tt.args.r)
// 			if status := rr.Code; status != tt.status {
// 				t.Errorf("Got wrong status code, expected %v but got %v",
// 					tt.status, status)
// 			}

// 			if tt.expBooks == nil {
// 				return
// 			}

// 			var books []storage.Book
// 			err := json.NewDecoder(rr.Body).Decode(&books)
// 			if err != nil {
// 				t.Errorf("returned data is somehow broken: %v", err)
// 				return
// 			}

// 			// It was not working because tt.expBooks was of type []Book and books of type Books
// 			if !reflect.DeepEqual(books, tt.expBooks) {
// 				t.Errorf("Got wrong data in responce, expected \n%+v\n but got \n%+v\n",
// 					tt.expBooks, books)
// 			}
// 		})
// 	}
// }

// func Test_booksHandlerByID(t *testing.T) {
// 	expBooks := storage.GetBooks()
// 	if len(expBooks) == 0 {
// 		t.Fatal("There is no books in tested storage")
// 	}
// 	bookurl := "/books/" + expBooks[0].ID

// 	reqGET := httptest.NewRequest("GET", bookurl, nil)

// 	reqDELGood := httptest.NewRequest("DELETE", bookurl, nil)

// 	reqDELBad := httptest.NewRequest("DELETE", "/books/badid", nil)

// 	type args struct {
// 		r *http.Request
// 	}
// 	tests := []struct {
// 		name   string
// 		status int
// 		expect *storage.Book
// 		args   args
// 	}{
// 		{
// 			"testGET",
// 			http.StatusOK,
// 			nil,
// 			args{
// 				reqGET,
// 			},
// 		},
// 		{"testDELGood", http.StatusAccepted, nil, args{reqDELGood}},
// 		{"testDELBad", http.StatusNotFound, nil, args{reqDELBad}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rr := httptest.ResponseRecorder{}
// 			booksHandlerByID(&rr, tt.args.r)
// 			if status := rr.Code; status != tt.status {
// 				t.Errorf("Got wrong status code, expected %v but got %v",
// 					tt.status, status)
// 			}

// 			if tt.expect == nil {

// 			}

// 			book := storage.Book{}
// 			err := json.NewDecoder(rr.Body).Decode(&book)
// 			if err != nil {
// 				t.Errorf("returned data is somehow broken: %v", err)
// 				return
// 			}

// 			if !reflect.DeepEqual(book, expBooks[0]) {
// 				t.Errorf("Got wrong data in responce, expected %v\n but got %v\n",
// 					book, expBooks[0])
// 			}

// 		})
// 	}
// }
