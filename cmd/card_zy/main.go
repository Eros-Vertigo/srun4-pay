package main

import (
	_ "srun4-pay/init/common"
	"srun4-pay/init/init_card_zy"
)

// 一卡通-正元
func main() {
	var Zy init_card_zy.ZyConfig
	Zy.LoadYaml()
	Zy.Routers()
}
