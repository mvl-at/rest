package simple_test

import (
	"fmt"
	"github.com/mvl-at/rest/simple"
	"testing"
)

func TestMembers(t *testing.T) {
	simple.OpenDatabase()
	fmt.Println(simple.Members())
}