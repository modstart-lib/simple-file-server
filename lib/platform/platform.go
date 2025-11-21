package platform

import (
	"github.com/shirou/gopsutil/v3/host"
	log "github.com/sirupsen/logrus"
)

type PlatformOSType string

const (
	PlatformOSType_Unknown PlatformOSType = "Unknown"
	PlatformOSType_Linux   PlatformOSType = "Linux"
	PlatformOSType_Windows PlatformOSType = "Windows"
	PlatformOSType_OSX     PlatformOSType = "OSX"
)

type PlatformArchType string

const (
	PlatformArchType_Unknown PlatformArchType = "Unknown"
	PlatformArchType_AMD64   PlatformArchType = "AMD64"
	PlatformArchType_ARM64   PlatformArchType = "ARM64"
)

type PlatformFamilyType string

const (
	PlatformFamilyType_Unknown      PlatformFamilyType = "Unknown"
	PlatformFamilyType_Linux        PlatformFamilyType = "Linux"
	PlatformFamilyType_Linux_Ubuntu PlatformFamilyType = "Ubuntu"
	PlatformFamilyType_Linux_CentOS PlatformFamilyType = "CentOS"
	PlatformFamilyType_OSX          PlatformFamilyType = "OSX"
	PlatformFamilyType_Windows      PlatformFamilyType = "Windows"
)

var PlatformOS PlatformOSType
var PlatformArch PlatformArchType
var PlatformFamily PlatformFamilyType

func init() {
	hostInfo, _ := host.Info()
	// os
	switch hostInfo.OS {
	case "linux":
		PlatformOS = PlatformOSType_Linux
		switch hostInfo.PlatformFamily {
		//case "debian":
		//	PlatformFamily = PlatformFamilyType_Ubuntu
		case "rhel":
			PlatformFamily = PlatformFamilyType_Linux_CentOS
		default:
			log.Warn("Unknown Family: ", hostInfo.PlatformFamily)
			PlatformFamily = PlatformFamilyType_Linux
		}
	case "windows":
		PlatformOS = PlatformOSType_Windows
		PlatformFamily = PlatformFamilyType_Windows
	case "darwin":
		PlatformOS = PlatformOSType_OSX
		PlatformFamily = PlatformFamilyType_OSX
	default:
		log.Warn("Unknown OS: ", hostInfo.OS)
		PlatformOS = PlatformOSType_Unknown
	}
	// arch
	switch hostInfo.KernelArch {
	case "x86_64":
		PlatformArch = PlatformArchType_AMD64
	case "arm64":
		PlatformArch = PlatformArchType_ARM64
	default:
		log.Warn("Unknown Arch: ", hostInfo.KernelArch)
		PlatformArch = PlatformArchType_Unknown
	}
	// distro
}
