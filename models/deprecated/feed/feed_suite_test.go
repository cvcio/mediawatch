package feed_test

import (
	"context"
	"testing"
	"time"

	"github.com/cvcio/mediawatch/models/deprecated/feed"

	"github.com/cvcio/mediawatch/internal/tests"
	"github.com/cvcio/mediawatch/pkg/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var test *tests.Test

var dbConn *db.MongoDB

func TestFeed(t *testing.T) {
	test = tests.New(true)
	dbConn = test.DB
	defer test.TearDown()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Feed DB Model Suite")
}

var _ = Describe("Feed", func() {
	var (
		ctx      context.Context
		now      time.Time
		newf     *feed.Feed
		newf1    *feed.Feed
		newf2    *feed.Feed
		createdF *feed.Feed
		err      error
	)
	BeforeSuite(func() {
		ctx = context.Background()
		now = time.Now()
		newf = &feed.Feed{
			Name:       "Demo Feed",
			Country:    "Greece",
			ScreenName: "demo",
			Status:     "active",
		}
		newf1 = &feed.Feed{
			Name:       "Demo Feed 1",
			Country:    "Greece",
			ScreenName: "demo1",
			Status:     "active",
		}
		newf2 = &feed.Feed{
			Name:       "Demo Feed 2",
			Country:    "Greece",
			ScreenName: "demo2",
			Status:     "active",
		}
	})

	Describe("Create Feed", func() {
		It("should be able to create a feed", func() {
			createdF, err = feed.Create(ctx, dbConn, newf, now)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should be able to create a feed 1", func() {
			_, err = feed.Create(ctx, dbConn, newf1, now)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should be able to create a feed 2", func() {
			_, err = feed.Create(ctx, dbConn, newf2, now)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Get Geed", func() {
		It("should be able to get feed by id", func() {
			_, err = feed.Get(ctx, dbConn, createdF.ID.Hex())
			Expect(err).NotTo(HaveOccurred())
		})
		It("should NOT be able to get feed by invalid id", func() {
			_, err = feed.Get(ctx, dbConn, "asdfasdfasdf")
			Expect(err).To(HaveOccurred())
		})
		It("should be able to get feed by screenname", func() {
			_, err = feed.ByScreenName(ctx, dbConn, createdF.ScreenName)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should not be able to get feed by wrong screenname", func() {
			_, err = feed.ByScreenName(ctx, dbConn, "demoNotExist")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("List feeds", func() {
		It("should be able to get feeds list", func() {
			feeds, err := feed.List(ctx, dbConn)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(feeds.Data)).Should(BeNumerically(">", 2))
		})
		It("should be able to get feeds list with limit", func() {
			feeds, err := feed.List(ctx, dbConn, feed.Limit(1))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(feeds.Data)).Should(BeNumerically("==", 1))
		})
		It("should be able to get feeds list with skip", func() {
			feeds, err := feed.List(ctx, dbConn, feed.Offset(1))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(feeds.Data)).Should(BeNumerically("==", 2))
		})
	})
	Describe("Delete Geed", func() {
		It("should be able to delete feed ", func() {
			err = feed.Delete(ctx, dbConn, createdF.ID.Hex())
			Expect(err).NotTo(HaveOccurred())
		})
		It("should NOT be able to get deleted feed by id", func() {
			_, err = feed.Get(ctx, dbConn, createdF.ID.Hex())
			Expect(err).To(HaveOccurred())
		})
		It("should NOT be able to get feed by screenname", func() {
			_, err = feed.ByScreenName(ctx, dbConn, createdF.ScreenName)
			Expect(err).To(HaveOccurred())
		})
		It("should NOT be able to delete feed by invalid id", func() {
			err = feed.Delete(ctx, dbConn, "asdfasdf")
			Expect(err).To(HaveOccurred())
		})
		It("should NOT be able to get feeds list with deleted", func() {
			feeds, err := feed.List(ctx, dbConn)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(feeds.Data)).Should(BeNumerically("==", 2))
		})
		It("should be able to get deleted feeds list only", func() {
			feeds, err := feed.List(ctx, dbConn, feed.Deleted())
			Expect(err).NotTo(HaveOccurred())
			Expect(len(feeds.Data)).Should(BeNumerically("==", 1))
		})
		It("should be able to get deleted feed bu screenname", func() {
			_, err := feed.ByScreenName(ctx, dbConn, "demo", feed.Deleted())
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
