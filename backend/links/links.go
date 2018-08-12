package links

var quicLinksKey = []byte("quickLinks")
var svc LinkSvc

type QuickLink struct {
	Name     string
	Link     string
	Visitted int64
}

func (ql *QuickLink) FullName() string {
	if ql.Name == "" {
		return ql.Link
	}
	return ql.Name
}

func SetLinkSvc(s LinkSvc) {
	svc = s
}

type LinkSvc interface {
	Links() []*QuickLink
	Add(link *QuickLink) error
	Delete(id string) error
	DeleteAll() error
}

// Links return  all exiting links
func Links() []*QuickLink {
	return svc.Links()
}

//AddQuickLinks add the new links to existing list
func Add(link *QuickLink) error {
	return svc.Add(link)
}

func DeleteAll() error {
	return svc.DeleteAll()
}
