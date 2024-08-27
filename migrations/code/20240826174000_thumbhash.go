package code

import (
	"context"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"

	"github.com/Vilsol/slox"
	"github.com/alitto/pond"
	"github.com/galdor/go-thumbhash"
	"github.com/lab259/go-migration"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
	"github.com/satisfactorymodding/smr-api/util/converter"
)

func init() {
	migration.NewCodeMigration(
		func(ctxInt interface{}) error {
			ctx := ctxInt.(context.Context)

			pool := pond.New(32, 0, pond.MinWorkers(32))

			// Calculate for all mods
			mods, err := db.From(ctx).Mod.Query().Select(mod.FieldID, mod.FieldLogo).Where(mod.LogoThumbhashIsNil()).All(ctx)
			if err != nil {
				return err
			}

			for i, m := range mods {
				if i%50 == 0 {
					slox.Info(ctx, "generated thumbhash for n mods", slog.Int("n", i))
				}

				pool.Submit(func() {
					if m.Logo == "" {
						return
					}

					resp, err := http.Get(m.Logo)
					if err != nil {
						slox.Error(ctx, "invalid url", slog.String("mod_id", m.ID), slog.String("logo", m.Logo), slog.Any("err", err))
						return
					}

					defer resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						slox.Error(ctx, "invalid url", slog.String("mod_id", m.ID), slog.String("logo", m.Logo), slog.Any("err", err))
						return
					}

					data, err := io.ReadAll(resp.Body)
					if err != nil {
						slox.Error(ctx, "invalid url", slog.String("mod_id", m.ID), slog.String("logo", m.Logo), slog.Any("err", err))
						return
					}

					imageData, err := converter.DecodeAny(data)
					if err != nil {
						slox.Error(ctx, "failed decoding image", slog.String("mod_id", m.ID), slog.String("logo", m.Logo), slog.Any("err", err))
						return
					}

					hash := thumbhash.EncodeImage(imageData)
					thumbHash := base64.StdEncoding.EncodeToString(hash)

					if err := m.Update().SetLogoThumbhash(thumbHash).Exec(ctx); err != nil {
						slox.Error(ctx, "failed saving thumbhash", slog.String("mod_id", m.ID), slog.String("logo", m.Logo), slog.Any("err", err))
						return
					}
				})
			}

			// Calculate for all users
			users, err := db.From(ctx).User.Query().Select(user.FieldID, user.FieldAvatar).Where(user.AvatarThumbhashIsNil()).All(ctx)
			if err != nil {
				return err
			}

			for i, u := range users {
				if i%50 == 0 {
					slox.Info(ctx, "generated thumbhash for n users", slog.Int("n", i))
				}

				pool.Submit(func() {
					if u.Avatar == "" {
						return
					}

					resp, err := http.Get(u.Avatar)
					if err != nil {
						slox.Error(ctx, "invalid url", slog.String("user_id", u.ID), slog.String("avatar", u.Avatar), slog.Any("err", err))
						return
					}

					defer resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						slox.Error(ctx, "invalid url", slog.String("user_id", u.ID), slog.String("avatar", u.Avatar), slog.Any("err", err))
						return
					}

					data, err := io.ReadAll(resp.Body)
					if err != nil {
						slox.Error(ctx, "invalid url", slog.String("user_id", u.ID), slog.String("avatar", u.Avatar), slog.Any("err", err))
						return
					}

					imageData, err := converter.DecodeAny(data)
					if err != nil {
						slox.Error(ctx, "failed decoding image", slog.String("user_id", u.ID), slog.String("avatar", u.Avatar), slog.Any("err", err))
						return
					}

					hash := thumbhash.EncodeImage(imageData)
					thumbHash := base64.StdEncoding.EncodeToString(hash)

					if err := u.Update().SetAvatarThumbhash(thumbHash).Exec(ctx); err != nil {
						slox.Error(ctx, "failed saving thumbhash", slog.String("user_id", u.ID), slog.String("avatar", u.Avatar), slog.Any("err", err))
						return
					}
				})
			}

			return nil
		},
	)
}
