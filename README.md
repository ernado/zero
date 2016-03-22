# cydev/zero
Check if golang struct is empty

[![Build Status](https://travis-ci.org/cydev/zero.svg?branch=master)](https://travis-ci.org/cydev/zero)
[![Build Status](https://travis-ci.org/cydev/zero.svg?branch=master)](https://travis-ci.org/cydev/zero)

``` go
package main

import (
        "fmt"

        "github.com/cydev/zero"
)

type Structure struct {
        ID int
}

func ExampleStructure() {
        zeroStructure := Structure{}
        zeroStructurePointer := &zeroStructure
        nonZero := Structure{ID: 1}
        nonZeroPointer := &nonZero
        fmt.Println(zero.IsZero(zeroStructure))        // true
        fmt.Println(zero.IsZero(zeroStructurePointer)) // true
        fmt.Println(zero.IsZero(nonZero))              // false
        fmt.Println(zero.IsZero(nonZeroPointer))       // false
        // Output:
        // true
        // true
        // false
        // false
}

func main() {
        ExampleStructure()
}
```
