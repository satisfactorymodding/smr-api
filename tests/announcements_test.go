package tests

import (
	"strconv"
	"testing"

	"github.com/MarvinJWendt/testza"

	"github.com/satisfactorymodding/smr-api/config"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/migrations"
)

func init() {
	migrations.SetMigrationDir("../migrations")
	config.SetConfigDir("../")
	postgres.EnableDebug()
	db.EnableDebug()
}

func TestAnnouncements(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	// Run Twice to detect any cache issues
	for i := 0; i < 2; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var announcementID string

			t.Run("Create", func(t *testing.T) {
				createAnnouncement := authRequest(`mutation {
					createAnnouncement(announcement: {
						importance: Alert,
						message: "Hello World"
					}) {
						id
					}
				}`, token)

				var createAnnouncementResponse struct {
					CreateAnnouncement generated.Announcement
				}
				testza.AssertNoError(t, client.Run(ctx, createAnnouncement, &createAnnouncementResponse))
				testza.AssertNotEqual(t, "", createAnnouncementResponse.CreateAnnouncement.ID)

				announcementID = createAnnouncementResponse.CreateAnnouncement.ID
			})

			t.Run("Query One", func(t *testing.T) {
				queryAnnouncement := authRequest(`query ($id: AnnouncementID!) {
					getAnnouncement(announcementId: $id) {
						id
						message
						importance
					}
				}`, token)
				queryAnnouncement.Var("id", announcementID)

				var queryAnnouncementResponse struct {
					GetAnnouncement generated.Announcement
				}
				testza.AssertNoError(t, client.Run(ctx, queryAnnouncement, &queryAnnouncementResponse))
				testza.AssertEqual(t, announcementID, queryAnnouncementResponse.GetAnnouncement.ID)
				testza.AssertEqual(t, "Hello World", queryAnnouncementResponse.GetAnnouncement.Message)
				testza.AssertEqual(t, generated.AnnouncementImportanceAlert, queryAnnouncementResponse.GetAnnouncement.Importance)
			})

			t.Run("Update", func(t *testing.T) {
				updateAnnouncement := authRequest(`mutation ($id: AnnouncementID!) {
					updateAnnouncement(
						announcementId: $id,
						announcement: {
							importance: Fix,
							message: "Foo Bar"
						}
					) {
						id
					}
				}`, token)
				updateAnnouncement.Var("id", announcementID)

				var updateAnnouncementResponse struct {
					UpdateAnnouncement generated.Announcement
				}
				testza.AssertNoError(t, client.Run(ctx, updateAnnouncement, &updateAnnouncementResponse))
			})

			t.Run("Query Many", func(t *testing.T) {
				queryAnnouncements := authRequest(`query {
					getAnnouncements {
						id
						message
						importance
					}
				}`, token)

				var queryAnnouncementsResponse struct {
					GetAnnouncements []generated.Announcement
				}
				testza.AssertNoError(t, client.Run(ctx, queryAnnouncements, &queryAnnouncementsResponse))
				testza.AssertEqual(t, 1, len(queryAnnouncementsResponse.GetAnnouncements))
				testza.AssertEqual(t, announcementID, queryAnnouncementsResponse.GetAnnouncements[0].ID)
				testza.AssertEqual(t, "Foo Bar", queryAnnouncementsResponse.GetAnnouncements[0].Message)
				testza.AssertEqual(t, generated.AnnouncementImportanceFix, queryAnnouncementsResponse.GetAnnouncements[0].Importance)
			})

			t.Run("Query By Importance", func(t *testing.T) {
				getAnnouncementsByImportance := authRequest(`query {
					getAnnouncementsByImportance(importance: Info) {
						id
						message
						importance
					}
				}`, token)

				var getAnnouncementsByImportanceResponse struct {
					GetAnnouncements []generated.Announcement
				}
				testza.AssertNoError(t, client.Run(ctx, getAnnouncementsByImportance, &getAnnouncementsByImportanceResponse))
				testza.AssertEqual(t, 0, len(getAnnouncementsByImportanceResponse.GetAnnouncements))
			})

			t.Run("Delete", func(t *testing.T) {
				deleteAnnouncement := authRequest(`mutation ($id: AnnouncementID!) {
					deleteAnnouncement(announcementId: $id)
				}`, token)
				deleteAnnouncement.Var("id", announcementID)

				var deleteAnnouncementResponse struct {
					DeleteAnnouncement bool
				}
				testza.AssertNoError(t, client.Run(ctx, deleteAnnouncement, &deleteAnnouncementResponse))
				testza.AssertTrue(t, deleteAnnouncementResponse.DeleteAnnouncement)
			})
		})
	}
}
