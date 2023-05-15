package prompts

BlueSky: [{
	role: "user"
	content: ##"""
		Why is the sky blue?
		"""##
}]

Planets: {

	context: "My name is Miles. You are an astronomer, knowledgeable about the solar system."
	examples: [{
		input:  "How many moons does Mars have?"
		output: "The planet Mars has two moons, Phobos and Deimos."
	}, {
		input:  "What is the largest planet in the solar system?"
		output: "Jupiter is the largest planet in our solar system."
	}]

	messages: [{
		role:    "user"
		content: "How many planets are there in the solar system?"
	}, {
		role:    "assistant"
		content: "There are eight planets in our solar system: Mercury, Venus, Earth, Mars, Jupiter, Saturn, Uranus, and Neptune."
	}, {
		role:    "user"
		content: "What happened to Pluto?"
	}]

}
