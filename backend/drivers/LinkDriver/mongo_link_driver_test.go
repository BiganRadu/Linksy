package LinkDriver

import (
	"testing"
	"time"

	"backend/models"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func newTestDriver(coll *mongo.Collection) *MongoLinkDriver {
	return &MongoLinkDriver{
		collection: coll,
		nowFunc: func() time.Time {
			return time.Unix(0, 0)
		},
	}
}

func TestGetLinkByID_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("success", func(mt *mtest.T) {
		d := newTestDriver(mt.Coll)

		// FindOne returns a cursor with one doc.
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "links.links", mtest.FirstBatch,
				bson.D{
					{"_id", "id123"},
					{"title", "Test Title"},
					{"member_email", "m@e"},
				},
			),
		)

		got, err := d.GetLinkByID("id123")
		require.NoError(mt, err)
		require.NotNil(mt, got)
		require.Equal(mt, "id123", got.ID)
		require.Equal(mt, "Test Title", got.Title)
	})
}

func TestGetLinkByID_NotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("not_found", func(mt *mtest.T) {
		d := newTestDriver(mt.Coll)

		// Empty first batch -> ErrNoDocuments.
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "links.links", mtest.FirstBatch),
		)

		got, err := d.GetLinkByID("missing")
		require.Error(mt, err)
		require.Equal(mt, mongo.ErrNoDocuments, err)
		require.Nil(mt, got)
	})
}

func TestGetLinksForMember_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("member_links", func(mt *mtest.T) {
		d := newTestDriver(mt.Coll)

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "links.links", mtest.FirstBatch,
				bson.D{{"_id", "a"}, {"title", "A"}, {"member_email", "u@x"}},
				bson.D{{"_id", "b"}, {"title", "B"}, {"member_email", "u@x"}},
			),
			// Server closes cursor (empty next batch)
			mtest.CreateCursorResponse(0, "links.links", mtest.NextBatch),
		)

		got, err := d.GetLinksForMember("u@x")
		require.NoError(mt, err)
		require.Len(mt, got, 2)
		require.Equal(mt, "a", got[0].ID)
		require.Equal(mt, "b", got[1].ID)
	})
}

func TestGetLinksForMember_Error(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("member_links_error", func(mt *mtest.T) {
		d := newTestDriver(mt.Coll)

		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{Code: 123, Message: "find failed"},
			),
		)

		got, err := d.GetLinksForMember("u@x")
		require.Error(mt, err)
		require.Nil(mt, got)
	})
}

func TestGetQRsForMember_Filter(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("qrs_only", func(mt *mtest.T) {
		d := newTestDriver(mt.Coll)

		// Only docs with has_qr:true should be returned by the driver query.
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "links.links", mtest.FirstBatch,
				bson.D{{"_id", "qr1"}, {"title", "QR1"}, {"member_email", "u@x"}, {"has_qr", true}},
				bson.D{{"_id", "qr2"}, {"title", "QR2"}, {"member_email", "u@x"}, {"has_qr", true}},
			),
			mtest.CreateCursorResponse(0, "links.links", mtest.NextBatch),
		)

		got, err := d.GetQRsForMember("u@x")
		require.NoError(mt, err)
		require.Len(mt, got, 2)
		require.Equal(mt, "qr1", got[0].ID)
		require.Equal(mt, true, got[0].HasQR)
	})
}

func TestUpsertLink_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("upsert_ok", func(mt *mtest.T) {
		d := newTestDriver(mt.Coll)

		mt.AddMockResponses(
			// ReplaceOne response (update style)
			mtest.CreateSuccessResponse(
				bson.E{Key: "n", Value: 1},
				bson.E{Key: "nModified", Value: 1},
				bson.E{Key: "ok", Value: 1},
			),
		)

		link := &models.Link{ID: "idU", Title: "Upsert"}
		err := d.UpsertLink(link)
		require.NoError(mt, err)
	})
}

func TestUpsertLink_Error(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("upsert_err", func(mt *mtest.T) {
		d := newTestDriver(mt.Coll)

		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{Code: 50, Message: "write failed"},
			),
		)

		err := d.UpsertLink(&models.Link{ID: "x"})
		require.Error(mt, err)
	})
}

func TestDeleteLink_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("delete_ok", func(mt *mtest.T) {
		d := newTestDriver(mt.Coll)

		mt.AddMockResponses(
			mtest.CreateSuccessResponse(
				bson.E{Key: "n", Value: 1},
				bson.E{Key: "ok", Value: 1},
			),
		)

		err := d.DeleteLink("gone")
		require.NoError(mt, err)
	})
}

func TestDeleteLink_Error(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("delete_err", func(mt *mtest.T) {
		d := newTestDriver(mt.Coll)

		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{Code: 66, Message: "delete failed"},
			),
		)

		err := d.DeleteLink("gone")
		require.Error(mt, err)
	})
}
