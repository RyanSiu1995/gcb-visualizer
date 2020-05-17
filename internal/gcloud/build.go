package gcloud

import (
	"context"

	cloudbuildClient "cloud.google.com/go/cloudbuild/apiv1"
	cloudbuildpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

// GetRemoteBuild returns the remote build object
func GetRemoteBuild(project string, buildID string) (*cloudbuildpb.Build, error) {
	ctx := context.Background()

	client, err := cloudbuildClient.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	req := &cloudbuildpb.GetBuildRequest{
		ProjectId: project,
		Id:        buildID,
	}

	resp, err := client.GetBuild(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
