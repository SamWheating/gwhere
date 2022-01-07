package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/storage/v1"
)

// Get the project Number from the bucket Metadata
// For whatever reason this isn't surfaced in the storage API, so we have to use the Resource Manager API
func getProjectNumber(ctx context.Context, bucketName string) string {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	service, err := storage.NewService(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bucketService := storage.NewBucketsService(service)

	resp, err := bucketService.Get(bucketName).Do()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	projectNumber := strconv.FormatUint(resp.ProjectNumber, 10)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return projectNumber
}

// Returns the project name from the numeric project ID
func getProjectIDFromNumber(ctx context.Context, projectNumber string) string {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req := cloudresourcemanagerService.Projects.List()

	project := req.Filter("projectNumber=" + projectNumber)
	resp, err := project.Do()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(resp.Projects) < 1 {
		fmt.Printf("No project found for project Number %s\n", projectNumber)
		os.Exit(1)
	}
	return resp.Projects[0].ProjectId
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("gwhere - find the GCP project ID associated with a cloud storage bucket.")
		fmt.Println("usage: gwhere <bucket>")
		os.Exit(0)
	}

	if len(os.Args) != 2 {
		fmt.Printf("Exactly 1 argument required, %d provided\n", len(os.Args)-1)
		os.Exit(0)
	}

	ctx := context.Background()

	bucketName := os.Args[1]
	bucketName = strings.TrimPrefix(bucketName, "gs://")

	projectNumber := getProjectNumber(ctx, bucketName)
	projectID := getProjectIDFromNumber(ctx, projectNumber)

	fmt.Println(projectID)
	os.Exit(0)
}
