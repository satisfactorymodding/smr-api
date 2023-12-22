package gql

import (
	"github.com/satisfactorymodding/smr-api/generated"
)

type Resolver struct{}

func (r *Resolver) Mod() generated.ModResolver {
	return &modResolver{r}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) User() generated.UserResolver {
	return &userResolver{r}
}

func (r *Resolver) UserMod() generated.UserModResolver {
	return &userModResolver{r}
}

func (r *Resolver) VersionTarget() generated.VersionTargetResolver {
	return &versionTargetResolver{r}
}

func (r *Resolver) Version() generated.VersionResolver {
	return &versionResolver{r}
}

func (r *Resolver) GetMods() generated.GetModsResolver {
	return &getModsResolver{r}
}

func (r *Resolver) GetMyMods() generated.GetMyModsResolver {
	return &getMyModsResolver{r}
}

func (r *Resolver) GetVersions() generated.GetVersionsResolver {
	return &getVersionsResolver{r}
}

func (r *Resolver) GetMyVersions() generated.GetMyVersionsResolver {
	return &getMyVersionsResolver{r}
}

func (r *Resolver) Guide() generated.GuideResolver {
	return &guideResolver{r}
}

func (r *Resolver) GetGuides() generated.GetGuidesResolver {
	return &getGuidesResolver{r}
}

func (r *Resolver) GetSMLVersions() generated.GetSMLVersionsResolver {
	return &getSMLVersionsResolver{r}
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
