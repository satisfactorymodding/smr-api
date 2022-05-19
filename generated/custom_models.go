package generated

import "github.com/99designs/gqlgen/graphql"

type NewMod struct {
	Name             string          `json:"name" validate:"required,min=3,max=32"`
	ShortDescription string          `json:"short_description" validate:"required,min=16,max=128"`
	FullDescription  *string         `json:"full_description"`
	Logo             *graphql.Upload `json:"logo"`
	SourceURL        *string         `json:"source_url"`
	ModReference     string          `json:"mod_reference"`
	Hidden           *bool           `json:"hidden"`
	TagIDs           []string        `json:"tagIDs" validate:"dive,min=3,max=24"`
	Links            []string        `json:"links"`
}

type UpdateMod struct {
	Name             *string         `json:"name" validate:"omitempty,min=3,max=32"`
	ShortDescription *string         `json:"short_description" validate:"omitempty,min=16,max=128"`
	FullDescription  *string         `json:"full_description"`
	Logo             *graphql.Upload `json:"logo"`
	SourceURL        *string         `json:"source_url"`
	ModReference     *string         `json:"mod_reference"`
	Hidden           *bool           `json:"hidden"`
	Authors          []UpdateUserMod `json:"authors"`
	TagIDs           []string        `json:"tagIDs" validate:"dive,min=3,max=24"`
	Links            []string        `json:"links"`
}
