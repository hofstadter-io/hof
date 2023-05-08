package chat

var pretextString = `
TYPES: [string, int, bool, float, uuid, datetime, email, url]
RELATIONS: [belongs-to, has-one, has-many, many-to-many]

SCHEMA: """
{
  "Datamodel": {
    "Name": "<application-name>",
    "Models": {
      "<model-name>": {
        "<field-name>": "<field-type>",
        "<field-name>": "<field-type>",
        "$relations": {
          "<relation-name>": {
            "name": "<relation-name>",
            "type": "<relation-type>",
            "model": "<model-name>"
          },
          "<relation-name>": {
            ...
          }
        }
      },
      "<model-name>": {
        ...
      }
    }
  }
}
"""

Your task is to modify the original JSON object according to the instructions.
You should only output the relevant changes to the JSON object between a pair of triple backticks with no explanation or other words.
The JSON should conform to the SCHEMA and follow the GUIDELINES.
Once you have generated the JSON, reread the GUILDLINES and make any corrections.

Use the following GUIDELINES when performing the task:
- An application is described by a Datamodel and normally has users and many other models.
- The Datamodel should have a Name field set to the application name.
- The Datamodel is a loose representation of a SQL database, models are tables and fields are columns.
- The Datamodel is composed of Models and Models are composed of Fields.
- Models should appear directly within the Datamodel and NEVER within another Model.
- Do not place any Models within User. You may only modify User to add Relations.
- You are allowed to make assumptions about the need for new models or fields if the instructions seem to imply their need.
- You should try to keep the number of models concise and not introduce unnecessary duplication of information.
- Field values must come from the TYPES list.
- The common database fields are id, created_at, updated_at, and deleted_at.
- When adding a new model, include the common database fields, unless instructed otherwise.
- If the instructions do not specify the field type, you should make your best guess.
- If a field can be calculated by a SQL query on another table, don't add it.
- Models can have relations between them. If you make a relation, there must be a model for both sides.
- If a user has something, this implies a new Model and Relation. It is up to you to determine the correct relation type.
- The relations type must come from the RELATIONS list.
- "has-many" and "many-to-many" relations should be named as the plural of the model they refer to.
- "many-to-many" relations require an extra model to hold the linking information.
- "$relations" is a special field on a model, represented as an object with the fields "name", "type", and "model".
- If there are any syntax errors in the original JSON, you should also fix those.
`

var promptTemplate = `
JSON: """
{{ .code }}
"""
`

/*
INSTRUCTIONS: """
{{ .inst }}
"""
`
*/

var initialCode = `
{
  "Datamodel": {
    "Name": "<application-name>",
    "Models": {
      "User": {
        "id": "uuid",
        "created_at": "datetime",
        "updated_at": "datetime",
        "deleted_at": "datetime",
        "email": "email",
        "active": "bool",
        "verified": "bool"
      }
    }
  }
}
`
