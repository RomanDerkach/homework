package api

import (
	"testing"

	"github.com/RomanDerkach/homework/storage"
)

func Test_indexByID(t *testing.T) {
	type args struct {
		id    string
		books []storage.Book
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"found", args{"C97376B9-6C2E-41E5-9DBE-2E82C0EF114B", storage.GetBooksData()}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
