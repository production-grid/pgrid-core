package notifications

import (
	"bytes"
	"os"
	"text/template"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

//TemplateType defines the types notification templates supported
type TemplateType string

// Enumerates the template types
const (
	TemplateTypeSMS       TemplateType = "sms"
	TemplateTypeSubject   TemplateType = "subject"
	TemplateTypePlainBody TemplateType = "plain-body"
	TemplateTypeHTMLBody  TemplateType = "html-body"
)

// LoadTemplate resolves the correct notification template
func LoadTemplate(templateType TemplateType, def *applications.EventDef, event *applications.Event, rcpt *applications.Recipient) (string, error) {

	//TODO work in a cache system at some point

	var templateResult string
	var err error

	if event.Session.TenantID != "" {
		//custom tenant first
		templateResult, err = loadFile("/" + event.Session.TenantID + "/" + def.Key + "/" + resolveFileName(templateType, true))
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
		//default tenant next
		if templateResult == "" {
			templateResult, err = loadFile("/" + event.Session.TenantID + "/" + def.Key + "/" + resolveFileName(templateType, false))
			if err != nil && !os.IsNotExist(err) {
				return "", err
			}
		}
	}

	//custom non-tenant
	if templateResult == "" {
		templateResult, err = loadFile("/" + def.Key + "/" + resolveFileName(templateType, true))
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
	}

	//default non-tenant
	if templateResult == "" {
		templateResult, err = loadFile("/" + def.Key + "/" + resolveFileName(templateType, true))
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
	}

	//default non-tenant
	if templateResult == "" {
		templateResult, err = loadFile("/" + resolveFileName(templateType, false))
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
	}

	mergedResult, err := Merge(templateResult, def, event, rcpt)

	if err != nil {
		return "", err
	}

	return mergedResult, nil

}

//Merge merges the transient data with the template
func Merge(tpl string, def *applications.EventDef, event *applications.Event, rcpt *applications.Recipient) (string, error) {

	t, err := template.New("").Parse(tpl)
	if err != nil {
		logging.Errorf("Failed to read template %v: %v", tpl, err)
	}

	bytes := &bytes.Buffer{}

	data := make(map[string]interface{})
	data["def"] = *def
	data["event"] = *event
	data["rcpt"] = *rcpt

	err = t.Execute(bytes, data)

	if err != nil {
		return "", err
	}

	return string(bytes.Bytes()), nil

}

func resolveFileName(templateType TemplateType, custom bool) string {

	lastBit := "default"

	if custom {
		lastBit = "custom"
	}

	switch templateType {
	case TemplateTypeSMS:
		return "sms-" + lastBit + ".txt"
	case TemplateTypeSubject:
		return "subject-" + lastBit + ".txt"
	case TemplateTypePlainBody:
		return "body-" + lastBit + ".txt"
	case TemplateTypeHTMLBody:
		return "body-" + lastBit + ".html"
	}

	return ""

}

func loadFile(path string) (string, error) {

	file, err := applications.CurrentApplication.TemplateLoader.String("events" + path)

	if err != nil {
		return "", err
	}

	return file, nil

}
