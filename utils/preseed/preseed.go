package redfish

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"gopkg.in/yaml.v2"
	"regexp"
	"strconv"
//	"errors"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// Log handlder initialisation
func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// type FileSystem struct {
// 	MountPoint   string `yaml:"mountpoint"`
// 	FsType       string `yaml:"fstype"`
// 	MountOptions string `yaml:"mount-options"`
// }
//
// type Partitions struct {
// 	Name       string     `yaml:"disk"`
// 	Size       string     `yaml:"size"`
// 	MaxSize    string
// 	Primary    bool  `yaml:"primary"`
// 	Bootable   bool       `yaml:"bootable"`
// 	FileSystem FileSystem `yaml:"filesystem"`
// }
//
// type PhysicalDevice struct {
// 	Disk       string       `yaml:"disk"`
// 	Partitions []Partitions `yaml:"partitions"`
// }
//
// type Storage struct {
// 	PhysicalDevices []PhysicalDevice `yaml:"physical-devices"`
// 	OSDisk   string
// 	Hostname string
// 	NoOfPartitions int
// }
//
// // HardwareProfile sets the hardware raid levels and nics requested by user
// // on the host.
// type HostProfile struct {
// 	Storage Storage `yaml:"storage"`
// }
//
// type Platform struct {
//         GrubConfig string `yaml:"grub_config"`
// }

type hostProfile struct {
		hp	*host.HostProfile
}

// func (hp *HostProfile) readFromYaml(storageProfileYaml string) (hostProfile *HostProfile, err error) {
// 	hprofile  :=  HostProfile {}
//
// 	if len(storageProfileYaml)  != 0 {
// 		hprofile.Storage = Storage{ }
// 		 err = yaml.Unmarshal([]byte(storageProfileYaml), &hprofile.Storage)
// 	}
// 	if err != nil{
// 		Error.Println("Failed to read yaml file")
// 	}
// 	return &hprofile, err
//
// }

func New(hp string) *hostProfile {
		hostProf, _ := host.GetHostProfile(hp)
		return &hostProfile(hp: hostProf)
}

func (hp *HostProfile) CreatePreseedfile (hostname string , templatepath string, preeseedpath string) (err error) {

	/*
		file, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", output, ":", err)
		}

		multi := io.MultiWriter(file, os.Stdout)
	*/
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	Info.Println("Inside CreateaPreseeedfile()")

	hostprofile := hp.hp

	if err != nil {
		Error.Println("Failed to read storage profile from config file")
		return err
	}

	//check the structure for correctness and update the values in MB or -1 as the case is
	if len(hostprofile.Storage.PhysicalDevices)  != 0 {

	for x , phydevice := range hostprofile.Storage.PhysicalDevices{
		hostprofile.Storage.OSDisk = phydevice.Disk
		hostprofile.Storage.NoOfPartitions = len(phydevice.Partitions)
		for y, partition := range phydevice.Partitions {
            disksize, maxDiskSpace, err := getDiskSpaceinMB(partition.Size)
			if err != nil {
				return err
			}

			hostprofile.Storage.PhysicalDevices[x].Partitions[y].Size = disksize
			hostprofile.Storage.PhysicalDevices[x].Partitions[y].MaxSize = maxDiskSpace
		}

	}
}
	//Update the OSDisk and Hostname into the hostprofile structure

	hostprofile.Storage.Hostname = hostname

	tmpl, err := template.ParseFiles(templatepath)
	if err != nil {
		return err
	}

	f, err := os.Create(preeseedpath)
    if err != nil {
       Error.Println("Failed to create file: ")
       return err
    }
	err = tmpl.Execute(f, *hostprofile)

	return err
}

func getDiskSpaceinMB ( diskspace string ) ( diskspaceinMB string, maxdiskSizeinMB string, err error ) {
	//check if there is diskspace listed with a numeber followed by g
	// if there is > in front of the
	re := regexp.MustCompile(`(>*)(\d+)([a-z])`)
	t := re.FindSubmatch([]byte(diskspace))
	if len(t) == 4 {
		disksizeGB, err := strconv.Atoi(string(t[2]))
		disksizeMB := strconv.Itoa(disksizeGB*1024)
        var maxdiskSizeinMB string
		if string(t[1]) == ">" {
            maxdiskSizeinMB = "-1"
		}else {
			maxdiskSizeinMB = disksizeMB
		}
         return  disksizeMB, maxdiskSizeinMB, err
	}
	return "", "", err

}

func (hp *HostProfile) CreateGrubfile(grubfilepath string, grubtemplatepath string) (err error) {

        platProfile := hp.hp.Platform
        if err != nil {
                Error.Println("Failed to read storage profile from config file")
                return err
        }

        tmpl, err := template.ParseFiles(grubtemplatepath)
        if err != nil {
                return err
        }
        f, err := os.Create(grubfilepath)
        if err != nil {
                Error.Println("Failed to create file: ")
                return err
        }
        err = tmpl.Execute(f, *platProfile)
        return err

}

// func readFromGrubYaml(platformProfileYaml string) ( platform *Platform, err  error) {
//         platProfile  :=  Platform {}
//
//         if len(platformProfileYaml)  != 0 {
//                 err = yaml.Unmarshal([]byte(platformProfileYaml), &platProfile)
//         }
//         if err != nil{
//                 Error.Println("Failed to read yaml file")
//         }
//         return &platProfile, err
// }
