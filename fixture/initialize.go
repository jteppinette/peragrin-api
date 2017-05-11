package fixture

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"gitlab.com/peragrin/api/models"
)

var (
	jteppinette = &models.Account{Email: "jteppinette@jteppinette.com"}
	sajohnson   = &models.Account{Email: "sajohnson@sajohnson.com"}
	bjones      = &models.Account{Email: "bjones@bjones.com"}

	midtown = &models.Community{Name: "Midtown Atlanta"}
	decatur = &models.Community{Name: "Decatur Georgia"}

	midtownAtlantaChamber = &models.Organization{Name: "Metro Atlanta Chamber", Address: "191 Peachtree Tower, 191 Peachtree St NE #3400, Atlanta, GA 30303", Longitude: -84.38642, Latitude: 33.759115}
	bobbyDoddStadium      = &models.Organization{Name: "Bobby Dodd Stadium", Address: "North Avenue NW, Atlanta, GA 30313", Longitude: -84.3903448, Latitude: 33.7712937}
	emoryPublix           = &models.Organization{Name: "Publix Super Market at Emory Commons", Address: "2155 N Decatur Rd, Decatur, GA 30033", Longitude: -84.30444, Latitude: 33.79023}

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
		{decatur, emoryPublix, false},
	}

	posts = []struct {
		organization *models.Organization
		items        []*models.Post
	}{
		{midtownAtlantaChamber, []*models.Post{{Content: "Hey! This is your midtown atlanta chamber, live on Peragrin!"}}},
		{bobbyDoddStadium, []*models.Post{{Content: "Hey everyone! Come down to Bobby Dodd for the game today!"}, {Content: "Thats a win!"}}},
		{emoryPublix, []*models.Post{{Content: "We have great specials today on subs! Come check it out!"}}},
	}
)

// Initialize loads fixture data intot the current database. Any previously
// uploaded data will be deleted.
func Initialize(client *sqlx.DB) error {

	// Clear out the database.
	if _, err := client.Exec(`
			DELETE FROM Community;
			DELETE FROM Organization;
			DELETE FROM Account;
			DELETE FROM Post;
			DELETE FROM Membership;
			DELETE FROM Operator;
	`); err != nil {
		return err
	}

	for i := range communities {
		if err := communities[i].Save(client); err != nil {
			return err
		}
		log.Infof("Created Community: %+v", communities[i])
	}
	for i := range accounts {
		account := accounts[i]
		if err := account.SetPassword(strings.Split(account.Email, "@")[0]); err != nil {
			return err
		}
		if err := account.Save(client); err != nil {
			return err
		}
		log.Infof("Created Account: %+v", account)
	}
	for i := range organizations {
		if err := organizations[i].Save(client); err != nil {
			return err
		}
		log.Infof("Created Organization: %+v", organizations[i])
	}

	for _, membership := range memberships {
		if err := membership.organization.AddMembership(membership.community.ID, membership.isAdministrator, client); err != nil {
			return err
		}
	}

	for _, operator := range operators {
		if err := operator.account.AddOperator(operator.organization.ID, client); err != nil {
			return err
		}
	}

	for _, post := range posts {
		for _, item := range post.items {
			item.OrganizationID = post.organization.ID
			if err := item.Save(client); err != nil {
				return err
			}
			log.Infof("Created Post: %+v", item)
		}
	}

	return nil
}
