package links

var quicLinksKey = []byte("quickLinks")
var svc LinkSvc

type QuickLink struct {
	ID       string
	Name     string
	Link     string
	Type     string
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
	LinksByType(typ string) []*QuickLink
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

func LinksByType(typ string) []*QuickLink {
	return svc.LinksByType(typ)
}

func Delete(id string) error {
	return svc.Delete(id)
}
