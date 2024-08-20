// Code generated by ent, DO NOT EDIT.

package runtime

import (
	"time"

	"github.com/satisfactorymodding/smr-api/db/schema"
	"github.com/satisfactorymodding/smr-api/generated/ent/announcement"
	"github.com/satisfactorymodding/smr-api/generated/ent/guide"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/satisfactoryversion"
	"github.com/satisfactorymodding/smr-api/generated/ent/tag"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
	"github.com/satisfactorymodding/smr-api/generated/ent/usergroup"
	"github.com/satisfactorymodding/smr-api/generated/ent/usersession"
	"github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/generated/ent/versiondependency"
	"github.com/satisfactorymodding/smr-api/generated/ent/versiontarget"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	announcementMixin := schema.Announcement{}.Mixin()
	announcementMixinHooks2 := announcementMixin[2].Hooks()
	announcement.Hooks[0] = announcementMixinHooks2[0]
	announcementMixinInters2 := announcementMixin[2].Interceptors()
	announcement.Interceptors[0] = announcementMixinInters2[0]
	announcementMixinFields0 := announcementMixin[0].Fields()
	_ = announcementMixinFields0
	announcementMixinFields1 := announcementMixin[1].Fields()
	_ = announcementMixinFields1
	announcementFields := schema.Announcement{}.Fields()
	_ = announcementFields
	// announcementDescCreatedAt is the schema descriptor for created_at field.
	announcementDescCreatedAt := announcementMixinFields1[0].Descriptor()
	// announcement.DefaultCreatedAt holds the default value on creation for the created_at field.
	announcement.DefaultCreatedAt = announcementDescCreatedAt.Default.(func() time.Time)
	// announcementDescUpdatedAt is the schema descriptor for updated_at field.
	announcementDescUpdatedAt := announcementMixinFields1[1].Descriptor()
	// announcement.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	announcement.DefaultUpdatedAt = announcementDescUpdatedAt.Default.(func() time.Time)
	// announcement.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	announcement.UpdateDefaultUpdatedAt = announcementDescUpdatedAt.UpdateDefault.(func() time.Time)
	// announcementDescID is the schema descriptor for id field.
	announcementDescID := announcementMixinFields0[0].Descriptor()
	// announcement.DefaultID holds the default value on creation for the id field.
	announcement.DefaultID = announcementDescID.Default.(func() string)
	guideMixin := schema.Guide{}.Mixin()
	guideMixinHooks2 := guideMixin[2].Hooks()
	guide.Hooks[0] = guideMixinHooks2[0]
	guideMixinInters2 := guideMixin[2].Interceptors()
	guide.Interceptors[0] = guideMixinInters2[0]
	guideMixinFields0 := guideMixin[0].Fields()
	_ = guideMixinFields0
	guideMixinFields1 := guideMixin[1].Fields()
	_ = guideMixinFields1
	guideFields := schema.Guide{}.Fields()
	_ = guideFields
	// guideDescCreatedAt is the schema descriptor for created_at field.
	guideDescCreatedAt := guideMixinFields1[0].Descriptor()
	// guide.DefaultCreatedAt holds the default value on creation for the created_at field.
	guide.DefaultCreatedAt = guideDescCreatedAt.Default.(func() time.Time)
	// guideDescUpdatedAt is the schema descriptor for updated_at field.
	guideDescUpdatedAt := guideMixinFields1[1].Descriptor()
	// guide.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	guide.DefaultUpdatedAt = guideDescUpdatedAt.Default.(func() time.Time)
	// guide.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	guide.UpdateDefaultUpdatedAt = guideDescUpdatedAt.UpdateDefault.(func() time.Time)
	// guideDescName is the schema descriptor for name field.
	guideDescName := guideFields[1].Descriptor()
	// guide.NameValidator is a validator for the "name" field. It is called by the builders before save.
	guide.NameValidator = guideDescName.Validators[0].(func(string) error)
	// guideDescShortDescription is the schema descriptor for short_description field.
	guideDescShortDescription := guideFields[2].Descriptor()
	// guide.ShortDescriptionValidator is a validator for the "short_description" field. It is called by the builders before save.
	guide.ShortDescriptionValidator = guideDescShortDescription.Validators[0].(func(string) error)
	// guideDescViews is the schema descriptor for views field.
	guideDescViews := guideFields[4].Descriptor()
	// guide.DefaultViews holds the default value on creation for the views field.
	guide.DefaultViews = guideDescViews.Default.(int)
	// guideDescID is the schema descriptor for id field.
	guideDescID := guideMixinFields0[0].Descriptor()
	// guide.DefaultID holds the default value on creation for the id field.
	guide.DefaultID = guideDescID.Default.(func() string)
	modMixin := schema.Mod{}.Mixin()
	modMixinHooks2 := modMixin[2].Hooks()
	mod.Hooks[0] = modMixinHooks2[0]
	modMixinInters2 := modMixin[2].Interceptors()
	mod.Interceptors[0] = modMixinInters2[0]
	modMixinFields0 := modMixin[0].Fields()
	_ = modMixinFields0
	modMixinFields1 := modMixin[1].Fields()
	_ = modMixinFields1
	modFields := schema.Mod{}.Fields()
	_ = modFields
	// modDescCreatedAt is the schema descriptor for created_at field.
	modDescCreatedAt := modMixinFields1[0].Descriptor()
	// mod.DefaultCreatedAt holds the default value on creation for the created_at field.
	mod.DefaultCreatedAt = modDescCreatedAt.Default.(func() time.Time)
	// modDescUpdatedAt is the schema descriptor for updated_at field.
	modDescUpdatedAt := modMixinFields1[1].Descriptor()
	// mod.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	mod.DefaultUpdatedAt = modDescUpdatedAt.Default.(func() time.Time)
	// mod.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	mod.UpdateDefaultUpdatedAt = modDescUpdatedAt.UpdateDefault.(func() time.Time)
	// modDescName is the schema descriptor for name field.
	modDescName := modFields[0].Descriptor()
	// mod.NameValidator is a validator for the "name" field. It is called by the builders before save.
	mod.NameValidator = modDescName.Validators[0].(func(string) error)
	// modDescShortDescription is the schema descriptor for short_description field.
	modDescShortDescription := modFields[1].Descriptor()
	// mod.ShortDescriptionValidator is a validator for the "short_description" field. It is called by the builders before save.
	mod.ShortDescriptionValidator = modDescShortDescription.Validators[0].(func(string) error)
	// modDescApproved is the schema descriptor for approved field.
	modDescApproved := modFields[6].Descriptor()
	// mod.DefaultApproved holds the default value on creation for the approved field.
	mod.DefaultApproved = modDescApproved.Default.(bool)
	// modDescViews is the schema descriptor for views field.
	modDescViews := modFields[7].Descriptor()
	// mod.DefaultViews holds the default value on creation for the views field.
	mod.DefaultViews = modDescViews.Default.(uint)
	// modDescHotness is the schema descriptor for hotness field.
	modDescHotness := modFields[8].Descriptor()
	// mod.DefaultHotness holds the default value on creation for the hotness field.
	mod.DefaultHotness = modDescHotness.Default.(uint)
	// modDescPopularity is the schema descriptor for popularity field.
	modDescPopularity := modFields[9].Descriptor()
	// mod.DefaultPopularity holds the default value on creation for the popularity field.
	mod.DefaultPopularity = modDescPopularity.Default.(uint)
	// modDescDownloads is the schema descriptor for downloads field.
	modDescDownloads := modFields[10].Descriptor()
	// mod.DefaultDownloads holds the default value on creation for the downloads field.
	mod.DefaultDownloads = modDescDownloads.Default.(uint)
	// modDescDenied is the schema descriptor for denied field.
	modDescDenied := modFields[11].Descriptor()
	// mod.DefaultDenied holds the default value on creation for the denied field.
	mod.DefaultDenied = modDescDenied.Default.(bool)
	// modDescModReference is the schema descriptor for mod_reference field.
	modDescModReference := modFields[13].Descriptor()
	// mod.ModReferenceValidator is a validator for the "mod_reference" field. It is called by the builders before save.
	mod.ModReferenceValidator = modDescModReference.Validators[0].(func(string) error)
	// modDescHidden is the schema descriptor for hidden field.
	modDescHidden := modFields[14].Descriptor()
	// mod.DefaultHidden holds the default value on creation for the hidden field.
	mod.DefaultHidden = modDescHidden.Default.(bool)
	// modDescID is the schema descriptor for id field.
	modDescID := modMixinFields0[0].Descriptor()
	// mod.DefaultID holds the default value on creation for the id field.
	mod.DefaultID = modDescID.Default.(func() string)
	satisfactoryversionMixin := schema.SatisfactoryVersion{}.Mixin()
	satisfactoryversionMixinFields0 := satisfactoryversionMixin[0].Fields()
	_ = satisfactoryversionMixinFields0
	satisfactoryversionFields := schema.SatisfactoryVersion{}.Fields()
	_ = satisfactoryversionFields
	// satisfactoryversionDescEngineVersion is the schema descriptor for engine_version field.
	satisfactoryversionDescEngineVersion := satisfactoryversionFields[1].Descriptor()
	// satisfactoryversion.DefaultEngineVersion holds the default value on creation for the engine_version field.
	satisfactoryversion.DefaultEngineVersion = satisfactoryversionDescEngineVersion.Default.(string)
	// satisfactoryversion.EngineVersionValidator is a validator for the "engine_version" field. It is called by the builders before save.
	satisfactoryversion.EngineVersionValidator = satisfactoryversionDescEngineVersion.Validators[0].(func(string) error)
	// satisfactoryversionDescID is the schema descriptor for id field.
	satisfactoryversionDescID := satisfactoryversionMixinFields0[0].Descriptor()
	// satisfactoryversion.DefaultID holds the default value on creation for the id field.
	satisfactoryversion.DefaultID = satisfactoryversionDescID.Default.(func() string)
	tagMixin := schema.Tag{}.Mixin()
	tagMixinHooks2 := tagMixin[2].Hooks()
	tag.Hooks[0] = tagMixinHooks2[0]
	tagMixinInters2 := tagMixin[2].Interceptors()
	tag.Interceptors[0] = tagMixinInters2[0]
	tagMixinFields0 := tagMixin[0].Fields()
	_ = tagMixinFields0
	tagMixinFields1 := tagMixin[1].Fields()
	_ = tagMixinFields1
	tagFields := schema.Tag{}.Fields()
	_ = tagFields
	// tagDescCreatedAt is the schema descriptor for created_at field.
	tagDescCreatedAt := tagMixinFields1[0].Descriptor()
	// tag.DefaultCreatedAt holds the default value on creation for the created_at field.
	tag.DefaultCreatedAt = tagDescCreatedAt.Default.(func() time.Time)
	// tagDescUpdatedAt is the schema descriptor for updated_at field.
	tagDescUpdatedAt := tagMixinFields1[1].Descriptor()
	// tag.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	tag.DefaultUpdatedAt = tagDescUpdatedAt.Default.(func() time.Time)
	// tag.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	tag.UpdateDefaultUpdatedAt = tagDescUpdatedAt.UpdateDefault.(func() time.Time)
	// tagDescName is the schema descriptor for name field.
	tagDescName := tagFields[0].Descriptor()
	// tag.NameValidator is a validator for the "name" field. It is called by the builders before save.
	tag.NameValidator = tagDescName.Validators[0].(func(string) error)
	// tagDescDescription is the schema descriptor for description field.
	tagDescDescription := tagFields[1].Descriptor()
	// tag.DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	tag.DescriptionValidator = tagDescDescription.Validators[0].(func(string) error)
	// tagDescID is the schema descriptor for id field.
	tagDescID := tagMixinFields0[0].Descriptor()
	// tag.DefaultID holds the default value on creation for the id field.
	tag.DefaultID = tagDescID.Default.(func() string)
	userMixin := schema.User{}.Mixin()
	userMixinHooks2 := userMixin[2].Hooks()
	user.Hooks[0] = userMixinHooks2[0]
	userMixinInters2 := userMixin[2].Interceptors()
	user.Interceptors[0] = userMixinInters2[0]
	userMixinFields0 := userMixin[0].Fields()
	_ = userMixinFields0
	userMixinFields1 := userMixin[1].Fields()
	_ = userMixinFields1
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userMixinFields1[0].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userMixinFields1[1].Descriptor()
	// user.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	user.DefaultUpdatedAt = userDescUpdatedAt.Default.(func() time.Time)
	// user.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	user.UpdateDefaultUpdatedAt = userDescUpdatedAt.UpdateDefault.(func() time.Time)
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[0].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = userDescEmail.Validators[0].(func(string) error)
	// userDescUsername is the schema descriptor for username field.
	userDescUsername := userFields[1].Descriptor()
	// user.UsernameValidator is a validator for the "username" field. It is called by the builders before save.
	user.UsernameValidator = userDescUsername.Validators[0].(func(string) error)
	// userDescBanned is the schema descriptor for banned field.
	userDescBanned := userFields[4].Descriptor()
	// user.DefaultBanned holds the default value on creation for the banned field.
	user.DefaultBanned = userDescBanned.Default.(bool)
	// userDescRank is the schema descriptor for rank field.
	userDescRank := userFields[5].Descriptor()
	// user.DefaultRank holds the default value on creation for the rank field.
	user.DefaultRank = userDescRank.Default.(int)
	// userDescGithubID is the schema descriptor for github_id field.
	userDescGithubID := userFields[6].Descriptor()
	// user.GithubIDValidator is a validator for the "github_id" field. It is called by the builders before save.
	user.GithubIDValidator = userDescGithubID.Validators[0].(func(string) error)
	// userDescGoogleID is the schema descriptor for google_id field.
	userDescGoogleID := userFields[7].Descriptor()
	// user.GoogleIDValidator is a validator for the "google_id" field. It is called by the builders before save.
	user.GoogleIDValidator = userDescGoogleID.Validators[0].(func(string) error)
	// userDescFacebookID is the schema descriptor for facebook_id field.
	userDescFacebookID := userFields[8].Descriptor()
	// user.FacebookIDValidator is a validator for the "facebook_id" field. It is called by the builders before save.
	user.FacebookIDValidator = userDescFacebookID.Validators[0].(func(string) error)
	// userDescID is the schema descriptor for id field.
	userDescID := userMixinFields0[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() string)
	usergroupMixin := schema.UserGroup{}.Mixin()
	usergroupMixinHooks2 := usergroupMixin[2].Hooks()
	usergroup.Hooks[0] = usergroupMixinHooks2[0]
	usergroupMixinInters2 := usergroupMixin[2].Interceptors()
	usergroup.Interceptors[0] = usergroupMixinInters2[0]
	usergroupMixinFields0 := usergroupMixin[0].Fields()
	_ = usergroupMixinFields0
	usergroupMixinFields1 := usergroupMixin[1].Fields()
	_ = usergroupMixinFields1
	usergroupFields := schema.UserGroup{}.Fields()
	_ = usergroupFields
	// usergroupDescCreatedAt is the schema descriptor for created_at field.
	usergroupDescCreatedAt := usergroupMixinFields1[0].Descriptor()
	// usergroup.DefaultCreatedAt holds the default value on creation for the created_at field.
	usergroup.DefaultCreatedAt = usergroupDescCreatedAt.Default.(func() time.Time)
	// usergroupDescUpdatedAt is the schema descriptor for updated_at field.
	usergroupDescUpdatedAt := usergroupMixinFields1[1].Descriptor()
	// usergroup.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	usergroup.DefaultUpdatedAt = usergroupDescUpdatedAt.Default.(func() time.Time)
	// usergroup.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	usergroup.UpdateDefaultUpdatedAt = usergroupDescUpdatedAt.UpdateDefault.(func() time.Time)
	// usergroupDescUserID is the schema descriptor for user_id field.
	usergroupDescUserID := usergroupFields[0].Descriptor()
	// usergroup.UserIDValidator is a validator for the "user_id" field. It is called by the builders before save.
	usergroup.UserIDValidator = usergroupDescUserID.Validators[0].(func(string) error)
	// usergroupDescGroupID is the schema descriptor for group_id field.
	usergroupDescGroupID := usergroupFields[1].Descriptor()
	// usergroup.GroupIDValidator is a validator for the "group_id" field. It is called by the builders before save.
	usergroup.GroupIDValidator = usergroupDescGroupID.Validators[0].(func(string) error)
	// usergroupDescID is the schema descriptor for id field.
	usergroupDescID := usergroupMixinFields0[0].Descriptor()
	// usergroup.DefaultID holds the default value on creation for the id field.
	usergroup.DefaultID = usergroupDescID.Default.(func() string)
	usersessionMixin := schema.UserSession{}.Mixin()
	usersessionMixinHooks2 := usersessionMixin[2].Hooks()
	usersession.Hooks[0] = usersessionMixinHooks2[0]
	usersessionMixinInters2 := usersessionMixin[2].Interceptors()
	usersession.Interceptors[0] = usersessionMixinInters2[0]
	usersessionMixinFields0 := usersessionMixin[0].Fields()
	_ = usersessionMixinFields0
	usersessionMixinFields1 := usersessionMixin[1].Fields()
	_ = usersessionMixinFields1
	usersessionFields := schema.UserSession{}.Fields()
	_ = usersessionFields
	// usersessionDescCreatedAt is the schema descriptor for created_at field.
	usersessionDescCreatedAt := usersessionMixinFields1[0].Descriptor()
	// usersession.DefaultCreatedAt holds the default value on creation for the created_at field.
	usersession.DefaultCreatedAt = usersessionDescCreatedAt.Default.(func() time.Time)
	// usersessionDescUpdatedAt is the schema descriptor for updated_at field.
	usersessionDescUpdatedAt := usersessionMixinFields1[1].Descriptor()
	// usersession.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	usersession.DefaultUpdatedAt = usersessionDescUpdatedAt.Default.(func() time.Time)
	// usersession.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	usersession.UpdateDefaultUpdatedAt = usersessionDescUpdatedAt.UpdateDefault.(func() time.Time)
	// usersessionDescToken is the schema descriptor for token field.
	usersessionDescToken := usersessionFields[0].Descriptor()
	// usersession.TokenValidator is a validator for the "token" field. It is called by the builders before save.
	usersession.TokenValidator = usersessionDescToken.Validators[0].(func(string) error)
	// usersessionDescID is the schema descriptor for id field.
	usersessionDescID := usersessionMixinFields0[0].Descriptor()
	// usersession.DefaultID holds the default value on creation for the id field.
	usersession.DefaultID = usersessionDescID.Default.(func() string)
	versionMixin := schema.Version{}.Mixin()
	versionMixinHooks2 := versionMixin[2].Hooks()
	version.Hooks[0] = versionMixinHooks2[0]
	versionMixinInters2 := versionMixin[2].Interceptors()
	version.Interceptors[0] = versionMixinInters2[0]
	versionMixinFields0 := versionMixin[0].Fields()
	_ = versionMixinFields0
	versionMixinFields1 := versionMixin[1].Fields()
	_ = versionMixinFields1
	versionFields := schema.Version{}.Fields()
	_ = versionFields
	// versionDescCreatedAt is the schema descriptor for created_at field.
	versionDescCreatedAt := versionMixinFields1[0].Descriptor()
	// version.DefaultCreatedAt holds the default value on creation for the created_at field.
	version.DefaultCreatedAt = versionDescCreatedAt.Default.(func() time.Time)
	// versionDescUpdatedAt is the schema descriptor for updated_at field.
	versionDescUpdatedAt := versionMixinFields1[1].Descriptor()
	// version.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	version.DefaultUpdatedAt = versionDescUpdatedAt.Default.(func() time.Time)
	// version.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	version.UpdateDefaultUpdatedAt = versionDescUpdatedAt.UpdateDefault.(func() time.Time)
	// versionDescVersion is the schema descriptor for version field.
	versionDescVersion := versionFields[1].Descriptor()
	// version.VersionValidator is a validator for the "version" field. It is called by the builders before save.
	version.VersionValidator = versionDescVersion.Validators[0].(func(string) error)
	// versionDescDownloads is the schema descriptor for downloads field.
	versionDescDownloads := versionFields[4].Descriptor()
	// version.DefaultDownloads holds the default value on creation for the downloads field.
	version.DefaultDownloads = versionDescDownloads.Default.(uint)
	// versionDescApproved is the schema descriptor for approved field.
	versionDescApproved := versionFields[7].Descriptor()
	// version.DefaultApproved holds the default value on creation for the approved field.
	version.DefaultApproved = versionDescApproved.Default.(bool)
	// versionDescHotness is the schema descriptor for hotness field.
	versionDescHotness := versionFields[8].Descriptor()
	// version.DefaultHotness holds the default value on creation for the hotness field.
	version.DefaultHotness = versionDescHotness.Default.(uint)
	// versionDescDenied is the schema descriptor for denied field.
	versionDescDenied := versionFields[9].Descriptor()
	// version.DefaultDenied holds the default value on creation for the denied field.
	version.DefaultDenied = versionDescDenied.Default.(bool)
	// versionDescModReference is the schema descriptor for mod_reference field.
	versionDescModReference := versionFields[11].Descriptor()
	// version.ModReferenceValidator is a validator for the "mod_reference" field. It is called by the builders before save.
	version.ModReferenceValidator = versionDescModReference.Validators[0].(func(string) error)
	// versionDescHash is the schema descriptor for hash field.
	versionDescHash := versionFields[16].Descriptor()
	// version.HashValidator is a validator for the "hash" field. It is called by the builders before save.
	version.HashValidator = func() func(string) error {
		validators := versionDescHash.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(hash string) error {
			for _, fn := range fns {
				if err := fn(hash); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// versionDescID is the schema descriptor for id field.
	versionDescID := versionMixinFields0[0].Descriptor()
	// version.DefaultID holds the default value on creation for the id field.
	version.DefaultID = versionDescID.Default.(func() string)
	versiondependencyMixin := schema.VersionDependency{}.Mixin()
	versiondependencyMixinHooks1 := versiondependencyMixin[1].Hooks()
	versiondependency.Hooks[0] = versiondependencyMixinHooks1[0]
	versiondependencyMixinInters1 := versiondependencyMixin[1].Interceptors()
	versiondependency.Interceptors[0] = versiondependencyMixinInters1[0]
	versiondependencyMixinFields0 := versiondependencyMixin[0].Fields()
	_ = versiondependencyMixinFields0
	versiondependencyFields := schema.VersionDependency{}.Fields()
	_ = versiondependencyFields
	// versiondependencyDescCreatedAt is the schema descriptor for created_at field.
	versiondependencyDescCreatedAt := versiondependencyMixinFields0[0].Descriptor()
	// versiondependency.DefaultCreatedAt holds the default value on creation for the created_at field.
	versiondependency.DefaultCreatedAt = versiondependencyDescCreatedAt.Default.(func() time.Time)
	// versiondependencyDescUpdatedAt is the schema descriptor for updated_at field.
	versiondependencyDescUpdatedAt := versiondependencyMixinFields0[1].Descriptor()
	// versiondependency.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	versiondependency.DefaultUpdatedAt = versiondependencyDescUpdatedAt.Default.(func() time.Time)
	// versiondependency.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	versiondependency.UpdateDefaultUpdatedAt = versiondependencyDescUpdatedAt.UpdateDefault.(func() time.Time)
	// versiondependencyDescCondition is the schema descriptor for condition field.
	versiondependencyDescCondition := versiondependencyFields[2].Descriptor()
	// versiondependency.ConditionValidator is a validator for the "condition" field. It is called by the builders before save.
	versiondependency.ConditionValidator = versiondependencyDescCondition.Validators[0].(func(string) error)
	versiontargetMixin := schema.VersionTarget{}.Mixin()
	versiontargetMixinFields0 := versiontargetMixin[0].Fields()
	_ = versiontargetMixinFields0
	versiontargetFields := schema.VersionTarget{}.Fields()
	_ = versiontargetFields
	// versiontargetDescID is the schema descriptor for id field.
	versiontargetDescID := versiontargetMixinFields0[0].Descriptor()
	// versiontarget.DefaultID holds the default value on creation for the id field.
	versiontarget.DefaultID = versiontargetDescID.Default.(func() string)
}

const (
	Version = "v0.14.0"                                         // Version of ent codegen.
	Sum     = "h1:EO3Z9aZ5bXJatJeGqu/EVdnNr6K4mRq3rWe5owt0MC4=" // Sum of ent codegen.
)
