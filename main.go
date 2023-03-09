// For didactic purposes only.
// The idea is to prove that an interface value wraps a copy of the concrete value it has been assigned to it by comparing the addresses of the various values involved.
// It uses the unsafe package. Again, for didactic purposes only.
//
// In this program we create a value of type MyType and store it in the variable myTypeValue. The value has a field, Name, which is of type string.
// We define also a variable aVar of type interface{} and we assign the value stored in myTypeValue to aVar.
// The assignement operation creates a copy of the value of type MyType and it is the copy that is wrapped within the interface{} value.
// To prove that it is really a copy, we compare the addresses of the value of MyType and the pointer to the value wrpaeed within the interface{} value
// and we see that they are different.
// To extract the pointer to the value wrapped by the interface we use the unsafe package and the possibility it gives us to convert an unsafe.Pointer to a
// pointer of any type.
//
// We compare also the addresses of the fields Name extracted from both the value of MyType and the interface{} value
// and we see that they are different.
// The 2 different string values share the same backing byte array, and we prove this comparing their addresses.
// We then assign a new string value to the field Name of the value of type MyType stored in the variable myTypeValue.
// After this new assignement, the 2 string values of the fields Name have different backing byte arrays.
//
// To be able to handle the addresses within the interface{} value and within the string values we use the unsafe package.
// Again this example is only for didactical purposes.
// Inspired by https://stackoverflow.com/a/57698257/5699993
//
// To try this program run the following command
// go run .

package main

import (
	"fmt"
	"unsafe"
)

type MyType struct{ Name string }

// eface is a type struct that has the same memory layout as any interface value
type eface struct {
	typ, val unsafe.Pointer
}

// toEface receives a value of interface{} type and returns a pointer to a value of type eface which overlays the same memory space as the
// interface value passed in as argument.
func toEface(arg interface{}) *eface {
	return (*eface)(unsafe.Pointer(&arg))
}

// ifaceConcreteValPtr receives a value of interface{} type and returns a pointer to the concrete value wrapped by the interface{} value.
func ifaceConcreteValPtr(arg interface{}) unsafe.Pointer {
	ptrToEfaceVal := toEface(arg)
	return ptrToEfaceVal.val
}

// toMyType takes an unsafe pointer and convert it to a pointer of type *MyType and then dereferences such pointer.
// The net result is that toMyType returns a value of type MyType whose pointer is the value passed in as argument.
func toMyType(unsPtr unsafe.Pointer) MyType {
	return *(*MyType)(unsPtr)
}

// estring is a type struct that has the same memory layout as any string value
// it is used to extract the pointer to the byte array which is the backing array of a string
type estring struct {
	ptrToByteArray unsafe.Pointer
	length         int
}

// toPtrToByteArray returns the pointer to the backing byte array of the string passed in as argument
func toPtrToByteArray(s string) unsafe.Pointer {
	return (*estring)(unsafe.Pointer(&s)).ptrToByteArray
}

// toLength returns the int value representing the length of the string passed in as argument
func toLength(s string) int {
	return (*estring)(unsafe.Pointer(&s)).length
}

func main() {
	var aVar any

	type MyType struct {
		Name string
	}

	myTypeValue := MyType{Name: "MyName"}

	// the assignement of myTypeValue to aVar creates a copy of the value referenced by myTypeValue. This copy is what is wrapped within the interface{} value.
	aVar = myTypeValue

	aVarAsMyType := toMyType(ifaceConcreteValPtr(aVar))

	fmt.Printf("value %v - type %T \n", aVarAsMyType, aVarAsMyType)
	fmt.Printf("address of myTypeValue: %p \n", &myTypeValue)
	fmt.Printf("address of the concrete value held by the interface type variable aVar: %p \n", &aVarAsMyType)
	// the equality check &myTypeValue == &aVarAsMyType raises a compiler error "invalid operation: &myTypeValue == &aVarAsMyType (mismatched types *MyType and *MyType)"
	// for this reason the check is done converting the addressed to uintptr
	fmt.Printf("Are the 2 fields Name the same string value? %v \n", uintptr(unsafe.Pointer(&myTypeValue)) == uintptr(unsafe.Pointer(&aVarAsMyType)))

	fmt.Print("\n\n")

	fmt.Printf("address of Name field of myTypeValue: %p \n", &myTypeValue.Name)
	fmt.Printf("address of Name field of the concrete value held by the interface type variable aVar: %p \n", &aVarAsMyType.Name)
	fmt.Printf("Are the 2 fields Name the same string value? %v \n", &myTypeValue.Name == &aVarAsMyType.Name)

	fmt.Print("\n\n")

	fmt.Println("Address of the byte array backing the string value held by the Name field of myTypeValue", toPtrToByteArray(myTypeValue.Name))
	fmt.Println("Lenght of the byte array of Name field of myTypeValue", toLength(myTypeValue.Name))
	fmt.Println("Address of the byte array backing the string value held by the Name field of aVarAsMyType", toPtrToByteArray(aVarAsMyType.Name))
	fmt.Println("Lenght of the byte array of Name field of aVarAsMyType", toLength(aVarAsMyType.Name))
	fmt.Printf("Do the 2 fields Name share the same backing byte array? %v \n", toPtrToByteArray(myTypeValue.Name) == toPtrToByteArray(aVarAsMyType.Name))

	fmt.Print("\n\n")

	fmt.Println("Now we assign a new string value to the field Name of the value of type MyType referenced by the variable myTypeValue")
	myTypeValue.Name = "Another name"
	fmt.Println("Address of Name field of myTypeValue after change", toPtrToByteArray(myTypeValue.Name))
	fmt.Println("Lenght of Name field of myTypeValue after change", toLength(myTypeValue.Name))
	fmt.Printf("Do the 2 fields Name share the same backing byte array after the change? %v \n", toPtrToByteArray(myTypeValue.Name) == toPtrToByteArray(aVarAsMyType.Name))
}
