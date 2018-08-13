package links

import (
	"encoding/json"
	"sort"

	"github.com/rs/xid"

	"github.com/sirupsen/logrus"

	"github.com/golovers/kiki/backend/db"
)

type simLinkSvc struct{}

func NewSimLinkSvc() LinkSvc {
	return &simLinkSvc{}
}

// Links return  all exiting links
func (svc *simLinkSvc) Links() []*QuickLink {
	data, err := db.Get(quicLinksKey)
	if err != nil {
		logrus.Errorf("failed to load quick links: %s", err)
		return []*QuickLink{}
	}
	var rs []*QuickLink
	err = json.Unmarshal(data, &rs)
	if err != nil {
		logrus.Errorf("failed to load quick links: %s", err)
		return rs
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
	links := svc.Links()
	links = append(links, link)
	data, err := json.Marshal(links)
	if err != nil {
		return err
	}
	db.Put(quicLinksKey, data)
	return nil
}

func (svc *simLinkSvc) Delete(id string) error {
	links := svc.Links()
	svc.DeleteAll()
	for _, l := range links {
		if l.ID != id {
			svc.Add(l)
		}
	}
	return nil
}

func (svc *simLinkSvc) DeleteAll() error {
	return db.Delete(quicLinksKey)
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
