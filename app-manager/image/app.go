package image

import (
	"fmt"
	"os/exec"
)
type Command struct{
	Args []string 							`json:"args"`
	DesiredExitcode int  					`json:"desired_exitcode"`
}

type BuildProcedure struct {
	Order    	   []string             	`json:"order"`
	Commands 	   map[string]Command 		`json:"commands"`

	Dependencies   map[string][]string  	`json:"dependencies"`
}
func (procedure *BuildProcedure) Build() error {
	for dep,test := range procedure.Dependencies {
		bytes,err := exec.Command(dep,test...).Output()
		if err != nil {
			return fmt.Errorf("error testing dependencies %s : %s \n %s",dep,err.Error(),bytes)
		} else {
			fmt.Printf("%s dependencies installed\n",dep)
			fmt.Println(string(bytes))
		}
	}

	for _,commandName := range procedure.Order {
		command := procedure.Commands[commandName]
		bytes,err := exec.Command(command.Args[0],command.Args[1:]...).Output()
		if err != nil && command.DesiredExitcode != -1{
			return fmt.Errorf("error building application %s",err.Error())
		} else {
			fmt.Printf("%s dependencies installed\n",commandName)
			fmt.Println(string(bytes))
		}
	}

	return nil
}







type Manifest struct {
	Needbuild     	bool           			`json:"need_build"`
	BuildPocedure 	BuildProcedure 			`json:"build_procedure"`

	Version         bool           			`json:"version"`
	UpdatePocedure 	BuildProcedure 			`json:"update_procedure"`

	Needrmv         bool           			`json:"need_rmv"`
	RmvPocedure 	BuildProcedure 			`json:"rmv_procedure"`

	ContainerFolder string   				`json:"folder"`
	Executable      string   				`json:"executable"`

	Execute struct {
		ExecuteArgs     []string 				`json:"args"`
		Env             map[string]string 		`json:"envs"`
	}`json:"execute"`
}


type InvocationMode string
const (
	ModeSingleton     = "SINGLETON"
	ModeUserAppDaemon = "DAEMON"
	ModeUserApp		  = "APPLICATION"
)

type ApplicationImage struct {
	Id   int    							`json:"id"`
	Name string 							`json:"name"`
	Mode InvocationMode						`json:"mode"`

	Manifest Manifest    					`json:"manifest"`
	Metadata interface{} 					`json:"metadata"`
}

func (desired *ApplicationImage) Apply(now *ApplicationImage) (err error) {
	if desired == nil && now == nil {
		return nil
	} else if desired == nil {
		if now.Manifest.Needrmv {
			return now.Manifest.RmvPocedure.Build()
		}
	} else if now == nil {
		if desired.Manifest.Needbuild {
			return desired.Manifest.BuildPocedure.Build()
		}
	}

	if desired.Manifest.Version != now.Manifest.Version {
		err = desired.Manifest.UpdatePocedure.Build()
		if err != nil {
			return 
		}
	}





	return nil
}