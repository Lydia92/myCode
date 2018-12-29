package main

import (
	"bufio"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"io"
	"os"
	"strings"
)

func main() {
	var dataSourceName string

	//var grants string
	s := readConf("/home/liyingdi/db.conf")
	d := "/home/liyingdi/info.sql"
	dd := "/home/liyingdi/DB.sql"
	for _, v := range s {
		fmt.Println(v)
		dataSourceName = v[0] + ":" + v[1] + "@tcp(" + v[2] + ":" + v[3] + ")" + "/"
		dataSourceName = strings.Replace(dataSourceName, " ", "", -1)
		writeResult(d, dataSourceName)
		writeResult(dd, dataSourceName)
		str := getDB(dataSourceName)

		for _, v := range str {

			writeResult(dd, v)
		}
		writeResult(dd, "\n")
		var engine *xorm.Engine
		var err error
		engine, err = xorm.NewEngine("mysql", dataSourceName)
		if err != nil {
			fmt.Println("input ocur some error")
		}
		defer engine.Close()
		res, err := engine.Query("select user,host from mysql.user where user not like 'mysql.%' and user!='' and host !='::1' and host not like '%localdomain%'")
		if err != nil {
			fmt.Println("sql error")
			panic(err)
		}
		for _, v1 := range res {
			query := fmt.Sprintf("show grants for %s@'%s';", v1["user"], v1["host"])
			res1, err := engine.Query(query)
			if err != nil {
				fmt.Println(err)
			}
			for _, v2 := range res1 {
				//fmt.Println(v2)
				for _, vv := range v2 {
					if split(string(vv)) != "" {
						str := fmt.Sprintf("%s\t%s\t%s", v[2], v[3], split(string(vv)))

						writeResult(d, str)
					}

				}
			}

		}
		writeResult(d, "\n")

	}
}
func split(str string) string {
	var grants string
	var database string
	var table string
	var user string
	var hosts string
	var option string
	var results string
	if !strings.Contains(str, "PROXY ON") {
		//fmt.Println(str)

		srr := strings.Split(str, " ON ")[0]

		grants = strings.Split(srr, "GRANT")[1]
		//fmt.Println(grants)
		srr1 := strings.Split(str, " ON ")[1]
		database = strings.Split(srr1, "TO")[0]
		table = strings.Split(database, ".")[0]
		database = strings.Split(database, ".")[0]

		str1 := strings.Split(srr1, "TO")[1]
		user = strings.Split(str1, "@")[0]
		hosts = strings.Split(str1, "@")[1]
		hosts = strings.Split(hosts, "' ")[0]
		hosts = strings.Split(hosts, "'")[1]

		//fmt.Println(len(strings.Split(hosts,"WITH")))
		if strings.HasSuffix(str, "GRANT OPTION") {
			option = "1"
		} else {
			option = "0"
		}
		results = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", grants, database, table, user, hosts, option)
		//	fmt.Println(results)
	}
	return results
}

func writeResult(fileName string, str string) {
	//fileName:="E:\\go\\src\\info.sql"
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(str + "\n")
	//fmt.Println()
}
func getDB(dataSourceName string) []string {
	var engine *xorm.Engine
	showRes := make([]string, 0)
	var err error
	engine, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println("input ocur some error")
	}
	defer engine.Close()
	res, err := engine.Query("show databases")
	if err != nil {
		fmt.Println("sql error")
		panic(err)
	}
	for _, kv := range res {
		for _, vv := range kv {
			showRes = append(showRes, string(vv))
		}
	}
	return showRes
}

func readConf(fileName string) [][]string {
	var str []string
	var con [][]string
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("open file err")
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')

		line = strings.TrimSpace(line)
		str = strings.SplitAfter(line, " ")
		con = append(con, str)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("read file error")
				panic(err)
			}
		}

	}
	return con
}
