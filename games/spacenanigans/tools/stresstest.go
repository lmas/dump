package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gobwas/ws"
	"github.com/lmas/spacenanigans/wsclient"
)

const (
	urlAuth      = "/auth/login"
	urlWebsocket = "/game/ws"
)

var (
	fHost = flag.String("host", "http://localhost:8000", "addr to game host")
	fUser = flag.String("user", "tester", "login as user")
	fPass = flag.String("pass", "tester", "login with password")
	fBots = flag.Int("bots", 10, "number of stress bots to run")
	fSeed = flag.Int64("seed", 0, "seed to use for the RNG (default current time)")

	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func main() {
	flag.Parse()
	if *fSeed == 0 {
		*fSeed = time.Now().Unix()
	}
	rand.Seed(*fSeed)
	//var headers http.Header
	headers := http.Header{}
	cookie, err := doAuth(*fHost, *fUser, *fPass)
	if err != nil {
		panic(err)
	}
	headers.Set("Cookie", cookie.String())

	for i := *fBots; i > 0; i-- {
		go func() {
			err := NewClient(*fHost, headers)
			if err != nil {
				panic(err)
			}
		}()
	}
	select {} // block forever
}

////////////////////////////////////////////////////////////////////////////////

func doAuth(host, user, pass string) (*http.Cookie, error) {
	cli := &http.Client{
		Timeout: 60 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects
			return http.ErrUseLastResponse
		},
	}

	u, err := url.Parse(host + urlAuth)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	data := url.Values{
		"user": []string{user},
		"pass": []string{pass},
	}
	resp, err := cli.Post(host+urlAuth, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 303 {
		return nil, fmt.Errorf("bad response status: %s", resp.Status)
	}
	return resp.Cookies()[0], nil
}

////////////////////////////////////////////////////////////////////////////////

var dialer = &ws.Dialer{
	Timeout:         60 * time.Second,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewClient(host string, headers http.Header) error {
	u, err := url.Parse(host + urlWebsocket)
	if err != nil {
		return err
	}
	if u.Scheme != "ws" {
		u.Scheme = "ws"
	}

	dialer.Header = ws.HandshakeHeaderHTTP(headers)
	conn, _, _, err := dialer.Dial(context.Background(), u.String())
	if err != nil {
		return fmt.Errorf("error connecting: %s", err)
	}
	c := wsclient.New(conn, logger)
	defer c.Close()
	go func() {
		for {
			_, err := c.Read()
			if err != nil {
				logger.Printf("read error: %s\n", err)
				break
			}
		}
	}()

	time.Sleep(1 * time.Second)
	for {
		switch rand.Intn(3) {
		case 0:
			time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

		case 1:
			time.Sleep(2 * time.Second)
			if rand.Intn(100) < 20 {
				var msg []string
				for m := rand.Intn(5); m > 0; m-- {
					msg = append(msg, Messages[rand.Intn(len(Messages))])
				}
				c.Write(wsclient.PacketMessage, strings.Join(msg, ". "))
			}

		case 2:
			for i := rand.Intn(10); i > 0; i-- {
				dir := rand.Intn(4)
				c.Write(wsclient.PacketMove, dir)
				sleep := rand.Intn(10) + 1
				time.Sleep(time.Duration(sleep*100) * time.Millisecond)
				c.Write(wsclient.PacketStopMove, nil)
			}
		}
	}
	return nil
}

var Messages = []string{
	"hello world! xD",
	"hmm....",
	"fuuuq i almost blowd meself up o.0'",
	"yeah this is a verry long and messy message that I gotta send through this test",
	"[Insert your commit message here. Be sure to make it descriptive.]",
	"Derp. Fix missing constant post rename",
	"It's secret!",
	"This Is Why We Don't Push To Production On Fridays",
	"bug fix",
	"Apparently works-for-me is a crappy excuse.",
	"Either Hot Shit or Total Bollocks",
	"Spinning up the hamster...",
	"Not sure why",
	"various changes",
	"Copy-paste to fix previous copy-paste",
	"This is my code. My code is amazing.",
	"COMMIT ALL THE FILES!",
	"One little whitespace gets its very own commit! Oh, life is so erratic!",
	"Another commit to keep my CAN streak going. ",
	"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	"/sigh ",
	"Same as last commit with changes",
	"Reticulating Splines",
	"Gathering Goblins",
	"Lifting Weights",
	"Pushing Pixels",
	"Formulating Plan",
	"Taking Break",
	"Herding Ducks",
	"Feeding Developers",
	"Fishing for Change",
	"Searching for Dancers",
	"Waking Up Gnomes",
	"Playing Chess",
	"Building Igloos",
	"Converting Celsius",
	"Scanning Power Level",
	"Delivering Presents",
	"Finding Dragon Balls",
	"Firing Lasers",
	"Party Rocking",
	"Walking up to the club",
	"Righting wrongs",
	"Building Lego",
	"Assembling Avengers",
	"Turning Down for What",
	"Reaching 88mph",
	"Pondering Existence",
	"Battling Robots",
	"Smashing Pots",
	"Stomping Goombas",
	"Doing Donuts",
	"Entering Danger Zone",
	"Talking to Mom",
	"Chasing Squirrels",
	"Setting Phasers to Stun",
	"Doing Macarena",
	"Dropping Bass",
	"Removing Biebers",
	"Performing Magic",
	"Autotuning Kanye",
	"Waxing Legs",
	"Invading Space",
	"Levelling Up",
	"Generating Map",
	"Conquering France",
	"Piloting Tardis",
	"Destroying Deathstar",
	"Typing Letters",
	"Making Code",
	"Running Marathon",
	"Shooting Pucks",
	"Kicking Field Goals",
	"Fighting Bad Guys",
	"Driving Batmobile",
	"Warming Up Kryptonite",
	"Popping Popcorn",
	"Creating Hashes",
	"Spawning Boss",
	"Evaluating Life Choices",
	"Eating Ramen",
	"Re-heating Leftovers",
	"Petting Kittens",
	"Walking Puppies",
	"Catching Z’s",
	"Jumping Rope",
	"Declaring Variables",
	"Yessing Doge",
	"Recycling Memes",
	"Tipping Fedora",
	"Walking Runway",
	"Counting to Ten",
	"Booting Native Client",
	"Launching App",
	"Drawing Icons",
	"Reading Instructions",
	"Finding Screws",
	"Completing Puzzles",
	"Generating Volume Slider",
	"Brightening Orange",
	"Ordering Pizza",
	"You Look Good Today",
	"Clearing Screen",
	"Stirring Pot",
	"Mashing Potatoes",
	"Banishing Evil",
	"Taking Selfies",
	"Accelerating Disks",
	"Benching Network",
	"Rocking Out",
	"Grinding Mage",
	"Studying Calculus",
	"Playing N64",
	"Racing GoKarts",
	"Defeating Creepers",
	"Blowing Game Cartridge",
	"Choosing Pikachu",
	"Postponing Half Life 3",
	"Rushing Zergs",
	"Rescuing Hostages",
	"Typing Konami Code",
	"Building Snowman",
	"Letting it Snow",
	"Burning HDMI Cords",
	"Applying Filters",
	"Taking Screenshot",
	"Shaving Mustache",
	"Growing Beard",
	"Baking Muffins",
	"Iterating Javascript",
	"Attracting Venture Capital",
	"Disrupting Industry",
	"Tweeting Hashtags",
	"Encrypting Lines",
	"Obfuscating C",
	"Enhancing License Plate",
	"Running Diagnostic",
	"Warming Hyperdrive",
	"Calibrating Positions",
	"Calculating Percentages",
	"Revoking Licenses",
	"Shedding Core",
	"Dampening Gravity",
	"Increasing Power",
	"Checking Sensors",
	"Indexing RSS",
	"Programming PCI",
	"Determining USB Position",
	"Connecting to Bus",
	"Inverting Ports",
	"Bypassing Capacitor",
	"Reversing Bandwidth Throttle",
	"Testing AI",
	"Virtualizing Microchip",
	"Emulating Playstation",
	"Synthesizing Drivers",
	"Structuring Chlorophyll",
	"Watering Plants",
	"Ingesting Caffeine",
	"Chugging Redbull",
	"Parsing System",
	"Navigating Arrays",
	"Searching Google",
	"Overflowing Stack",
	"Compiling Binaries",
	"Answering Emails",
	"Migrating CSS",
	"Backing Up Primaries",
	"Rendering Dialogs",
	"Reading RSS",
	"Compressing Data",
	"Rejecting Cloud",
	"Evaluating Weissman Score",
	"Purging Local Storage",
	"Leaking Memory",
	"Scripting Python",
	"Grunting Ruby",
	"Benching RAM",
	"Determining Auxiliaries",
	"Jiggling Internet",
	"Ejecting Floppy",
	"Fluctuating Objects",
	"Spiking Reactor Core",
	"Firing Bosons",
	"Testing Processor",
	"Debugging Prompts",
	"Connecting Floats",
	"Rounding Integers",
	"Pronouncing Gigawatt",
	"Inverting Transponders",
	"Bypassing Silicon",
	"Raising Funds",
	"Caching Logs",
	"Dithering Broadband",
	"Eating Poutine",
	"Rolling Rims to Win",
	"Begging for Change",
	"Chasing Waterfalls",
	"Pumping Gas",
	"Emptying Pipes",
	"Hitting Piñata",
	"Unleashing Freedom",
	"Airbrushing Actors",
	"FIling Taxes",
	"Powering Mitochondria",
	"Calculating Qi charge",
	"Completing Geometry",
	"Turning in Algebra",
	"Solving for X",
	"Benching Wattage",
	"Kludging Playback Bar",
	"Stringifying Json",
	"Consuming Spaghetti Code",
	"Deleting Comments",
	"Transitioning to Django",
	"Learning to Code",
	"Battling Feature Creep",
	"Losing Flappy Bird",
	"Celebrating Good Times",
	"Sharpening Pencils",
	"Automating Processes",
	"Attacking Godzilla",
	"Carbonating Soda",
	"Thinking of Witty Text",
}
