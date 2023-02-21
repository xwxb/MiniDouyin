package favor

import (
	"log"
	"testing"
)

func TestUnFav(t *testing.T) {
	res, err := UnFav(1, 36)
	if err != nil {
		log.Println(err)
	}

	log.Println(res)
}
