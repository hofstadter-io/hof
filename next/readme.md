# next hof docs

This is the code for the "next" version of our docs using NextJS.

### local 

```
# setup
npm i

# run dev
npm run dev

# local build
npm run build
npm run start
```


### dagger

We use [Dagger](https://dagger.io) for CI
build, test, publish, and deploy.

```
# build
npm run dagger -- build

# test (tbd, broken link checker for sure)

# publish
npm run dagger -- publish --version <tag>

# deploy next
npm run dagger -- deploy --version <tag> --name hof-next-docs

# deploy prod
npm run dagger -- deploy --version <tag> --name hof-prod
```
