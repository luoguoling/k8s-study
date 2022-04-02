package main
import "fmt"
import "os"
import "log"
import "time"
//判断文件夹是否存在
func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}
//写日志
func Writelog(path,command string){
    filepath := path + "/cmd.txt"
    fmt.Println(filepath)
    logFile,err := os.OpenFile(filepath,os.O_RDWR|os.O_CREATE|os.O_APPEND,0777)
    if err != nil{
        log.Fatalln("读取文件日志失败!!!",err)
    }
    defer logFile.Close()
    logger := log.New(logFile,"\r",log.Ldate|log.Ltime)
    logger.Print(command)

}
func main(){
    command := os.Args[1]
    timestamp := fmt.Sprintf(time.Now().Format("2006-01-02"))
    path := fmt.Sprintf("/usr/bin/.hist/%s",timestamp)
    exist,_ := PathExists(path)
    if exist{
       // Writelog(path,command)
        //return
    } else {
        err := os.Mkdir(path,os.ModePerm)
        fmt.Println(err)
    }
    Writelog(path,command)
}
