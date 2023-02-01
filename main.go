package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nedpals/supabase-go"
	"github.com/pigeatgarlic/VirtualOS-daemon/child-process"
	"github.com/pigeatgarlic/VirtualOS-daemon/system"
	"github.com/pigeatgarlic/oauth2l"
)



type Daemon struct {
	childprocess *childprocess.ChildProcesses
	shutdown     chan bool
}

type Cred struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Supbase struct {
	URL string `json:"url"`
	KEY string `json:"key"`
}

func main() {
	daemon := Daemon{
		shutdown:               make(chan bool),
	}

	result,err := os.ReadFile("./supabase.json");
	if err != nil {
		fmt.Printf("%s",err.Error());
		return
	}

	var conf Supbase
	err = json.Unmarshal(result,&conf)
	if err != nil {
		fmt.Printf("%s",err.Error());
		return
	}


	if err != nil {
		fmt.Printf("wrong cred %s",err.Error());
	}

	var cred Cred
	var dofetch bool
	var account_id string

	f, err := os.OpenFile("./cache.secret.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		dofetch = true;
	} else {
		bytes := make([]byte,1000)
		count,_ := f.Read(bytes)	
		err = json.Unmarshal(bytes[:count],&cred)
		dofetch = (err != nil);
	}

	if dofetch {
		sysinf := system.GetInfor()
		account,err := oauth2l.StartAuth(sysinf)
		if err != nil {
			fmt.Printf("%s\n",err.Error())
			return;
		}

		cred.Username = account.Username 
		cred.Password = account.Password 
		bytes,_ := json.Marshal(cred)
		if err != nil {
			fmt.Printf("%s",err.Error())
		}

		_,err = f.Write(bytes);
		if err != nil {
			fmt.Printf("%s",err.Error())
		}

		if err := f.Close(); err != nil {
			fmt.Printf("%s",err.Error())
		}
	}


	supabase_client := supabase.CreateClient(conf.URL,conf.KEY);
	detail,err := supabase_client.Auth.SignIn(context.Background(),supabase.UserCredentials{
		Email: cred.Username,
		Password: cred.Password,
	})
	if err != nil {
		fmt.Printf("%s",err.Error())
		return
	}

	fmt.Printf("signin with username %s\n",detail.User.Email);
	account_id = detail.User.ID
	defer func ()  {
		supabase_client.Auth.SignOut(context.Background(),detail.AccessToken);
	}()

	var data interface{}
	err = supabase_client.DB.From("worker_profile").Select("*").Eq("account_id",account_id).Execute(&data)
	if err != nil {
		fmt.Printf("%s",err.Error())
		return;
	} 

	val,_ := json.MarshalIndent(data,"","  ")
	fmt.Printf("registered new worker: %s\n",string(val))

	go func() {
		for {
			var result interface{}
			if res := supabase_client.DB.From("worker_profile").Update(struct{ 
				LastUpdate string `json:"last_update"`
			}{
				LastUpdate: time.Now().Format(time.RFC3339Nano),
			}).Eq("account_id",account_id).Execute(&result); res != nil {
				fmt.Printf("error ping %s\n",res.Error())
			}

			time.Sleep(time.Second)
		}
	}()

	go func() {
		var raw_worker_profile interface{}
		for {
			time.Sleep(time.Second)
			err = supabase_client.DB.From("worker_profile").Select("metadata").Eq("account_id",account_id).Execute(&raw_worker_profile)
			if err != nil {
				fmt.Printf("error sync %s\n",err.Error())
			}
			// last_ud := raw_worker_profile.([]interface{})[0].(map[string]interface{})["metadata"];
		}
	}()





	<-daemon.shutdown
}
