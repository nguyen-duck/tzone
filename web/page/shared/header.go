package shared

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Header() g.Node {
	return h.Header(h.ID("header"), h.Class("row"),
		h.Div(h.Class("wrapper clearfix"),
			h.Div(h.Class("top-bar clearfix"),
				// HAMBURGER MENU
				h.Button(h.Type("button"), h.Role("button"), h.Aria("label", "Toggle Navigation"), h.Class("lines-button minus"),
					h.Span(h.Class("lines")),
				),

				// LOGO
				h.Div(h.ID("logo"),
					h.A(h.Href("/"),
						h.Object(h.Type("image/svg+xml"), g.Attr("data", "https://fdn.gsmarena.com/vv/assets12/i/logo.svg"),
							h.Img(h.Src("https://fdn.gsmarena.com/vv/assets12/i/logo-fallback.gif"), h.Alt("GSMArena.com")),
						),
						h.Span(g.Text("GSMArena.com")),
					),
				),

				// NAV
				h.Div(h.ID("nav"), h.Role("main"),
					h.Form(h.Action("res.php3"), h.Method("get"), h.ID("topsearch"),
						h.Input(h.Type("text"), h.Placeholder("Search"), h.TabIndex("201"), g.Attr("accesskey", "s"), h.ID("topsearch-text"), h.Name("sSearch"), h.AutoComplete("off")),
						h.Span(h.ID("quick-search-button"),
							h.Input(h.Type("submit"), h.Value("Go")),
							h.I(h.Class("head-icon icomoon-liga icon-search-left")),
						),
					),
				),

				// SOCIAL CONNECT
				h.Div(h.ID("social-connect"),
					h.A(h.Href("tipus.php3"),
						h.I(h.Class("head-icon icon-tip-us icomoon-liga")), h.Br(), h.Span(h.Class("icon-count"), g.Text("Tip us")),
					),
					h.Span(h.Class("bar")),

					h.A(h.Href("https://www.youtube.com/channel/UCbLq9tsbo8peV22VxbDAfXA?sub_confirmation=1"), h.Class("yt-icon"), h.Target("_blank"), h.Rel("noopener"),
						h.I(h.Class("head-icon icon-soc-youtube icomoon-liga")), h.Br(), h.Span(h.Class("icon-count"), g.Text("2.0m")),
					),

					h.A(h.Href("https://www.instagram.com/gsmarenateam/"), h.Target("_blank"), h.Rel("noopener"),
						h.I(h.Class("head-icon icon-instagram icomoon-liga")), h.Span(h.Class("icon-count"), g.Text("150k")),
					),

					h.A(h.Href("rss-news-reviews.php3"),
						h.I(h.Class("head-icon icon-soc-rss2 icomoon-liga")), h.Br(), h.Span(h.Class("icon-count"), g.Text("RSS")),
					),

					h.A(h.Href("https://www.arenaev.com/"), h.Target("_blank"), h.Rel("noopener"),
						h.I(h.Class("head-icon icon-specs-car icomoon-liga")), h.Br(), h.Span(h.Class("icon-count"), g.Text("EV")),
					),

					h.A(h.Href("https://merch.gsmarena.com/"), h.Target("_blank"), h.Rel("noopener"),
						h.I(h.Class("head-icon icon-cart icomoon-liga")), h.Br(), h.Span(h.Class("icon-count"), g.Text("Merch")),
					),

					h.Span(h.Class("bar")),

					h.A(h.Href("#"), g.Attr("onclick", "return false;"), h.Class("login-icon"), h.ID("login-active"),
						h.I(h.Class("head-icon icon-login")), h.Br(), h.Span(h.Class("icon-count"), h.Style("right:4px;"), g.Text("Log in")),
					),

					h.Span(h.Class("tooltip"), h.ID("login-popup2"),
						h.Form(h.Action("login.php3"), h.Method("post"),
							h.Input(h.Type("Hidden"), h.Name("sSource"), h.Value("MA%3D%3D")),
							h.P(g.Text("Login")),
							h.Label(h.For("email")),
							h.Input(h.Type("email"), h.ID("email"), h.Name("sEmail"), h.MaxLength("50"), h.Value(""), h.Required(), h.AutoComplete("false")),

							h.Label(h.For("upass")),
							h.Input(h.Type("password"), h.ID("upass"), h.Name("sPassword"), h.Placeholder("Your password"), h.MaxLength("20"), h.Pattern("\\S{6,}"), h.Required(), h.AutoComplete("false")),

							h.Input(h.Class("button"), h.Type("submit"), h.Value("Log in"), h.ID("nick-submit")),
						),
						h.A(h.Class("forgot"), h.Href("forgot.php3"), g.Text("I forgot my password")),
					),
					h.A(h.Href("/signup"), h.Class("signup-icon no-margin-right"),
						h.I(h.Class("head-icon icon-user-plus")), h.Span(h.Class("icon-count"), g.Text("Sign up")),
					),
				),
			),

			// MENU
			h.Ul(h.ID("menu"), h.Class("main-menu-list"),
				h.Li(h.A(h.Href("/"), g.Text("Home"))),
				h.Li(h.A(h.Href("news.php3"), g.Text("News"))),
				h.Li(h.A(h.Href("reviews.php3"), g.Text("Reviews"))),
				h.Li(h.A(h.Href("videos.php3"), g.Text("Videos"))),
				h.Li(h.A(h.Href("news.php3?sTag=Featured"), g.Text("Featured"))),
				h.Li(h.A(h.Href("search.php3"), g.Text("Phone Finder"))),
				h.Li(h.A(h.Href("deals.php3"), h.Style("position: relative;"), g.Text("Deals"))),
				h.Li(h.A(h.Href("https://merch.gsmarena.com/"), h.Target("_blank"), h.Style("position: relative;"),
					g.Text("Merch"), h.Span(h.Class("icon-count"), h.Style("top: 3px; right: 5px;"), g.Text("New")),
				)),
				h.Li(h.A(h.Href("network-bands.php3"), g.Text("Coverage"))),
				h.Li(h.A(h.Href("contact.php3"), g.Text("Contact"))),
			),
		),
	)
}
