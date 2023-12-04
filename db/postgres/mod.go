package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/patrickmn/go-cache"
)

func GetModByID(ctx context.Context, modID string) *Mod {
	cacheKey := "GetModById_" + modID
	if mod, ok := dbCache.Get(cacheKey); ok {
		return mod.(*Mod)
	}

	return GetModByIDNoCache(ctx, modID)
}

func GetModByIDNoCache(ctx context.Context, modID string) *Mod {
	var mod Mod
	DBCtx(ctx).Preload("Tags").Preload("Versions.Targets").Find(&mod, "id = ?", modID)

	if mod.ID == "" {
		return nil
	}

	dbCache.Set("GetModById_"+modID, &mod, cache.DefaultExpiration)

	return &mod
}

func GetModsByID(ctx context.Context, modIds []string) []Mod {
	cacheKey := "GetModsById_" + strings.Join(modIds, ":")
	if mods, ok := dbCache.Get(cacheKey); ok {
		return mods.([]Mod)
	}

	var mods []Mod
	DBCtx(ctx).Preload("Tags").Find(&mods, "id in (?)", modIds)

	if len(modIds) != len(mods) {
		return nil
	}

	dbCache.Set(cacheKey, mods, cache.DefaultExpiration)

	return mods
}

func GetModCount(ctx context.Context, search string, unapproved bool) int64 {
	cacheLey := "GetModCount_" + search + "_" + fmt.Sprint(unapproved)
	if count, ok := dbCache.Get(cacheLey); ok {
		return count.(int64)
	}

	var modCount int64
	query := DBCtx(ctx).Model(Mod{}).Where("approved = ? AND denied = ?", !unapproved, false)

	if search != "" {
		query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.ReplaceAll(search, " ", " & "))
	}

	query.Count(&modCount)

	dbCache.Set(cacheLey, modCount, cache.DefaultExpiration)

	return modCount
}

func IncrementModViews(ctx context.Context, mod *Mod) {
	// TODO unignore
	// DBCtx(ctx).Model(mod).Update("views", mod.Views+1)
}

func GetMods(ctx context.Context, limit int, offset int, orderBy string, order string, search string, unapproved bool) []Mod {
	cacheKey := "GetMods_" + fmt.Sprint(limit) + "_" + fmt.Sprint(offset) + "_" + orderBy + "_" + order + "_" + search + "_" + fmt.Sprint(unapproved)
	if mods, ok := dbCache.Get(cacheKey); ok {
		return mods.([]Mod)
	}

	var mods []Mod
	query := DBCtx(ctx).Limit(limit).Offset(offset)

	if orderBy == "last_version_date" {
		query = query.Order("case when last_version_date is null then 1 else 0 end, last_version_date")
	} else {
		query = query.Order(orderBy + " " + order)
	}

	query = query.Where("approved = ? AND denied = ?", !unapproved, false)

	if search != "" {
		query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.ReplaceAll(search, " ", " & "))
	}

	query.Find(&mods)

	dbCache.Set(cacheKey, mods, cache.DefaultExpiration)

	return mods
}

func GetModByIDOrReference(ctx context.Context, modIDOrReference string) *Mod {
	cacheKey := "GetModByIDOrReference_" + modIDOrReference
	if mod, ok := dbCache.Get(cacheKey); ok {
		return mod.(*Mod)
	}

	var mod Mod
	DBCtx(ctx).Preload("Tags").Preload("Versions.Targets").Find(&mod, "mod_reference = ? OR id = ?", modIDOrReference, modIDOrReference)

	if mod.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &mod, cache.DefaultExpiration)

	return &mod
}
