package appmanager

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/nedpals/supabase-go"
	"github.com/pigeatgarlic/VirtualOS-daemon/app-manager/app"
	"github.com/pigeatgarlic/VirtualOS-daemon/app-manager/image"
)

type AppManager struct {
	images_mutex *sync.Mutex
	images map[string]image.ApplicationImage

	apps_mutex *sync.Mutex
	apps   map[int]app.WorkerApplication
}






func NewAppManger(supabase *supabase.Client) AppManager {
	apps := AppManager{
		apps: make(map[int]app.WorkerApplication),
		images: make(map[string]image.ApplicationImage),
		apps_mutex: &sync.Mutex{},
		images_mutex: &sync.Mutex{},
	}


	go func ()  {
		for {
			var result interface{}
			err := supabase.DB.From("application_image").Select("*").Execute(&result)
			if err != nil{
				fmt.Printf("error fetching appg img %s\n",err.Error())
			}

			for _,img_ := range result.([]interface{}) {
				img := image.ApplicationImage{}

				img_data,_ := json.Marshal(img_)
				err := json.Unmarshal(img_data,&img)
				if err != nil{
					fmt.Printf("error fetching appg img %s\n",err.Error())
					continue
				}
				
				apps.ApplyNewImage(&img)
			}
		}
	}()

	go func ()  {
		for {
			var result interface{}
			err := supabase.DB.From("worker_application").Select("*").Execute(result)
			if err != nil{
				fmt.Printf("error fetching appg img %s\n",err.Error())
			}

			for _,img_ := range result.([]interface{}) {
				app := app.WorkerApplication{}

				img_data,_ := json.Marshal(img_)
				err := json.Unmarshal(img_data,&app)
				if err != nil{
					fmt.Printf("error fetching appg img %s\n",err.Error())
					continue
				}

				apps.ApplyNewApp(&app)
			}
		}
	}()

	return apps
}



func (manager *AppManager) ApplyNewApp(app *app.WorkerApplication) {
	manager.apps_mutex.Lock()
	defer manager.apps_mutex.Unlock()

}
func (manager *AppManager) ApplyNewImage(img *image.ApplicationImage) {
	manager.apps_mutex.Lock()
	defer manager.apps_mutex.Unlock()
	// prev := apps.images[img.Name]
		

}