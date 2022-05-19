package models

import (
	"reflect"
	"strconv"

	"github.com/satisfactorymodding/smr-api/generated"

	"github.com/pkg/errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/go-playground/validator.v9"
)

var dataValidator = validator.New()

type VersionFilter struct {
	Limit   *int                     `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset  *int                     `json:"offset" validate:"omitempty,min=0"`
	OrderBy *generated.VersionFields `json:"order_by"`
	Order   *generated.Order         `json:"order"`
	Search  *string                  `json:"search" validate:"omitempty,min=3"`
	Ids     []string                 `json:"ids" validate:"omitempty,max=100"`
	Fields  []string                 `json:"-"`
}

func (f *VersionFilter) IsDefault(ignoreLimits bool) bool {
	return ((f.Limit != nil && *f.Limit == 10) || ignoreLimits) &&
		f.Offset != nil && *f.Offset == 0 &&
		f.Ids == nil &&
		f.Order != nil && *f.Order == generated.OrderDesc &&
		f.OrderBy != nil && *f.OrderBy == generated.VersionFieldsCreatedAt
}

func DefaultVersionFilter() *VersionFilter {
	limit := 10
	offset := 0
	order := generated.OrderDesc
	orderBy := generated.VersionFieldsCreatedAt
	return &VersionFilter{
		Limit:   &limit,
		Offset:  &offset,
		Ids:     nil,
		Order:   &order,
		OrderBy: &orderBy,
		Fields:  nil,
	}
}

func (f *VersionFilter) AddField(name string) {
	switch name {
	case "id",
		"mod_id",
		"version",
		"sml_version",
		"changelog",
		"downloads",
		"stability",
		"approved",
		"updated_at",
		"created_at",
		"metadata",
		"size",
		"hash":
		f.Fields = append(f.Fields, name)
	case "link":
		f.Fields = append(f.Fields, "key")
	}
}

func (f *VersionFilter) Hash() (string, error) {
	hash, err := hashstructure.Hash(f, hashstructure.FormatV2, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash VersionFilter")
	}
	return strconv.FormatUint(hash, 10), nil
}

func ProcessVersionFilter(filter map[string]interface{}) (*VersionFilter, error) {
	base := DefaultVersionFilter()

	if filter == nil {
		return base, nil
	}

	if err := ApplyChanges(filter, base); err != nil {
		return nil, err
	}

	if err := dataValidator.Struct(base); err != nil {
		return nil, errors.Wrap(err, "failed to validate VersionFilter")
	}

	return base, nil
}

type ModFilter struct {
	Limit      *int                 `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset     *int                 `json:"offset" validate:"omitempty,min=0"`
	OrderBy    *generated.ModFields `json:"order_by"`
	Order      *generated.Order     `json:"order"`
	Search     *string              `json:"search" validate:"omitempty,min=3"`
	Ids        []string             `json:"ids" validate:"omitempty,max=100"`
	References []string             `json:"references" validate:"omitempty,max=100"`
	Hidden     *bool                `json:"hidden"`
	Fields     []string             `json:"-"`
	TagIDs     []string             `json:"tagIDs" validate:"omitempty,max=100"`
}

func DefaultModFilter() *ModFilter {
	limit := 10
	offset := 0
	order := generated.OrderDesc
	orderBy := generated.ModFieldsCreatedAt
	return &ModFilter{
		Limit:   &limit,
		Offset:  &offset,
		Ids:     nil,
		Order:   &order,
		OrderBy: &orderBy,
		Fields:  nil,
	}
}

func (f *ModFilter) AddField(name string) {
	switch name {
	case "id",
		"name",
		"short_description",
		"full_description",
		"logo",
		"source_url",
		"creator_id",
		"approved",
		"views",
		"downloads",
		"hotness",
		"popularity",
		"updated_at",
		"created_at",
		"last_version_date",
		"mod_reference",
		"hidden":
		f.Fields = append(f.Fields, "mods."+name)
	}
}

func (f *ModFilter) Hash() (string, error) {
	hash, err := hashstructure.Hash(f, hashstructure.FormatV2, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash ModFilter")
	}
	return strconv.FormatUint(hash, 10), nil
}

func ProcessModFilter(filter map[string]interface{}) (*ModFilter, error) {
	base := DefaultModFilter()

	if filter == nil {
		return base, nil
	}

	if err := ApplyChanges(filter, base); err != nil {
		return nil, err
	}

	if err := dataValidator.Struct(base); err != nil {
		return nil, errors.Wrap(err, "failed to validate ModFilter")
	}

	return base, nil
}

type GuideFilter struct {
	Limit   *int                   `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset  *int                   `json:"offset" validate:"omitempty,min=0"`
	OrderBy *generated.GuideFields `json:"order_by"`
	Order   *generated.Order       `json:"order"`
	Search  *string                `json:"search" validate:"omitempty,min=3"`
	Ids     []string               `json:"ids" validate:"omitempty,max=100"`
	TagIDs  []string               `json:"tagIDs" validate:"omitempty,max=100"`
}

func (f GuideFilter) Hash() (string, error) {
	hash, err := hashstructure.Hash(f, hashstructure.FormatV2, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash GuideFilter")
	}
	return strconv.FormatUint(hash, 10), nil
}

func DefaultGuideFilter() *GuideFilter {
	limit := 10
	offset := 0
	order := generated.OrderDesc
	orderBy := generated.GuideFieldsCreatedAt
	return &GuideFilter{
		Limit:   &limit,
		Offset:  &offset,
		Ids:     nil,
		Order:   &order,
		OrderBy: &orderBy,
	}
}

func ProcessGuideFilter(filter map[string]interface{}) (*GuideFilter, error) {
	base := DefaultGuideFilter()

	if filter == nil {
		return base, nil
	}

	if err := ApplyChanges(filter, base); err != nil {
		return nil, err
	}

	if err := dataValidator.Struct(base); err != nil {
		return nil, errors.Wrap(err, "failed to validate GuideFilter")
	}

	return base, nil
}

type SMLVersionFilter struct {
	Limit   *int                        `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset  *int                        `json:"offset" validate:"omitempty,min=0"`
	OrderBy *generated.SMLVersionFields `json:"order_by"`
	Order   *generated.Order            `json:"order"`
	Search  *string                     `json:"search" validate:"omitempty,min=3"`
	Ids     []string                    `json:"ids" validate:"omitempty,max=100"`
}

func DefaultSMLVersionFilter() *SMLVersionFilter {
	limit := 10
	offset := 0
	order := generated.OrderDesc
	orderBy := generated.SMLVersionFieldsCreatedAt
	return &SMLVersionFilter{
		Limit:   &limit,
		Offset:  &offset,
		Ids:     nil,
		Order:   &order,
		OrderBy: &orderBy,
	}
}

func ProcessSMLVersionFilter(filter map[string]interface{}) (*SMLVersionFilter, error) {
	base := DefaultSMLVersionFilter()

	if filter == nil {
		return base, nil
	}

	if err := ApplyChanges(filter, base); err != nil {
		return nil, err
	}

	if err := dataValidator.Struct(base); err != nil {
		return nil, errors.Wrap(err, "failed to validate SMLVersionFilter")
	}

	return base, nil
}

func ApplyChanges(changes interface{}, to interface{}) error {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		TagName:     "json",
		Result:      to,
		ZeroFields:  true,
		DecodeHook: func(a reflect.Type, b reflect.Type, v interface{}) (interface{}, error) {
			if reflect.PtrTo(b).Implements(reflect.TypeOf((*graphql.Unmarshaler)(nil)).Elem()) {
				resultType := reflect.New(b)
				result := resultType.MethodByName("UnmarshalGQL").Call([]reflect.Value{reflect.ValueOf(v)})
				err, _ := result[0].Interface().(error)
				return resultType.Elem().Interface(), err
			}

			return v, nil
		},
	})

	if err != nil {
		return errors.Wrap(err, "failed to create decoder")
	}

	return errors.Wrap(dec.Decode(changes), "failed to decode changes")
}

type BootstrapVersionFilter struct {
	Limit   *int                              `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset  *int                              `json:"offset" validate:"omitempty,min=0"`
	OrderBy *generated.BootstrapVersionFields `json:"order_by"`
	Order   *generated.Order                  `json:"order"`
	Search  *string                           `json:"search" validate:"omitempty,min=3"`
	Ids     []string                          `json:"ids" validate:"omitempty,max=100"`
}

func DefaultBootstrapVersionFilter() *BootstrapVersionFilter {
	limit := 10
	offset := 0
	order := generated.OrderDesc
	orderBy := generated.BootstrapVersionFieldsCreatedAt
	return &BootstrapVersionFilter{
		Limit:   &limit,
		Offset:  &offset,
		Ids:     nil,
		Order:   &order,
		OrderBy: &orderBy,
	}
}

func ProcessBootstrapVersionFilter(filter map[string]interface{}) (*BootstrapVersionFilter, error) {
	base := DefaultBootstrapVersionFilter()

	if filter == nil {
		return base, nil
	}

	if err := ApplyChanges(filter, base); err != nil {
		return nil, err
	}

	if err := dataValidator.Struct(base); err != nil {
		return nil, errors.Wrap(err, "failed to validate BootstrapVersionFilter")
	}

	return base, nil
}

type SMLLinkFilter struct {
	Limit   *int                     `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset  *int                     `json:"offset" validate:"omitempty,min=0"`
	OrderBy *generated.SMLLinkFields `json:"order_by"`
	Order   *generated.Order         `json:"order"`
	Search  *string                  `json:"search" validate:"omitempty,min=3"`
	Ids     []string                 `json:"ids" validate:"omitempty,max=100"`
}

func DefaultSMLLinkFilter() *SMLLinkFilter {
	limit := 10
	offset := 0
	order := generated.OrderDesc
	orderBy := generated.SMLLinkFieldsPlatform
	return &SMLLinkFilter{
		Limit:   &limit,
		Offset:  &offset,
		Ids:     nil,
		Order:   &order,
		OrderBy: &orderBy,
	}
}

func ProcessSMLLinkFilter(filter map[string]interface{}) (*SMLLinkFilter, error) {
	base := DefaultSMLLinkFilter()

	if filter == nil {
		return base, nil
	}

	if err := ApplyChanges(filter, base); err != nil {
		return nil, err
	}

	if err := dataValidator.Struct(base); err != nil {
		return nil, errors.Wrap(err, "failed to validate SMLLinkFilter")
	}

	return base, nil
}
