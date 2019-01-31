package tools

import (
	"fmt"
)

func ExamplePatchPUSH1() {
	input := "606060405260068060106000396000f3606060405200"
	//                    07                  ^00
	fmt.Println(Patch(input))
	// Output:
	// 00606060405260078060106000396000f300606060405200
}
func ExamplePatchPUSH2() {
	input := "6060604052aabb616fff8606060405200"
	//                        7000 ^00
	fmt.Println(Patch(input))
	// Output:
	// 006060604052aabb616fff800606060405200
}
func ExamplePatchDelegate() {
	input := "6504032353da71506060604052"
	fmt.Println(Patch(input))
	// Output:
	// 006504032353da71506060604052
}
