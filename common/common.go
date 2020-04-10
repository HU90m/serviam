package common

import (
    "log"
    "os"
    "os/exec"
    "regexp"
)


//---------------------------------------------------------------------------
// Shared Constants
//---------------------------------------------------------------------------
const INDENT = "\t"


//---------------------------------------------------------------------------
// Shared Functions
//---------------------------------------------------------------------------
//
// Panics if passed an error
//
func CheckErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

//
// Checks a directory exist. If one doesn't, it makes it.
//
func CheckDir(path string) {
	var info os.FileInfo
	var err error

	info, err = os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			log.Fatalf("'%s' is not a directory.", path)
		}
	} else {
		if os.IsNotExist(err) {
			log.Printf("the directory '%s' does not exist.\n", path)
			log.Printf("making '%s' directory.\n", path)
			os.MkdirAll(path, 0755)
		} else {
			log.Fatal(err)
		}
	}
}

//
// Removes characters that are not in the posix filename standard
//
func PosixFileName(input string) string {
    reg := regexp.MustCompile("([^\\w.-])+")
    return reg.ReplaceAllString(input, "")
}

//
// Creates a file containing the bytes given
//
func SaveBlob(blob []byte, location string) {
	var err error
	var file_p *os.File
	file_p, err = os.Create(location)
	CheckErr(err)
	_, err = file_p.Write(blob)
	CheckErr(err)
	err = file_p.Close()
	CheckErr(err)
}

//
// Displays Image using sxiv
//
func DisplayImage(image_location string) {
	var cmd *exec.Cmd
	log.Printf("Displaying '%s'.\n", image_location)
	cmd = exec.Command("sxiv", image_location)
	bytes, err := cmd.CombinedOutput()
	os.Stdout.Write(bytes)
	CheckErr(err)
}

