// Resources are pretty simple
#Resources: [string]: #Resource
#Resource: {
	Model:  #Model
	Name:   Model.Name
	Routes: #Routes
}

// We map from Datamodle to Resource
#DatamodelToResources: {
	Datamodel: #Datamodel
	Resources: #Resources & {
		for n, M in Datamodel.Models {
			// sub-value with same name as label as the model
			"\(n)": {
				// Same model and name
				Model: M
				Name:  M.Name

				// The default CRUD routes
				Routes: [{
					Name:   "\(M.Name)Create"
					Path:   ""
					Method: "POST"
				}, {
					Name: "\(M.Name)Read"
					Path: ""
					Params: ["\(strings.ToLower(M.Index))"]
					Method: "GET"
				}, {
					Name:   "\(M.Name)Update"
					Path:   ""
					Method: "PATCH"
				}, {
					Name: "\(M.Name)Delete"
					Path: ""
					Params: ["\(strings.ToLower(M.Index))"]
					Method: "DELETE"
				}, ...] // left open so you can add custom routes
			}
		}
	}
}
