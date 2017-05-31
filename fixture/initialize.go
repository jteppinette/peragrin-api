package fixture

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	minio "github.com/minio/minio-go"
	"gitlab.com/peragrin/api/common"
	"gitlab.com/peragrin/api/models"
)

var (
	jteppinette = &models.Account{Email: "jteppinette@jteppinette.com"}
	sajohnson   = &models.Account{Email: "sajohnson@sajohnson.com"}
	bjones      = &models.Account{Email: "bjones@bjones.com"}

	midtown = &models.Community{Name: "Midtown Atlanta"}
	decatur = &models.Community{Name: "Decatur Georgia"}

	midtownAtlantaChamber = &models.Organization{
		Name: "Metro Atlanta Chamber",
		Address: models.Address{
			Street: "191 Peachtree St NE #3400", City: "Atlanta", State: "GA", Country: "United States", Zip: "30303",
		},
		Lon:     -84.38642,
		Lat:     33.759115,
		Email:   "contact@metroatlantachamber.com",
		Phone:   "(678) 390-2910",
		Website: "https://midtownatlanta.com",
		Logo:    "metro-atlanta-chamber.png",
	}
	bobbyDoddStadium = &models.Organization{
		Name: "Bobby Dodd Stadium",
		Address: models.Address{
			Street: "North Avenue NW", City: "Atlanta", State: "GA", Country: "United States", Zip: "30313",
		},
		Lon:     -84.3903448,
		Lat:     33.7712937,
		Email:   "contact-us@bobby-dodd-stadium.com",
		Phone:   "(770) 320-3202",
		Website: "gt.edu",
		Logo:    "bobby-dodd-stadium.png",
	}
	emoryPublix = &models.Organization{
		Name: "Publix Super Market at Emory Commons",
		Address: models.Address{
			Street: "2155 N Decatur Rd", City: "Decatur", State: "GA", Country: "United States", Zip: "30033",
		},
		Lon:     -84.30444,
		Lat:     33.79023,
		Email:   "contact@publix.com",
		Phone:   "(770) 402-2309",
		Website: "publix.com",
	}

	accounts = []*models.Account{
		jteppinette, sajohnson, bjones,
	}

	communities = []*models.Community{
		midtown, decatur,
	}

	organizations = []*models.Organization{
		midtownAtlantaChamber, bobbyDoddStadium, emoryPublix,
	}

	operators = []struct {
		account      *models.Account
		organization *models.Organization
	}{
		{jteppinette, midtownAtlantaChamber},
		{sajohnson, bobbyDoddStadium},
		{bjones, emoryPublix},
	}

	memberships = []struct {
		community       *models.Community
		organization    *models.Organization
		isAdministrator bool
	}{
		{midtown, midtownAtlantaChamber, true},
		{midtown, bobbyDoddStadium, false},
		{decatur, emoryPublix, true},
	}

	posts = []struct {
		organization *models.Organization
		items        []*models.Post
	}{
		{midtownAtlantaChamber, []*models.Post{{Content: "Hey! This is your midtown atlanta chamber, live on Peragrin!"}}},
		{bobbyDoddStadium, []*models.Post{{Content: "Hey everyone! Come down to Bobby Dodd for the game today!"}, {Content: "Thats a win!"}}},
		{emoryPublix, []*models.Post{{Content: "We have great specials today on subs! Come check it out!"}}},
	}

	promotions = []struct {
		organization *models.Organization
		items        []*models.Promotion
	}{
		{
			bobbyDoddStadium, []*models.Promotion{
				{Name: "10% Off Food Purchases", Description: "All food purchases will be 10% off for members!"},
				{Name: "5% Off Jerseys", Description: "Home team jerseys will be 10% off for members!", Exclusions: "Seasons Ticket Holder Required", Expiration: common.JSONNullTime{pq.NullTime{Time: time.Now().AddDate(0, 3, 0), Valid: true}}, IsSingleUse: true},
			},
		},
		{
			emoryPublix, []*models.Promotion{
				{Name: "15% Off Publix Subs", Description: "All members, come enjoy 15% off our subs during this limited time offer!", Expiration: common.JSONNullTime{pq.NullTime{Time: time.Now().AddDate(0, 0, 7), Valid: true}}},
			},
		},
	}

	hours = models.Hours{
		{Weekday: time.Sunday, Start: 900, Close: 1700},
		{Weekday: time.Monday, Start: 900, Close: 1700},
		{Weekday: time.Tuesday, Start: 900, Close: 1700},
		{Weekday: time.Wednesday, Start: 900, Close: 1700},
		{Weekday: time.Thursday, Start: 900, Close: 1700},
	}
)

// Initialize loads fixture data intot the current database. Any previously
// uploaded data will be deleted.
func Initialize(db *sqlx.DB, store *minio.Client, dir string) error {

	const bucket = "peragrin"
	const location = "us-east-1"

	if err := store.MakeBucket(bucket, location); err != nil {
		if exists, err := store.BucketExists(bucket); err != nil || !exists {
			return err
		}
	}

	if _, err := db.Exec(`
			DELETE FROM Community;
			DELETE FROM Organization;
			DELETE FROM Account;
			DELETE FROM Post;
			DELETE FROM CommunityOrganization;
			DELETE FROM AccountOrganization;
	`); err != nil {
		return err
	}

	for i := range accounts {
		account := accounts[i]
		if err := account.SetPassword(strings.Split(account.Email, "@")[0]); err != nil {
			return err
		}
		if err := account.Save(db); err != nil {
			return err
		}
	}

	for _, operator := range operators {
		if err := operator.organization.Create(operator.account.ID, db); err != nil {
			return err
		}
		if err := hours.Set(operator.organization.ID, db); err != nil {
			return err
		}

		if logo := operator.organization.Logo; logo != "" {
			file, err := os.Open(path.Join(dir, logo))
			if err != nil {
				return err
			}
			defer file.Close()
			if err := operator.organization.UploadLogo(file, store); err != nil {
				return err
			}
		}
	}

	for _, membership := range memberships {
		if membership.isAdministrator {
			if err := membership.community.Create(membership.organization.ID, db); err != nil {
				return err
			}
		} else {
			co := models.CommunityOrganization{CommunityID: membership.community.ID, OrganizationID: membership.organization.ID}
			if err := co.Create(db); err != nil {
				return err
			}
		}
	}

	for _, post := range posts {
		for _, item := range post.items {
			item.OrganizationID = post.organization.ID
			if err := item.Save(db); err != nil {
				return err
			}
		}
	}

	for _, promotion := range promotions {
		for _, item := range promotion.items {
			item.OrganizationID = promotion.organization.ID
			if err := item.Save(db); err != nil {
				return err
			}
		}
	}

	return nil
}
