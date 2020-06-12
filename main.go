package WumpagotchiAIO

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Wumpus struct {
	Credits		int
	Name		string
	Color		int
	Age			int
	Health    	int
	Hunger  	int
	Energy    	int
	Happiness 	int
	Sick      	bool
	Sleeping  	bool
	Left      	bool
}

func init() {
	flag.StringVar(&DiscordToken, "t", "", "Discord API Token")
	flag.Parse()
}

var DiscordToken string
var CommandPrefix = "w."
var wump *discordgo.Session
var err error
var db *sql.DB

func main() {
	if DiscordToken == "" {
		fmt.Println("==WumpagotchiAIO Error==\nYour start command should be as follows:\nWumpagotchiAIO -t <Discord API Token>")
		os.Exit(0)
	}

	db, _ = sql.Open("sqlite3", "./wump.sqlite")
	wump, err = discordgo.New("Bot " + DiscordToken)
	if err != nil {
		fmt.Println("==WumpagotchiAIO Error==\n" + err.Error())
		os.Exit(0)
	}

	wump.AddHandler(loginLogic)

	err = wump.Open()
	if err != nil {
		fmt.Println()
	}

	go timespell()
	go agespell()

	fmt.Println("WumpagotchiAIO Online\nRunning until a termination signal is recieved ...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	wump.Close()
}

func loginLogic(session *discordgo.Session, event *discordgo.Ready) {
	for true {
		query := datastore.NewQuery("User")
		var wumpus []Wumpus
		var err error
		_, err = gcp.GetAll(ctx, query, &wumpus)
		if err != nil {
			fmt.Println(err.Error())
		}
		session.UpdateStatus(0, "with "+strconv.Itoa(len(wumpus))+" Wumpi")
		time.Sleep(60 * time.Second)
	}
}

func sendMessage(session *discordgo.Session, event *discordgo.MessageCreate, channel string, message string) {
	sentMessage, err := session.ChannelMessageSend(channel, message)
	if err != nil {
		fmt.Println(err)
	}
	err = session.ChannelMessageDelete(event.ChannelID, event.ID)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(10 * time.Second)
	session.ChannelMessageDelete(sentMessage.ChannelID, sentMessage.ID)
}

func sendEmbed(session *discordgo.Session, event *discordgo.MessageCreate, channel string, embed *discordgo.MessageEmbed) {
	sentMessage, err := session.ChannelMessageSendEmbed(channel, embed)
	if err != nil {
		fmt.Println(err)
	}
	err = session.ChannelMessageDelete(event.ChannelID, event.ID)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(15 * time.Second)
	session.ChannelMessageDelete(sentMessage.ChannelID, sentMessage.ID)
}

func timespell() {
	for range time.Tick(time.Hour * 2) {
		updateList()
		fmt.Println("Updated Wumpi")
		for w := 0; w+1 <= len(wumpus); w++ {
			rand.Seed(time.Now().UnixNano())
			//First Check if the Wumpus left so we don't run anything else if they have
			if wumpus[w].Left == true {
				continue
			}
			if wumpus[w].Energy > 8 && wumpus[w].Happiness > 8 && wumpus[w].Health > 8 && wumpus[w].Hunger > 8 && wumpus[w].Sick == false && wumpus[w].Age > 1 {
				wumpus[w].Credits += 10
			}
			//Have a 10% chance of making the Wumpus sick
			if rand.Float32() <= 0.10 {
				wumpus[w].Sick = true
			}
			//Next check if they are sick and if they are reduce all stats by 1
			if wumpus[w].Sick == true {
				wumpus[w].Health--
				wumpus[w].Energy--
				wumpus[w].Hunger--
				wumpus[w].Happiness--
			}
			//Next have a 25% chance of reducing happiness
			if rand.Float32() <= 0.25 {
				wumpus[w].Happiness--
			}
			//Have a 50% chance of reducing hunger by 1
			if rand.Float32() <= 0.50 {
				wumpus[w].Hunger--
			}
			//Check if they are sleeping if they are add 4 to energy and 1 to health
			//then if have a 75% chance of reducing energy by 1
			//If they are not sleeping and their energy is at or below 0 mark them as sleeping
			if wumpus[w].Sleeping == true {
				wumpus[w].Energy += 4
			} else if rand.Float32() <= 0.75 {
				wumpus[w].Energy--
			}
			if wumpus[w].Energy <= 0 {
				wumpus[w].Sleeping = true
			}
			//Check if hunger is @ or below 0 if so reduce health and happiness by 1
			if wumpus[w].Hunger <= 0 {
				wumpus[w].Health--
				wumpus[w].Happiness--
			}

			//Make sure that all values are set to 0 before checking for happiness and writing to GCP Datastore
			if wumpus[w].Health < 0 {
				wumpus[w].Health = 0
			}
			if wumpus[w].Hunger < 0 {
				wumpus[w].Hunger = 0
			}
			if wumpus[w].Energy < 0 {
				wumpus[w].Energy = 0
			}
			if wumpus[w].Happiness < 0 {
				wumpus[w].Happiness = 0
			}
			if wumpus[w].Health > 10 {
				wumpus[w].Health = 10
			}
			if wumpus[w].Hunger > 10 {
				wumpus[w].Hunger = 10
			}
			if wumpus[w].Energy > 10 {
				wumpus[w].Energy = 10
			}
			if wumpus[w].Happiness > 10 {
				wumpus[w].Happiness = 10
			}
			//Check if the happiness is @ or below 0 and then have a 50% chance if both are true mark the wumpus as Left
			if wumpus[w].Happiness <= 0 && rand.Float32() < 0.50 {
				wumpus[w].Left = true
			}

			if wumpus[w].Health == 0 {
				wumpus[w].Left = true
			}

			userKey := datastore.NameKey("User", keys[w].Name, nil)
			if _, err := gcp.Put(ctx, userKey, &wumpus[w]); err != nil {
				fmt.Println("==Warning==\nFailed to update Wumpus in Datastore")
				break
			}
		}
	}
}

func agespell() {
	for range time.Tick(time.Hour * 24) {
		updateList()
		fmt.Println("Aged Wumpi")
		for w := 0; w+1 <= len(wumpus); w++ {
			if wumpus[w].Left == true {
				continue
			}
			wumpus[w].Age++
			if wumpus[w].Age >= 14 {
				wumpus[w].Age = 14
				wumpus[w].Left = true
			}
			userKey := datastore.NameKey("User", keys[w].Name, nil)
			if _, err := gcp.Put(ctx, userKey, &wumpus[w]); err != nil {
				fmt.Println("==Warning==\nFailed to update Wumpus in Datastore")
				break
			}
		}
	}
}

func updateList() {
	query := datastore.NewQuery("User")
	wumpus = nil
	keys, err = gcp.GetAll(ctx, query, &wumpus)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err.Error())
	}
}