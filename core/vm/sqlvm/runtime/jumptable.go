package runtime

var jumpTable = [256]OpFunction{
	// 0x60
	LOAD: opLoad,
}
