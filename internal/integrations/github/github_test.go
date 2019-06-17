package github

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestGetTravisCIBuildURLForRef_Success(t *testing.T) {
	defer gock.Off()
	ctx := context.Background()

	mockResponse, _ := ioutil.ReadFile("testdata/check_runs_response_success.json")
	gock.New("https://api.github.com").
		Get("/repos/Wattpad/highlander/commits/2c76895cdb9f3ff5100ecf93a7a6c6747aaeda8c/check-runs").
		Reply(200).
		JSON(mockResponse)

	github := New(ctx, "Wattpad", "fake-access-token")
	url, _ := github.GetTravisCIBuildURLForRef(ctx, "highlander", "2c76895cdb9f3ff5100ecf93a7a6c6747aaeda8c")

	assert.Equal(t, "https://travis-ci.com/Wattpad/highlander/builds/115827260", url)
}

func TestGetTravisCIBuildURLForRef_Empty(t *testing.T) {
	defer gock.Off()
	ctx := context.Background()

	mockResponse, _ := ioutil.ReadFile("testdata/check_runs_response_empty.json")
	gock.New("https://api.github.com").
		Get("/repos/Wattpad/highlander/commits/2c76895cdb9f3ff5100ecf93a7a6c6747aaeda8c/check-runs").
		Reply(200).
		JSON(mockResponse)

	github := New(ctx, "Wattpad", "fake-access-token")
	_, err := github.GetTravisCIBuildURLForRef(ctx, "highlander", "2c76895cdb9f3ff5100ecf93a7a6c6747aaeda8c")

	assert.EqualError(t, errTravisCIBuildNotFound, err.Error())
}

func TestGetTravisCIBuildURLForRef_500(t *testing.T) {
	defer gock.Off()
	ctx := context.Background()

	gock.New("https://api.github.com").
		Get("/repos/Wattpad/highlander/commits/2c76895cdb9f3ff5100ecf93a7a6c6747aaeda8c/check-runs").
		Reply(500)

	github := New(ctx, "Wattpad", "fake-access-token")
	_, err := github.GetTravisCIBuildURLForRef(ctx, "highlander", "2c76895cdb9f3ff5100ecf93a7a6c6747aaeda8c")

	assert.NotEmpty(t, err)
}
