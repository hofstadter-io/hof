# notes

once we start pushing this package out to more of the code base
we might want to consider using a hof/runtime to house the common
pieces across the major features, loading and enriching the CUE
as a whole, then being able to work with any of the sub-features

For now, we are only using this setup in datamodel as a POC.
Gen should be easy to update after that, flow with the rewrite
