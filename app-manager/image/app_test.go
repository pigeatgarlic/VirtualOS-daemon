package image

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRm(t *testing.T) {
	app := &ApplicationImage{
		Id: 0,
		Name: "test",

		Manifest: Manifest{
			Needbuild: true,
			BuildPocedure: BuildProcedure{
				Dependencies: map[string][]string{
					"git.exe" : { "version" },
					"dotnet.exe": { "--version"},
					"go.exe": { "version"},
				},
				Commands: map[string]Command{
					"check": {
						Args: []string{"powershell.exe","Remove-Item","hellogitworld","-Recurse:$true","-Force","-Confirm:$false" },
						DesiredExitcode: -1,
					},
					"clone": {
						Args: []string{"git","clone","https://github.com/githubtraining/hellogitworld"},
						DesiredExitcode: -1,
					},
				},
				Order: []string{"check"},
			},
		},
	}

	err := app.Apply(nil)
	if err != nil {
		t.Error(err)
		return
	}
	bytes,_ := json.MarshalIndent(app,"","   ")
	fmt.Printf("tested manifest: \n%s",bytes)
}
func TestApp(t *testing.T) {
	app := &ApplicationImage{
		Id: 0,
		Name: "test",

		Manifest: Manifest{
			Needbuild: true,
			BuildPocedure: BuildProcedure{
				Dependencies: map[string][]string{
					"git.exe" : { "version" },
					"dotnet.exe": { "--version"},
					"go.exe": { "version"},
				},
				Commands: map[string]Command{
					"check": {
						Args: []string{"powershell.exe","Remove-Item","hellogitworld","-Recurse:$true","-Force","-Confirm:$false" },
						DesiredExitcode: -1,
					},
					"clone": {
						Args: []string{"git","clone","https://github.com/githubtraining/hellogitworld"},
						DesiredExitcode: -1,
					},
				},
				Order: []string{"check","clone"},
			},
		},
	}

	err := app.Apply(nil)
	if err != nil {
		t.Error(err)
		return
	}
	bytes,_ := json.MarshalIndent(app,"","   ")
	fmt.Printf("tested manifest: \n%s",bytes)
}