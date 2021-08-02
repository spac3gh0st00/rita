package explodeddns

import (
	"github.com/activecm/rita/resources"
	"github.com/globalsign/mgo/bson"
)

//Results returns hostnames and their subdomain/ lookup statistics from the database.
//limit and noLimit control how many results are returned.
func Results(res *resources.Resources, limit int, noLimit bool) ([]Result, error) {
	ssn := res.DB.Session.Copy()
	defer ssn.Close()

	var explodedDNSResults []Result

	explodedDNSQuery := []bson.M{
		{"$unwind": "$dat"},
		{"$project": bson.M{"domain": 1, "subdomain_count": 1, "visited": "$dat.visited"}},
		{"$group": bson.M{
			"_id":             "$domain",
			"visited":         bson.M{"$sum": "$visited"},
			"subdomain_count": bson.M{"$first": "$subdomain_count"},
		}},
		{"$project": bson.M{
			"_id":             0,
			"domain":          "$_id",
			"visited":         1,
			"subdomain_count": 1,
		}},
		{"$sort": bson.M{"visited": -1}},
		{"$sort": bson.M{"subdomain_count": -1}},
	}

	if !noLimit {
		explodedDNSQuery = append(explodedDNSQuery, bson.M{"$limit": limit})
	}

	err := ssn.DB(res.DB.GetSelectedDB()).C(res.Config.T.DNS.ExplodedDNSTable).Pipe(explodedDNSQuery).AllowDiskUse().All(&explodedDNSResults)

	return explodedDNSResults, err

}
