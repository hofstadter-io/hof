package common

import (
	"context"
	"fmt"

	"google.golang.org/adk/session"
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
)

type CreatePayload struct {
	User string `json:"user,omitempty"`

	Title   string `json:"title,omitempty"`
	Agent   string `json:"agent,omitempty"`
	Model   string `json:"model,omitempty"`
	EnvName string `json:"envName,omitempty"`

	Environ *environ.EnvironCreateOptions `json:"environ,omitempty"`
}

func SessionCreate(ctx context.Context, ar Runtime, payload CreatePayload) (session.Session, error) {
	// initial state
	initialState := make(map[string]any)
	if payload.Title != "" {
		initialState["title"] = payload.Title
	}

	initialState["agent"] = payload.Agent
	initialState["model"] = payload.Model
	initialState["envName"] = payload.EnvName

	pe := payload.Environ
	if pe == nil {
		pe = new(environ.EnvironCreateOptions)
	}

	// maybe attach an environment
	if pe.FromUri == "" && payload.EnvName != "" {
		// fmt.Println("searching for env:", payload.EnvName)
		for _, e := range ar.GetAgenticConfig().Environs {
			// fmt.Printf(" ? %#+v\n", e)
			if e.Name == payload.EnvName {
				if e.SpecValue.Exists() {
					pe.EnvValue = e.SpecValue
				} else if e.Spec.From != "" {
					pe.FromUri = "oci://" + e.Spec.From
				}
				break
			}
		}

		env := environ.Client()
		envUri, err := env.Create(pe)
		if err != nil {
			return nil, fmt.Errorf("in 'session.create' while creating env: %w", err)
		}
		// will these empty strings get deleted? (vs nil to delete, make sure delete is correct)
		initialState["initEnv"] = pe
		initialState["origEnv"] = string(envUri)
		initialState["currEnv"] = string(envUri)
	}

	// include any client level state
	// maps.Copy(initialState, c.State)

	// create our session
	resp, err := ar.GetSessionService().Create(ctx, &session.CreateRequest{
		AppName: ar.GetAppName(),
		UserID:  payload.User,
		State:   initialState,
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating session: %w", err)
	}
	return resp.Session, nil
}
