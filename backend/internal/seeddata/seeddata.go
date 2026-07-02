// Package seeddata is the SINGLE SOURCE OF TRUTH for the sample/seed dataset.
//
// It is intentionally dependency-free (plain data) so that BOTH the seeder
// (cmd/seed) and the end-to-end tests (test/e2e) import the exact same
// competitions/winners inventory — the e2e assertions therefore always match
// what the seeder inserted.
//
// The content here was transcribed from the cloned frontend's previously
// HARD-CODED mock (src/lib/data.ts + the winners page). Nothing here is
// invented business logic — it is SAMPLE data that now lives in Postgres +
// MinIO instead of in the frontend components. Fields the backend schema does
// not model (badge text, accent colour, homepage section) are NOT stored here;
// they are pure presentation and live in the frontend (src/lib/presentation.ts),
// keyed by the same Slug (= slugify(Title)).
package seeddata

// Namespace is the fixed UUIDv5 namespace used by the seeder to derive stable
// IDs from natural keys (competition slug, user email, …). Because the IDs are
// deterministic, re-running the seeder upserts the same rows — it never
// duplicates. DO NOT change this value or a re-seed will create fresh rows.
const Namespace = "b07b0000-5eed-5eed-5eed-000000000001"

// Statuses used by the competition service (see internal/competition/entity).
const (
	StatusLive   = "live"
	StatusClosed = "closed"
	StatusDraft  = "draft"
)

// WinnersArchiveSlug is the single CLOSED competition that owns every seeded
// draw/winner. Keeping winners on one archived competition means they do not
// pollute the live competition grids while still being real draw rows the
// public draw endpoint (and the admin draw list) can return.
const WinnersArchiveSlug = "botb-winners-archive"

// Competition is one sample competition. Money is in integer pence to match the
// competitions.ticket_price_pence column. EndsInHours is turned into ends_at =
// now()+EndsInHours by the seeder; the frontend derives the "ENDS …" badge and
// the homepage section purely from that timestamp (plus price==0 => free comp,
// title contains "INSTANT WINS" => instant-wins section).
type Competition struct {
	Slug             string // = slugify(Title); stable natural key (UNIQUE)
	Title            string
	Prize            string
	Description      string
	TicketPricePence int64
	TicketsTotal     int64
	TicketsSold      int64
	Status           string
	EndsInHours      int    // ends_at = now + this many hours
	Image            string // filename under <assets>/images/comps/
	Featured         bool   // shown in the homepage "featured" hero grid
}

// Winner is one sample winner. The seeder creates, per winner: a user, a ticket
// on the winners-archive competition, an avatar media row (owner_type=user),
// and a DRAWN draw linking user+ticket+prize. WonForPence is the tiny amount
// the ticket cost ("won for 2p"). Image is a filename under images/winners/.
type Winner struct {
	Name        string
	Email       string
	Prize       string
	WonForPence int64
	Image       string // filename under <assets>/images/winners/
	AgoLabel    string // human "revealed" label, e.g. "today"
}

// Competitions is the distinct sample competition set (deduped from the mock's
// repeated cards). Titles are unique so slugify(Title) yields a unique slug.
var Competitions = []Competition{
	// --- featured hero grid ---
	{Slug: "1-2m-home-in-zone-1", Title: "£1.2M HOME IN ZONE 1", Prize: "£1.2M London home in Zone 1", Description: "Win a £1.2M home in London Zone 1, mortgage-free.", TicketPricePence: 125, TicketsTotal: 2700000, TicketsSold: 2497500, Status: StatusLive, EndsInHours: 6, Image: "13344.webp", Featured: true},
	{Slug: "500k-instant-wins", Title: "£500K+ INSTANT WINS", Prize: "£500,000+ in instant wins", Description: "Over £500,000 of instant-win prizes plus a guaranteed end draw.", TicketPricePence: 149, TicketsTotal: 235000, TicketsSold: 64674, Status: StatusLive, EndsInHours: 30, Image: "13548.webp", Featured: true},
	{Slug: "audi-r8-for-21p", Title: "AUDI R8 FOR 21P!", Prize: "Audi R8 supercar", Description: "Win an Audi R8 for as little as 21p a ticket.", TicketPricePence: 21, TicketsTotal: 20000, TicketsSold: 12400, Status: StatusLive, EndsInHours: 54, Image: "13558.webp", Featured: true},
	{Slug: "2-5m-instant-wins", Title: "£2.5M+ INSTANT WINS", Prize: "£2.5M+ in instant wins", Description: "Our biggest instant-win competition — over £2.5M in prizes.", TicketPricePence: 187, TicketsTotal: 1500000, TicketsSold: 1050000, Status: StatusLive, EndsInHours: 78, Image: "13551.webp", Featured: true},
	{Slug: "evoque-mini-for-9p", Title: "EVOQUE + MINI FOR 9P!", Prize: "Range Rover Evoque + Mini", Description: "Two cars, one draw — Evoque plus a Mini, from 9p.", TicketPricePence: 12, TicketsTotal: 20000, TicketsSold: 3600, Status: StatusLive, EndsInHours: 264, Image: "13740.webp", Featured: true},

	// --- ends today ---
	{Slug: "iphone-17-and-1249-prizes", Title: "IPHONE 17 & 1,249+ PRIZES!", Prize: "iPhone 17 + 1,249 prizes", Description: "iPhone 17 headlines 1,249+ guaranteed prizes.", TicketPricePence: 89, TicketsTotal: 100000, TicketsSold: 61000, Status: StatusLive, EndsInHours: 5, Image: "11196.webp"},
	{Slug: "1k-house-tickets", Title: "1K HOUSE TICKETS!", Prize: "1,000 house draw tickets", Description: "Bag 1,000 tickets into our dream-home draws.", TicketPricePence: 32, TicketsTotal: 8499, TicketsSold: 3924, Status: StatusLive, EndsInHours: 4, Image: "13831-v1.webp"},
	{Slug: "rattan-dining-set", Title: "RATTAN DINING SET", Prize: "Rattan garden dining set", Description: "A premium rattan dining set for the garden.", TicketPricePence: 7, TicketsTotal: 9600, TicketsSold: 4512, Status: StatusLive, EndsInHours: 6, Image: "13095.webp"},
	{Slug: "mystery-cash-prize", Title: "MYSTERY CASH PRIZE", Prize: "Mystery cash prize", Description: "A mystery cash sum — could be big.", TicketPricePence: 24, TicketsTotal: 10000, TicketsSold: 5100, Status: StatusLive, EndsInHours: 7, Image: "11121.webp"},
	{Slug: "ninja-autobarista-pro", Title: "NINJA AUTOBARISTA PRO", Prize: "Ninja Autobarista Pro", Description: "The Ninja Autobarista Pro coffee machine.", TicketPricePence: 7, TicketsTotal: 10000, TicketsSold: 3100, Status: StatusLive, EndsInHours: 8, Image: "13763.webp"},
	{Slug: "1000-house-tickets-1p", Title: "1000 HOUSE TICKETS 1P", Prize: "1,000 house draw tickets", Description: "1,000 house-draw tickets for a penny a go.", TicketPricePence: 2, TicketsTotal: 10000, TicketsSold: 8900, Status: StatusLive, EndsInHours: 3, Image: "13831-v1.webp"},

	// --- ends tomorrow ---
	{Slug: "5g-gold-bar", Title: "5G GOLD BAR", Prize: "5g gold bar", Description: "A genuine 5g investment gold bar.", TicketPricePence: 7, TicketsTotal: 10000, TicketsSold: 4300, Status: StatusLive, EndsInHours: 28, Image: "9098-gold.webp"},
	{Slug: "2250-cash", Title: "£2,250 CASH", Prize: "£2,250 tax-free cash", Description: "£2,250 straight into your account.", TicketPricePence: 7, TicketsTotal: 12000, TicketsSold: 5766, Status: StatusLive, EndsInHours: 30, Image: "10287-2250.webp"},
	{Slug: "iphone-17-pro-max", Title: "IPHONE 17 PRO MAX", Prize: "iPhone 17 Pro Max", Description: "The latest iPhone 17 Pro Max.", TicketPricePence: 37, TicketsTotal: 10000, TicketsSold: 2950, Status: StatusLive, EndsInHours: 32, Image: "9275-iphone.webp"},
	{Slug: "500-house-tickets", Title: "500 HOUSE TICKETS!", Prize: "500 house draw tickets", Description: "500 tickets into the dream-home draws.", TicketPricePence: 62, TicketsTotal: 1899, TicketsSold: 772, Status: StatusLive, EndsInHours: 34, Image: "13831-v1.webp"},

	// --- instant wins ---
	{Slug: "500k-midweek-instant-wins", Title: "£500K+ MIDWEEK INSTANT WINS", Prize: "£500,000+ midweek instant wins", Description: "A midweek instant-win blowout — over £500K in prizes.", TicketPricePence: 99, TicketsTotal: 400000, TicketsSold: 79500, Status: StatusLive, EndsInHours: 96, Image: "13734.webp"},

	// --- ends soon ---
	{Slug: "samsung-galaxy-book6", Title: "SAMSUNG GALAXY BOOK6", Prize: "Samsung Galaxy Book6", Description: "The Samsung Galaxy Book6 laptop.", TicketPricePence: 5, TicketsTotal: 10000, TicketsSold: 1000, Status: StatusLive, EndsInHours: 60, Image: "13767.webp"},
	{Slug: "macbook-neo", Title: "MACBOOK NEO", Prize: "MacBook Neo", Description: "A brand-new MacBook Neo.", TicketPricePence: 5, TicketsTotal: 10000, TicketsSold: 2850, Status: StatusLive, EndsInHours: 62, Image: "11668.webp"},
	{Slug: "toshiba-tv", Title: "TOSHIBA TV", Prize: "Toshiba Fire TV", Description: "A big-screen Toshiba Fire TV.", TicketPricePence: 2, TicketsTotal: 10000, TicketsSold: 3950, Status: StatusLive, EndsInHours: 64, Image: "10106-toshiba.webp"},
	{Slug: "nintendo-bundle", Title: "NINTENDO BUNDLE", Prize: "Nintendo console bundle", Description: "A Nintendo console plus games bundle.", TicketPricePence: 37, TicketsTotal: 10000, TicketsSold: 1950, Status: StatusLive, EndsInHours: 84, Image: "8941-nintendo.webp"},
	{Slug: "mystery-lifestyle", Title: "MYSTERY LIFESTYLE", Prize: "Mystery lifestyle prize", Description: "A mystery lifestyle prize bundle.", TicketPricePence: 25, TicketsTotal: 10000, TicketsSold: 2950, Status: StatusLive, EndsInHours: 86, Image: "11402.webp"},
	{Slug: "1250-cash", Title: "£1,250 CASH", Prize: "£1,250 tax-free cash", Description: "£1,250 in tax-free cash.", TicketPricePence: 2, TicketsTotal: 10000, TicketsSold: 3550, Status: StatusLive, EndsInHours: 88, Image: "9211-money.webp"},
	{Slug: "5000-cash", Title: "£5,000 CASH", Prize: "£5,000 tax-free cash", Description: "£5,000 straight to your bank.", TicketPricePence: 5, TicketsTotal: 10000, TicketsSold: 1550, Status: StatusLive, EndsInHours: 108, Image: "8941-1000cash.webp"},
	{Slug: "1k-amazon-voucher", Title: "£1K AMAZON VOUCHER", Prize: "£1,000 Amazon voucher", Description: "£1,000 to spend on Amazon.", TicketPricePence: 13, TicketsTotal: 10000, TicketsSold: 800, Status: StatusLive, EndsInHours: 110, Image: "12931.webp"},
	{Slug: "puremate-ac", Title: "PUREMATE AC", Prize: "PureMate air conditioner", Description: "A PureMate portable air-conditioning unit.", TicketPricePence: 7, TicketsTotal: 10000, TicketsSold: 1900, Status: StatusLive, EndsInHours: 112, Image: "13932.webp"},
	{Slug: "ultimate-botb-pass", Title: "ULTIMATE BOTB PASS", Prize: "Ultimate BOTB pass", Description: "The ultimate pass — entries across every draw.", TicketPricePence: 13, TicketsTotal: 10000, TicketsSold: 1400, Status: StatusLive, EndsInHours: 132, Image: "11364.webp"},
	{Slug: "mystery-tech-prize", Title: "MYSTERY TECH PRIZE", Prize: "Mystery tech prize", Description: "A mystery tech prize.", TicketPricePence: 19, TicketsTotal: 10000, TicketsSold: 700, Status: StatusLive, EndsInHours: 134, Image: "10873.webp"},
	{Slug: "1750-cash", Title: "£1,750 CASH", Prize: "£1,750 tax-free cash", Description: "£1,750 tax-free.", TicketPricePence: 11, TicketsTotal: 10000, TicketsSold: 1050, Status: StatusLive, EndsInHours: 136, Image: "10168.webp"},
	{Slug: "dream-car", Title: "DREAM CAR", Prize: "Dream car of your choice", Description: "Win the dream car of your choice, plus gold.", TicketPricePence: 113, TicketsTotal: 50000, TicketsSold: 25000, Status: StatusLive, EndsInHours: 138, Image: "13732.webp", Featured: true},
	{Slug: "10g-gold-bar", Title: "10G GOLD BAR", Prize: "10g gold bar", Description: "A 10g investment gold bar.", TicketPricePence: 2, TicketsTotal: 10000, TicketsSold: 200, Status: StatusLive, EndsInHours: 156, Image: "11858.webp"},
	{Slug: "shark-fan-bundle", Title: "SHARK FAN BUNDLE", Prize: "Shark cooling fan bundle", Description: "A Shark cooling fan bundle.", TicketPricePence: 4, TicketsTotal: 10000, TicketsSold: 50, Status: StatusLive, EndsInHours: 158, Image: "13936.webp"},
	{Slug: "800-cash", Title: "£800 CASH", Prize: "£800 tax-free cash", Description: "£800 in tax-free cash.", TicketPricePence: 2, TicketsTotal: 10000, TicketsSold: 300, Status: StatusLive, EndsInHours: 160, Image: "10415.webp"},
	{Slug: "lifestyle-competition", Title: "LIFESTYLE COMPETITION", Prize: "Lifestyle prize bundle", Description: "A lifestyle prize bundle.", TicketPricePence: 44, TicketsTotal: 10000, TicketsSold: 500, Status: StatusLive, EndsInHours: 162, Image: "13982.webp"},

	// --- free comps (ticket price 0) ---
	{Slug: "free-world-cup-tickets", Title: "FREE WORLD CUP TICKETS!", Prize: "World Cup tickets", Description: "Win World Cup tickets for free.", TicketPricePence: 0, TicketsTotal: 50000, TicketsSold: 22000, Status: StatusLive, EndsInHours: 6, Image: "13883.webp"},
	{Slug: "free-250-cash", Title: "FREE £250 CASH", Prize: "£250 cash", Description: "Win £250 cash for free.", TicketPricePence: 0, TicketsTotal: 30000, TicketsSold: 12000, Status: StatusLive, EndsInHours: 54, Image: "9211-money.webp"},
	{Slug: "free-mystery-prize", Title: "FREE MYSTERY PRIZE", Prize: "Mystery prize", Description: "A free mystery prize.", TicketPricePence: 0, TicketsTotal: 30000, TicketsSold: 9000, Status: StatusLive, EndsInHours: 78, Image: "10873.webp"},
}

// Winners is the sample winners feed (from the mock `winners` list + the
// winners page `extraWinners`). Each becomes a DRAWN draw on the winners
// archive competition.
var Winners = []Winner{
	{Name: "Mark H.", Email: "mark.h@botb-winners.example", Prize: "£500 Cash", WonForPence: 2, Image: "kfi-1.webp", AgoLabel: "7 hours ago"},
	{Name: "Daniel J.", Email: "daniel.j@botb-winners.example", Prize: "Oura Ring 5 + Case", WonForPence: 2, Image: "kfi-2.webp", AgoLabel: "today"},
	{Name: "Barrie M.", Email: "barrie.m@botb-winners.example", Prize: "Toshiba Fire TV", WonForPence: 1, Image: "kfi-3.webp", AgoLabel: "today"},
	{Name: "Jaisa O.", Email: "jaisa.o@botb-winners.example", Prize: "750 House Tickets", WonForPence: 5, Image: "kfi-4.webp", AgoLabel: "today"},
	{Name: "Steven P.", Email: "steven.p@botb-winners.example", Prize: "£1,000 Cash", WonForPence: 3, Image: "kfi-5.webp", AgoLabel: "today"},
	{Name: "Laura B.", Email: "laura.b@botb-winners.example", Prize: "iPhone 17 Pro", WonForPence: 4, Image: "kfi-6.webp", AgoLabel: "yesterday"},
	{Name: "Connor M.", Email: "connor.m@botb-winners.example", Prize: "£2,250 Cash", WonForPence: 2, Image: "kfi-7.webp", AgoLabel: "yesterday"},
	{Name: "Priya S.", Email: "priya.s@botb-winners.example", Prize: "MacBook Neo", WonForPence: 5, Image: "kfi-8.webp", AgoLabel: "yesterday"},
	{Name: "Gary T.", Email: "gary.t@botb-winners.example", Prize: "Nintendo Bundle", WonForPence: 3, Image: "kfi-9.webp", AgoLabel: "yesterday"},
	{Name: "Emma W.", Email: "emma.w@botb-winners.example", Prize: "£800 Cash", WonForPence: 1, Image: "kfi-10.webp", AgoLabel: "yesterday"},
	{Name: "Tom R.", Email: "tom.r@botb-winners.example", Prize: "10g Gold Bar", WonForPence: 2, Image: "kfi-11.webp", AgoLabel: "2 days ago"},
	{Name: "Sofia L.", Email: "sofia.l@botb-winners.example", Prize: "Mystery Tech Prize", WonForPence: 4, Image: "kfi-12.webp", AgoLabel: "2 days ago"},
	{Name: "Aisha K.", Email: "aisha.k@botb-winners.example", Prize: "Rolex Submariner", WonForPence: 5, Image: "kfi-1.webp", AgoLabel: "2 days ago"},
	{Name: "Liam D.", Email: "liam.d@botb-winners.example", Prize: "£5,000 Cash", WonForPence: 3, Image: "kfi-3.webp", AgoLabel: "3 days ago"},
	{Name: "Chloe F.", Email: "chloe.f@botb-winners.example", Prize: "PlayStation 6 Bundle", WonForPence: 2, Image: "kfi-5.webp", AgoLabel: "3 days ago"},
	{Name: "Raj P.", Email: "raj.p@botb-winners.example", Prize: "Audi RS3 Carbon Black", WonForPence: 6, Image: "kfi-7.webp", AgoLabel: "4 days ago"},
}
