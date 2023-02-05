package app





type WorkerAppManifest struct {
	ExecuteArgs     []string 				`json:"args"`
	Env             map[string]string 		`json:"envs"`
}


type WorkerApplication struct {
	ID               int			   `json:"id"`

	SessionID        int  			   `json:"session_id"`
	App              int  			   `json:"app"`

	Manifest  		 WorkerAppManifest `json:"manifest"`
}