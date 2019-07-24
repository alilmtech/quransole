package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/alilmtech/quranc"

	"go.etcd.io/bbolt"
)

func main() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var client quranc.QuranAPI = quranc.New()

	db, err := bbolt.Open("/tmp/quransole.db", os.ModePerm, bbolt.DefaultOptions)
	if err == nil {
		boltedClient, err := quranc.BoltCache(client, db)
		if err == nil {
			client = boltedClient
		}
	}
	defer db.Close()

	chapters, err := client.Chapters(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	chosenSurah := chapters[rand.Intn(len(chapters))]

	chapterVerses, err := client.Verses(context.Background(), chosenSurah.ID, quranc.VersesLimit(1))
	if err != nil {
		log.Printf("surah=%d err=%q", chosenSurah.ChapterNumber, err)
		return
	}

	chosenAyahNumber := rand.Intn(chosenSurah.VersesCount)
	chosenAyah, err := client.Verse(context.Background(), chosenSurah.ID, chapterVerses[0].ID+chosenAyahNumber)
	if err != nil {
		log.Printf("surah=%d ayah=%d err=%q", chosenSurah.ChapterNumber, chosenAyahNumber, err)
		return
	}

	fmt.Printf("Surah %d: %s - %s (%s)\n", chosenSurah.ChapterNumber, chosenSurah.NameArabic, chosenSurah.NameComplex, chosenSurah.TranslatedName.Name)
	fmt.Println("Ayah:", chosenAyah.VerseNumber, " Page:", chosenAyah.PageNumber)
	fmt.Println(chosenAyah.TextMadani)
	var words []string
	for _, word := range chosenAyah.Words {
		if word.Translation.LanguageName != "english" {
			continue
		}
		words = append(words, word.Translation.Text)
	}
	fmt.Println(strings.Join(words, " "))
}
