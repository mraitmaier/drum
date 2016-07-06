//
package main

import (
    "fmt"
    //"github.com/mraitmaier/drum"
)

const (
    path1 = "fixtures/pattern_1.splice"
    path2 = "fixtures/pattern_2.splice"
    path3 = "fixtures/pattern_3.splice"
    path4 = "fixtures/pattern_4.splice"
    path5 = "fixtures/pattern_5.splice"
)

var paths = []string{path1, path2, path3, path4, path5}

func main() {

    var p *Pattern
    var err error
    for _, f := range paths {
        if p, err = DecodeFile(f); err != nil {
            fmt.Printf(err.Error())
        }
        fmt.Println(p.String())
    }
}
