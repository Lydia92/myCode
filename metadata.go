package main

import (
	"Athena/util/encrypt"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
		_ "github.com/go-sql-driver/mysql"
	"time"
)

//原表信息
type InformationTable struct {
	Table_schema    string
	Table_name      string
	Engine          string
	Row_format      string
	Table_rows      uint64
	Avg_row_length  uint64
	Data_length     uint64
	Max_data_length uint64
	Index_length    uint64
	Data_free       uint64
	Auto_increment  uint64
	Table_collation string
	Table_comment   string
	Create_time     time.Time
	Update_time     time.Time
	Check_time      time.Time
}
type MysqlMetadataTables struct {
	Id             uint64    `orm:"auto" description:"主键"`
	NodeAddr       string    `orm:"size(15)" description:"表示实例地址"`
	NodePort       uint16    `orm:"default(3001)" description:"表示实例端口, node_addr+node_port 表示唯一一个实例"`
	TableSchema    string    `orm:"size(64)" description:"库名"`
	TableName      string    `orm:"size(64)" description:"表名"`
	DbEngine       string    `orm:"size(64)" description:"引擎"`
	RowFormat      string    `orm:"size(10)" description:"行格式"`
	TableRows      uint64    `orm:"size(21)" description:"行数"`
	AvgRowLength   uint64    `orm:"size(21)" description:"平均行长度"`
	MaxDataLength  uint64    `orm:"size(21)" description:"最大行长度的"`
	DataLength     uint64    `orm:"size(21)" description:"数据长度"`
	IndexLength    uint64    `orm:"size(21)" description:"索引长度"`
	DataFree       uint64    `orm:"size(21)" description:"空闲"`
	ChipSize       uint64    `orm:"size(21)" description:"碎片"`
	AutoIncrement  uint64    `orm:"size(21)" description:"下一个自增的值"`
	TableCollation string    `orm:"size(32)" description:"字符集"`
	CreateTime     time.Time `orm:"type(datetime)" description:"表创建时间"`
	UpdateTime     time.Time `orm:"null; type(datetime)" description:"表更新时间"`
	CheckTime      time.Time `orm:"null; type(datetime)" description:"表检测时间"`
	TableComment   string    `orm:"size(2048)" description:"表注释"`
	TableMd5       string    `orm:"type(char)" description:"HOST+PORT+db_name+tab_name生成MD5,判断是否变化过"`
}

//库名
type Database struct {
	TableSchema string
}


func init(){
	user:="root"
	pass:="123456"
	address:="192.168.160.133"
	port:=3306
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, pass, address, port, "aa")
	//注册数据库，连接对应的生产环境
	orm.RegisterDataBase("default", "mysql", dns)
	orm.RegisterModel(new(MysqlMetadataTables),)
	orm.RunSyncdb("default", false, true)





}



//获取系统元表信息
func GetMetaTable(DB1 orm.Ormer)[]InformationTable{
	var informationTable []InformationTable
	//var res []map[string]string
	sql := "SELECT table_schema, table_name, engine, row_format, table_rows, avg_row_length,data_length," +
		" max_data_length,index_length, data_free, auto_increment,table_collation, table_comment," +
		" create_time, update_time, check_time FROM information_schema.tables " +
		"where table_schema not in ('sys', 'test', 'information_schema', 'performance_schema', 'mysql')"
	DB1.Raw(sql).QueryRows(&informationTable)
	  return  informationTable


}



func main(){
	CompareMetaTable()
}
//插入元表信息
func (mysqlMetadataTables *MysqlMetadataTables)InsertTableInfo(Host, Table_schema, Table_name, Engine, Row_format string, Table_rows, Avg_row_length,
Data_length, Max_data_length, Index_length, Data_free, Auto_increment uint64,
	Table_collation, Table_comment, Md5 string, Create_time, Update_time, check_time time.Time, port uint16) {
	DB := orm.NewOrm()
	mysqlMetadataTables = &MysqlMetadataTables{
		NodeAddr:       Host,
		NodePort:       port,
		TableSchema:    Table_schema,
		TableName:      Table_name,
		DbEngine:       Engine,
		RowFormat:      Row_format,
		TableRows:      Table_rows,
		AvgRowLength:   Avg_row_length,
		MaxDataLength:  Max_data_length,
		DataLength:     Data_length,
		IndexLength:    Index_length,
		DataFree:       Data_free,
		ChipSize:       0,
		AutoIncrement:  Auto_increment,
		TableCollation: Table_collation,
		CreateTime:     Create_time,
		UpdateTime:     Update_time,
		CheckTime:      check_time,
		TableComment:   Table_comment,
		TableMd5:       Md5,
	}
	DB.Insert(mysqlMetadataTables)
}

//对比元表数据与平台数据，更新新数据
func CompareMetaTable() {
	DB1 := orm.NewOrm()
	infomation:=MysqlMetadataTables{}
		informationTable:=GetMetaTable(DB1)
		for _,v:=range informationTable{
			fmt.Println("-==================",v.Table_name)
			md5 := fmt.Sprintf("%s%d%s%s", "192.168.160.163", 3306, v.Table_schema, v.Table_name)
			md5 = strings.TrimSpace(md5)
			md5 = encrypt.GenerateMD5(md5)
			infomation.InsertTableInfo("192.168.160.162", v.Table_schema, v.Table_name,
				v.Engine, v.Row_format, v.Table_rows,
				v.Avg_row_length, v.Data_length, v.Max_data_length,
				v.Index_length, v.Data_free, v.Auto_increment,
				v.Table_collation, v.Table_comment, md5, v.Create_time,
				v.Update_time, v.Check_time, 3388)
		}

}




