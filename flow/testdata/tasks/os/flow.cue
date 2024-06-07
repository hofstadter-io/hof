tasks: {
  @flow()
	r: { filename: "in.txt", contents: string } @task(os.ReadFile)
}

