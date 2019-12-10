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
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shamsher31/goimgtype"
)

//
func PathIsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
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
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	sFolderPath, errFolderPath := MakePath(sPath)

	if errFolderPath != nil {
		fmt.Println("MakeFolder: " + errFolderPath.Error())
		return false, errFolderPath
	}

	errFolderPath = os.Mkdir(sFolderPath, os.ModePerm)

	if errFolderPath != nil {
		fmt.Println("MakeFolder: " + errFolderPath.Error())
		return false, errFolderPath
	}

	return true, nil
}

func SaveBase64AsImage(imageContent, targetPath string) (bool, error) {
	if imageContent == "" {
		var errMsg = "SaveBase64AsImage : image content is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if targetPath == "" {
		var errMsg = "SaveBase64AsImage : target file path is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if PathIsExist(targetPath) {
		bDelete := DeleteFile(targetPath)
		if bDelete == false {
			var errMsg = "SaveBase64AsImage : target Path already exist and cannot delete"
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}

	if strings.Contains(imageContent, "data:") == false || strings.Contains(imageContent, ";base64,") == false {
		var errMsg = "SaveBase64AsImage : Image Content Format Error"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var base64Index = strings.Index(imageContent, ";base64,")
	var base64Image = imageContent[base64Index+8:]

	decodedImage, errDecode := base64.StdEncoding.DecodeString(base64Image)
	if errDecode != nil {
		var errMsg = "SaveBase64AsImage : Cannot Decode Base64 Image"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}
	err2 := ioutil.WriteFile(targetPath, decodedImage, 0666)

	if err2 != nil {
		var errMsg = "SaveBase64AsImage : Cannot Save image"
		fmt.Println(errMsg)
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
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(image)

	datatype, err2 := imgtype.Get(imagePath)
	if err2 != nil {
		var errMsg = "ReadImageAsBase64: Cannot get image type"
		fmt.Println(errMsg)
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
		fmt.Println("PathIsImage: " + err2.Error())
		return false
	}
	return true
}

func GetImageType(base64Image string) (string, error) {
	if base64Image == "" {
		var errMsg = "Get Image Type: base64Image is empty"
		fmt.Println(errMsg)
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
				fmt.Println(errMsg)
				return "", errors.New(errMsg)
			}
		} else {
			var errMsg = "Get Image Type : Cannot get image type"
			fmt.Println(errMsg)
			return "", errors.New(errMsg)
		}
	}

	var errMsg = "Get Image Type : Cannot get image type"
	fmt.Println(errMsg)
	return "", errors.New(errMsg)
}

func MakePath(sPath string) (string, error) {
	sfolder, sfile := filepath.Split(sPath)

	if sfolder == "" || sfile == "" {
		var errMsg = "MakePath: folder or file name is empty folder " + sfolder + " file " + sfile
		fmt.Println(errMsg)
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
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	targetExist := PathIsExist(linkFolder)

	if targetExist {
		var errMsg = "Make Soft Link 4 Folder:linkFolder Already Exist " + linkFolder
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	sysType := runtime.GOOS

	var mkLinkCmd *exec.Cmd

	if "windows" == sysType {
		mkLinkCmd = exec.Command("cmd", "/c", "mklink /j "+linkFolder+"  "+srcFolder)
	} else if "linux" == sysType || "darwin" == sysType {
		mkLinkCmd = exec.Command("bash", "-c", "ln -s "+srcFolder+" "+linkFolder)
	} else { //Not support other platforms now
		return false, nil
	}

	_, err := mkLinkCmd.Output()

	if err != nil {
		fmt.Println("MakeSoftLink: " + err.Error())
		return false, err
	}
	return true, nil
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		fmt.Println("CopyFile: " + err.Error())
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		var errMsg = "CopyFile " + src + "is not a regular file"
		fmt.Println("CopyFile: " + errMsg)
		return 0, errors.New(errMsg)
	}

	source, err := os.Open(src)
	if err != nil {
		fmt.Println("CopyFile: " + err.Error())
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		fmt.Println("CopyFile: " + err.Error())
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func MoveFile(src, dst string) (int64, error) {
	iCopy, errCopy := CopyFile(src, dst)

	if errCopy != nil {
		fmt.Println("MoveFile: " + errCopy.Error())
		return 0, errCopy
	}

	errRemove := os.Remove(src)

	if errRemove != nil {
		fmt.Println("MoveFile: " + errRemove.Error())
		return 0, errRemove
	}

	return iCopy, nil
}

func DeleteFile(filePath string) bool {
	errRemove := os.Remove(filePath)

	if errRemove != nil {
		return false
	}
	return true
}

func Try2FindSpFile(siteFolderPath string) (string, error) {
	if PathIsExist(siteFolderPath) == false {
		var errMsg = "Try2FindSpFile: Site Folder not exist"
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var spCount int
	spCount = 0
	var spFileName string
	spFileName = ""

	files, _ := ioutil.ReadDir(siteFolderPath)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sp") {
			spFileName = f.Name()
			spCount++
			if spCount > 1 {
				var errMsg = "Try2FindSpFile: More than 1 .sp file"
				fmt.Println(errMsg)
				return "", errors.New(errMsg)
			}
		}
	}
	return spFileName, nil
}
