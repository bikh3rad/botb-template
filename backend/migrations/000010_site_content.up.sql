-- Minimal site-copy key-value store: public read, admin write. Seeded from the
-- values that were previously HARDCODED in the frontend (src/lib/data.ts and
-- src/lib/presentation.ts) so switching the frontend to this table changes
-- nothing visually.
CREATE TABLE IF NOT EXISTS site_content (
    key        TEXT PRIMARY KEY,
    value      TEXT        NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO site_content (key, value) VALUES
('hero.slides', $$[
  {"badge":"APP EXCLUSIVE","title":"WIN FREE TECH IN APP!","subtitle":"Free Apple Airpods, Nintendo Switch 2s & Ninja Creamis","image":"/images/hero/slide-13734.webp"},
  {"badge":"ENDS SOON","title":"A LUXURY £1.2M HOME IN ZONE 1","subtitle":"Last 7 days to win this London home for just £1","image":"/images/hero/slide-13344.webp"},
  {"badge":"DREAM CAR COMPETITION","title":"THE GOLDEN BOOT!","subtitle":"6 cars with gold, who will win the coveted Golden Boot!","image":"/images/hero/slide-13728.webp"},
  {"badge":"SUPERCAR","title":"WIN A DEFENDER D350 X!","subtitle":"The ultimate defender for only 20p!","image":"/images/hero/slide-13736.webp"},
  {"badge":"PRIZE EVERY TIME","title":"SUMMER FESTIVAL INSTANT WINS!","subtitle":"A prize every time to kick-start your summer!","image":"/images/hero/slide-13551.webp"},
  {"badge":"OPEN FOR ENTRIES!","title":"LIFESTYLE COMPETITION","subtitle":"Win cars, cash, tech, watches, and so much more!","image":"/images/hero/slide-13982.webp"},
  {"badge":"THE CLOCK'S TICKING","title":"WIMBLEDON WONDERS INSTANT WINS!","subtitle":"Game, set, match! Epic prizes for £1.19!","image":"/images/hero/slide-13548.webp"},
  {"badge":"YOU'VE STILL GOT TIME","title":"WIN AN AUDI RS3 CARBON BLACK","subtitle":"Drive away in the ultimate hot hatch for 6p!","image":"/images/hero/slide-13589.webp"}
]$$),
('hero.stats', $$[
  {"value":"26 Years","label":"UK's No.1","est":true},
  {"value":"Over 721k+","label":"Winners"},
  {"value":"£160M+","label":"in Prizes Won"}
]$$),
('winners.count', '9,700 winners'),
('trustband', $${"est":"EST. 1999 — £160M+ IN PRIZES","guarantee":"Guaranteed winners every week","dreamCar":"Dream Car »","firstVisit":"First Visit? »"}$$),
('justlaunched', $${"title":"Just Launched!","subtitle":"Win cars, bikes, tech or cash!","cta":"Enter Now"}$$),
('footer.description', $$["Winvia Entertainment PLC (formerly known as Best of the Best Limited) operates skilled prize competitions resulting in the allocation of prizes in accordance with the Terms and Conditions of the website.","Win a brand-new car, take the cash alternative, or win Competition Credit in the BOTB Dream Car Competition. There are over 150 new car prizes to choose from, and the closest person in the skilled Spot the Ball game wins the car or a life-changing amount of cash!","And don't miss out on our Instant Win, Lifestyle, Luxury car, House and Free Competitions to win cars, motorbikes, holidays, watches, tech, cash, and more life-changing prizes. Simply choose your tickets and check out."]$$)
ON CONFLICT (key) DO NOTHING;
