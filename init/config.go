package init

import (
	"fmt"
	"github.com/Eros-Vertigo/srun4-config/config"
)

func init() {
	temp, err := config.GetConfig("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(temp)
}
