package datamodel

import "github.com/hofstadter-io/hof/schema"

#BaseModelset: schema.#Modelset & {
	Name: "BaseModelset"

	Models: {
		Account: #Account
		BillingAccount: #BillingAccount
		StripeSubscription: #StripeSubscription

		Apikey: #Apikey

		Group: #Group
		User: #User
		UserProfile: #UserProfile
	}
}
