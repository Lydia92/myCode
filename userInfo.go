package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"io"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)
func main(){
	var dataSourceName string
	s:=readConf("E:\\go\\src\\DB.conf")
	//d:="E:\\go\\src\\info.sql"
    for _,v:=range s{
		dataSourceName=v[0]+":"+v[1]+"@tcp(" +v[2]+":"+v[3]+")"+"/"
		dataSourceName= strings.Replace(dataSourceName, " ", "", -1)

		/*str :=getDB(dataSourceName)

		writeResult(d,dataSourceName)
		for _,v:=range str{
			writeResult(d,v)
		}*/
		var engine *xorm.Engine
		showRes:=make([]string,0)
		var err error
		engine, err = xorm.NewEngine("mysql", dataSourceName)
		if err != nil {
			fmt.Println("input ocur some error")
		}
		defer engine.Close()
		res,err :=engine.Query("select user,host from mysql.user where user not like 'mysql.%'")
		if err	!=nil{
			fmt.Println("sql error")
			panic(err)
		}
		for _,v :=range res{
			for _,vv:=range v{
				showRes=append(showRes,string(vv))

				//fmt.Println(k,string(vv))
			}

		}
		//fmt.Println(showRes)
		for i,v:=range showRes{
			fmt.Println(i,v)
		}
	}
}


/*func writeResult(fileName string,str string){

	file,err:=os.OpenFile(fileName,os.O_WRONLY|os.O_APPEND|os.O_CREATE,0666)
	if err!=nil{
		panic(err)
	}
		_,err=file.WriteString(str+"\n")
}
func getDB(dataSourceName string)[]string{
	var engine *xorm.Engine
	showRes:=make([]string,0)
	var err error
	engine, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println("input ocur some error")
	}
	defer engine.Close()
	res,err :=engine.Query("show databases")
	if err	!=nil{
		fmt.Println("sql error")
		panic(err)
	}
	for _,kv:=range res{
		for _,vv:=range kv{
			showRes=append(showRes, string(vv))
			}
		}
	return showRes
	}*/

func readConf(fileName string)[][]string{
	var str []string
	var con [][]string
	file,err:=os.OpenFile(fileName,os.O_RDWR|os.O_APPEND,0666)
	if 	err!=nil{
		fmt.Println("open file err")
	}
	defer file.Close()
	buf:=bufio.NewReader(file)
	for{
		line,err:=buf.ReadString('\n')

		line=strings.TrimSpace(line)
		str=strings.SplitAfter(line, " ")
		con=append(con,str)
		if err!=nil {
			if err==io.EOF {
				break
			}else{
				fmt.Println("read file error")
				panic(err)
			}
		}

	}
	return con
}
