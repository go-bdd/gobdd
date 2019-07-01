package gobdd

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func printErrorf(format string, params ...interface{}) {
	fmt.Print(Red(fmt.Sprintf(format, params...)))
}
