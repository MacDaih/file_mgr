package logs

import ( 
	"os"
	"log"
)
const (
	LOG = "./log.txt"
	NAME = "FILE_MGR =>"
) 

func WriteLog(trace string,info *string, crash error) {
	f, err := os.OpenFile(LOG, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed on logs init at %s : %s", LOG, err)
	}
	defer f.Close()
	logger := log.New(f, NAME , log.LstdFlags)
	if crash != nil {
		logger.Printf("%s : %v \n", trace, crash)
	} else {
		logger.Printf("%s : %s \n", trace, &info)
	}
}