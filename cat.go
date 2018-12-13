package main

import (
	"fmt"
	"io/ioutil"
)
func main(){
	readFile("/data/go/cpu.go")

}

func readFile(name string){
    con,err:=ioutil.ReadFile(name)
    if err!=nil{
    	fmt.Println("something ERR")
	}
	res:=string(con)

	fmt.Println(res)
}
