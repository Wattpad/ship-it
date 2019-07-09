package ecrconsumer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseSHA(t *testing.T) {
	assert.Equal(t, parseSHA("sha256:the-tag"), "the-tag")
	assert.Equal(t, parseSHA("malformed"), "")
}

func TestParseMessage(t *testing.T) {
	inputJSON := `
{
	"detail": {
		"eventTime": "2019-06-28T17:42:49Z",
		"responseElements": {
			"image": {
				"repositoryName": "writer-dashboard",
				"imageManifest": "{\n   \"schemaVersion\": 2,\n   \"mediaType\": \"application/vnd.docker.distribution.manifest.v2+json\",\n   \"config\": {\n      \"mediaType\": \"application/vnd.docker.container.image.v1+json\",\n      \"size\": 2316,\n      \"digest\": \"sha256:5a1b9e65745aadb59af6669eafb6057f238f118621771bc4000698741847ee72\"\n   },\n   \"layers\": [\n      {\n         \"mediaType\": \"application/vnd.docker.image.rootfs.diff.tar.gzip\",\n         \"size\": 45339350,\n         \"digest\": \"sha256:6f2f362378c5a6fd915d96d11dda1e0223ccf213bf121ace56ae0f6616ea1dc8\"\n      },\n      {\n         \"mediaType\": \"application/vnd.docker.image.rootfs.diff.tar.gzip\",\n         \"size\": 3672402,\n         \"digest\": \"sha256:82a2520544ca213d45bd63751ea90f529e1993b98cddfbd63e806091c1ba625c\"\n      },\n      {\n         \"mediaType\": \"application/vnd.docker.image.rootfs.diff.tar.gzip\",\n         \"size\": 4436039,\n         \"digest\": \"sha256:7f97b11819ed7e1453daf98f2125d921dbd38e517e0092d0a665e2569d28658e\"\n      }\n   ]\n}",
				"registryId": "723255503624",
				"imageId": {
					"imageDigest": "sha256:d8ac457c16d1e20172d771e5244c75bdcfb93e33e49fa85eb7e25d75ca7c74c1",
					"imageTag": "latest"
				}
			}
		}
	}
}
`

	deployTime, err := time.Parse(time.RFC3339, "2019-06-28T17:42:49Z")
	assert.NoError(t, err)

	expectedMessage := SQSMessage{
		Detail: Detail{
			EventTime: deployTime,
			Response: ResponseElements{
				Image: ImageData{
					RepositoryName: "writer-dashboard",
					ID: ImageID{
						Digest: "sha256:d8ac457c16d1e20172d771e5244c75bdcfb93e33e49fa85eb7e25d75ca7c74c1",
					},
				},
			},
		},
	}

	inputMessage, err := parseMsg(inputJSON)
	assert.NoError(t, err)
	assert.Exactly(t, expectedMessage, *inputMessage)
}

func TestMakeImage(t *testing.T) {
	assert.Exactly(t, Image{
		Registry: "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "ship-it",
		Tag: "shipped",
	}, makeImage("ship-it", "shipped"))
}