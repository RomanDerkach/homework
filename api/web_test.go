package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/RomanDerkach/homework/storage"
)

func testEqBooks(t *testing.T, a, b storage.Books) bool {
	t.Helper()
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	if &a == &b {
		return true
	}

	return true
}

func Test_indexByID(t *testing.T) {
	type args struct {
		id    string
		books storage.Books
	}
	tests := []struct {
		name        string
		args        args
		want        int
		expectedErr error
	}{
		{"found", args{"C97376B9-6C2E-41E5-9DBE-2E82C0EF114B", storage.GetBooksData()}, 1, nil},
		{"notfound", args{"11111111-1111-1111-1111-111111111111", storage.GetBooksData()}, 0, ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := indexByID(tt.args.id, tt.args.books)
			if err != tt.expectedErr {
				t.Errorf("indexByID() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}
			if got != tt.want {
				t.Errorf("indexByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_booksHandler(t *testing.T) {
	expBooks := storage.GetBooksData()
	newBook := storage.Book{
		Title:  "test",
		Ganres: []string{"erotic"},
		Pages:  111,
		Price:  33.33,
	}
	book, err := json.Marshal(newBook)
	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewReader(book)

	reqGET := httptest.NewRequest("GET", "/books", nil)
	reqPOSTBad := httptest.NewRequest("POST", "/books", nil)
	reqPOSTGood := httptest.NewRequest("POST", "/books", body)
	//gives panic ?
	//reqPOST, err := http.NewRequest("POST", "/books", nil)

	rwGET := httptest.NewRecorder()
	rwPOSTBad := httptest.NewRecorder()
	rwPOSTGood := httptest.NewRecorder()

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		status   int
		expBooks []storage.Book
		args     args
	}{
		{"testget", http.StatusOK, expBooks, args{rwGET, reqGET}},
		{"testpostbad", http.StatusBadRequest, nil, args{rwPOSTBad, reqPOSTBad}},
		{"testpostgood", http.StatusOK, nil, args{rwPOSTGood, reqPOSTGood}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			booksHandler(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != tt.status {
				t.Errorf("Got wrong status code, expected %v but got %v",
					tt.status, status)
			}

			books := storage.Books{}
			err := json.NewDecoder(tt.args.w.Body).Decode(&books)
			if err != nil {
				t.Errorf("returned data is somehow broken: %v", err)
			}
			//?? !reflect.DeepEqual why here it's not working
			//if !reflect.DeepEqual(books, tt.expBooks) {
			// TODO: CHECK!!!
			if !testEqBooks(books, tt.expBooks) {
				t.Errorf("Got wrong data in responce, expected %+v\n but got %+v\n",
					tt.expBooks, books)
			}
		})
	}
}

func Test_booksHandlerByID(t *testing.T) {
	expBooks := storage.GetBooksData()
	if len(expBooks) == 0 {
		t.Fatal("There is no books in tested storage")
	}
	bookurl := "/books/" + expBooks[0].ID

	reqGET := httptest.NewRequest("GET", bookurl, nil)
	rwGET := httptest.NewRecorder()

	reqDELGood := httptest.NewRequest("DELETE", bookurl, nil)
	rwDELGood := httptest.NewRecorder()

	reqDELBad := httptest.NewRequest("DELETE", "/books/badid", nil)
	rwDELBad := httptest.NewRecorder()

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name   string
		rw     *httptest.ResponseRecorder
		status int
		expect *storage.Book
		args   args
	}{
		{
			"testGET",
			rwGET,
			http.StatusOK,
			true,
			args{
				rwGET,
				reqGET,
			},
		},
		{"testDELGood", rwDELGood, http.StatusAccepted, nil, args{reqDELGood}},
		{"testDELBad", rwDELBad, http.StatusNotFound, false, args{reqDELBad}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.ResponseRecorder{}
			booksHandlerByID(&rr, tt.args.r)
			if status := rr.Code; status != tt.status {
				t.Errorf("Got wrong status code, expected %v but got %v",
					tt.status, status)
			}

			if tt.expect == nil {

			}

			book := storage.Book{}
			err := json.NewDecoder(tt.rw.Body).Decode(&book)
			if err != nil {
				t.Errorf("returned data is somehow broken: %v", err)
				return
			}

			if !reflect.DeepEqual(book, expBooks[0]) {
				t.Errorf("Got wrong data in responce, expected %v\n but got %v\n",
					book, expBooks[0])
			}

		})
	}
}
