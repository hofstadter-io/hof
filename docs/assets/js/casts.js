function addCasts() {
	var casts = ["create", "module"]

	casts.forEach(c => {
		var cast = document.getElementById(c + "-cast");

		console.log(cast);

		if (!!cast) {
			AsciinemaPlayer.create(`/casts/${c}.cast`, cast, {
				autoPlay: true,
				loop: true,
			});
		}
	})
}

addCasts()
