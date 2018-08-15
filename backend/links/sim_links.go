package links

import (
	"encoding/json"
	"sort"

	"github.com/rs/xid"

	"github.com/sirupsen/logrus"

	"github.com/golovers/kiki/backend/db"
)

var linksDB db.Database

type simLinkSvc struct{}

func NewSimLinkSvc() LinkSvc {
	linksDB = db.Table("links")
	return &simLinkSvc{}
}

// Links return  all exiting links
func (svc *simLinkSvc) Links() []*QuickLink {
	rs := make([]*QuickLink, 0)
	for it := linksDB.NewIterator(); it.Next(); {
		var link QuickLink
		if err := json.Unmarshal(it.Value(), &link); err != nil {
			logrus.Errorf("failed to unmarshal val: %v", err)
		}
		rs = append(rs, &link)
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Visitted < rs[j].Visitted
	})
	for _, l := range rs {
		logrus.Infof("name: %s, link: %s", l.FullName(), l.Link)
	}
	return rs
}

//Add add the new links to existing list
func (svc *simLinkSvc) Add(link *QuickLink) error {
	link.ID = xid.New().String()
	data, err := json.Marshal(link)
	if err != nil {
		return err
	}
	return linksDB.Put([]byte(link.ID), data)
}

func (svc *simLinkSvc) Delete(id string) error {
	return linksDB.Delete([]byte(id))
}

func (svc *simLinkSvc) LinksByType(typ string) []*QuickLink {
	links := make([]*QuickLink, 0)
	for _, l := range svc.Links() {
		if l.Type == typ {
			links = append(links, l)
		}
	}
	return links
}
