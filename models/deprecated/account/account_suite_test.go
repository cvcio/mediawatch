package account_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cvcio/mediawatch/internal/tests"
	"github.com/cvcio/mediawatch/models/deprecated/account"
	"github.com/cvcio/mediawatch/pkg/auth"
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
	RunSpecs(t, "Account DB Model Suite")
}

var _ = Describe("Accounts - Admin actions", func() {
	var (
		newAdminUser account.NewAccount
		newUser      account.NewAccount
		userUser     *account.Account
		// userUser     *account.Account
		updU   account.UpdateAccount
		now    time.Time
		ctx    context.Context
		err    error
		savedU *account.Account
	)

	BeforeSuite(func() {
		ctx = context.Background()
		newAdminUser = account.NewAccount{
			Email: "press@mediawatch.io",
			//Roles:           []string{auth.RoleAdmin},
			//Organization:    "org1",
			Password:        "mediawatch",
			PasswordConfirm: "mediawatch",
		}
		newUser = account.NewAccount{
			Email:           "user@mediawatch.io",
			Password:        "mediawatch",
			PasswordConfirm: "mediawatch",
		}
		now = time.Now()
		updU = account.UpdateAccount{
			FirstName:       tests.StringPointer("FirstName"),
			LastName:        tests.StringPointer("LastName"),
			Email:           tests.StringPointer("user2@mediawatch.io"),
			Organization:    tests.StringPointer("Upd org"),
			Password:        tests.StringPointer("newpass1"),
			PasswordConfirm: tests.StringPointer("newpass1"),
		}
	})

	Describe("Create user", func() {
		It("should be able to create an admin user", func() {
			_, err = account.Create(ctx, dbConn, &newAdminUser, now)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should be able to create an admin user", func() {
			userUser, err = account.Create(ctx, dbConn, &newUser, now)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Get User", func() {
		It("should be able to get a user by id", func() {
			savedU, err := account.Get(ctx, dbConn, userUser.ID.Hex())
			Expect(err).NotTo(HaveOccurred())
			Expect(userUser.Email).To(Equal(savedU.Email))
		})
		It("should not be able to get user with invalid id", func() {
			_, err := account.Get(ctx, dbConn, "asdfasdf")
			Expect(err).To(HaveOccurred())
		})

		It("should be able to get a user by email", func() {
			savedU, err = account.ByEmail(ctx, dbConn, userUser.Email)
			Expect(err).NotTo(HaveOccurred())
			Expect(userUser.Email).To(Equal(savedU.Email))
		})
		It("should not be able to get a non existing user by email", func() {
			_, err := account.ByEmail(ctx, dbConn, "non@user.com")
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("Update User", func() {
		It("should be able to update a user's email and name", func() {
			nowUpd := time.Now()
			err := account.Update(ctx, dbConn, savedU.ID.Hex(), &updU, nowUpd)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should be able to see updates to mail and name", func() {
			updated, err := account.Get(ctx, dbConn, savedU.ID.Hex())
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Email).To(Equal(*updU.Email))
		})
	})
	Describe("List users", func() {
		It("should be able to get user's list", func() {
			_, err := account.List(ctx, dbConn)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("Delete User", func() {
		It("should be able to delete user", func() {
			err := account.Delete(ctx, dbConn, userUser.ID.Hex())
			Expect(err).NotTo(HaveOccurred())
		})
		It("should not be able to get user back", func() {
			_, err := account.Get(ctx, dbConn, userUser.ID.Hex())
			Expect(err).To(HaveOccurred())
		})
	})
})

// token generator in a specific way.
type mockTokenGenerator struct{}

// GenerateToken implements the TokenGenerator interface. It returns a "token"
// that includes some information about the claims it was passed.
func (mockTokenGenerator) GenerateToken(claims auth.Claims) (string, error) {
	return fmt.Sprintf("sub:%q iss:%d", claims.Subject, claims.IssuedAt), nil
}

func (mockTokenGenerator) ParseClaims(tknStr string) (auth.Claims, error) {

	return auth.Claims{}, nil
}

var _ = Describe("Accounts - User actions", func() {
	var (
		newUser  account.NewAccount
		userUser *account.Account
		now      time.Time
		ctx      context.Context
		err      error
		tknGen   mockTokenGenerator
	)

	BeforeEach(func() {
		now = time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)

		ctx = context.Background()
		newUser = account.NewAccount{
			Email:           "user@mediawatch.io",
			Password:        "mediawatch",
			PasswordConfirm: "mediawatch",
		}
		// now = time.Now()
	})
	Describe("Register user", func() {
		It("should be able to register", func() {
			userUser, err = account.Create(ctx, dbConn, &newUser, now)
			Expect(err).NotTo(HaveOccurred())
			Expect(userUser.CreatedAt).To(Equal(now))
		})
	})
	Describe("Authenticate User", func() {
		It("should not be able to authenticate if user not exist", func() {
			_, err := account.ByEmail(ctx, dbConn, "notexist@user.com")
			Expect(err).To(HaveOccurred())
		})
		It("should not be able to authenticate with wrong password", func() {
			gotUser, err := account.ByEmail(ctx, dbConn, newUser.Email)
			Expect(err).NotTo(HaveOccurred())
			err = account.PasswordOK(ctx, gotUser, "wrong password")
			Expect(err).To(HaveOccurred())
		})
		It("should be able to authenticate and get a token", func() {
			gotUser, err := account.ByEmail(ctx, dbConn, newUser.Email)
			Expect(err).NotTo(HaveOccurred())
			err = account.PasswordOK(ctx, gotUser, newUser.Password)
			Expect(err).NotTo(HaveOccurred())
			tkn, err := account.Authenticate(ctx, tknGen, now, gotUser)
			Expect(err).NotTo(HaveOccurred())
			Expect(tkn.AccessToken).To(Equal(fmt.Sprintf("sub:%q iss:1538352000", gotUser.ID.Hex())))
		})
	})
})
