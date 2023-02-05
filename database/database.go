package database

import "github.com/pigeatgarlic/VirtualOS-daemon/app-manager/image"

type Database interface {
	GetWorkerProfile()

	GetApplicationImage() *map[string]image.ApplicationImage

	GetWorkerSession() 
}