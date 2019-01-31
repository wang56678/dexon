package tools

import (
	"encoding/binary"
	"encoding/hex"
	"regexp"
)

func PatchBinary(input []byte) []byte {
	if input == nil {
		return nil
	}
	stringCode := hex.EncodeToString(input)
	stringCode = Patch(stringCode)
	result, _ := hex.DecodeString(stringCode)
	return result
}
func patchPattern1(input string) string {
	// Case 1
	// PUSH1 v1 DUP1 PUSH1 v2 PUSH1 v3 CODECOPY
	// we have to set v1 = v1 + 1
	r, _ := regexp.Compile("60..8060..60..39")
	loc := r.FindStringIndex(input)
	if len(loc) > 0 {
		valLoc := loc[0]/2 + 1
		bString, _ := hex.DecodeString(input)
		bString[valLoc]++
		input = hex.EncodeToString(bString)
	}
	return input
}
func patchPattern2(input string) string {
	// Case 2
	// PUSH2 v1 v2 DUP1 PUSH2 v3 PUSH1 v4 CODECOPY
	// we have to set BigEndian(v1,v2)++
	r, _ := regexp.Compile("61....8061....60..39")
	loc := r.FindStringIndex(input)
	if len(loc) > 0 {
		valLoc := loc[0]/2 + 1
		bString, _ := hex.DecodeString(input)
		val := binary.BigEndian.Uint16(bString[valLoc : valLoc+2])
		val++
		binary.BigEndian.PutUint16(bString[valLoc:], val)
		input = hex.EncodeToString(bString)
	}
	return input
}
func patchPattern3(input string) string {
	// Case 3
	// PUSH1 v1 DUP1 PUSH3 v2 v3 v4 DUP4 CODECOPY
	// we have to set BigEndian(v2,v3,v4)++
	r, _ := regexp.Compile("60..8062......8339")
	loc := r.FindStringIndex(input)
	if len(loc) > 0 {
		valLoc := loc[0]/2 + 4
		bString, _ := hex.DecodeString(input)
		toConvert := append([]byte{0}, bString[valLoc:valLoc+3]...)
		tmpBinary := make([]byte, 4)
		val := binary.BigEndian.Uint32(toConvert)
		val--
		binary.BigEndian.PutUint32(tmpBinary, val)
		bString[valLoc] = tmpBinary[1]
		bString[valLoc+1] = tmpBinary[2]
		bString[valLoc+2] = tmpBinary[3]
		input = hex.EncodeToString(bString)
	}
	return input
}
func patchPattern4(input string) string {
	// Case 4
	// PUSH2 v1 v2 DUP1 PUSH2 v3 DUP4 CODECOPY
	// we have to set BigEndian(v1,v2)+2
	r, _ := regexp.Compile("61....8061....8339")
	loc := r.FindStringIndex(input)
	if len(loc) > 0 {
		valLoc := loc[0]/2 + 1
		bString, _ := hex.DecodeString(input)
		val := binary.BigEndian.Uint16(bString[valLoc : valLoc+2])
		val = val + 2
		binary.BigEndian.PutUint16(bString[valLoc:], val)
		input = hex.EncodeToString(bString)
	}
	return input
}
func patchPattern5(input string) string {
	// Case 5
	// PUSH1 v1 DUP1 PUSH2 v2 v3 DUP4 CODECOPY
	// we have to set BigEndian(v2,v3)++
	r, _ := regexp.Compile("60..8061....8339")
	loc := r.FindStringIndex(input)
	if len(loc) > 0 {
		valLoc := loc[0]/2 + 4
		bString, _ := hex.DecodeString(input)
		val := binary.BigEndian.Uint16(bString[valLoc : valLoc+2])
		val = val + 1
		binary.BigEndian.PutUint16(bString[valLoc:], val)
		input = hex.EncodeToString(bString)
	}
	return input
}
func patchPattern6(input string) string {
	// Case 6
	// SUB DUP1 PUSH2 v2 v3 DUP4 CODECOPY
	// we have to set BigEndian(v2,v3)++
	r, _ := regexp.Compile("038061....8339")
	loc := r.FindStringIndex(input)
	if len(loc) > 0 {
		valLoc := loc[0]/2 + 3
		bString, _ := hex.DecodeString(input)
		val := binary.BigEndian.Uint16(bString[valLoc : valLoc+2])
		val = val + 1
		binary.BigEndian.PutUint16(bString[valLoc:], val)
		input = hex.EncodeToString(bString)
	}
	return input
}
func addPrefix(input string) string {
	r, _ := regexp.Compile("6060604052")
	loc := r.FindAllStringIndex(input, -1)
	if len(loc) > 0 {
		for i, offset := range loc {
			insertLoc := offset[0] + i*2
			if !isDelegateCall(input, insertLoc) {
				input = input[:insertLoc] + "00" + input[insertLoc:]
			}
		}
		if input[0:2] != "00" {
			input = "00" + input
		}
	}
	return input
}
func isDelegateCall(input string, codeStartAt int) bool {
	delegateCallPrefix := "6504032353da7150"
	prefixLen := len(delegateCallPrefix)
	newStartAt := codeStartAt - prefixLen
	if newStartAt >= 0 {
		if input[newStartAt:newStartAt+prefixLen] == delegateCallPrefix {
			return true
		}
	}
	return false
}

// Patch patch lagacy bytecode(hex string) to new vm interface compatible bytecode
func Patch(input string) string {
	if len(input) == 0 {
		return input
	}
	input = patchPattern1(input)
	input = patchPattern2(input)
	input = patchPattern3(input)
	input = patchPattern4(input)
	input = patchPattern5(input)
	input = patchPattern6(input)
	result := addPrefix(input)
	return result
}
