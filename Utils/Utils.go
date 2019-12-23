package Utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shamsher31/goimgtype"
)

//
func PathIsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		//Logger.Println(err.Error())
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GUID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		Logger.Println("GUID" + err.Error())
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func CurrentTime() string {
	t := time.Now()
	str := t.Format("2006-01-02 15:04:05")

	return str
}

func PathIsMarkdown(filePath string) bool {
	if filePath == "" {
		return false
	}

	ext := filepath.Ext(filePath)

	if ext == ".md" || ext == ".markdown" || ext == ".mdown" || ext == ".mmd" {
		return true
	}
	return false
}

func MakeFolder(sPath string) (bool, error) {
	if sPath == "" {
		var errMsg = "MakeFolder: sPath is empty"
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	sFolderPath, errFolderPath := MakePath(sPath)

	if errFolderPath != nil {
		Logger.Println("MakeFolder: " + errFolderPath.Error())
		return false, errFolderPath
	}

	errFolderPath = os.Mkdir(sFolderPath, os.ModePerm)

	if errFolderPath != nil {
		Logger.Println("MakeFolder: " + errFolderPath.Error())
		return false, errFolderPath
	}

	return true, nil
}

func SaveBase64AsImage(imageContent, targetPath string) (bool, error) {
	if imageContent == "" {
		var errMsg = "SaveBase64AsImage : image content is empty"
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if targetPath == "" {
		var errMsg = "SaveBase64AsImage : target file path is empty"
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if PathIsExist(targetPath) {
		bDelete := DeleteFile(targetPath)
		if bDelete == false {
			var errMsg = "SaveBase64AsImage : target Path already exist and cannot delete"
			Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}

	if strings.Contains(imageContent, "data:") == false || strings.Contains(imageContent, ";base64,") == false {
		var errMsg = "SaveBase64AsImage : Image Content Format Error"
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var base64Index = strings.Index(imageContent, ";base64,")
	var base64Image = imageContent[base64Index+8:]

	decodedImage, errDecode := base64.StdEncoding.DecodeString(base64Image)
	if errDecode != nil {
		var errMsg = "SaveBase64AsImage : Cannot Decode Base64 Image"
		Logger.Println(errMsg)
		Logger.Println(errDecode.Error())
		return false, errors.New(errMsg)
	}
	err2 := ioutil.WriteFile(targetPath, decodedImage, 0666)

	if err2 != nil {
		var errMsg = "SaveBase64AsImage : Cannot Save image"
		Logger.Println(errMsg)
		Logger.Println(err2.Error())
		return false, errors.New(errMsg)
	}

	return true, nil
}

func ReadImageAsBase64(imagePath string) (string, error) {

	var retImage string
	retImage = ""

	image, errRead := ioutil.ReadFile(imagePath)

	if errRead != nil {
		var errMsg = "ReadImageAsBase64: Read Fail"
		Logger.Println(errMsg)
		Logger.Println(errRead.Error())
		return "", errors.New(errMsg)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(image)

	datatype, err2 := imgtype.Get(imagePath)
	if err2 != nil {
		var errMsg = "ReadImageAsBase64: Cannot get image type"
		Logger.Println(errMsg)
		Logger.Println(err2.Error())
		return "", errors.New(errMsg)
	} else {
		retImage = "data:" + datatype + ";base64," + imageBase64
	}

	return retImage, nil
}

func PathIsImage(filePath string) bool {

	if filePath == "" {
		return false
	}

	_, err2 := imgtype.Get(filePath)
	if err2 != nil {
		Logger.Println("PathIsImage: " + err2.Error())
		return false
	}
	return true
}

func GetImageType(base64Image string) (string, error) {
	if base64Image == "" {
		var errMsg = "Get Image Type: base64Image is empty"
		Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var datatypeParts = strings.Split(base64Image, ";") //Get data:image/png
	if len(datatypeParts) > 1 {
		var datatypePart = datatypeParts[0]
		var datatypes = strings.Split(datatypePart, ":") //Get image/png
		if len(datatypes) == 2 {
			var datatype = datatypes[1]
			var subTypes = strings.Split(datatype, "/") //Get png
			if len(subTypes) == 2 {
				return subTypes[1], nil
			} else {
				var errMsg = "Get Image Type : Cannot get image type"
				Logger.Println(errMsg)
				return "", errors.New(errMsg)
			}
		} else {
			var errMsg = "Get Image Type : Cannot get image type"
			Logger.Println(errMsg)
			return "", errors.New(errMsg)
		}
	}

	var errMsg = "Get Image Type : Cannot get image type"
	Logger.Println(errMsg)
	return "", errors.New(errMsg)
}

func MakePath(sPath string) (string, error) {
	sfolder, sfile := filepath.Split(sPath)

	if sfolder == "" || sfile == "" {
		var errMsg = "MakePath: folder or file name is empty folder " + sfolder + " file " + sfile
		Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	sfolder = filepath.Clean(sfolder)

	if !PathIsExist(sfolder) {
		os.MkdirAll(sfolder, os.ModePerm)
	}

	return filepath.Join(sfolder, sfile), nil

}

func MakeSoftLink4Folder(srcFolder, linkFolder string) (bool, error) {
	srcExist := PathIsExist(srcFolder)

	if !srcExist {
		var errMsg = "Make Soft Link 4 Folder: SrcFolder Not Exist " + srcFolder
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	targetExist := PathIsExist(linkFolder)

	if targetExist {
		var errMsg = "Make Soft Link 4 Folder:linkFolder Already Exist " + linkFolder
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	errLink := os.Symlink(srcFolder, linkFolder)

	if errLink != nil {
		Logger.Println("MakeSoftLink: " + errLink.Error())
		return false, errLink
	}
	return true, nil
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		Logger.Println("CopyFile: " + err.Error())
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		var errMsg = "CopyFile " + src + "is not a regular file"
		Logger.Println("CopyFile: " + errMsg)
		return 0, errors.New(errMsg)
	}

	source, err := os.Open(src)
	if err != nil {
		Logger.Println("CopyFile: " + err.Error())
		return 0, err
	}
	defer source.Close()

	if PathIsExist(dst) == false {
		MakePath(dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		Logger.Println("CopyFile: " + err.Error())
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func CopyFileWithConfirm(src, dst string) (bool, error) {
	var msg = dst + " already exist, replace it?"
	var confirm = UserConfirm(msg)

	if confirm {
		_, errCopy := CopyFile(src, dst)
		if errCopy != nil {
			var errMsg = "SiteModule.AddFile fail " + errCopy.Error()
			Logger.Println(errMsg)
			return false, errCopy
		}
		return true, nil
	}
	return false, errors.New("User Ingore it")

}

func CopyFolder(src, dst string, addForce bool) (bool, error) {
	if src == "" {
		var errMsg = "Utils.CopyFolder: Src is empty"
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if dst == "" {
		var errMsg = "Utils.CopyFolder: Dst is empty"
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if PathIsExist(src) == false {
		var errMsg = "Utils.CopyFolder: " + src + " not exist"
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if PathIsDir(src) == false {
		var errMsg = "Utils.CopyFolder: " + src + " is not folder"
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	files, _ := ioutil.ReadDir(src)

	for _, f := range files {
		if f.IsDir() {
			var srcFolderPath = filepath.Join(src, f.Name())
			var dstFolderPath = filepath.Join(dst, f.Name())

			CopyFolder(srcFolderPath, dstFolderPath, addForce)
		} else {
			var srcFilePath = filepath.Join(src, f.Name())
			var dstFilePath = filepath.Join(dst, f.Name())

			if addForce {
				CopyFile(srcFilePath, dstFilePath)
			} else {
				CopyFileWithConfirm(srcFilePath, dstFilePath)
			}
		}
	}

	return true, nil
}

func MoveFile(src, dst string) (int64, error) {
	iCopy, errCopy := CopyFile(src, dst)

	if errCopy != nil {
		Logger.Println("MoveFile: " + errCopy.Error())
		return 0, errCopy
	}

	errRemove := os.Remove(src)

	if errRemove != nil {
		Logger.Println("MoveFile: " + errRemove.Error())
		return 0, errRemove
	}

	return iCopy, nil
}

func DeleteFile(filePath string) bool {
	err := os.Remove(filePath)

	if err != nil {
		Logger.Println(err.Error())
	}

	if PathIsExist(filePath) {
		return false
	}

	return true
}

func DeleteFolder(folderPath string) bool {
	if folderPath == "" {
		var errMsg = "Utils.DeleteFolder: folderPath is empty"
		Logger.Println(errMsg)
		return false
	}

	if PathIsExist(folderPath) == false {
		var errMsg = "Utils.DeleteFolder: " + folderPath + " not exist"
		Logger.Println(errMsg)
		return false
	}

	if PathIsDir(folderPath) == false {
		var errMsg = "Utils.DeleteFolder: " + folderPath + " is not folder"
		Logger.Println(errMsg)
		return false
	}

	bError := os.RemoveAll(folderPath)

	if bError != nil {
		return false
	}

	return true
}

func ClearFolder(folderPath string) bool {
	if folderPath == "" {
		var errMsg = "Utils.DeleteFolder: folderPath is empty"
		Logger.Println(errMsg)
		return false
	}

	if PathIsExist(folderPath) == false {
		var errMsg = "Utils.DeleteFolder: " + folderPath + " not exist"
		Logger.Println(errMsg)
		return false
	}

	if PathIsDir(folderPath) == false {
		var errMsg = "Utils.DeleteFolder: " + folderPath + " is not folder"
		Logger.Println(errMsg)
		return false
	}

	files, _ := ioutil.ReadDir(folderPath)

	for _, f := range files {
		var fPath = filepath.Join(folderPath, f.Name())
		os.RemoveAll(fPath)
	}

	return true
}

func Try2FindSpFile(siteFolderPath string) (string, error) {
	if PathIsExist(siteFolderPath) == false {
		var errMsg = "Try2FindSpFile: Site Folder not exist"
		Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var spCount int
	spCount = 0
	var spFileName string
	spFileName = ""

	files, errReadDir := ioutil.ReadDir(siteFolderPath)
	if errReadDir != nil {
		Logger.Println(errReadDir.Error())
		Logger.Println("Try2FindSpFile: Cannot read Dir")
		return "", errReadDir
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sp") {
			spFileName = f.Name()
			spCount++
			if spCount > 1 {
				var errMsg = "Try2FindSpFile: More than 1 .sp file"
				Logger.Println(errMsg)
				return "", errors.New(errMsg)
			}
		}
	}
	return spFileName, nil
}

var Logger *log.Logger

func InitLogger() {
	file, err := os.Create("ipsc.log")
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}

	Logger = log.New(file, "", log.Llongfile)
}

func PathIsFile(filePath string) bool {
	f, err := os.Stat(filePath)

	if err != nil {
		Logger.Println(err.Error())
		return false
	}

	if f.IsDir() == false {
		return true
	}

	return false
}

func PathIsDir(filePath string) bool {
	f, err := os.Stat(filePath)

	if err != nil {
		Logger.Println(err.Error())
		return false
	}

	return f.IsDir()
}

func UserConfirm(msg string) bool {
	fmt.Println(msg)
	fmt.Println("Enter Y/N (Default Y): ")

	var str string
	str = ""
	fmt.Scanf("%s", &str)

	if str == "" {
		return true
	}

	str = strings.ToUpper(str)

	if str == "Y" {
		return true
	}

	return false
}
