package main

import (
	"os"
	"bufio"
	"io"
	"strings"
	"fmt"
)

func getBinlogFilePos(file *os.File) {
	//ret := make(map[string]string)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		// 表示为非 gtid 模式
		if strings.Contains(line, "CHANGE MASTER TO") {

			line = line[strings.Index(line,"'")+1:len(line)]
			s := strings.Split(line, "'")
			masterFile := s[0]
			po := strings.Split(s[1], ",")
			po = strings.Split(po[1], ";")
			po = strings.Split(po[0], "=")
			fmt.Println("-------",masterFile,po[1])
			break
		}
	}


}
func main(){
	file,_:=os.Open("/data/192.168.160.132_3388aa_20190409_110221/192.168.160.132_3388aa_20190409_110221.sql")
	//var file *os.File

	getBinlogFilePos(file)
}
