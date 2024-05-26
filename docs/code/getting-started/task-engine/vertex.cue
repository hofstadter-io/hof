package examples

import "strings"

vertex_chat: {
  @flow() // define a flow

  steps: {

    // task: get auth from external command
    gcp: {
      @task(os.Exec)
      cmd: ["gcloud", "auth", "print-access-token"]
      stdout: string
      key:    strings.TrimSpace(stdout)
    }

    // task: api call via reusable task
    call: _gemini & {
      apikey: gcp.key

      msg: "What is the CUE language?"

      resp: body: _
    }

    // task: print text to std output
    out: {
      @task(os.Stdout)
      text: call.final.text
    }

  }
}

// reusable task
_gemini: {
  @task(api.Call)

  model: string | *"gemini-1.0-pro-002:generateContent"

  msg: string
  apikey: string
  prompt: string | *"You are a model which is direct and concise when responding."

  req: {
    host: "https://us-central1-aiplatform.googleapis.com"
    path: "/v1/projects/hof-io--develop/locations/us-central1/publishers/google/models/\(model)"
    headers: {
      "Content-Type": "application/json"
      Authorization:  "Bearer \(apikey)"
    }
    data: {
      systemInstruction: {
        role: "MODEL"
        parts: [{
          text: prompt
        }]
      }

      contents: [{
        role: "USER"
        parts: [{
          text: msg
        }]
      }]
    }
    method: "POST"
  }

  resp: {
    body: _
  }

  // task-local ETL
  final: {
    cand: resp.body.candidates[0]
    text: cand.content.parts[0].text
  }
}