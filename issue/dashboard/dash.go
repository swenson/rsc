// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dashboard implements the issue dashboard for an
// upcoming Go release.
package dashboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"code.google.com/p/rsc/appfs/fs"
	"code.google.com/p/rsc/issue"
)

type Point struct {
	Time  time.Time
	Yes   int
	Maybe int
}

func (p Point) JSDate() template.JS {
	yy, mm, dd := p.Time.Date()
	h, m, s := p.Time.Clock()
	return template.JS(fmt.Sprintf("new Date(%d, %d, %d, %d, %d, %d)", yy, mm-1, dd, h, m, s))
}

func Update(ctxt *fs.Context, client *http.Client, version string) error {
	prefix := strings.Map(func(r rune) rune {
		if 'A' <= r && r <= 'Z' {
			return r - 'A' + 'a'
		}
		if 'a' <= r && r <= 'z' || '0' <= r && r <= '9' {
			return r
		}
		return -1
	}, version)

	label := strings.Map(func(r rune) rune {
		if r == ' ' {
			return -1
		}
		return r
	}, version)

	graphFile := "/issue-dashboard/" + prefix + ".graph"
	htmlFile := "/issue-dashboard/" + label

	var graph []Point
	data, _, err := ctxt.Read(graphFile)
	if err == nil {
		if err := json.Unmarshal(data, &graph); err != nil {
			return fmt.Errorf("unmarshal dashboard graph: %v", err)
		}
	}

	yes, err := issue.Search("go", "open", "label:"+label, false, client)
	if err != nil {
		return fmt.Errorf("searching for %s issues: %v", version, err)
	}
	maybe, err := issue.Search("go", "open", "label:"+label+"Maybe", false, client)
	if err != nil {
		return fmt.Errorf("searching for %sMaybe issues: %v", label, err)
	}

	graph = append(graph, Point{time.Now(), len(yes), len(maybe)})
	data, err = json.Marshal(graph)
	if err != nil {
		return fmt.Errorf("marshal dashboard graph: %v", err)
	}
	if err := ctxt.Write(graphFile, data); err != nil {
		return fmt.Errorf("writing dashboard graph: %v", err)
	}

	byDir := map[string][]*issue.Issue{}
	for _, p := range append(yes, maybe...) {
		dir := p.Summary
		if i := strings.Index(dir, ":"); i >= 0 {
			dir = dir[:i]
		}
		if i := strings.Index(dir, ","); i >= 0 {
			dir = dir[:i]
		}
		byDir[dir] = append(byDir[dir], p)
	}

	var small []Point
	now := time.Now()
	day := -1
	for _, p := range graph {
		if p.Maybe == 0 {
			continue
		}
		d := p.Time.Day()
		if d != day || now.Sub(p.Time) < 3*24*time.Hour {
			day = d
			small = append(small, p)
		}
	}

	tmplData := struct {
		Version string
		Label   string
		Graph   []Point
		Issues  map[string][]*issue.Issue
	}{
		Version: version,
		Label:   label,
		Graph:   small,
		Issues:  byDir,
	}

	var buf bytes.Buffer
	tmpl, err := template.New("main").
		Funcs(template.FuncMap{
		"hasLabel":  hasLabel,
		"hasStatus": hasStatus,
	}).
		Parse(dashTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %v", err)
	}
	if err := tmpl.Execute(&buf, &tmplData); err != nil {
		return fmt.Errorf("executing template: %v", err)
	}
	if err := ctxt.Write(htmlFile, buf.Bytes()); err != nil {
		return fmt.Errorf("writing html: %v", err)
	}
	return nil
}

func hasStatus(p *issue.Issue, status string) bool {
	return p.Status == status
}

func hasLabel(p interface{}, name string) string {
	switch p := p.(type) {
	case *issue.Issue:
		for _, l := range p.Label {
			if l == name {
				return name
			}
			if strings.HasSuffix(name, "-") && strings.HasPrefix(l, name) {
				return l[len(name):]
			}
		}
	case []*issue.Issue:
		if strings.HasPrefix(name, "-") {
			for _, q := range p {
				if hasLabel(q, name) != "" {
					return ""
				}
			}
			return "ok"
		}
		for _, q := range p {
			if s := hasLabel(q, name); s != "" {
				return s
			}
		}
	}
	return ""
}

var dashTemplate = `<html>
  <head>
    <script type="text/javascript" src="https://www.google.com/jsapi"></script>
    <script type="text/javascript">
      google.load("visualization", "1", {packages:["corechart"]});
      google.setOnLoadCallback(drawCharts);
      function drawCharts() {
        var data = new google.visualization.DataTable();
        data.addColumn('datetime', 'Date');
        data.addColumn('number', '{{.Version}}');
        data.addColumn('number', '{{.Version}} + Maybe');
        var one = 1;
        data.addRows([
{{range .Graph}}          [{{.JSDate}}, {{.Yes}}, {{.Yes}}+{{.Maybe}}],
{{end}}
        ])
        var options = {
          width: 800, height: 400,
          title: '{{.Version}} Issues',
          strictFirstColumnType: true,
          vAxis: {minValue: 0, maxValue: 299},
          vAxes: {0: {title: 'Open Issues'}}
        };
        var chart = new google.visualization.AreaChart(document.getElementById('open_div'));
        chart.draw(data, options);
      }
    </script>
    <script type="text/javascript" src="https://ajax.googleapis.com/ajax/libs/jquery/1.8.2/jquery.min.js"></script>
    <script>
      var onlySuggest = false;
      var hideMaybe = true;
      function rehide() {
        $("tr").show();
        if(onlySuggest) {
          $("tr.nosuggest").hide();
        }
        if(hideMaybe) {
          $("tr.maybe").hide();
        }
        
        if(onlySuggest && !hideMaybe) {
          $("#suggest").html("suggested issues only");
        } else {
          $("#suggest").html("<a href='javascript:dosuggest()'>show suggested issues only</a>");
        }
        if(!onlySuggest && hideMaybe) {
          $("#yes").html("{{.Version}} issues only");
        } else {
          $("#yes").html("<a href='javascript:doyes()'>show {{.Version}} issues only</a>");
        }
        if(!onlySuggest && !hideMaybe) {
          $("#maybe").html("all issues");
        } else {
          $("#maybe").html("<a href='javascript:domaybe()'>show all issues</a>");
        }
      }
      function dosuggest() {
        onlySuggest = true;
        hideMaybe = false;
        window.location.hash = "s";
        rehide();
      }
      function doyes() {
        onlySuggest = false;
        hideMaybe = true;
        window.location.hash = "";
        rehide();
      }
      function domaybe() {
        onlySuggest = false;
        hideMaybe = false;
        window.location.hash = "m";
        rehide();
      }
      function start() {
        if (window.location.hash == "s" || window.location.hash == "#s") {
          dosuggest();
        } else if (window.location.hash == "m" || window.location.hash == "#m") {
          domaybe();
        } else {
          doyes();
        }
      }
    </script>
    
    <style>
      td.dir {font-weight: bold;}
      td.suggest {padding-left: 1em;}
      .size {font-family: sans-serif; font-size: 70%; text-align: center;}
      tr.maybe {color: #aaa;}
      tr.suggest {}
      h1 {font-size: 120%;}
      a {color: #000;}
      tr.maybe a {color: #aaa;}
      .key, .key td {font-family: sans-serif; font-size: 90%;}
    </style>
  </head>

  <body onload="start()">
    <h1>{{.Version}}: Open Issues</h1>

    <div id="open_div"></div>
    
    <div class="key">
    Key:
    <table>
    <tr><td class="suggest"><td class="size">S</td><td>small change: less than 30 minutes (e.g. doc fix)
    <tr><td class="suggest"><td class="size">M</td><td>medium change: less than 2 hours (e.g. small feature/fix + tests)
    <tr><td class="suggest"><td class="size">L</td><td>large change: less than 8 hours
    <tr><td class="suggest"><td class="size">XL</td><td>extra large change: more than one day
    <tr><td class="suggest"><td>&#x261e;</td><td>suggested for people looking for work
    </table>
    </div>
    <br><br>
    
    <span id="suggest"></span> | <span id="yes"></span> | <span id="maybe"></span>

    <br><br>
    <table>
    {{range $dir, $list := .Issues}}
      <tr class="{{if hasLabel $list $.Label}}yes{{else}}maybe{{end}} {{if hasLabel $list "Suggested"}}suggest{{else}}nosuggest{{end}}"><td class="dir" colspan="4">{{$dir}}
      {{range $list}}
        <tr class="{{if hasLabel . $.Label}}yes{{else}}maybe{{end}} {{if hasLabel . "Suggested"}}suggest{{else}}nosuggest{{end}}">
          <td class="suggest">{{if hasLabel . "Suggested"}}&#x261e;{{end}}
          <td class="size">{{hasLabel . "Size-"}}
          <td class="num">{{.ID}}
          <td class="title"><a href="http://golang.org/issue/{{.ID}}">{{.Summary}}</a>
            {{if hasLabel . "{{.Version}}Maybe"}}[maybe]{{end}}
            {{if hasStatus . "Started"}}[<i>started by {{.Owner}}</i>]{{end}}
      {{end}}
    {{end}}
    </table>
  </body>
</html>
`
