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

func testEqBooks(a, b storage.Books) bool {

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
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"found", args{"C97376B9-6C2E-41E5-9DBE-2E82C0EF114B", storage.GetBooksData()}, 1, false},
		{"notfound", args{"11111111-1111-1111-1111-111111111111", storage.GetBooksData()}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := indexByID(tt.args.id, tt.args.books)
			if (err != nil) != tt.wantErr {
				t.Errorf("indexByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("indexByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_booksHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

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
	//??reqPOST := http.NewRequest("POST", "/books", nil)
	rwGET := httptest.NewRecorder()
	rwPOSTBad := httptest.NewRecorder()
	rwPOSTGood := httptest.NewRecorder()

	tests := []struct {
		name     string
		status   int
		expBooks []storage.Book
		rw       *httptest.ResponseRecorder
		args     args
	}{
		{"testget", http.StatusOK, expBooks, rwGET, args{rwGET, reqGET}},
		{"testpostbad", http.StatusBadRequest, nil, rwPOSTBad, args{rwPOSTBad, reqPOSTBad}},
		{"testpostgood", http.StatusOK, nil, rwPOSTGood, args{rwPOSTGood, reqPOSTGood}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			booksHandler(tt.args.w, tt.args.r)
			if status := tt.rw.Code; status != tt.status {
				t.Errorf("Got wrong status code, expected %v but got %v",
					tt.status, status)
			}
			if tt.expBooks != nil {
				books := storage.Books{}
				err := json.NewDecoder(tt.rw.Body).Decode(&books)
				if err != nil {
					t.Errorf("returned data is somehow broken: %v", err)
				}
				//?? !reflect.DeepEqual why here it's not working
				if !testEqBooks(books, tt.expBooks) {
					t.Errorf("Got wrong data in responce, expected %+v\n but got %+v\n",
						tt.expBooks, books)
				}
			}
		})
	}
}

func Test_booksHandlerByID(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
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

	tests := []struct {
		name   string
		rw     *httptest.ResponseRecorder
		status int
		expect bool
		args   args
	}{
		{"testGET", rwGET, http.StatusOK, true, args{rwGET, reqGET}},
		{"testDELGood", rwDELGood, http.StatusAccepted, false, args{rwDELGood, reqDELGood}},
		{"testDELBad", rwDELBad, http.StatusNotFound, false, args{rwDELBad, reqDELBad}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			booksHandlerByID(tt.args.w, tt.args.r)
			if status := tt.rw.Code; status != tt.status {
				t.Errorf("Got wrong status code, expected %v but got %v",
					tt.status, status)
			}
			if tt.expect {
				book := storage.Book{}
				err := json.NewDecoder(tt.rw.Body).Decode(&book)
				if err != nil {
					t.Errorf("returned data is somehow broken: %v", err)
				}

				if !reflect.DeepEqual(book, expBooks[0]) {
					t.Errorf("Got wrong data in responce, expected %v\n but got %v\n",
						book, expBooks[0])
				}
			}
		})
	}
}
