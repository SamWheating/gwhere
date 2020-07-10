package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iterator"
	"os"
	"strings"
)

// Looks at a bucket's ACLs to find where a bucket might be.
func getCandidateProjects(bucketName string) map[string]bool {

	projectNumbers := map[string]bool{}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	attrs, err := client.Bucket(bucketName).Attrs(ctx)

	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	for _, rule := range attrs.ACL {
		if rule.ProjectTeam != nil {
			projectNumbers[rule.ProjectTeam.ProjectNumber] = true
		}
	}

	return projectNumbers
}

// Returns the project name from the numeric project ID
func getProjectIDsFromNumber(projectIDs map[string]bool) map[string]bool {

	projectNames := map[string]bool{}

	ctx := context.Background()

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req := cloudresourcemanagerService.Projects.List()

	for projectID := range projectIDs {
		project := req.Filter("projectNumber=" + projectID)
		resp, err := project.Do()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(resp.Projects) > 0 {
			projectNames[resp.Projects[0].ProjectId] = true
		}
	}
	return projectNames
}

func confirmBucketExists(bucketName string, projectIDs map[string]bool) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for projectID := range projectIDs {
		it := client.Buckets(ctx, projectID)
		it.Prefix = bucketName
		for {
			bucketAttrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if bucketAttrs.Name == bucketName {
				fmt.Println(projectID)
				os.Exit(0)
			}
		}
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Exactly 1 argument required, %d provided\n", len(os.Args)-1)
		os.Exit(1)
	}

	bucketName := os.Args[1]
	bucketName = strings.TrimPrefix(bucketName, "gs://")

	projectNumbers := getCandidateProjects(bucketName)
	projectIDs := getProjectIDsFromNumber(projectNumbers)

	confirmBucketExists(bucketName, projectIDs)

	fmt.Printf("Could not find project for bucket gs://%s", bucketName)
	os.Exit(1)

}
