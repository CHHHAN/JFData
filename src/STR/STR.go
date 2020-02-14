package STR

import (
	"math/rand"
	"time"
	"strings"
)

func Cut(str string, s string) []string {
	return strings.Split(str,s)
}

func HitCode(number int) string {

	numberArray := [62]string{"1","2","3","4","5","6","7","8","9","0","q","a","z","w","s","x","e","d","c","r","f","v","t","g","b","y","h","n","u","j","m","i","k","o","l","p","Q","A","Z","W","S","X","E","D","C","R","F","V","T","G","B","Y","H","N","U","J","M","I","K","O","L","P"}
	rand.Seed(time.Now().Unix())

	var autoString string
	var i int

	for i = 0; i < number; i++ {

		autoString = autoString + numberArray[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(61)]
	}

	return autoString

}

func Is(str string,is string) bool {
	return strings.Contains(str, is)
}