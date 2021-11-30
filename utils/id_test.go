package utils_test

import (
	"fmt"
	. "github.com/starship-cloud/starship-iac/testing"
	"github.com/starship-cloud/starship-iac/utils"
	"strings"
	"testing"
)

func TestIdGen(t *testing.T) {
	for i:=0; i < 10; i++ {
		t.Run("test-" + fmt.Sprint(i), func(t *testing.T){
			strUserId :=  utils.GenUserId()
			fmt.Println("user id: " + strUserId)
			exp := "user-"
			Assert(t, strings.Contains(strUserId, exp), "exp %q to be contained in %q", exp, strUserId)
		})
	}



}
