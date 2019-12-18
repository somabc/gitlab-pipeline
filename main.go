package main

import (
	"flag"
	"fmt"
	"gitlab.com/schehata/gitlab-pipeline-trigger/gitlab"
	"log"
	"os"
	"strings"
	"time"
)

// Configuration models the configurations needed to run the
// triggers, and/or extra parameters that might be passed to a pipeline
type Configuration struct {
	Host          string
	APIToken      string
	TriggerToken  string
	TargetBranch  string
	ProjectID     int
	WaitForStatus bool
	ExtraParams   map[string]string
}

func parseFlags() (conf Configuration) {
	flag.StringVar(&conf.Host, "h", "https://gitlab.com", "Gitlab Host URL (default: https://gitlab.com)")
	flag.StringVar(&conf.APIToken, "a", "", "API Token")
	flag.StringVar(&conf.TriggerToken, "t", "", "Trigger Token")
	flag.StringVar(&conf.TargetBranch, "b", "master", "Target Branch (default: master)")
	flag.IntVar(&conf.ProjectID, "p", 0, "Project ID")
	flag.BoolVar(&conf.WaitForStatus, "w", false, "Wait for pipeline to finish")
	conf.ExtraParams = make(map[string]string)

	flag.Parse()
	extraParams := flag.Args()
	// loop over extra parameters and add them as variables[key]=value to Post Form
	for _, param := range extraParams {
		splitted := strings.Split(param, "=")
		if len(splitted) < 2 {
			err := fmt.Errorf("failed to parse argument %s, use syntax: key=value", param)
			fmt.Println(err)
			os.Exit(1)
		}
		key := fmt.Sprintf("variables[%s]", splitted[0])
		conf.ExtraParams[key] = splitted[1]
	}

	return
}

func isConfigurationValid(conf Configuration) (err error) {
	if conf.APIToken == "" {
		err = fmt.Errorf("please set the API Token, pass -a $API_TOKEN to the command line")
	}
	if conf.TriggerToken == "" {
		err = fmt.Errorf("please set the Trigger Token, pass -t $TRIGGER_TOKEN to the command line")
	}
	if conf.ProjectID == 0 {
		err = fmt.Errorf("please set the Project ID, pass -p $PROJECT_ID to the command line")
	}

	return
}

func main() {

	conf := parseFlags()
	err := isConfigurationValid(conf)
	if err != nil {
		log.Fatal(err)
	}

	client := gitlab.New(conf.Host, conf.APIToken)
	pipeline, err := client.TriggerPipeline(conf.ProjectID, conf.TriggerToken, conf.TargetBranch, conf.ExtraParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Triggered New Pipeline: ", pipeline.WebURL)
	fmt.Printf("By user: %s(%s)\n", pipeline.User.Name, pipeline.User.Username)
	fmt.Println("Created At:", pipeline.CreatedAt)

	if !conf.WaitForStatus {
		return
	}

	finished := false
	fmt.Println("Getting ready to check pipeline status every 5 seconds...")
	for finished == false {
		pipeline, err = client.CheckPipelineStatus(pipeline.ID, conf.ProjectID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(".")
		switch pipeline.Status {
		case "failed":
			finished = true
			log.Printf("pipeline %d failed\n", pipeline.ID)
			break
		case "cancled":
			finished = true
			log.Printf("pipeline %d cancled\n", pipeline.ID)
			break
		case "success":
			finished = true
			log.Printf("pipeline %d succeeded", pipeline.ID)
			break
		default:
			time.Sleep(5 * time.Second)
		}
	}

	pipeline, err = client.CheckPipelineStatus(pipeline.ID, conf.ProjectID)
	if err != nil {
		log.Fatal(err)
	}
}
