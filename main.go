package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/go-xorm/xorm"
)

//var WG sync.WaitGroup

type connectInfo struct {
	Username      string
	Pass          string
	Host          string
	Port          int
	SourceConn    *xorm.Engine
	DeractionConn *xorm.Engine
}
type Userinfo struct {
	Id         int8
	Node_addr  string
	Node_port  int
	Ownership  string
	Account    string
	Passwd     string
	Priv       string
	Db_name    string
	Table_name string
	Login_addr string
	Is_grant   int
	Validity   int8
	Role       int
	Md5        string
	CreateTime time.Time `xorm:"created"`
	UpdateTime time.Time `xorm:"updated"`
}
type MysqlMetadataTables struct {
	Node_addr       string
	Node_port       int
	Table_schema    string
	Table_name      string
	Db_engine       string
	Row_format      string
	Table_rows      int
	Avg_row_length  int
	Max_data_length int
	Data_length     int
	Index_length    int
	Data_free       int
	Chip_size       int
	Auto_increment  int
	Table_collation string
	Create_time     time.Time
	Update_time     time.Time
	Check_time      time.Time
	Table_comment   string
	Table_md5       string
}
type MysqlMetadataColumns struct {
	Node_addr      string
	Node_port      int
	Table_schema   string
	Table_name     string
	Column_name    string
	Column_type    string
	Collation_name string
	Is_nullable    string
	Column_key     string
	Column_default string
	Extra          string
	Col_privileges string
	Column_comment string
	Column_md5     string
}
type MysqlMetadataIndexs struct {
	Node_addr     string
	Node_port     int
	Table_schema  string
	Table_name    string
	Column_name   string
	Non_unique    int
	Index_name    string
	Seq_in_index  int
	Cardinality   int
	Nullable      string
	Index_type    string
	Index_comment string
	Index_md5     string
}
func NewConnectInfo(Susername, Spass, Shost string, Sport int, Sdatabase string,
	Dusername, Dpass, Dhost string, Dport int, Ddatabase string) *connectInfo {
	//初始化结构体
	con := new(connectInfo)
	//拼接连接字符串
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", Susername, Spass, Shost, Sport, Sdatabase)
	dataSourceName = strings.Replace(dataSourceName, " ", "", -1)
	dataSourceName1 := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", Dusername, Dpass, Dhost, Dport, Ddatabase)
	dataSourceName1 = strings.Replace(dataSourceName1, " ", "", -1)
	var err error
	//连接数据库
	con.SourceConn, err = xorm.NewEngine("mysql", dataSourceName)
	con.DeractionConn, err = xorm.NewEngine("mysql", dataSourceName1)
	if err != nil {
		panic(err)
	}

	//返回结构体
	return con
}

//实现方法
func SqlQuery(e *xorm.Engine, sql string) []map[string][]byte {
	str, err := e.Query(sql)
	if err != nil {
		panic(err)
	}
	return str
}
func (con *connectInfo) GetColumns(host string, port int) {
	var mysqlMetadataColumns MysqlMetadataColumns
	sql := "SELECT table_schema,table_name,column_name,column_type,collation_name,is_nullable," +
		"column_key,column_default,extra,privileges,column_comment FROM information_schema.columns " +
		"where table_schema not in ('sys', 'test', 'information_schema', 'performance_schema', 'mysql')"
	res := SqlQuery(con.SourceConn, sql)
	for _, v := range res {

		md5 := fmt.Sprintf("%s%d%s%s%s%s%s%s%s%s%s%s%s%", host, port, string(v["table_schema"]),
			string(v["table_name"]), string(v["column_name"]), string(v["column_type"]), string(v["collation_name"]),
			string(v["is_nullable"]), string(v["column_key"]), string(v["column_default"]), string(v["extra"]),
			string(v["col_privileges"]), string(v["column_comment"]))
		md5 = strings.Replace(md5, " ", "", -1)
		md5 = MD5(md5)

		mysqlMetadataColumns.Node_addr = host
		mysqlMetadataColumns.Node_port = port
		mysqlMetadataColumns.Table_schema = string(v["table_schema"])
		mysqlMetadataColumns.Table_name = string(v["table_name"])
		mysqlMetadataColumns.Column_name = string(v["column_name"])
		mysqlMetadataColumns.Column_type = string(v["column_type"])
		mysqlMetadataColumns.Collation_name = string(v["collation_name"])
		mysqlMetadataColumns.Is_nullable = string(v["is_nullable"])
		mysqlMetadataColumns.Column_key = string(v["column_key"])
		mysqlMetadataColumns.Column_default = string(v["column_default"])
		mysqlMetadataColumns.Extra = string(v["extra"])
		mysqlMetadataColumns.Col_privileges = string(v["privileges"])
		mysqlMetadataColumns.Column_comment = string(v["column_comment"])
		mysqlMetadataColumns.Column_md5 = md5

		sql1 := fmt.Sprintf("select id,Node_addr,Node_port,table_schema,table_name,Column_name,Column_md5 "+
			"from aa.mysql_metadata_columns "+
			" where Node_addr='%s' and Node_port=%d and  table_schema='%s' and table_name='%s' and Column_name='%s' ",
			host, port, string(v["table_schema"]), string(v["table_name"]), string(v["column_name"]))
		//查询平台数据
		res1 := SqlQuery(con.DeractionConn, sql1)

		ok, _ := con.DeractionConn.SQL(sql1).Exist()
		if !ok {
			_, err := con.DeractionConn.Insert(&mysqlMetadataColumns)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("insert inito mysqlMetadataColumns ok !")
			}

		} else {
			for _, v1 := range res1 {
				if md5 == string(v1["md5"]) {
					fmt.Printf("%s table not changed \n", mysqlMetadataColumns.Table_name)
				} else {
					_, err := con.DeractionConn.Exec("update aa.mysql_metadata_columns set column_name=?,"+
						"column_type=?,collation_name=?,is_nullable=?,column_key=?,column_default=?,extra=?,"+
						"col_privileges=?,column_comment=?,column_md5=? where id=?",
						mysqlMetadataColumns.Column_name, mysqlMetadataColumns.Column_type,
						mysqlMetadataColumns.Collation_name, mysqlMetadataColumns.Is_nullable,
						mysqlMetadataColumns.Column_key, mysqlMetadataColumns.Column_default,
						mysqlMetadataColumns.Extra, mysqlMetadataColumns.Col_privileges, mysqlMetadataColumns.Column_comment,
						mysqlMetadataColumns.Column_md5, v1["id"])
					if err != nil {
						panic(err)
					}

				}
			}
		}
	}
}
func (con *connectInfo) GetIndexs(host string, port int) {
	var mysqlMetadataIndexs MysqlMetadataIndexs
	sql := "select table_schema,table_name,column_name,non_unique,index_name,seq_in_index,cardinality," +
		"nullable,index_type,index_comment from information_schema.statistics " +
		"where table_schema not in ('sys', 'test', 'information_schema', 'performance_schema', 'mysql')"
	res := SqlQuery(con.SourceConn, sql)
	for _, v := range res {

		mysqlMetadataIndexs.Node_addr = host
		mysqlMetadataIndexs.Node_port = port
		mysqlMetadataIndexs.Table_schema = string(v["table_schema"])
		mysqlMetadataIndexs.Table_name = string(v["table_name"])
		mysqlMetadataIndexs.Column_name = string(v["column_name"])
		mysqlMetadataIndexs.Non_unique = GetInt(string(v["non_unique"]))
		mysqlMetadataIndexs.Index_name = string(v["index_name"])
		mysqlMetadataIndexs.Seq_in_index = GetInt(string(v["seq_in_index"]))
		mysqlMetadataIndexs.Cardinality = GetInt(string(v["cardinality"]))
		mysqlMetadataIndexs.Nullable = string(v["nullable"])
		mysqlMetadataIndexs.Index_type = string(v["index_type"])
		mysqlMetadataIndexs.Index_comment = string(v["index_comment"])
		md5 := fmt.Sprintf("%s%d%s%s%s%s%s%s%s%s%s%s", host, port, string(v["table_schema"]),
			string(v["table_name"]), string(v["column_name"]), string(v["non_unique"]), string(v["index_name"]),
			string(v["seq_in_index"]), string(v["cardinality"]), string(v["nullable"]), string(v["index_type"]),
				string(v["index_comment"]))
		md5=MD5(md5)

		mysqlMetadataIndexs.Index_md5 = md5
		sql1 := fmt.Sprintf("select id,table_schema,table_name from aa.mysql_metadata_indexs "+
			"where node_addr='%s' and  node_port=%d and table_schema='%s' and table_name='%s'",
			host, port, string(v["table_schema"]), string(v["table_name"]))
		//查询平台数据
		res1 := SqlQuery(con.DeractionConn, sql1)

		ok, _ := con.DeractionConn.SQL(sql1).Exist()
		if !ok {
			_, err := con.DeractionConn.Insert(&mysqlMetadataIndexs)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("insert inito mysqlMetadataTtables ok !")
			}

		} else {
			for _, v1 := range res1 {
				if md5 == string(v1["md5"]) {
					fmt.Printf("%s table not changed \n", mysqlMetadataIndexs.Table_name)
				} else {
					_, err := con.DeractionConn.Exec("update aa.mysql_metadata_indexs set column_name=?,"+
						"non_unique=?,index_name=?,seq_in_index=?,cardinality=?,nullable=?,index_type=?,"+
						"index_comment=?,index_md5=? where id=?",
						mysqlMetadataIndexs.Column_name, mysqlMetadataIndexs.Non_unique, mysqlMetadataIndexs.Index_name,
						mysqlMetadataIndexs.Seq_in_index, mysqlMetadataIndexs.Cardinality, mysqlMetadataIndexs.Nullable,
						mysqlMetadataIndexs.Index_type,mysqlMetadataIndexs.Index_comment,mysqlMetadataIndexs.Index_md5,
							v1["id"])
					if err != nil {
						panic(err)
					}

				}
			}
		}

	}
}
func (con *connectInfo) GetTable(host string, port int) {
	var mysqlMetadataTtables MysqlMetadataTables
	sql := "SELECT table_schema, table_name, engine, row_format, table_rows, avg_row_length,data_length," +
		" max_data_length,index_length, data_free, auto_increment,table_collation, table_comment," +
		" create_time, update_time, check_time FROM information_schema.tables " +
		"where table_schema not in ('sys', 'test', 'information_schema', 'performance_schema', 'mysql')"
	res := SqlQuery(con.SourceConn, sql)

	for _, v := range res {
		mysqlMetadataTtables.Table_schema = string(v["table_schema"])
		mysqlMetadataTtables.Table_name = string(v["table_name"])
		mysqlMetadataTtables.Db_engine = string(v["engine"])
		mysqlMetadataTtables.Row_format = string(v["row_format"])
		mysqlMetadataTtables.Node_addr = host
		mysqlMetadataTtables.Node_port = port
		mysqlMetadataTtables.Table_rows = GetInt(string(v["table_rows"]))
		mysqlMetadataTtables.Avg_row_length = GetInt(string(v["avg_row_length"]))
		mysqlMetadataTtables.Max_data_length = GetInt(string(v["max_data_length"]))
		mysqlMetadataTtables.Data_length = GetInt(string(v["data_length"]))
		mysqlMetadataTtables.Index_length = GetInt(string(v["index_length"]))
		mysqlMetadataTtables.Data_free = GetInt(string(v["data_free"]))
		mysqlMetadataTtables.Chip_size = 0
		mysqlMetadataTtables.Auto_increment = GetInt(string(v["auto_increment"]))
		mysqlMetadataTtables.Table_collation = string(v["table_collation"])
		mysqlMetadataTtables.Create_time = GetDatetime(string(v["create_time"]))
		mysqlMetadataTtables.Update_time = GetDatetime(string(v["create_time"]))
		mysqlMetadataTtables.Check_time = GetDatetime(string(v["create_time"]))
		mysqlMetadataTtables.Table_comment = string(v["table_comment"])

		md5 := fmt.Sprintf("%s%d%s%s", host, port, string(v["table_schema"]), string(v["table_name"]))
		md5 = strings.Replace(md5, " ", "", -1)
		md5 = MD5(md5)
		mysqlMetadataTtables.Table_md5 = md5

		sql1 := fmt.Sprintf("select id,table_schema,table_name from aa.mysql_metadata_tables "+
			"where node_addr='%s' and  node_port=%d and table_schema='%s' and table_name='%s'",
			host, port, string(v["table_schema"]), string(v["table_name"]))
		//查询平台数据
		res1 := SqlQuery(con.DeractionConn, sql1)

		ok, _ := con.DeractionConn.SQL(sql1).Exist()
		if !ok {
			_, err := con.DeractionConn.Insert(&mysqlMetadataTtables)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("insert inito mysqlMetadataTtables ok !")
			}

		} else {
			for _, v1 := range res1 {
				if md5 == string(v1["md5"]) {
					fmt.Printf("%s table not changed \n", mysqlMetadataTtables.Table_name)
				} else {
					_, err := con.DeractionConn.Exec("update aa.mysql_metadata_tables set db_engine=?,"+
						" row_format=?,avg_row_length=?,max_data_length=?,data_length=?,index_length=?,data_free=?,"+
						"chip_size=?,auto_increment=?,table_collation=?,update_time=?,table_md5=? where id=?",
						mysqlMetadataTtables.Db_engine, mysqlMetadataTtables.Row_format,
						mysqlMetadataTtables.Avg_row_length, mysqlMetadataTtables.Max_data_length,
						mysqlMetadataTtables.Data_length, mysqlMetadataTtables.Index_length,
						mysqlMetadataTtables.Data_free, mysqlMetadataTtables.Chip_size,
						mysqlMetadataTtables.Auto_increment, mysqlMetadataTtables.Table_collation,
						mysqlMetadataTtables.Update_time, mysqlMetadataTtables.Table_md5, v1["id"])
					if err != nil {
						panic(err)
					}

				}
			}
		}

		//fmt.Println(host, port, string(v["table_schema"]), string(v["table_name"]), string(v["table_rows"]))
	}
}
func (con *connectInfo) UserHost(host string, port int) {
	var user map[string]string
	var cmdb_user Userinfo
	sql := "select user,host from mysql.user where user not like 'mysql.%' and user!=''" +
		" and host !='::1' and host not like '%localdomain%'"
	res := SqlQuery(con.SourceConn, sql)
	for _, v := range res {
		query := fmt.Sprintf("show grants for %s@'%s';", v["user"], v["host"])
		str := SqlQuery(con.SourceConn, query)
		for _, v := range str {
			for _, v1 := range v {
				if len(split(string(v1))) != 0 {
					user = split(string(v1))
				}
				query1 := fmt.Sprintf("select id,node_addr,node_port,account,login_addr,md5 "+
					"from aa.userinfo where  node_addr='%s' and node_port=%d "+
					"and account=%s and login_addr ='%s'", host, port, user["user"], user["hosts"])

				ss := fmt.Sprintf("%s%s%s%s%s%s", user["user"], user["grants"], user["database"],
					user["table"], user["hosts"], user["option"])
				ss = strings.Replace(ss, " ", "", -1)
				ss = strings.Replace(ss, "`", "", -1)
				ss = strings.Replace(ss, "'", "", -1)
				ss = MD5(ss)
				cmdb_user.Account = user["user"]
				cmdb_user.Account = strings.Replace(cmdb_user.Account, "'", "", -1)
				cmdb_user.Account = strings.Replace(cmdb_user.Account, " ", "", -1)
				cmdb_user.Db_name = user["database"]
				cmdb_user.CreateTime = time.Now()
				cmdb_user.Is_grant = GetInt(user["option"])
				cmdb_user.Login_addr = user["hosts"]
				cmdb_user.Node_addr = host
				cmdb_user.Node_port = port
				cmdb_user.Ownership = "admin"
				cmdb_user.Passwd = "123456"
				cmdb_user.Priv = user["grants"]
				cmdb_user.Role = 4
				cmdb_user.Table_name = user["tables"]
				cmdb_user.UpdateTime = time.Now()
				cmdb_user.Md5 = ss
				rr := SqlQuery(con.DeractionConn, query1)

				ok, _ := con.DeractionConn.SQL(query1).Exist()
				if !ok {
					//fmt.Println(query1)
					affected, err := con.DeractionConn.Insert(&cmdb_user)
					if err != nil {
						panic(err)
					} else {
						fmt.Println("insert inito userinfo ok !")
					}
					fmt.Println(affected)

				} else {
					for _, v2 := range rr {
						if ss == string(v2["md5"]) {
							fmt.Printf("%s grant not changed \n", user["user"])
						} else {
							_, err := con.DeractionConn.Exec("update aa.userinfo set priv = ?,db_name = ? ,"+
								"table_name = ?,update_time = ? ,md5 = ?"+
								" where id = ? ", cmdb_user.Priv, cmdb_user.Db_name, cmdb_user.Table_name,
								cmdb_user.UpdateTime, ss, v2["id"])
							if err != nil {
								panic(err)
							}

						}
					}
				}

			}

		}

	}

}
func GetInt(s string) int {
	if len(s) == 0 {
		return 0
	}
	r, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("--:%s:--", s)
		panic(fmt.Sprintf("%v---%s---", err, s))
	}
	return r
}
func GetDatetime(s string) time.Time {
	datetime, _ := time.Parse("2006-01-02 15:04:05", s)
	datetime.Format("2006-01-02 15:04:05")
	return datetime
}
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}
func split(str string) map[string]string {
	var grants string
	var database string
	var table string
	var user string
	var hosts string
	var option string
	//var results string
	res := make(map[string]string)
	if !strings.Contains(str, "PROXY ON") && !strings.Contains(str, "USAGE") {
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
		//hosts= strings.Replace(hosts, " ", "", -1)

		//fmt.Println(len(strings.Split(hosts,"WITH")))
		if strings.HasSuffix(str, "GRANT OPTION") {
			option = "1"
		} else {
			option = "0"
		}
		//results=fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",grants,database,table,user,hosts,option)
		res["grants"] = grants
		res["database"] = database
		res["table"] = table
		res["user"] = user
		res["hosts"] = hosts
		res["option"] = option

		//fmt.Println(res)
	} else {

		return nil
	}
	return res
}
func readConf(fileName string) [][]string {
	var str []string
	//var res []string
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
func main() {
	ss := readConf("E:\\go\\src\\DB.conf")
	for _, res := range ss {
		port, err := strconv.Atoi(res[3])
		if err != nil {
			panic(err)
		}
		//fmt.Println(res)
		con := NewConnectInfo(res[0], res[1], res[2], port, "mysql",
			"root", "123456", "192.168.160.134",
			3306, "aa")
		//con.UserHost(res[2],port)
		con.GetTable(res[2], port)
		con.GetColumns(res[2], port)
		con.GetIndexs(res[2], port)
	}

	//WG.Add(10)
	/*for i:=0;i<=9;i++{
		go readConf()
	}
	WG.Wait()*/

}
