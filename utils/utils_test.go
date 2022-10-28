package utils

import (
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUtils(t *testing.T) {
	t.Run("Contains", func(t *testing.T) {
		cases := []struct {
			array []string
			value string
			want  bool
		}{
			{[]string{"ingredients", "recipes"}, "recipes", true},
			{[]string{"ingredients", "recipes"}, "nope", false},
		}

		for _, c := range cases {
			got := Contains(c.array, c.value)
			if got != c.want {
				t.Errorf("Contains(%v, \"%s\"), got (\"%t\"), want (\"%t\")", c.array, c.value, got, c.want)
			}
		}
	})

	t.Run("ParseDatetimeFromMongoID", func(t *testing.T) {
		t1, _ := time.Parse("2006-01-02T15:04:05.000Z", "2021-11-07T14:40:54.000Z")

		cases := []struct {
			id   string
			want time.Time
			err  error
		}{
			{"6187e576abc057dac3e7d5e4", t1, nil},
			{"hello", time.Unix(0, 0).UTC(), fmt.Errorf(ErrorInvalidObjectID)},
		}

		for _, c := range cases {
			got, err := ParseDatetimeFromMongoID(c.id)
			if got != c.want {
				t.Errorf("ParseDatetimeFromMongoID(\"%s\"), got (\"%s\", \"%s\"), want (\"%s\", \"%s\")", c.id, got, err, c.want, c.err)
			}
		}
	})

	t.Run("StringSliceFromBsonM", func(t *testing.T) {
		documents := []primitive.M{
			bson.M{
				"name": "Boba Fett",
				"age":  41,
			},
			bson.M{
				"name": "Din Djarin",
				"age":  39,
			},
		}

		cases := []struct {
			documents []primitive.M
			key       string
			want      []string
		}{
			{documents, "name", []string{"Boba Fett", "Din Djarin"}},
			{documents, "stuff", []string{}},
		}

		for _, c := range cases {
			got := StringSliceFromBsonM(c.documents, c.key)
			for i := range got {
				if got[i] != c.want[i] {
					t.Errorf("StringSliceFromBsonM(\"%s\", \"%s\")[%d], got (\"%s\"), want (\"%s\")", documents, c.key, i, got[i], c.want[i])
				}
			}
		}
	})
}
