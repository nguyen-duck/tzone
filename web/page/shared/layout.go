package shared

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func BaseLayout(title string, children g.Node) g.Node {
	return h.Doctype(
		h.HTML(g.Attr("xmlns", "http://www.w3.org/1999/xhtml"), g.Attr("xml:lang", "en-US"), h.Lang("en-US"),
			h.Head(
				h.TitleEl(g.Text(title)),
				h.Meta(h.Charset("utf-8")),
				h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1.0")),

				// Google Fonts
				h.Link(h.Rel("preconnect"), h.Href("https://fonts.googleapis.com")),
				h.Link(h.Rel("preconnect"), h.Href("https://fonts.gstatic.com"), g.Attr("crossorigin", "true")),
				h.Link(h.Rel("stylesheet"), h.Href("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Outfit:wght@400;600;700&display=swap")),

				// Lucide Icons (CDN)
				h.Script(h.Src("https://unpkg.com/lucide@latest")),

				// Global Premium Styles
				h.StyleEl(h.Type("text/css"), g.Raw(`
					:root {
						--primary: #2563eb;
						--primary-hover: #1d4ed8;
						--primary-light: #eff6ff;
						--accent: #d50000;
						--bg-main: #f8fafc;
						--bg-card: #ffffff;
						--text-main: #0f172a;
						--text-muted: #64748b;
						--border: #e2e8f0;
						--radius-sm: 8px;
						--radius-md: 12px;
						--radius-lg: 20px;
						--shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
						--shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
						--shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1);
						--font-sans: 'Inter', system-ui, sans-serif;
						--font-heading: 'Outfit', sans-serif;
					}

					* { margin: 0; padding: 0; box-sizing: border-box; }
					body { 
						font-family: var(--font-sans); 
						background-color: var(--bg-main); 
						color: var(--text-main); 
						line-height: 1.5;
						-webkit-font-smoothing: antialiased;
					}

					.l-container {
						max-width: 1200px;
						margin: 0 auto;
						padding: 0 20px;
					}

					a { text-decoration: none; color: inherit; transition: all 0.2s ease; }

					/* Header Overrides for Premium Look */
					#header {
						background: rgba(255, 255, 255, 0.8);
						backdrop-filter: blur(12px);
						-webkit-backdrop-filter: blur(12px);
						border-bottom: 1px solid var(--border);
						position: sticky;
						top: 0;
						z-index: 1000;
						padding: 12px 0;
					}

					#logo { display: flex; align-items: center; gap: 10px; font-family: var(--font-heading); font-weight: 700; font-size: 22px; }
					#logo img { height: 32px; width: auto; }
					#logo span { color: var(--text-main); }

					.main-menu-list {
						display: flex;
						list-style: none;
						gap: 24px;
						margin-top: 15px;
						font-size: 14px;
						font-weight: 500;
						color: var(--text-muted);
					}
					.main-menu-list a:hover { color: var(--primary); }

					#social-connect { display: flex; gap: 15px; align-items: center; }
					.icon-count { font-size: 11px; color: var(--text-muted); }

					#footer {
						background: var(--bg-card);
						border-top: 1px solid var(--border);
						padding: 60px 0;
						margin-top: 80px;
						text-align: center;
					}
					.footer-logo img { height: 40px; margin-bottom: 20px; opacity: 0.8; }
					#footmenu p { font-size: 13px; color: var(--text-muted); margin-bottom: 10px; }
					#footmenu a { margin: 0 10px; color: var(--text-muted); }
					#footmenu a:hover { color: var(--primary); }

					h1, h2, h3 { font-family: var(--font-heading); font-weight: 700; }
				`)),
				h.Meta(h.Name("description"), h.Content("GSMArena.com - The ultimate resource for GSM handset information")),
				h.Meta(h.Name("keywords"), h.Content("GSM,mobile,phone,Nokia,Sony,Apple,iPhone,Motorola,Alcatel,Xiaomi,Samsung,Oppo,Oneplus,cellphone,specifications,information,info,opinion,review,pictures,photos")),
				h.Link(h.Rel("alternate"), h.Type("application/rss+xml"), h.Title("GSMArena.com RSS feed"), h.Href("https://www.gsmarena.com/rss-news-reviews.php3")),
				h.Link(h.Rel("canonical"), h.Href("https://www.gsmarena.com")),
				h.Link(h.Rel("alternate"), g.Attr("media", "only screen and (max-width: 640px)"), h.Href("https://m.gsmarena.com")),
			),
			h.Body(
				h.Script(h.Type("text/javascript"), h.Src("https://fdn.gsmarena.com/vv/assets12/js/misc.js?v=128")),

				Header(),

				children,

				Footer(),

				h.Script(h.Type("text/javascript"), h.Src("https://fdn.gsmarena.com/vv/assets12/js/autocomplete.js?v=16")),
				h.Script(h.Type("text/javascript"), h.Lang("javascript"), g.Raw(`
AUTOCOMPLETE_LIST_URL = "/quicksearch-81833.jpg";
$gsm.addEventListener(document, "DOMContentLoaded", function() 
{
    new Autocomplete( "topsearch-text", "topsearch", true );
}
)
`)),
				h.Script(h.Type("text/javascript"), h.Src("https://fdn.gsmarena.com/vv/assets12/js/infinite-scroll.js?v=50")),
				h.Script(g.Raw("lucide.createIcons();")),
			),
		),
	)
}
