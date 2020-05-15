package schema

#Modelsets: [N=string]: #Modelset & { Name: N, ... }
#Modelset: {
  Name: string

  Tags: [...string]

  Models: #Model
  Modelsets: #Modelset

  Stores: #Store
}
