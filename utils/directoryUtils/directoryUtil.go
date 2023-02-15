package directoryUtils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

//删除指定目录下的空文件夹
func DeleteDir(rootDir string) {
	for {
		//扫描根路径下的所有文件夹
		filePaths, err := GetCurDirList(rootDir)
		if err != nil {
			fmt.Printf("获取根目录下层级目录[%v]失败,error info:%v\n", rootDir, err)
			return
		}
		if len(filePaths) == 0 {
			//fmt.Println("根目录下为空目录")
			time.Sleep(time.Second * 10)
			continue
		}
		fmt.Println("文件夹列表: ", filePaths)

		//去除当前日期的文件夹
		var fileDir []string
		curDir := time.Now().Format("20060102")
		for _, value := range filePaths {
			if value != curDir {
				fileDir = append(fileDir, value)
			}
		}
		if len(fileDir) == 0 {
			fmt.Println("没有需要删除的空目录")
			time.Sleep(time.Second * 10)
			continue
		}
		fmt.Println("过滤后文件夹列表: ", fileDir)

		//获取过滤后文件夹下所有目录列表
		var fileDirList []string
		for _, dir := range fileDir {
			fileDirList, err = GetDirList(rootDir + "\\" + dir)
			if err != nil {
				fmt.Println(fileDirList)
				continue
			}
		}
		fmt.Println("获取过滤后目录", fileDirList)

		//遍历删除空文件夹
		//for _, dirPath := range filelist {
		for i := len(fileDirList) - 1; i >= 0; i-- { //从里往外删除空文件夹
			var files []string
			files, _ = GetAllFile(fileDirList[i], files)
			if len(files) == 0 {
				fmt.Println("删除空文件夹目录", fileDirList[i])
				//删除空文件夹
				err := os.Remove(fileDirList[i])
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Printf("目录%v文件数量：%v \n", fileDirList[i], len(files))
			}
		}

		//获取前一天日期
		//d, _ := time.ParseDuration("-48h")
		//d1 := time.Now().Add(d).Format("20060102")
		//fmt.Println(d1)

		time.Sleep(time.Second * 10)
	}
}

// 获取目录下所有的文件夹，包括层级目录下
func GetDirList(dirpath string) ([]string, error) {
	var dir_list []string
	dir_err := filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			//if path != dirpath { //过滤1级文件夹
			//	dir_list = append(dir_list, path)
			//}
			dir_list = append(dir_list, path)
			return nil
		}
		return nil
	})
	return dir_list, dir_err
}

//获取指定目录下所有的文件名，包括层级目录下
func GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := os.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			//添加文件限制，防止一下子获取全部文件
			if len(s) > 10000 {
				fmt.Println("======10000===")
				return s, err
			}
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}


//获取指定目录下的所有目录，不进入下一级目录搜索
func GetCurDirList(dirPth string) ([]string, error) {
	var files []string
	dir, err := os.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		if fi.IsDir() {
			files = append(files, fi.Name())
		}
	}
	return files, nil
}

// 目录是否存在
func IsDirExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
