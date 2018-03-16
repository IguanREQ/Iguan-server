package subscriber

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"iguan/event"
	"iguan/logs"
)

// TODO: move it to config
const CliRoot = "/var/opt/iguan/handlers"

const (
	DestTypeCli  DestType = iota + 1
	DestTypeHttp
)

var (
	ErrAlreadyExists = errors.New("Subscriber already exists")
	ErrInvalidParams = errors.New("Invalid params")
)

type DestPath = string
type DestType = uint8

type SubjectNotifyInfo struct {
	DestType DestType `json:"destType"`
	DestPath DestPath `json:"destPath"`
}

func (s *SubjectNotifyInfo) Valid() bool {
	switch s.DestType {
	case DestTypeCli:
		if _, err := os.Stat(s.DestPath); os.IsNotExist(err) {
			// path/to/whatever does not exist
		}
	case DestTypeHttp:
		u, err := url.ParseRequestURI(s.DestPath)
		if err == nil && u.IsAbs() {
			return true
		}
	}
	return false
}

type Subscriber struct {
	destType   DestType
	destPath   DestPath
	sourceTag  string
	eventMasks []string
	createdAt  time.Time

	domains []*DomainSubscribers
}

func newSubscriber() *Subscriber {
	return &Subscriber{
		destType:   0,
		destPath:   "",
		sourceTag:  "",
		eventMasks: nil,
		createdAt:  time.Now(),
		domains:    nil,
	}
}

func (s *Subscriber) fire(e *event.Event) error {
	switch s.destType {
	case DestTypeCli:
		logs.Error("Not implemented")
	case DestTypeHttp:
		logs.Info("Fire event by HTTP")
		buf, err := json.Marshal(e)
		if err != nil {
			logs.Error("Json marshal error: %s", err)
			return err
		}
		_, err = http.DefaultClient.Post(s.destPath, "application/json", bytes.NewReader(buf))
		if err != nil {
			logs.Error("HTTP post error: %s", err)
			return err
		}
		logs.Info("Event fired successfully")
	default:
		return ErrInvalidParams
	}
	return nil
}

type DomainSubscribers struct {
	childs        map[string]*DomainSubscribers
	subExact      []*Subscriber
	subChildLeafs []*Subscriber
	subAllChilds  []*Subscriber
}

var subscribers map[DestPath]*Subscriber
var rootToken *DomainSubscribers
var mu sync.RWMutex

func init() {
	rootToken = NewDomain()
}

func NewDomain() *DomainSubscribers {
	return &DomainSubscribers{
		childs:        make(map[string]*DomainSubscribers),
		subExact:      nil,
		subChildLeafs: nil,
		subAllChilds:  nil,
	}
}

func (et *DomainSubscribers) getChild(token string) *DomainSubscribers {
	if v, ok := et.childs[token]; ok {
		return v
	} else {
		return nil
	}
}

func (et *DomainSubscribers) setSubscriber(h *[]*Subscriber, sub *Subscriber) {
	for _, s := range *h {
		if s.destPath == sub.destPath {
			return
		}
	}
	*h = append(*h, sub)

	sub.domains = append(sub.domains, et)
}

func (et *DomainSubscribers) addSubscriber(mask *string, sub *Subscriber) {
	res := strings.SplitN(*mask, ".", 2)
	if len(res) == 2 {
		*mask = res[1]
	} else {
		*mask = ""
	}
	part := res[0]
	switch part {
	case "#":
		et.setSubscriber(&et.subAllChilds, sub)
	case "*":
		et.setSubscriber(&et.subChildLeafs, sub)
	default:
		ch := et.getChild(part)
		if ch == nil {
			ch = NewDomain()
			et.childs[part] = ch
		}
		if *mask == "" {
			ch.setSubscriber(&ch.subExact, sub)
		} else {
			ch.addSubscriber(mask, sub)
		}
	}
}

func Fire(e *event.Event) error {
	mu.RLock()
	defer mu.RUnlock()

	sub := rootToken

	eventDomains := strings.Split(e.Body.Name, ".")
	domainsCount := len(eventDomains)
	for i, group := range eventDomains {
		sub = sub.getChild(group)
		if sub == nil {
			return nil
		}
		// TODO: firing by goroutines
		switch i {
		case domainsCount - 2:
			for _, s := range sub.subChildLeafs {
				s.fire(e)
			}
		case domainsCount - 1:
			for _, s := range sub.subExact {
				s.fire(e)
			}
		}

		for _, s := range sub.subAllChilds {
			s.fire(e)
		}
	}

	return nil
}

func Register(sourceTag string, eventMask string, si *SubjectNotifyInfo) error {
	mu.Lock()
	defer mu.Unlock()

	if !si.Valid() {
		return ErrInvalidParams
	}

	sub, ok := subscribers[si.DestPath]
	if ok {
		sub.eventMasks = append(sub.eventMasks, eventMask)
	} else {
		sub = newSubscriber()
		sub.destPath = si.DestPath
		sub.destType = si.DestType
		sub.eventMasks = []string{eventMask}
		sub.sourceTag = sourceTag
	}

	for _, mask := range sub.eventMasks {
		m := mask
		rootToken.addSubscriber(&m, sub)
	}
	return nil
}

/*
func Unregister(s *SubjectNotifyInfo) error {
	return nil
}

func UnregisterAll() error {
	return nil
}
*/
