package slackbot

import (
	"reflect"
	"testing"

	"github.com/nlopes/slack"

	v1 "k8s.io/api/core/v1"

	jenkinsv1 "github.com/jenkins-x/jx/pkg/apis/jenkins.io/v1"

	slackapp "github.com/jenkins-x-labs/slack/pkg/apis/slack/v1alpha1"
	slackappapi "github.com/jenkins-x-labs/slack/pkg/apis/slack/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreateSlackBot(t *testing.T) {
	secretName := "test_secret"
	testToken := "123abc"
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Data: map[string][]byte{
			"token": []byte(testToken),
		},
	}
	fakeclient := fake.NewSimpleClientset(secret)

	clients := &Clients{
		KubeClient: fakeclient,
	}

	tests := []struct {
		name     string
		slackBot *slackapp.SlackBot
		want     *slack.Client
		wantErr  bool
	}{
		{name: "good_slack_bot", slackBot: getSlackBot(secretName), want: slack.New(testToken), wantErr: false},
		{name: "good_slack_bot", slackBot: getSlackBot("does_not_exist"), want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateSlackBot(clients, tt.slackBot)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSlackBot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("CreateSlackBot() got = nil")
			}
			if !tt.wantErr && !reflect.DeepEqual(got.SlackClient, tt.want) {
				t.Errorf("CreateSlackBot() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getSlackBot(secretName string) *slackappapi.SlackBot {
	return &slackappapi.SlackBot{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test_slack_bot",
		},
		Spec: slackappapi.SlackBotSpec{
			TokenReference: jenkinsv1.ResourceReference{
				Kind: "Secret",
				Name: secretName,
			},
		},
	}
}
