//sources: https://www.devdungeon.com/content/working-files-go
//Renamer Program with info and copy functions (flags determine action and how action occurs)
	//Pass in path to wanted directory
	//Flags: -r = Rename files with options: addDate, addDir, addInc, remAll, remNums, remSpec
		//-i = Get file information with options: all, name
		//-c = Copy files in directory with options: copy

package main

import (
	"flag"
	"fmt"
	"io"
    "log"
    "os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"strconv"
)

var (
	dirPath string
	oldPath string
	newPath string
	newName string
	fileName string
	fileType string
	typeOnly string
	useName string
	fileInfo os.FileInfo
	sepChar string = "_"
	err      error
	choice = "none"
	root = "none"
	counter = 0;
)

//Walk Rename Function
func walkRename(files *[]string, choice string) filepath.WalkFunc {
	//Split flags apart and save in slice
	c := strings.Split(choice, ",")
	var types []string
	types = strings.Split(typeOnly, ",")
	name := useName
	//Go through everything in the directory, only working with files
    return func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Fatal(err)
		}
		if info.IsDir() == false {
			//Save the file info in holders
			oldPath = dirPath + "\\" + info.Name()
			index := 0
			s := []byte(".");
			for i := len(info.Name())-1; i > 0; i-- {
				if info.Name()[i] == s[0] {
					index = i
					break
				}
			}
			fileName = info.Name()[0:index]
			fileType = info.Name()[index:]

			for _, ty := range types {
				if ty == fileType || ty == "all"{

					if name != "none" {
						addUse(name, info)
					}
					//Loop through choices, calling functions when set by flag
					for _, val := range c {
						if (val == "remAll" && name == "none")|| (val == " remAll" && name == "none"){
							removeAll(info)
						}
						if val == "remNums" || val == " remNums"{
							removeNum(info)
						}
						if val == "remSpec" || val == " remSpec"{
							removeSpecChars(info)
						}
						if val == "addDate" || val == " addDate"{
							addDate(info)
						}
						if val == "addDir" || val == " addDir"{
							addDir(info)
						}
						if val == "addInc" || val == " addInc"{
							addInc(info)
						}
					}
		
					//Remove proceeding '_' if present
					s = []byte("_");
					if newName[0] == s[0] {
						newName = newName[1:]
					} 
					//Add fileType at end of new file name
					newName += fileType
					//Rename uses paths, so create new path with new file name
					newPath = dirPath + "\\" + newName
					//Rename the old file with the new name
					err := os.Rename(oldPath, newPath)
					if err != nil {
						log.Fatal(err)
					} else {
						//Show successful change
						fmt.Printf("%s changed to %s \n", info.Name(), newName)
					}
				}
			}
		}
        return nil
    }
}

//Walk Copy Function
func walkCopy(files *[]string, choice string) filepath.WalkFunc {
    return func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Fatal(err)
		}
		if info.IsDir() == false {
			copyBasic(info)
		}
        return nil
    }
}

//RENAME ADD FUNCTIONS
//Rename Add Date Function
func addDate(fileInfo os.FileInfo){
	crtTime := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	fullTime := time.Unix(0, crtTime.CreationTime.Nanoseconds()).String()
	crtDate := fullTime[0:10]
	newName = fileName + sepChar + crtDate
	fileName = newName
}
//Rename Add Directory Function
func addDir(fileInfo os.FileInfo){
	dir := dirPath
	index := 0
	s := []byte("\\");
	for i := len(dir)-1; i > 0; i-- {
		if dir[i] == s[0] {
			index = i
			break
		}
	}
	index++
	substring := dir[index:]
	dir = substring
	newName = fileName + sepChar + dir
	fileName = newName
}
//Rename Add Incrementor Function
func addInc(fileInfo os.FileInfo){
	counter++
	newName = fileName + sepChar + fmt.Sprint(counter)
	fileName = newName
}
//Rename Add Incrementor Function
func addUse(name string, fileInfo os.FileInfo){
	newName = name;
	fileName = newName
}

//RENAME REMOVE FUNCTIONS
//Rename Remove Numbers Function
func removeAll(fileInfo os.FileInfo){
	fileName = ""
}
//Rename Remove Numbers Function
func removeNum(fileInfo os.FileInfo){
	nameSlice := strings.Split(fileName,"")
	for i := len(nameSlice)-1; i >= 0; i-- {
		if _, err := strconv.Atoi(nameSlice[i]); err == nil {
			if i != len(nameSlice)-1 {
				for j := i; j < len(nameSlice)-1; j++{
					nameSlice[j] = nameSlice[j+1]
				}
			} 
			nameSlice = nameSlice[:len(nameSlice)-1]   
		}
	}
	newName = strings.Join(nameSlice, "")
	fileName = newName
}
//Rename Remove Special Characters Function
func removeSpecChars(fileInfo os.FileInfo){
	specChars := []byte{'!','@','#','$','%','^','&','*','(',')','_','-','[',']',';','.'}

	nameSlice := strings.Split(fileName,"")
	for i := len(nameSlice)-1; i >= 0; i-- {
		if contains(specChars, nameSlice[i]){
			if i != len(nameSlice)-1 {
				for j := i; j < len(nameSlice)-1; j++{
					nameSlice[j] = nameSlice[j+1]
				}
			} 
			nameSlice = nameSlice[:len(nameSlice)-1]  
		}
	}

	newName = strings.Join(nameSlice, "")
	fileName = newName
}

//Contains Function
func contains(s []byte, e string) bool {
    for _, a := range s {
		b := string(a);
        if e == b {
            return true
        }
    }
    return false
}


//COPY FUNCTIONS
//Copy Basic Function
func copyBasic(fileInfo os.FileInfo){
	index := 0
	s := []byte(".");
	for i := len(fileInfo.Name())-1; i > 0; i-- {
		if fileInfo.Name()[i] == s[0] {
			index = i
			break
		}
	}

	fileName := fileInfo.Name()[0:index]
	fileType := fileInfo.Name()[index:]
	newName := fileName + "_copy" + fileType

	originalFile, err := os.Open(fileName + fileType)
    if err != nil {
        log.Fatal(err)
    }
    defer originalFile.Close()

	 newFile, err := os.Create(newName)
	 if err != nil {
		 log.Fatal(err)
	 }
	 defer newFile.Close()

	 bytesWritten, err := io.Copy(newFile, originalFile)
	 if err != nil {
		 log.Fatal(err)
	 }

	 fmt.Printf("Copied %d bytes.", bytesWritten)
	 err = newFile.Sync()
	 if err != nil {
		 log.Fatal(err)
	 }
}


//Main Function
func main() {
	//Setup for Flags
	renamePtr := flag.String("r", "none", "rename file with additional info")
	usePtr := flag.String("u", "none", "name to use for rename if desired")
	sepPtr := flag.String("s", "_", "character to use as name seperator")
	typePtr := flag.String("t", "all", "type of file to rename")
	copyPtr := flag.String("c", "none", "copy file")
	flag.Parse()
	//Get and save the location for the files to manipulate
	root = strings.Join(flag.Args()," ")
	dirPath = root
	os.Chdir(dirPath)
	//Look for flags
	var files []string
	if (*sepPtr != "none"){
		sepChar = *sepPtr
	}
	if (*usePtr != "none"){
		useName = *usePtr
	}
	if (*typePtr != "all"){
		typeOnly = *typePtr
	}else{
		typeOnly = "all"
	}
	if (*copyPtr != "none"){
		err := filepath.Walk(root, walkCopy(&files, *copyPtr))
		if err != nil {
			panic(err)
		}
	}
	if (*renamePtr != "none"){
		err := filepath.Walk(root, walkRename(&files, *renamePtr))
		if err != nil {
			panic(err)
		}
	}
}