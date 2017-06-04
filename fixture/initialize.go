package fixture

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	minio "github.com/minio/minio-go"
	"gitlab.com/peragrin/api/models"
)

var (
	kCone     = &models.Account{Email: "kathleen@billkaelin.com"}
	gYeremian = &models.Account{Email: "gilbert@communitashospitality.com"}
	sDoty     = &models.Account{Email: "shaun@bantamandbiddy.com"}
	kPeak     = &models.Account{Email: "kevin@peak.com"}
	gCameli   = &models.Account{Email: "george@cameli.com"}
	kWalker   = &models.Account{Email: "kelsey&walker.com"}
	jDelp     = &models.Account{Email: "jeff@fcsministries.org"}
	cBarrow   = &models.Account{Email: "clintbarrow@thiscompany.com"}
	tRogers   = &models.Account{Email: "tyler@kingofpops.net"}
	missy     = &models.Account{Email: "missy@missy.com"}
	brenda    = &models.Account{Email: "brenda@360media.net"}
	natasha   = &models.Account{Email: "natasha@vicbrands.com"}
	aSmith    = &models.Account{Email: "andrea@ladybirdatlanta.com"}
	anna      = &models.Account{Email: "anna@saviprovisions.com"}
	leah      = &models.Account{Email: "leah@kalemecrazy.net"}
	jamie     = &models.Account{Email: "jamie.saye26@gmail.com"}

	atlantaBeltLine = &models.Community{
		Name: "Atlanta BeltLine",
		Lon:  -84.3669705,
		Lat:  33.7561718,
		Zoom: 12,
	}

	trekker       = &models.Membership{Name: "Trekker", Description: "Trekker"}
	explorer      = &models.Membership{Name: "Explorer", Description: "Explorer"}
	pathfinder    = &models.Membership{Name: "Pathfinder", Description: "Pathfinder"}
	railrunner    = &models.Membership{Name: "Railrunner", Description: "Railrunner"}
	groundbreaker = &models.Membership{Name: "Groundbreaker", Description: "Groundbreaker"}
	trailblazer   = &models.Membership{Name: "Trailblazer", Description: "Trailblazer"}
	bridgebuilder = &models.Membership{Name: "Bridgebuilder", Description: "Bridgebuilder"}

	memberships = []struct {
		community *models.Community
		items     []*models.Membership
	}{
		{
			atlantaBeltLine,
			[]*models.Membership{trekker, explorer, pathfinder, railrunner, groundbreaker, trailblazer, bridgebuilder},
		},
	}

	patrons = []struct {
		membership *models.Membership
		models.Account
	}{
		{trekker, models.Account{Email: "josheppinette@josheppinette.com"}},
		{trekker, models.Account{Email: "blanchardsawyer@comveyer.com"}},
		{trekker, models.Account{Email: "dotsoncraig@comveyer.com"}},
		{trekker, models.Account{Email: "kayesanford@comveyer.com"}},
		{trekker, models.Account{Email: "glovertucker@comveyer.com"}},
		{trekker, models.Account{Email: "blackenglish@comveyer.com"}},
		{trekker, models.Account{Email: "gabriellenorton@comveyer.com"}},
		{trekker, models.Account{Email: "edwinavaldez@comveyer.com"}},
		{trekker, models.Account{Email: "holdencoffey@comveyer.com"}},
		{trekker, models.Account{Email: "bauercurtis@comveyer.com"}},
		{trekker, models.Account{Email: "dalemoran@comveyer.com"}},
		{trekker, models.Account{Email: "maureennoel@comveyer.com"}},
		{trekker, models.Account{Email: "cobbmarsh@comveyer.com"}},
		{trekker, models.Account{Email: "serenajones@comveyer.com"}},
		{trekker, models.Account{Email: "clementsfaulkner@comveyer.com"}},
		{trekker, models.Account{Email: "hendricksshepherd@comveyer.com"}},
		{trekker, models.Account{Email: "montoyadrake@comveyer.com"}},
		{trekker, models.Account{Email: "steelecastaneda@comveyer.com"}},
		{trekker, models.Account{Email: "saundershicks@comveyer.com"}},
		{trekker, models.Account{Email: "contrerasdejesus@comveyer.com"}},
		{trekker, models.Account{Email: "riddlekemp@comveyer.com"}},
		{trekker, models.Account{Email: "conradlawson@comveyer.com"}},
		{trekker, models.Account{Email: "howemcgee@comveyer.com"}},
		{trekker, models.Account{Email: "paulamartin@comveyer.com"}},
		{trekker, models.Account{Email: "mullenmoses@comveyer.com"}},
		{trekker, models.Account{Email: "jacobsjennings@comveyer.com"}},
		{trekker, models.Account{Email: "olsenpratt@comveyer.com"}},
		{trekker, models.Account{Email: "halemyers@comveyer.com"}},
		{trekker, models.Account{Email: "amelialindsey@comveyer.com"}},
		{trekker, models.Account{Email: "eulamills@comveyer.com"}},
		{trekker, models.Account{Email: "adrianhowe@comveyer.com"}},
		{trekker, models.Account{Email: "chasitybuchanan@comveyer.com"}},
		{trekker, models.Account{Email: "adasimpson@comveyer.com"}},
		{trekker, models.Account{Email: "tashaphelps@comveyer.com"}},
		{trekker, models.Account{Email: "hollieoliver@comveyer.com"}},
		{trekker, models.Account{Email: "warnerwoodard@comveyer.com"}},
		{trekker, models.Account{Email: "simonwall@comveyer.com"}},
		{trekker, models.Account{Email: "barkertaylor@comveyer.com"}},
		{trekker, models.Account{Email: "kirbyhaney@comveyer.com"}},
		{trekker, models.Account{Email: "minniepittman@comveyer.com"}},
		{trekker, models.Account{Email: "annaparrish@comveyer.com"}},
		{trekker, models.Account{Email: "hortonosborne@comveyer.com"}},
		{trekker, models.Account{Email: "kelseybartlett@comveyer.com"}},
		{trekker, models.Account{Email: "rowenaphillips@comveyer.com"}},
		{trekker, models.Account{Email: "mamiemarquez@comveyer.com"}},
		{trekker, models.Account{Email: "martinahuber@comveyer.com"}},
		{trekker, models.Account{Email: "jordanrichardson@comveyer.com"}},
		{trekker, models.Account{Email: "irmamcpherson@comveyer.com"}},
		{trekker, models.Account{Email: "evangelinecalhoun@comveyer.com"}},
		{trekker, models.Account{Email: "mitzieaton@comveyer.com"}},
	}

	geoJSONOverlays = []struct {
		community *models.Community
		models.GeoJSONOverlay
	}{
		{
			community: atlantaBeltLine,
			GeoJSONOverlay: models.GeoJSONOverlay{
				Name: "belt-line-zones.geojson",
				Style: types.JSONText([]byte(`
					{
						"property": "BPA_Segmen",
						"base": {"weight": 3, "color": "white", "opacity": 1, "fillOpacity": 0.3},
						"values": {
							"Northside": {"fillColor": "#00aef4"},
							"Northeast": {"fillColor": "#8bc932"},
							"Southeast": {"fillColor": "#0061c2"},
							"Southwest": {"fillColor": "#8bc932"},
							"Westside": {"fillColor": "#0061c2"}
						}
					}
				`)),
			},
		},
		{
			community: atlantaBeltLine,
			GeoJSONOverlay: models.GeoJSONOverlay{
				Name:  "belt-line.geojson",
				Style: types.JSONText([]byte(`{"base": {"weight": 5, "color": "#726a6a"}}`)),
			},
		},
	}

	atlantaBeltLinePartnership = &models.Organization{
		Name: "Atlanta BeltLine Partnership",
		Address: models.Address{
			Street: "112 Krog St NE #14", City: "Atlanta", State: "GA", Country: "United States", Zip: "30307",
		},
		Lon:     -84.3669705,
		Lat:     33.7561718,
		Email:   "info@atlbeltlinepartnership.org",
		Phone:   "(404) 446-4404",
		Website: "beltline.org",
		Icon:    "atlanta-belt-line-partnership.png",
		Logo:    "atlanta-belt-line-partnership.png",
	}

	tenthAndPiedmont = &models.Organization{
		Name: "10th & Piedmont",
		Address: models.Address{
			Street: "991 Piedmont Ave NE", City: "Atlanta", State: "GA", Country: "United States", Zip: "30309",
		},
		Lon:      -81.371266557291,
		Lat:      40.8058975055397,
		Email:    "info@10thandpiedmont.com",
		Phone:    "(404) 602-5510",
		Website:  "http://www.10thp.com/",
		Logo:     "tenth-and-piedmont.png",
		Category: models.Resturaunt,
	}
	bantamAndBiddy = &models.Organization{
		Name: "Bantam & Biddy",
		Address: models.Address{
			Street: "1544 Piedmont Ave NE #301", City: "Atlanta", State: "GA", Country: "United States", Zip: "30324",
		},
		Lon:      -84.3688771243521,
		Lat:      33.7984687036734,
		Email:    "shaun@bantamandbiddy.com",
		Phone:    "(404) 907-3469",
		Website:  "http://www.bantamandbiddy.com/",
		Logo:     "bantam-and-biddy.png",
		Category: models.Resturaunt,
	}
	cajaPopcorn = &models.Organization{
		Name: "CaJa Popcorn",
		Address: models.Address{
			Street: "2333 Peachtree Rd", City: "Atlanta", State: "GA", Country: "United States", Zip: "30305",
		},
		Lon:      -84.3696849,
		Lat:      33.844624,
		Email:    "contact@cajapopcorn.com",
		Phone:    "(404) 846-2156",
		Website:  "http://www.cajapopcorn.com/",
		Logo:     "caja-popcorn.png",
		Category: models.Resturaunt,
	}
	camelisPizza = &models.Organization{
		Name: "Cameli's Pizza",
		Address: models.Address{
			Street: "337 Moreland Ave NE", City: "Atlanta", State: "GA", Country: "United States", Zip: "30307",
		},
		Lon:      -84.3491539,
		Lat:      33.757426,
		Email:    "info@camelispizza.com",
		Phone:    "(404) 249-9020",
		Website:  "http://www.camelispizza.com/",
		Logo:     "camelis-pizza.png",
		Category: models.Resturaunt,
	}
	chickABiddy = &models.Organization{
		Name: "Chick-a-Biddy",
		Address: models.Address{
			Street: "264 19th St NW", City: "Atlanta", State: "GA", Country: "United States", Zip: "30363",
		},
		Lon:      -84.39713372,
		Lat:      33.79346104,
		Email:    "kelsey@lizlapiduspr.com",
		Phone:    "(404) 688-1466",
		Website:  "http://www.chickabiddyatl.com/",
		Logo:     "chick-a-biddy.png",
		Category: models.Resturaunt,
	}
	communityGroundsCoffeeshop = &models.Organization{
		Name: "Community Grounds Coffeeshop",
		Address: models.Address{
			Street: "1297 McDonough Blvd SE", City: "Atlanta", State: "GA", Country: "United States", Zip: "30315",
		},
		Lon:      -84.3829909,
		Lat:      33.717947,
		Email:    "jeff@fcsministries.org",
		Phone:    "(404) 586-0692",
		Website:  "https://communitygrounds.com",
		Logo:     "community-grounds-coffeeshop.png",
		Category: models.Resturaunt,
	}
	frogsCantina = &models.Organization{
		Name: "F.R.O.G.S. Cantina",
		Address: models.Address{
			Street: "931 Monroe Dr", City: "Atlanta", State: "GA", Country: "United States", Zip: "30308",
		},
		Lon:      -84.36819,
		Lat:      33.780192,
		Email:    "clintbarrow@thiscompany.com",
		Phone:    "(404) 607-9967",
		Website:  "http://www.frogsmidtown.com/",
		Logo:     "frogs-cantina.png",
		Category: models.Resturaunt,
	}
	gsMidtown = &models.Organization{
		Name: "G's Midtown",
		Address: models.Address{
			Street: "219 10th St NE", City: "Atlanta", State: "GA", Country: "United States", Zip: "30309",
		},
		Lon:      -84.3823024,
		Lat:      33.7816376,
		Email:    "gilbert@communitashospitality.com",
		Phone:    "(404) 872-8012",
		Website:  "http://www.gsmidtown.com/",
		Logo:     "gs-midtown.png",
		Category: models.Resturaunt,
	}

	accounts = []*models.Account{
		kCone, gYeremian, sDoty, kPeak, gCameli, kWalker, jDelp, cBarrow, tRogers, missy, brenda, natasha, aSmith, anna, leah, jamie,
	}

	communities = []*models.Community{
		atlantaBeltLine,
	}

	organizations = []*models.Organization{
		atlantaBeltLinePartnership, tenthAndPiedmont, bantamAndBiddy, cajaPopcorn, camelisPizza, chickABiddy, communityGroundsCoffeeshop, frogsCantina, gsMidtown,
	}

	operators = []struct {
		account      *models.Account
		organization *models.Organization
	}{
		{kCone, atlantaBeltLinePartnership},
		{gYeremian, tenthAndPiedmont},
		{sDoty, bantamAndBiddy},
		{kPeak, cajaPopcorn},
		{gCameli, camelisPizza},
		{kWalker, chickABiddy},
		{jDelp, communityGroundsCoffeeshop},
		{cBarrow, frogsCantina},
		{gYeremian, gsMidtown},
	}

	members = []struct {
		community       *models.Community
		organization    *models.Organization
		isAdministrator bool
	}{
		{atlantaBeltLine, atlantaBeltLinePartnership, true},
		{atlantaBeltLine, tenthAndPiedmont, false},
		{atlantaBeltLine, bantamAndBiddy, false},
		{atlantaBeltLine, cajaPopcorn, false},
		{atlantaBeltLine, camelisPizza, false},
		{atlantaBeltLine, chickABiddy, false},
		{atlantaBeltLine, communityGroundsCoffeeshop, false},
		{atlantaBeltLine, frogsCantina, false},
		{atlantaBeltLine, gsMidtown, false},
	}

	promotions = []struct {
		organization *models.Organization
		items        []*models.Promotion
	}{
		{
			tenthAndPiedmont, []*models.Promotion{
				{Name: "10% Off", Exclusions: "alcohol excluded, dinner and dine-in only, one discount per table, not combined with other offers, not valid for special events"},
			},
		},
		{
			bantamAndBiddy, []*models.Promotion{
				{Name: "10% Off Food Purchases", Exclusions: "alcohol excluded"},
			},
		},
		{
			cajaPopcorn, []*models.Promotion{
				{Name: "10% Off"},
			},
		},
		{
			camelisPizza, []*models.Promotion{
				{Name: "15% Off", Exclusions: "Dine in only, one discount per card, not combined with other offers."},
			},
		},
		{
			chickABiddy, []*models.Promotion{
				{Name: "10% Off"},
			},
		},
		{
			communityGroundsCoffeeshop, []*models.Promotion{
				{Name: "10% Off"},
			},
		},
		{
			frogsCantina, []*models.Promotion{
				{Name: "Buy One Get One 'BeltLine Margaritas'", Exclusions: "1 per table; 1 visit per day"},
			},
		},
		{
			gsMidtown, []*models.Promotion{
				{Name: "10% Off", Exclusions: "alcohol excluded, dinner and dine-in only, one discount per table, not combined with other offers, not valid for special events"},
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
			DELETE FROM Account;
			DELETE FROM Organization;
			DELETE FROM Hours;
			DELETE FROM Promotion;
			DELETE FROM Community;
			DELETE FROM GeoJSONOverlay;
			DELETE FROM Post;
			DELETE FROM AccountOrganization;
			DELETE FROM CommunityOrganization;
			DELETE FROM Membership;
			DELETE FROM AccountMembership;
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
			file, err := os.Open(path.Join(dir, "logos", logo))
			if err != nil {
				return err
			}
			defer file.Close()
			if err := operator.organization.UploadLogo(file, store); err != nil {
				return err
			}
		}
	}

	for _, member := range members {
		if member.isAdministrator {
			if err := member.community.Create(member.organization.ID, db); err != nil {
				return err
			}
			if icon := member.organization.Icon; icon != "" {
				file, err := os.Open(path.Join(dir, "icons", icon))
				if err != nil {
					return err
				}
				defer file.Close()
				if err := member.organization.UploadIcon(file, store); err != nil {
					return err
				}
			}
		} else {
			co := models.CommunityOrganization{CommunityID: member.community.ID, OrganizationID: member.organization.ID}
			if err := co.Create(db); err != nil {
				return err
			}
		}
	}

	for _, overlay := range geoJSONOverlays {
		file, err := os.Open(path.Join(dir, "geojson", overlay.Name))
		if err != nil {
			return err
		}
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&overlay.Data); err != nil {
			return err
		}
		if err := overlay.Create(overlay.community.ID, db); err != nil {
			return err
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

	for _, membership := range memberships {
		for _, item := range membership.items {
			if err := item.Save(membership.community.ID, db); err != nil {
				return err
			}
		}
	}

	for _, patron := range patrons {
		if err := patron.SetPassword(strings.Split(patron.Email, "@")[0]); err != nil {
			return err
		}
		if err := patron.CreateWithMembership(patron.membership.ID, db); err != nil {
			return err
		}
	}

	return nil
}
