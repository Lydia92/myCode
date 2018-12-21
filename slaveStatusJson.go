package main
import(
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	//"strings"
	"encoding/json"
	"strconv"
)
func main() {
	var engine *xorm.Engine
	var err error
	engine, err = xorm.NewEngine("mysql", "root:123456@tcp(192.168.160.133:3306)/aa")
	if err != nil {
		fmt.Println("input ocur some error")
	} else {
		fmt.Println(" ok")
	}
	defer engine.Close()

	/*	err = engine.Ping()
		if err != nil {
			fmt.Println(err)
			return
		}*/
	//res,err :=engine.Query("select * from tt where id=1")
	res,err :=engine.Query("show slave status")
	if err	!=nil{
		fmt.Println("sql error")
		return
	}
	ss:=slaveStatus(res)
	fmt.Println(ss)
	/*s,err:=json.Marshal(ss)
	if err!=nil {
		fmt.Println("error")
	}
	fmt.Println(string(s))*/

}
func slaveStatus(res []map[string][]byte)string{
	result:=make(map[string]string)
	rr:=make(map[string]string)
	//re:=make(map[string]interface{})
	for _,v :=range res {
		//fmt.Println(k)
		//fmt.Println(v)

		for k,vv:=range v{
			//result=vv
			result=map[string]string{k:string(vv)}
			switch k {
			case "Master_User":
				rr["master_user"]=result["Master_User"]

			case "Master_Host":
				rr["Master_Host"]=result["Master_Host"]
			case "Master_Port":
				rr["Master_Port"]=result["Master_Port"]
			case "Connect_Retry":
				rr["Connect_Retry"]=result["Connect_Retry"]
			case "Master_Log_File":
				rr["Master_Log_File"]=result["Master_Log_File"]
			case "Read_Master_Log_Pos":
				rr["Read_Master_Log_Pos"]=result["Read_Master_Log_Pos"]
			case "Relay_Master_Log_File":
				rr["Relay_Master_Log_File"]=result["Relay_Master_Log_File"]
			case "Slave_IO_Running":
				rr["Slave_IO_Running"]=result["Slave_IO_Running"]
			case "Slave_SQL_Running":
				rr["Slave_SQL_Running"]=result["Slave_SQL_Running"]
			case "Exec_Master_Log_Pos":
				rr["Exec_Master_Log_Pos"]=result["Exec_Master_Log_Pos"]
			case "Relay_Log_Space":
				rr["Relay_Log_Space"]=result["Relay_Log_Space"]
			case "Master_UUID":
				rr["Master_UUID"]=result["Master_UUID"]
			case "Seconds_Behind_Master":
				rr["Seconds_Behind_Master"]=result["Seconds_Behind_Master"]
			default:
				continue
			}
			//fmt.Println(result)
			//fmt.Println("result['Master_User']:",result["Master_User"])
			//result={ k:string(vv }
			//fmt.Printf("%s:",k)
			//fmt.Printf(" %s ",string(vv))
		}

	}
	if rr["Master_Log_File"]==rr["Relay_Master_Log_File"]{

		//fmt.Println("master:", rr["Master_Log_File"])

		if  result["Exec_Master_Log_Pos"]==result["Read_Master_Log_Pos"]{

			rr["Seconds_Behind_M"]="0"
		}else{
			//ss:=strconv.Itoa(rr["Read_Master_Log_Pos"])-strconv.Itoa(rr["Exec_Master_Log_Pos"])
			Read_Master_Log_Pos,err :=strconv.Atoi(rr["Read_Master_Log_Pos"])
			Exec_Master_Log_Pos,err :=strconv.Atoi(rr["Exec_Master_Log_Pos"])
			if err!=nil{
				fmt.Println("fail")
			}

			Seconds_Behind_M:=strconv.Itoa(Read_Master_Log_Pos-Exec_Master_Log_Pos)
			rr["Seconds_Behind_M"]=Seconds_Behind_M

		}
		//fmt.Println(Master_Log_File)

	}else {
		Master_Log_File,err :=strconv.Atoi(rr["Master_Log_File"][4:len(rr["Master_Log_File"])])
		Relay_Master_Log_File,err :=strconv.Atoi(rr["Relay_Master_Log_File"][4:len(rr["Relay_Master_Log_File"])])
		if err!=nil	{
			fmt.Println("fail")

		}
		Seconds_Behind_M:=strconv.Itoa(Master_Log_File-Relay_Master_Log_File)

		rr["Seconds_Behind_M"]=Seconds_Behind_M
	}
	//re=rr
	s,err:=json.Marshal(rr)
	if err!=nil{
		fmt.Println("error")
	}
	ss:=string(s)
	return ss
}
