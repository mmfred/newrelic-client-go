// +build integration

package alerts

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationChannel(t *testing.T) {
	t.Parallel()

	var (
		testChannelEmail = Channel{
			Name: "integration-test-email",
			Type: "email",
			Configuration: ChannelConfiguration{
				Recipients:            "devtoolkittest@newrelic.com",
				IncludeJSONAttachment: "true",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelOpsGenie = Channel{
			Name: "integration-test-opsgenie",
			Type: "opsgenie",
			Configuration: ChannelConfiguration{
				APIKey:     "abc123",
				Teams:      "dev-toolkit",
				Tags:       "tag1,tag2",
				Recipients: "devtoolkittest@newrelic.com",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelSlack = Channel{
			Name: "integration-test-slack",
			Type: "slack",
			Configuration: ChannelConfiguration{
				URL:     "https://example-org.slack.com",
				Channel: "test-channel",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelVictorops = Channel{
			Name: "integration-test-victorops",
			Type: "victorops",
			Configuration: ChannelConfiguration{
				Key:      "abc123",
				RouteKey: "/route-name",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelWebhook = Channel{
			Name: "integration-test-webhook",
			Type: "webhook",
			Configuration: ChannelConfiguration{
				BaseURL:     "https://test.com",
				PayloadType: "application/json",
				Headers: MapStringInterface{
					"x-test-header": "test-header",
				},
				Payload: MapStringInterface{
					"account_id": "123",
				},
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelWebhookEmptyHeadersAndPayload = Channel{
			Name: "integration-test-webhook-empty-headers-and-payload",
			Type: "webhook",
			Configuration: ChannelConfiguration{
				BaseURL: "https://test.com",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelWebhookWeirdHeadersAndPayload = Channel{
			Name: "integration-test-webhook-weird-headers-and-payload",
			Type: "webhook",
			Configuration: ChannelConfiguration{
				BaseURL: "https://test.com",
				Headers: MapStringInterface{
					"": "",
				},
				Payload: MapStringInterface{
					"": "",
				},
				PayloadType: "application/json",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		// Currently the v2 API has minimal validation on the data
		// structure for Headers and Payload, so we need to test
		// as many scenarios as possible.
		testChannelWebhookComplexHeadersPayload = Channel{
			Name: "integration-test-webhook",
			Type: "webhook",
			Configuration: ChannelConfiguration{
				BaseURL:     "https://test.com",
				PayloadType: "application/json",
				Headers: MapStringInterface{
					"x-test-header": "test-header",
					"object": map[string]interface{}{
						"key": "value",
						"nestedObject": map[string]interface{}{
							"k": "v",
						},
					},
				},
				Payload: MapStringInterface{
					"account_id": "123",
					"array":      []interface{}{"string", 2},
					"object": map[string]interface{}{
						"key": "value",
						"nestedObject": map[string]interface{}{
							"k": "v",
						},
					},
				},
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		channels = []Channel{
			testChannelEmail,
			testChannelOpsGenie,
			testChannelSlack,
			testChannelVictorops,
			testChannelWebhook,
			testChannelWebhookEmptyHeadersAndPayload,
			testChannelWebhookWeirdHeadersAndPayload,
			testChannelWebhookComplexHeadersPayload,
		}
	)

	client := newIntegrationTestClient(t)

	for _, channel := range channels {
		// Test: Create
		created, err := client.CreateChannel(channel)

		require.NoError(t, err)
		require.NotNil(t, created)

		// Test: Read
		read, err := client.GetChannel(created.ID)

		require.NoError(t, err)
		require.NotNil(t, read)

		// Test: Delete
		deleted, err := client.DeleteChannel(read.ID)

		require.NoError(t, err)
		require.NotNil(t, deleted)
	}
}

func TestListChannels(t *testing.T) {
	fmt.Printf("\n\nAPI KEY: %+v \n\n", os.Getenv("NEWRELIC_API_KEY"))

	client := newIntegrationTestClient(t)
	count := 400

	for i := 0; i < count; i++ {
		channel := Channel{
			Name: fmt.Sprintf("test-channel-%d", i),
			Type: "email",
			Configuration: ChannelConfiguration{
				Recipients:            "sblue@newrelic.com",
				IncludeJSONAttachment: "true",
			},
		}

		ch, err := client.CreateChannel(channel)

		if err != nil {
			t.Log(err)
		} else {
			fmt.Printf("\nCreated channel: %+v\n", ch.Name)
		}
	}

	channels, _ := client.ListChannels()

	fmt.Printf("\nChannels count: %+v\n", len(channels))

	// for _, ch := range channels {
	// 	match := strings.HasPrefix(ch.Name, "test-channel")

	// 	if match == true {
	// 		c, err := client.DeleteChannel(ch.ID)

	// 		if err != nil {
	// 			t.Log(err)
	// 		}

	// 		t.Logf("Deleted channel: %+v", c.Name)
	// 	}
	// }
}

func TestDeleteChannels(t *testing.T) {
	client := newIntegrationTestClient(t)

	channels, _ := client.ListChannels()

	for _, ch := range channels {
		match := strings.HasPrefix(ch.Name, "test-channel")

		if match == true {
			_, err := client.DeleteChannel(ch.ID)

			t.Logf("\nDeleting channel: %+v\n", ch.Name)

			if err != nil {
				t.Log(err)
			}
		}
	}
}
