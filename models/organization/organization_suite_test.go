package organization_test

import (
	"context"
	"testing"
	"time"

	"github.com/cvcio/mediawatch/models/organization"

	"github.com/cvcio/mediawatch/internal/tests"
	"github.com/cvcio/mediawatch/pkg/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var test *tests.Test

var dbConn *db.MongoDB

func TestAccount(t *testing.T) {
	test = tests.New(true)
	dbConn = test.DB
	defer test.TearDown()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Organization DB Models Suite")
}

var _ = Describe("Organization", func() {
	var (
		// newUser account.NewAccount
		// userUser *account.Account
		err     error
		ctx     context.Context
		now     time.Time
		newOrg  organization.Organization
		newOrg2 organization.Organization
		newOrg3 organization.Organization
		org     *organization.Organization
		// org2    *organization.Organization
		// org3    *organization.Organization
	)
	BeforeSuite(func() {
		ctx = context.Background()
		newOrg = organization.Organization{
			Name:    "Org1",
			Country: "Greece",
			Email:   "user@org1.com",
		}
		newOrg2 = organization.Organization{
			Name:    "Org2",
			Country: "Greece",
			Email:   "user@org2.com",
		}
		newOrg3 = organization.Organization{
			Name:    "Org3",
			Country: "Greece",
			Email:   "user@org3.com",
		}
		now = time.Now()
	})

	Describe("Create Organization", func() {
		It("should be able to create an organization", func() {
			org, err = organization.Create(ctx, dbConn, &newOrg, now)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should be able to create an organization", func() {
			_, err = organization.Create(ctx, dbConn, &newOrg2, now)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should be able to create an organization", func() {
			_, err = organization.Create(ctx, dbConn, &newOrg3, now)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("Get organization by ID", func() {
		It("should be able to get org by id", func() {
			_, err = organization.Get(ctx, dbConn, org.ID.Hex())
			Expect(err).NotTo(HaveOccurred())
		})
		It("should not be able to get org with invalid id", func() {
			_, err := organization.Get(ctx, dbConn, "asdfasdf")
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("Get list of organizations", func() {
		It("should be able to get a list of organizations", func() {
			_, err := organization.List(ctx, dbConn)
			Expect(err).NotTo(HaveOccurred())

		})
	})
})
