package simple_test

import (
	"fmt"
	"github.com/mvl-at/rest/simple"
	"testing"
)

func TestMembers(t *testing.T) {
	simple.OpenDatabase()
	members := []simple.DBO{simple.MemberGroup{}}
	err := simple.QueryData(simple.MemberQuery, &members, 7)
	fmt.Println(err)
	fmt.Println(members)
}