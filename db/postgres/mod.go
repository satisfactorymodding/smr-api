package postgres

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

func GetModByID(ctx context.Context, modID string) *Mod {
	cacheKey := "GetModById_" + modID
	if mod, ok := dbCache.Get(cacheKey); ok {
		return mod.(*Mod)
	}

	return GetModByIDNoCache(modID, ctx)
}

func GetModByIDNoCache(modID string, ctx *context.Context) *Mod {
	var mod Mod
	DBCtx(ctx).Preload("Tags").Find(&mod, "id = ?", modID)

	if mod.ID == "" {
		return nil
	}

	dbCache.Set("GetModById_"+modID, &mod, cache.DefaultExpiration)

	return &mod
}

func GetModByReference(ctx context.Context, modReference string) *Mod {
	cacheKey := "GetModByReference_" + modReference
	if mod, ok := dbCache.Get(cacheKey); ok {
		return mod.(*Mod)
	}

	var mod Mod
	DBCtx(ctx).Preload("Tags").Find(&mod, "mod_reference = ?", modReference)

	if mod.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &mod, cache.DefaultExpiration)

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

func DeleteMod(ctx context.Context, modID string) {
	DBCtx(ctx).Delete(Mod{}, "id = ?", modID)
	DBCtx(ctx).Delete(Version{}, "mod_id = ?", modID)
	DBCtx(ctx).Delete(UserMod{}, "mod_id = ?", modID)
}

func GetModCount(ctx context.Context, search string, unapproved bool) int64 {
	cacheLey := "GetModCount_" + search + "_" + fmt.Sprint(unapproved)
	if count, ok := dbCache.Get(cacheLey); ok {
		return count.(int64)
	}

	var modCount int64
	query := DBCtx(ctx).Model(Mod{}).Where("approved = ? AND denied = ?", !unapproved, false)

	if search != "" {
		query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.Replace(search, " ", " & ", -1))
	}

	query.Count(&modCount)

	dbCache.Set(cacheLey, modCount, cache.DefaultExpiration)

	return modCount
}

func GetModCountNew(ctx context.Context, filter *models.ModFilter, unapproved bool) int64 {
	hash, err := filter.Hash()
	cacheKey := ""
	if err == nil {
		cacheKey = "GetModCountNew_" + hash + "_" + fmt.Sprint(unapproved)
		if count, ok := dbCache.Get(cacheKey); ok {
			return count.(int64)
		}
	}

	var modCount int64
	NewModQuery(ctx, filter, unapproved, true).Count(&modCount)

	if cacheKey != "" {
		dbCache.Set(cacheKey, modCount, cache.DefaultExpiration)
	}

	return modCount
}

func IncrementModViews(ctx context.Context, mod *Mod) {
	DBCtx(ctx).Model(mod).Update("views", mod.Views+1)
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
		query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.Replace(search, " ", " & ", -1))
	}

	query.Find(&mods)

	dbCache.Set(cacheKey, mods, cache.DefaultExpiration)

	return mods
}

func GetModsNew(ctx context.Context, filter *models.ModFilter, unapproved bool) []Mod {
	hash, err := filter.Hash()
	cacheKey := ""
	if err == nil {
		cacheKey = "GetModsNew_" + hash + "_" + fmt.Sprint(unapproved)
		if mods, ok := dbCache.Get(cacheKey); ok {
			return mods.([]Mod)
		}
	}

	var mods []Mod
	NewModQuery(ctx, filter, unapproved, false).Find(&mods)

	if cacheKey != "" {
		dbCache.Set(cacheKey, mods, cache.DefaultExpiration)
	}

	return mods
}

func CreateMod(ctx context.Context, mod *Mod) (*Mod, error) {
	// Allow only new 4 mods per 24h

	mod.ID = util.GenerateUniqueID()

	var mods []Mod
	DBCtx(ctx).Order("created_at asc").Find(&mods, "creator_id = ? AND created_at > ?", mod.CreatorID, time.Now().Add(time.Hour*24*-1))

	currentAvailable := float64(4)
	lastModTime := time.Now()
	for _, mod := range mods {
		currentAvailable--
		if mod.CreatedAt.After(lastModTime) {
			diff := mod.CreatedAt.Sub(lastModTime)
			currentAvailable = math.Min(4, currentAvailable+diff.Hours()/6)
		}
		lastModTime = mod.CreatedAt
	}

	if currentAvailable < 1 {
		timeToWait := time.Until(lastModTime.Add(time.Hour * 6)).Minutes()
		return nil, fmt.Errorf("please wait %.0f minutes to post another mod", timeToWait)
	}

	DBCtx(ctx).Create(&mod)
	DBCtx(ctx).Create(&UserMod{
		Role:   "creator",
		ModID:  mod.ID,
		UserID: mod.CreatorID,
	})

	return mod, nil
}

func NewModQuery(ctx context.Context, filter *models.ModFilter, unapproved bool, count bool) *gorm.DB {
	query := DBCtx(ctx)

	if count {
		query = query.Model(Mod{})
	}

	query = query.Where("approved = ? AND denied = ?", !unapproved, false)
	query.Preload("Tags")
	if filter != nil {
		if filter.Search != nil && *filter.Search != "" {
			cleanSearch := strings.Replace(strings.TrimSpace(*filter.Search), " ", " & ", -1)
			sub := DBCtx(ctx).Table("mods")
			sub = sub.Select("id, (similarity(name, ?) * 2 + similarity(short_description, ?) + similarity(full_description, ?) * 0.5) as s", cleanSearch, cleanSearch, cleanSearch)

			query = query.Joins("INNER JOIN (?) AS t1 on t1.id = mods.id", sub)
			query = query.Where("t1.s > 0.2")

			if !count && *filter.OrderBy == generated.ModFieldsSearch {
				query = query.Order("t1.s DESC")
			}
		}

		if !count {
			query = query.Limit(*filter.Limit).
				Offset(*filter.Offset)

			if *filter.OrderBy != generated.ModFieldsSearch {
				if string(*filter.OrderBy) == "last_version_date" {
					query = query.Order("case when last_version_date is null then 1 else 0 end, last_version_date " + string(*filter.Order))
				} else {
					query = query.Order("mods." + string(*filter.OrderBy) + " " + string(*filter.Order))
				}
			}
		}

		if filter.Hidden == nil || !(*filter.Hidden) {
			query = query.Where("hidden = false")
		}

		if filter.Ids != nil && len(filter.Ids) > 0 {
			query = query.Where("mods.id in (?)", filter.Ids)
		} else if filter.References != nil && len(filter.References) > 0 {
			query = query.Where("mod_reference in (?)", filter.References)
		}

		if filter.Fields != nil && len(filter.Fields) > 0 {
			query = query.Select(filter.Fields)
		}

		if filter.TagIDs != nil && len(filter.TagIDs) > 0 {
			query.Joins("INNER JOIN mod_tags on mod_tags.tag_id in ? AND mod_tags.mod_id = mods.id", filter.TagIDs)
		}
	}

	return query.Debug()
}

func ClearModTags(modID string, ctx *context.Context) error {
	r := DBCtx(ctx).Where("mod_id = ?", modID).Delete(&ModTag{})
	return r.Error
}

func SetModTags(modID string, tagIDs []string, ctx *context.Context) error {
	for _, tag := range tagIDs {
		err := AddModTag(modID, tag, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func ResetModTags(modID string, tagIDs []string, ctx *context.Context) error {
	err := ClearModTags(modID, ctx)
	if err != nil {
		return err
	}
	err = SetModTags(modID, tagIDs, ctx)
	if err != nil {
		return err
	}
	return nil
}

func AddModTag(modID string, tagID string, ctx *context.Context) error {
	r := DBCtx(ctx).Create(&ModTag{ModID: modID, TagID: tagID})
	return r.Error
}

func RemoveModTag(modID string, tagID string, ctx *context.Context) error {
	r := DBCtx(ctx).Delete(&ModTag{ModID: modID, TagID: tagID})
	return r.Error
}
