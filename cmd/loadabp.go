package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"gitlab.com/peragrin/api/db"
	"gitlab.com/peragrin/api/models"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx/types"
	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func loadABP() error {
	dbClient, err := db.Client(viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	storeClient, err := minio.New(viper.GetString("STORE_ENDPOINT"), viper.GetString("STORE_ACCESS_KEY"), viper.GetString("STORE_SECRET_KEY"), viper.GetBool("STORE_SECURE"))
	if err != nil {
		log.Fatal(err)
	}

	const bucket = "peragrin"
	const location = "us-east-1"
	var dir = viper.GetString("DIR")

	if err := storeClient.MakeBucket(bucket, location); err != nil {
		if exists, err := storeClient.BucketExists(bucket); err != nil || !exists {
			return err
		}
	}

	if _, err := dbClient.Exec(`
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
			DELETE FROM AccountPromotion;
	`); err != nil {
		return err
	}

	kathleen := &models.Account{Email: "kathleen@billkaelin.com"}
	if err := kathleen.SetPassword("password"); err != nil {
		return err
	}
	if err := kathleen.Save(dbClient); err != nil {
		return err
	}

	atlantaBeltLinePartnership := &models.Organization{
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
	if err := atlantaBeltLinePartnership.Create(kathleen.ID, dbClient); err != nil {
		return err
	}
	if logo := atlantaBeltLinePartnership.Logo; logo != "" {
		file, err := os.Open(path.Join(dir, "logos", logo))
		if err != nil {
			return err
		}
		defer file.Close()
		if err := atlantaBeltLinePartnership.UploadLogo(file, storeClient); err != nil {
			return err
		}
	}
	if icon := atlantaBeltLinePartnership.Icon; icon != "" {
		file, err := os.Open(path.Join(dir, "icons", icon))
		if err != nil {
			return err
		}
		defer file.Close()
		if err := atlantaBeltLinePartnership.UploadIcon(file, storeClient); err != nil {
			return err
		}
	}

	atlantaBeltLine := &models.Community{
		Name: "Atlanta BeltLine",
		Lon:  -84.3669705,
		Lat:  33.7561718,
		Zoom: 12,
	}
	if err := atlantaBeltLine.Create(atlantaBeltLinePartnership.ID, dbClient); err != nil {
		return err
	}

	trekker := &models.Membership{Name: "Trekker", Description: "Trekker"}
	explorer := &models.Membership{Name: "Explorer", Description: "Explorer"}
	pathfinder := &models.Membership{Name: "Pathfinder", Description: "Pathfinder"}
	railrunner := &models.Membership{Name: "Railrunner", Description: "Railrunner"}
	groundbreaker := &models.Membership{Name: "Groundbreaker", Description: "Groundbreaker"}
	trailblazer := &models.Membership{Name: "Trailblazer", Description: "Trailblazer"}
	bridgebuilder := &models.Membership{Name: "Bridgebuilder", Description: "Bridgebuilder"}
	memberships := []*models.Membership{trekker, explorer, pathfinder, railrunner, groundbreaker, trailblazer, bridgebuilder}
	for _, membership := range memberships {
		if err := membership.Save(atlantaBeltLine.ID, dbClient); err != nil {
			return err
		}
	}

	geoJSONOverlays := []*models.GeoJSONOverlay{
		&models.GeoJSONOverlay{
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
		&models.GeoJSONOverlay{
			Name:  "belt-line.geojson",
			Style: types.JSONText([]byte(`{"base": {"weight": 5, "color": "#726a6a"}}`)),
		},
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
		if err := overlay.Create(atlantaBeltLine.ID, dbClient); err != nil {
			return err
		}
	}

	file, err := os.Open(path.Join(dir, "csv", "abp.csv"))
	if err != nil {
		return err
	}
	defer file.Close()

	r := csv.NewReader(file)

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	for _, row := range records[1:] {
		latlon := strings.Split(row[7], ",")
		lat, _ := strconv.ParseFloat(latlon[0], 64)
		lon, _ := strconv.ParseFloat(latlon[1], 64)
		o := models.Organization{
			Category: row[0],
			Name:     row[3],
			Address: models.Address{
				Street:  row[8],
				City:    row[9],
				State:   row[10],
				Zip:     row[11],
				Country: "United States of America",
			},
			Lat:     lat,
			Lon:     lon,
			Email:   row[12],
			Phone:   row[13],
			Website: row[14],
			Logo:    fmt.Sprintf("%s.png", strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.ToLower(strings.TrimSpace(row[3])), " ", "-", -1), "é", "e", -1), "&", "and", -1), ".", "", -1), "(", "", -1), ")", "", -1), "'", "", -1), "/", "", -1), ":", "", -1), "@", "at", -1)),
		}

		operator := models.Account{Email: o.Email}
		a, _ := models.GetAccountByEmail(o.Email, dbClient)
		if a == nil {
			operator.SetPassword("password")
			if err := operator.Save(dbClient); err != nil {
				return err
			}
		} else {
			operator = *a
		}

		if err := o.Create(operator.ID, dbClient); err != nil {
			return err
		}
		co := models.CommunityOrganization{CommunityID: atlantaBeltLine.ID, OrganizationID: o.ID}
		if err := co.Create(dbClient); err != nil {
			return err
		}

		if logo := o.Logo; logo != "" {
			file, err := os.Open(path.Join(dir, "logos", logo))
			if err != nil {
				return err
			}
			defer file.Close()
			if err := o.UploadLogo(file, storeClient); err != nil {
				return err
			}
		}

		if strings.ToLower(row[18]) != "unclaimed" {
			h := models.Hours{}
			for _, hours := range strings.Split(row[18], "\n") {
				if string(hours) == "" {
					continue
				}
				fields := strings.Split(hours, ", ")
				weekday, _ := strconv.Atoi(string(fields[0]))
				start, _ := strconv.Atoi(string(fields[1]))
				close, _ := strconv.Atoi(string(fields[2]))

				h = append(h, models.Hour{Weekday: time.Weekday(weekday), Start: start, Close: close})
			}
			if err := h.Set(o.ID, dbClient); err != nil {
				return err
			}
		}

		promotionNames := strings.Split(row[4], "\n")
		promotionDescriptions := strings.Split(row[5], "\n")
		promotionExclusions := strings.Split(row[6], "\n")
		for i := range promotionNames {
			promotion := models.Promotion{
				OrganizationID: o.ID,
				Name:           promotionNames[i],
			}
			if len(promotionDescriptions) >= (i + 1) {
				promotion.Description = promotionDescriptions[i]
			}
			if len(promotionExclusions) >= (i + 1) {
				promotion.Exclusions = promotionExclusions[i]
			}
			if err := promotion.Save(dbClient); err != nil {
				return err
			}
		}

		fmt.Printf("+%v\n", o)
	}

	return nil
}

// LoadABP is a cobra command that loads a given CSV files data into the database.
var LoadABP *cobra.Command

func init() {
	LoadABP = &cobra.Command{
		Use: "loadabp",
		Run: func(_ *cobra.Command, args []string) {
			log.Info("load abp - this will remove all data from the database")

			if err := loadABP(); err != nil {
				log.Fatal(err)
			}

			log.Info("completed successfully")
		},
	}

	LoadABP.PersistentFlags().StringP("dir", "", "", "absolute file path to the data dir")
	viper.BindPFlag("DIR", LoadABP.PersistentFlags().Lookup("dir"))
}
