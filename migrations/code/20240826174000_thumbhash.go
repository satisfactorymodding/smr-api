package code

import (
	"context"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"

	"github.com/Vilsol/slox"
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

			// Calculate for all mods
			mods, err := db.From(ctx).Mod.Query().Select(mod.FieldID, mod.FieldLogo).Where(mod.LogoThumbhashIsNil()).All(ctx)
			if err != nil {
				return err
			}

			for _, m := range mods {
				if m.Logo == "" {
					continue
				}

				resp, err := http.Get(m.Logo)
				if err != nil {
					slox.Error(ctx, "invalid url", slog.String("mod_id", m.ID), slog.String("logo", m.Logo), slog.Any("err", err))
					continue
				}

				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					slox.Error(ctx, "invalid url", slog.String("mod_id", m.ID), slog.String("logo", m.Logo), slog.Any("err", err))
					continue
				}

				data, err := io.ReadAll(resp.Body)
				if err != nil {
					slox.Error(ctx, "invalid url", slog.String("mod_id", m.ID), slog.String("logo", m.Logo), slog.Any("err", err))
					continue
				}

				imageData, err := converter.DecodeAny(data)
				if err != nil {
					slox.Error(ctx, "failed decoding image", slog.String("mod_id", m.ID), slog.String("logo", m.Logo), slog.Any("err", err))
					continue
				}

				hash := thumbhash.EncodeImage(imageData)
				thumbHash := base64.StdEncoding.EncodeToString(hash)

				if err := m.Update().SetLogoThumbhash(thumbHash).Exec(ctx); err != nil {
					return err
				}
			}

			// Calculate for all users
			users, err := db.From(ctx).User.Query().Select(user.FieldID, user.FieldAvatar).Where(user.AvatarThumbhashIsNil()).All(ctx)
			if err != nil {
				return err
			}

			for _, u := range users {
				if u.Avatar == "" {
					continue
				}

				resp, err := http.Get(u.Avatar)
				if err != nil {
					slox.Error(ctx, "invalid url", slog.String("user_id", u.ID), slog.String("avatar", u.Avatar), slog.Any("err", err))
					continue
				}

				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					slox.Error(ctx, "invalid url", slog.String("user_id", u.ID), slog.String("avatar", u.Avatar), slog.Any("err", err))
					continue
				}

				data, err := io.ReadAll(resp.Body)
				if err != nil {
					slox.Error(ctx, "invalid url", slog.String("user_id", u.ID), slog.String("avatar", u.Avatar), slog.Any("err", err))
					continue
				}

				imageData, err := converter.DecodeAny(data)
				if err != nil {
					slox.Error(ctx, "failed decoding image", slog.String("user_id", u.ID), slog.String("avatar", u.Avatar), slog.Any("err", err))
					continue
				}

				hash := thumbhash.EncodeImage(imageData)
				thumbHash := base64.StdEncoding.EncodeToString(hash)

				if err := u.Update().SetAvatarThumbhash(thumbHash).Exec(ctx); err != nil {
					return err
				}
			}

			return nil
		},
	)
}
